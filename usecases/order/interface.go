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

// Queue producer interface
type QueueProducer interface {
	PostOrderNotification(input OrderWebhookDtoInput) error
}

type Queue interface {
	QueueProducer
}
