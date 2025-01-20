package order

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/usecases/announcement"
	mock_announcement "github.com/Vractos/kloni/usecases/announcement/mock"
	common "github.com/Vractos/kloni/usecases/common"
	common_mock "github.com/Vractos/kloni/usecases/common/mock"
	"github.com/Vractos/kloni/usecases/order"
	mock_order "github.com/Vractos/kloni/usecases/order/mock"
	"github.com/Vractos/kloni/usecases/store"
	mock_store "github.com/Vractos/kloni/usecases/store/mock"
	"github.com/Vractos/kloni/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// TestProcessWebhook tests the ProcessWebhook method of OrderService.
// It verifies:
// 1. Successful webhook processing
// 2. Error handling when posting to queue fails
func TestProcessWebhook(t *testing.T) {
	tests := []struct {
		name       string
		input      order.OrderWebhookDtoInput
		setupMocks func(*Mocks)
		wantErr    bool
		errMessage string
	}{
		{
			name: "successful webhook processing",
			input: order.OrderWebhookDtoInput{
				ID:            "test-id",
				Resource:      "/orders/20210101000000",
				UserID:        1,
				Attempts:      0,
				Topic:         "orders_v2",
				ApplicationID: 1,
				Received:      "2022-10-30T16:19:20.129Z",
				Sent:          "2022-10-30T16:19:20.106Z",
			},
			setupMocks: func(m *Mocks) {
				m.mockOrderQueue.EXPECT().PostOrderNotification(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error posting to queue",
			input: order.OrderWebhookDtoInput{
				ID:            "test-id",
				Resource:      "/orders/20210101000000",
				UserID:        1,
				Attempts:      0,
				Topic:         "orders_v2",
				ApplicationID: 1,
				Received:      "2022-10-30T16:19:20.129Z",
				Sent:          "2022-10-30T16:19:20.106Z",
			},
			setupMocks: func(m *Mocks) {
				m.mockOrderQueue.EXPECT().PostOrderNotification(gomock.Any()).Return(errors.New("queue error"))
				m.mockLogger.EXPECT().Error(
					"Error to post order notification",
					gomock.Any(),
					zap.String("notification_id", "test-id"),
					zap.Int("user_id", 1),
					zap.Int("attempts", 0),
					zap.String("sent", "2022-10-30T16:19:20.106Z"),
				)
			},
			wantErr:    true,
			errMessage: "error to post order notification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			tt.setupMocks(mocks)

			err := orderService.ProcessWebhook(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMessage {
				t.Errorf("ProcessWebhook() error = %v, want %v", err, tt.errMessage)
			}
		})
	}
}

// TestProcessOrder tests the ProcessOrder method of OrderService.
// Test scenarios include:
// 1. Standard order processing
// 2. Orders with variations
// 3. Orders with multiple items
// 4. Orders with items lacking SKUs
// 5. Various edge cases and error conditions
func TestProcessOrder(t *testing.T) {

	// -------------------------------------------------------
	// ------------------------ Mocks ------------------------
	// -------------------------------------------------------
	accountId := entity.ID(uuid.New())
	seccondAccountId := entity.ID(uuid.New())
	defaultOrderMessage := order.OrderMessage{
		Store:         "1",
		OrderId:       "20210101000000",
		Attempts:      0,
		ReceiptHandle: "test-receipt-handle",
	}

	defaultMeliCredentials := &[]store.Credentials{
		{
			ID: accountId,
			MeliCredential: &common.MeliCredential{
				AccessToken: "test-access-token",
				UserID:      "1",
			},
		},
		{
			ID: seccondAccountId,
			MeliCredential: &common.MeliCredential{
				AccessToken: "test-access-token-2",
				UserID:      "2",
			},
		},
	}

	nilCredentials := &[]store.Credentials{
		{
			ID: accountId,
			MeliCredential: &common.MeliCredential{
				AccessToken: "missing",
			},
		},
	}

	defaultMeliOrder := &common.MeliOrder{
		ID:          "20210101000000",
		DateCreated: "2022-10-30T16:19:20.129Z",
		Status:      common.Paid,
		Items: []common.OrderItem{
			{
				ID:       "1",
				Title:    "test-title",
				Sku:      "test-sku",
				Quantity: 1,
			},
		},
	}

	defaultMeliAnnouncementsClones := []announcement.Announcements{
		{
			AccountID: accountId,
			Announcements: &[]common.MeliAnnouncement{
				{
					ID:       "1",
					Title:    "test-title",
					Quantity: 1,
					Price:    1.0,
					Sku:      "test-sku",
				},
				{
					ID:       "2",
					Title:    "test-title2",
					Quantity: 1,
					Price:    1.0,
					Sku:      "test-sku",
				},
			},
		},
	}

	// -------------------------------------------------------
	// ---------------- Scenarios -> New order ---------------
	// -------------------------------------------------------

	newOrderScenarios := []struct {
		name                     string
		accountId                entity.ID
		orderMessage             order.OrderMessage
		meliOrder                *common.MeliOrder
		meliCredentials          *[]store.Credentials
		rootCredentials          *store.Credentials
		orderAnnouncementsClones [][]announcement.Announcements
		odr                      *entity.Order
		OrderMatcher             *OrderMatcher
	}{
		{
			name:         "default scenario",
			accountId:    accountId,
			orderMessage: defaultOrderMessage,
			meliOrder: &common.MeliOrder{
				ID:          "20210101000000",
				DateCreated: "2022-10-30T16:19:20.129Z",
				Status:      common.Paid,
				Items: []common.OrderItem{
					{
						ID:       "1",
						Title:    "test-title",
						Sku:      "test-sku",
						Quantity: 1,
					},
				},
			},
			meliCredentials: defaultMeliCredentials,
			rootCredentials: &(*defaultMeliCredentials)[0],
			orderAnnouncementsClones: [][]announcement.Announcements{
				{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "1",
								Title:    "test-title",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
							{
								ID:       "2",
								Title:    "test-title2",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
						},
					},
				},
			},
			odr: &entity.Order{
				AccountID:     accountId,
				MarketplaceID: "20210101000000",
				Status:        "paid",
				Items: []entity.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "test-sku",
					},
				},
			},
			OrderMatcher: &OrderMatcher{},
		},
		{
			name:         "default scenario with variation",
			accountId:    accountId,
			orderMessage: defaultOrderMessage,
			meliOrder: &common.MeliOrder{
				ID:          "20210101000000",
				DateCreated: "2022-10-30T16:19:20.129Z",
				Status:      common.Paid,
				Items: []common.OrderItem{
					{
						ID:       "1",
						Title:    "test-title1",
						Sku:      "test-sku",
						Quantity: 1,
					},
				},
			},
			meliCredentials: defaultMeliCredentials,
			rootCredentials: &(*defaultMeliCredentials)[0],
			orderAnnouncementsClones: [][]announcement.Announcements{
				{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "1",
								Title:    "test-title1",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
							{
								ID:       "2",
								Title:    "test-title2",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
								Variations: []struct {
									ID                int
									AvailableQuantity int
								}{
									{
										ID:                222,
										AvailableQuantity: 1,
									},
								},
							},
						},
					},
				},
			},
			odr: &entity.Order{
				AccountID:     accountId,
				MarketplaceID: "20210101000000",
				Status:        "paid",
				Items: []entity.OrderItem{
					{
						Title:    "test-title1",
						Quantity: 1,
						Sku:      "test-sku",
					},
				},
			},
			OrderMatcher: &OrderMatcher{},
		},
		{
			name:         "default scenario where the sold item is a variation",
			accountId:    accountId,
			orderMessage: defaultOrderMessage,
			meliOrder: &common.MeliOrder{
				ID:          "20210101000000",
				DateCreated: "2022-10-30T16:19:20.129Z",
				Status:      common.Paid,
				Items: []common.OrderItem{
					{
						ID:          "1",
						Title:       "test-title1",
						Sku:         "test-sku",
						Quantity:    1,
						VariationID: 111,
					},
				},
			},
			meliCredentials: defaultMeliCredentials,
			rootCredentials: &(*defaultMeliCredentials)[0],
			orderAnnouncementsClones: [][]announcement.Announcements{
				{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "1",
								Title:    "test-title1",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
								Variations: []struct {
									ID                int
									AvailableQuantity int
								}{
									{
										ID:                111,
										AvailableQuantity: 1,
									},
								},
							},
							{
								ID:       "2",
								Title:    "test-title2",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
						},
					},
				},
			},
			odr: &entity.Order{
				AccountID:     accountId,
				MarketplaceID: "20210101000000",
				Status:        "paid",
				Items: []entity.OrderItem{
					{
						Title:       "test-title1",
						Quantity:    1,
						Sku:         "test-sku",
						VariationID: 111,
					},
				},
			},
			OrderMatcher: &OrderMatcher{},
		},
		{
			name:         "more than one item",
			accountId:    accountId,
			orderMessage: defaultOrderMessage,
			meliOrder: &common.MeliOrder{
				ID:          "20210101000000",
				DateCreated: "2022-10-30T16:19:20.129Z",
				Status:      common.Paid,
				Items: []common.OrderItem{
					{
						ID:       "1",
						Title:    "test-title",
						Sku:      "test-sku",
						Quantity: 1,
					},
					{
						ID:       "10",
						Title:    "test-title-another-item",
						Sku:      "test-sku-another-item",
						Quantity: 1,
					},
				},
			},
			meliCredentials: defaultMeliCredentials,
			rootCredentials: &(*defaultMeliCredentials)[0],
			orderAnnouncementsClones: [][]announcement.Announcements{
				{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "1",
								Title:    "test-title",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
							{
								ID:       "2",
								Title:    "test-title2",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
						},
					},
				},
				{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "10",
								Title:    "test-title-another-item",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku-another-item",
							},
							{
								ID:       "11",
								Title:    "test-title-another-item2",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku-another-item",
							},
						},
					},
				},
			},
			odr: &entity.Order{
				AccountID:     accountId,
				MarketplaceID: "20210101000000",
				Status:        "paid",
				Items: []entity.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "test-sku",
					},
					{
						Title:    "test-title-another-item",
						Quantity: 1,
						Sku:      "test-sku-another-item",
					},
				},
			},
			OrderMatcher: &OrderMatcher{},
		},
		{
			name:         "more than one item and one doesn't have an SKU",
			accountId:    accountId,
			orderMessage: defaultOrderMessage,
			meliOrder: &common.MeliOrder{
				ID:          "20210101000000",
				DateCreated: "2022-10-30T16:19:20.129Z",
				Status:      common.Paid,
				Items: []common.OrderItem{
					{
						ID:       "1",
						Title:    "test-title",
						Quantity: 1,
						Sku:      "",
					},
					{
						ID:       "10",
						Title:    "test-title-another-item",
						Quantity: 1,
						Sku:      "test-sku-another-item",
					},
				},
			},
			meliCredentials: defaultMeliCredentials,
			rootCredentials: &(*defaultMeliCredentials)[0],
			orderAnnouncementsClones: [][]announcement.Announcements{
				{},
				{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "11",
								Title:    "test-title-another-item2",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku-another-item",
							},
							{
								ID:       "10",
								Title:    "test-title-another-item",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku-another-item",
							},
						},
					},
				},
			},
			odr: &entity.Order{
				AccountID:     accountId,
				MarketplaceID: "20210101000000",
				Status:        "paid",
				Items: []entity.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "",
					},
					{
						Title:    "test-title-another-item",
						Quantity: 1,
						Sku:      "test-sku-another-item",
					},
				},
			},
			OrderMatcher: &OrderMatcher{},
		},
	}

	// -------------------------------------------------------
	// ------------------- Test Execution --------------------
	// -------------------------------------------------------
	for _, tt := range newOrderScenarios {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			tt.OrderMatcher.expected = tt.odr
			mocks.mockOrderCache.EXPECT().GetOrder(tt.orderMessage.OrderId).Return(nil, nil)
			mocks.mockOrderRepo.EXPECT().GetOrder(tt.orderMessage.OrderId).Return(nil, nil)
			mocks.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(tt.orderMessage.Store).Return(tt.meliCredentials, nil)
			mocks.mockMercadoLivre.EXPECT().FetchOrder(tt.orderMessage.OrderId, tt.rootCredentials.AccessToken).Return(tt.meliOrder, nil)
			order.RemoveDuplicateItemsTest(&tt.meliOrder.Items)

			for i, item := range tt.meliOrder.Items {
				if item.Sku == "" {
					mocks.mockLogger.EXPECT().Warn("The product doesn't have sku", zap.String("order_id", tt.orderMessage.OrderId), zap.String("announcement_id", item.ID))
					continue
				}

				mocks.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(item.Sku, tt.meliCredentials).Return(&tt.orderAnnouncementsClones[i], nil)
			}

			orderItemsIds := make([]string, len(tt.meliOrder.Items))

			for i, item := range tt.meliOrder.Items {
				orderItemsIds[i] = item.ID
			}

			for _, itmClns := range tt.orderAnnouncementsClones {
				for _, acc := range itmClns {
					var currentCredentials *store.Credentials
					for i := 0; i < len(*tt.meliCredentials); i++ {
						if acc.AccountID == (*tt.meliCredentials)[i].ID {
							currentCredentials = &(*tt.meliCredentials)[i]
						}
					}

					for _, ann := range *acc.Announcements {
						if utils.Contains(&orderItemsIds, ann.ID) {
							continue
						}

						var soldQuantity int
						for _, item := range tt.meliOrder.Items {
							if item.Sku == ann.Sku {
								soldQuantity = 1
							}
						}

						if ann.Variations != nil {
							for _, variation := range ann.Variations {
								mocks.mockAnnUseCase.EXPECT().UpdateQuantity(
									ann.ID, variation.AvailableQuantity-soldQuantity, *currentCredentials, variation.ID).Return(nil)
							}
							continue
						}
						mocks.mockAnnUseCase.EXPECT().UpdateQuantity(
							ann.ID, ann.Quantity-soldQuantity, *currentCredentials).Return(nil)
					}
				}
			}

			mocks.mockOrderRepo.EXPECT().RegisterOrder(gomock.All(tt.OrderMatcher)).Return(nil)
			mocks.mockOrderCache.EXPECT().SetOrder(gomock.All(tt.OrderMatcher)).Return(nil)
			mocks.mockOrderQueue.EXPECT().DeleteOrderNotification(tt.orderMessage.ReceiptHandle).Return(nil)

			err := orderService.ProcessOrder(tt.orderMessage)
			if err != nil {
				t.Errorf("Error processing order: %v", err)
			}
		})
	}

	// -------------------------------------------------------
	// ---------- Scenarios -> Order already exists ----------
	// -------------------------------------------------------
	orderAlreadyExistsScenarios := []struct {
		name         string
		orderMessage order.OrderMessage
		mockCall     func(m *Mocks)
		errMessage   string
	}{
		{
			name: "order already exists on cache",
			mockCall: func(m *Mocks) {
				orderStatus := entity.Paid
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(&orderStatus, nil),
					m.mockOrderQueue.EXPECT().DeleteOrderNotification(defaultOrderMessage.ReceiptHandle).Return(nil),
				)
			},
			orderMessage: defaultOrderMessage,
		},
		{
			name: "order already exists on db",
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(&entity.Order{}, nil),
					m.mockOrderQueue.EXPECT().DeleteOrderNotification(defaultOrderMessage.ReceiptHandle).Return(nil),
				)
			},
			orderMessage: defaultOrderMessage,
		},
		{
			name: "error getting order from cache",
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, errors.New("error getting order from cache")),
					m.mockLogger.EXPECT().Warn(
						"Fail to retrieve order from cache",
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.Error(errors.New("error getting order from cache")),
					),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(gomock.Any()).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(gomock.Any(), gomock.Any()).Return(&common.MeliOrder{}, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(gomock.Any(), gomock.Any()).Return(&[]announcement.Announcements{}, nil).AnyTimes(),
					m.mockAnnUseCase.EXPECT().UpdateQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
					m.mockOrderRepo.EXPECT().RegisterOrder(gomock.Any()).Return(nil),
					m.mockOrderCache.EXPECT().SetOrder(gomock.Any()).Return(nil),
					m.mockOrderQueue.EXPECT().DeleteOrderNotification(defaultOrderMessage.ReceiptHandle).Return(nil),
				)
			},
			orderMessage: defaultOrderMessage,
		},
		{
			name: "error getting order from db",
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, errors.New("error getting order from db")),
					m.mockLogger.EXPECT().Error(
						"Fail to retrieve order from the DB",
						errors.New("error getting order from db"),
						zap.String("order_id", defaultOrderMessage.OrderId),
					),
				)
			},
			orderMessage: defaultOrderMessage,
			errMessage:   "error getting order from db",
		},
	}

	// -------------------------------------------------------
	// ------------------- Test Execution --------------------
	// -------------------------------------------------------
	for _, tt := range orderAlreadyExistsScenarios {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			if tt.mockCall != nil {
				tt.mockCall(mocks)
			}

			err := orderService.ProcessOrder(tt.orderMessage)
			if (err != nil) != (tt.errMessage != "") {
				t.Errorf("Error processing order: %v", err)
			}
			if tt.errMessage != "" && err.Error() != tt.errMessage {
				t.Errorf("Wrong error message: %v", err)
			}
		})
	}

	// -------------------------------------------------------
	// --------- Scenarios -> Error retrieving data ----------
	// -------------------------------------------------------
	errRetrievingDataFromMeli := []struct {
		name         string
		orderMessage order.OrderMessage
		mockCall     func(m *Mocks)
		errMessage   string
	}{
		{
			name:         "error retrieving meli credentials",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(gomock.Any()).Return(nil, errors.New("error retrieving meli credentials")),
					m.mockLogger.EXPECT().Error(
						"Error in retrieving Meli credentials during order processing",
						errors.New("error retrieving meli credentials"),
						zap.String("store_id", defaultOrderMessage.Store),
					),
				)
			},
			errMessage: "credentials not found",
		},
		{
			name:         "error converting credentials to map",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(gomock.Any()).Return(nilCredentials, nil),
					m.mockLogger.EXPECT().Error(
						"Error in converting credentials to map",
						gomock.Any(),
						zap.String("store_id", defaultOrderMessage.Store),
					),
				)
			},
			errMessage: "credentials not found",
		},
		{
			name:         "error fetching order from meli",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(gomock.Any()).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New("error getting order from meli")),
					m.mockLogger.EXPECT().Error(
						"Error to fetch the order",
						errors.New("error getting order from meli"),
						zap.String("order_id", defaultOrderMessage.OrderId),
					),
				)
			},
			errMessage: "error processing order",
		},
		{
			name:         "error retrieving announcements and the error type doesn't match any of the expected",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				annErr := errors.New("error retrieving announcements")
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Error(
						"Error in retrieving the order product clones",
						annErr,
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
				)
			},
			errMessage: "error processing order",
		},
		{
			name:         "error retrieving announcements and unable to retry",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				annErr := &announcement.AnnouncementError{
					Message:       "error retrieving announcements",
					IsAbleToRetry: false,
					Sku:           defaultMeliOrder.Items[0].Sku,
				}
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Warn(
						"Fail in retrieving the order product clones",
						zap.Error(annErr),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockLogger.EXPECT().Error(
						"Error in retrieving the order product clones",
						annErr,
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
				)
			},
			errMessage: fmt.Sprintf("error processing order"),
		},
		{
			name:         "error retrieving announcements and able to retry - success on retry",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				annErr := &announcement.AnnouncementError{
					Message:       "error retrieving announcements",
					IsAbleToRetry: true,
					Sku:           defaultMeliOrder.Items[0].Sku,
				}

				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Warn(
						"Fail in retrieving the order product clones",
						zap.Error(annErr),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockLogger.EXPECT().Info(
						"Retrying to retrieve order products clones...",
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(&defaultMeliAnnouncementsClones, nil),

					m.mockAnnUseCase.EXPECT().UpdateQuantity(
						(*defaultMeliAnnouncementsClones[0].Announcements)[1].ID,
						(*defaultMeliAnnouncementsClones[0].Announcements)[1].Quantity-defaultMeliOrder.Items[0].Quantity,
						(*defaultMeliCredentials)[0],
					).Return(nil),
					m.mockOrderRepo.EXPECT().RegisterOrder(gomock.Any()).Return(nil),
					m.mockOrderCache.EXPECT().SetOrder(gomock.Any()).Return(nil),
					m.mockOrderQueue.EXPECT().DeleteOrderNotification(defaultOrderMessage.ReceiptHandle).Return(nil),
				)
			},
		},
		{
			name:         "error retrieving announcements and able to retry - error on retry",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				annErr := &announcement.AnnouncementError{
					Message:       "error retrieving announcements",
					IsAbleToRetry: true,
					Sku:           defaultMeliOrder.Items[0].Sku,
				}

				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Warn(
						"Fail in retrieving the order product clones",
						zap.Error(annErr),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockLogger.EXPECT().Info(
						"Retrying to retrieve order products clones...",
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Error(
						"Error in retrieving the order product clones",
						annErr,
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
				)
			},
			errMessage: fmt.Sprintf("error processing order"),
		},
	}

	// -------------------------------------------------------
	// ------------------- Test Execution --------------------
	// -------------------------------------------------------
	for _, tt := range errRetrievingDataFromMeli {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			if tt.mockCall != nil {
				tt.mockCall(mocks)
			}

			err := orderService.ProcessOrder(tt.orderMessage)
			if (err != nil) != (tt.errMessage != "") {
				t.Errorf("Error processing order: %v", err)
			}
			if tt.errMessage != "" && err.Error() != tt.errMessage {
				t.Errorf("Wrong error message: %v", err)
			}
		})
	}

	// -------------------------------------------------------
	// --------- Scenarios -> Error updating data ------------
	// -------------------------------------------------------
	// TODO: Fix this test
	errUpdatingAnnouncements := []struct {
		name         string
		orderMessage order.OrderMessage
		mockCall     func(m *Mocks)
		errMessage   string
	}{
		{
			name:         "error updating announcements",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				anns := []announcement.Announcements{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "1",
								Title:    "test-title",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
							{
								ID:       "2",
								Title:    "test-title2",
								Quantity: 1,
								Price:    1.0,
								Sku:      "test-sku",
							},
						},
					},
				}
				odrErro := &order.OrderError{
					Message: "Error updating announcements",
					AnnouncementsError: []announcement.Announcements{
						{
							AccountID: accountId,
							Announcements: &[]common.MeliAnnouncement{
								{
									ID:       (*anns[0].Announcements)[1].ID,
									Title:    (*anns[0].Announcements)[1].Title,
									Sku:      (*anns[0].Announcements)[1].Sku,
									Quantity: (*anns[0].Announcements)[1].Quantity - defaultMeliOrder.Items[0].Quantity,
								},
							},
						},
					},
				}

				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(&anns, nil),
					m.mockAnnUseCase.EXPECT().UpdateQuantity(
						(*anns[0].Announcements)[1].ID,
						(*anns[0].Announcements)[1].Quantity-defaultMeliOrder.Items[0].Quantity,
						(*defaultMeliCredentials)[0],
					).Return(errors.New("error updating announcements")),
					m.mockLogger.EXPECT().Error(
						odrErro.Message,
						odrErro,
					),
				)
			},
			errMessage: "error syncing quantities",
		},
	}

	// -------------------------------------------------------
	// ------------------- Test Execution --------------------
	// -------------------------------------------------------
	for _, tt := range errUpdatingAnnouncements {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			if tt.mockCall != nil {
				tt.mockCall(mocks)
			}

			err := orderService.ProcessOrder(tt.orderMessage)
			if (err != nil) != (tt.errMessage != "") {
				t.Errorf("Error processing order: %v", err)
			}
			if tt.errMessage != "" && err.Error() != tt.errMessage {
				t.Errorf("Wrong error message:\n Got: %v\n Expect: %v", err, tt.errMessage)
			}
		})
	}

	// -------------------------------------------------------
	// --------- Scenarios -> Error registering order --------
	// -------------------------------------------------------
	errRegisteringOrder := []struct {
		name         string
		orderMessage order.OrderMessage
		mockCall     func(m *Mocks)
		errMessage   string
	}{
		{
			name:         "error registering order",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(&[]announcement.Announcements{}, nil),
					m.mockAnnUseCase.EXPECT().UpdateQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
					m.mockOrderRepo.EXPECT().RegisterOrder(gomock.Any()).Return(errors.New("error registering order")),
					m.mockLogger.EXPECT().Error(
						"Fail to store the order",
						errors.New("error registering order"),
						zap.String("order_id", defaultOrderMessage.OrderId),
					),
				)
			},
			errMessage: "couldn't store order",
		},
		{
			name:         "error setting order on cache",
			orderMessage: defaultOrderMessage,
			mockCall: func(m *Mocks) {
				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts(defaultMeliOrder.Items[0].Sku, defaultMeliCredentials).Return(&[]announcement.Announcements{}, nil),
					m.mockAnnUseCase.EXPECT().UpdateQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
					m.mockOrderRepo.EXPECT().RegisterOrder(gomock.Any()).Return(nil),
					m.mockOrderCache.EXPECT().SetOrder(gomock.Any()).Return(errors.New("error setting order on cache")),
					m.mockLogger.EXPECT().Warn(
						"Fail to cache the order",
						zap.String("order_id", defaultOrderMessage.OrderId),
					),
					m.mockOrderQueue.EXPECT().DeleteOrderNotification(defaultOrderMessage.ReceiptHandle).Return(nil),
				)
			},
		},
	}

	// -------------------------------------------------------
	// ------------------- Test Execution --------------------
	// -------------------------------------------------------
	for _, tt := range errRegisteringOrder {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			if tt.mockCall != nil {
				tt.mockCall(mocks)
			}

			err := orderService.ProcessOrder(tt.orderMessage)
			if (err != nil) != (tt.errMessage != "") {
				t.Errorf("Error processing order: %v", err)
			}
			if tt.errMessage != "" && err.Error() != tt.errMessage {
				t.Errorf("Wrong error message:\n Got: %v\n Expect: %v", err, tt.errMessage)
			}
		})
	}
}

