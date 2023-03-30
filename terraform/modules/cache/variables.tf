variable "subnet_id" {
  type = string
  nullable = false
}

variable "redis_security_group_ids" {
  type = set(string)
  nullable = false
}