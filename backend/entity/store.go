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
	ID              uuid.UUID
	Email           string
	Name            string
	MeliCredentials MeliCredentials
}

func NewUser() (*Store, error) {
	return &Store{}
}
