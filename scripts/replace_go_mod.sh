#!/bin/bash


set -Eeuo pipefail

if [ "$#" -lt 2 ]; then
    echo "Please give the file to update and the custom path"
    echo "Example: $(basename $0) <go.mod.updated> ./work_dir"
    exit
fi

FILE=$1
CUSTOM_PATH=$2
REPO_NAME=zbc-go-ethereum

CUSTOM_PATH=$(echo "$CUSTOM_PATH" | sed 's/[\/&]/\\&/g')

echo "FILE: $FILE"
echo "CUSTOM_PATH: $CUSTOM_PATH"

grep -n "zama.ai/$REPO_NAME" "$FILE" | cut -d: -f1 | xargs -I{} sed -i -e "{}s/=>.*/=> $CUSTOM_PATH\/$REPO_NAME/" "$FILE"

# fhevm go
sed -i -E "/PLACEHOLDER_FHEVM_GO/s|// PLACEHOLDER_FHEVM_GO.*|$(grep 'fhevm-go' "$FILE" ) => $CUSTOM_PATH/fhevm-go|" "$FILE" 