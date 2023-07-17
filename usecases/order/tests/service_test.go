package order

import (
	"testing"

	mock_announcement "github.com/Vractos/dolly/usecases/announcement/mock"
	common_mock "github.com/Vractos/dolly/usecases/common/mock"
	"github.com/Vractos/dolly/usecases/order"
	"github.com/Vractos/dolly/usecases/order/mock"
	"github.com/Vractos/dolly/usecases/store/mock"
	gomock "github.com/golang/mock/gomock"
)

func TestOrderUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockOrderQueue := mock_order.NewMockQueue(ctrl)
	mockMercadoLivre := common_mock.NewMockMercadoLivre(ctrl)
	mockStoreUseCase := mock_store.NewMockUseCase(ctrl)
	mockAnnUseCase := mock_announcement.NewMockUseCase()
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

	t.Run("")
}
