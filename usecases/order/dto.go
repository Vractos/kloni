package order

import "time"

type OrderWebhookDtoInput struct {
	ID            string    `json:"_id"`
	Resource      string    `json:"resource"`
	UserID        int       `json:"user_id"`
	Topic         string    `json:"topic"`
	ApplicationID int64     `json:"application_id"`
	Attempts      int       `json:"attempts"`
	Sent          time.Time `json:"sent"`
	Received      time.Time `json:"received"`
}
