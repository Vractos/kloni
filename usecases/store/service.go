package store

import (
	"github.com/Vractos/dolly/entity"
)

type StoreService struct {
	repo Repository
}

func NewStoreService(repository Repository) *StoreService {
	return &StoreService{repo: repository}
}

func (s *StoreService) RegisterStore(input RegisterStoreDtoInput) (entity.ID, error) {
	store, err := entity.NewStore(input.Email, input.Name)
	if err != nil {
		return store.ID, err
	}
	return s.repo.Create(store)
}

func (s *StoreService) RegisterMeliCredentials(input RegisterMeliCredentialsDtoInput) error {
	panic("unimplemented")
}
