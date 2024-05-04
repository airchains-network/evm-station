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
```shell
    /bin/bash ./scripts/local-setup.sh
```
⚠️ **Warning:**
This command will delete old data at `~/.evmstation` and `./build` directories. Also it will `delete keys of DA and Junction`. So make sure those wallets don't have balance, or export the keys before running this command.

### Init sequencer
```shell
  build/bin/evmstationd sequencer init --home "$HOMEDIR" --daRpc "mock-rpc" --daKey "mockKey" --daType "mock" --junctionRpc "http://0.0.0.0:26657" --junctionKeyName j-key
```

### Get Details of Sequencer, Balance of Junction, and DA
```shell
  build/bin/evmstationd sequencer details
  build/bin/evmstationd sequencer balance junction
  build/bin/evmstationd sequencer balance da  # currently not create or implimented
````

## Contributing
Contributions are greatly appreciated. You can make contributions by creating issues, fixing bugs, or suggesting new features. Feel free to fork this repository and create pull requests to affect changes.


## License
This project is licensed under the MIT license - see the [LICENSE](LICENSE) file for more information.

## Contact
For any inquiries or constructive feedback, please contact this email contact@airchains.io
