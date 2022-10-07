package store

import "github.com/Vractos/dolly/entity"

type RegisterStoreDtoInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RegisterMeliCredentialsDtoInput struct {
	Code  string    `json:"code"`
	Store entity.ID `json:"store_id"`
}
