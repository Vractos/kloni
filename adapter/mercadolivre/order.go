package mercadolivre

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Vractos/dolly/usecases/common"
	"github.com/go-playground/validator/v10"
)

// FetchOrder implements common.MercadoLivre
func (m *MercadoLivre) FetchOrder(orderId string, accessToken string) (*common.MeliOrder, error) {
	url := fmt.Sprintf("%s/oauth/token", m.Endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(resp.Body)
		log.Println("Error to fetch order: " + string(b))
		return nil, errors.New(string(b))
	}

	order := &Order{}
	if err := json.NewDecoder(resp.Body).Decode(order); err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	if resp.StatusCode == http.StatusPartialContent {
		err := m.Validate.Struct(order)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				log.Println(err.Field() + "was not provided")
			}
			return nil, errors.New("partial content")
		}
	}

	items := make([]common.OrderItem, len(order.OrderItems))
	for i, o := range order.OrderItems {
		items[i].ID = o.Item.ID
		items[i].Title = o.Item.Title
		items[i].Quantity = o.Quantity
		items[i].Sku = o.Item.SellerSku
	}

	return &common.MeliOrder{
		ID:          strconv.FormatUint(order.ID, 10),
		DateCreated: order.DateCreated,
		Status:      common.OrderStatus(order.Status),
		Items:       items,
	}, nil

}
