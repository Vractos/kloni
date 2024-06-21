package announcement

import (
	"github.com/Vractos/kloni/entity"
)

type CloneAnnouncementDtoInput struct {
	RootID string    `json:"root_id"`
	Titles []string  `json:"titles"`
	Store  entity.ID `json:"store_id"`
}

type GetAnnouncementsDtoInput struct {
	Sku string `json:"sku"`
}

type ImportAnnouncementDtoInput struct {
	AnnouncementID string    `json:"announcement_id"`
	StoreOrigin    entity.ID `json:"store_id_origin"`
	StoreDestiny   entity.ID `json:"store_id_destiny"`
}
