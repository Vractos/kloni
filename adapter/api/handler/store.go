package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	mdw "github.com/Vractos/dolly/adapter/api/middleware"
	"github.com/Vractos/dolly/adapter/api/presenter"
	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/pkg/contexttools"
	"github.com/Vractos/dolly/pkg/metrics"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func registerStore(service store.UseCase, logger metrics.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding store"
		input := &store.RegisterStoreDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			logger.Error("Error to decode body", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		input.Email = strings.Replace(input.Email, "\n", "", -1)
		input.Email = strings.Replace(input.Email, "\r", "", -1)

		input.Name = strings.Replace(input.Name, "\n", "", -1)
		input.Name = strings.Replace(input.Name, "\r", "", -1)

		id, err := service.RegisterStore(*input)
		if err != nil {
			logger.Error(
				"Fail to register the store",
				err,
				zap.String("name", input.Name),
				zap.String("email", input.Email),
			)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		output := &presenter.Store{
			ID: id,
		}

		w.WriteHeader(http.StatusCreated)
		logger.Info(
			"A store was created",
			zap.String("email", input.Email),
			zap.String("name", input.Name),
			zap.String("id", id.String()),
		)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	}
}

func registerMeliCredentials(service store.UseCase, logger metrics.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		erroMessage := "Error to get credentials"
		input := &store.RegisterMeliCredentialsDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			logger.Error("Error to decode body", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(erroMessage))
			return
		}

		storeId, err := contexttools.RetrieveStoreIDFromCtx(r.Context())
		if err != nil {
			logger.Error("Fail to retrieve the storeID from the context", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(erroMessage))
			return
		}
		escapedStoreId := strings.Replace(storeId, "\n", "", -1)
		escapedStoreId = strings.Replace(escapedStoreId, "\r", "", -1)
		store, err := entity.StringToID(escapedStoreId)
		if err != nil {
			logger.Error("Fail to convert storeID from a string to an entity ID", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(erroMessage))
			return
		}

		input.Store = store

		err = service.RegisterMeliCredentials(*input)
		if err != nil {
			logger.Error("Fail to register meli's credentials", err, zap.String("store_id", storeId))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(erroMessage))
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func MakeStoreHandlers(r chi.Router, service store.UseCase, logger metrics.Logger) {
	r.Route("/store", func(r chi.Router) {
		r.Post("/", registerStore(service, logger))
		// That isn't the best approach to passing an auth middleware to only one route inside the route maker.
		// TODO Improve how this middleware is passed
		r.With(mdw.EnsureValidToken(logger)).With(mdw.AddStoreIDToCtx).Post("/meli-credentials", registerMeliCredentials(service, logger))
	})
}
