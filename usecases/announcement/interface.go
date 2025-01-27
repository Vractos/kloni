package announcement

import (
	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/usecases/common"
	"github.com/Vractos/kloni/usecases/store"
)

type Announcements struct {
	AccountID     entity.ID
	AccountName   string
	Announcements *[]common.MeliAnnouncement
}

type UseCase interface {
	// Retrieve announcements from a specific account
	RetrieveAnnouncements(sku string, credentials store.Credentials) (*[]common.MeliAnnouncement, error)
	// Retrieve announcements from all accounts that have the same SKU
	RetrieveAnnouncementsFromAllAccounts(sku string, credentials *[]store.Credentials) (*[]Announcements, error)
	UpdateQuantity(id string, quantity int, credentials store.Credentials, variationIDs ...int) error
	CloneAnnouncement(input CloneAnnouncementDtoInput, credentials *[]store.Credentials) error
	ImportAnnouncement(input ImportAnnouncementDtoInput, credentials *[]store.Credentials) error
}
