# EVM station

## Introduction

The Ethereum Virtual Machine (EVM)-based Cosmos Chain represents a groundbreaking platform in the blockchain sphere, specifically designed for the creation and deployment of decentralized applications (DApps) and smart contracts. This platform is a harmonious blend of the scalability provided by the Cosmos network and the adaptability and widespread acceptance of the Ethereum Virtual Machine (EVM). It is engineered to cater to blockchain developers who are seeking a robust and flexible environment for their innovations.

At its core, this platform offers a unique proposition: it combines the high-performance infrastructure of Cosmos, known for its inter-blockchain communication and scalability, with the powerful and versatile programming capabilities of the EVM. This integration allows developers to build sophisticated and efficient DApps and smart contracts that can leverage the best of both worlds.


## Prerequisites
The project requires:

- [Go](https://golang.org/dl/) (Version 1.22 or later)
- [jq](https://stedolan.github.io/jq/download/): A lightweight and flexible command-line JSON processor.

## Getting Started
- To begin using this project, firstly clone this repository to your local machine. 
```shell
    git clone https://github.com/airchains-network/evm-station
    cd evm-station;
    go mod tidy;
```

## EVM node setup
- To setup the evm node, execute the following command:

⚠️ **Warning:**
This command will delete old data at `~/.evmstation` and `./build` directories. Also it will `delete keys of DA and Junction`. So make sure those wallets don't have balance, or export the keys before running this command.
```shell
    /bin/bash ./scripts/local-setup.sh
```

### Init sequencer
For Testing: `Mock DA`
```shell
  HOMEDIR=$HOME/.evmstationd
  build/bin/evmstationd tracks init --home "$HOMEDIR" --daRpc "mock-rpc" --daKey "mockKey" --daType "mock" --junctionRpc "https://junction-testnet-rpc.synergynodes.com:443" --junctionKeyName j-key
```
Alternative: `For Eigen DA`
```shell 
  HOMEDIR=$HOME/.evmstationd
build/bin/evmstationd tracks init --home "$HOMEDIR" --daRpc "disperser-holesky.eigenda.xyz" --daKey "9430d5ad8ea52329be63afe66a8c8d5e0ba75bf0de0cbd41aa30fadf5f575ec24cff557777e20a0578ec4fedc66274c37fe5d25ed4c4a09cb73b1ddc15349bb4" --daType "eigen" --junctionRpc "http://0.0.0.0:26657" --junctionKeyName j-key
```

### Get Details of Sequencer, Balance of Junction, and DA
```shell
  build/bin/evmstationd tracks details
  build/bin/evmstationd tracks balance junction
  build/bin/evmstationd tracks balance da
```
- Fund tracks account 

### Create Station
To create a station in Junction, run the following command:
```bash
  build/bin/evmstationd tracks create-station --info "some info"
```
By default, the `--track` parameter uses the address created during sequencer initialization in above steps.
Alternatively, if you want to specify a different or multiple track address, use the following command format:
```bash
  build/bin/evmstationd tracks create-station --info "some info" --tracks ["<track_address-1>","<track_address-2>"]
```

### Start the node
To start the node, run the following command:
```bash
  /bin/bash ./scripts/local-start.sh
```

### Start the Sequencer 
To start the tracks, run the following command:
```bash
  build/bin/evmstationd tracks start
```
### Testing 
To test the node, run the following command:
```bash
  go run ./cmd/evmstationd/main.go tracks start
```
### Setup Track


Make sure to replace `<track_address>` with the appropriate address, including your own address created in the previous steps.
**Note:** Ensure you use the correct track address, including yours created in the previous steps.

## Contributing
Contributions are greatly appreciated. You can make contributions by creating issues, fixing bugs, or suggesting new features. Feel free to fork this repository and create pull requests to affect changes.

## License
This project is licensed under the MIT license - see the [LICENSE](LICENSE) file for more information.

## Contact
For any inquiries or constructive feedback, please contact this email contact@airchains.io