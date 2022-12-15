package mercadolivre

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Vractos/dolly/usecases/common"
)

func (m *MercadoLivre) GetAnnouncementsIDsViaSKU(sku string, userId string, accessToken string) ([]string, error) {
	url := fmt.Sprintf("%s/users/%s/items/search?seller_sku=%s", m.Endpoint, userId, sku)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New("Error to fetch clones" + string(b))
	}

	queryAnnouncementResult := &QueryAnnouncementViaSku{}
	if err := json.NewDecoder(resp.Body).Decode(queryAnnouncementResult); err != nil {
		return nil, err
	}

	if queryAnnouncementResult.Results == nil {
		return nil, errors.New("announcements not found")
	}

	return queryAnnouncementResult.Results, nil
}

func (m *MercadoLivre) GetAnnouncements(ids []string, accessToken string) (*[]common.MeliAnnouncement, error) {
	url := fmt.Sprintf("%s/items?ids=%s", m.Endpoint, strings.Join(ids, ","))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New("Error to fetch announcements" + string(b))
	}

	queryAnnouncementsResult := &AnnouncementsMultiGet{}
	if err := json.NewDecoder(resp.Body).Decode(queryAnnouncementsResult); err != nil {
		return nil, err
	}

	meliAnnouncement := make([]common.MeliAnnouncement, len(ids))
	for i, a := range *queryAnnouncementsResult {
		if a.Code != http.StatusOK {
			return nil, errors.New("error to fetch all clones. " + a.Body.ID + a.Body.Error)
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
