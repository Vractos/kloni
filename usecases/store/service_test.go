package store

import (
	"testing"
	"time"

	"github.com/Vractos/dolly/entity"
	common "github.com/Vractos/dolly/usecases/common"
	common_mock "github.com/Vractos/dolly/usecases/common/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestStoreUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStoreRepo := NewMockRepository(ctrl)
	mockMercadoLivre := common_mock.NewMockMercadoLivre(ctrl)
	mockLogger := common_mock.NewMockLogger(ctrl)
	mockUseCase := NewMockUseCase(ctrl)

	t.Run("register store", func(t *testing.T) {
		storeInput := RegisterStoreDtoInput{
			Name:  "Test Store",
			Email: "store@teststore.xyz",
		}

		str, err := entity.NewStore(storeInput.Email, storeInput.Name)
		if err != nil {
			t.Errorf("Error creating store instance: %v", err)
		}

		mockStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(str)).Return(str.ID, nil)

		storeService := NewStoreService(mockStoreRepo, mockMercadoLivre, mockLogger)
		store, err := storeService.RegisterStore(storeInput)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if store == uuid.Nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("register meli credentials", func(t *testing.T) {
		inputMeliCredentials := RegisterMeliCredentialsDtoInput{
			Store: entity.ID(uuid.New()),
			Code:  "test-code",
		}

		mockMercadoLivre.EXPECT().RegisterCredential(inputMeliCredentials.Code).Return(&common.MeliCredential{}, nil)
		mockStoreRepo.EXPECT().RegisterMeliCredential(inputMeliCredentials.Store, gomock.AssignableToTypeOf(&common.MeliCredential{})).Return(nil)

		storeService := NewStoreService(mockStoreRepo, mockMercadoLivre, mockLogger)
		err := storeService.RegisterMeliCredentials(inputMeliCredentials)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("retrieve meli credentials from store id", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		credentials := common.MeliCredential{
			AccessToken:  "test-access-token",
			ExpiresIn:    21600,
			UserID:       "test-user-id",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now(),
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromStoreID(storeId).Return(&credentials, nil)
		mockUseCase.EXPECT().RefreshMeliCredential(storeId, credentials.RefreshToken).Return(&Credentials{}, nil).Times(0)

		storeService := NewStoreService(mockStoreRepo, mockMercadoLivre, mockLogger)
		_, err := storeService.RetrieveMeliCredentialsFromStoreID(storeId)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

	})
}
