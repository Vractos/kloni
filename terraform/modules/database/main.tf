resource "aws_db_subnet_group" "default" {
  name = var.db_name
  subnet_ids = var.subnet_ids

  tags = {
    Name = "${var.db_name} DB subnet group"
  }
}

resource "aws_db_instance" "postgres_db" {
  identifier = var.db_name
  allocated_storage = 10
  db_name = var.db_name
  engine = "postgres"
  engine_version = "14.5"
  instance_class = "db.t3.micro"
  username = var.db_username
  password = var.db_password
  backup_retention_period = 5
  backup_window = "00:00-01:00"
  copy_tags_to_snapshot = true
  storage_encrypted = true
  performance_insights_enabled = true
  enabled_cloudwatch_logs_exports = ["postgresql"]
  maintenance_window = "Sun:01:00-Sun:02:00"
  publicly_accessible = false
  final_snapshot_identifier = var.final_snapshot_identifier
  db_subnet_group_name = aws_db_subnet_group.default.name
  vpc_security_group_ids = var.db_security_group_ids
  

  tags = {
    Name = "${var.project}-postgres-db"
  }
  
}