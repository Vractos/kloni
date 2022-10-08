package mercadolivre

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Vractos/dolly/usecases/store"
)

type MercadoLivreStore struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Endpoint     string
	HttpClient   *http.Client
}

func NewMercadoLivreStore(clientId, clientSecret, redirectUrl, endpoint string) *MercadoLivreStore {
	return &MercadoLivreStore{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUrl:  redirectUrl,
		Endpoint:     endpoint,
		HttpClient:   &http.Client{},
	}
}

// RegisterCredential implements store.MercadoLivre
func (m *MercadoLivreStore) RegisterCredential(code string) (*store.MeliCredential, error) {
	url := fmt.Sprintf("%s/oauth/token", m.Endpoint)
	bodyRequest := map[string]interface{}{
		"client_id":     m.ClientId,
		"client_secret": m.ClientSecret,
		"code":          code,
		"redirect_uri":  m.RedirectUrl,
		"grant_type":    "authorization_code",
	}

	jsonBody, err := json.Marshal(bodyRequest)
	if err != nil {
		log.Printf("Failed to parse the issuer url: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	} else if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(b))
		return nil, errors.New(string(b))
	}

	defer resp.Body.Close()

	credentials := &Credentials{}
	if err := json.NewDecoder(resp.Body).Decode(credentials); err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	return &store.MeliCredential{
		AccessToken:  credentials.AccessToken,
		UserID:       credentials.UserID,
		RefreshToken: credentials.RefreshToken,
		ExpiresIn:    credentials.ExpiresIn,
		UpdatedAt:    time.Now().UTC(),
	}, nil
}

// RefreshCredentials implements store.MercadoLivre
func (*MercadoLivreStore) RefreshCredentials(refreshToken string) (*store.MeliCredential, error) {
	panic("unimplemented")
}
