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

func (a *AnnouncementService) UpdateQuantity(id string, newQuantity int, credentials store.Credentials) error {
	err := a.meli.UpdateQuantity(newQuantity, id, credentials.MeliAccessToken)
	if err != nil {
		return &AnnouncementError{
			Message:        "Error to update quantity",
			AnnouncementID: id,
			IsAbleToRetry:  true,
		}
	}
	return nil
}
