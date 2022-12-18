package order

import (
	"errors"
	"log"

	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/usecases/announcement"
	"github.com/Vractos/dolly/usecases/common"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/Vractos/dolly/utils"
)

type OrderService struct {
	queue    Queue
	meli     common.MercadoLivre
	store    store.UseCase
	announce announcement.UseCase
	repo     Repository
	cache    Cache
}

func NewOrderService(
	queue Queue,
	mercadolivre common.MercadoLivre,
	storeUseCase store.UseCase,
	announceUseCase announcement.UseCase,
	repository Repository,
	cache Cache,
) *OrderService {
	return &OrderService{
		queue:    queue,
		meli:     mercadolivre,
		store:    storeUseCase,
		announce: announceUseCase,
		repo:     repository,
		cache:    cache,
	}
}

func (o *OrderService) ProcessWebhook(input OrderWebhookDtoInput) error {
	if err := o.queue.PostOrderNotification(input); err != nil {
		log.Println(err.Error())
		return errors.New("error to post order notification")
	}
	return nil
}

func removeDuplicateItens(items *[]common.OrderItem) {
	var unique []common.OrderItem
	type key struct{ Sku string }
	m := make(map[key]int)
	for _, item := range *items {
		k := key{item.Sku}
		if i, ok := m[k]; ok {
			unique[i].Quantity = unique[i].Quantity + item.Quantity
		} else {
			m[k] = len(unique)
			unique = append(unique, item)
		}
	}
	*items = unique
}

func (o *OrderService) ProcessOrder(order OrderMessage) error {
	// --------------------------------------
	// --- VERIFYING IF IT'S A NEW ORDER ---
	// --------------------------------------

	status, err := o.cache.GetOrder(order.OrderId)
	if err != nil {
		log.Println(err.Error())
	}

	if status != nil {
		o.queue.DeleteOrderNotification(order.ReceiptHandle)
		return nil
	}

	odrSaved, err := o.repo.GetOrder(order.OrderId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if odrSaved != nil {
		o.queue.DeleteOrderNotification(order.ReceiptHandle)
		return nil
	}

	// --------------------------------------
	// --------- GETTING CREDENTIALS --------
	// --------------------------------------
	credentials, err := o.store.RetrieveMeliCredentialsFromMeliUserID(order.Store)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error to process order - get credentials")
	}

	// --------------------------------------
	// --------- GETTING ORDER DATA ---------
	// --------------------------------------
	orderData, err := o.meli.FetchOrder(order.OrderId, credentials.MeliAccessToken)
	if err != nil {
		return err
	}

	// -----------------------------------------------------
	// ---- DEFINING ANNOUNCEMENTS THAT MUST BE CHANGED ----
	// -----------------------------------------------------

	// It contains announcements IDs that Mercado Livre has already updated
	itsOk := make([]string, len(orderData.Items))
	for i, item := range orderData.Items {
		itsOk[i] = item.ID
	}

	removeDuplicateItens(&orderData.Items)

	var anns []common.OrderItem
	for _, item := range orderData.Items {
		clones, err := o.announce.RetrieveAnnouncements(item.Sku, *credentials)
		if err != nil {
			var annErr *announcement.AnnouncementError
			if errors.As(err, &annErr) {
				log.Println(annErr.Error())
				// Retry
				if annErr.IsAbleToRetry {
					clns, err := o.announce.RetrieveAnnouncements(item.Sku, *credentials)
					if err != nil {
						log.Println(err.Error())
						return err
					}
					for _, cln := range *clns {
						if !utils.Contains(&itsOk, cln.ID) {
							anns = append(anns, common.OrderItem{
								ID:       cln.ID,
								Title:    cln.Title,
								Sku:      cln.Sku,
								Quantity: cln.Quantity - item.Quantity,
							})
						}
					}
				} else {
					return err
				}
			} else {
				log.Println(err.Error())
				return err
			}
		}
		for _, cln := range *clones {
			if !utils.Contains(&itsOk, cln.ID) {
				anns = append(anns, common.OrderItem{
					ID:       cln.ID,
					Title:    cln.Title,
					Sku:      cln.Sku,
					Quantity: cln.Quantity - item.Quantity,
				})
			}
		}
	}

	// -----------------------------------------
	// --------- CHANGING THE QUANTITY ---------
	// -----------------------------------------

	// It contains announcements that were not possible to change the quantity,
	// but will be retried
	var toChangeQuantity []struct {
		id       string
		quantity int
	}

	// It contains announcements that were not possible to change the quantity
	var canNotChangeQuantity []struct {
		id       string
		quantity int
	}

	for _, ann := range anns {
		if err := o.announce.UpdateQuantity(ann.ID, ann.Quantity, *credentials); err != nil {
			var annErr *announcement.AnnouncementError
			if errors.As(err, &annErr) && annErr.IsAbleToRetry {
				toChangeQuantity = append(toChangeQuantity, struct {
					id       string
					quantity int
				}{ann.ID, ann.Quantity})
			} else {
				canNotChangeQuantity = append(canNotChangeQuantity, struct {
					id       string
					quantity int
				}{ann.ID, ann.Quantity})
			}
		}
	}

	// TODO Turn into a goroutine
	// if utils.PercentOf(len(toChangeQuantity), len(orderData.Items)) >= 40.0 {
	for _, ann := range toChangeQuantity {
		if err := o.announce.UpdateQuantity(ann.id, ann.quantity, *credentials); err != nil {
			canNotChangeQuantity = append(canNotChangeQuantity, struct {
				id       string
				quantity int
			}{ann.id, ann.quantity})
		}
	}
	// }

	if utils.PercentOf(len(canNotChangeQuantity), len(orderData.Items)) >= 20.0 {
		return &OrderError{
			Message:            "Couldn't change the announcements",
			AnnouncementsError: canNotChangeQuantity,
		}
	}

	// -------------------------------------------
	// --------- STORING ORDER IN THE DB ---------
	// -------------------------------------------
	orderItems := make([]entity.OrderItem, len(orderData.Items))
	for i, item := range orderData.Items {
		orderItems[i] = entity.OrderItem{
			Title:    item.Title,
			Quantity: item.Quantity,
			Sku:      item.Sku,
		}
	}

	odr, err := entity.NewOrder(credentials.StoreID, orderData.ID, orderItems, entity.OrderStatus(orderData.Status))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err := o.repo.RegisterOrder(odr); err != nil {
		log.Println(err.Error())
		return errors.New("couldn't store order")
	}

	// -------------------------------------------
	// ------------- CACHING ORDER ---------------
	// -------------------------------------------

	if err := o.cache.SetOrder(odr); err != nil {
		log.Println(err.Error())
	}
	o.queue.DeleteOrderNotification(order.ReceiptHandle)
	return nil
}
