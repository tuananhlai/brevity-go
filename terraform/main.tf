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

  name                    = "brevity-vpc"
  cidr                    = local.default_vpc_cidr
  azs                     = slice(data.aws_availability_zones.available.names, 0, 2)
  public_subnets          = ["10.0.0.0/20", "10.0.16.0/20"]
  private_subnets         = ["10.0.32.0/20", "10.0.48.0/20"]
  enable_dns_hostnames    = true
  enable_dns_support      = true
  map_public_ip_on_launch = true
}
