output "sqs_fifo_queue_url" {
  value = aws_sqs_queue.sqs_fifo_queue.url
}