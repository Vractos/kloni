package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vractos/dolly/infrastructure/api/middleware"
	"github.com/Vractos/dolly/infrastructure/api/presenter"
	"github.com/Vractos/dolly/pkg/contexts"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/go-chi/chi/v5"
)

func registerStore(service store.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding store"
		input := &store.RegisterStoreDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		id, err := service.RegisterStore(*input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		output := &presenter.Store{
			ID: id,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	}
}

func getStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeId, err := contexts.RetrieveStoreIDFromCtx(r.Context())
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(storeId))
	}
}

func MakeStoreHandlers(r chi.Router, service store.UseCase) {
	r.Route("/store", func(r chi.Router) {
		r.Post("/", registerStore(service))
		r.With(middleware.EnsureValidToken()).With(middleware.AddStoreIDToCtx).Get("/", getStore())
	})
}
