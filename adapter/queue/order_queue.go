package queue

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/Vractos/dolly/usecases/order"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type OrderSQSQueue struct {
	Client *sqs.Client
	Url    string
}

func NewOrderQueue(client *sqs.Client, url string) *OrderSQSQueue {
	return &OrderSQSQueue{Client: client, Url: url}
}

// PostOrderNotification implements order.Queue
func (q *OrderSQSQueue) PostOrderNotification(input order.OrderWebhookDtoInput) error {
	msgBody, err := json.Marshal(input)
	if err != nil {
		log.Printf("Failed to marshal input into json: %v", err)
		return err
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
		QueueUrl:               &q.Url,
		MessageBody:            aws.String(string(msgBody)),
		MessageDeduplicationId: aws.String(input.ID),
		MessageGroupId:         aws.String("order-notification"),
	}

	resp, err := q.Client.SendMessage(context.TODO(), mgsInput)
	if err != nil {
		log.Println("Got an error sending the order message:")
		log.Panicln(err)
		return err
	}

	log.Printf("Sent message with ID: %s", *resp.MessageId)
	return nil
}

// ReadOrderNotification implements order.Queue
func (q *OrderSQSQueue) ConsumeOrderNotification() []order.OrderMessage {
	getMsgInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            &q.Url,
		WaitTimeSeconds:     int32(20),
		MaxNumberOfMessages: 10,
	}

	resp, err := q.Client.ReceiveMessage(context.TODO(), getMsgInput)
	if err != nil {
		log.Println("Got an error receiving the order message")
		log.Panicln(err)
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
		QueueUrl:      &q.Url,
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := q.Client.DeleteMessage(context.TODO(), dMInput)
	if err != nil {
		log.Println("Got an error deleting the order message:")
		log.Println(err)
		return err
	}

	return nil
}
