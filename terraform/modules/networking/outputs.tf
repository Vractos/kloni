output "vpc_id" {
  value = aws_vpc.main.id
}

output "eip_id" {
  value = aws_eip.elastic_ip.id
}