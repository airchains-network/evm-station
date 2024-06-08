#!/bin/bash
CHAINID="${CHAIN_ID:-stationevm_1234-1}"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
HOMEDIR="$HOME/.evmosd"
TRACE=""
BASEFEE=1000000000
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json
VAL_KEY="mykey"
# ./build/bin/evmstationd keys export dev0 --keyring-backend test --unsafe --unarmored-hex
./build/bin/evmstationd  keys   unsafe-export-eth-key "$VAL_KEY"  --keyring-backend "$KEYRING"