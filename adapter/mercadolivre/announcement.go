package mercadolivre

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Vractos/kloni/usecases/common"
	"go.uber.org/zap"
)

func (m *MercadoLivre) GetAnnouncementsIDsViaSKU(sku string, userId string, accessToken string) ([]string, error) {
	urlPath := fmt.Sprintf("%s/users/%s/items/search?seller_sku=%s", m.Endpoint, userId, sku)

	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		m.Logger.Error(
			"Error to make a request to Mercado Livre",
			err,
			zap.String("sku", sku),
			zap.String("user_id", userId),
			zap.String("path", "/"+urlPath),
		)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		queryAnnouncementsError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(queryAnnouncementsError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Couldn't retrieve any announcements IDs",
			zap.String("sku", sku),
			zap.String("meli_message", queryAnnouncementsError.Message),
			zap.String("meli_erro", queryAnnouncementsError.Error),
			zap.Any("cause", queryAnnouncementsError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return nil, errors.New("error to fetch clones")
	}

	queryAnnouncementResult := &QueryAnnouncementViaSku{}
	if err := json.NewDecoder(resp.Body).Decode(queryAnnouncementResult); err != nil {
		return nil, err
	}

	return queryAnnouncementResult.Results, nil
}

func (m *MercadoLivre) GetAnnouncements(ids []string, accessToken string) (*[]common.MeliAnnouncement, error) {
	urlPath := fmt.Sprintf("%s/items?ids=%s", m.Endpoint, strings.Join(ids, ","))

	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		m.Logger.Error(
			"Error to make a request to Mercado Livre",
			err,
			zap.Strings("announcements_ids", ids),
			zap.String("path", "/"+urlPath),
		)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		queryAnnouncementsError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(queryAnnouncementsError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Couldn't retrieve any announcements",
			zap.Strings("announcements_ids", ids),
			zap.String("meli_message", queryAnnouncementsError.Message),
			zap.String("meli_erro", queryAnnouncementsError.Error),
			zap.Any("cause", queryAnnouncementsError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return nil, errors.New("error to fetch announcements")
	}

	queryAnnouncementsResult := &AnnouncementsMultiGet{}
	if err := json.NewDecoder(resp.Body).Decode(queryAnnouncementsResult); err != nil {
		return nil, err
	}

	meliAnnouncement := make([]common.MeliAnnouncement, len(ids))
	for i, a := range *queryAnnouncementsResult {
		if a.Code != http.StatusOK {
			m.Logger.Warn(
				"Fail to fetch all clones",
				zap.String("announcement_id", a.Body.ID),
				zap.String("error", a.Body.Error),
			)
			return nil, errors.New("error to fetch all clones")
		}
		var sku string
		for _, v := range a.Body.Attributes {
			if v.ID == "SELLER_SKU" {
				sku = v.ValueName
				break
			}
		}
		meliAnnouncement[i] = common.MeliAnnouncement{
			ID:           a.Body.ID,
			Title:        a.Body.Title,
			Quantity:     a.Body.AvailableQuantity,
			Price:        a.Body.Price,
			ThumbnailURL: a.Body.Thumbnail,
			Sku:          sku,
			Link:         a.Body.Permalink,
		}

	}

	return &meliAnnouncement, nil
}

func (m *MercadoLivre) UpdateQuantity(quantity int, announcementId, accessToken string) error {
	urlPath := fmt.Sprintf("%s/items/%s", m.Endpoint, announcementId)
	bodyRequest := map[string]interface{}{
		"available_quantity": quantity,
	}

	jsonBody, err := json.Marshal(bodyRequest)
	if err != nil {
		m.Logger.Error(
			"Fail to encode the request body",
			err,
		)
		return err
	}

	req, err := http.NewRequest(http.MethodPut, urlPath, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		m.Logger.Error(
			"Error to make a request to Mercado Livre",
			err,
			zap.String("announcement_id", announcementId),
			zap.String("path", "/"+urlPath),
		)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		updateAnnouncementsError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(updateAnnouncementsError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return err
		}
		m.Logger.Warn(
			"Couldn't update the quantity",
			zap.String("announcement_id", announcementId),
			zap.String("meli_message", updateAnnouncementsError.Message),
			zap.String("meli_erro", updateAnnouncementsError.Error),
			zap.Any("cause", updateAnnouncementsError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return errors.New("fail to update the quantity")
	}

	return nil
}

// GetAnnouncement implements common.MercadoLivre
func (m *MercadoLivre) GetAnnouncement(id string) (*common.MeliAnnouncement, error) {
	urlPath := fmt.Sprintf("%s/items/%s", m.Endpoint, id)

	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
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
			zap.String("announcement_id", id),
			zap.String("path", "/"+urlPath),
		)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		queryAnnouncementError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(queryAnnouncementError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Couldn't retrieve announcement",
			zap.String("announcements_id", id),
			zap.String("meli_message", queryAnnouncementError.Message),
			zap.String("meli_erro", queryAnnouncementError.Error),
			zap.Any("cause", queryAnnouncementError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return nil, errors.New("error to fetch announcement")
	}

	aR := &Announcement{}
	if err := json.NewDecoder(resp.Body).Decode(aR); err != nil {
		return nil, err
	}

	pics := make([]string, len(aR.Pictures))
	for i, p := range aR.Pictures {
		pics[i] = p.URL
	}

	return &common.MeliAnnouncement{
		ID:            aR.ID,
		Title:         aR.Title,
		Quantity:      aR.AvailableQuantity,
		Price:         aR.Price,
		CategoryID:    aR.CategoryID,
		Condition:     aR.Condition,
		ListingTypeID: aR.ListingTypeID,
		Channels:      aR.Channels,
		Pictures:      pics,
		SaleTerms: []struct {
			ID          string
			Name        string
			ValueID     interface{}
			ValueName   string
			ValueStruct struct {
				Number int
				Unit   string
			}
			Values []struct {
				ID     interface{}
				Name   string
				Struct struct {
					Number int
					Unit   string
				}
			}
			ValueType string
		}(aR.SaleTerms),
		Attributes: []struct {
			ID          string
			Name        string
			ValueID     string
			ValueName   string
			ValueStruct interface{}
			Values      []struct {
				ID     string
				Name   string
				Struct interface{}
			}
			AttributeGroupID   string
			AttributeGroupName string
			ValueType          string
		}(aR.Attributes),
	}, nil
}

func (m *MercadoLivre) GetDescription(id string) (*string, error) {
	urlPath := fmt.Sprintf("%s/items/%s/description", m.Endpoint, id)

	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
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
			zap.String("announcement_id", id),
			zap.String("path", "/"+urlPath),
		)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		getDescriptionError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(getDescriptionError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Couldn't retrieve description",
			zap.String("announcements_id", id),
			zap.String("meli_message", getDescriptionError.Message),
			zap.String("meli_erro", getDescriptionError.Error),
			zap.Any("cause", getDescriptionError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return nil, errors.New("error to fetch description")
	}

	description := &Description{}
	if err := json.NewDecoder(resp.Body).Decode(description); err != nil {
		return nil, err
	}

	return &description.PlainText, nil
}

// AddDescription implements common.MercadoLivre
func (m *MercadoLivre) AddDescription(description string, announcementId string, accessToken string) error {
	urlPath := fmt.Sprintf("%s/items/%s/description", m.Endpoint, announcementId)
	bodyRequest := map[string]interface{}{
		"plain_text": description,
	}

	jsonBody, err := json.Marshal(bodyRequest)
	if err != nil {
		m.Logger.Error(
			"Fail to encode the request body",
			err,
		)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		m.Logger.Error(
			"Error to make a request to Mercado Livre",
			err,
			zap.String("announcement_id", announcementId),
			zap.String("path", "/"+urlPath),
		)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		addAnnouncementsError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(addAnnouncementsError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return err
		}
		m.Logger.Warn(
			"Fail to add a description to the announcement",
			zap.String("announcement_id", announcementId),
			zap.String("meli_message", addAnnouncementsError.Message),
			zap.String("meli_erro", addAnnouncementsError.Error),
			zap.Any("cause", addAnnouncementsError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return errors.New("error to add description")
	}

	return nil
}

// PublishAnnouncement implements common.MercadoLivre
func (m *MercadoLivre) PublishAnnouncement(announcement []byte, accessToken string) (ID *string, err error) {
	urlPath := fmt.Sprintf("%s/items", m.Endpoint)

	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer(announcement))
	if err != nil {
		m.Logger.Error(
			"Fail to encode the request body",
			err,
		)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
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
		publishAnnouncementsError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(publishAnnouncementsError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Fail to add a description to the announcement",
			zap.String("meli_message", publishAnnouncementsError.Message),
			zap.String("meli_erro", publishAnnouncementsError.Error),
			zap.Any("cause", publishAnnouncementsError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return nil, errors.New("error to publish announcement")
	}

	result := &Announcement{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	return &result.ID, nil

}
