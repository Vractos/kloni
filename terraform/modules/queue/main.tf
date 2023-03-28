resource "aws_sqs_queue" "sqs_fifo_queue_deadletter" {
  name                       = "${var.sqs_queue_name}Deadletter.fifo"
  # name_prefix                       = "${var.sqs_queue_name}Deadetter"
  fifo_queue = true
  visibility_timeout_seconds = 60
  receive_wait_time_seconds  = 20
  sqs_managed_sse_enabled = false

  tags = {
    Name        = "${var.project}-${var.sqs_queue_name}-dll-sqs-fifo"
    Environment = var.environment
  }
}

resource "aws_sqs_queue_redrive_allow_policy" "fifo-queue-dll" {
  queue_url = aws_sqs_queue.sqs_fifo_queue_deadletter.url

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue"
    sourceQueueArns = [aws_sqs_queue.sqs_fifo_queue.arn]
  })

  depends_on = [
    aws_sqs_queue.sqs_fifo_queue_deadletter,
    aws_sqs_queue.sqs_fifo_queue
  ]
}

resource "aws_sqs_queue" "sqs_fifo_queue" {
  name                       = "${var.sqs_queue_name}.fifo"
  # name_prefix                       = "${var.sqs_queue_name}"
  fifo_queue = true
  visibility_timeout_seconds = var.fifo_queue_visibility_timeout
  receive_wait_time_seconds  = 20
  sqs_managed_sse_enabled = false

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.sqs_fifo_queue_deadletter.arn
    maxReceiveCount = 7
  })
  depends_on = [
    aws_sqs_queue.sqs_fifo_queue_deadletter
  ]

  tags = {
    Name        = "${var.project}-${var.sqs_queue_name}-sqs-fifo"
    Environment = var.environment
  }
}

resource "aws_sqs_queue_policy" "default" {
  queue_url = aws_sqs_queue.sqs_fifo_queue.url

  policy = <<POLICY
   {
   "Version": "2012-10-17",
   "Id": "sqs-default-policy",
   "Statement": [{
      "Sid":"Default_Send_Receive",
      "Effect": "Allow",
      "Principal": {
         "AWS": "${var.sqs_queue_allowed_user}"
      },
      "Action": [
         "sqs:SendMessage",
         "sqs:ReceiveMessage",
         "sqs:DeleteMessage"
      ],
      "Resource": "${aws_sqs_queue.sqs_fifo_queue.arn}"
   }]
}
  POLICY
}

resource "aws_sqs_queue_policy" "default-dll" {
  queue_url = aws_sqs_queue.sqs_fifo_queue_deadletter.url

  policy = <<POLICY
   {
   "Version": "2012-10-17",
   "Id": "sqs-dll-default-policy",
   "Statement": [{
      "Sid":"Default_Send_Receive",
      "Effect": "Allow",
      "Principal": {
         "AWS": "${var.sqs_queue_allowed_user}"
      },
      "Action": [
         "sqs:SendMessage",
         "sqs:ReceiveMessage",
         "sqs:DeleteMessage"
      ],
      "Resource": "${aws_sqs_queue.sqs_fifo_queue_deadletter.arn}"
   }]
}
  POLICY
}