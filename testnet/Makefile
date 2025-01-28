
export MNEMONIC ?= stumble tilt business detect father ticket major inner awake jeans name vibrant tribe pause crunch sad wine muscle hidden pumpkin inject segment rocket silver
export DOCKER_IMAGE ?= ghcr.io/zama-ai/ethermint-node:v0.3.0-1
export CHAINID ?= ethermint_8009-1
export AWS_REGION ?= eu-west-2
# TODO: different on linux
export OBJC_DISABLE_INITIALIZE_FORK_SAFETY = YES

PUBLIC_NODE_IP_EXPR = cd terraform/ && terraform output -raw vpn_instance_public_ip
LOCAL_WG_TUNNEL_DEVICE ?= utun5

.PHONY: check_ssh_key
check_ssh_key:
	ssh-add -L | grep 'ssh-'

.PHONY: check_env
check_env:
ifndef DOCKER_PULL_GITHUB_TOKEN
	$(error DOCKER_PULL_GITHUB_TOKEN is undefined, needed for pulling image from gcr in remote host)
endif
ifndef DOCKER_PULL_GITHUB_USER
	$(error DOCKER_PULL_GITHUB_USER is undefined, needed for pulling image from gcr in remote host)
endif

.PHONY: login_public_node
login_public_node: check_ssh_key terraform/terraform.tfstate
	ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ubuntu@$(shell $(PUBLIC_NODE_IP_EXPR))

.PHONY: setup_periodic_tests
setup_periodic_tests: check_ssh_key terraform/terraform.tfstate
	cat setup_periodic_tests.sh | \
	  ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
	  ubuntu@$(shell $(PUBLIC_NODE_IP_EXPR)) '/bin/bash'

terraform/terraform.tfvars:
	echo 'ssh_public_key = "$(shell ssh-add -L)" \
	region = "$(AWS_REGION)"' > terraform/terraform.tfvars

.PHONY: localnet
localnet: localnet/fhevm-keys
	./setup_genesis.sh

localnet/fhevm-keys:
	mkdir -p localnet/fhevm-keys
	fhevm-tfhe-cli generate-keys -d ./localnet/fhevm-keys

.PHONY: generate_genesis
generate_genesis: localnet

.PHONY: upload_keys
upload_keys: localnet/fhevm-keys
	rsync -avz --progress \
	 -e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \
	 localnet/fhevm-keys/ \
	 ubuntu@$(shell $(PUBLIC_NODE_IP_EXPR)):~/fhevm-keys/

.PHONY: upload_genesis
upload_genesis:
	./upload_genesis.sh

.PHONY: terraform_apply
terraform_apply: check_ssh_key terraform/terraform.tfvars
	cd terraform && terraform apply -auto-approve

.PHONY: terraform_plan
terraform_plan: check_ssh_key terraform/terraform.tfvars
	cd terraform && terraform plan

.PHONY: full_setup
full_setup: check_env check_ssh_key generate_genesis
	$(MAKE) terraform_apply
	# TODO: find better option to wait
	sleep 20 # wait for servers to boot
	$(MAKE) provision_vpn_server
	$(MAKE) start_wireguard_vpn
	$(MAKE) provision_servers
	$(MAKE) provision_faucet

ansible/vpn.privkey:
	wg genkey > ansible/vpn.privkey
	chmod 600 ansible/vpn.privkey

ansible/vpn.pubkey: ansible/vpn.privkey
	cat ansible/vpn.privkey | wg pubkey > ansible/vpn.pubkey

vpn-client.conf: ansible/vpn.privkey FORCE
	./generate_vpn_file.sh \
	  ansible/vpn.privkey \
	  ansible/remote-vpn-server-keys/wg0.pubkey \
	  $(shell $(PUBLIC_NODE_IP_EXPR)) > vpn-client.conf

FORCE: ;

.PHONY: start_wireguard_vpn
start_wireguard_vpn: vpn-client.conf stop_wireguard_vpn
	which wg wireguard-go
	-sudo wireguard-go $(LOCAL_WG_TUNNEL_DEVICE)
	-sudo wg setconf $(LOCAL_WG_TUNNEL_DEVICE) vpn-client.conf
	-sudo ifconfig utun5 inet 172.49.0.10 172.49.0.100 netmask 255.255.255.0
	-sudo route add -net 10.0.0.0/24 -host 172.49.0.10

.PHONY: stop_wireguard_vpn
stop_wireguard_vpn:
	-rm -f vpn-client.conf
	-sudo route delete -net 10.0.0.0/24 -host 172.49.0.10
	-sudo rm -f /var/run/wireguard/$(LOCAL_WG_TUNNEL_DEVICE).sock

.PHONY: teardown
teardown: stop_wireguard_vpn
	rm -rf ./localnet
	cd terraform && terraform destroy -auto-approve

.PHONY: provision_vpn_server
provision_vpn_server: check_env ansible/vpn.pubkey
	cd ansible && ansible-playbook --limit '$(shell $(PUBLIC_NODE_IP_EXPR))' --tags vpn provision-playbook.yml

.PHONY: provision_servers
provision_servers: check_env
	cd ansible && ansible-playbook --skip-tags vpn,faucet provision-playbook.yml

# for some reason when run in flow with ansible faucet doesn't work
# only later
.PHONY: provision_faucet
provision_faucet: check_env
	cd ansible && ansible-playbook --tags faucet provision-playbook.yml

.PHONY: install_ansible_dependencies
install_ansible_dependencies:
	cd ansible
	ansible-galaxy collection install ansible.posix community.general community.docker prometheus.prometheus
	ansible-galaxy role install geerlingguy.docker
