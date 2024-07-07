package store

import (
	"time"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/usecases/common"
	"go.uber.org/zap"
)

type StoreService struct {
	repo   Repository
	meli   common.MercadoLivre
	logger common.Logger
}

func NewStoreService(repository Repository, mercadolivre common.MercadoLivre, logger common.Logger) *StoreService {
	return &StoreService{
		repo:   repository,
		meli:   mercadolivre,
		logger: logger,
	}
}

func (s *StoreService) RegisterStore(input RegisterStoreDtoInput) (entity.ID, error) {
	store, err := entity.NewStore(input.Email, input.Name)
	if err != nil {
		s.logger.Error("Error in the generation of a new store", err)
		return store.ID, err
	}
	return s.repo.Create(store)
}

func (s *StoreService) RegisterMeliCredentials(input RegisterMeliCredentialsDtoInput) error {
	credentials, err := s.meli.RegisterCredential(input.Code)
	if err != nil {
		s.logger.Error(
			"Fail to register meli's credentials",
			err,
			zap.String("store_id", input.Store.String()),
		)
		return err
	}
	if err := s.repo.RegisterMeliCredential(input.Store, credentials); err != nil {
		s.logger.Error(
			"Fail to store meli's credentials",
			err,
			zap.String("store_id", input.Store.String()),
		)
		return err
	}
	return nil
}

// Checks if the credentials are still valid
func (s *StoreService) validateCredentials(credentials *Credentials) (*Credentials, error) {
	timeNowUTC := time.Now().UTC()
	if timeNowUTC.Sub(credentials.UpdatedAt.UTC()).Hours() >= 5 {
		credentialsData, err := s.RefreshMeliCredential(credentials.ID, credentials.RefreshToken)
		if err != nil {
			s.logger.Error(
				"Fail to refresh meli's credentials",
				err,
				zap.String("account_id", credentials.ID.String()),
			)
			return nil, err
		}

		return credentialsData, nil
	} else {
		return credentials, nil
	}
}

func (s *StoreService) RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]Credentials, error) {
	credentials, err := s.repo.RetrieveMeliCredentialsFromStoreID(id)
	if err != nil {
		s.logger.Error(
			"Fail to retrieve meli credentials via the store ID",
			err,
			zap.String("store_id", id.String()),
		)
		return nil, err
	}

	// Check if the credentials are still valid, if not, refresh them
	// Also, format the data to be returned
	for i, credential := range *credentials {
		credentialsData, err := s.validateCredentials(&credential)
		if err != nil {
			s.logger.Error(
				"Fail to validate meli's credentials",
				err,
				zap.String("account_id", credential.ID.String()),
			)
		}

		(*credentials)[i] = Credentials{
			ID:          credential.ID,
			AccountName: credential.AccountName,
			MeliCredential: &common.MeliCredential{
				AccessToken: credentialsData.AccessToken,
				UserID:      credentialsData.UserID,
			},
		}
	}

	return credentials, nil
}

func (s *StoreService) RetrieveMeliCredentialsFromMeliUserID(id string) (*[]Credentials, error) {
	credentials, err := s.repo.RetrieveMeliCredentialsFromMeliUserID(id)
	if err != nil {
		s.logger.Error(
			"Fail to retrieve meli credentials via the meli user ID",
			err,
			zap.String("user_id", id),
		)
		return nil, err
	}

	// Check if the credentials are still valid, if not, refresh them
	// Also, format the data to be returned
	for i, credential := range *credentials {
		credentialsData, err := s.validateCredentials(&credential)
		if err != nil {
			s.logger.Error(
				"Fail to validate meli's credentials",
				err,
				zap.String("account_id", credential.ID.String()),
			)
		}

		(*credentials)[i] = Credentials{
			ID:          credential.ID,
			OwnerID:     credential.OwnerID,
			AccountName: credential.AccountName,
			MeliCredential: &common.MeliCredential{
				AccessToken: credentialsData.AccessToken,
				UserID:      credentialsData.UserID,
			},
		}
	}

	return credentials, nil

}

func (s *StoreService) RefreshMeliCredential(accountId entity.ID, refreshToken string) (*Credentials, error) {
	credentials, err := s.meli.RefreshCredentials(refreshToken)
	if err != nil {
		s.logger.Error(
			"Fail to refresh meli's credentials",
			err,
			zap.String("account_id", accountId.String()),
		)
		return nil, err
	}

	s.logger.Info("Meli's credentials were updated", zap.String("account_id", accountId.String()))

	if err := s.repo.UpdateMeliCredentials(accountId, credentials); err != nil {
		s.logger.Error(
			"Fail to update meli's credentials",
			err,
			zap.String("account_id", accountId.String()),
		)
		return nil, err
	}

	return &Credentials{
		ID: accountId,
		MeliCredential: &common.MeliCredential{
			AccessToken: credentials.AccessToken,
			UserID:      credentials.UserID,
		},
	}, nil
}

// Exported for testing purposes
var ValidateCredentialsTest = (*StoreService).validateCredentials
