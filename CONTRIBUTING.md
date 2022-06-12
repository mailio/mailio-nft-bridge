# Contributing

Contributions to `mailio-nft` are welcome from anyone.

 I'm ([@igorrendulic](https://www.github.com/igorrendulic)) developing this from my own interest. As a consequence, I reserve discretionary veto rights on feature additions, architectural decisions, etc., 

## Contribution lifecycle

Before spending lots of time working on an issue, consider asking me for feedback via [Discord](https://discord.gg/uzVbJA46E3), the [issue tracker](https://github.com/mailio/mailio-nft-bridge/issues). I would love to help make your contributions more successful!

Pull requests (internal or external) will be reviewed by me. 

We include a PR template to serve as a **guideline**, not a **rule**, to communicate the code culture I wish to maintain in this repository.

## Style

When in doubt, defer to the [Effective Go](https://go.dev/doc/effective_go) document [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments) as "style guides".

# Prerequisites (MacOS instructions)

### Install golang

```
brew install golang
```

### Install [golangci-lint](https://golangci-lint.run):

```
brew install golangci-lint
brew upgrade golangci-lint
```

### Make sure GOPATH is set:

```
echo $GOPATH
```

You should see `$HOME/go`.

# Building and Testing

To build:

```shell
go build ./...
```

To run tests:

```shell
go test ./...
```

# Viewing Godocs website

```shell
godoc --http :6060
```

and navigate to http://localhost:6060/pkg/github.com/statechannels/go-nitro/

# Pre PR checks:

Please execute the following on any branch that's ready for review or merge.

### format:

```shell
gofmt -w .
```

### lint:

```shell
golangci-lint run
```

### remove unused dependencies:

```shell
go mod tidy
```

# Debugging Tests

VS code is used to debug tests. To start a debugging session in VS code:

- Ensure you have the [go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) installed
- Open the test file.
- Open the `Run and Debug` section.
- Run the `Debug Test` configuration.

With the extension it is also possible to start a debugging session right from a test function.