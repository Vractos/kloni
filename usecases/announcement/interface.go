package announcement

import (
	"github.com/Vractos/dolly/usecases/common"
	"github.com/Vractos/dolly/usecases/store"
)

type UseCase interface {
	RetrieveAnnouncements(sku string, credentials store.Credentials) (*[]common.MeliAnnouncement, error)
	UpdateQuantity(id string, quantity int, credentials store.Credentials) error
	CloneAnnouncement(input CloneAnnouncementDtoInput) error
}
