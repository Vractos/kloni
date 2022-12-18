package announcement

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/usecases/common"
	"github.com/Vractos/dolly/usecases/store"
)

type AnnouncementService struct {
	meli  common.MercadoLivre
	store store.UseCase
}

func NewAnnouncementService(mercadolivre common.MercadoLivre, storeUseCase store.UseCase) *AnnouncementService {
	return &AnnouncementService{
		meli:  mercadolivre,
		store: storeUseCase,
	}
}

func (a *AnnouncementService) RetrieveAnnouncements(sku string, credentials store.Credentials) (*[]common.MeliAnnouncement, error) {
	annIDs, err := a.meli.GetAnnouncementsIDsViaSKU(sku, credentials.MeliUserID, credentials.MeliAccessToken)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	} else if len(annIDs) == 0 {
		return nil, &AnnouncementError{
			Message: "Not found",
			Sku:     sku,
		}
	}

	if len(annIDs) >= 20 {
		// TODO Chunk slice
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
	log.Println(id)
	return nil
}

func (a *AnnouncementService) getAnnouncement(id string) (*common.MeliAnnouncement, error) {
	ann, err := a.meli.GetAnnouncement(id)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("erro to get the root announcement")
	}
	descrp, err := a.meli.GetDescription(id)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("erro to get description")
	}
	ann.Description = *descrp
	return ann, nil
}

// CloneAnnouncement implements UseCase
func (a *AnnouncementService) CloneAnnouncement(input CloneAnnouncementDtoInput) error {
	credentials, err := a.store.RetrieveMeliCredentialsFromStoreID(input.Store)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error to clone the announcement - get credentials")
	}
	ann, err := a.getAnnouncement(input.RootID)
	if err != nil {
		return err
	}

	newAnns := make([][2]entity.Announcement, len(input.Titles))
	for i, t := range input.Titles {
		nAnn, err := entity.NewAnnouncement(ann)
		if err != nil {
			return err
		}
		nAnn.ChangeTitle(t)
		newAnns[i][0] = *nAnn
		nAnn.GenerateClassic()
		newAnns[i][1] = *nAnn
	}

	for _, ans := range newAnns {
		for _, an := range ans {
			jsonAnn, err := json.Marshal(an)
			if err != nil {
				log.Println(err.Error())
				return errors.New("error to marshal announcement json")
			}

			rAnn, err := a.meli.PublishAnnouncement(jsonAnn, credentials.MeliAccessToken)
			if err != nil {
				log.Println(err.Error())
				return errors.New("error to publish clone")
			}

			log.Println(*rAnn)

			err = a.meli.AddDescription(ann.Description, *rAnn, credentials.MeliAccessToken)
			if err != nil {
				log.Println("Error to add description: " + err.Error())
			}
		}
	}

	return nil
}
