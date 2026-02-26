# migration-manager

[![Go Report Card](https://goreportcard.com/badge/github.com/anoideaopen/migration-manager)](https://goreportcard.com/report/github.com/anoideaopen/migration-manager)
[![Go Reference](https://pkg.go.dev/badge/github.com/anoideaopen/migration-manager.svg)](https://pkg.go.dev/github.com/anoideaopen/migration-manager)
![GitHub License](https://img.shields.io/github/license/anoideaopen/migration-manager)

[![Go Verify Build](https://github.com/anoideaopen/migration-manager/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/anoideaopen/migration-manager/actions/workflows/go.yml)
[![Security vulnerability scan](https://github.com/anoideaopen/migration-manager/actions/workflows/vulnerability-scan.yml/badge.svg?branch=main)](https://github.com/anoideaopen/migration-manager/actions/workflows/vulnerability-scan.yml)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/anoideaopen/migration-manager/main)
![GitHub Tag](https://img.shields.io/github/v/tag/anoideaopen/migration-manager)

## TOC

- [migration-manager](#migration-manager)
    - [TOC](#toc)
    - [Description](#description)
    - [Build](#build)
      - [Go](#go)
    - [Configuration yaml file](#configuration-yaml-file)
    - [Run](#run)
      - [Export](#export)
      - [Import](#import)
    - [License](#license)
    - [Links](#links)

## Description

A utility for ensuring the migration of the state from the old to the new fabric #hlf#tool#migration#go#

------
## Build

### Go
```shell
CGO_ENABLED=0 go build -v 
```

------
## Configuration yaml file
```yaml
# snapshot is the path to the directory where the files with the HLF state will be located
snapshot: ./tmp

# The HLF section defines the parameters for connecting clients to HLF,
# for state requests and transaction execution
hlf:
  # config defines the path to the HLF SDK connection description
  config:    "/opt/services/etc/client.yaml"
  # org - user's organization.
  org:       org1
  # user - the name of the user from whom the connection is being made.
  user:      User1
  # channel - the channel where the chaincode is deployed.
  channel:   fiat
  # chaincode - the name of the chaincode.
  chaincode: fiat
  # exectimeout defines the maximum execution time of the Invoke or Query HLF call.
  # If this parameter is not set, the default value of 2 minutes will be applied.
  exectimeout: 2m
```
------

## Run
```shell
./migration-manager import -c migration.yaml
```
or

You can also [redefine](https://github.com/spf13/viper#working-with-environment-variables) the values from the config with env variables with the `MIGRATION_` prefix.
```shell
export MIGRATION_SNAPSHOT="/etc/data/snap" &&
export MIGRATION_HLF_CHANNEL="fiat" &&
export MIGRATION_HLF_CHAINCODE="fiat" &&
./migration-manager export -c migration.yaml
```

### Export
Получить стейт чанками (по 1000 в чанке) без фиксации в HLF 
```shell
./migration-manager export -c migration.yaml -e 1000
```
Получить стейт чанками (по 1000 в чанке) с фиксацией в HLF. Таким способом получить только ключи запрещено.
```shell
./migration-manager export -c migration.yaml  -e 1000 -i
```

### Import
Загрузить стейт в новый HLF
```shell
./migration-manager import -c migration.yaml
```

## License

Apache-2.0

## Links

* [origin](https://github.com/anoideaopen/migration-manager)
