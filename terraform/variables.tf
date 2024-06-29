variable "project" {
  type    = string
  default = "kloni"
}

variable "region" {
  type    = string
  default = "us-east-1"
}

variable "sdk_user_name" {
  type     = string
  nullable = false
}

variable "db_username" {
  type      = string
  nullable  = false
  sensitive = true
}

variable "db_password" {
  type      = string
  nullable  = false
  sensitive = true
}

variable "public_ip" {
  type      = list(any)
  sensitive = true
}

variable "ssh_public_key" {
  sensitive = true
  type      = string
}

variable "vpn_server_certificate_arn" {
  type     = string
  nullable = true
}

variable "vpn_client_root_certificate_chain_arn" {
  type     = string
  nullable = true
}

variable "tailscale_ami" {
  description = "AMI ID for the tailscale subnet router"
  type        = string
}

variable "tailscale_authkey" {
  description = "tailscale authentication key"
  type        = string
}

variable "tailscale_cidr" {
  description = "tailscale CIDR block"
  type        = string
  default     = "100.64.0.0/10"
}
