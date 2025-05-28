terraform {
  cloud {
    organization = "andy-learn-terraform"
    workspaces {
      name = "brevity-go"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
