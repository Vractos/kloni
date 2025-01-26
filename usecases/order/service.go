// Package order implements the order processing and management functionality
package order

import (
	"errors"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/usecases/announcement"
	"github.com/Vractos/kloni/usecases/common"
	"github.com/Vractos/kloni/usecases/store"
	"github.com/Vractos/kloni/utils"
	"go.uber.org/zap"
)

// Error definitions for order-related operations
var (
	// ErrOrderNotFound is returned when an order cannot be found
	ErrOrderNotFound = errors.New("order not found")
	// ErrInvalidQuantity is returned when an order has an invalid quantity
	ErrInvalidQuantity = errors.New("invalid quantity")
	// ErrInvalidSKU is returned when an order has an invalid SKU
	ErrInvalidSKU = errors.New("invalid SKU")
	// ErrPostingOrderNotification is returned when there's an error posting an order notification
	ErrPostingOrderNotification = errors.New("error to post order notification")
	// ErrProcessingOrder is returned when there's an error processing an order
	ErrProcessingOrder = errors.New("error processing order")
	// ErrSyncingQuantities is returned when there's an error syncing quantities between cloned items
	ErrSyncingQuantities = errors.New("error syncing quantities")
	// ErrCredentialsNotFound is returned when store credentials cannot be found
	ErrCredentialsNotFound = errors.New("credentials not found")
)

type SyncContext struct {
	Item                common.OrderItem
	Credentials         *store.Credentials
	AllCredentials      *[]store.Credentials
	CredentialsHashMap  map[interface{}]store.Credentials
	ProcessedItems      map[string]bool
	ProcessedVariations map[string][]int
}

// OrderService handles all order-related operations including processing orders,
// managing webhooks, and synchronizing quantities across cloned items.
type OrderService struct {
	queue    Queue                // Queue service for handling order notifications
	meli     common.MercadoLivre  // Mercado Livre API client
	store    store.UseCase        // Store management use case
	announce announcement.UseCase // Announcement management use case
	repo     Repository           // Order repository for data persistence
	cache    Cache                // Cache service for temporary data storage
	logger   common.Logger        // Logger for error and info logging
}

// NewOrderService creates a new instance of OrderService with all required dependencies.
// It follows the dependency injection pattern to provide flexibility and testability.
//
// Parameters:
//   - queue: Queue service for handling order notifications
//   - mercadolivre: Mercado Livre API client
//   - storeUseCase: Store management use case
//   - announceUseCase: Announcement management use case
//   - repository: Order repository for data persistence
//   - cache: Cache service for temporary data storage
//   - logger: Logger for error and info logging
//
// Returns:
//   - *OrderService: A new instance of OrderService
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

// ProcessWebhook handles incoming order webhooks from Mercado Livre.
// It posts the order notification to a queue for asynchronous processing.
//
// Parameters:
//   - input: OrderWebhookDtoInput containing the webhook data
//
// Returns:
//   - error: ErrPostingOrderNotification if there's an error posting to the queue, nil otherwise
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
		return ErrPostingOrderNotification
	}
	return nil
}

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
func (o *OrderService) ProcessOrder(order OrderMessage) error {
	if precessed, err := o.orderExists(order); err != nil {
		return err
	} else if precessed {
		return nil
	}

	credentials, allCredentials, credMap, err := o.getStoreCredentials(order.Store)
	if err != nil {
		o.logger.Error("Error in retrieving store credentials", err, zap.String("store_id", order.Store))
		return err
	}

	orderData, err := o.fetchOrderData(order.OrderId, credentials.AccessToken)
	if err != nil {
		return err
	}

	// Track processed items and variations to avoid duplicate processing
	processedItems := make(map[string]bool)
	processedVariations := make(map[string][]int)
	for _, item := range orderData.Items {
		if item.VariationID == 0 {
			processedItems[item.ID] = true
		} else {
			sku := item.Sku
			if _, exists := processedVariations[sku]; !exists {
				processedVariations[sku] = make([]int, 0)
			}
			processedVariations[sku] = append(processedVariations[sku], item.VariationID)
		}
	}

	for _, item := range orderData.Items {
		if item.Sku == "" {
			o.logger.Warn("The product doesn't have sku",
				zap.String("order_id", orderData.ID),
				zap.String("announcement_id", item.ID),
			)
			continue
		}

		ctx := &SyncContext{
			Item:                item,
			Credentials:         credentials,
			AllCredentials:      allCredentials,
			CredentialsHashMap:  credMap,
			ProcessedItems:      processedItems,
			ProcessedVariations: processedVariations,
		}

		if err := o.syncItemQuantities(ctx); err != nil {
			return err
		}
	}

	// ------------------------------------
	// --------- STORE ORDER IN DB --------
	// ------------------------------------
	orderItems := make([]entity.OrderItem, len(orderData.Items))
	for i, item := range orderData.Items {
		orderItems[i] = entity.OrderItem{
			Title:       item.Title,
			Quantity:    item.Quantity,
			Sku:         item.Sku,
			VariationID: item.VariationID,
		}
	}

	odr, err := entity.NewOrder(credentials.ID, order.OrderId, orderItems, entity.OrderStatus(orderData.Status))
	if err != nil {
		o.logger.Error("Fail to generate the order entity", err, zap.String("order_id", orderData.ID))
		return err
	}

	if err := o.repo.RegisterOrder(odr); err != nil {
		o.logger.Error("Fail to store the order", err, zap.String("order_id", orderData.ID))
		return errors.New("couldn't store order")
	}

	// ------------------------------
	// --------- CACHE ORDER --------
	// ------------------------------
	if err := o.cache.SetOrder(odr); err != nil {
		o.logger.Warn("Fail to cache the order", zap.String("order_id", orderData.ID))
	}

	o.queue.DeleteOrderNotification(order.ReceiptHandle)
	return nil
}

