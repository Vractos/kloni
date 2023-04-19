terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.31.0"
    }
  }
}

provider "aws" {
  region = var.region
  default_tags {
    tags = {
      Environment = "${terraform.workspace}"
      Project = var.project
    }
  }
}

module "network" {
  count = (terraform.workspace != "dev") ? 1 : 0
  source       = "./modules/networking"
  project = var.project
  my_public_ip = var.public_ip
}

module "database" {
  count = (terraform.workspace != "dev") ? 1 : 0
  depends_on = [
    module.network
  ]
  project = var.project
  source                = "./modules/database"
  subnet_ids            = module.network.private_subnets
  db_username           = var.db_username
  db_password           = var.db_password
  db_security_group_ids = module.network.database_security_group_id
}

module "queue" {
  source                 = "./modules/queue"
  project = var.project
  sqs_queue_name         = "orders"
  sqs_queue_allowed_user = data.aws_iam_user.sdk_user.arn
}

module "computing" {
  project = var.project
  count = (terraform.workspace != "dev") ? 1 : 0
  sdk_username = data.aws_iam_user.sdk_user.user_name
  source                 = "./modules/computing"
  eip_id                 = module.network.eip_id
  ami_default_public_key = var.ssh_public_key
  server_public_subnet   = module.network.public_subnet
  server_security_group  = module.network.server_security_group_id

}

module "cache" {
  count = (terraform.workspace != "dev") ? 1 : 0
  project = var.project
  depends_on = [
    module.network
  ]
  source                   = "./modules/cache"
  subnet_id                = module.network.private_subnets[0]
  redis_security_group_ids = module.network.redis_security_group_id
}

resource "aws_cloudformation_stack" "network" {
  count = (terraform.workspace == "prod") ? 1 : 0
  name         = "dolly-vpn"
  template_url = var.openvpn_template_url
  capabilities = ["CAPABILITY_IAM"]

  parameters = {
    KeyName      = "dafault-key",
    VpcId        = module.network.vpc_id
    SubnetId     = module.network.public_subnet
    InstanceType = "t2.micro"
    InstanceName = "Dolly VPN"
  }

}
