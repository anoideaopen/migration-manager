# AGENTS.md — migration-manager

Go CLI tool for migrating Hyperledger Fabric state between channels/chaincodes. Uses cobra, viper, fork of fabric-sdk-go.

## Build & dev

```shell
CGO_ENABLED=0 go build -v
```

Lint is **not** runnable locally — only via CI (`golangci/golangci-lint-action v9` with `.golangci.yml`). Before commit, run:

```shell
go mod tidy && go fix ./... && go fmt ./... && git diff --exit-code
```

## Test

```shell
go test -count 1 ./...
```

Only `cfg/` package has tests (config parsing). No race flag, no integration tests.

## CI order (must match)

1. check-cyrillic-comments — fails if any `.go` file contains Cyrillic chars
2. validate-go — `go mod tidy` + `go fmt ./...` (+ diff check)
3. golangci-lint
4. go-unit-test

## CLI

```
migration-manager export -c migration.yaml -e 1000   # invoke (commits to HLF)
migration-manager get   -c migration.yaml -e 1000   # query only
migration-manager import -c migration.yaml           # upload snapshot to HLF
```

`-e` range: 100–10000, default 1000. Retries up to 10 times (1s sleep).

## Config

YAML via viper, env prefix `MIGRATION_` (e.g. `MIGRATION_HLF_CHANNEL`). Snapshot dir is cleaned on every `export`/`get` run.

## Key dependency

```shell
# go.mod replace — uses forked SDK
replace github.com/hyperledger/fabric-sdk-go => github.com/anoideaopen/fabric-sdk-go v0.1.1
```

## Release

Tags `v*` trigger Docker multi-arch build (`scientificideas/migration-manager`) + GitHub release (linux/darwin/windows × amd64/arm64).

## Packages

| Path | Role |
|------|------|
| `cmd/` | Cobra commands (root, export, get, import) |
| `cfg/` | Config struct + viper loader, only unit tests |
| `core/` | SDK init (`InitChCli`), proto marshal helpers |

Entrypoint: `main.go` → `cmd.Execute()`.
