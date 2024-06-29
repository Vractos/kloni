resource "aws_ec2_client_vpn_endpoint" "vpn_endpoint" {
  description            = "Client VPN endpoint for database access"
  server_certificate_arn = var.server_certificate_arn
  client_cidr_block      = var.client_cidr_block

  authentication_options {
    type                       = "certificate-authentication"
    root_certificate_chain_arn = var.client_root_certificate_chain_arn
  }

  connection_log_options {
    enabled = false
  }

  tags = {
    Name = "${var.project}-client-vpn-endpoint"
  }
}

resource "aws_ec2_client_vpn_network_association" "vpn_subnet_association" {
  count                  = length(var.subnet_ids)
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.vpn_endpoint.id
  subnet_id              = var.subnet_ids[count.index]
}

resource "aws_ec2_client_vpn_authorization_rule" "vpn_auth_rule" {
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.vpn_endpoint.id
  target_network_cidr    = var.vpc_cidr
  authorize_all_groups   = true
}

resource "aws_security_group_rule" "allow_vpn_to_db" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  security_group_id        = var.db_security_group_id
  source_security_group_id = aws_ec2_client_vpn_endpoint.vpn_endpoint.security_group_id
}
