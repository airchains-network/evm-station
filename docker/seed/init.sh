# SPDX-License-Identifier: BUSL-1.1
#
# Copyright (C) 2023, Berachain Foundation. All rights reserved.
# Use of this software is govered by the Business Source License included
# in the LICENSE file of this repository and at www.mariadb.com/bsl11.
#
# ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
# TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
# VERSIONS OF THE LICENSED WORK.
#
# THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
# LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
# LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
#
# TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
# AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
# EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
# TITLE.

KEY="brick"
CHAINID="berachain-666"
MONIKER="brickchain"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
HOMEDIR="data/.polard"
TRACE=""
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

if [ "$(ls -A $HOMEDIR)" ]; then
    echo "$HOMEDIR is not empty"
    evmstationd start --pruning=nothing "$TRACE" --log_level $LOGLEVEL --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR"
else
    echo "$HOMEDIR is empty, creating a new network"
    
    evmstationd init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

    jq '.app_state["staking"]["params"]["bond_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["crisis"]["constant_fee"]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["evm"]["params"]["evm_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["mint"]["params"]["mint_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.consensus["params"]["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS"
    mv "$TMP_GENESIS" "$GENESIS"

    evmstationd config set client keyring-backend $KEYRING --home "$HOMEDIR"
    evmstationd config set client chain-id "$CHAINID" --home "$HOMEDIR"

    evmstationd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"

    evmstationd genesis add-genesis-account $KEY 100000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"

    # evmstationd genesis add-genesis-account cosmos1yrene6g2zwjttemf0c65fscg8w8c55w58yh8rl 100000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"

    evmstationd genesis gentx $KEY 1000000000000000000000abera --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"

    evmstationd genesis collect-gentxs --home "$HOMEDIR"

    evmstationd genesis validate-genesis --home "$HOMEDIR"

    evmstationd start --pruning=nothing "$TRACE" --log_level $LOGLEVEL --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR"
    evmstationd start --pruning=nothing '' --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home data/.evmstationd
fi