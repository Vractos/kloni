package store

import "github.com/Vractos/dolly/backend/entity"

type RegisterStoreDtoInput struct {
	Email    string
	Password string
	Name     string
}

type RegisterStoreDtoOutput struct {
	ID          entity.ID
	Email       string
	Name        string
	ErroMessage string
}
