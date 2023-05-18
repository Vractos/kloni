package entity

import (
	"testing"
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
