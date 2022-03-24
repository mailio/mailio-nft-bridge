# Mailio NFT Bridge

Mailio NFT server is a bridge between Mailio server and the NFT contract on Blockchain. It's sole purpose is to answer these two questions: 

- has the user been registered for over a month?
- does the user engage with mailio platform?
- does the user know the main keywords emphasised in the content?

If the answer to those questions is yes then user is given the requested NFT free of charge. 

# Usage

## Create admin user

You must create a user when the server is stopped and after you create a `conf.yaml` file. While server is running creating an admin is not possible.

Run the command:

```
go run scripts/make_user.go --email test@example.com -password mypass -config conf.yaml
```

# Development

## conf.yaml

To run the program `conf.yaml` file is required.

```yml
version: 1.0
port: 8080
title: "Mailio NFT Server"
description: "Mailio NFT Server for communication with https://mail.io and MailioNFT Smart Contract"
mode: debug # "debug": or "release"
swagger: true # false disables it
auth_token:
  enabled: false
  header: "authkey"
  token: "abc"

jwt_token:
  enabled: true
  secret_key: "abcedf"

datastore_path: "./data"

admins:
  - email: "admin@mail.io"
    password: "123456"
  - email: "someone@mail.io"
    password: "123456"
```

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