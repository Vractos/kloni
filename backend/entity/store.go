package entity

import "github.com/google/uuid"

type Store struct {
	ID              string
	Email           string
	Name            string
	Password        string
	MeliCredentials struct {
		AccessToken  string
		TokenType    string
		ExpiresIn    int
		Scope        string
		UserID       int
		RefreshToken string
	}
}

func NewUser(email, password, name string) (*Store, error) {
	store := Store{
		ID:       uuid.NewString(),
		Email:    email,
		Password: password,
		Name:     name,
	}
	return &store, nil
}
