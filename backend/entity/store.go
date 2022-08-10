package entity

import "github.com/google/uuid"

type MeliCredentials struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	Scope        string
	UserID       int
	RefreshToken string
}

type Store struct {
	ID              string
	Email           string
	Name            string
	Password        string
	MeliCredentials MeliCredentials
}

func NewUser() (*Store, error) {
	store := Store{
		ID: uuid.NewString(),
	}
	return &store, nil
}
