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

// type OrderMessage {

// }

// Queue producer interface
type QueueProducer interface {
	PostOrderNotification(input OrderWebhookDtoInput) error
}

type QueueConsumer interface {
	ConsumeOrderNotification()
}

type Queue interface {
	QueueProducer
	QueueConsumer
}
