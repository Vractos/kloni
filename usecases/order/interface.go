package order

type UseCase interface {
	ProcessWebhook(input OrderWebhookDtoInput) error
}

/*
#########################################
#########################################
------------------QUEUE------------------
#########################################
#########################################
*/

type OrderMessage struct {
	Store    string
	OrderId  string
	Attempts int
}

// Queue producer interface
type QueueProducer interface {
	PostOrderNotification(input OrderWebhookDtoInput) error
}

type QueueConsumer interface {
	ConsumeOrderNotification() []OrderMessage
}

type Queue interface {
	QueueProducer
	QueueConsumer
}
