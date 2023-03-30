resource "aws_elasticache_subnet_group" "redis_subnet" {
  name = "dolly-redis-subnet"
  subnet_ids = [var.subnet_id]
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id = "dolly-redis"
  engine = "redis"
  node_type = "cache.t3.micro"
  num_cache_nodes = 1
  parameter_group_name = "default.redis6.x"
  engine_version = "6.2"
  subnet_group_name = aws_elasticache_subnet_group.redis_subnet.name
  security_group_ids = var.redis_security_group_ids
}