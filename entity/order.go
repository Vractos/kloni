package entity

import (
	"time"
)

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

func (s OrderStatus) String() string {
	switch s {
	case Unknown:
		return ""
	case Confirmed:
		return "confirmed"
	case PaymentRequired:
		return "payment_required"
	case PaymentInProcess:
		return "payment_in_process"
	case PartiallyPaid:
		return "partially_paid"
	case Paid:
		return "paid"
	case PartiallyRefunded:
		return "partially_refunded"
	case PendingCancel:
		return "pending_cancel"
	case Cancelled:
		return "cancelled"
	case Invalid:
		return "invalid"
	}
	return ""
}

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
