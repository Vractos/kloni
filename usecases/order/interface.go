package order

import "github.com/Vractos/dolly/entity"

type UseCase interface {
	ProcessWebhook(input OrderWebhookDtoInput) error
	ProcessOrder(order OrderMessage) error
}

/*
#########################################
#########################################
------------------QUEUE------------------
#########################################
#########################################
*/

type OrderMessage struct {
	Store         string
	OrderId       string
	Attempts      int
	ReceiptHandle string
}

// Queue producer interface
type QueueProducer interface {
	PostOrderNotification(input OrderWebhookDtoInput) error
}

type QueueConsumer interface {
	ConsumeOrderNotification() []OrderMessage
	DeleteOrderNotification(receiptHandle string) error
}

type Queue interface {
	QueueProducer
	QueueConsumer
}

/*
#########################################
#########################################
---------------REPOSITORY---------------
#########################################
#########################################
*/

type RepoWriter interface {
	RegisterOrder(o *entity.Order) error
}

type RepoReader interface {
}

type Repository interface {
	RepoWriter
	RepoReader
}

/*
#########################################
#########################################
------------------CACHE------------------
#########################################
#########################################
*/

type CacheWriter interface {
	SetOrder(orderId string) error
}

type CacheReader interface {
}

type Cache interface {
	CacheWriter
	CacheReader
}
