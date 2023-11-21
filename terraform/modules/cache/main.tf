resource "aws_elasticache_subnet_group" "redis_subnet" {
  name = "kloni-redis-subnet"
  subnet_ids = [var.subnet_id]
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id = "kloni-redis"
  engine = "redis"
  node_type = var.node_type
  num_cache_nodes = var.num_cache_nodes
  parameter_group_name = var.parameter_group_name
  engine_version = var.engine_version
  subnet_group_name = aws_elasticache_subnet_group.redis_subnet.name
  security_group_ids = var.redis_security_group_ids
}