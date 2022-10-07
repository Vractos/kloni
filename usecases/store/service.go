package store

import (
	"github.com/Vractos/dolly/entity"
)

type StoreService struct {
	repo Repository
	meli MercadoLivre
}

func NewStoreService(repository Repository, mercadolivre MercadoLivre) *StoreService {
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
