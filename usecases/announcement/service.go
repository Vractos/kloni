package announcement

import (
	"encoding/json"
	"errors"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/common"
	"github.com/Vractos/kloni/usecases/store"
	"github.com/Vractos/kloni/utils"
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
		annIDsChunk := utils.Chunk(annIDs, 20)
		anns := make([]common.MeliAnnouncement, len(annIDs))
		for i, ids := range annIDsChunk {
			annsRes, err := a.meli.GetAnnouncements(ids, credentials.MeliAccessToken)
			if err != nil {
				cErr := &AnnouncementError{
					Message:       "Error to retrieve announcements",
					IsAbleToRetry: true,
				}
				a.logger.Error(cErr.Message, err, zap.Strings("announcements_ids", annIDs))
				return nil, cErr
			}
			start := i * 20
			end := start + 20
			if end >= len(anns) {
				end = len(anns)
			}
			copy(anns[start:end], *annsRes)
		}
		return &anns, err
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

func (a *AnnouncementService) UpdateQuantity(id string, newQuantity int, credentials store.Credentials, variationIDs ...int) error {
	err := a.meli.UpdateQuantity(newQuantity, id, credentials.MeliAccessToken, variationIDs...)
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

func (a *AnnouncementService) getAnnouncement(id string, credentials store.Credentials) (*common.MeliAnnouncement, error) {
	ann, err := a.meli.GetAnnouncement(id, credentials.MeliAccessToken)
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
	ann, err := a.getAnnouncement(input.RootID, *credentials)
	if err != nil {
		a.logger.Error(
			"Error in retrieving root announcement during the cloning process",
			err,
			zap.String("announcement_id", input.RootID))
		return errors.New("error to clone the announcement - get root announcement")
	}

	newAnns := make([]entity.Announcement, len(input.Titles)+1)
	classicAnn, err := entity.NewAnnouncement(ann)
	if err != nil {
		a.logger.Error("Error in the generation of a classic announcement during the cloning process", err, zap.String("announcement_id", ann.ID))
		return err
	}
	classicAnn.GenerateClassic()
	newAnns[0] = *classicAnn

	for i := 0; i < len(input.Titles); i++ {
		nAnn, err := entity.NewAnnouncement(ann)
		if err != nil {
			a.logger.Error("Error in the generation of a new announcement during the cloning process", err, zap.String("announcement_id", ann.ID))
			return err
		}
		nAnn.ChangeTitle(input.Titles[i])
		newAnns[i+1] = *nAnn
	}

	for _, ans := range newAnns {
		jsonAnn, err := json.Marshal(ans)
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

	return nil
}

// TODO: Implement this properly
func (a *AnnouncementService) ImportAnnouncement(input ImportAnnouncementDtoInput) error {
	originCredentials, err := a.store.RetrieveMeliCredentialsFromStoreID(input.StoreOrigin)
	if err != nil {
		a.logger.Error("Error in retrieving meli credentials, from origin store, during the importing process", err, zap.String("store_id", input.StoreOrigin.String()))
		return errors.New("error to import the announcement - get credentials")
	}

	destinyCredentials, err := a.store.RetrieveMeliCredentialsFromStoreID(input.StoreDestiny)
	if err != nil {
		a.logger.Error("Error in retrieving meli credentials, from destiny store, during the importing process", err, zap.String("store_id", input.StoreDestiny.String()))
		return errors.New("error to import the announcement - get credentials")
	}

	announcement, err := a.getAnnouncement(input.AnnouncementID, *originCredentials)
	if err != nil {
		a.logger.Error("Error in retrieving announcement during the importing process", err, zap.String("announcement_id", announcement.ID))
		return errors.New("error to import the announcement - get announcement")
	}

	if announcement.Quantity == 0 {
		return errors.New("error to import the announcement - quantity is 0")
	}

	newAnn, err := entity.NewAnnouncement(announcement)
	if err != nil {
		a.logger.Error("Error in the generation of a new announcement during the cloning process", err, zap.String("announcement_id", input.AnnouncementID))
		return err
	}

	jsonAnn, err := json.Marshal(newAnn)
	if err != nil {
		a.logger.Error("Error to marshal announcement json", err, zap.String("announcement_id", input.AnnouncementID))
		return errors.New("error to marshal announcement json")
	}

	rAnn, err := a.meli.PublishAnnouncement(jsonAnn, destinyCredentials.MeliAccessToken)
	if err != nil {
		cErr := &AnnouncementError{
			Message: "Error to publish an announcement",
		}
		a.logger.Error(cErr.Message, err, zap.String("announcement_id", input.AnnouncementID))
		return errors.New("error to publish clone")
	}

	a.logger.Info("Imported", zap.String("new_announcement_id", *rAnn))

	err = a.meli.AddDescription(announcement.Description, *rAnn, destinyCredentials.MeliAccessToken)
	if err != nil {
		a.logger.Error("Error to add description", err, zap.String("announcement_id", *rAnn))
	}

	return nil
}
