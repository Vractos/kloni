package announcement

import (
	"github.com/Vractos/kloni/entity"
)

type CloneAnnouncementDtoInput struct {
	RootID          string      `json:"root_id"`
	Titles          []string    `json:"titles"`
	RootAccountID   entity.ID   `json:"account_id"`
	DestinyAccounts []entity.ID `json:"destiny_accounts"`
}

type GetAnnouncementsDtoInput struct {
	Sku string `json:"sku"`
}

type ImportAnnouncementDtoInput struct {
	AnnouncementID string    `json:"announcement_id"`
	AccountOrigin  entity.ID `json:"account_id_origin"`
	AccountDestiny entity.ID `json:"account_id_destiny"`
}
