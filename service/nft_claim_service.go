package service

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"github.com/jinzhu/copier"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/util"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
)

type NftClaimService struct {
	environment *model.Environment
}

func NewNftClaimService(environment *model.Environment) *NftClaimService {
	return &NftClaimService{
		environment: environment,
	}
}

/**
 * Returns balance in WEI
 */
func (ecs *NftClaimService) GetBalance() (big.Int, error) {
	privKey := lc.Conf.BlockchainConfig.MailioNFTBrokerPrivateKey
	privateKey, err := crypto.HexToECDSA(privKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return big.Int{}, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	if err != nil {
		return *big.NewInt(0), err
	}
	balance, err := ecs.environment.EthClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		lc.Log.Error("failed to get balance", err)
		return *big.NewInt(0), err
	}
	return *balance, nil
}

// verify users signature in order to claim an NFT
func verifySignature(fromAddress, signatureHex string, catalogId string) error {

	sd := apitypes.TypedData{}
	copier.Copy(&sd, &model.SignerData) // deep copy object structure (too complex to not many data points required)

	sd.Message["catalogId"] = catalogId
	sd.Message["wallet"] = fromAddress
	sd.Domain.Salt = lc.Conf.BlockchainConfig.EIP712TypedData.Salt
	sd.Domain.Name = lc.Conf.BlockchainConfig.EIP712TypedData.Name
	sd.Domain.VerifyingContract = lc.Conf.BlockchainConfig.MailioNFTContractAddress
	sd.Domain.Version = lc.Conf.BlockchainConfig.EIP712TypedData.Version
	sd.Domain.ChainId = math.NewHexOrDecimal256(int64(lc.Conf.BlockchainConfig.DefaultChainId))

	domainMap := sd.Domain.Map()
	delete(domainMap, "salt") // remove salt from domain map

	typedDataHash, tdhErr := sd.HashStruct(sd.PrimaryType, sd.Message)
	if tdhErr != nil {
		lc.Log.Error("failed to hash struct: primaryType: ", sd.PrimaryType, tdhErr)
		return tdhErr
	}
	domainSeparator, dsErr := sd.HashStruct("EIP712Domain", domainMap)
	if dsErr != nil {
		lc.Log.Error("failed to hash struct: domain: ", domainMap, dsErr)
		return dsErr
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	msg := crypto.Keccak256Hash(rawData)

	signature := hexutil.MustDecode(signatureHex)
	if len(signature) != 65 {
		return errors.New("invalid signature length")
	}

	if signature[64] != 27 && signature[64] != 28 {
		return errors.New("invalid recovery id")
	}
	signature[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	pubKeyRaw, err := crypto.Ecrecover(msg.Bytes(), signature)
	if err != nil {
		lc.Log.Error("failed to recover public key", err)
		return err
	}
	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		lc.Log.Error("failed to unmarshal public key", err)
		return err
	}
	userAddress := common.HexToAddress(fromAddress)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	fmt.Printf("user addr: %v\n", userAddress)
	fmt.Printf("recovered addr: %v\n", recoveredAddr)

	if !bytes.Equal(userAddress.Bytes(), recoveredAddr.Bytes()) {
		return errors.New("invalid signature")
	}
	return nil
}

// actual minting of the new Mailio NFT
// throws ErrSignature if signature is invalid
// throws ErrExists if NFT already claimed by user for this category
func (ecs *NftClaimService) MintForUser(claim *model.Claim, catalog *model.Catalog) (*types.Transaction, *model.Claim, error) {
	// validate signature
	signatureErr := verifySignature(claim.WalletAddress, claim.Signature, catalog.ID)
	if signatureErr != nil {
		lc.Log.Error("failed to verify signature", signatureErr)
		return nil, nil, model.ErrSignature
	}

	// validate keywords
	isKeywordMatch := ecs.CheckKeywordsMatch(claim.Keywords, catalog.Keywords)
	if !isKeywordMatch {
		return nil, nil, model.ErrKeyword
	}
	// upload NFT to IPFS (JSON File)
	erc20JsonFile := model.Erc721Json{
		Name:            catalog.Name,
		Description:     catalog.Description,
		Image:           "ipfs://" + catalog.ImageLink,
		ExternalUrl:     catalog.ContentLink,
		BackgroundColor: "212529",
		Attributes: []*model.Erc721Attribute{
			{
				"display_type": "boost_number",
				"trait_type":   "informed",
				"value":        5,
			},
		},
	}

	if catalog.VideoLink != "" {
		erc20JsonFile.YoutubeUrl = catalog.VideoLink
	}

	erc20Json, jErr := json.Marshal(erc20JsonFile)
	if jErr != nil {
		lc.Log.Error("failed to marshal ERC721 JSON", jErr)
		return nil, nil, jErr
	}
	jsponReader := bytes.NewReader(erc20Json)
	uploadedIpfs, upErr := util.UploadToIPFSAndPin(ecs.environment.IpfsInfuraClient, claim.WalletAddress+"_"+claim.CatalogId+".json", jsponReader)
	if upErr != nil {
		lc.Log.Error("failed to upload ERC721 JSON to IPFS", upErr)
		return nil, nil, upErr
	}
	if len(uploadedIpfs) == 0 {
		return nil, nil, errors.New("failed to upload ERC721 JSON to IPFS")
	}
	// always take the first item (link to file, not the folder)
	tokenURI := "ipfs://" + uploadedIpfs[0].Hash

	// convert catalogId to bytes ([12]byte)
	catalogID, err := xid.FromString(catalog.ID)
	if err != nil {
		lc.Log.Error("failed to parse catalog id", err)
		return nil, nil, err
	}

	// validate if claim already exists for the catalog and users wallet
	_, cErr := ecs.GetClaim(catalog.ID, claim.WalletAddress)
	if cErr == nil {
		// if not not found error or any other errors then claim exists
		return nil, nil, model.ErrExists
	}

	// every transaction is signed with the private key (in our case private key of broker - the peyee of the transactions)
	privateKey, err := crypto.HexToECDSA(lc.Conf.BlockchainConfig.MailioNFTBrokerPrivateKey)
	if err != nil {
		lc.Log.Error("failed to parse private key", err)
		return nil, nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		lc.Log.Error("error casting public key to ECDSA")
		return nil, nil, errors.New("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// we also need to figure out the gas price and the nonce
	nonce, err := ecs.environment.EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		lc.Log.Error("failed to get nonce", err)
		return nil, nil, err
	}

	gasPrice, err := ecs.environment.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		lc.Log.Error("failed to get gas price", err)
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chainID, err := ecs.environment.EthClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, aErr := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if aErr != nil {
		lc.Log.Error("failed to create tx options", aErr)
		return nil, nil, aErr
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units Expected gas cost for our contract is about 170000 units
	auth.GasPrice = gasPrice
	auth.Context = ctx
	auth.From = fromAddress

	to := common.HexToAddress(claim.WalletAddress)
	// get the uri of the NFT from the catalog
	tx, smErr := ecs.environment.NftContract.SafeMint(auth, to, tokenURI, catalogID)
	if smErr != nil {
		lc.Log.Error("failed to call contract method SafeMint: ", smErr)
		return nil, nil, smErr
	}

	// store minted tx to database
	cl := &model.Claim{
		CatalogId:     catalog.ID,
		TxHash:        tx.Hash().Hex(),
		TokenUri:      tokenURI,
		Signature:     claim.Signature,
		WalletAddress: claim.WalletAddress,
		GasPrice:      tx.GasPrice().Uint64(),
		Created:       time.Now().UnixMilli(),
	}
	claimed, claimErr := ecs.PutClaimedNFT(cl)
	if claimErr != nil {
		lc.Log.Error("failed to put claim", claimErr)
		return tx, nil, nil
	}
	lc.Log.Info("Succesfully minted new mailio NFT at transaction", tx.Hash().Hex())
	return tx, claimed, nil
}

// PutClaim only inserts a new claim to the database
func (ecs *NftClaimService) PutClaimedNFT(claim *model.Claim) (*model.Claim, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()

	claim.Created = time.Now().UnixMilli()
	id := claim.WalletAddress + "_" + claim.CatalogId // unique per catalogId

	m, err := util.MarshalToBytes(claim)
	err = ecs.environment.DB.Put(ctx, util.CreateKey(model.ClaimTable, id), m)
	if err != nil {
		lc.Log.Error("failed to create new catalog", err)
		return nil, err
	}
	return claim, nil
}

// returns stored claim if exists, or ErrNotFound error
func (ecs *NftClaimService) GetClaim(catalogId string, walletAddress string) (*model.Claim, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()
	key := util.CreateKey(model.ClaimTable, walletAddress+"_"+catalogId)
	m, err := ecs.environment.DB.Get(ctx, key)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrNotFound
		}
		lc.Log.Error("failed to get claim", err)
		return nil, err
	}
	claimMap, err := util.UnmarshalFromBytes(m)
	if err != nil {
		lc.Log.Error("failed to unmarshal catalog", err)
		return nil, err
	}
	var clm model.Claim
	err = mapstructure.Decode(claimMap, &clm)
	return &clm, err
}

// lists all claims for a given catalogId order descending by created date
func (ecs *NftClaimService) ListClaims(limit int) ([]*model.Claim, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()

	q := query.Query{
		Limit:  limit,
		Prefix: "/" + model.ClaimTable,
		Orders: []query.Order{query.OrderByKeyDescending{}},
	}

	qRes, err := ecs.environment.DB.Query(ctx, q)
	defer qRes.Close()
	if err != nil {
		lc.Log.Error("failed to list all catalogs", err)
		return nil, err
	}

	claims := []*model.Claim{}
	res, err := qRes.Rest()
	for _, r := range res {
		clBytes, err := util.UnmarshalFromBytes(r.Value)
		if err != nil {
			lc.Log.Error("failed to unmarshal catalog", err)
			return nil, err
		}
		var claim model.Claim
		mapstructure.Decode(clBytes, &claim)
		claims = append(claims, &claim)
	}
	return claims, nil
}

// list all claims belonding to single wallet
func (ecs *NftClaimService) ListClaimsByUser(wallet string, limit int) ([]*model.Claim, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()

	q := query.Query{
		Limit:  limit,
		Prefix: "/" + model.ClaimTable,
		Orders: []query.Order{query.OrderByKeyDescending{}},
		Filters: []query.Filter{
			query.FilterKeyPrefix{
				Prefix: "/" + model.ClaimTable + "/" + wallet,
			},
		},
	}

	qRes, err := ecs.environment.DB.Query(ctx, q)
	defer qRes.Close()
	if err != nil {
		lc.Log.Error("failed to list all catalogs", err)
		return nil, err
	}

	claims := []*model.Claim{}
	res, err := qRes.Rest()
	for _, r := range res {
		clBytes, err := util.UnmarshalFromBytes(r.Value)
		if err != nil {
			lc.Log.Error("failed to unmarshal catalog", err)
			return nil, err
		}
		var claim model.Claim
		mapstructure.Decode(clBytes, &claim)
		claims = append(claims, &claim)
	}
	return claims, nil
}

// checks if keywords match word by word (order of the words doesn't matter)
func (ecs *NftClaimService) CheckKeywordsMatch(claimKeywords []model.ClaimKeyword, catalogKeywords string) bool {

	if catalogKeywords == "" {
		return false
	}
	if len(claimKeywords) == 0 {
		return false
	}

	keywordsFromCatalog := strings.Split(catalogKeywords, ",")
	if len(keywordsFromCatalog) != len(claimKeywords) {
		return false
	}

	for _, keyword := range keywordsFromCatalog {
		foundMatch := false
		for _, catalogKeyword := range claimKeywords {
			if catalogKeyword.Word == strings.Trim(keyword, " ") {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			return false
		}
	}

	return true
}

// reads and parses the transaction log
func (ecs *NftClaimService) ReadClaimedTransactionLogs(walletAddress string, limit int) ([]*model.ClaimPreview, error) {
	claims, err := ecs.ListClaimsByUser(walletAddress, limit)
	if err != nil {
		return nil, err
	}
	// return empty array
	if len(claims) == 0 {
		return []*model.ClaimPreview{}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()

	claimPreviews := []*model.ClaimPreview{}
	for _, claim := range claims {
		// in case transaction doesn't exist yet return 404
		receipt, err := ecs.environment.EthClient.TransactionReceipt(ctx, common.HexToHash(claim.TxHash))
		if err != nil {
			lc.Log.Error("failed to get transaction receipt", err)
			return nil, model.ErrNotFound
		}

		// contract ABI
		logs := receipt.Logs

		claimPreview := &model.ClaimPreview{
			Claim:    *claim,
			TxStatus: receipt.Status,
		}

		for _, l := range logs {
			if l.Address.Hex() == lc.Conf.BlockchainConfig.MailioNFTProxyAddress {
				if len(l.Topics) == 4 { // 0 = method (SafeMint), 1 = from, 2 = to, 3 = tokenURI
					// The event log function signature hash
					// tokenID is always the last param (output from function safeMint) topic[3]
					tokenId := l.Topics[3]
					bigIntTokenId := tokenId.Big()
					claimPreview.TokenId = bigIntTokenId.Uint64()
				}
			}
		}
		claimPreviews = append(claimPreviews, claimPreview)
	}
	return claimPreviews, nil
}
