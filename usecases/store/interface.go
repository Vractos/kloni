package store

import (
	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/usecases/common"
)

type Credentials struct {
	ID          entity.ID
	OwnerID     entity.ID
	AccountName string
	*common.MeliCredential
}

// UseCase interface
type UseCase interface {
	RegisterStore(input RegisterStoreDtoInput) (entity.ID, error)
	RegisterMeliCredentials(input RegisterMeliCredentialsDtoInput) error
	RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]Credentials, error)
	RetrieveMeliCredentialsFromMeliUserID(id string) (*Credentials, error)
	RefreshMeliCredential(storeId entity.ID, refreshToken string) (*Credentials, error)
}

/*
#########################################
#########################################
---------------REPOSITORY---------------
#########################################
#########################################
*/

// Repository reader interface
type RepoReader interface {
	Get(id string) (*entity.Store, error)
	// Retrieves all meli credentials from a store
	RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]Credentials, error)
	RetrieveMeliCredentialsFromMeliUserID(accountId string) (*Credentials, error)
}

// Repository writer interface
type RepoWriter interface {
	Create(e *entity.Store) (entity.ID, error)
	RegisterMeliCredential(id entity.ID, c *common.MeliCredential) error
	UpdateMeliCredentials(accountId entity.ID, c *common.MeliCredential) error
	Update(e *entity.Store) error
	Delete(id entity.ID) error
}

// Repository interface
type Repository interface {
	RepoReader
	RepoWriter
}
