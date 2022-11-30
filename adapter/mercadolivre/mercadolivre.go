package mercadolivre

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type MercadoLivre struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Endpoint     string
	Validate     *validator.Validate
	HttpClient   *http.Client
}

func NewMercadoLivre(clientId, clientSecret, redirectUrl, endpoint string, validator *validator.Validate) *MercadoLivre {
	return &MercadoLivre{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUrl:  redirectUrl,
		Endpoint:     endpoint,
		Validate:     validator,
		HttpClient:   &http.Client{Timeout: time.Second * 60},
	}
}
