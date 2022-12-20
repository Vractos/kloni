package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vractos/dolly/adapter/api/presenter"
	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/pkg/contexttools"
	"github.com/Vractos/dolly/usecases/announcement"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/go-chi/chi/v5"
)

func cloneAnnouncement(service announcement.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error to clone the announcement"
		input := &announcement.CloneAnnouncementDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		storeId, err := contexttools.RetrieveStoreIDFromCtx(r.Context())
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		store, err := entity.StringToID(storeId)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		input.Store = store

		err = service.CloneAnnouncement(*input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getAnnouncements(announce announcement.UseCase, store store.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error to get announcement"
		sku := chi.URLParam(r, "sku")

		storeId, err := contexttools.RetrieveStoreIDFromCtx(r.Context())
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}

		id, err := entity.StringToID(storeId)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		credential, err := store.RetrieveMeliCredentialsFromStoreID(id)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		anns, err := announce.RetrieveAnnouncements(sku, *credential)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		output := make([]*presenter.Announcement, len(*anns))
		for i, ann := range *anns {
			output[i] = &presenter.Announcement{
				ID:           ann.ID,
				Title:        ann.Title,
				Quantity:     ann.Quantity,
				Price:        ann.Price,
				ThumbnailURL: ann.ThumbnailURL,
				Sku:          ann.Sku,
				Link:         ann.Link,
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

func MakeAnnouncementHandlers(r chi.Router, announceService announcement.UseCase, storeService store.UseCase) {
	r.Route("/announcement", func(r chi.Router) {
		r.Post("/", cloneAnnouncement(announceService))
		r.Get("/{sku}", getAnnouncements(announceService, storeService))
	})
}
