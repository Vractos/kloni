output "tailscale_subnet_router_id" {
  description = "ID of the tailscale subnet router EC2 instance"
  value       = aws_instance.tailscale_subnet_router.id
}

output "tailscale_subnet_router_private_ip" {
  description = "Private IP of the tailscale subnet router EC2 instance"
  value       = aws_instance.tailscale_subnet_router.private_ip
}

output "tailscale_security_group_id" {
  description = "ID of the tailscale security group"
  value       = aws_security_group.tailscale_sg.id
}
