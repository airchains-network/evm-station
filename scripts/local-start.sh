#!/bin/bash

#Normal With CometBFT
#./build/bin/evmstationd start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "./.tmp/evmstationd"
HOMEDIR="$HOME/.evmstationd"
#First Implementation of StationBFT
./build/bin/evmstationd  start --rpc.laddr tcp://0.0.0.0:26667 --p2p.laddr "tcp://0.0.0.0:26666" --grpc.address "localhost:9900"  --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR"