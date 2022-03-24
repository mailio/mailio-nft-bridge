package main

import (
	"os"

	leveldb "github.com/ipfs/go-ds-leveldb"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
)

// setupEnvironment init of datastore
func setupEnvironment(confg *lc.Config) *model.Environment {
	db := setupDatastore(confg)
	env := &model.Environment{
		DB: db,
	}

	// create admin users if missing
	return env
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

// tearDownEnvironemnt closes the leveldb datastore
func tearDownEnvironment(env *model.Environment) {
	if env.DB != nil {
		env.DB.Close()
	}
}
