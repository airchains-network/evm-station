#!/bin/bash

TRACE="--trace"
ETHERMINTD="ethermintd"

LOGLEVEL="info"
# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
$ETHERMINTD start --pruning=nothing $TRACE --log_level $LOGLEVEL \
        --minimum-gas-prices=0.0001aphoton \
        --json-rpc.gas-cap=50000000 \
        --json-rpc.api eth,txpool,net,web \
	--rpc.laddr tcp://0.0.0.0:26657
