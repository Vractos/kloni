variable "vpc_name" {
  type = string
  default = "dolly"
}

variable "my_public_ip" {
  description = "Your public IP address"
  type = list
  sensitive = true
}

variable "project" {
  type = string
  default = "dolly"
}