package order

import (
	"errors"
	"log"
)

type OrderService struct {
	queue Queue
}

func NewOrderService(queue Queue) *OrderService {
	return &OrderService{queue: queue}
}

func (o *OrderService) ProcessWebhook(input OrderWebhookDtoInput) error {
	if err := o.queue.PostOrderNotification(input); err != nil {
		log.Println(err.Error())
		return errors.New("error to post order notification")
	}
	return nil
}
