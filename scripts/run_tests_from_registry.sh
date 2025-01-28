#!/bin/bash

cp .env.example .env


PRIVATE_KEY_ORACLE_DEPLOYER=$(grep PRIVATE_KEY_ORACLE_DEPLOYER .env | cut -d '"' -f 2)

npx hardhat task:computePredeployAddress --private-key "$PRIVATE_KEY_ORACLE_DEPLOYER"

echo "computePredeployAddress"

PRIVATE_KEY_ORACLE_RELAYER=$(grep PRIVATE_KEY_ORACLE_RELAYER .env | cut -d '"' -f 2)

ORACLE_CONTRACT_PREDEPLOY_ADDRESS=$(grep ORACLE_CONTRACT_PREDEPLOY_ADDRESS oracle/.env.oracle | cut -d '=' -f2)
  
sleep 30

npx hardhat compile:specific --contract lib

npx hardhat compile:specific --contract oracle

TEST_CONTAINER_NAME=ethermintnode0 npx hardhat task:launchFhevm

sleep 60

TEST_CONTAINER_NAME=ethermintnode0  npm run test:inband

