package store

import (
	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/usecases/common"
)

type Credentials struct {
	ID          entity.ID
	OwnerID     entity.ID
	AccountName *string
	*common.MeliCredential
}

// UseCase interface
type UseCase interface {
	RegisterStore(input RegisterStoreDtoInput) (entity.ID, error)
	RegisterMeliCredentials(input RegisterMeliCredentialsDtoInput) error
	// Retrieve all the meli credentials from a store
	RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]Credentials, error)
	// Retrieve all meli credentials from a meli user id
	RetrieveMeliCredentialsFromMeliUserID(id string) (*[]Credentials, error)
	RefreshMeliCredential(accountId entity.ID, refreshToken string) (*Credentials, error)
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
	// Retrieves meli credentials from a meli user id
	//
	// Example of return:
	//
	// 	Credentials {
	// 		ID: "000000-000000-000000-000000-000000",
	// 		OwnerID: "000000-000000-000000-000000-000000",
	// 		AccountName: "Name",
	// 		AccessToken: "111111-111111-111111-111111-111111",
	// 		RefreshToken: "222222-222222-222222-222222-222222",
	// 		UpdateAt: "2021-01-01T00:00:00Z",
	// 	}
	//
	RetrieveMeliCredentialsFromMeliUserID(accountId string) (*[]Credentials, error)
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
