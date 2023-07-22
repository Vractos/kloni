package order

import (
	"fmt"
	"testing"

	"github.com/Vractos/dolly/entity"
	mock_announcement "github.com/Vractos/dolly/usecases/announcement/mock"
	common "github.com/Vractos/dolly/usecases/common"
	common_mock "github.com/Vractos/dolly/usecases/common/mock"
	"github.com/Vractos/dolly/usecases/order"
	mock_order "github.com/Vractos/dolly/usecases/order/mock"
	"github.com/Vractos/dolly/usecases/store"
	mock_store "github.com/Vractos/dolly/usecases/store/mock"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

type OrderMatcher struct {
	expected *entity.Order
}

func (o *OrderMatcher) Matches(x interface{}) bool {
	order, ok := x.(*entity.Order)
	if !ok {
		return false
	}

	return cmp.Equal(order, o.expected, cmpopts.IgnoreFields(entity.Order{}, "ID", "DateCreated"), cmpopts.IgnoreFields(entity.OrderItem{}, "ID"))
}

func (o *OrderMatcher) String() string {
	return fmt.Sprintf("matches order %v", o.expected)
}

func TestOrderUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockOrderQueue := mock_order.NewMockQueue(ctrl)
	mockMercadoLivre := common_mock.NewMockMercadoLivre(ctrl)
	mockStoreUseCase := mock_store.NewMockUseCase(ctrl)
	mockAnnUseCase := mock_announcement.NewMockUseCase(ctrl)
	mockOrderRepo := mock_order.NewMockRepository(ctrl)
	mockOrderCache := mock_order.NewMockCache(ctrl)
	mockLogger := common_mock.NewMockLogger(ctrl)

	orderService := order.NewOrderService(
		mockOrderQueue,
		mockMercadoLivre,
		mockStoreUseCase,
		mockAnnUseCase,
		mockOrderRepo,
		mockOrderCache,
		mockLogger,
	)

	t.Run("process webhook", func(t *testing.T) {
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

		mockOrderQueue.EXPECT().PostOrderNotification(notification).Return(nil)

		err := orderService.ProcessWebhook(notification)
		if err != nil {
			t.Errorf("Error processing webhook: %v", err)
		}
	})

	t.Run("process a new order", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		orderMessage := order.OrderMessage{
			Store:         "1",
			OrderId:       "20210101000000",
			Attempts:      0,
			ReceiptHandle: "test-receipt-handle",
		}
		meliCredentials := &store.Credentials{
			StoreID:         storeId,
			MeliAccessToken: "test-access-token",
			MeliUserID:      "1",
		}
		meliOrder := &common.MeliOrder{
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

		// Including the announcement from the order
		orderAnnouncementsClones := []common.MeliAnnouncement{
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
			{
				ID:       "3",
				Title:    "test-title3",
				Quantity: 1,
				Price:    1.0,
				Sku:      "test-sku",
			},
		}

		odr := &entity.Order{
			StoreID:       storeId,
			MarketplaceID: "20210101000000",
			Status:        "paid",
			Items: []entity.OrderItem{
				{
					Title:    "test-title",
					Quantity: 1,
					Sku:      "test-sku",
				},
			},
		}
		odrMatcher := &OrderMatcher{expected: odr}

		mockOrderCache.EXPECT().GetOrder(orderMessage.OrderId).Return(nil, nil)
		mockOrderRepo.EXPECT().GetOrder(orderMessage.OrderId).Return(nil, nil)
		mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(orderMessage.Store).Return(meliCredentials, nil)
		mockMercadoLivre.EXPECT().FetchOrder(orderMessage.OrderId, meliCredentials.MeliAccessToken).Return(meliOrder, nil)
		mockAnnUseCase.EXPECT().RetrieveAnnouncements(meliOrder.Items[0].Sku, *meliCredentials).Return(&orderAnnouncementsClones, nil)
		for i := 1; i < len(orderAnnouncementsClones); i++ {
			mockAnnUseCase.EXPECT().UpdateQuantity(
				orderAnnouncementsClones[i].ID, 0, *meliCredentials).Return(nil)
		}
		mockOrderRepo.EXPECT().RegisterOrder(gomock.All(odrMatcher)).Return(nil)
		mockOrderCache.EXPECT().SetOrder(gomock.All(odrMatcher)).Return(nil)

		mockOrderQueue.EXPECT().DeleteOrderNotification(orderMessage.ReceiptHandle).Return(nil)

		err := orderService.ProcessOrder(orderMessage)
		if err != nil {
			t.Errorf("Error processing order: %v", err)
		}
	})
}

// func TestOrderUseCaseErros(t *testing.T) {
//   ctrl := gomock.NewController(t)

// 	mockOrderQueue := mock_order.NewMockQueue(ctrl)
// 	mockMercadoLivre := common_mock.NewMockMercadoLivre(ctrl)
// 	mockStoreUseCase := mock_store.NewMockUseCase(ctrl)
// 	mockAnnUseCase := mock_announcement.NewMockUseCase(ctrl)
// 	mockOrderRepo := mock_order.NewMockRepository(ctrl)
// 	mockOrderCache := mock_order.NewMockCache(ctrl)
// 	mockLogger := common_mock.NewMockLogger(ctrl)

// 	orderService := order.NewOrderService(
// 		mockOrderQueue,
// 		mockMercadoLivre,
// 		mockStoreUseCase,
// 		mockAnnUseCase,
// 		mockOrderRepo,
// 		mockOrderCache,
// 		mockLogger,
// 	)
//   t.Run("process a new order - error to retrieve credentials", func(t *testing.T) {
//     orderMessage := order.OrderMessage{
//       Store:         "1",
//       OrderId:       "20210101000000",
//       Attempts:      0,
//       ReceiptHandle: "test-receipt-handle",
//     }

//     mockOrderCache.EXPECT().GetOrder(orderMessage.OrderId).Return(nil, nil)
//     mockOrderRepo.EXPECT().GetOrder(orderMessage.OrderId).Return(nil, nil)
//   }

// }

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
}
