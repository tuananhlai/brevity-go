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

  enable_ipv6                                    = true
  public_subnet_assign_ipv6_address_on_creation  = true
  private_subnet_assign_ipv6_address_on_creation = true
  public_subnet_ipv6_prefixes                    = [0, 1]
  private_subnet_ipv6_prefixes                   = [2, 3]
}

// For some stupid reason, ECS endpoint doesn't support IPv6, so a NAT instance
// is still necessary to register EC2 instances with ECS.
// If they have fixed that issue, the following command will return an IPv6 address.
// > dig AAAA +short ecs.us-east-1.amazonaws.com
# module "fck-nat" {
#   source  = "RaJiska/fck-nat/aws"
#   version = "~> 1.0"

#   instance_type      = "t2.micro"
#   name               = "brevity-fck-nat-instance"
#   vpc_id             = module.vpc.vpc_id
#   subnet_id          = module.vpc.public_subnets[0]
#   use_spot_instances = true

#   update_route_tables = true
#   route_tables_ids = {
#     for index, rt_id in toset(module.vpc.private_route_table_ids) :
#     rt_id => rt_id
#   }
# }
