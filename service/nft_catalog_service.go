package service

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/onchain"
	"github.com/mailio/mailio-nft-server/util"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
)

type NftCatalogService struct {
	environment          *model.Environment
	mailioNftContractAbi *abi.ABI
}

func NewNftCatalog(environment *model.Environment) *NftCatalogService {
	loadedAbi, err := onchain.ReadMailioNftContractAbi()
	if err != nil {
		lc.Log.Error("failed to read embedded MailioNftContractAbi", err)
		panic("failed to read embedded MailioNftContractAbi")
	}
	mailioNftAbi, err := abi.JSON(bytes.NewReader(loadedAbi))
	if err != nil {
		lc.Log.Error("failed to load MailioNftContractAbi", err)
		panic("failed to load MailioNftContractAbi")
	}
	return &NftCatalogService{
		environment:          environment,
		mailioNftContractAbi: &mailioNftAbi,
	}
}

// PutCatalog usperts a new catalog (insert is exists, or update existing)
func (nc *NftCatalogService) PutCatalog(catalog *model.Catalog) (*model.Catalog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()
	id := util.GenerateRandomID()
	if catalog.ID != "" {
		id = catalog.ID
	}
	catalog.ID = id
	catalog.Modified = time.Now().UnixMilli()
	catalog.Created = time.Now().UnixMilli()

	m, err := util.MarshalToBytes(catalog)
	err = nc.environment.DB.Put(ctx, util.CreateKey(model.CatalogTable, id), m)
	if err != nil {
		lc.Log.Error("failed to create new catalog", err)
		return nil, err
	}
	return catalog, nil
}

// GetCatalog returns a catalog from datastore by ID
func (nc *NftCatalogService) GetCatalog(id string) (*model.Catalog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()
	key := util.CreateKey(model.CatalogTable, id)
	m, err := nc.environment.DB.Get(ctx, key)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrNotFound
		}
		lc.Log.Error("failed to get catalog", err)
		return nil, err
	}
	catalogMap, err := util.UnmarshalFromBytes(m)
	if err != nil {
		lc.Log.Error("failed to unmarshal catalog", err)
		return nil, err
	}
	var cat model.Catalog
	err = mapstructure.Decode(catalogMap, &cat)

	// querying smart contract to get number of claimed NFTs
	xidID, _ := xid.FromString(cat.ID)
	var catId [12]byte
	copy(catId[:], xidID.Bytes())
	categoryCountBigInt, cErr := nc.environment.NftContract.CategoryTokenCount(&bind.CallOpts{}, catId)
	if cErr != nil {
		lc.Log.Error("failed to retrieve category count on blockchain", cErr)
		// ignores the error (it will fail within contract if more than 100 claimed)
	}
	tokensClaimed, convErr := strconv.Atoi(strconv.FormatUint(categoryCountBigInt.Uint64(), 10))
	if convErr != nil {
		lc.Log.Error("failed to convert category count to int", convErr)
		// ignores the error (it will fail within contract if more than 100 claimed)
	}
	cat.NftTokensUsed = tokensClaimed

	return &cat, err
}

// ListAllCatalogs retruns all catalogs from datastore
func (nc *NftCatalogService) ListAllCatalogs(limit int) ([]*model.Catalog, error) {

	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()

	q := query.Query{
		Limit:  limit,
		Prefix: "/" + model.CatalogTable,
		Orders: []query.Order{query.OrderByKeyDescending{}},
	}

	qRes, err := nc.environment.DB.Query(ctx, q)
	defer qRes.Close()
	if err != nil {
		lc.Log.Error("failed to list all catalogs", err)
		return nil, err
	}

	catalogs := []*model.Catalog{}
	res, err := qRes.Rest()
	for _, r := range res {
		cat, err := util.UnmarshalFromBytes(r.Value)
		if err != nil {
			lc.Log.Error("failed to unmarshal catalog", err)
			return nil, err
		}
		var catalog model.Catalog
		mapstructure.Decode(cat, &catalog)
		catalogs = append(catalogs, &catalog)
	}
	return catalogs, nil
}
