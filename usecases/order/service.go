package order

import (
	"errors"
	"strconv"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/usecases/announcement"
	"github.com/Vractos/kloni/usecases/common"
	"github.com/Vractos/kloni/usecases/store"
	"github.com/Vractos/kloni/utils"
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
		if i, ok := m[k]; ok && item.Sku != "" {
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
		o.logger.Warn("Fail to retrieve order from cache", zap.String("order_id", order.OrderId), zap.Error(err))
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
		// TODO: Checks if the order was successfully deleted from the queue
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
		if item.VariationID == 0 {
			itsOk[i] = item.ID
		} else {
			itsOk[i] = strconv.Itoa(item.VariationID)
		}
	}

	removeDuplicateItens(&orderData.Items)

	var anns []common.MeliAnnouncement
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
							anns = append(anns, common.MeliAnnouncement{
								ID:       cln.ID,
								Title:    cln.Title,
								Sku:      cln.Sku,
								Quantity: cln.Quantity - item.Quantity,
							})
						}
					}
					continue
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
			if cln.Variations != nil {
				ann := common.MeliAnnouncement{
					ID:    cln.ID,
					Title: cln.Title,
					Sku:   cln.Sku,
				}
				for _, variation := range cln.Variations {
					if variation.AvailableQuantity > 0 && !utils.Contains(&itsOk, strconv.Itoa(variation.ID)) {
						ann.Variations = append(ann.Variations, struct {
							ID                int
							AvailableQuantity int
						}{
							ID:                variation.ID,
							AvailableQuantity: variation.AvailableQuantity - item.Quantity,
						})
					}
				}
				if len(ann.Variations) > 0 {
					anns = append(anns, ann)
				}
				continue
			}
			if !utils.Contains(&itsOk, cln.ID) {
				anns = append(anns, common.MeliAnnouncement{
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
	var toChangeQuantity []common.MeliAnnouncement

	// It contains announcements that were not possible to change the quantity
	var canNotChangeQuantity []common.MeliAnnouncement

	for _, ann := range anns {
		if ann.Variations != nil {
			for _, variation := range ann.Variations {
				if err := o.announce.UpdateQuantity(ann.ID, variation.AvailableQuantity, *credentials, variation.ID); err != nil {
					var annErr *announcement.AnnouncementError

					annFailed := ann

					annFailed.Variations = []struct {
						ID                int
						AvailableQuantity int
					}{variation}

					if errors.As(err, &annErr) && annErr.IsAbleToRetry {
						toChangeQuantity = append(toChangeQuantity, annFailed)
					} else {
						canNotChangeQuantity = append(canNotChangeQuantity, annFailed)
					}
				}
			}
			continue
		}

		if err := o.announce.UpdateQuantity(ann.ID, ann.Quantity, *credentials); err != nil {
			var annErr *announcement.AnnouncementError
			if errors.As(err, &annErr) && annErr.IsAbleToRetry {
				toChangeQuantity = append(toChangeQuantity, ann)
			} else {
				canNotChangeQuantity = append(canNotChangeQuantity, ann)
			}
		}
	}

	// TODO Turn into a goroutine
	for _, ann := range toChangeQuantity {
		if ann.Variations != nil {
			for _, variation := range ann.Variations {
				if err := o.announce.UpdateQuantity(ann.ID, variation.AvailableQuantity, *credentials, variation.ID); err != nil {
					canNotChangeQuantity = append(canNotChangeQuantity, ann)
				}
			}
			continue
		}

		if err := o.announce.UpdateQuantity(ann.ID, ann.Quantity, *credentials); err != nil {
			canNotChangeQuantity = append(canNotChangeQuantity, ann)
		}
	}

	if canNotChangeQuantity != nil {
		oErr := &OrderError{
			Message:            "Couldn't change the announcements",
			AnnouncementsError: canNotChangeQuantity,
		}
		o.logger.Error(
			oErr.Message,
			oErr,
			zap.Any("Announcements that were not possible to change the quantity", oErr.AnnouncementsError),
		)
		return oErr
	}
	// -------------------------------------------
	// --------- STORING ORDER IN THE DB ---------
	// -------------------------------------------
	orderItems := make([]entity.OrderItem, len(orderData.Items))
	for i, item := range orderData.Items {
		orderItems[i] = entity.OrderItem{
			Title:       item.Title,
			Quantity:    item.Quantity,
			Sku:         item.Sku,
			VariationID: item.VariationID,
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

// Exporting for testing purposes
var RemoveDuplicateItensTest = removeDuplicateItens
