package entity

import (
	"testing"
)

func TestStore(t *testing.T) {
	t.Run("creates store entity", func(t *testing.T) {
		email := "store@test.com"
		name := "Test Store"
		s, err := NewStore(email, name)
		if err != nil {
			t.Fatalf("unexpected error creating store: %v", err)
		}

		if s.ID.String() == "" {
			t.Errorf("got %v, want not nil", s.ID)
		}
		if s.Email != email {
			t.Errorf("got %s, want %s", s.Email, email)
		}
		if s.Name != name {
			t.Errorf("got %s, want %s", s.Name, name)
		}
	})
}
