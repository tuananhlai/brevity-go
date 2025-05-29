// https://aws.amazon.com/blogs/security/use-iam-roles-to-connect-github-actions-to-actions-in-aws/
// https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services#configuring-the-role-and-trust-policy

locals {
  github_repo = "tuananhlai/brevity-go"
}

resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com"
  ]
}

// IAM Role for Github Actions workflow to assume. It allows deployment of ECR images and ECS services.
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
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${local.github_repo}:*"
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
      // Allow Github Actions to push new Docker images to the ECR repository.
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
        Sid : "Logs",
        Effect : "Allow",
        Action : [
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams",
          "logs:GetLogEvents"
        ],
        Resource : "*"
      },
      // Allow Github Actions to create new task definitions and update ECS services for continuous deployment.
      // https://github.com/aws-actions/amazon-ecs-deploy-task-definition?tab=readme-ov-file#permissions
      {
        Sid    = "RegisterTaskDefinition",
        Effect = "Allow",
        Action = [
          "ecs:RegisterTaskDefinition"
        ],
        Resource = "*"
      },
      {
        Sid    = "PassRolesInTaskDefinition",
        Effect = "Allow",
        Action = [
          "iam:PassRole"
        ],
        Resource = [
          module.ecs_task_execution_role.iam_role_arn
        ]
      },
      {
        Sid    = "DeployService",
        Effect = "Allow",
        Action = [
          "ecs:UpdateService",
          "ecs:DescribeServices"
        ],
        Resource = [
          "${aws_ecs_service.backend.id}"
        ]
      }
    ]
  })
}

output "github_actions_iam_role" {
  value = {
    arn = aws_iam_role.github_actions.arn
  }
}
