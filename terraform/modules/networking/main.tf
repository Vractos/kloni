resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "${var.vpc_name}-vpc"
    Environment = var.environment
  }
}

resource "aws_internet_gateway" "main_vpc_internet_gateway" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.internet_gateway_name}-igw"
    Environment = var.environment
  }
}

resource "aws_eip" "elastic_ip" {
  vpc = true
  depends_on = [aws_internet_gateway.main_vpc_internet_gateway]
  tags = {
    Name = "${var.elastic_ip_name}-eip"
    Environment = var.environment
  }
}

resource "aws_subnet" "public" {
  vpc_id = aws_vpc.main.id

  cidr_block = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.aws_az.names[0]
  map_public_ip_on_launch = true
  
  tags = {
    Name = "${var.public_subnet_name}-public-subnet"
    Environment = var.environment
  }
}


resource "aws_subnet" "private" {
  count = 2
  vpc_id = aws_vpc.main.id

  cidr_block = "10.0.${count.index+2}.0/24"
  availability_zone = data.aws_availability_zones.aws_az.names[count.index+1]
  map_public_ip_on_launch = false
  
  tags = {
    Name = "${var.private_subnet_name}${count.index+1}-private-subnet"
    Environment = var.environment
  }
}

resource "aws_route_table" "public_subnet_route_table" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main_vpc_internet_gateway.id
  }

  tags = {
    Name = "${var.public_subnet_route_table_name}-public-subnet-route-table"
    Environment = var.environment
  }
}

resource "aws_route_table" "private_subnet_route_table" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.private_subnet_route_table_name}-private-subnet-route-table"
    Environment = var.environment
  }
}

resource "aws_route_table_association" "public_subnet_route_table_association" {
  subnet_id = aws_subnet.public.id
  route_table_id = aws_route_table.public_subnet_route_table.id
}

resource "aws_route_table_association" "private_subnet_route_table_association" {
  count = 2

  subnet_id = element(aws_subnet.private[*].id, count.index)
  route_table_id = aws_route_table.private_subnet_route_table.id
}
resource "aws_security_group" "server_security_group" {
  name        = "server_security_group"
  description = "Allow traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "Allow all traffic through HTTP"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    description = "Allow all traffic through HTTPS (TLS)"
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow SSH from my computer"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = var.my_public_ip
  }

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "server_security_group"
  }
}

resource "aws_security_group" "database_security_group" {
  name        = "database_security_group"
  description = "Allow Datababse traffic"
  vpc_id      = aws_vpc.main.id

    ingress {
    description     = "Allow Postgres traffic from only the server sg"
    from_port       = "5432"
    to_port         = "5432"
    protocol        = "tcp"
    security_groups = [aws_security_group.server_security_group.id]
  }

  tags = {
    Name = "database_security_group"
  } 
}
