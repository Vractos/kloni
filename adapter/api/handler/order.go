package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vractos/dolly/usecases/order"
	"github.com/go-chi/chi/v5"
)

func receiveMeliOrderNotification(service order.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &order.OrderWebhookDtoInput{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := service.ProcessWebhook(*input); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func MakeOrderHandlers(r chi.Router, service order.UseCase) {
	r.Route("/order", func(r chi.Router) {
		r.Post("/meli-notification", receiveMeliOrderNotification(service))
	})
}
