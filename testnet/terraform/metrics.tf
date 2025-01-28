resource "aws_instance" "metrics" {
  ami = "${var.ami}"
  instance_type = "c5a.4xlarge"
  associate_public_ip_address = false
  private_ip = "10.0.0.80"
  subnet_id = aws_subnet.private_subnet.id
  key_name = aws_key_pair.machine_key.id
  iam_instance_profile = aws_iam_instance_profile.logs_ec2_instance_profile.name

  security_groups = [aws_security_group.worker_sg.id]

  root_block_device {
    volume_size = "128"
    volume_type = "gp3"
    # we need good throughput disk performance
    # to display dashboards quickly
    iops = 16000
    throughput = 1000
  }

  # forces recreation every time if not here for some reason
  lifecycle {
    ignore_changes = [
      security_groups
    ]
  }

  tags = {
    Name = "${var.environment} Metrics node"
  }

  # make sure VPN/NAT node is setup before worker nodes
  depends_on = [
    aws_instance.validator_testnet_vpn
  ]
}
