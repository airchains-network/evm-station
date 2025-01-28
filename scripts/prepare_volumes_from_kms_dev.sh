#!/bin/bash

# This bash script creates global fhe keys
# and copy them to the right folder in volumes directory.
# It accepts 
# - the version of kms-dev as the first parameter
# - the LOCAL_BUILD_PUBLIC_KEY_PATH as the second optional parameter.
# - the LOCAL_BUILD_PRIVATE_KEY_PATH as the third optional parameter.

set -Eeuo pipefail

if [ "$#" -lt 1 ]; then
    echo "Usage: $(basename "$0") <kms-dev version> [LOCAL_BUILD_PUBLIC_KEY_PATH] [LOCAL_BUILD_PRIVATE_KEY_PATH]"
    echo "Example: $(basename "$0") v0.3.2 $(PWD)/running_node/node1/.ethermintd/zama/keys/network-fhe-keys /PATH_TO_KMS_KEYS"
    exit
fi

KMS_DEV_VERSION=$1
BINARY_NAME="/app/kms/bin/kms-gen"
DOCKER_IMAGE=ghcr.io/zama-ai/kms-dev:"$KMS_DEV_VERSION"
CURRENT_FOLDER=$PWD


KEYS_FULL_PATH=$CURRENT_FOLDER/res/keys
mkdir -p $KEYS_FULL_PATH

if [ "$#" -ge 3 ]; then
    LOCAL_BUILD_PUBLIC_KEY_PATH=$2
    LOCAL_BUILD_PRIVATE_KEY_PATH=$3
    NETWORK_KEYS_PUBLIC_PATH="${LOCAL_BUILD_PUBLIC_KEY_PATH}"
    NETWORK_KEYS_PRIVATE_PATH="${LOCAL_BUILD_PRIVATE_KEY_PATH}"
else
    NETWORK_KEYS_PUBLIC_PATH="./volumes/network-public-fhe-keys"
    NETWORK_KEYS_PRIVATE_PATH="./volumes/network-private-fhe-keys"
fi

mkdir -p "$KEYS_FULL_PATH"

docker run -v "$PWD/res/keys:/keys" "$DOCKER_IMAGE" "$BINARY_NAME" /keys/

echo "$KEYS_FULL_PATH"

echo "###########################################################"
echo "Keys creation is done, they are stored in $KEYS_FULL_PATH"
echo "###########################################################"



echo "$NETWORK_KEYS_PUBLIC_PATH"
echo "$NETWORK_KEYS_PRIVATE_PATH"

MANDATORY_KEYS_LIST=('sks.bin' 'cks.bin' 'pks.bin' 'default-software-keys.bin')

for key in "${MANDATORY_KEYS_LIST[@]}"; do
    if [ ! -f "$KEYS_FULL_PATH/$key" ]; then
        echo "#####ATTENTION######"
        echo "$key does not exist in $KEYS_FULL_PATH!"
        echo "####################"
        exit
    fi
done

echo "###########################################################"
echo "All the required keys exist in $KEYS_FULL_PATH"
echo "###########################################################"

mkdir -p $NETWORK_KEYS_PUBLIC_PATH
mkdir -p $NETWORK_KEYS_PRIVATE_PATH

key="sks.bin"
echo "Copying $key to $NETWORK_KEYS_PUBLIC_PATH, please wait ..."
cp $KEYS_FULL_PATH/$key $NETWORK_KEYS_PUBLIC_PATH/sks

key="pks.bin"
echo "Copying $key to $NETWORK_KEYS_PUBLIC_PATH, please wait ..."
cp $KEYS_FULL_PATH/$key $NETWORK_KEYS_PUBLIC_PATH/pks

key="cks.bin"
echo "Copying $key to $NETWORK_KEYS_PRIVATE_PATH, please wait ..."
cp $KEYS_FULL_PATH/$key $NETWORK_KEYS_PRIVATE_PATH/cks
cp $KEYS_FULL_PATH/$key $NETWORK_KEYS_PRIVATE_PATH

key="default-software-keys.bin"
echo "Copying $key to $NETWORK_KEYS_PRIVATE_PATH, please wait ..."
cp $KEYS_FULL_PATH/$key $NETWORK_KEYS_PRIVATE_PATH