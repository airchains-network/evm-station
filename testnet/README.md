# Terraform test environment

## Introduction

This allows running example testnet with 5 nodes, could be scaled up to 90.

To customize behaviour look into `variables.tf` file.

To not prompt public key prepare tfvars file with your public key:
```
echo 'ssh_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFCo2+r+BHQXwUn0hgS2ZEP79RUmz9T6JUKfMw4f1iAe deividas@Deividass-MacBook-Pro.local"' > terraform.tfvars
```

## Requirements

### Software

`sshuttle`, `terraform`, `jq`

### Credentials

Github token and user to allow pulling docker image from `ghcr.io/zama-ai/ethermint-node`

## Running terraform

Run plan
```
terraform plan
```

Provision the resources
```
terraform apply
```

## Connecting to nodes

Once terraform apply succeeds, get public ip of the single vpn node with:
```
cat terraform.tfstate | grep 'public_ip": "[0-9]'
```
Example output:
```
"public_ip": "35.177.70.69",
```

You can ssh to this node from which you'll reach all the other nodes in the setup. Only this node has public ip.

## VPN

To communicate/develop with other nodes you can use `sshuttle`
```
sshuttle -r ubuntu@35.177.70.69 10.0.0.0/24
```

After this you can ssh to the 5 (or more if you configured) worker nodes
```
ssh ubuntu@10.0.0.10
ssh ubuntu@10.0.0.11
ssh ubuntu@10.0.0.12
ssh ubuntu@10.0.0.13
ssh ubuntu@10.0.0.14
```

All worker nodes have only private ip and don't expose anything to the internet.

However, they have public internet access to download packages and etc.

All traffic is allowed between private nodes.

## Testing

To spin up all validators run

```
make full_setup
```

To destroy environment
```
make teardown
```

## Running multiple environments

Create file `terraform/terraform.tfvars` to override default values:

```
# your ssh key
ssh_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFCo2+r+BHQXwUn0hgS2ZEP79RUmz9T6JUKfMw4f1iAe"
# make sure azs belong to the same region, and extra_az is different from first az
region = "eu-central-1"
az = "eu-central-1a"
elb_extra_az = "eu-central-1b"
# make sure DNS names are unique globally for your environment``
explorer_dns_name = "explorer.staging.zama.ai"
rpc_dns_name = "rpc.staging.zama.ai"
ws_dns_name = "ws.staging.zama.ai"
grafana_dns_name = "grafana.staging.zama.ai"
chain_name = "ZAMA Staging Devnet"
# pick AMI for the region from https://cloud-images.ubuntu.com/locator/ec2/
ami = "ami-04d9351fa78a6efea"
# make sure to change environment to your unique environment
# it should have no spaces and be all lowercase letters
environment = "staging"
```
