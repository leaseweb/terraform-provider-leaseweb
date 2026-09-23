# Terraform Leaseweb Provider

[![Go Reference](https://pkg.go.dev/badge/github.com/leaseweb/terraform-provider-leaseweb.svg)](https://pkg.go.dev/github.com/leaseweb/terraform-provider-leaseweb)

A Terraform provider to manage Leaseweb resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.9
- [Go](https://golang.org/doc/install) >= 1.23
- [Node](https://nodejs.org) >= 22 (LTS)
- [pnPM](https://pnpm.io/) >= 8.15.9

All requirements are also satisfied by the included [docker-compose.yml](docker-compose.yml).

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

To run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Run mock server

To test against API specifications, a mock server can be run

```shell
docker-compose up -d
```

Make sure to use `localhost` in your requests as caddy does not know what to do with `127.0.0.1`

```shell

curl -i http://localhost:8080/publicCloud/v1/instances --header 'x-lsw-auth: tralala'
```

## First steps

To install relevant git hooks run

```bash
pnpm i
```

### Commits

All commits must adhere to the [conventional commit spec](https://www.conventionalcommits.org/en/v1.0.0/).
For the acceptance tests to run properly make sure to run ```docker-compose up -d```
before committing anything or the commit will fail.

### API Stability

All underlying APIs are now stable and production ready. The Terraform plugin is no longer in beta.

## Architecture

Code architecture is discussed in [ARCHITECTURE.md](ARCHITECTURE.md).
