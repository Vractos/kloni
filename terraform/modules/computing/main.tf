resource "aws_key_pair" "default_key" {
  key_name   = "dafault-key"
  public_key = var.ami_default_public_key

  tags = {
    "Description" = "Created to use on the Dolly Project"
  }
}

resource "aws_iam_role" "ec2_sqs_role" {
  name = "ec2_access_sqs"
  description = "Allows EC2 instances to call SQS"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "sts:AssumeRole"
        ],
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "sqs_access_policy" {
  name = "sqs_access_policy"
  role = aws_iam_role.ec2_sqs_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
        {
            Effect = "Allow",
            Action = [
                "sqs:DeleteMessage",
                "sqs:GetQueueUrl",
                "sqs:ChangeMessageVisibility",
                "sqs:ReceiveMessage",
                "sqs:SendMessage"
            ],
            Resource = "arn:aws:sqs:us-east-1:068065400013:*"
        }
    ]
  })
}

resource "aws_iam_instance_profile" "dolly_instance_profile" {
  name = "dolly_server_profile"
  role = aws_iam_role.ec2_sqs_role.name
}

resource "aws_instance" "dolly_server" {
  ami                    = var.ami
  instance_type          = "t2.micro"
  key_name               = aws_key_pair.default_key.key_name
  monitoring             = true
  subnet_id              = var.server_public_subnet
  vpc_security_group_ids = [var.server_security_group]

  iam_instance_profile = aws_iam_instance_profile.dolly_instance_profile.name


  tags = {
    Name        = "${var.instance_name} Server"
    Environment = var.environment
  }
}

resource "aws_eip_association" "eip_assoc" {
  instance_id   = aws_instance.dolly_server.id
  allocation_id = var.eip_id
}