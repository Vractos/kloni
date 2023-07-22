package order

import (
	"errors"

	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/usecases/announcement"
	"github.com/Vractos/dolly/usecases/common"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/Vractos/dolly/utils"
	"go.uber.org/zap"
)

type OrderService struct {
	queue    Queue
	meli     common.MercadoLivre
	store    store.UseCase
	announce announcement.UseCase
	repo     Repository
	cache    Cache
	logger   common.Logger
}

func NewOrderService(
	queue Queue,
	mercadolivre common.MercadoLivre,
	storeUseCase store.UseCase,
	announceUseCase announcement.UseCase,
	repository Repository,
	cache Cache,
	logger common.Logger,
) *OrderService {
	return &OrderService{
		queue:    queue,
		meli:     mercadolivre,
		store:    storeUseCase,
		announce: announceUseCase,
		repo:     repository,
		cache:    cache,
		logger:   logger,
	}
}

func (o *OrderService) ProcessWebhook(input OrderWebhookDtoInput) error {
	if err := o.queue.PostOrderNotification(input); err != nil {
		o.logger.Error(
			"Error to post order notification",
			err,
			zap.String("notification_id", input.ID),
			zap.Int("user_id", input.UserID),
			zap.Int("attempts", input.Attempts),
			zap.String("sent", input.Sent),
		)
		return errors.New("error to post order notification")
	}
	return nil
}

// removeDuplicateItens removes duplicate items from the given slice of OrderItem pointers.
// It modifies the original slice by removing duplicates and summing their quantities.
// The function uses the Sku field of each OrderItem to determine if it is a duplicate.
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
		o.logger.Warn("Fail to retrieve order cache", zap.String("order_id", order.OrderId), zap.Error(err))
	}

	if status != nil {
		o.queue.DeleteOrderNotification(order.ReceiptHandle)
		return nil
	}

	odrSaved, err := o.repo.GetOrder(order.OrderId)
	if err != nil {
		o.logger.Error("Fail to retrieve order from the DB", err, zap.String("order_id", order.OrderId))
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
		o.logger.Error("Error in retrieving Meli credentials during order processing", err, zap.String("order_id", order.OrderId))
		return errors.New("error to process order - get credentials")
	}

	// --------------------------------------
	// --------- GETTING ORDER DATA ---------
	// --------------------------------------
	orderData, err := o.meli.FetchOrder(order.OrderId, credentials.MeliAccessToken)
	if err != nil {
		o.logger.Error("Error to fetch the order", err, zap.String("order_id", order.OrderId))
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
		if item.Sku == "" {
			o.logger.Warn("The product doesn't have sku", zap.String("order_id", order.OrderId), zap.String("announcement_id", item.ID))
			continue
		}
		clones, err := o.announce.RetrieveAnnouncements(item.Sku, *credentials)
		if err != nil {
			var annErr *announcement.AnnouncementError
			if errors.As(err, &annErr) {
				o.logger.Warn("Fail in retrieving the order product clones", zap.Error(err), zap.String("order_id", order.OrderId), zap.String("sku", item.Sku))
				// Retry
				if annErr.IsAbleToRetry {
					o.logger.Info("Retrying to retrieve order products clones...", zap.String("order_id", order.OrderId), zap.String("sku", item.Sku))
					clns, err := o.announce.RetrieveAnnouncements(item.Sku, *credentials)
					if err != nil {
						o.logger.Error("Error in retrieving the order product clones", err, zap.String("order_id", order.OrderId))
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
					o.logger.Error("Error in retrieving the order product clones", err, zap.String("order_id", order.OrderId), zap.String("sku", item.Sku))
					return err
				}
			} else {
				o.logger.Error("Error in retrieving the order product clones", err, zap.String("order_id", order.OrderId), zap.String("sku", item.Sku))
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
		oErr := &OrderError{
			Message:            "Couldn't change the announcements",
			AnnouncementsError: canNotChangeQuantity,
		}

		o.logger.Error(oErr.Message, err, zap.Reflect("Announcements that were not possible to change the quantity", canNotChangeQuantity))
		return oErr
	}

	if canNotChangeQuantity != nil {
		oErr := &OrderError{
			Message:            "Couldn't change the announcements",
			AnnouncementsError: canNotChangeQuantity,
		}
		o.logger.Warn(oErr.Message, zap.Reflect("Announcements that were not possible to change the quantity", canNotChangeQuantity))
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
		o.logger.Error("Fail to generate the order entity", err, zap.String("order_id", orderData.ID))
		return err
	}

	if err := o.repo.RegisterOrder(odr); err != nil {
		o.logger.Error("Fail to store the order", err, zap.String("order_id", orderData.ID))
		return errors.New("couldn't store order")
	}

	// -------------------------------------------
	// ------------- CACHING ORDER ---------------
	// -------------------------------------------

	if err := o.cache.SetOrder(odr); err != nil {
		o.logger.Warn("Fail to cache the order", zap.String("order_id", orderData.ID))
	}
	o.queue.DeleteOrderNotification(order.ReceiptHandle)
	return nil
}

// Exports for testing purposes
var RemoveDuplicateItensTest = removeDuplicateItens
