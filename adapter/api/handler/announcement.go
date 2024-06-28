package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Vractos/kloni/adapter/api/presenter"
	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/pkg/contexttools"
	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/announcement"
	"github.com/Vractos/kloni/usecases/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func cloneAnnouncement(service announcement.UseCase, store store.UseCase, logger metrics.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error to clone the announcement"
		input := &announcement.CloneAnnouncementDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			logger.Error("Error to decode body", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		storeId, err := contexttools.RetrieveStoreIDFromCtx(r.Context())
		if err != nil {
			logger.Error("Fail to retrieve the storeID from the context", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		strUUID, err := entity.StringToID(storeId)
		if err != nil {
			logger.Error("Fail to convert storeID from a string to an entity ID", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		credentials, err := store.RetrieveMeliCredentialsFromStoreID(strUUID)
		if err != nil {
			logger.Error("Couldn't retrieve meli's credentials", err, zap.String("store_id", storeId))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		err = service.CloneAnnouncement(*input, credentials)
		if err != nil {
			logger.Error("Error to clone announcement", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getAnnouncements(announce announcement.UseCase, store store.UseCase, logger metrics.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error to get announcement"
		notFound := "Announcement not found"
		sku := chi.URLParam(r, "sku")

		storeId, err := contexttools.RetrieveStoreIDFromCtx(r.Context())
		if err != nil {
			logger.Error("Fail to retrieve the storeID from the context", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}

		id, err := entity.StringToID(storeId)
		if err != nil {
			logger.Error("Fail to convert storeID from a string to an entity ID", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		credential, err := store.RetrieveMeliCredentialsFromStoreID(id)
		if err != nil {
			logger.Error("Couldn't retrieve meli's credentials", err, zap.String("store_id", storeId))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		anns, err := announce.RetrieveAnnouncementsFromAllAccounts(sku, credential)
		if err != nil {
			logger.Error("Fail to retrieve announcements", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		} else if anns == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(notFound))
			return
		}

		output := []presenter.Announcement{}
		for _, ann := range *anns {
			if ann.Announcements == nil {
				continue
			}
			account := struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   ann.AccountID.String(),
				Name: ann.AccountName,
			}

			for _, an := range *ann.Announcements {
				output = append(output, presenter.Announcement{
					ID:           an.ID,
					Title:        an.Title,
					Price:        an.Price,
					ThumbnailURL: an.ThumbnailURL,
					Quantity:     an.Quantity,
					Sku:          an.Sku,
					Link:         an.Link,
					Account:      account,
				})
			}
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	}
}

func MakeAnnouncementHandlers(r chi.Router, announceService announcement.UseCase, storeService store.UseCase, logger metrics.Logger) {
	r.Route("/announcement", func(r chi.Router) {
		r.Post("/", cloneAnnouncement(announceService, storeService, logger))
		r.Get("/{sku}", getAnnouncements(announceService, storeService, logger))
	})
}
