package model

import leveldb "github.com/ipfs/go-ds-leveldb"

type Environment struct {
	DB *leveldb.Datastore
}
