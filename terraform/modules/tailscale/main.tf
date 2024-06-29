resource "aws_instance" "tailscale_subnet_router" {
  ami                    = var.ami
  instance_type          = "t3.micro"
  key_name               = var.key_name
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [aws_security_group.tailscale_sg.id]

  user_data = <<-EOF
              #!/bin/bash
              curl -fsSL https://tailscale.com/install.sh | sh
              tailscale up --authkey=${var.tailscale_authkey} --advertise-routes=${var.vpc_cidr}
              EOF

  tags = {
    Name = "${var.project}-tailscale-subnet-router"
  }
}

resource "aws_security_group" "tailscale_sg" {
  name        = "${var.project}-tailscale-security-group"
  description = "Security group for tailscale subnet router"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 41641
    to_port     = 41641
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-tailscale-sg"
  }
}

resource "aws_route" "tailscale_route" {
  route_table_id         = var.private_route_table_id
  destination_cidr_block = var.tailscale_cidr
  instance_id            = aws_instance.tailscale_subnet_router.id
}
