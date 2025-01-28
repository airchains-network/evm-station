#!/bin/bash

DENOM=aphoton

MONIKER="orchestrator"

if [ -z "$CHAINID" ];
then
    echo "\$CHAINID environment variable is undefined"
    exit 1
fi

# Orchestrator account
KEY="orchestrator"
if [ -z "$MNEMONIC" ];
then
	echo \$MNEMONIC is undefined
	exit 1
fi

LOCALNET_DIR=$(pwd)/localnet
BUILD_DIR=$LOCALNET_DIR/build
STARTING_IP=10.0.0.10

# TODO uncomment this when issue https://github.com/evmos/ethermint/issues/1579 is solved
DATA_DIR=$BUILD_DIR/node4/ethermintd
LOCALNET_DIR_DOCKER=/localnet
BUILD_DIR_DOCKER=$LOCALNET_DIR_DOCKER/build
DATA_DIR_DOCKER=$BUILD_DIR_DOCKER/node4/ethermintd
FHEVM_KEYS_DOCKER=$LOCALNET_DIR_DOCKER/fhevm-keys
CHAIND="docker run -e FHEVM_GO_KEYS_DIR=$FHEVM_KEYS_DOCKER -v $LOCALNET_DIR:$LOCALNET_DIR_DOCKER -i --rm $DOCKER_IMAGE /usr/bin/ethermintd"

CONF_DIR=$DATA_DIR/config
GENESIS=$CONF_DIR/genesis.json
TEMP_GENESIS=$CONF_DIR/tmp_genesis.json
CONFIG=$CONF_DIR/config.toml
# TODO: workout gas supply for genesis
VALIDATOR_COUNT=4
FULL_NODE_COUNT=3
FULL_NODE_FIRST_IDX=$(( VALIDATOR_COUNT + 1 ))
FULL_NODE_LAST_IDX=$(( FULL_NODE_FIRST_IDX + FULL_NODE_COUNT ))

if [ -z "$DOCKER_IMAGE" ];
then
	echo "Environment variable \$DOCKER_IMAGE is undefined"
	exit 1
fi

if [ -d "$BUILD_DIR" ];
then
	echo "localnet directory $LOCALNET_DIR already exists, to regenerate genesis files remove it"
	exit 1
fi

# create necessary directory for orchestrator node
mkdir -p "$DATA_DIR"

$CHAIND testnet init-files --v $VALIDATOR_COUNT -o $BUILD_DIR_DOCKER --keyring-backend=test --starting-ip-address "$STARTING_IP" --chain-id "$CHAINID"

echo "Create and add Orchestrator keys"
echo "$MNEMONIC" | $CHAIND keys add "$KEY" --home "$DATA_DIR_DOCKER" --no-backup --chain-id "$CHAINID" --keyring-backend test --recover
echo "Init $CHAINID with moniker=$MONIKER and chain-id=$CHAINID"
$CHAIND init "$MONIKER" --chain-id "$CHAINID" --home "$DATA_DIR_DOCKER"

echo "Prepare genesis..."
echo "- Set gas limit in genesis"
jq '.consensus_params["block"]["max_gas"]="10000000"' "$GENESIS" > "$TEMP_GENESIS" && mv "$TEMP_GENESIS" "$GENESIS"

echo "- Set $DENOM as denom"
sed -i.bak "s/aphoton/$DENOM/g" $GENESIS
sed -i.bak "s/stake/$DENOM/g" $GENESIS

# Change proposal periods to pass within a reasonable time for local testing
sed -i.bak 's/"max_deposit_period": "172800s"/"max_deposit_period": "30s"/g' "$GENESIS"
sed -i.bak 's/"voting_period": "172800s"/"voting_period": "30s"/g' "$GENESIS"
# Change proposal required quorum to 15%, so with the orchestrator vote the proposals pass
sed -i.bak 's/"quorum": "0.334000000000000000"/"quorum": "0.150000000000000000"/g' "$GENESIS"

