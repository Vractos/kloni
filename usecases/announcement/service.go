package announcement

import (
	"encoding/json"
	"errors"

	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/pkg/metrics"
	"github.com/Vractos/dolly/usecases/common"
	"github.com/Vractos/dolly/usecases/store"
	"go.uber.org/zap"
)

type AnnouncementService struct {
	meli   common.MercadoLivre
	store  store.UseCase
	logger metrics.Logger
}

func NewAnnouncementService(mercadolivre common.MercadoLivre, storeUseCase store.UseCase, logger metrics.Logger) *AnnouncementService {
	return &AnnouncementService{
		meli:   mercadolivre,
		store:  storeUseCase,
		logger: logger,
	}
}

func (a *AnnouncementService) RetrieveAnnouncements(sku string, credentials store.Credentials) (*[]common.MeliAnnouncement, error) {
	annIDs, err := a.meli.GetAnnouncementsIDsViaSKU(sku, credentials.MeliUserID, credentials.MeliAccessToken)
	if err != nil {
		cErr := &AnnouncementError{
			Message: "Error to retrieve the announcements IDs",
			Sku:     sku,
		}
		a.logger.Error(cErr.Message, err,
			zap.String("sku", cErr.Sku),
		)
		return nil, cErr
	} else if len(annIDs) == 0 {
		return nil, nil
	}

	if len(annIDs) >= 20 {
		// TODO Chunk slice
	}

	anns, err := a.meli.GetAnnouncements(annIDs, credentials.MeliAccessToken)
	if err != nil {
		cErr := &AnnouncementError{
			Message:       "Error to retrieve announcements",
			IsAbleToRetry: true,
		}
		a.logger.Error(cErr.Message, err, zap.Strings("announcements_ids", annIDs))
		return nil, cErr
	}
	return anns, err
}

func (a *AnnouncementService) UpdateQuantity(id string, newQuantity int, credentials store.Credentials) error {
	err := a.meli.UpdateQuantity(newQuantity, id, credentials.MeliAccessToken)
	if err != nil {
		cErr := &AnnouncementError{
			Message:        "Error to update quantity",
			AnnouncementID: id,
			IsAbleToRetry:  true,
		}
		a.logger.Error(cErr.Message, err, zap.String("announcement_id", id))
		return cErr
	}
	a.logger.Info("Quantity updated", zap.String("announcement_id", id), zap.Int("new_quantity", newQuantity))
	return nil
}

func (a *AnnouncementService) getAnnouncement(id string) (*common.MeliAnnouncement, error) {
	ann, err := a.meli.GetAnnouncement(id)
	if err != nil {
		cErr := &AnnouncementError{
			Message:        "Error to retrieve root announcement",
			AnnouncementID: id,
		}
		a.logger.Error(cErr.Message, err, zap.String("announcement_id", id))
		return nil, cErr
	}

	descrp, err := a.meli.GetDescription(id)
	if err != nil {
		a.logger.Error("Error to retrieve description", err, zap.String("announcement_id", id))
		return nil, errors.New("erro to get description")
	}
	ann.Description = *descrp
	return ann, nil
}

// CloneAnnouncement implements UseCase
func (a *AnnouncementService) CloneAnnouncement(input CloneAnnouncementDtoInput) error {
	credentials, err := a.store.RetrieveMeliCredentialsFromStoreID(input.Store)
	if err != nil {
		a.logger.Error("Error in retrieving meli credentials during the cloning process", err, zap.String("store_id", input.Store.String()))
		return errors.New("error to clone the announcement - get credentials")
	}
	ann, err := a.getAnnouncement(input.RootID)
	if err != nil {
		a.logger.Error(
			"Error in retrieving root announcement during the cloning process",
			err,
			zap.String("announcement_id", input.RootID))
		return errors.New("error to clone the announcement - get root announcement")
	}

	newAnns := make([][2]entity.Announcement, len(input.Titles))
	for i, t := range input.Titles {
		nAnn, err := entity.NewAnnouncement(ann)
		if err != nil {
			a.logger.Error("Error in the generation of a new announcement during the cloning process", err, zap.String("announcement_id", ann.ID))
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
				a.logger.Error("Error to marshal announcement json", err, zap.String("announcement_id", input.RootID))
				return errors.New("error to marshal announcement json")
			}

			rAnn, err := a.meli.PublishAnnouncement(jsonAnn, credentials.MeliAccessToken)
			if err != nil {
				cErr := &AnnouncementError{
					Message: "Error to publish an announcement",
				}
				a.logger.Error(cErr.Message, err, zap.String("announcement_id", input.RootID))
				return errors.New("error to publish clone")
			}

			a.logger.Info("New clone", zap.String("new_announcement_id", *rAnn))

			err = a.meli.AddDescription(ann.Description, *rAnn, credentials.MeliAccessToken)
			if err != nil {
				a.logger.Error("Error to add description", err, zap.String("announcement_id", *rAnn))
			}
		}
	}

	return nil
}
