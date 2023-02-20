variable "environment" {
  type = string
}

variable "vpc_name" {
  type = string
  default = "dolly"
}

variable "internet_gateway_name" {
  type = string
  default = "dolly"
}

variable "elastic_ip_name" {
  type = string
  default = "dolly"
}

variable "public_subnet_name" {
  type = string
  default = "dolly"
}

variable "private_subnet_name" {
  type = string
  default = "dolly"
}

variable "public_subnet_route_table_name" {
  type = string
  default = "dolly"
}

variable "private_subnet_route_table_name" {
  type = string
  default = "dolly"
}

variable "my_public_ip" {
  description = "Your public IP address"
  type = string
  sensitive = true
}