package store

import (
	"github.com/Vractos/dolly/entity"
)

type Reader interface {
	Get(id string) (*entity.Store, error)
}

type Writer interface {
	Create(e *entity.Store) (entity.ID, error)
	Update(e *entity.Store) error
	Delete(id entity.ID) error
	RegisterMeliCredentials(id entity.ID) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	RegisterStore(input RegisterStoreDtoInput) (entity.ID, error)
	RegisterMeliCredentials(input RegisterMeliCredentialsDtoInput) error
}
