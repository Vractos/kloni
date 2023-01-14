output "vpc_id" {
  value = aws_vpc.main.id
}

output "eip_id" {
  value = aws_eip.elastic_ip.id
}

output "private_subnets" {
  description = "List of IDs of private subnets"
  value       = aws_subnet.private[*].id
}

output "database_security_group_id" {
  description = "Database security group ID"
  value = aws_security_group.database_security_group[*].id
}