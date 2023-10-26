package mercadolivre

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Vractos/kloni/usecases/common"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// FetchOrder implements common.MercadoLivre
func (m *MercadoLivre) FetchOrder(orderId string, accessToken string) (*common.MeliOrder, error) {
	urlPath := fmt.Sprintf("%s/orders/%s", m.Endpoint, orderId)

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
			zap.String("order_id", orderId),
			zap.String("path", "/"+urlPath),
		)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		queryOrderError := &MeliError{}
		if err := json.NewDecoder(resp.Body).Decode(queryOrderError); err != nil {
			m.Logger.Error(
				"Error to decode response body",
				err,
			)
			return nil, err
		}
		m.Logger.Warn(
			"Couldn't retrieve the order",
			zap.String("order_id", orderId),
			zap.String("meli_message", queryOrderError.Message),
			zap.String("meli_erro", queryOrderError.Error),
			zap.Any("cause", queryOrderError.Cause),
			zap.Int("status_code", resp.StatusCode),
		)
		return nil, errors.New("error to fetch order")
	}

	order := &Order{}
	if err := json.NewDecoder(resp.Body).Decode(order); err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusPartialContent {
		err := m.Validate.Struct(order)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				m.Logger.Warn(
					err.Field()+"was not provided",
					zap.String("order_id", orderId),
				)
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
