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



## Running the project

- To setup the project, execute the following command:
```shell
    /bin/bash ./scripts/local-setup.sh 
```

```shell
sudo tee /etc/systemd/system/evmosd.service > /dev/null <<EOF
[Unit]
Description=evmosd node
After=network-online.target
[Service]
User=$USER
WorkingDirectory=$HOME/.evmosd
ExecStart=/root/evm-station/build/station-evm start \
--metrics "" \
--log_level "info" \
--json-rpc.api eth,txpool,personal,net,debug,web3 \
--chain-id "stationevm_1234-1"
Restart=on-failure
RestartSec=5
LimitNOFILE=65535
[Install]
WantedBy=multi-user.target
EOF
```

```shell
sudo systemctl daemon-reload
sudo systemctl enable evmosd
sudo systemctl restart evmosd
```

```shell
sudo journalctl -u evmosd -fo cat
```

You can stop it with CTRL-C

```shell
sed -i 's/address = "127.0.0.1:8545"/address = "0.0.0.0:8545"/' /root/.evmosd/config/app.toml
```

- To Get the Private Keys of the EVM Chain, execute the following command:
```shell
    /bin/bash ./scripts/local-keys.sh 
```



## Contributing
Contributions are greatly appreciated. You can make contributions by creating issues, fixing bugs, or suggesting new features. Feel free to fork this repository and create pull requests to affect changes.


## License
This project is licensed under the MIT license - see the [LICENSE](LICENSE) file for more information.

## Contact
For any inquiries or constructive feedback, please contact this email contact@airchains.io
