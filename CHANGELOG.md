<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [v0.5.1] - 2024-16-07

This is the Q2 release of fhEVM-native

Ethermint-node docker image: ghcr.io/zama-ai/ethermint-node:v0.5.1

Ethermint-node developer docker image: ghcr.io/zama-ai/ethermint-dev-node:v0.5.1

## Notes

- Get rid of protected storage
- Fully support of async decryption and reencryption
- Tested with KMS centralized v0.7.1

For build:

|      Name       |    Type    |    version    |
| :-------------: | :--------: | :-----------: |
| zbc-go-ethereum | repository |    v0.5.0     |
|  zbc-ethermint  | repository |    v0.5.1     |
|    fhevm-go     | repository |    v0.5.0     |
|     tfhe-rs     | repository |    v0.6.3     |
|    kms-core     |   docker   |    v0.7.1     |
|    gateway      |   docker   |    v0.7.1     |
|   connector     |   docker   |    v0.7.1     |

For e2e test:

| Name  |    Type    | version  |
| :---: | :--------: | :------: |
| fhevm | repository | v0.5.4-0 |


## [v0.4.0] - 2024-16-04

This version is the same as v0.3.0-8 but has an updated README and tutorial for metrics. This is the Q1 release

## [v0.3.0-8] - 2024-16-04

Ethermint-node docker image: ghcr.io/zama-ai/ethermint-node:v0.4.2

Ethermint-node developer docker image: ghcr.io/zama-ai/ethermint-dev-node:v0.3.0-8

## Notes

- New features added like 64 bits support, rotr, eaddress
- Async decryption (only) is supported
- Sync decryption and reencryption are deprecated

For build:

|      Name       |    Type    |    version    |
| :-------------: | :--------: | :-----------: |
| zbc-go-ethereum | repository |    v0.1.8     |
|  zbc-ethermint  | repository |   v0.3.0-8    |
|    fhevm-go     | repository |   v0.2.0-9    |
|     tfhe-rs     | repository |     0.5.2     |
|    kms-sync     |   docker   |   099f79a19   |
|    kms-async    |   docker   | v0.3.3-pre-13 |
|     oracle      |   docker   |    v0.1.1     |
|   aggregator    |   docker   |    6068bbf    |

For e2e test:

| Name  |    Type    | version  |
| :---: | :--------: | :------: |
| fhevm | repository | v0.4.0-8 |



## [v0.2.3] - 2024-18-02

Ethermint-node docker image: ghcr.io/zama-ai/ethermint-node:v0.2.3

Ethermint-node developer docker image: ghcr.io/zama-ai/ethermint-dev-node:v0.2.3

## Notes

- Allow solidity version `pragma solidity ^0.8.20;`

For build:

|      Name       |    Type    |   version    |
| :-------------: | :--------: | :----------: |
| zbc-go-ethereum | repository |    v0.1.5    |
|  zbc-ethermint  | repository |    v0.2.3    |
|    fhevm-go     | repository |    v0.1.5    |
|     tfhe-rs     | repository |    0.4.1     |
|       kms       |   docker   |    v0.1.2    |


For e2e test:

|      Name      |    Type    |  version  |
| :------------: | :--------: |  :-----:  |
| fhevm-solidity | repository |  v0.3.0   |
| fhevm-tfhe-cli | repository |  v0.2.3   |


## [v0.2.2] - 2024-10-02

Ethermint-node docker image: ghcr.io/zama-ai/ethermint-node:v0.2.2

Ethermint-node developer docker image: ghcr.io/zama-ai/ethermint-dev-node:v0.2.2

## Notes

- Fix random bit generation
- Add tutorial to run debugger mode
- Use native cmux from tfhe-rs which is very **efficient** and **cheap**

For build:

|      Name       |    Type    |   version    |
| :-------------: | :--------: | :----------: |
| zbc-go-ethereum | repository |    v0.1.5    |
|  zbc-ethermint  | repository |    v0.2.2    |
|    fhevm-go     | repository |    v0.1.5    |
|     tfhe-rs     | repository |    0.4.1     |
|       kms       |   docker   |    v0.1.2    |


For e2e test:

|      Name      |    Type    |  version  |
| :------------: | :--------: |  :-----:  |
| fhevm-solidity | repository |  v0.3.0   |
| fhevm-tfhe-cli | repository |  v0.2.2   |

## [v0.2.0] - 2024-01-02

Ethermint-node docker image: ghcr.io/zama-ai/ethermint-node:v0.2.0

Ethermint-node developer docker image: ghcr.io/zama-ai/ethermint-dev-node:v0.2.0


## Notes

- Bump tfhe-rs version to v0.4.1
- Bump fhevm-tfhe-cli version to v0.2.2
- Three docker-compose exist:
    - docker.compose.local.yml: used to build the ethermint node image and for the ci e2e test
    - docker.compose.single.validator.yml: self-sufficient single validator testnet composed of two services:
        1) the validator running docker image from registry
        2) kms service to allow a self-sufficient single validator testnet
    - docker.compose.validator.yml: only running a validator node, the KMS endpoint relies on ```KMS_ENDPOINT_ADDR``` env variable
- Private key (cks) is/should be stored in the KMS (without Key Management System)
- All tests from fhevm (solidity) are passing **except rand** (will be fixed soon)

## Breaking changes

- Use of kms for
    - reencrypt: example of usage is for balanceOf
    - decrypt: we still need them for governor e2e test and all the unit test of fhevm


For build:

|      Name       |    Type    |   version    |
| :-------------: | :--------: | :----------: |
| zbc-go-ethereum | repository | v0.1.2-alpha |
|  zbc-ethermint  | repository |    v0.2.0    |
|    fhevm-go     | repository | v0.1.2-alpha |
|     tfhe-rs     | repository |    0.4.1     |
|       kms       |   docker   |    v0.1.2    |


For e2e test:

|      Name      |    Type    |                 version                  |
| :------------: | :--------: | :--------------------------------------: |
| fhevm-solidity | repository | d7e0e96468356f910678151a54cbe0784f2a7ff2 |
| fhevm-tfhe-cli | repository |                  v0.2.2                  |

## [v0.1.1] - 2023-12-21

Ethermint-node docker image: ghcr.io/zama-ai/ethermint-node:v0.1.1

Ethermint-node developer docker image: ghcr.io/zama-ai/ethermint-dev-node:v0.1.1


## Notes

- First version with a fully integrated fhevm-go
- Use the last version of go-ethereum, i.e. `v1.13.5`
- Private key (cks) is still stored in the validator (without Key Management System)
- All tests from fhevm (solidity) are passing


For build:

|      Name       |    Type    | version |
| :-------------: | :--------: | :-----: |
| zbc-go-ethereum | repository | v0.1.1  |
|  zbc-ethermint  | repository | v0.1.1  |
|    fhevm-go     | repository | v0.1.1  |
|    tfhe-rs      | repository | 0.3.1   |


For e2e test:

|      Name      |    Type    |                 version                  |
| :------------: | :--------: | :--------------------------------------: |
| fhevm-solidity | repository | d7e0e96468356f910678151a54cbe0784f2a7ff2 |
| fhevm-tfhe-cli | repository |                  v0.2.1                  |


