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
	cfg.YamlConfig   `yaml:",inline"`
	DatastorePath    string               `yaml:"datastore_path"`
	EtherscanConfig  EtherscanSubConfig   `yaml:"etherscan"`
	BlockchainConfig BlockchainSubConfig  `yaml:"blockchain"`
	ReCaptchaV3      ReCaptchaV3SubConfig `yaml:"recaptcha"`
}

type EtherscanSubConfig struct {
	Endpoint                 string `yaml:"endpoint"`
	ApiKey                   string `yaml:"api_key"`
	MailioNftContractAddress string `yaml:"mailio_nft_contract_address"`
}

type BlockchainSubConfig struct {
	DefaultChainId            int                      `yaml:"default_chain_id"`
	MailioNFTProxyAddress     string                   `yaml:"mailio_nft_proxy"`
	MailioNFTContractAddress  string                   `yaml:"mailio_nft_contract"`
	MailioNFTBrokerPrivateKey string                   `yaml:"broker_private_key"`
	Endpoint                  string                   `yaml:"endpoint"`
	InfuraKey                 string                   `yaml:"infura_key"`
	InfuraSecret              string                   `yaml:"infura_secret"`
	InfuraIpfsApiEndpoint     string                   `yaml:"infura_ipfs_api_endpoint"`
	InfuraIpfsGateway         string                   `yaml:"infura_ipfs_gateway"`
	EIP712TypedData           EIP712TypedDataSubConfig `yaml:"eip712_typed_data"`
}

type EIP712TypedDataSubConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Salt    string `yaml:"salt"`
}

type ReCaptchaV3SubConfig struct {
	Secret string `yaml:"secret"`
	Host   string `yaml:"host"`
}

func init() {
	l, err := mclog.NewEntry2ZapLogger("mailio-nft-server")
	if err != nil {
		panic("failed to initialize logging")
	}
	Log = l
}
