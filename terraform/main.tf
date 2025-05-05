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

// == Database ==

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
  name            = "brevity-db-sg-"
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
  name_prefix = "brevity-db-subnet-group-"
  subnet_ids  = module.vpc.private_subnets
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

// == ECS ==

module "ecs_ec2_iam" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"
  version = "~> 5.0"

  create_role             = true
  role_name_prefix        = "brevity-ecs-ec2-"
  create_instance_profile = true

  create_custom_role_trust_policy = true
  custom_role_trust_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  custom_role_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
  ]
}

module "ecs_ec2_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  vpc_id          = module.vpc.vpc_id
  name            = "brevity-ecs-ec2-sg-"
  use_name_prefix = true

  // TODO: Make the CIDR block more restrictive.
  ingress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "all"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  egress_with_cidr_blocks = [
    {
      protocol         = "all"
      from_port        = 0
      to_port          = 0
      cidr_blocks      = "0.0.0.0/0"
      ipv6_cidr_blocks = "::0/0"
    }
  ]
}

data "aws_ami" "amz_linux_2023" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name = "name"
    // Note that we must use an AMI optimized for ECS instead of
    // a general purpose AMI like the ones you see when launching a new EC2 instance.
    values = ["al2023-ami-ecs-hvm-*-kernel-6.1-x86_64"]
  }
}

resource "aws_ecs_cluster" "default" {
  name = "brevity-ecs-cluster"
}

resource "aws_launch_template" "ecs_lt" {
  name_prefix            = "brevity-ecs-lt-"
  image_id               = data.aws_ami.amz_linux_2023.id
  instance_type          = "t2.micro"
  vpc_security_group_ids = [module.ecs_ec2_sg.security_group_id]

  iam_instance_profile {
    arn = module.ecs_ec2_iam.iam_instance_profile_arn
  }

  // https://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_container_instance.html#linux-liw-advanced-details
  user_data = base64encode(<<EOF
#!/bin/bash
echo ECS_CLUSTER=${aws_ecs_cluster.default.name} >> /etc/ecs/ecs.config
    EOF
  )

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "brevity-ecs-instance"
    }
  }
}

resource "aws_autoscaling_group" "ecs_asg" {
  name_prefix         = "brevity-ecs-asg-"
  vpc_zone_identifier = module.vpc.private_subnets
  min_size            = 0
  desired_capacity    = 1
  max_size            = 3

  launch_template {
    id      = aws_launch_template.ecs_lt.id
    version = "$Latest"
  }

  tag {
    key                 = "AmazonECSManaged"
    value               = true
    propagate_at_launch = true
  }
}

resource "aws_ecr_repository" "default" {
  name = "brevity-ecs-repo"
}

module "ecs_alb_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  vpc_id          = module.vpc.vpc_id
  name            = "brevity-ecs-alb-sg-"
  use_name_prefix = true

  // TODO: Make the CIDR block more restrictive.
  ingress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "all"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  egress_with_cidr_blocks = [
    {
      protocol         = "all"
      from_port        = 0
      to_port          = 0
      cidr_blocks      = "0.0.0.0/0"
      ipv6_cidr_blocks = "::0/0"
    }
  ]
}

resource "aws_lb" "ecs" {
  load_balancer_type = "application"
  security_groups    = [module.ecs_alb_sg.security_group_id]
  subnets            = module.vpc.public_subnets
}

resource "aws_lb_target_group" "ecs" {
  port        = 80
  protocol    = "HTTP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"
}

resource "aws_lb_listener" "ecs" {
  port              = 80
  protocol          = "HTTP"
  load_balancer_arn = aws_lb.ecs.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ecs.arn
  }
}

// == Github Actions ==
// https://aws.amazon.com/blogs/security/use-iam-roles-to-connect-github-actions-to-actions-in-aws/

resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com"
  ]
}

resource "aws_iam_role" "github_actions" {
  name_prefix = "brevity-github-actions-"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRoleWithWebIdentity"
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
            "token.actions.githubusercontent.com:sub" = "repo: tuananhlai/brevity-go:ref:refs/heads/main",
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "github_actions" {
  name_prefix = "ecr-ecs-deploy-"
  role        = aws_iam_role.github_actions.name

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid    = "ECRAccess",
        Effect = "Allow",
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload"
        ],
        Resource = "*"
      },
      {
        Sid    = "ECSAccess",
        Effect = "Allow",
        Action = [
          "ecs:DescribeServices",
          "ecs:DescribeTaskDefinition",
          "ecs:RegisterTaskDefinition",
          "ecs:UpdateService"
        ],
        Resource = "*"
      },
      {
        Sid : "Logs",
        Effect : "Allow",
        Action : [
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams",
          "logs:GetLogEvents"
        ],
        Resource : "*"
      }
    ]
  })
}

output "primary_db" {
  value = {
    password = nonsensitive(random_password.db_password.result),
    username = aws_db_instance.primary.username
    host     = aws_db_instance.primary.address
    port     = aws_db_instance.primary.port
  }
}

output "ecr" {
  value = {
    name = aws_ecr_repository.default.name
    url  = aws_ecr_repository.default.repository_url
  }
}

output "github_actions" {
  value = {
    arn = aws_iam_role.github_actions.arn
  }
}
