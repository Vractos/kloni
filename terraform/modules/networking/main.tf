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
    Name = "${var.public_subnet_name}-subnet"
    Environment = var.environment
  }
}


resource "aws_subnet" "private" {
  vpc_id = aws_vpc.main.id

  cidr_block = "10.0.2.0/24"
  availability_zone = data.aws_availability_zones.aws_az.names[1]
  map_public_ip_on_launch = false
  
  tags = {
    Name = "${var.private_subnet_name}-subnet"
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
  subnet_id = aws_subnet.private.id
  route_table_id = aws_route_table.private_subnet_route_table.id
}