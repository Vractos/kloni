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

module "database" {
  depends_on = [
    module.network
  ]
  source = "./modules/database"
  environment = var.environment
  subnet_ids = module.network.private_subnets
  db_username = var.db_username
  db_password = var.db_password
  db_security_group_ids = module.network.database_security_group_id
}

module "queue" {
  source = "./modules/queue"
  environment = var.environment
  sqs_queue_name = "orders"
  sqs_queue_allowed_user = data.aws_iam_user.sdk_user.arn
}