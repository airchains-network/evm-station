#!/bin/sh

CLIENT_PRIVATE_KEY=$1
REMOTE_VPN_PUBLIC_KEY=$2
REMOTE_VPN_PUBLIC_IP=$3

cat <<EOF
[Interface]
PrivateKey = $(cat $CLIENT_PRIVATE_KEY)
#Address = 172.49.0.10/24
ListenPort = 51877

[Peer]
PublicKey = $(cat $REMOTE_VPN_PUBLIC_KEY)
AllowedIPs = 172.49.0.100/32,10.0.0.0/24
Endpoint = $REMOTE_VPN_PUBLIC_IP:51820
PersistentKeepalive = 25
EOF
