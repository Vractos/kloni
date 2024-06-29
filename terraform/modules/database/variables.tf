variable "subnet_ids" {
  type     = set(string)
  nullable = false
}

variable "db_name" {
  type    = string
  default = "dolly"
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

variable "final_snapshot_identifier" {
  type    = string
  default = "dolly-final-snapshot"
}

variable "db_security_group_ids" {
  type     = set(string)
  nullable = false
}

variable "project" {
  type    = string
  default = "dolly"
}
