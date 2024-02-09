#!/bin/bash

command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }


# Check if ../.tmp directory exists
if [ -d "../.tmp" ]; then
    read -p "The evm station directory already exists. Do you want to continue? (y/n) " -n 1 -r
    echo    # Move to a new line
    if [[ $REPLY =~ ^[Nn]$ ]]
    then
        echo "Running evmstationd start command."
        ../build/bin/evmstationd start --pruning=nothing "" --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "../.tmp/polard"
        exit 0
    fi
fi


rm -rf ../.tmp
rm -rf ../build
cd ../
make clean
make build
cd scripts
declare -a KEYS
KEYS[0]="dev0"
KEYS[1]="dev1"
KEYS[2]="dev2"
echo "KEYS: ${KEYS[@]}"

CHAINID="nooob-69420"
MONIKER="localtestnet"
PersistentPeers="id1@ip1:port1,id2@ip2:port2" # Example format; replace with your actual peers
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# Set dedicated home directory for the ../build/bin/evmstationd instance
HOMEDIR="../.tmp/polard"
# to trace evm
#TRACE="--trace"
TRACE=""

# Path variables
CONFIG_TOML=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json



#EVM Config

#This is a EIP 155 Standard where every chain needs a unique chain ID
EVMCHAINID="69420"

# used to exit on first error (any non-zero exit code)
set -e

	
	../build/bin/evmstationd init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

	# Set client config
	../build/bin/evmstationd config set client keyring-backend $KEYRING --home "$HOMEDIR"
	../build/bin/evmstationd config set client chain-id "$CHAINID" --home "$HOMEDIR"

	# If keys exist they should be deleted
	for KEY in "${KEYS[@]}"; do
		../build/bin/evmstationd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
	done

	# Change parameter token denominations to abera
	jq '.app_state["staking"]["params"]["bond_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["crisis"]["constant_fee"]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["mint"]["params"]["mint_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.consensus["params"]["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"


# Change the Default EVM Config

sed -i "/\[polaris\.polar\.chain\]/!b;n;c chain-id = \"$EVMCHAINID\"" ../.tmp/polard/config/app.toml
## Change exactly  persistent peers in config.toml



	# Allocate genesis accounts (cosmos formatted addresses)
	for KEY in "${KEYS[@]}"; do
		../build/bin/evmstationd genesis add-genesis-account $KEY 100000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"
	done

	# Test Account
	# absurd surge gather author blanket acquire proof struggle runway attract cereal quiz tattoo shed almost sudden survey boring film memory picnic favorite verb tank
	# 0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
	../build/bin/evmstationd genesis add-genesis-account cosmos1yrene6g2zwjttemf0c65fscg8w8c55w58yh8rl 69000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"

	# Sign genesis transaction
	../build/bin/evmstationd genesis gentx ${KEYS[0]} 1000000000000000000000abera --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"
	## In case you want to create multiple validators at genesis
	## 1. Back to `../build/bin/evmstationd keys add` step, init more keys
	## 2. Back to `../build/bin/evmstationd add-genesis-account` step, add balance for those
	## 3. Clone this ~/.../build/bin/evmstationd home directory into some others, let's say `~/.cloned../build/bin/evmstationd`
	## 4. Run `gentx` in each of those folders
	## 5. Copy the `gentx-*` folders under `~/.cloned../build/bin/evmstationd/config/gentx/` folders into the original `~/.../build/bin/evmstationd/config/gentx`

	# Collect genesis tx
	../build/bin/evmstationd genesis collect-gentxs --home "$HOMEDIR"

	# Run this to ensure everything worked and that the genesis file is setup correctly
	../build/bin/evmstationd genesis validate-genesis --home "$HOMEDIR"

	if [[ $1 == "pending" ]]; then
		echo "pending mode is on, please wait for the first block committed."
	fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)m
../build/bin/evmstationd start --pruning=nothing "$TRACE" --log_level $LOGLEVEL --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR"