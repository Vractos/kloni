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
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestProcessWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := newMocks(ctrl)
	orderService := mocks.newOrderService()

	notification := order.OrderWebhookDtoInput{
		ID:            "test-id",
		Resource:      "/orders/20210101000000",
		UserID:        1,
		Attempts:      0,
		Topic:         "orders_v2",
		ApplicationID: 1,
		Received:      "2022-10-30T16:19:20.129Z",
		Sent:          "2022-10-30T16:19:20.106Z",
	}

	t.Run("process webhook", func(t *testing.T) {
		mocks.mockOrderQueue.EXPECT().PostOrderNotification(notification).Return(nil)

		err := orderService.ProcessWebhook(notification)
		if err != nil {
			t.Errorf("Error processing webhook: %v", err)
		}
	})

	t.Run("process webhook - error posting to queue", func(t *testing.T) {
		mocks.mockOrderQueue.EXPECT().PostOrderNotification(notification).Return(errors.New("error posting to queue"))
		mocks.mockLogger.EXPECT().Error(
			"Error to post order notification",
			errors.New("error posting to queue"),
			zap.String("notification_id", notification.ID),
			zap.Int("user_id", notification.UserID),
			zap.Int("attempts", notification.Attempts),
			zap.String("sent", notification.Sent),
		)

		err := orderService.ProcessWebhook(notification)
		if err == nil {
			t.Errorf("Error testing error on processing webhook: %v", err)
		}
		if err.Error() != "error to post order notification" {
			t.Errorf("Wrong error message: %v", err)
		}
	})
}

func TestProcessOrder(t *testing.T) {
	accountId := entity.ID(uuid.New())
	defaultOrderMessage := order.OrderMessage{
		Store:         "1",
		OrderId:       "20210101000000",
		Attempts:      0,
		ReceiptHandle: "test-receipt-handle",
	}
	defaultMeliCredentials := &store.Credentials{
		StoreID:         accountId,
		MeliAccessToken: "test-access-token",
		MeliUserID:      "1",
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
	defaultMeliAnnouncementsClones := [][]common.MeliAnnouncement{
		{
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
	}

	newOrderScenarios := []struct {
		name                     string
		accountId                entity.ID
		orderMessage             order.OrderMessage
		meliOrder                *common.MeliOrder
		meliCredentials          *store.Credentials
		orderAnnouncementsClones [][]common.MeliAnnouncement
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
			orderAnnouncementsClones: [][]common.MeliAnnouncement{
				{
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
			orderAnnouncementsClones: [][]common.MeliAnnouncement{
				{
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
			orderAnnouncementsClones: [][]common.MeliAnnouncement{
				{
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
			orderAnnouncementsClones: [][]common.MeliAnnouncement{
				{
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
				{
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
			orderAnnouncementsClones: [][]common.MeliAnnouncement{
				{},
				{
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

	for _, tt := range newOrderScenarios {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mocks := newMocks(ctrl)
			orderService := mocks.newOrderService()

			tt.OrderMatcher.expected = tt.odr
			mocks.mockOrderCache.EXPECT().GetOrder(tt.orderMessage.OrderId).Return(nil, nil)
			mocks.mockOrderRepo.EXPECT().GetOrder(tt.orderMessage.OrderId).Return(nil, nil)
			mocks.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(tt.orderMessage.Store).Return(tt.meliCredentials, nil)
			mocks.mockMercadoLivre.EXPECT().FetchOrder(tt.orderMessage.OrderId, tt.meliCredentials.MeliAccessToken).Return(tt.meliOrder, nil)
			order.RemoveDuplicateItensTest(&tt.meliOrder.Items)
			for i, item := range tt.meliOrder.Items {
				if item.Sku == "" {
					mocks.mockLogger.EXPECT().Warn("The product doesn't have sku", zap.String("order_id", tt.orderMessage.OrderId), zap.String("announcement_id", item.ID))
					continue
				}
				mocks.mockAnnUseCase.EXPECT().RetrieveAnnouncements(item.Sku, *tt.meliCredentials).Return(&tt.orderAnnouncementsClones[i], nil)
			}
			for _, anns := range tt.orderAnnouncementsClones {
				for i := 1; i < len(anns); i++ {
					if anns[i].Variations != nil {
						for _, variation := range anns[i].Variations {
							mocks.mockAnnUseCase.EXPECT().UpdateQuantity(
								anns[i].ID, variation.AvailableQuantity-anns[0].Quantity, *tt.meliCredentials, variation.ID).Return(nil)
						}
						continue
					}
					mocks.mockAnnUseCase.EXPECT().UpdateQuantity(
						anns[i].ID, anns[i].Quantity-anns[0].Quantity, *tt.meliCredentials).Return(nil)
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
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(gomock.Any()).Return(&store.Credentials{}, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(gomock.Any(), gomock.Any()).Return(&common.MeliOrder{}, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(gomock.Any(), gomock.Any()).Return(&[]common.MeliAnnouncement{}, nil).AnyTimes(),
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
						zap.String("order_id", defaultOrderMessage.OrderId),
					),
				)
			},
			errMessage: "error to process order - get credentials",
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
			errMessage: "error getting order from meli",
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
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, defaultMeliCredentials.MeliAccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Error(
						"Error in retrieving the order product clones",
						annErr,
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
				)
			},
			errMessage: "error retrieving announcements",
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
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, defaultMeliCredentials.MeliAccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Warn(
						"Fail in retrieving the order product clones",
						zap.Error(annErr),
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockLogger.EXPECT().Error(
						"Error in retrieving the order product clones",
						annErr,
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
				)
			},
			errMessage: fmt.Sprintf("Message: error retrieving announcements - Announcement SKU: %s", defaultMeliOrder.Items[0].Sku),
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
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, defaultMeliCredentials.MeliAccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Warn(
						"Fail in retrieving the order product clones",
						zap.Error(annErr),
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockLogger.EXPECT().Info(
						"Retrying to retrieve order products clones...",
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(&defaultMeliAnnouncementsClones[0], nil),

					m.mockAnnUseCase.EXPECT().UpdateQuantity(
						defaultMeliAnnouncementsClones[0][1].ID,
						defaultMeliAnnouncementsClones[0][1].Quantity-defaultMeliOrder.Items[0].Quantity,
						*defaultMeliCredentials,
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
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, defaultMeliCredentials.MeliAccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Warn(
						"Fail in retrieving the order product clones",
						zap.Error(annErr),
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockLogger.EXPECT().Info(
						"Retrying to retrieve order products clones...",
						zap.String("order_id", defaultOrderMessage.OrderId),
						zap.String("sku", defaultMeliOrder.Items[0].Sku),
					),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(nil, annErr),
					m.mockLogger.EXPECT().Error(
						"Error in retrieving the order product clones",
						annErr,
						zap.String("order_id", defaultOrderMessage.OrderId),
					),
				)
			},
			errMessage: fmt.Sprintf("Message: error retrieving announcements - Announcement SKU: %s", defaultMeliOrder.Items[0].Sku),
		},
	}

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
				anns := []common.MeliAnnouncement{
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
				}
				odrErro := &order.OrderError{
					Message: "Couldn't change the announcements",
					AnnouncementsError: []common.MeliAnnouncement{
						{
							ID:       anns[1].ID,
							Title:    anns[1].Title,
							Sku:      anns[1].Sku,
							Quantity: anns[1].Quantity - defaultMeliOrder.Items[0].Quantity,
						},
					},
				}

				gomock.InOrder(
					m.mockOrderCache.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockOrderRepo.EXPECT().GetOrder(defaultOrderMessage.OrderId).Return(nil, nil),
					m.mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(defaultOrderMessage.Store).Return(defaultMeliCredentials, nil),
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, defaultMeliCredentials.MeliAccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(&anns, nil),
					m.mockAnnUseCase.EXPECT().UpdateQuantity(
						anns[1].ID,
						anns[1].Quantity-defaultMeliOrder.Items[0].Quantity,
						*defaultMeliCredentials,
					).Return(errors.New("error updating announcements")),
					m.mockLogger.EXPECT().Error(
						"Couldn't change the announcements",
						odrErro,
						zap.Any(
							"Announcements that were not possible to change the quantity",
							odrErro.AnnouncementsError,
						),
					),
				)
			},
			errMessage: "Couldn't change the announcements",
		},
	}

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
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, defaultMeliCredentials.MeliAccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(&[]common.MeliAnnouncement{}, nil),
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
					m.mockMercadoLivre.EXPECT().FetchOrder(defaultOrderMessage.OrderId, defaultMeliCredentials.MeliAccessToken).Return(defaultMeliOrder, nil),
					m.mockAnnUseCase.EXPECT().RetrieveAnnouncements(defaultMeliOrder.Items[0].Sku, *defaultMeliCredentials).Return(&[]common.MeliAnnouncement{}, nil),
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

func TestSupportingFuncs(t *testing.T) {
	t.Run("test remove duplicated items", func(t *testing.T) {
		items := &[]common.OrderItem{
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
				Quantity: 1,
				Sku:      "test-sku-other",
			},
		}
		expected := &[]common.OrderItem{
			{
				Title:    "test-title",
				Quantity: 3,
				Sku:      "test-sku",
			},
			{
				Title:    "test-title",
				Quantity: 1,
				Sku:      "test-sku-other",
			},
		}

		order.RemoveDuplicateItensTest(items)

		if len(*items) != 2 {
			t.Errorf("Error removing duplicated items: %v", items)
		}
		if !cmp.Equal(items, expected) {
			t.Errorf("Error removing duplicated items, diff:\n%v", cmp.Diff(items, expected))
		}
	})

	t.Run("test remove duplicated items - items don't have an SKU", func(t *testing.T) {
		items := &[]common.OrderItem{
			{
				Title:    "test-title",
				Quantity: 1,
				Sku:      "",
			},
			{
				Title:    "test-title-3",
				Quantity: 3,
				Sku:      "",
			},
			{
				Title:    "test-title-2",
				Quantity: 1,
				Sku:      "test-sku",
			},
		}
		expected := &[]common.OrderItem{
			{
				Title:    "test-title",
				Quantity: 1,
				Sku:      "",
			},
			{
				Title:    "test-title-3",
				Quantity: 3,
				Sku:      "",
			},
			{
				Title:    "test-title-2",
				Quantity: 1,
				Sku:      "test-sku",
			},
		}

		order.RemoveDuplicateItensTest(items)

		if len(*items) != 3 {
			t.Errorf("Error removing duplicated items: %v", items)
		}
		if !cmp.Equal(items, expected) {
			t.Errorf("Error removing duplicated items, diff:\n%v", cmp.Diff(items, expected))
		}
	})
}

// ============================ Custom matchers ============================
type OrderMatcher struct {
	expected *entity.Order
}

func (o *OrderMatcher) Matches(x interface{}) bool {
	order, ok := x.(*entity.Order)
	if !ok {
		return false
	}

	return cmp.Equal(order, o.expected, cmpopts.IgnoreFields(entity.Order{}, "ID", "StoreID", "DateCreated"), cmpopts.IgnoreFields(entity.OrderItem{}, "ID"))
}

func (o *OrderMatcher) String() string {
	return fmt.Sprintf("matches order %v", o.expected)
}

// ================================== Mocks =================================
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
