package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vractos/kloni/pkg/contexttools"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func TestCustomClaimsValidate(t *testing.T) {
	t.Run("valid claims", func(t *testing.T) {
		claims := CustomClaims{
			InternalID: "store-123",
			Scope:      "read:announcements",
		}
		err := claims.Validate(context.Background())
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("empty internal ID", func(t *testing.T) {
		claims := CustomClaims{
			InternalID: "",
			Scope:      "read:announcements",
		}
		err := claims.Validate(context.Background())
		if err == nil {
			t.Error("expected error for empty internal_id, got nil")
		}
	})
}

func TestHasScope(t *testing.T) {
	claims := CustomClaims{
		InternalID: "store-123",
		Scope:      "read:announcements write:store admin",
	}

	tests := []struct {
		name     string
		scope    string
		expected bool
	}{
		{"scope present", "read:announcements", true},
		{"scope present - write", "write:store", true},
		{"scope present - admin", "admin", true},
		{"scope absent", "delete:store", false},
		{"empty scope", "", false},
		{"partial match", "read", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := claims.HasScope(tt.scope)
			if result != tt.expected {
				t.Errorf("HasScope(%q) = %v, want %v", tt.scope, result, tt.expected)
			}
		})
	}
}

func TestAddStoreIDToCtx(t *testing.T) {
	t.Run("extracts internal_id and adds to context", func(t *testing.T) {
		expectedStoreID := "store-uuid-123"

		customClaims := &CustomClaims{
			InternalID: expectedStoreID,
			Scope:      "read:all",
		}

		validatedClaims := &validator.ValidatedClaims{
			CustomClaims: customClaims,
		}

		var capturedStoreID string
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedStoreID = r.Context().Value(contexttools.ContextKeyStoreId).(string)
			w.WriteHeader(http.StatusOK)
		})

		handler := AddStoreIDToCtx(nextHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), jwtmiddleware.ContextKey{}, validatedClaims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if capturedStoreID != expectedStoreID {
			t.Errorf("expected store ID %q, got %q", expectedStoreID, capturedStoreID)
		}
	})
}
