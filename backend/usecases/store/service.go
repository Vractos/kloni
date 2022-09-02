package store

import (
	"github.com/Vractos/dolly/backend/entity"
)

type StoreService struct {
	// repo Repository
}

// func NewStoreService(repository Repository) *StoreService {
// 	return &StoreService{repo: repository}
// }
func NewStoreService() *StoreService {
	return &StoreService{}
}

func (s *StoreService) RegisterStore(input RegisterStoreDtoInput) (entity.ID, error) {
	store, err := entity.NewStore(input.Email, input.Name, input.Password)
	if err != nil {
		return store.ID, err
	}
	// return s.repo.Create(store)
	return store.ID, nil
}
