package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vractos/dolly/adapter/api/presenter"
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

func registerMeliCredentials(service store.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		erroMessage := "Error to get credentials"
		input := &store.RegisterMeliCredentialsDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(erroMessage))
		}

		err = service.RegisterMeliCredentials(*input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(erroMessage))
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func MakeStoreHandlers(r chi.Router, service store.UseCase) {
	r.Route("/store", func(r chi.Router) {
		r.Post("/", registerStore(service))
		r.Post("/meli-credentials", registerMeliCredentials(service))
	})
}
