package presenter

import (
	"github.com/Vractos/dolly/entity"
)

type Store struct {
	ID          entity.ID `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	ErroMessage string    `json:"error,omitempty"`
}
