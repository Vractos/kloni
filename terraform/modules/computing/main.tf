resource "aws_key_pair" "default_key" {
  key_name = "dafault-key"
  public_key = var.ami_default_public_key

  tags = {
    "Description" = "Created to use on the Dolly Project"
  }
}

resource "aws_instance" "dolly_server" {
  ami = var.ami
  instance_type = "t2.micro"
  key_name = aws_key_pair.default_key.key_name
  monitoring = true
  subnet_id = var.server_public_subnet
  vpc_security_group_ids = [var.server_security_group]

  
  tags = {
    Name = "${var.instance_name} Server"
    Environment = var.environment
  }
}

resource "aws_eip_association" "eip_assoc" {
  instance_id = aws_instance.dolly_server.id
  allocation_id = var.eip_id
}