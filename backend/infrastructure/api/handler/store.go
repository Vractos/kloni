package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vractos/dolly/backend/infrastructure/api/presenter"
	"github.com/Vractos/dolly/backend/usecases/store"
)

func registerStore(service store.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
			return
		}
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
			ID:    id,
			Email: input.Email,
			Name:  input.Name,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	}
}

func MakeStoreHandlers(service store.UseCase) http.HandlerFunc {

}
