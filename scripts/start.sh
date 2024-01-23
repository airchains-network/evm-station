#!/bin/bash

#Normal With CometBFT
#./build/bin/polard start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "./.tmp/polard"

#First Implementation of  StationBFT
../build/bin/polard-v1 start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "../.tmp/polard"