package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestEntity(t *testing.T) {
	t.Run("generate a new ID", func(t *testing.T) {
		e := NewID()
		if e == uuid.Nil {
			t.Errorf("got %v, want not nil", e)
		}
	})

	t.Run("convert a string to an entity ID", func(t *testing.T) {
		s := "123e4567-e89b-12d3-a456-426655440000"
		e, err := StringToID(s)
		if err != nil {
			t.Fatalf("unexpected error converting string to entity ID: %v", err)
		}
		if e.String() != s {
			t.Errorf("got %s, want %s", e, s)
		}
	})
}
