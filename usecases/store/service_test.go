package store

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestRegisterStore(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStoreRepo := NewMockRepository(ctrl)
	mockStoreRepo.EXPECT

}
