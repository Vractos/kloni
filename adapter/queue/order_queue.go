package queue

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"

	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/order"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.uber.org/zap"
)

type OrderSQSQueue struct {
	client *sqs.Client
	url    string
	logger metrics.Logger
}

func NewOrderQueue(client *sqs.Client, url string, logger metrics.Logger) *OrderSQSQueue {
	return &OrderSQSQueue{
		client: client,
		url:    url,
		logger: logger,
	}
}

// PostOrderNotification implements order.Queue
func (q *OrderSQSQueue) PostOrderNotification(input order.OrderWebhookDtoInput) error {
	msgBody, err := json.Marshal(input)
	if err != nil {
		q.logger.Error(
			"Failed to marshal input into json",
			err,
			zap.String("notification_id", input.ID),
		)
		return errors.New("failed to handle input")
	}

	mgsInput := &sqs.SendMessageInput{
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Store": {
				DataType:    aws.String("Number.int"),
				StringValue: aws.String(strconv.Itoa(input.UserID)),
			},
			"ResourcePath": {
				DataType:    aws.String("String"),
				StringValue: aws.String(input.Resource),
			},
			"Attempts": {
				DataType:    aws.String("Number.int"),
				StringValue: aws.String(strconv.Itoa(input.Attempts)),
			},
			"Send": {
				DataType:    aws.String("String"),
				StringValue: aws.String(input.Sent),
			},
			"Received": {
				DataType:    aws.String("String"),
				StringValue: aws.String(input.Received),
			},
		},
		QueueUrl:               &q.url,
		MessageBody:            aws.String(string(msgBody)),
		MessageDeduplicationId: aws.String(input.ID),
		MessageGroupId:         aws.String("order-notification"),
	}

	resp, err := q.client.SendMessage(context.TODO(), mgsInput)
	if err != nil {
		q.logger.Error(
			"Failure to send the order message",
			err,
			zap.String("notification_id", input.ID),
		)
		return err
	}

	q.logger.Info("Sent message", zap.String("message_id", *resp.MessageId))
	return nil
}

// ReadOrderNotification implements order.Queue
func (q *OrderSQSQueue) ConsumeOrderNotification() []order.OrderMessage {
	getMsgInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            &q.url,
		WaitTimeSeconds:     int32(20),
		MaxNumberOfMessages: 10,
	}

	resp, err := q.client.ReceiveMessage(context.TODO(), getMsgInput)
	if err != nil {
		q.logger.Error(
			"Got an error receiving the order message",
			err,
		)
		return nil
	} else if resp.Messages == nil {
		return nil
	}

	orderMessages := make([]order.OrderMessage, len(resp.Messages))
	for i, e := range resp.Messages {
		orderMessages[i].Store = *e.MessageAttributes["Store"].StringValue
		orderMessages[i].OrderId = regexp.MustCompile(`\w+$`).FindString(*e.MessageAttributes["ResourcePath"].StringValue)
		orderMessages[i].Attempts, _ = strconv.Atoi(*e.MessageAttributes["Attempts"].StringValue)
		orderMessages[i].ReceiptHandle = *e.ReceiptHandle
	}

	return orderMessages
}

// DeleteOrderNotification implements order.Queue
func (q *OrderSQSQueue) DeleteOrderNotification(receiptHandle string) error {

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      &q.url,
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := q.client.DeleteMessage(context.TODO(), dMInput)
	if err != nil {
		q.logger.Error(
			"Got an error deleting the order message",
			err,
			zap.String("receipt_handle", receiptHandle),
		)
		return err
	}

	return nil
}