echo "- Allocate genesis accounts"
GENESIS_KEY="$($CHAIND keys show $KEY -a --home $DATA_DIR_DOCKER --keyring-backend test | grep -v INFO)"
# need to sleep here on MAC because file changes are not reflected instantly?
sleep 1
echo $CHAIND add-genesis-account $GENESIS_KEY 100000000000000000000000000000000$DENOM \
--home $DATA_DIR_DOCKER --keyring-backend test
$CHAIND add-genesis-account $GENESIS_KEY 100000000000000000000000000000000$DENOM \
--home $DATA_DIR_DOCKER --keyring-backend test

echo "- Sign genesis transaction"
$CHAIND gentx $KEY 100000000000000000000$DENOM --keyring-backend test --home $DATA_DIR_DOCKER --chain-id $CHAINID

echo "- Add all other validators genesis accounts"
for i in $(find $BUILD_DIR/gentxs -name "*.json")
do
    address=$(cat "$i" | jq '.body.messages[0].delegator_address'|tr -d '"')
    $CHAIND add-genesis-account  "$address" 100000000000000000000000000$DENOM --home $DATA_DIR_DOCKER --keyring-backend test
    [ $? -eq 0 ] && echo "$address added" || echo "$address failed"
done

# add gentx to gentxs dir
cp $CONF_DIR/gentx/*.json $BUILD_DIR/gentxs/node4.json

echo "- Collect genesis tx"
$CHAIND collect-gentxs --gentx-dir $BUILD_DIR_DOCKER/gentxs --home $DATA_DIR_DOCKER

echo "- Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
$CHAIND validate-genesis --home $DATA_DIR_DOCKER

echo "- Distribute final genesis.json to all validators"
TARGET_NODE_COUNT=10
for i in $(ls $BUILD_DIR | grep 'node');do
    # TODO uncomment this when issue https://github.com/evmos/ethermint/issues/1579 is solved
    # cp $GENESIS $BUILD_DIR/$i/$CHAIND/config/genesis.json
    cp $GENESIS $BUILD_DIR/$i/ethermintd/config/genesis.json
    [ $? -eq 0 ] && echo "$i: genesis updated successfully" || echo "$i: genesis update failed"
    cp $CONF_DIR/client.toml $BUILD_DIR/$i/ethermintd/config/client.toml
    cp multi-node-start.sh $BUILD_DIR/$i/ethermintd/
    # in future good faucet will be installed in the image
    # but image is not released yet at the time of this writing
    NODE_NO=$( echo "$i" | sed -E 's/^node//g' )
    LOGS_STREAM=val-$TARGET_NODE_COUNT
    cp ../faucet.py $BUILD_DIR/$i/ethermintd/

    TARGET_NODE_COUNT=$(( TARGET_NODE_COUNT + 1 ))
done

echo "copy config.toml to get the seeds"
# TODO uncomment this when issue https://github.com/evmos/ethermint/issues/1579 is solved
# cp $BUILD_DIR/node0/$CHAIND/config/config.toml $CONFIG
cp $BUILD_DIR/node0/ethermintd/config/config.toml $CONFIG
sed -i.bak 's/moniker = \"node0\"/moniker = \"orchestrator\"/g' $CONFIG

echo "copy app.toml to have same config on all nodes"
# TODO uncomment this when issue https://github.com/evmos/ethermint/issues/1579 is solved
# cp $BUILD_DIR/node0/$CHAIND/config/config.toml $CONF_DIR/app.toml
cp $BUILD_DIR/node0/ethermintd/config/app.toml $CONF_DIR/app.toml

# create full node directory without validator private keys
TARGET_NODE_COUNT=$(( FULL_NODE_FIRST_IDX + 10 ))
END_IDX=$(( $FULL_NODE_LAST_IDX - 1 ))
TARGET_DIR=$BUILD_DIR/rpcnode
echo $TARGET_DIR
cp -r $BUILD_DIR/node0 $TARGET_DIR
rm $BUILD_DIR/rpcnode/ethermintd/key_seed.json
rm $BUILD_DIR/rpcnode/ethermintd/config/node_key.json
rm $BUILD_DIR/rpcnode/ethermintd/config/priv_validator_key.json
rm $BUILD_DIR/rpcnode/ethermintd/data/priv_validator_state.json
rm -rf $BUILD_DIR/rpcnode/ethermintd/keyring-test
