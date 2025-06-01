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
      protocol    = "tcp"
      from_port   = 5432
      to_port     = 5432
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  ingress_with_ipv6_cidr_blocks = [
    {
      protocol         = "tcp"
      from_port        = 5432
      to_port          = 5432
      ipv6_cidr_blocks = "::/0"
    }
  ]

  egress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  egress_with_ipv6_cidr_blocks = [
    {
      rule             = "all-all"
      ipv6_cidr_blocks = "::/0"
    }
  ]
}

resource "aws_db_subnet_group" "primary" {
  name_prefix = "brevity-db-subnet-group-"
  subnet_ids  = module.vpc.private_subnets
}

resource "aws_db_instance" "primary" {
  identifier_prefix      = "brevity-"
  engine                 = "postgres"
  engine_version         = "17.4"
  instance_class         = "db.t4g.micro"
  allocated_storage      = 20
  db_name                = "brevity"
  username               = "postgres"
  password               = random_password.db_password.result
  apply_immediately      = true
  skip_final_snapshot    = true
  vpc_security_group_ids = [module.db_sg.security_group_id]
  db_subnet_group_name   = aws_db_subnet_group.primary.name
  # network_type           = "DUAL"
}

output "primary_db" {
  value = {
    password = nonsensitive(random_password.db_password.result),
    username = aws_db_instance.primary.username
    host     = aws_db_instance.primary.address
    port     = aws_db_instance.primary.port
  }
}
