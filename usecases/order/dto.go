package order

type OrderWebhookDtoInput struct {
	ID            string `json:"_id"`
	Resource      string `json:"resource"`
	UserID        int    `json:"user_id"`
	Topic         string `json:"topic"`
	ApplicationID int64  `json:"application_id"`
	Attempts      int    `json:"attempts"`
	// Can be converted to time
	Sent string `json:"sent"`
	// Can be converted to time
	Received string `json:"received"`
}
