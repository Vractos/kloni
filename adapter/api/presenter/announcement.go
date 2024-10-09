package presenter

type Announcement struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	ThumbnailURL string  `json:"picture"`
	Sku          string  `json:"sku"`
	Link         string  `json:"link"`
	Account      struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
}
