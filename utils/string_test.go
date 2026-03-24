package utils

import "testing"

func TestGetOrDefault(t *testing.T) {
	t.Run("non-nil string returns value", func(t *testing.T) {
		s := "hello"
		result := GetOrDefault(&s, "default")
		if result != "hello" {
			t.Errorf("expected %q, got %q", "hello", result)
		}
	})

	t.Run("nil string returns default", func(t *testing.T) {
		result := GetOrDefault(nil, "default")
		if result != "default" {
			t.Errorf("expected %q, got %q", "default", result)
		}
	})

	t.Run("empty string returns empty", func(t *testing.T) {
		s := ""
		result := GetOrDefault(&s, "default")
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})
}
