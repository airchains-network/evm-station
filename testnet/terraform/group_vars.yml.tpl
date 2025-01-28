testnet_blockchain_node_docker_image: ghcr.io/zama-ai/ethermint-node:v0.3.0-3
testnet_blockchain_kms_docker_image: ghcr.io/zama-ai/kms:v0.1.4
fhevm_version: v0.4.0-0
testnet_chain_id: ethermint_8009-1
numeric_chain_id: 8009
aws_logs_region: eu-west-2
aws_logs_group: zbc-devnet-logs
blockscout_version: 5.3.3
chain_name: ${chain_name}
currency_name: Zama
currency_symbol: ZAMA
explorer_dns_name: ${explorer_dns_name}
grafana_image: grafana/grafana:10.2.4-ubuntu
# TODO: move out password somewhere else?
grafana_admin_password: eTKgEqY9CtpqQAm4eT7OD
grafana_tempo_image: grafana/tempo:r134-e63f49f

faucet_test_accounts:
  - '0xa5e1defb98EFe38EBb2D958CEe052410247F4c80'
  - '0xfCefe53c7012a075b8a711df391100d9c431c468'
  - '0xa44366bAA26296c1409AD1e284264212029F02f1'
  - '0xc1d91b49A1B3D1324E93F86778C44a03f1063f1b'
  - '0x305F1F471e9baCFF2b3549F9601f9A4BEafc94e1'

# generate with this gen.js script inside fhevm repo
#
# const ethers = require('ethers')
# const bip39 = require('bip39')
# const phrase = bip39.generateMnemonic()
# const mnemonic = ethers.Mnemonic.fromPhrase(phrase)
# console.log('  - mnemonic: ' + phrase)
# console.log('    addresses:')
# const wallet = ethers.HDNodeWallet.fromMnemonic(mnemonic, `m/44'/60'/0'/0`)
# for (var i = 0; i < 5; i++) {
#   console.log("    - '" + wallet.deriveChild(i).address + "'")
# }
#
# then run 'node gen.js'
stress_test_accounts:
  - mnemonic: inspire engage system neither recipe bridge slow approve together jeans hospital side
    addresses:
    - '0xf8DD17aC1994dCa321C41bE71E48102130594785'
    - '0x3484D600694629874DC0470899a88a3819035d6D'
    - '0x31b7da659422A088160baDfE3523aC1e1c25CEac'
    - '0x5937A19Ff0980e6E08cEB36A4DDb141087B6b489'
    - '0x10d680eCe02863e97905d3C7f8d16CF665e67eF5'
  - mnemonic: ask stem speak home deny december crane raven sell boat artefact sign
    addresses:
    - '0xafef55e4aD14Cd0788c6737559b864D5F57f59e4'
    - '0xe149A046Ca6530634A8500207B2023672fBa71DC'
    - '0x0aFA3F55210CEE2b35352C6fA981323cDCA11747'
    - '0x5b391e0B87E2f78899D040c2bBfD14161be36c36'
    - '0x2Da77EeF4Fcc680F76830e75cB900e44a4c019A6'
  - mnemonic: join tissue afford sun run element spend west zero cupboard silent gentle
    addresses:
    - '0x8d19000880d698C82DBD36a4173973Ca8B62F0CF'
    - '0xBa30C698c57fD95f802dE4369371f4EC56D48478'
    - '0xd797CcbEA2998d76e2B1BaBFbF6AF2a87953441c'
    - '0x7a2Fe2adBd8Cf8e769dCcF1Ec334c985e39FcA8a'
    - '0x522B2da19C6477682dbf8174E8de5b0AA65AF128'
  - mnemonic: diesel label dutch eternal soap manual palace vast olive fix ocean giant
    addresses:
    - '0x5DC0b3F53c11bf67E97024371211177e7807145f'
    - '0x47669B6AEE629EAf1bd9F287f7d69D616250A4f0'
    - '0xaFBD6ECDB80eE54C4abe328b88C2D3937385726d'
    - '0x6455d710fc7717aDd24a7907Fb7511f234118b7F'
    - '0x3147880C5F7B818eF33E7d20897B9378d2AA5ff9'
  - mnemonic: trash tattoo ridge phrase february primary brief assault slice basket pyramid month
    addresses:
    - '0x3f261EE36ca0491BC9bF005fEd63f569110667e6'
    - '0x9FF7B2361D4b45FE323d54fC378BcA935CC29aB3'
    - '0x22eeA6cb7588048653C35ad0740FADDB0D9A5B9F'
    - '0xF34644922a164Ec69CE5c3fe8bE7bD5b834BA9a4'
    - '0x59fa427CCbF46d9830B51271699E9f339dBA48B6'
  - mnemonic: border shove offer label sister submit nerve save page define bind you
    addresses:
    - '0x7ebaA94617E334d72600ebe6a4C1d005d3fA20dF'
    - '0xaCe889a561d641336ee21011C406766154939A23'
    - '0x07cE906e3f28248750c49E9220A1DdeE1ED516a8'
    - '0x8d517473c54347561Ba79919a60a1938AE3b84EB'
    - '0x6DbDEAA1BcD68c2b25dF5EdD7a1b2aeB8CEc31C4'
