terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "4.31.0"
    }
  }
}

provider "aws" {
  region = var.region
}

module "network" {
  source = "./modules/networking"
  environment = var.environment
}

module "queue" {
  source = "./modules/queue"
  environment = var.environment
  sqs_queue_name = "orders"
  sqs_queue_allowed_user = data.aws_iam_user.sdk_user.id
}