// TestHelperFunctions tests utility functions used in order processing.
// Current tests:
// - RemoveDuplicateItems: Verifies correct handling of:
//   - Duplicate items with same SKU
//   - Items with different SKUs
//   - Items with empty SKUs
//   - Mix of items with and without SKUs
func TestHelperFunctions(t *testing.T) {
	t.Run("test remove duplicate items", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []common.OrderItem
			expected []common.OrderItem
		}{
			{
				name: "duplicate items with same SKU",
				input: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "test-sku",
					},
					{
						Title:    "test-title",
						Quantity: 2,
						Sku:      "test-sku",
					},
				},
				expected: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 3,
						Sku:      "test-sku",
					},
				},
			},
			{
				name: "items with different SKUs",
				input: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "test-sku-1",
					},
					{
						Title:    "test-title",
						Quantity: 2,
						Sku:      "test-sku-2",
					},
				},
				expected: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "test-sku-1",
					},
					{
						Title:    "test-title",
						Quantity: 2,
						Sku:      "test-sku-2",
					},
				},
			},
			{
				name: "items with empty SKUs",
				input: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "",
					},
					{
						Title:    "test-title",
						Quantity: 2,
						Sku:      "",
					},
				},
				expected: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "",
					},
					{
						Title:    "test-title",
						Quantity: 2,
						Sku:      "",
					},
				},
			},
			{
				name: "mix of items with and without SKUs",
				input: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 1,
						Sku:      "test-sku",
					},
					{
						Title:    "test-title",
						Quantity: 2,
						Sku:      "test-sku",
					},
					{
						Title:    "test-title",
						Quantity: 3,
						Sku:      "",
					},
				},
				expected: []common.OrderItem{
					{
						Title:    "test-title",
						Quantity: 3,
						Sku:      "test-sku",
					},
					{
						Title:    "test-title",
						Quantity: 3,
						Sku:      "",
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				items := tt.input
				order.RemoveDuplicateItemsTest(&items)
				if !cmp.Equal(items, tt.expected) {
					t.Errorf("RemoveDuplicateItems() = %v, want %v", items, tt.expected)
				}
			})
		}
	})
}

