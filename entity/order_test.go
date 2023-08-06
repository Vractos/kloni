package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestOrderStatus(t *testing.T) {
	tests := []struct {
		name string
		s    OrderStatus
		want string
	}{
		{
			name: "Unknown",
			s:    Unknown,
			want: "",
		},
		{
			name: "Confirmed",
			s:    Confirmed,
			want: "confirmed",
		},
		{
			name: "PaymentRequired",
			s:    PaymentRequired,
			want: "payment_required",
		},
		{
			name: "PaymentInProcess",
			s:    PaymentInProcess,
			want: "payment_in_process",
		},
		{
			name: "PartiallyPaid",
			s:    PartiallyPaid,
			want: "partially_paid",
		},
		{
			name: "Paid",
			s:    Paid,
			want: "paid",
		},
		{
			name: "PartiallyRefunded",
			s:    PartiallyRefunded,
			want: "partially_refunded",
		},
		{
			name: "PendingCancel",
			s:    PendingCancel,
			want: "pending_cancel",
		},
		{
			name: "Cancelled",
			s:    Cancelled,
			want: "cancelled",
		},
		{
			name: "Invalid",
			s:    Invalid,
			want: "invalid",
		},
		{
			name: "Invalid IO",
			s:    OrderStatus("Invalid IO"),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("OrderStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

type NewOrderArguments struct {
	StoreID       ID
	MarketplaceID string
	Items         []OrderItem
	Status        OrderStatus
}

func TestNewOrder(t *testing.T) {
	storeID, err := StringToID("d9b5f6d6-3c3b-4a3e-8e6d-4d1a1fbd3f9c")
	if err != nil {
		t.Errorf("Error creating store ID: %v", err)
	}

	tests := []struct {
		name  string
		input NewOrderArguments
		want  *Order
	}{
		{
			name: "New Order",
			input: NewOrderArguments{
				StoreID:       storeID,
				MarketplaceID: "MLB200004894",
				Items:         []OrderItem{{Title: "Item 1", Quantity: 1, Sku: "SKU-1"}},
				Status:        "confirmed",
			},
			want: &Order{
				StoreID:       storeID,
				MarketplaceID: "MLB200004894",
				DateCreated:   time.Now().UTC(),
				Items: []OrderItem{
					{
						Title:    "Item 1",
						Quantity: 1,
						Sku:      "SKU-1",
					},
				},
				Status: Confirmed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrder(tt.input.StoreID, tt.input.MarketplaceID, tt.input.Items, tt.input.Status)
			if err != nil {
				t.Errorf("Error creating order: %v", err)
			}

			if got.ID == uuid.Nil {
				t.Errorf("Order ID is nil")
			}
			if got.StoreID != tt.want.StoreID {
				t.Errorf("Order StoreID = %v, want %v", got.StoreID, tt.want.StoreID)
			}
			if got.MarketplaceID != tt.want.MarketplaceID {
				t.Errorf("Order MarketplaceID = %v, want %v", got.MarketplaceID, tt.want.MarketplaceID)
			}
			if diff := got.DateCreated.Sub(tt.want.DateCreated); diff > time.Second {
				t.Errorf("Order DateCreated = %v, want %v", got.DateCreated.Format(time.RFC3339), tt.want.DateCreated.Format(time.RFC3339))
			}
			// Asserting that the order items are equal is a bit more complex.
			// We need to compare the length of the slices and then compare each item.
			if len(got.Items) != len(tt.want.Items) {
				t.Errorf("Order it has missing items = %v, want %v", len(got.Items), len(tt.want.Items))
			}
			for i, item := range got.Items {
				if item.ID == uuid.Nil {
					t.Errorf("Order Item ID is empty")
				}
				if item.Title != tt.want.Items[i].Title {
					t.Errorf("Order Item Title = %v, want %v", item.Title, tt.want.Items[i].Title)
				}
				if item.Quantity != tt.want.Items[i].Quantity {
					t.Errorf("Order Item Quantity = %v, want %v", item.Quantity, tt.want.Items[i].Quantity)
				}
				if item.Sku != tt.want.Items[i].Sku {
					t.Errorf("Order Item Sku = %v, want %v", item.Sku, tt.want.Items[i].Sku)
				}
			}
			if got.Status != tt.want.Status {
				t.Errorf("Order Status = %v, want %v", got.Status, tt.want.Status)
			}
		})
	}
}
