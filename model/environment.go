package model

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-resty/resty/v2"
	leveldb "github.com/ipfs/go-ds-leveldb"
	nft "github.com/mailio/mailio-nft-server/onchain/mailionft"
)

type Environment struct {
	DB               *leveldb.Datastore
	EthClient        *ethclient.Client
	NftContract      *nft.Mailionft
	IpfsInfuraClient *resty.Client
}
