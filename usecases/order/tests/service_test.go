package order

import (
	"testing"

	"github.com/Vractos/dolly/entity"
	mock_announcement "github.com/Vractos/dolly/usecases/announcement/mock"
	common "github.com/Vractos/dolly/usecases/common"
	common_mock "github.com/Vractos/dolly/usecases/common/mock"
	"github.com/Vractos/dolly/usecases/order"
	"github.com/Vractos/dolly/usecases/order/mock"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/Vractos/dolly/usecases/store/mock"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

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

		mockOrderQueue.EXPECT().DeleteOrderNotification(orderMessage.ReceiptHandle).Return(nil)
		mockOrderCache.EXPECT().GetOrder(orderMessage.OrderId).Return(nil, nil)
		mockOrderRepo.EXPECT().GetOrder(orderMessage.OrderId).Return(nil, nil)
		mockStoreUseCase.EXPECT().RetrieveMeliCredentialsFromMeliUserID(orderMessage.Store).Return(meliCredentials, nil)
		mockMercadoLivre.EXPECT().FetchOrder(orderMessage.OrderId, meliCredentials.MeliAccessToken).Return(&common.MeliOrder{}, nil)

		err := orderService.ProcessOrder(input)
		if err != nil {
			t.Errorf("Error processing order: %v", err)
		}
	})

}
