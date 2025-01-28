#!/bin/env python3

import json
import os
import subprocess
import sys

if len(sys.argv) < 2:
    print("Pass address for tokens in the first argument")

FAUCET_DEST_ADDRESS = sys.argv[1]
FAUCET_WALLET_NAME = os.getenv("FAUCET_WALLET_NAME", default="orchestrator")
FAUCET_AMOUNT = os.getenv("FAUCET_AMOUNT", default="1000000000000000000")
DENOM = os.getenv("DENOM", default="aphoton")

IS_ZAMA_TESTNET = os.getenv('ZAMA_TESTNET')

def get_faucet_address():
    keys_list_command = ['ethermintd', '--output=json', 'keys', 'list']
    if IS_ZAMA_TESTNET:
        keys_list_command.append('--keyring-backend=test')

    # Run the ethermintd command and capture the output
    output_bytes = subprocess.check_output(keys_list_command)

    # Convert bytes to string
    output_str = output_bytes.decode("utf-8")

    # Split the lines of the output
    output_lines = output_str.split("\n")

    # Filter out lines starting with "INFO"
    filtered_lines = [line for line in output_lines if not line.startswith("INFO")]

    # Join the filtered lines back into a single string
    filtered_output_str = "\n".join(filtered_lines)

    # Parse the JSON from the filtered output
    addresses = json.loads(filtered_output_str)

    for address in addresses:
        if address["name"] == FAUCET_WALLET_NAME:
            return address["address"]
    return None


def get_bech32_addr(ethereum_address):
    output = subprocess.check_output(
        ["ethermintd", "debug", "addr", ethereum_address]
    ).decode("utf-8")
    bech32 = next((x for x in output.splitlines() if x.startswith("Bech32 Acc:")))
    return bech32.split(": ")[1]


faucet_address = get_faucet_address()
if faucet_address is None:
    print("Faucet account not found with name " + FAUCET_WALLET_NAME)
    sys.exit(1)

dst_bech_addr = get_bech32_addr(FAUCET_DEST_ADDRESS)
maybe_keyring_backend = ''
if IS_ZAMA_TESTNET:
    maybe_keyring_backend = '--keyring-backend=test'
os.system(
    f"ethermintd --output=json tx bank send {maybe_keyring_backend} {faucet_address} {dst_bech_addr} \
        {FAUCET_AMOUNT}{DENOM} --from {FAUCET_WALLET_NAME} \
        --gas-prices 1000000000{DENOM} -y"
)