output "sqs_fifo_queue_url" {
  value = module.queue.sqs_fifo_queue_url
}

output "database_endpoint" {
  description = "The endpoint of the database"
  sensitive = true
  value = length(module.database) > 0 ? module.database[0].database_endpoint : null
}
output "server_public_ip" {
  value = length(module.network) > 0 ? module.network[0].server_public_ip : null
  sensitive = true
}