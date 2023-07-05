package store

import (
	"time"

	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/usecases/common"
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

func (s *StoreService) RetrieveMeliCredentialsFromStoreID(id entity.ID) (*Credentials, error) {
	credentials, err := s.repo.RetrieveMeliCredentialsFromStoreID(id)
	if err != nil {
		s.logger.Error(
			"Fail to retrieve meli credentials via the store ID",
			err,
			zap.String("store_id", id.String()),
		)
		return nil, err
	}

	timeNowUTC := time.Now().UTC()
	if timeNowUTC.Sub(credentials.UpdatedAt.UTC()).Hours() >= 5 {
		credentialsData, err := s.RefreshMeliCredential(id, credentials.RefreshToken)
		if err != nil {
			s.logger.Error(
				"Fail to refresh meli's credentials",
				err,
				zap.String("store_id", id.String()),
			)
			return nil, err
		}

		return credentialsData, nil
	}

	return &Credentials{
		StoreID:         id,
		MeliAccessToken: credentials.AccessToken,
		MeliUserID:      credentials.UserID,
	}, nil
}

func (s *StoreService) RetrieveMeliCredentialsFromMeliUserID(id string) (*Credentials, error) {
	storeId, credentials, err := s.repo.RetrieveMeliCredentialsFromMeliUserID(id)
	if err != nil {
		s.logger.Error(
			"Fail to retrieve meli credentials via the meli user ID",
			err,
			zap.String("user_id", id),
		)
		return nil, err
	}

	timeNowUTC := time.Now().UTC()
	if timeNowUTC.Sub(credentials.UpdatedAt.UTC()).Hours() >= 5 {
		credentialsData, err := s.RefreshMeliCredential(*storeId, credentials.RefreshToken)
		if err != nil {
			s.logger.Error(
				"Fail to refresh meli's credentials",
				err,
				zap.String("user_id", id),
			)
			return nil, err
		}

		return credentialsData, nil
	}

	return &Credentials{
		StoreID:         *storeId,
		MeliAccessToken: credentials.AccessToken,
		MeliUserID:      id,
	}, nil
}

func (s *StoreService) RefreshMeliCredential(storeId entity.ID, refreshToken string) (*Credentials, error) {
	credentials, err := s.meli.RefreshCredentials(refreshToken)
	if err != nil {
		s.logger.Error(
			"Fail to refresh meli's credentials",
			err,
			zap.String("store_id", storeId.String()),
		)
		return nil, err
	}

	s.logger.Info("Meli's credentials were updated", zap.String("store_id", storeId.String()))

	if err := s.repo.UpdateMeliCredentials(storeId, credentials); err != nil {
		s.logger.Error(
			"Fail to update meli's credentials",
			err,
			zap.String("store_id", storeId.String()),
		)
		return nil, err
	}

	return &Credentials{
		StoreID:         storeId,
		MeliAccessToken: credentials.AccessToken,
		MeliUserID:      credentials.UserID,
	}, nil
}
