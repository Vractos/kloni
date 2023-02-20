output "sqs_fifo_queue_url" {
  value = module.queue.sqs_fifo_queue_url
}

output "database_endpoint" {
  description = "The endpoint of the database"
  sensitive = true
  value = module.database.database_endpoint
}

output "database_port" {
  description = "The port of the database"
  value = module.database.database_port
}

output "server_public_ip" {
  value = module.network.server_public_ip
  sensitive = true
}