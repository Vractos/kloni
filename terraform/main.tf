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
      Project     = var.project
    }
  }
}

module "network" {
  source       = "./modules/networking"
  count        = (terraform.workspace != "dev") ? 1 : 0
  project      = var.project
  my_public_ip = var.public_ip
}

moved {
  from = module.network
  to   = module.network[0]
}

locals {
  network_outputs = length(module.network) > 0 ? {
    vpc_id                     = module.network[0].vpc_id
    eip_id                     = module.network[0].eip_id
    server_public_ip           = module.network[0].server_public_ip
    private_subnets            = module.network[0].private_subnets
    public_subnet              = module.network[0].public_subnet
    database_security_group_id = module.network[0].database_security_group_id
    redis_security_group_id    = module.network[0].redis_security_group_id
    server_security_group_id   = module.network[0].server_security_group_id
    vpc_cidr                   = module.network[0].vpc_cidr
    private_route_table_id     = module.network[0].private_route_table_id
    vpc_enable_dns_hostnames   = module.network[0].vpc_enable_dns_hostnames
    } : {
    vpc_id                     = null
    eip_id                     = null
    server_public_ip           = null
    private_subnets            = null
    public_subnet              = null
    database_security_group_id = null
    redis_security_group_id    = null
    server_security_group_id   = null
    vpc_cidr                   = null
    private_route_table_id     = null
    vpc_enable_dns_hostnames   = null
  }
}

module "database" {
  source = "./modules/database"
  count  = (terraform.workspace != "dev") ? 1 : 0
  depends_on = [
    module.network
  ]
  project               = var.project
  subnet_ids            = local.network_outputs.private_subnets
  db_username           = var.db_username
  db_password           = var.db_password
  db_security_group_ids = local.network_outputs.database_security_group_id
}

moved {
  from = module.database
  to   = module.database[0]
}

module "queue" {
  source                 = "./modules/queue"
  project                = var.project
  sqs_queue_name         = (terraform.workspace == "prod") ? "orders" : "order-${terraform.workspace}"
  sqs_queue_allowed_user = data.aws_iam_user.sdk_user.arn
}

module "computing" {
  source                 = "./modules/computing"
  project                = var.project
  count                  = (terraform.workspace != "dev") ? 1 : 0
  sdk_account_id         = split(":", data.aws_iam_user.sdk_user.arn)[4]
  eip_id                 = local.network_outputs.eip_id
  ami_default_public_key = var.ssh_public_key
  server_public_subnet   = local.network_outputs.public_subnet
  server_security_group  = local.network_outputs.server_security_group_id
  region                 = var.region
}

moved {
  from = module.computing
  to   = module.computing[0]
}

# module "cache" {
#   source = "./modules/cache"
#   count  = (terraform.workspace != "dev") ? 1 : 0
#   # count   = (terraform.workspace != "dev") ? 0 : 1
#   project = var.project
#   depends_on = [
#     module.network
#   ]
#   subnet_id                = local.network_outputs.private_subnets[0]
#   redis_security_group_ids = local.network_outputs.redis_security_group_id
# }

# moved {
#   from = module.cache
#   to   = module.cache[0]
# }
