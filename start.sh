#!/bin/bash

CHAINID="${CHAIN_ID:-aircaiks_9000-1}"
MONIKER="localtestnet"
# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# Set dedicated home directory for the evmosd instance
HOMEDIR="$HOME/.evmosd"
# to trace evm
#TRACE="--trace"
TRACE=""

# feemarket params basefee
BASEFEE=1000000000

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json
	VAL_KEY="mykey"

./build/station-evm start \
	--metrics "$TRACE" \
	--log_level $LOGLEVEL \
	--json-rpc.api eth,txpool,personal,net,debug,web3 \
	--home "$HOMEDIR" \
	--chain-id "$CHAINID"



##	# Allocate genesis accounts (cosmos formatted addresses)
##	./build/station-evm add-genesis-account "$(./build/station-evm keys show "$VAL_KEY" -a --keyring-backend "$KEYRING" --home "$HOMEDIR")" 100000000000000000000000000aevmos --keyring-backend "$KEYRING" --home "$HOMEDIR"
##
#
#	# Sign genesis transaction
#	./build/station-evm gentx "$VAL_KEY" 1000000000000000000000aevmos --gas-prices ${BASEFEE}aevmos --keyring-backend "$KEYRING" --chain-id "$CHAINID" --home "$HOMEDIR"
#
#	# Collect genesis tx
#	./build/station-evm collect-gentxs --home "$HOMEDIR"
#
#	# Run this to ensure everything worked and that the genesis file is setup correctly
#	./build/station-evm validate-genesis --home "$HOMEDIR"