resource "aws_key_pair" "machine_key" {
  key_name   = "${var.environment}-validator-testnet-key-1"
  public_key = var.ssh_public_key
}

resource "local_file" "inventory_txt" {
  content = templatefile("${path.module}/inventory.txt.tpl",
    {
      vpn_public_ip = aws_instance.validator_testnet_vpn.public_ip
    }
  )
  filename = "${path.module}/../ansible/inventory.txt"
  file_permission = "0644"
}

resource "local_file" "ansible_group_vars" {
  content = templatefile("${path.module}/group_vars.yml.tpl",
    {
      chain_name = var.chain_name
      explorer_dns_name = var.explorer_dns_name
    }
  )
  filename = "${path.module}/../ansible/group_vars/all.yml"
  file_permission = "0644"
}
