
[validators]
10.0.0.10 localnet_conf_dir=node0
10.0.0.11 localnet_conf_dir=node1
10.0.0.12 localnet_conf_dir=node2
10.0.0.13 localnet_conf_dir=node3
10.0.0.14 localnet_conf_dir=node4

[rpc_archival]
10.0.0.30 localnet_conf_dir=rpcnode
10.0.0.31 localnet_conf_dir=rpcnode
10.0.0.32 localnet_conf_dir=rpcnode

[kms]
10.0.0.50

[block_explorer]
10.0.0.70

[metrics]
10.0.0.80

[vpn]
${vpn_public_ip} ansible_ssh_user=ubuntu
