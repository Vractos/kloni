package store

import (
	"github.com/Vractos/dolly/backend/entity"
)

type Reader interface {
	Get(id string) (*entity.Store, error)
}

type Writer interface {
	Create(e *entity.Store) (entity.ID, error)
	Update(e *entity.Store) error
	Delete(id entity.ID) error
}

type StoreRepository interface {
	Reader
	Writer
}
