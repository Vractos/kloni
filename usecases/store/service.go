package store

import (
	"time"

	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/usecases/common"
)

type StoreService struct {
	repo Repository
	meli common.MercadoLivre
}

func NewStoreService(repository Repository, mercadolivre common.MercadoLivre) *StoreService {
	return &StoreService{repo: repository, meli: mercadolivre}
}

func (s *StoreService) RegisterStore(input RegisterStoreDtoInput) (entity.ID, error) {
	store, err := entity.NewStore(input.Email, input.Name)
	if err != nil {
		return store.ID, err
	}
	return s.repo.Create(store)
}

func (s *StoreService) RegisterMeliCredentials(input RegisterMeliCredentialsDtoInput) error {
	credentials, err := s.meli.RegisterCredential(input.Code)
	if err != nil {
		return err
	}
	if err := s.repo.RegisterMeliCredential(input.Store, credentials); err != nil {
		return err
	}
	return nil
}

func (s *StoreService) RetrieveMeliCredentialsFromStoreID(id entity.ID) (*Credentials, error) {
	credentials, err := s.repo.RetrieveMeliCredentialsFromStoreID(id)
	if err != nil {
		return nil, err
	}

	timeNowUTC := time.Now().UTC()
	if timeNowUTC.Sub(credentials.UpdatedAt.UTC()).Hours() >= 5 {
		credentialsData, err := s.RefreshMeliCredential(id, credentials.RefreshToken)
		if err != nil {
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
		return nil, err
	}

	timeNowUTC := time.Now().UTC()
	if timeNowUTC.Sub(credentials.UpdatedAt.UTC()).Hours() >= 5 {
		credentialsData, err := s.RefreshMeliCredential(*storeId, credentials.RefreshToken)
		if err != nil {
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
		return nil, err
	}

	if err := s.repo.UpdateMeliCredentials(storeId, credentials); err != nil {
		return nil, err
	}

	return &Credentials{
		StoreID:         storeId,
		MeliAccessToken: credentials.AccessToken,
		MeliUserID:      credentials.UserID,
	}, nil
}
