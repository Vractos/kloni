package order

import (
	"errors"
	"log"

	"github.com/Vractos/dolly/usecases/announcement"
	"github.com/Vractos/dolly/usecases/common"
	"github.com/Vractos/dolly/usecases/store"
)

type OrderService struct {
	queue        Queue
	meli         common.MercadoLivre
	store        store.UseCase
	announcement announcement.UseCase
}

func NewOrderService(queue Queue, mercadolivre common.MercadoLivre, storeUseCase store.UseCase) *OrderService {
	return &OrderService{queue: queue, meli: mercadolivre}
}

func (o *OrderService) ProcessWebhook(input OrderWebhookDtoInput) error {
	if err := o.queue.PostOrderNotification(input); err != nil {
		log.Println(err.Error())
		return errors.New("error to post order notification")
	}
	return nil
}

func (o *OrderService) ProcessOrder(order OrderMessage) error {
	//TODO: Verify if it's a new order
	///Redis query
	///Postgres query
	///Delete the message from the queue (if it isn't new)

	credentials, err := o.store.RetrieveMeliCredentialsFromMeliUserID(order.Store)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error to process order - get credentials")
	}
	orderData, err := o.meli.FetchOrder(order.OrderId, credentials.MeliAccessToken)
	if err != nil {
		return err
	}
	o.announcement.RemoveQuantity(orderData.Items)
}
