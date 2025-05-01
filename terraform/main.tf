provider "aws" {
  region = "us-east-1"
}

locals {
  default_vpc_cidr = "10.0.0.0/16"
}

data "aws_availability_zones" "available" {
  state = "available"
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.21.0"

  name            = "BrevityVPC"
  cidr            = local.default_vpc_cidr
  azs             = slice(data.aws_availability_zones.available.names, 0, 2)
  public_subnets  = ["10.0.0.0/20", "10.0.16.0/20"]
  private_subnets = ["10.0.32.0/20", "10.0.48.0/20"]
}

resource "random_password" "db_password" {
  length  = 16
  special = true
  upper   = true
  numeric = true
}

resource "aws_db_instance" "primary" {
  identifier_prefix   = "brevity-db"
  engine              = "postgres"
  instance_class      = "db.t4g.micro"
  allocated_storage   = 20
  skip_final_snapshot = true
  db_name             = "brevity"
  username            = "postgres"
  password            = random_password.db_password.result
  apply_immediately   = true
}

output "primary" {
  value = {
    password = nonsensitive(random_password.db_password.result),
    username = aws_db_instance.primary.username
    address  = aws_db_instance.primary.address
  }
}
