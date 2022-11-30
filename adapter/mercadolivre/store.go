package mercadolivre

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Vractos/dolly/usecases/common"
)

// RegisterCredential implements common.MercadoLivre
func (m *MercadoLivre) RegisterCredential(code string) (*common.MeliCredential, error) {
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
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		log.Println("Error to register credentials: " + string(b))
		return nil, errors.New(string(b))
	}

	credentials := &Credentials{}
	if err := json.NewDecoder(resp.Body).Decode(credentials); err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	return &common.MeliCredential{
		AccessToken:  credentials.AccessToken,
		UserID:       strconv.Itoa(credentials.UserID),
		RefreshToken: credentials.RefreshToken,
		ExpiresIn:    credentials.ExpiresIn,
		UpdatedAt:    time.Now().UTC(),
	}, nil
}

// RefreshCredentials implements common.MercadoLivre
func (*MercadoLivre) RefreshCredentials(refreshToken string) (*common.MeliCredential, error) {
	panic("unimplemented")
}
