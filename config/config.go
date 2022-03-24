package config

import (
	cfg "github.com/chryscloud/go-microkit-plugins/config"
	mclog "github.com/chryscloud/go-microkit-plugins/log"
)

// Conf global config
var Conf Config

// Log global wide logging
var Log mclog.Logger

// Config - embedded global config definition
type Config struct {
	cfg.YamlConfig `yaml:",inline"`
	DatastorePath  string `yaml:"datastore_path"`
}

func init() {
	l, err := mclog.NewEntry2ZapLogger("mailio-nft-server")
	if err != nil {
		panic("failed to initialize logging")
	}
	Log = l
}
