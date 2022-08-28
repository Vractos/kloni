package store

import "github.com/Vractos/dolly/backend/entity"

type StoreService struct {
	repo StoreRepository
}

func NewStoreService(repository StoreRepository) *StoreService {
	return &StoreService{repo: repository}
}

func (s *StoreService) RegisterStore(input RegisterStoreDtoInput) (RegisterStoreDtoOutput, error) {
	store, err := entity.NewStore(input.Email, input.Name, input.Password)
	if err != nil {
		return RegisterStoreDtoOutput{}, err
	}
	return RegisterStoreDtoOutput{ID: store.ID, Email: store.Email, Name: store.Name, ErroMessage: ""}, nil
}
