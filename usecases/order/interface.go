package order

import "github.com/Vractos/kloni/entity"

type UseCase interface {
	// ProcessWebhook handles incoming order webhooks from Mercado Livre.
	// It posts the order notification to a queue for asynchronous processing.
	//
	// Parameters:
	//   - input: OrderWebhookDtoInput containing the webhook data
	//
	// Returns:
	//   - error: ErrPostingOrderNotification if there's an error posting to the queue, nil otherwise
	ProcessWebhook(input OrderWebhookDtoInput) error
	// ProcessOrder handles the complete order processing workflow.
	// It validates the order, retrieves necessary credentials, fetches order data,
	// and synchronizes quantities across cloned items.
	// It also registers the order in the database and caches it.
	//
	// Parameters:
	//   - order: OrderMessage containing the order details to process
	//
	// Returns:
	//   - error: Various error types depending on the failure point, nil on success
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
	GetOrder(orderMarketplaceId string) (*entity.Order, error)
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
	SetOrder(o *entity.Order) error
}

type CacheReader interface {
	GetOrder(orderId string) (*entity.OrderStatus, error)
}

type Cache interface {
	CacheWriter
	CacheReader
}
