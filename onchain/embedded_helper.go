package onchain

import "embed"

//go:embed abi/MailioNFT.abi
var contractAbi embed.FS

func ReadMailioNftContractAbi() ([]byte, error) {
	return contractAbi.ReadFile("abi/MailioNFT.abi")
}
