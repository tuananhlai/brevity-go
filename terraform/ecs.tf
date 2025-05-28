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
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  ingress_with_ipv6_cidr_blocks = [
    {
      rule             = "all-all"
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

// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/container-instance-eni.html
resource "aws_ecs_account_setting_default" "awsvpcTrunking" {
  name  = "awsvpcTrunking"
  value = "enabled"
}

resource "aws_ecr_repository" "default" {
  name = "brevity-ecs-repo"
}

module "ecs_alb_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  vpc_id          = module.vpc.vpc_id
  name            = "brevity-ecs-alb-sg"
  use_name_prefix = true

  // TODO: Make the CIDR block more restrictive.
  ingress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  ingress_with_ipv6_cidr_blocks = [
    {
      rule             = "all-all"
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

resource "aws_lb" "ecs" {
  load_balancer_type = "application"
  security_groups    = [module.ecs_alb_sg.security_group_id]
  subnets            = module.vpc.public_subnets
  ip_address_type    = "dualstack-without-public-ipv4"
}

resource "aws_lb_target_group" "ecs" {
  port        = 80
  protocol    = "HTTP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"

  health_check {
    path = "/health/liveness"
  }
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

resource "aws_launch_template" "ecs_lt" {
  name_prefix            = "brevity-ecs-lt-"
  image_id               = data.aws_ami.amz_linux_2023.id
  instance_type          = "t2.micro"
  vpc_security_group_ids = [module.ecs_ec2_sg.security_group_id]

  instance_market_options {
    market_type = "spot"
  }

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
  desired_capacity    = 2
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

resource "aws_ecs_capacity_provider" "default" {
  name = "brevity-ecs-asg-capacity-provider"

  auto_scaling_group_provider {
    auto_scaling_group_arn = aws_autoscaling_group.ecs_asg.arn
  }
}

resource "aws_ecs_cluster_capacity_providers" "default" {
  cluster_name       = aws_ecs_cluster.default.name
  capacity_providers = [aws_ecs_capacity_provider.default.name]

  default_capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.default.name
    weight            = 1
  }
}

module "ecs_task_execution_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"
  version = "~> 5.0"

  create_role      = true
  role_name_prefix = "brevity-ecs-execution-"

  create_custom_role_trust_policy = true
  custom_role_trust_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  custom_role_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  ]
}

module "ecs_service_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  vpc_id          = module.vpc.vpc_id
  name            = "brevity-ecs-service-sg"
  use_name_prefix = true

  // TODO: Make the CIDR block more restrictive.
  ingress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  ingress_with_ipv6_cidr_blocks = [
    {
      rule             = "all-all"
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

resource "aws_ecs_task_definition" "backend" {
  family             = "brevity"
  network_mode       = "awsvpc"
  execution_role_arn = module.ecs_task_execution_role.iam_role_arn

  // NOTE: The task definition will be updated using CI/CD, so we
  // don't want terraform to override the latest version.
  lifecycle {
    ignore_changes = [
      container_definitions
    ]
  }
  track_latest = true

  container_definitions = jsonencode([
    {
      name      = "nginx"
      image     = "nginx:latest"
      cpu       = 128
      memory    = 256
      essential = true
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
          protocol      = "tcp"
        }
      ]
    }
  ])
}

resource "aws_ecs_service" "backend" {
  name            = "brevity-backend-service"
  cluster         = aws_ecs_cluster.default.id
  task_definition = aws_ecs_task_definition.backend.arn
  desired_count   = 1
  // TODO: Remove when the application is running stably on production.
  force_delete = true

  network_configuration {
    subnets         = module.vpc.public_subnets
    security_groups = [module.ecs_service_sg.security_group_id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.ecs.arn
    container_name   = "nginx"
    container_port   = 80
  }

  capacity_provider_strategy {
    base              = 0
    capacity_provider = aws_ecs_capacity_provider.default.name
    weight            = 1
  }

  deployment_circuit_breaker {
    enable   = false
    rollback = false
  }

  deployment_controller {
    type = "ECS"
  }
}

output "ecr" {
  value = {
    name = aws_ecr_repository.default.name
    url  = aws_ecr_repository.default.repository_url
  }
}

output "alb" {
  value = {
    url = aws_lb.ecs.dns_name
  }
  description = "The main load balancer for the backend server instances."
}
