package mercadolivre

import (
	"net/http"
	"time"

	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/go-playground/validator/v10"
)

type MercadoLivre struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Endpoint     string
	Validate     *validator.Validate
	HttpClient   *http.Client
	Logger       metrics.Logger
}

func NewMercadoLivre(clientId, clientSecret, redirectUrl, endpoint string, validator *validator.Validate, logger metrics.Logger) *MercadoLivre {
	return &MercadoLivre{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUrl:  redirectUrl,
		Endpoint:     endpoint,
		Validate:     validator,
		HttpClient:   &http.Client{Timeout: time.Second * 60},
		Logger:       logger,
	}
}
