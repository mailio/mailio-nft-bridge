package main

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-resty/resty/v2"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/mailio/mailio-nft-server/config"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
	nft "github.com/mailio/mailio-nft-server/onchain/mailionft"
)

// setupEnvironment init of datastore
func setupEnvironment(confg *lc.Config) *model.Environment {
	db := setupDatastore(confg)
	ethClient := setupEthClient(confg)
	contract := loadContract(ethClient, confg.BlockchainConfig.MailioNFTProxyAddress)
	ipfsInfuraClient := setupIPFSInfuraClient()
	env := &model.Environment{
		DB:               db,
		EthClient:        ethClient,
		NftContract:      contract,
		IpfsInfuraClient: ipfsInfuraClient,
	}

	return env
}

// load Mailio NFT contract from the blockchain
func loadContract(client *ethclient.Client, address string) *nft.Mailionft {
	contractInstance, err := nft.NewMailionft(common.HexToAddress(address), client)
	if err != nil {
		panic(err)
	}

	return contractInstance
}

// setup ETH client (connection to the blockchain)
func setupEthClient(config *lc.Config) *ethclient.Client {
	cl, err := ethclient.Dial(config.BlockchainConfig.Endpoint)
	if err != nil {
		panic(err)
	}
	return cl
}

// setupDatastore tries to create a folder for levedb database and initializes the leveldb datastore
func setupDatastore(config *lc.Config) *leveldb.Datastore {
	err := os.MkdirAll(lc.Conf.DatastorePath, 0755)
	if err != nil {
		panic(err)
	}
	ds, err := leveldb.NewDatastore(lc.Conf.DatastorePath, &leveldb.Options{})
	if err != nil {
		panic(err)
	}
	return ds
}

func setupIPFSInfuraClient() *resty.Client {
	return resty.New().
		SetHostURL(config.Conf.BlockchainConfig.InfuraIpfsApiEndpoint).
		SetBasicAuth(config.Conf.BlockchainConfig.InfuraKey, config.Conf.BlockchainConfig.InfuraSecret)
}

// tearDownEnvironemnt closes the leveldb datastore
func tearDownEnvironment(env *model.Environment) {
	if env.DB != nil {
		env.DB.Close()
	}
}
