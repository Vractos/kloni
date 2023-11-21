package presenter

import (
	"github.com/Vractos/kloni/entity"
)

type Store struct {
	ID          entity.ID `json:"id"`
	ErroMessage string    `json:"error,omitempty"`
}
