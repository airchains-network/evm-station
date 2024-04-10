#!/bin/bash

#Normal With CometBFT
#./build/bin/evmstationd start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "./.tmp/evmstationd"
HOMEDIR="./.evmstation"
#First Implementation of  StationBFT
./build/bin/evmstationd  start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR"