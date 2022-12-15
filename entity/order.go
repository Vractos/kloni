package entity

import "time"

type OrderStatus string

const (
	Unknown           OrderStatus = ""
	Confirmed         OrderStatus = "confirmed"
	PaymentRequired   OrderStatus = "payment_required"
	PaymentInProcess  OrderStatus = "payment_in_process"
	PartiallyPaid     OrderStatus = "partially_paid"
	Paid              OrderStatus = "paid"
	PartiallyRefunded OrderStatus = "partially_refunded"
	PendingCancel     OrderStatus = "pending_cancel"
	Cancelled         OrderStatus = "cancelled"
	Invalid           OrderStatus = "invalid"
)

type OrderItem struct {
	Title    string
	Quantity int
	Sku      string
}

type Order struct {
	ID            ID
	MarketplaceID string
	DateCreated   time.Time
	Items         []OrderItem
	Status        OrderStatus
}

func NewOrder(id string, items []OrderItem, status OrderStatus) (*Order, error) {
	return &Order{
		ID:            NewID(),
		MarketplaceID: id,
		DateCreated:   time.Now().UTC(),
		Items:         items,
		Status:        status,
	}, nil
}
