variable "subnet_id" {
  type = string
  nullable = false
}

variable "redis_security_group_ids" {
  type = set(string)
  nullable = false
}

variable "node_type" {
  type = string
  default = "cache.t3.micro"
}

variable "num_cache_nodes" {
  type = number
  default = 1
}

variable "parameter_group_name" {
  type = string
  default = "default.redis6.x"
}

variable "engine_version" {
  type = string
  default = "6.2"
}

variable "project" {
  type = string
  default = "dolly"
}