package announcement

import (
	"github.com/Vractos/kloni/usecases/common"
	"github.com/Vractos/kloni/usecases/store"
)

type UseCase interface {
	RetrieveAnnouncements(sku string, credentials store.Credentials) (*[]common.MeliAnnouncement, error)
	UpdateQuantity(id string, quantity int, credentials store.Credentials, variationIDs ...int) error
	CloneAnnouncement(input CloneAnnouncementDtoInput) error
}
