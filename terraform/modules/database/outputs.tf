output "database_endpoint" {
  description = "The endpoint of the database"
  sensitive = true
  value = aws_db_instance.postgres_db.address
}

output "database_port" {
  description = "The port of the database"
  value = aws_db_instance.postgres_db.port
}