// orderExists verifies if an order has already been processed by checking both cache and repository storage.
// Returns true and deletes the notification if the order exists, false if it's a new order.
// This prevents duplicate order processing and ensures data consistency.
//
// Parameters:
//   - order: OrderMessage to validate
//
// Returns:
//   - error: Error if validation fails, nil if order is valid
func (o *OrderService) orderExists(order OrderMessage) (bool, error) {
	status, err := o.cache.GetOrder(order.OrderId)
	if err != nil {
		o.logger.Warn("Fail to retrieve order from cache", zap.String("order_id", order.OrderId), zap.Error(err))
	}

	if status != nil {
		o.queue.DeleteOrderNotification(order.ReceiptHandle)
		return true, nil
	}

	odrSaved, err := o.repo.GetOrder(order.OrderId)
	if err != nil {
		o.logger.Error("Fail to retrieve order from the DB", err, zap.String("order_id", order.OrderId))
		return false, err
	}

	if odrSaved != nil {
		o.queue.DeleteOrderNotification(order.ReceiptHandle)
		return true, nil
	}

	return false, nil
}

// getStoreCredentials retrieves all store credentials and returns the credentials
// for the specific store that the order belongs to, along with all credentials.
//
// Parameters:
//   - storeID: ID of the store to retrieve credentials for
//
// Returns:
//   - *store.Credentials: Credentials for the specific store
//   - *[]store.Credentials: All available store credentials
//   - error: ErrCredentialsNotFound or other errors
func (o *OrderService) getStoreCredentials(storeID string) (*store.Credentials, *[]store.Credentials, map[interface{}]store.Credentials, error) {
	credentials, err := o.store.RetrieveMeliCredentialsFromMeliUserID(storeID)
	if err != nil {
		o.logger.Error("Error in retrieving Meli credentials during order processing", err, zap.String("store_id", storeID))
		return nil, nil, nil, ErrCredentialsNotFound
	}

	credMap, err := utils.HashMap(credentials, "UserID")
	if err != nil {
		o.logger.Error("Error in converting credentials to map", err, zap.String("store_id", storeID))
		return nil, nil, nil, ErrCredentialsNotFound
	}

	rCredentials := findCredentialsByMeliUserID(storeID, credMap)
	if rCredentials == nil {
		return nil, nil, nil, ErrCredentialsNotFound
	}

	return rCredentials, credentials, credMap, nil
}

// fetchOrderData retrieves the order data from Mercado Livre API.
// It also removes duplicate items from the order.
//
// Parameters:
//   - orderID: ID of the order to fetch
//   - accessToken: Access token for Mercado Livre API
//
// Returns:
//   - *common.MeliOrder: Order data
//   - error: ErrProcessingOrder or other errors
func (o *OrderService) fetchOrderData(orderID string, accessToken string) (*common.MeliOrder, error) {
	orderData, err := o.meli.FetchOrder(orderID, accessToken)
	if err != nil {
		o.logger.Error("Error to fetch the order", err, zap.String("order_id", orderID))
		return nil, ErrProcessingOrder
	}

	removeDuplicateItems(&orderData.Items)
	return orderData, nil
}

