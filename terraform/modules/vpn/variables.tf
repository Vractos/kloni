variable "project" {
  description = "The name of the project"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs to associate with the Client VPN endpoint"
  type        = list(string)
}

variable "db_security_group_id" {
  description = "The ID of the database security group"
  type        = string
}

variable "client_cidr_block" {
  description = "The IPv4 CIDR block for the Client VPN endpoint"
  type        = string
  default     = "10.100.0.0/16"
}

variable "server_certificate_arn" {
  description = "The ARN of the ACM certificate for the VPN server"
  type        = string
}

variable "client_root_certificate_chain_arn" {
  description = "The ARN of the client root certificate chain"
  type        = string
}