// ============================ Custom matchers ============================

// OrderMatcher is a custom gomock matcher for Order entities.
// It compares orders while ignoring specific fields (ID, AccountID, MarketplaceID, DateCreated).
type OrderMatcher struct {
	expected *entity.Order
}

func (o *OrderMatcher) Matches(x interface{}) bool {
	order, ok := x.(*entity.Order)
	if !ok {
		return false
	}

	return cmp.Equal(order, o.expected, cmpopts.IgnoreFields(entity.Order{}, "ID", "AccountID", "MarketplaceID", "DateCreated"), cmpopts.IgnoreFields(entity.OrderItem{}, "ID"))
}

func (o *OrderMatcher) String() string {
	return fmt.Sprintf("matches order %v", o.expected)
}

// ================================== Mocks =================================

// Mocks holds all mock instances used in testing.
// Provides a centralized way to manage test dependencies.
type Mocks struct {
	mockOrderQueue   *mock_order.MockQueue
	mockMercadoLivre *common_mock.MockMercadoLivre
	mockStoreUseCase *mock_store.MockUseCase
	mockAnnUseCase   *mock_announcement.MockUseCase
	mockOrderRepo    *mock_order.MockRepository
	mockOrderCache   *mock_order.MockCache
	mockLogger       *common_mock.MockLogger
}

