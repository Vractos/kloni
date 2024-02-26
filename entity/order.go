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
	ID          ID
	Title       string
	Quantity    int
	Sku         string
	VariationID int
}

type Order struct {
	ID            ID
	StoreID       ID
	MarketplaceID string
	DateCreated   time.Time
	Items         []OrderItem
	Status        OrderStatus
}

func NewOrder(store ID, marketplace_id string, items []OrderItem, status OrderStatus) (*Order, error) {
	for i := range items {
		items[i].ID = NewID()
	}

	return &Order{
		ID:            NewID(),
		StoreID:       store,
		MarketplaceID: marketplace_id,
		DateCreated:   time.Now().UTC(),
		Items:         items,
		Status:        status,
	}, nil
}
