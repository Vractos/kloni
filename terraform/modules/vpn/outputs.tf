output "client_vpn_endpoint_id" {
  description = "The ID of the Client VPN endpoint"
  value       = aws_ec2_client_vpn_endpoint.vpn_endpoint.id
}

output "client_vpn_endpoint_arn" {
  description = "The ARN of the Client VPN endpoint"
  value       = aws_ec2_client_vpn_endpoint.vpn_endpoint.arn
}

output "client_vpn_endpoint_dns_name" {
  description = "The DNS name to be used by clients when establishing VPN sessions"
  value       = aws_ec2_client_vpn_endpoint.vpn_endpoint.dns_name
}
