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
		// credentialID -> Can have any value but must be a valid entity.ID
		credentialID := entity.ID(uuid.New())
		inputMeliCredentials := store.RegisterMeliCredentialsDtoInput{
			Store:       entity.ID(uuid.New()),
			Code:        "test-code",
			AccountName: "Test Account",
		}

		mockMercadoLivre.EXPECT().RegisterCredential(inputMeliCredentials.Code).Return(&common.MeliCredential{}, nil)
		mockStoreRepo.EXPECT().RegisterMeliCredential(gomock.AssignableToTypeOf(credentialID), inputMeliCredentials.Store, gomock.AssignableToTypeOf(&common.MeliCredential{}), inputMeliCredentials.AccountName).Return(nil)

		err := storeService.RegisterMeliCredentials(inputMeliCredentials)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("retrieve meli credentials from store id - one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId := entity.ID(uuid.New())
		accountName := "test-account-name"

		meliCredentials := &common.MeliCredential{
			UserID:       "test-user-id",
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId,
				AccountName:    &accountName,
				MeliCredential: meliCredentials,
			},
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId,
				AccountName: &accountName,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token",
					UserID:      "test-user-id",
				},
			},
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

	t.Run("retrieve meli credentials from store id - more than one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId1, accountId2 := entity.ID(uuid.New()), entity.ID(uuid.New())

		var (
			accountName1 = "test-account-name-1"
			accountName2 = "test-account-name-2"
		)

		meliCredentials1 := &common.MeliCredential{
			UserID:       "test-user-id-1",
			AccessToken:  "test-access-token-1",
			RefreshToken: "test-refresh-token-1",
			UpdatedAt:    time.Now().UTC(),
		}

		meliCredentials2 := &common.MeliCredential{
			UserID:       "test-user-id-2",
			AccessToken:  "test-access-token-2",
			RefreshToken: "test-refresh-token-2",
			UpdatedAt:    time.Now().UTC(),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId1,
				AccountName:    &accountName1,
				MeliCredential: meliCredentials1,
			},
			{
				ID:             accountId2,
				AccountName:    &accountName2,
				MeliCredential: meliCredentials2,
			},
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId1,
				AccountName: &accountName1,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-1",
					UserID:      "test-user-id-1",
				},
			},
			{
				ID:          accountId2,
				AccountName: &accountName2,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-2",
					UserID:      "test-user-id-2",
				},
			},
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

	t.Run("retrieve meli credentials from store id with expired credentials - one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId, accountName := entity.ID(uuid.New()), "test-account-name"

		meliCredentials := &common.MeliCredential{
			UserID:       "test-user-id",
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId,
				AccountName:    &accountName,
				MeliCredential: meliCredentials,
			},
		}

		refreshedMeliCredentials := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed",
			ExpiresIn:    21666,
			UserID:       meliCredentials.UserID,
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId,
				AccountName: &accountName,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-refreshed",
					UserID:      "test-user-id",
				},
			},
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromStoreID(storeId).Return(credentials, nil)
		mockMercadoLivre.EXPECT().RefreshCredentials((*credentials)[0].RefreshToken).Return(&refreshedMeliCredentials, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId, &refreshedMeliCredentials)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId.String()))

		cred, err := storeService.RetrieveMeliCredentialsFromStoreID(storeId)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from store id with expired credentials - more than one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId1, accountId2 := entity.ID(uuid.New()), entity.ID(uuid.New())
		accountName1, accountName2 := "test-account-name-1", "test-account-name-2"

		meliCredentials1 := &common.MeliCredential{
			UserID:       "test-user-id-1",
			AccessToken:  "test-access-token-1",
			RefreshToken: "test-refresh-token-1",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		meliCredentials2 := &common.MeliCredential{
			UserID:       "test-user-id-2",
			AccessToken:  "test-access-token-2",
			RefreshToken: "test-refresh-token-2",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId1,
				AccountName:    &accountName1,
				MeliCredential: meliCredentials1,
			},
			{
				ID:             accountId2,
				AccountName:    &accountName2,
				MeliCredential: meliCredentials2,
			},
		}

		refreshedMeliCredentials1 := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed-1",
			ExpiresIn:    21666,
			UserID:       meliCredentials1.UserID,
			RefreshToken: "test-refresh-token-1",
			UpdatedAt:    time.Now().UTC(),
		}

		refreshedMeliCredentials2 := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed-2",
			ExpiresIn:    21666,
			UserID:       meliCredentials2.UserID,
			RefreshToken: "test-refresh-token-2",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId1,
				AccountName: &accountName1,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-refreshed-1",
					UserID:      "test-user-id-1",
				},
			},
			{
				ID:          accountId2,
				AccountName: &accountName2,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-refreshed-2",
					UserID:      "test-user-id-2",
				},
			},
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromStoreID(storeId).Return(credentials, nil)
		mockMercadoLivre.EXPECT().RefreshCredentials((*credentials)[0].RefreshToken).Return(&refreshedMeliCredentials1, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId1, &refreshedMeliCredentials1)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId1.String()))

		mockMercadoLivre.EXPECT().RefreshCredentials((*credentials)[1].RefreshToken).Return(&refreshedMeliCredentials2, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId2, &refreshedMeliCredentials2)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId2.String()))

		cred, err := storeService.RetrieveMeliCredentialsFromStoreID(storeId)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from meli user id - one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId := entity.ID(uuid.New())
		id, accountName := "test-user-id", "test-account-name"

		meliCredentials := &common.MeliCredential{
			UserID:       id,
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId,
				OwnerID:        storeId,
				AccountName:    &accountName,
				MeliCredential: meliCredentials,
			},
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId,
				OwnerID:     storeId,
				AccountName: &accountName,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token",
					UserID:      "test-user-id",
				},
			},
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromMeliUserID(id).Return(credentials, nil)

		cred, err := storeService.RetrieveMeliCredentialsFromMeliUserID(id)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from meli user id - more than one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId1, accountId2 := entity.ID(uuid.New()), entity.ID(uuid.New())

		var (
			id  = "test-user-id"
			id2 = "test-user-id-2"
		)

		var (
			accountName1 = "test-account-name-1"
			accountName2 = "test-account-name-2"
		)

		meliCredentials1 := &common.MeliCredential{
			UserID:       id,
			AccessToken:  "test-access-token-1",
			RefreshToken: "test-refresh-token-1",
			UpdatedAt:    time.Now().UTC(),
		}

		meliCredentials2 := &common.MeliCredential{
			UserID:       id2,
			AccessToken:  "test-access-token-2",
			RefreshToken: "test-refresh-token-2",
			UpdatedAt:    time.Now().UTC(),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId1,
				OwnerID:        storeId,
				AccountName:    &accountName1,
				MeliCredential: meliCredentials1,
			},
			{
				ID:             accountId2,
				OwnerID:        storeId,
				AccountName:    &accountName2,
				MeliCredential: meliCredentials2,
			},
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId1,
				OwnerID:     storeId,
				AccountName: &accountName1,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-1",
					UserID:      id,
				},
			},
			{
				ID:          accountId2,
				OwnerID:     storeId,
				AccountName: &accountName2,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-2",
					UserID:      id2,
				},
			},
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromMeliUserID(id).Return(credentials, nil)

		cred, err := storeService.RetrieveMeliCredentialsFromMeliUserID(id)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from meli user id with expired credentials - one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId := entity.ID(uuid.New())
		id, accountName := "test-user-id", "test-account-name"

		meliCredentials := &common.MeliCredential{
			AccessToken:  "test-access-token",
			UserID:       id,
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId,
				OwnerID:        storeId,
				AccountName:    &accountName,
				MeliCredential: meliCredentials,
			},
		}

		refreshedMeliCredentials := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed",
			ExpiresIn:    21666,
			UserID:       id,
			UpdatedAt:    time.Now().UTC(),
			RefreshToken: "test-refresh-token",
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId,
				AccountName: &accountName,
				OwnerID:     storeId,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-refreshed",
					UserID:      id,
				},
			},
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromMeliUserID(id).Return(credentials, nil)
		mockMercadoLivre.EXPECT().RefreshCredentials((*credentials)[0].RefreshToken).Return(&refreshedMeliCredentials, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId, &refreshedMeliCredentials)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId.String()))

		cred, err := storeService.RetrieveMeliCredentialsFromMeliUserID(id)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("retrieve meli credentials from meli user id with expired credentials - more than one account", func(t *testing.T) {
		storeId := entity.ID(uuid.New())
		accountId1, accountId2 := entity.ID(uuid.New()), entity.ID(uuid.New())

		var (
			id  = "test-user-id"
			id2 = "test-user-id-2"
		)

		var (
			accountName1 = "test-account-name-1"
			accountName2 = "test-account-name-2"
		)

		meliCredentials1 := &common.MeliCredential{
			UserID:       id,
			AccessToken:  "test-access-token-1",
			RefreshToken: "test-refresh-token-1",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		meliCredentials2 := &common.MeliCredential{
			UserID:       id2,
			AccessToken:  "test-access-token-2",
			RefreshToken: "test-refresh-token-2",
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		}

		credentials := &[]store.Credentials{
			{
				ID:             accountId1,
				OwnerID:        storeId,
				AccountName:    &accountName1,
				MeliCredential: meliCredentials1,
			},
			{
				ID:             accountId2,
				OwnerID:        storeId,
				AccountName:    &accountName2,
				MeliCredential: meliCredentials2,
			},
		}

		refreshedMeliCredentials1 := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed-1",
			ExpiresIn:    21666,
			UserID:       id,
			UpdatedAt:    time.Now().UTC(),
			RefreshToken: "test-refresh-token-1",
		}

		refreshedMeliCredentials2 := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed-2",
			ExpiresIn:    21666,
			UserID:       id2,
			UpdatedAt:    time.Now().UTC(),
			RefreshToken: "test-refresh-token-2",
		}

		expectedCredentials := &[]store.Credentials{
			{
				ID:          accountId1,
				OwnerID:     storeId,
				AccountName: &accountName1,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-refreshed-1",
					UserID:      id,
				},
			},
			{
				ID:          accountId2,
				OwnerID:     storeId,
				AccountName: &accountName2,
				MeliCredential: &common.MeliCredential{
					AccessToken: "test-access-token-refreshed-2",
					UserID:      id2,
				},
			},
		}

		mockStoreRepo.EXPECT().RetrieveMeliCredentialsFromMeliUserID(id).Return(credentials, nil)
		mockMercadoLivre.EXPECT().RefreshCredentials((*credentials)[0].RefreshToken).Return(&refreshedMeliCredentials1, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId1, &refreshedMeliCredentials1)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId1.String()))

		mockMercadoLivre.EXPECT().RefreshCredentials((*credentials)[1].RefreshToken).Return(&refreshedMeliCredentials2, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId2, &refreshedMeliCredentials2)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId2.String()))

		cred, err := storeService.RetrieveMeliCredentialsFromMeliUserID(id)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})

	t.Run("refresh meli credentials", func(t *testing.T) {
		accountId := entity.ID(uuid.New())
		refreshToken := "test-refresh-token"

		refreshedCredentials := common.MeliCredential{
			AccessToken:  "test-access-token-refreshed",
			ExpiresIn:    21666,
			UserID:       "test-user-id",
			RefreshToken: "test-refresh-token",
			UpdatedAt:    time.Now().UTC(),
		}

		expectedCredentials := &store.Credentials{
			ID: accountId,
			MeliCredential: &common.MeliCredential{
				AccessToken: "test-access-token-refreshed",
				UserID:      "test-user-id",
			},
		}

		mockMercadoLivre.EXPECT().RefreshCredentials(refreshToken).Return(&refreshedCredentials, nil)
		mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId, &refreshedCredentials)
		mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId.String()))

		cred, err := storeService.RefreshMeliCredential(accountId, refreshToken)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !cmp.Equal(cred, expectedCredentials) {
			t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, expectedCredentials))
		}
	})
}
func TestValidateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStoreRepo := mock_store.NewMockRepository(ctrl)
	mockMercadoLivre := common_mock.NewMockMercadoLivre(ctrl)
	mockLogger := common_mock.NewMockLogger(ctrl)

	storeService := store.NewStoreService(mockStoreRepo, mockMercadoLivre, mockLogger)

	accountId := entity.ID(uuid.New())
	refreshToken := "test-refresh-token"
	accountName := "test-account-name"

	credentials := &store.Credentials{
		ID:          accountId,
		AccountName: &accountName,
		MeliCredential: &common.MeliCredential{
			AccessToken:  "test-access-token",
			UserID:       "test-user-id",
			RefreshToken: refreshToken,
			UpdatedAt:    time.Now().UTC().Add(-8 * time.Hour),
		},
	}

	refreshedCredentials := &common.MeliCredential{
		AccessToken:  "test-access-token-refreshed",
		UserID:       "test-user-id",
		RefreshToken: refreshToken,
		UpdatedAt:    time.Now().UTC(),
	}

	expectedCredentials := &store.Credentials{
		ID: accountId,
		MeliCredential: &common.MeliCredential{
			AccessToken: "test-access-token-refreshed",
			UserID:      "test-user-id",
		},
	}

	mockMercadoLivre.EXPECT().RefreshCredentials(refreshToken).Return(refreshedCredentials, nil)
	mockStoreRepo.EXPECT().UpdateMeliCredentials(accountId, refreshedCredentials)
	mockLogger.EXPECT().Info("Meli's credentials were updated", zap.String("account_id", accountId.String()))

	cred, err := store.ValidateCredentialsTest(storeService, credentials)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if !cmp.Equal(cred, expectedCredentials) {
		t.Errorf("Error refreshing and retrieving credentials. diff: %v", cmp.Diff(cred, refreshedCredentials))
	}
}