// syncItemQuantities synchronizes quantities for a specific item across all cloned announcements.
//
// Parameters:
//   - ctx: SyncContext containing the item and processing context
//
// Returns:
//   - error: Error if synchronization fails
func (o *OrderService) syncItemQuantities(
	ctx *SyncContext,
) error {
	clones, err := o.announce.RetrieveAnnouncementsFromAllAccounts(ctx.Item.Sku, ctx.AllCredentials)

	if err != nil {
		return o.handleAnnouncementError(err, ctx)
	}

	return o.updateCloneQuantities(clones, ctx)
}

// updateCloneQuantities updates the quantities of cloned items.
// It handles both simple items and items with variations.
//
// Parameters:
//   - clones: List of cloned announcements
//   - ctx: SyncContext containing the item and processing context
//
// Returns:
//   - error: ErrSyncingQuantities or other errors
func (o *OrderService) updateCloneQuantities(
	clones *[]announcement.Announcements,
	ctx *SyncContext,
) error {
	for _, cln := range *clones {
		announcements := []common.MeliAnnouncement{}

		credentials := findCredentialsByAccountID(cln.AccountID, ctx.CredentialsHashMap)

		for _, cl := range *cln.Announcements {
			if cl.Variations != nil {
				ann := o.handleVariationUpdate(cl, ctx.Item, ctx.ProcessedVariations)
				if ann != nil {
					announcements = append(announcements, *ann)
				}
			} else {
				ann := o.handleSimpleUpdate(cl, ctx.Item, ctx.ProcessedItems)
				if ann != nil {
					announcements = append(announcements, *ann)
				}
			}
		}

		for _, ann := range announcements {
			if ann.Variations != nil {
				for _, variation := range ann.Variations {
					if err := o.announce.UpdateQuantity(ann.ID, variation.AvailableQuantity, *credentials, variation.ID); err != nil {
						odrErr := &OrderError{
							Message: "Error updating announcements",
							AnnouncementsError: []announcement.Announcements{
								{
									AccountID: ctx.Credentials.ID,
									Announcements: &[]common.MeliAnnouncement{
										{
											ID:       ann.ID,
											Title:    ann.Title,
											Sku:      ann.Sku,
											Quantity: ann.Quantity,
											Variations: []struct {
												ID                int
												AvailableQuantity int
											}{
												{
													ID:                variation.ID,
													AvailableQuantity: variation.AvailableQuantity,
												},
											},
										},
									},
								},
							},
						}
						o.logger.Error("Error updating announcements", odrErr)
						return ErrSyncingQuantities
					}
				}
			} else {
				if err := o.announce.UpdateQuantity(ann.ID, ann.Quantity, *credentials); err != nil {
					odrErr := &OrderError{
						Message: "Error updating announcements",
						AnnouncementsError: []announcement.Announcements{
							{
								AccountID: ctx.Credentials.ID,
								Announcements: &[]common.MeliAnnouncement{
									{
										ID:       ann.ID,
										Title:    ann.Title,
										Sku:      ann.Sku,
										Quantity: ann.Quantity,
									},
								},
							},
						},
					}

					o.logger.Error("Error updating announcements", odrErr)
					return ErrSyncingQuantities
				}
			}
		}
	}
	return nil
}

// handleVariationUpdate processes quantity updates for items with variations.
// It checks if variations have already been processed and updates quantities accordingly.
//
// Parameters:
//   - cl: The announcement to update
//   - item: The order item being processed
//   - processedVariations: Map of already processed variations
//
// Returns:
//   - *common.MeliAnnouncement: Updated announcement or nil if no updates needed
func (o *OrderService) handleVariationUpdate(
	cl common.MeliAnnouncement,
	item common.OrderItem,
	processedVariations map[string][]int,
) *common.MeliAnnouncement {
	ann := common.MeliAnnouncement{
		ID:    cl.ID,
		Title: cl.Title,
		Sku:   cl.Sku,
	}

	hasUpdates := false
	for _, variation := range cl.Variations {
		// Skip if this variation was already processed in the order
		if variations, exists := processedVariations[item.Sku]; exists && utils.Contains(&variations, variation.ID) {
			continue
		}

		if variation.AvailableQuantity > 0 {
			newQuantity := variation.AvailableQuantity - item.Quantity
			ann.Variations = append(ann.Variations, struct {
				ID                int
				AvailableQuantity int
			}{
				ID:                variation.ID,
				AvailableQuantity: newQuantity,
			})
			hasUpdates = true
		}
	}

	if hasUpdates {
		return &ann
	}
	return nil
}

