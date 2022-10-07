package store

import (
	"time"

	"github.com/Vractos/dolly/entity"
)

// UseCase interface
type UseCase interface {
	RegisterStore(input RegisterStoreDtoInput) (entity.ID, error)
	RegisterMeliCredentials(input RegisterMeliCredentialsDtoInput) error
}

/*
#########################################
#########################################
---------------REPOSITORY---------------
#########################################
#########################################
*/

//Repository reader interface
type RepoReader interface {
	Get(id string) (*entity.Store, error)
}

//Repository writer interface
type RepoWriter interface {
	Create(e *entity.Store) (entity.ID, error)
	RegisterMeliCredential(id entity.ID, c *MeliCredential) error
	Update(e *entity.Store) error
	Delete(id entity.ID) error
}

// Repository interface
type Repository interface {
	RepoReader
	RepoWriter
}

/*
#########################################
#########################################
--------------Mercado Livre--------------
#########################################
#########################################
*/

/*
###################################
---------------Models-------------
###################################
*/
type MeliCredential struct {
	AccessToken  string
	ExpiresIn    int
	UserID       int
	RefreshToken string
	UpdateAt     time.Time
}

/*
###################################
-------------Interfaces------------
###################################
*/

// Mercado Livre reader interface
type MeliReader interface {
	// Implement
}

// Mercado Livre writer interface
type MeliWriter interface {
	RegisterCredential(code string) (*MeliCredential, error)
}

type MercadoLivre interface {
	MeliReader
	MeliWriter
}
