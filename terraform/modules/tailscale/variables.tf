variable "project" {
  description = "Project name"
  type        = string
}

variable "ami" {
  description = "AMI ID for the tailscale subnet router"
  type        = string
  default     = "ami-0557a15b87f6559cf"
}

variable "instance_type" {
  description = "Instance type for the tailscale subnet router"
  type        = string
  default     = "t3.micro"
}

variable "key_name" {
  description = "Key pair name for SSH access"
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID where the tailscale subnet router will be deployed"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where the tailscale subnet router will be deployed"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block of the VPC"
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

variable "private_route_table_id" {
  description = "ID of the private route table"
  type        = string
}
