#!/bin/bash

set -Eeuo pipefail

# in /config folder

ETHERMINT_NETWORK_KEYS_PATH=/root/.ethermintd/zama/keys/network-fhe-keys

mkdir -p $ETHERMINT_NETWORK_KEYS_PATH
cp /network-fhe-keys/* $ETHERMINT_NETWORK_KEYS_PATH

# init node
./setup.sh

# start the node
TRACE=""
LOGLEVEL="info"

ETHERMINTD="ethermintd"

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
$ETHERMINTD start --pruning=nothing $TRACE --log_level $LOGLEVEL \
        --minimum-gas-prices=0.0001aphoton \
        --json-rpc.gas-cap=50000000 \
        --json-rpc.api eth,txpool,net,web3 \
        --rpc.laddr "tcp://0.0.0.0:26657"
