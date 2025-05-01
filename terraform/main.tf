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
  version = "~> 5.0"

  name            = "brevity-vpc"
  cidr            = local.default_vpc_cidr
  azs             = slice(data.aws_availability_zones.available.names, 0, 2)
  public_subnets  = ["10.0.0.0/20", "10.0.16.0/20"]
  private_subnets = ["10.0.32.0/20", "10.0.48.0/20"]
}

resource "random_password" "db_password" {
  length  = 16
  special = false
  upper   = true
  numeric = true
}

module "db_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  vpc_id          = module.vpc.vpc_id
  name            = "brevity-db-sg"
  use_name_prefix = true

  // TODO: Make the CIDR block more restrictive.
  ingress_with_cidr_blocks = [
    {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  egress_with_cidr_blocks = [
    {
      protocol         = "-1"
      from_port        = 0
      to_port          = 0
      cidr_blocks      = "0.0.0.0/0"
      ipv6_cidr_blocks = "::0/0"
    }
  ]
}

resource "aws_db_subnet_group" "primary" {
  name       = "brevity-db-subnet-group"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_db_instance" "primary" {
  identifier_prefix      = "brevity-"
  engine                 = "postgres"
  engine_version         = "17.2"
  instance_class         = "db.t4g.micro"
  allocated_storage      = 20
  db_name                = "brevity"
  username               = "postgres"
  password               = random_password.db_password.result
  apply_immediately      = true
  skip_final_snapshot    = true
  vpc_security_group_ids = [module.db_sg.security_group_id]
  db_subnet_group_name   = aws_db_subnet_group.primary.name
}

output "primary" {
  value = {
    password = nonsensitive(random_password.db_password.result),
    username = aws_db_instance.primary.username
    address  = aws_db_instance.primary.address
  }
}
