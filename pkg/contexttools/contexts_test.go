package contexttools

import (
	"context"
	"testing"
)

func TestRetrieveStoreIDFromCtx(t *testing.T) {
	t.Run("valid store ID in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ContextKeyStoreId, "store-123")
		id, err := RetrieveStoreIDFromCtx(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != "store-123" {
			t.Errorf("expected %q, got %q", "store-123", id)
		}
	})

	t.Run("empty store ID returns error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ContextKeyStoreId, "")
		_, err := RetrieveStoreIDFromCtx(ctx)
		if err == nil {
			t.Fatal("expected error for empty store ID, got nil")
		}
	})

	t.Run("missing store ID panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for missing context value, got none")
			}
		}()
		ctx := context.Background()
		RetrieveStoreIDFromCtx(ctx)
	})
}
