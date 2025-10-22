provider "aws" {
  region = "us-east-1"
}

locals {
  default_vpc_cidr = "10.0.0.0/16"
  ssh_key_name     = "jump"
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
  map_public_ip_on_launch = true

  // I enabled IPv6 before, but it is causing more problems than it solves and
  // even AWS, the provider who is charging for IPv4, doesn't have full support
  // for it. So I've decided to disable everything related to IPv6.
  enable_ipv6                                                   = true
  public_subnet_assign_ipv6_address_on_creation                 = false
  private_subnet_assign_ipv6_address_on_creation                = false
  public_subnet_ipv6_prefixes                                   = [0, 1]
  private_subnet_ipv6_prefixes                                  = [2, 3]
  public_subnet_enable_dns64                                    = false
  private_subnet_enable_dns64                                   = false
  public_subnet_enable_resource_name_dns_aaaa_record_on_launch  = false
  private_subnet_enable_resource_name_dns_aaaa_record_on_launch = false
}

// For some stupid reason, ECS endpoint doesn't support IPv6, so a NAT instance
// is still necessary to register EC2 instances with ECS.
// If they have fixed that issue, the following command will return an IPv6 address.
// > dig AAAA +short ecs.us-east-1.amazonaws.com
module "fck-nat" {
  source  = "RaJiska/fck-nat/aws"
  version = "~> 1.0"

  instance_type = "t2.micro"
  name          = "brevity-fck-nat-instance"
  vpc_id        = module.vpc.vpc_id
  subnet_id     = module.vpc.public_subnets[0]
  use_ssh       = true
  ssh_key_name  = local.ssh_key_name
  ssh_cidr_blocks = {
    ipv4 = ["0.0.0.0/0"]
  }

  update_route_tables = true
  route_tables_ids = {
    for index, rt_id in toset(module.vpc.private_route_table_ids) :
    rt_id => rt_id
  }
}
