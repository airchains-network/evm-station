# Define subnet CIDR blocks
locals {
  validator_suffix = range(10, 29)
  rpc_suffix = range(30, 49)
}

resource "aws_iam_instance_profile" "logs_ec2_instance_profile" {
  name = "${var.environment}-logs-ec2-instance-profile"

  # Attach the existing IAM role to the instance profile
  role = "zbc-cloud-watch"
}

# Create the VPN server in a public subnet with public IP and Elastic IP and update its variable to use the EIP
resource "aws_instance" "validator_testnet_vpn" {
  ami = "${var.ami}"
  instance_type = "c5n.2xlarge"
  associate_public_ip_address = true
  private_ip = "10.0.1.100"
  subnet_id = aws_subnet.public_subnet.id
  source_dest_check = false

  key_name = aws_key_pair.machine_key.id
  iam_instance_profile = aws_iam_instance_profile.logs_ec2_instance_profile.name

  security_groups = [aws_security_group.vpn_sg.id]

  # forces recreation every time if not here for some reason
  lifecycle {
    ignore_changes = [
      security_groups
    ]
  }

  root_block_device {
    volume_size = 20
  }

  tags = {
    Name = "${var.environment} ZBC environment VPN server"
  }
}

output "vpn_instance_public_ip" {
  value = aws_instance.validator_testnet_vpn.public_ip
}

resource "aws_instance" "validator_testnet_kms" {
  ami = "${var.ami}"
  instance_type = "${var.worker_type}"
  associate_public_ip_address = false
  private_ip = "10.0.0.50"
  subnet_id = aws_subnet.private_subnet.id

  key_name = aws_key_pair.machine_key.id
  iam_instance_profile = aws_iam_instance_profile.logs_ec2_instance_profile.name

  security_groups = [aws_security_group.worker_sg.id]

  root_block_device {
    volume_size = "${var.worker_disk_size_gb}"
  }

  # forces recreation every time if not here for some reason
  lifecycle {
    ignore_changes = [
      security_groups
    ]
  }

  tags = {
    Name = "${var.environment} ZBC environment KMS server"
  }
}

resource "aws_instance" "validator_testnet_workers" {
  count = var.validator_count
  ami = "${var.ami}"
  instance_type = "${var.worker_type}"
  associate_public_ip_address = false
  private_ip = "10.0.0.${local.validator_suffix[count.index]}"
  subnet_id = aws_subnet.private_subnet.id

  key_name = aws_key_pair.machine_key.id
  iam_instance_profile = aws_iam_instance_profile.logs_ec2_instance_profile.name

  security_groups = [aws_security_group.worker_sg.id]

  root_block_device {
    volume_size = "${var.worker_disk_size_gb}"
    volume_type = "gp3"
    iops = 16000
    throughput = 500
  }

  # forces recreation every time if not here for some reason
  lifecycle {
    ignore_changes = [
      security_groups
    ]
  }

  tags = {
    Name = "${var.environment} ZBC environment Validator ${count.index}"
  }

  # make sure VPN/NAT node is setup before worker nodes
  depends_on = [
    aws_instance.validator_testnet_vpn
  ]
}

resource "aws_instance" "validator_testnet_rpc_nodes" {
  count = var.full_node_count
  ami = "${var.ami}"
  instance_type = "${var.worker_type}"
  associate_public_ip_address = false
  private_ip = "10.0.0.${local.rpc_suffix[count.index]}"
  subnet_id = aws_subnet.private_subnet.id

  key_name = aws_key_pair.machine_key.id
  iam_instance_profile = aws_iam_instance_profile.logs_ec2_instance_profile.name

  security_groups = [aws_security_group.worker_sg.id]

  root_block_device {
    volume_size = "${var.worker_disk_size_gb}"
    volume_type = "gp3"
    iops = 16000
    throughput = 500
  }

  # forces recreation every time if not here for some reason
  lifecycle {
    ignore_changes = [
      security_groups
    ]
  }

  tags = {
    Name = "${var.environment} ZBC environment RPC server ${count.index}"
  }

  # make sure VPN/NAT node is setup before worker nodes
  depends_on = [
    aws_instance.validator_testnet_vpn
  ]
}

resource "aws_instance" "block_explorer" {
  ami = "${var.ami}"
  instance_type = "c5a.4xlarge"
  associate_public_ip_address = false
  private_ip = "10.0.0.70"
  subnet_id = aws_subnet.private_subnet.id

  key_name = aws_key_pair.machine_key.id
  iam_instance_profile = aws_iam_instance_profile.logs_ec2_instance_profile.name

  security_groups = [aws_security_group.worker_sg.id]

  root_block_device {
    volume_size = "512"
  }

  # forces recreation every time if not here for some reason
  lifecycle {
    ignore_changes = [
      security_groups
    ]
  }

  tags = {
    Name = "${var.environment} ZBC environment Block Explorer"
  }

  # make sure VPN/NAT node is setup before worker nodes
  depends_on = [
    aws_instance.validator_testnet_vpn
  ]
}
