package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/order"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func receiveMeliOrderNotification(service order.UseCase, logger metrics.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &order.OrderWebhookDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			logger.Error("Error to decode body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := service.ProcessWebhook(*input); err != nil {
			logger.Error(
				"Fail to process webhook",
				err,
				zap.String("notification_id", input.ID),
				zap.Int("user_id", input.UserID),
				zap.Int("attempts", input.Attempts),
				zap.String("sent", input.Sent),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func MakeOrderHandlers(r chi.Router, service order.UseCase, logger metrics.Logger) {
	r.Route("/order", func(r chi.Router) {
		r.Post("/meli-notification", receiveMeliOrderNotification(service, logger))
	})
}
