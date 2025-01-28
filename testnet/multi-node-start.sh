#!/bin/bash


CHAIND=/usr/bin/ethermintd
DATA_DIR=/root/.ethermintd
CONFIG=$DATA_DIR/config/config.toml
APP_CONFIG=$DATA_DIR/config/app.toml

if [ -z "$CHAINID" ];
then
    echo "\$CHAINID environment variable is undefined"
    exit 1
fi

sed -i 's/prometheus = false/prometheus = true/g' $CONFIG
sed -i 's/enable-indexer = false/enable-indexer = true/g' $APP_CONFIG
sed -i 's/max-recv-msg-size = "10485760"/max-recv-msg-size = "104857600"/g' $APP_CONFIG
#perl -i -0pe 's/# Enable defines if the API server should be enabled.\nenable = false/# Enable defines if the API server should be enabled.\nenable = true/' $APP_CONFIG

sed -i 's/timeout_commit = "5s"/timeout_commit = "7s"/g' "$CONFIG"
# make sure the localhost IP is 0.0.0.0
sed -i 's/pprof_laddr = "localhost:6060"/pprof_laddr = "0.0.0.0:6060"/g' "$CONFIG"
sed -i 's/127.0.0.1/0.0.0.0/g' "$APP_CONFIG"
sed -i 's/localhost/0.0.0.0/g' "$APP_CONFIG"

# disable state sync
sed -i.bak 's/enable = true/enable = false/g' "$CONFIG"

# sed -i.bak 's/db_backend = "goleveldb"/db_backend = "rocksdb"/g' "$CONFIG"

# Change max_subscription to for bots workers
# toml-cli set $CONFIG rpc.max_subscriptions_per_client 500
# Change max_subscription to for bots workers
sed -i.bak 's/max_subscriptions_per_client = 5/max_subscriptions_per_client = 600/g' "$CONFIG"

sed -i.bak 's/indexer = "null"/indexer = "kv"/g' "$CONFIG"
sed -i.bak 's/namespace = "tendermint"/namespace = "cometbft"/g' "$CONFIG"

if [ ! -d /root/.ethermintd/keyring-test ]; then
    # this is full archival node in the testnet
    sed -i 's/pruning = "default"/pruning = "nothing"/g' "$APP_CONFIG"
    ARCHIVAL_NODE_ARGS="--json-rpc.api=eth,net,web3,debug,txpool --pruning=nothing"
fi

echo "running $CHAINID with extra flags $EXTRA_FLAGS"
echo "starting $CHAINID node in background ..."
echo "$CHAIND start $ARCHIVAL_NODE_ARGS --rpc.unsafe --keyring-backend test "$EXTRA_FLAGS" >"$DATA_DIR"/node.log"
$CHAIND start --rpc.unsafe \
--json-rpc.enable true --api.enable $ARCHIVAL_NODE_ARGS \
--keyring-backend test --chain-id $CHAINID $EXTRA_FLAGS