// handleSimpleUpdate processes quantity updates for simple items without variations.
//
// Parameters:
//   - cl: The announcement to update
//   - item: The order item being processed
//   - processedItems: Map of already processed items
//
// Returns:
//   - *common.MeliAnnouncement: Updated announcement or nil if no updates needed
func (o *OrderService) handleSimpleUpdate(
	cl common.MeliAnnouncement,
	item common.OrderItem,
	processedItems map[string]bool,
) *common.MeliAnnouncement {
	// Skip if this item was already processed in the order
	if processedItems[cl.ID] {
		return nil
	}

	if cl.ID != item.ID && cl.Quantity > 0 {
		newQuantity := cl.Quantity - item.Quantity
		return &common.MeliAnnouncement{
			ID:       cl.ID,
			Title:    cl.Title,
			Sku:      cl.Sku,
			Quantity: newQuantity,
		}
	}
	return nil
}

// handleAnnouncementError handles errors that occur during announcement retrieval.
// It implements retry logic for retriable errors.
//
// Parameters:
//   - err: The error that occurred
//   - ctx: SyncContext containing the item and processing context
//
// Returns:
//   - error: ErrProcessingOrder or other errors
func (o *OrderService) handleAnnouncementError(
	err error,
	ctx *SyncContext,
) error {
	var annErr *announcement.AnnouncementError
	if errors.As(err, &annErr) {
		o.logger.Warn("Fail in retrieving the order product clones", zap.Error(err), zap.String("sku", ctx.Item.Sku))
		if annErr.IsAbleToRetry {
			o.logger.Info("Retrying to retrieve order products clones...", zap.String("sku", ctx.Item.Sku))
			return o.retryAnnouncementRetrieval(ctx)
		}
	}
	o.logger.Error("Error in retrieving the order product clones", err, zap.String("sku", ctx.Item.Sku))
	return ErrProcessingOrder
}

// retryAnnouncementRetrieval attempts to retrieve announcements again after a failure.
//
// Parameters:
//   - ctx: SyncContext containing the item and processing context
//
// Returns:
//   - error: ErrProcessingOrder or other errors
func (o *OrderService) retryAnnouncementRetrieval(
	ctx *SyncContext,
) error {
	clones, err := o.announce.RetrieveAnnouncementsFromAllAccounts(ctx.Item.Sku, ctx.AllCredentials)
	if err != nil {
		o.logger.Error("Error in retrieving the order product clones", err, zap.String("sku", ctx.Item.Sku))
		return ErrProcessingOrder
	}
	return o.updateCloneQuantities(clones, ctx)
}

// removeDuplicateItems removes duplicate items from a slice of OrderItems.
// Items are considered duplicates if they have the same SKU.
// Quantities of duplicate items are summed together.
//
// Parameters:
//   - items: Pointer to the slice of OrderItems to process
func removeDuplicateItems(items *[]common.OrderItem) {
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

// findCredentialsByMeliUserID finds store credentials by Mercado Livre user ID.
//
// Parameters:
//   - meliUserId: Mercado Livre user ID to search for
//   - hashMap: Map of credentials indexed by user ID
//
// Returns:
//   - *store.Credentials: Found credentials or nil if not found
func findCredentialsByMeliUserID(meliUserId string, hashMap map[interface{}]store.Credentials) *store.Credentials {
	if val, ok := hashMap[meliUserId]; ok {
		return &val
	}
	return nil
}

// findCredentialsByAccountID finds store credentials by account ID.
//
// Parameters:
//   - accountID: Account ID to search for
//   - hashMap: Map of credentials indexed by user ID
//
// Returns:
//   - *store.Credentials: Found credentials or nil if not found
func findCredentialsByAccountID(accountID entity.ID, hashMap map[interface{}]store.Credentials) *store.Credentials {
	for _, v := range hashMap {
		if v.ID == accountID {
			return &v
		}
	}
	return nil
}

// RemoveDuplicateItemsTest is a test helper function that exposes removeDuplicateItems for testing
func RemoveDuplicateItemsTest(items *[]common.OrderItem) {
	removeDuplicateItems(items)
}
