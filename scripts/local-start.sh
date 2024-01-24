#!/bin/bash

#Normal With CometBFT
#./build/bin/evmstationd start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "./.tmp/evmstationd"

#First Implementation of  StationBFT
../build/bin/evmstationd-v1  start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "../.tmp/polard"