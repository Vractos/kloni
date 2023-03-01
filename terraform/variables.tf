variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "sdk_user_name" {
  type = string
  nullable = false
}

variable "db_username" {
  type = string
  nullable = false
  sensitive = true
}

variable "db_password" {
  type = string
  nullable = false
  sensitive = true
}

variable "public_ip" {
  type = list
  sensitive = true
}

variable "ssh_public_key" {
  sensitive = true
  type = string
}

variable "openvpn_template_url" {
  sensitive = true
  type = string
}