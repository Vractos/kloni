package announcement

import "github.com/Vractos/dolly/entity"

type UseCase interface {
	RetrieveAnnouncements(sku string) ([]entity.Announcement, error)
	RemoveQuantity(ids []string, quantity int) error
}
