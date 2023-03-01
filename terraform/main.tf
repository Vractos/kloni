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
  my_public_ip = var.public_ip
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

module "computing" {
  source = "./modules/computing"
  environment = var.environment
  eip_id = module.network.eip_id
  ami_default_public_key = var.ssh_public_key
  server_public_subnet = module.network.public_subnet
  server_security_group = module.network.server_security_group_id
  
}

resource "aws_cloudformation_stack" "network" {
  name = "dolly-vpn"
  template_url = var.openvpn_template_url
  capabilities = ["CAPABILITY_IAM"]

  parameters = {
    KeyName = "dafault-key",
    VpcId = module.network.vpc_id
    SubnetId = module.network.public_subnet
    InstanceType = "t2.micro"
    InstanceName = "Dolly VPN"
  }

}