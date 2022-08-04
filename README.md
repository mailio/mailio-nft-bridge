# Mailio NFT Bridge

![https://discord.gg/hXjFS2zWra](https://img.shields.io/static/v1?label=discord&message=developers&color=green)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mailio/mailio-nft-bridge)
![GitHub issues](https://img.shields.io/github/issues/mailio/mailio-nft-bridge)

Mailio NFT server is a bridge between Mailio server and the NFT contract on Blockchain. It's sole purpose is to answer this question:

- does the user know the main keywords emphasized in the content?

If the answer to those questions is yes then user is given the requested NFT free of charge.

# Prerequisities

There is couple of developer accounts needed to run the project. All of the services are free up to certain usage point.

## What you need

- Wallet (Bridge wallets private key, funded with some crypto).
- Etherscan api key on your target blockchain/sidechain
- Deployed [mailio-nft-contracts](https://github.com/mailio/mailio-nft-contracts)
- Infura account with IPFS and IPFS gateway set up (for uploading your NFTs)

Next thing you'll need is a `conf.yaml` configuration.

# Usage

## conf.yaml

It start with a configuration file.

````yml
## conf.yaml

To run the program `conf.yaml` file is required.

```yml
version: 1.0
port: 8080
title: "Mailio NFT Server"
description: "Mailio NFT Server"
mode: debug # "debug": or "release"
swagger: true # false disables it
auth_token: # not used
  enabled: false
  header: "authkey"
  token: "abc"

jwt_token:
  enabled: true
  secret_key: "abcedf" # create strong key

# datastore specific config
datastore_path: "./data"

# reCaptchaV3
recaptcha:
  secret: abcdef
  host: https://www.google.com/recaptcha/api/siteverify


# etherscan config
etherscan:
  mailio_nft_contract_address: "0xabc"
  api_key: abc
  endpoint: "https://api-testnet.polygonscan.com/api" # mainnnet: https://api.polygonscan.com/api

# blockchain specific config
blockchain:
  default_chain_id: 137 # 137 in production
  mailio_nft_proxy: "0xabc" # mailio NFT proxy contract address
  mailio_nft_contract: "0xabc" # mailio NFT contract address
  broker_private_key: "abc" # Broker wallet private key
  endpoint: "https://polygon-mumbai.g.alchemy.com/v2/zM-abc" # Access to blockchain node
  infura_key: "abc" # infura key
  infura_secret: "abc" # infura secret
  infura_ipfs_api_endpoint: "https://ipfs.infura.io:5001" # infura api endpoint
  infura_ipfs_gateway: "https://mailio.infura-ipfs.io" # inufura ipfs gateway
  eip712_typed_data: # building data for EIP-712 signature
    name: "Mailio Knowledge NFTs"
    version: "1.0"
    salt: "0xabc" # domain differentiator (for avoiding the same signature in multiple contracts)
````

## Create admin user

You must create a user when the server is stopped and after you create a `conf.yaml` file. While server is running creating an admin is not possible.

Run the command:

```
go run scripts/make_user.go --email test@example.com -password mypass -config conf.yaml
```

# Development

Run development server:

```
go run setup.go main.go --config conf.yaml
```

## Swagger

Run swagger on the code to generate update API docs:

```
swag init
swag fmt
```

Docs available at: `http://localhost:8080/swagger/index.html`

## Main dependencies

- [Swaggo](https://github.com/swaggo/swag)
- [IPFS leveldb datastore](https://github.com/ipfs/go-ds-leveldb)
- [go-microkit-plugins](https://github.com/chryscloud/go-microkit-plugins)
