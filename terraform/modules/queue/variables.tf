variable "sqs_queue_name" {
  type = string
  nullable = false
}

variable "fifo_queue_visibility_timeout" {
  type = number
  default = 30
}

variable "sqs_queue_allowed_user" {
  type = string
  nullable = false
}

variable "project" {
  type = string
  default = "kloni"
}