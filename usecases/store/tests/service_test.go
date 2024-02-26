package store

import (
	"testing"
	"time"

	"github.com/Vractos/kloni/entity"
	common "github.com/Vractos/kloni/usecases/common"
	common_mock "github.com/Vractos/kloni/usecases/common/mock"
	"github.com/Vractos/kloni/usecases/store"
	mock_store "github.com/Vractos/kloni/usecases/store/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestStoreUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStoreRepo := mock_store.NewMockRepository(ctrl)
	mockMercadoLivre := common_mock.NewMockMercadoLivre(ctrl)
	mockLogger := common_mock.NewMockLogger(ctrl)

	storeService := store.NewStoreService(mockStoreRepo, mockMercadoLivre, mockLogger)

	t.Run("register store", func(t *testing.T) {
		storeInput := store.RegisterStoreDtoInput{
			Name:  "Test Store",
			Email: "store@teststore.xyz",
		}

		str, err := entity.NewStore(storeInput.Email, storeInput.Name)
		if err != nil {
			t.Errorf("Error creating store instance: %v", err)
		}

		mockStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(str)).Return(str.ID, nil)

		store, err := storeService.RegisterStore(storeInput)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if store == uuid.Nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("register meli credentials", func(t *testing.T) {
		inputMeliCredentials := store.RegisterMeliCredentialsDtoInput{
			Store: entity.ID(uuid.New()),
			Code:  "test-code",
		}

		mockMercadoLivre.EXPECT().RegisterCredential(inputMeliCredentials.Code).Return(&common.MeliCredential{}, nil)
		mockStoreRepo.EXPECT().RegisterMeliCredential(inputMeliCredentials.Store, gomock.AssignableToTypeOf(&common.MeliCredential{})).Return(nil)

		err := storeService.RegisterMeliCredentials(inputMeliCredentials)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("retrieve meli credentials from store id", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		credentials := &common.MeliCredential{
			AccessToken:  "test-access-token",
			ExpiresIn:    21666,
			UserID:       "test-user-id",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &store.Credentials{
			StoreID:         storeId,
			MeliAccessToken: "test-access-token",
			MeliUserID:      "test-user-id",
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromStoreID(storeId).Return(credentials, nil)

		cred, err := storeService.RetrieveMeliCredentialsFromStoreID(storeId)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from store id with expired credentials", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		credentials := &common.MeliCredential{
			AccessToken:  "test-access-token",
			ExpiresIn:    21666,
			UserID:       "test-user-id",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		refreshedCredentials := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed",
			ExpiresIn:    21666,
			UserID:       credentials.UserID,
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &store.Credentials{
			StoreID:         storeId,
			MeliAccessToken: "test-access-token-refreshed",
			MeliUserID:      "test-user-id",
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromStoreID(storeId).Return(credentials, nil)
		mockMercadoLivre.EXPECT().RefreshCredentials(credentials.RefreshToken).Return(&refreshedCredentials, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(storeId, &refreshedCredentials)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("store_id", storeId.String()))

		cred, err := storeService.RetrieveMeliCredentialsFromStoreID(storeId)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from meli user id", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		id := "test-user-id"
		credentials := &common.MeliCredential{
			AccessToken:  "test-access-token",
			ExpiresIn:    21666,
			UserID:       id,
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &store.Credentials{
			StoreID:         storeId,
			MeliAccessToken: "test-access-token",
			MeliUserID:      id,
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromMeliUserID(id).Return(&storeId, credentials, nil)

		cred, err := storeService.RetrieveMeliCredentialsFromMeliUserID(id)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from meli user id with expired credentials", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		id := "test-user-id"
		credentials := &common.MeliCredential{
			AccessToken:  "test-access-token",
			ExpiresIn:    21666,
			UserID:       id,
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		refreshedCredentials := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed",
			ExpiresIn:    21666,
			UserID:       id,
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &store.Credentials{
			StoreID:         storeId,
			MeliAccessToken: "test-access-token-refreshed",
			MeliUserID:      id,
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromMeliUserID(id).Return(&storeId, credentials, nil)
		mockMercadoLivre.EXPECT().RefreshCredentials(credentials.RefreshToken).Return(&refreshedCredentials, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(storeId, &refreshedCredentials)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("store_id", storeId.String()))

		cred, err := storeService.RetrieveMeliCredentialsFromMeliUserID(id)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("refresh meli credentials", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		refreshToken := "test-refresh-token"

		refreshedCredentials := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed",
			ExpiresIn:    21666,
			UserID:       "test-user-id",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &store.Credentials{
			StoreID:         storeId,
			MeliAccessToken: "test-access-token-refreshed",
			MeliUserID:      "test-user-id",
		}

		mockMercadoLivre.EXPECT().RefreshCredentials(refreshToken).Return(&refreshedCredentials, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(storeId, &refreshedCredentials)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("store_id", storeId.String()))

		cred, err := storeService.RefreshMeliCredential(storeId, refreshToken)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})
}