func newMocks(ctrl *gomock.Controller) *Mocks {
	return &Mocks{
		mockOrderQueue:   mock_order.NewMockQueue(ctrl),
		mockMercadoLivre: common_mock.NewMockMercadoLivre(ctrl),
		mockStoreUseCase: mock_store.NewMockUseCase(ctrl),
		mockAnnUseCase:   mock_announcement.NewMockUseCase(ctrl),
		mockOrderRepo:    mock_order.NewMockRepository(ctrl),
		mockOrderCache:   mock_order.NewMockCache(ctrl),
		mockLogger:       common_mock.NewMockLogger(ctrl),
	}
}

func (m *Mocks) newOrderService() *order.OrderService {
	return order.NewOrderService(
		m.mockOrderQueue,
		m.mockMercadoLivre,
		m.mockStoreUseCase,
		m.mockAnnUseCase,
		m.mockOrderRepo,
		m.mockOrderCache,
		m.mockLogger,
	)
}

// TestProcessOrderErrors specifically tests error scenarios in order processing.
// Covers errors in:
// 1. Credential retrieval
// 2. Mercado Livre API calls
// 3. Announcement retrieval
// 4. Quantity updates
// 5. Order registration
func TestProcessOrderErrors(t *testing.T) {
	accountId := entity.ID(uuid.New())
	defaultOrderMessage := order.OrderMessage{
		Store:         "1",
		OrderId:       "20210101000000",
		Attempts:      0,
		ReceiptHandle: "test-receipt-handle",
	}

	defaultMeliCredentials := &[]store.Credentials{
		{
			ID: accountId,
			MeliCredential: &common.MeliCredential{
				AccessToken: "test-access-token",
				UserID:      "1",
			},
		},
	}

	tests := []struct {
		name       string
		setupMocks func(*Mocks)
		wantErr    bool
		errMessage string
	}{
		{
			name: "error getting credentials",
			setupMocks: func(m *Mocks) {
				m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(nil, errors.New("credentials error"))
				m.mockLogger.EXPECT().Error("Error in retrieving Meli credentials during order processing",
					gomock.Any(),
					zap.String("store_id", defaultOrderMessage.Store))
			},
			wantErr:    true,
			errMessage: "credentials not found",
		},
		{
			name: "error fetching order from Mercado Livre",
			setupMocks: func(m *Mocks) {
				m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil)
				m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(nil, errors.New("meli error"))
				m.mockLogger.EXPECT().Error("Error to fetch the order",
					gomock.Any(),
					zap.String("order_id", defaultOrderMessage.OrderId))
			},
			wantErr:    true,
			errMessage: "error processing order",
		},
		{
			name: "error retrieving announcements",
			setupMocks: func(m *Mocks) {
				meliOrder := &common.MeliOrder{
					ID:          defaultOrderMessage.OrderId,
					DateCreated: "2022-10-30T16:19:20.129Z",
					Status:      common.Paid,
					Items: []common.OrderItem{
						{
							ID:       "1",
							Title:    "test-title",
							Sku:      "test-sku",
							Quantity: 1,
						},
					},
				}

				m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil)
				m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(meliOrder, nil)
				m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts("test-sku", defaultMeliCredentials).Return(nil, errors.New("announcement error"))
				m.mockLogger.EXPECT().Error("Error in retrieving the order product clones",
					gomock.Any(),
					zap.String("sku", "test-sku"))
			},
			wantErr:    true,
			errMessage: "error processing order",
		},
		{
			name: "error updating quantity",
			setupMocks: func(m *Mocks) {
				meliOrder := &common.MeliOrder{
					ID:          defaultOrderMessage.OrderId,
					DateCreated: "2022-10-30T16:19:20.129Z",
					Status:      common.Paid,
					Items: []common.OrderItem{
						{
							ID:       "1",
							Title:    "test-title",
							Sku:      "test-sku",
							Quantity: 1,
						},
					},
				}

				announcements := &[]announcement.Announcements{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "2",
								Title:    "test-title",
								Sku:      "test-sku",
								Quantity: 2,
							},
						},
					},
				}

				m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil)
				m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(meliOrder, nil)
				m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts("test-sku", defaultMeliCredentials).Return(announcements, nil)
				m.mockAnnUseCase.EXPECT().UpdateQuantity("2", 1, (*defaultMeliCredentials)[0]).Return(errors.New("update error"))
				m.mockLogger.EXPECT().Error("Error updating announcements",
					gomock.Any(),
				)
			},
			wantErr:    true,
			errMessage: "error syncing quantities",
		},
		{
			name: "error registering order",
			setupMocks: func(m *Mocks) {
				meliOrder := &common.MeliOrder{
					ID:          defaultOrderMessage.OrderId,
					DateCreated: "2022-10-30T16:19:20.129Z",
					Status:      common.Paid,
					Items: []common.OrderItem{
						{
							ID:       "1",
							Title:    "test-title",
							Sku:      "test-sku",
							Quantity: 1,
						},
					},
				}

				announcements := &[]announcement.Announcements{
					{
						AccountID: accountId,
						Announcements: &[]common.MeliAnnouncement{
							{
								ID:       "2",
								Title:    "test-title",
								Sku:      "test-sku",
								Quantity: 2,
							},
						},
					},
				}

				m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil)
				m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil)
				m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, (*defaultMeliCredentials)[0].AccessToken).Return(meliOrder, nil)
				m.mockAnnUseCase.EXPECT().RetrieveAnnouncementsFromAllAccounts("test-sku", defaultMeliCredentials).Return(announcements, nil)
				m.mockAnnUseCase.EXPECT().UpdateQuantity("2", 1, (*defaultMeliCredentials)[0]).Return(nil)
				m.mockOrderRepo.EXPECT().RegisterOrder(gomock.Any()).Return(errors.New("register error"))
				m.mockLogger.EXPECT().Error("Fail to store the order",
					gomock.Any(),
					zap.String("order_id", defaultOrderMessage.OrderId))
			},
			wantErr:    true,
			errMessage: "couldn't store order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			tt.setupMocks(mocks)

			err := orderService.ProcessOrder(defaultOrderMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMessage {
				t.Errorf("ProcessOrder() error = %v, want %v", err, tt.errMessage)
			}
		})
	}
}
