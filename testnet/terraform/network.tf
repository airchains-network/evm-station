provider "aws" {
  region = "${var.region}"
}

# Initialize a new VPC
resource "aws_vpc" "validator_testnet_vpc" {
  cidr_block = "${var.cidr_vpc}"
}

# Public subnet
resource "aws_subnet" "public_subnet" {
  cidr_block = "${var.public_subnet}"
  vpc_id = aws_vpc.validator_testnet_vpc.id
  availability_zone = "${var.elb_extra_az}"
}

resource "aws_internet_gateway" "validator_testnet_igw" {
  vpc_id = aws_vpc.validator_testnet_vpc.id
}

resource "aws_eip" "nat_gw" {
  domain   = "vpc"
}

resource "aws_nat_gateway" "nat_gw" {
  allocation_id = aws_eip.nat_gw.id
  subnet_id     = aws_subnet.public_subnet.id

  tags = {
    Name = "${var.environment} NAT gateway"
  }

  depends_on = [aws_internet_gateway.validator_testnet_igw]
}

resource "aws_route_table" "public_subnet_route_table" {
  vpc_id = aws_vpc.validator_testnet_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.validator_testnet_igw.id
  }
}

resource "aws_route_table_association" "public_subnet" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_subnet_route_table.id
}

# Private subnet
resource "aws_subnet" "private_subnet" {
  cidr_block = "${var.private_subnet}"
  vpc_id = aws_vpc.validator_testnet_vpc.id
  availability_zone = "${var.az}"
}

resource "aws_route_table" "private_subnet_route_table" {
  vpc_id = aws_vpc.validator_testnet_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.nat_gw.id
  }
}

resource "aws_route_table_association" "private_subnet" {
  subnet_id      = aws_subnet.private_subnet.id
  route_table_id = aws_route_table.private_subnet_route_table.id
}
