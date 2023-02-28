variable "ami" {
  description = "Ubuntu 22.04 LTS"
  default = "ami-0557a15b87f6559cf"
}

variable "instance_name" {
  type = string
  default = "dolly"
}

variable "ami_default_public_key" {
  description = "A AMI Key pair"
}

variable "environment" {
  type = string
}

variable "eip_id" {
  description = "The Elastic IP to associate with the instance"
}

variable "server_public_subnet" {
  nullable = false
  type = string
}

variable "server_security_group" {
  nullable = false
  type = string
}