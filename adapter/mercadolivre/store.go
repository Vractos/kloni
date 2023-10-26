package mercadolivre

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Vractos/kloni/usecases/common"
	"go.uber.org/zap"
)

// RegisterCredential implements common.MercadoLivre
func (m *MercadoLivre) RegisterCredential(code string) (*common.MeliCredential, error) {
	urlPath := fmt.Sprintf("%s/oauth/token", m.Endpoint)
	bodyRequest := map[string]interface{}{
		"client_id":     m.ClientId,
		"client_secret": m.ClientSecret,
		"code":          code,
		"redirect_uri":  m.RedirectUrl,
		"grant_type":    "authorization_code",
	}

	jsonBody, err := json.Marshal(bodyRequest)
	if err != nil {
		m.Logger.Error(
			"Fail to encode the request body",
			err,
		)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		m.Logger.Error(
			"Error to make a request to Mercado Livre",
			err,
			zap.String("path", "/"+urlPath),
		)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		registerCredentialsError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(registerCredentialsError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Couldn't retrieve the order",
			zap.String("meli_message", registerCredentialsError.Message),
			zap.String("meli_erro", registerCredentialsError.Error),
			zap.Any("cause", registerCredentialsError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)

		return nil, errors.New("error to register credentials")
	}

	credentials := &Credentials{}
	if err := json.NewDecoder(resp.Body).Decode(credentials); err != nil {
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
func (m *MercadoLivre) RefreshCredentials(refreshToken string) (*common.MeliCredential, error) {
	urlPath := fmt.Sprintf("%s/oauth/token", m.Endpoint)
	bodyRequest := map[string]interface{}{
		"grant_type":    "refresh_token",
		"client_id":     m.ClientId,
		"client_secret": m.ClientSecret,
		"refresh_token": refreshToken,
	}

	jsonBody, err := json.Marshal(bodyRequest)
	if err != nil {
		m.Logger.Error(
			"Fail to encode the request body",
			err,
		)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		m.Logger.Error(
			"Error to make a request to Mercado Livre",
			err,
			zap.String("path", "/"+urlPath),
		)
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		refreshCredentialsError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(refreshCredentialsError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Couldn't retrieve the order",
			zap.String("meli_message", refreshCredentialsError.Message),
			zap.String("meli_erro", refreshCredentialsError.Error),
			zap.Any("cause", refreshCredentialsError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return nil, errors.New("fail to refresh credentials")
	}

	credentials := &Credentials{}
	if err := json.NewDecoder(resp.Body).Decode(credentials); err != nil {
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
