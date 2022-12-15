package announcement

import (
	"log"

	"github.com/Vractos/dolly/usecases/common"
	"github.com/Vractos/dolly/usecases/store"
)

type AnnouncementService struct {
	meli common.MercadoLivre
}

func NewAnnouncementService(mercadolivre common.MercadoLivre) *AnnouncementService {
	return &AnnouncementService{
		meli: mercadolivre,
	}
}

func (a *AnnouncementService) RetrieveAnnouncements(sku string, credentials store.Credentials) (*[]common.MeliAnnouncement, error) {
	annIDs, err := a.meli.GetAnnouncementsIDsViaSKU(sku, credentials.MeliUserID, credentials.MeliAccessToken)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if len(annIDs) >= 20 {
		// TODO: Chunk slice
	}
	anns, err := a.meli.GetAnnouncements(annIDs, credentials.MeliAccessToken)
	if err != nil {
		log.Println(err.Error())
		return nil, &AnnouncementError{
			Message:       "Error to get announcements",
			IsAbleToRetry: true,
		}
	}
	return anns, err
}

func (a *AnnouncementService) RemoveQuantity(id string, quantity int, credentials store.Credentials) error {
	log.Printf("Removing %v from %s...\n", quantity, id)
	return nil
	// panic("unimplemented")
}
