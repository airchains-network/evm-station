#!/bin/bash

# Check if a file name is provided as an argument
if [ $# -eq 0 ]; then
    echo "Usage: $0 <filename>"
    exit 1
fi

# Assign the first argument to a variable
file="$1"

# Check if the file exists
if [ ! -f "$file" ]; then
    echo "File does not exist: $file"
    exit 1
fi


# Extracting the version using grep and awk
version=$(grep 'ghcr.io/zama-ai/kms@' "$file" | awk -F':' '{print $3}' | tr -d '[:space:]')

echo $version
