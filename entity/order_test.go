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
	AccountID     ID
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
		name    string
		input   NewOrderArguments
		want    *Order
		wantErr bool
	}{
		{
			name: "Valid Order",
			input: NewOrderArguments{
				AccountID:     storeID,
				MarketplaceID: "MLB200004894",
				Items:         []OrderItem{{Title: "Item 1", Quantity: 1, Sku: "SKU-1"}},
				Status:        Confirmed,
			},
			want: &Order{
				AccountID:     storeID,
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
			wantErr: false,
		},
		{
			name: "Order with Multiple Items",
			input: NewOrderArguments{
				AccountID:     storeID,
				MarketplaceID: "MLB200004894",
				Items: []OrderItem{
					{Title: "Item 1", Quantity: 1, Sku: "SKU-1"},
					{Title: "Item 2", Quantity: 2, Sku: "SKU-2"},
				},
				Status: Confirmed,
			},
			want: &Order{
				AccountID:     storeID,
				MarketplaceID: "MLB200004894",
				DateCreated:   time.Now().UTC(),
				Items: []OrderItem{
					{Title: "Item 1", Quantity: 1, Sku: "SKU-1"},
					{Title: "Item 2", Quantity: 2, Sku: "SKU-2"},
				},
				Status: Confirmed,
			},
			wantErr: false,
		},
		{
			name: "Order with Variations",
			input: NewOrderArguments{
				AccountID:     storeID,
				MarketplaceID: "MLB200004894",
				Items: []OrderItem{
					{Title: "Item 1", Quantity: 1, Sku: "SKU-1", VariationID: 123},
					{Title: "Item 1", Quantity: 1, Sku: "SKU-1", VariationID: 456},
				},
				Status: Confirmed,
			},
			want: &Order{
				AccountID:     storeID,
				MarketplaceID: "MLB200004894",
				DateCreated:   time.Now().UTC(),
				Items: []OrderItem{
					{Title: "Item 1", Quantity: 1, Sku: "SKU-1", VariationID: 123},
					{Title: "Item 1", Quantity: 1, Sku: "SKU-1", VariationID: 456},
				},
				Status: Confirmed,
			},
			wantErr: false,
		},
		{
			name: "Order with Empty SKU",
			input: NewOrderArguments{
				AccountID:     storeID,
				MarketplaceID: "MLB200004894",
				Items: []OrderItem{
					{Title: "Item 1", Quantity: 1, Sku: ""},
				},
				Status: Confirmed,
			},
			want: &Order{
				AccountID:     storeID,
				MarketplaceID: "MLB200004894",
				DateCreated:   time.Now().UTC(),
				Items: []OrderItem{
					{Title: "Item 1", Quantity: 1, Sku: ""},
				},
				Status: Confirmed,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrder(tt.input.AccountID, tt.input.MarketplaceID, tt.input.Items, tt.input.Status)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewOrder() error = %v, wantErr %v", (err != nil), tt.wantErr)
				return
			}

			if tt.wantErr {
				if got != nil {
					t.Errorf("NewOrder() = %v, want nil", got)
				}
				return
			}

			if got.ID == uuid.Nil {
				t.Errorf("Order ID is nil")
			}

			if got.AccountID != tt.want.AccountID {
				t.Errorf("Order AccountID = %v, want %v", got.AccountID, tt.want.AccountID)
			}

			if got.MarketplaceID != tt.want.MarketplaceID {
				t.Errorf("Order MarketplaceID = %v, want %v", got.MarketplaceID, tt.want.MarketplaceID)
			}

			if diff := got.DateCreated.Sub(tt.want.DateCreated); diff > time.Second {
				t.Errorf("Order DateCreated = %v, want %v", got.DateCreated.Format(time.RFC3339), tt.want.DateCreated.Format(time.RFC3339))
			}

			if len(got.Items) != len(tt.want.Items) {
				t.Errorf("Order has wrong number of items = %v, want %v", len(got.Items), len(tt.want.Items))
				return
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
				if item.VariationID != tt.want.Items[i].VariationID {
					t.Errorf("Order Item VariationID = %v, want %v", item.VariationID, tt.want.Items[i].VariationID)
				}
			}

			if got.Status != tt.want.Status {
				t.Errorf("Order Status = %v, want %v", got.Status, tt.want.Status)
			}
		})
	}
}
