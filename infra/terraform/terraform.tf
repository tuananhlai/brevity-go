terraform {
  backend "s3" {
    # Key needs to be different for EVERY terraform project.
    key = "global/brevity/terraform.tfstate"

    bucket         = "opentofu-remote-state-986377"
    region         = "ap-northeast-3"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
