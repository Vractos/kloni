output "vpc_id" {
  value = aws_vpc.main.id
}

output "eip_id" {
  value = aws_eip.elastic_ip.id
}

output "server_public_ip" {
  value = aws_eip.elastic_ip.public_ip
  sensitive = true
}

output "private_subnets" {
  description = "List of IDs of private subnets"
  value       = aws_subnet.private[*].id
}
output "public_subnet" {
  description = "ID of public subnets"
  value       = aws_subnet.public.id
}

output "database_security_group_id" {
  description = "Database security group ID"
  value = aws_security_group.database_security_group[*].id
}

output "redis_security_group_id" {
  description = "Redis security group ID"
  value = aws_security_group.redis_security_group[*].id
}

output "server_security_group_id" {
  description = "Server security group ID"
  value = aws_security_group.server_security_group.id
}