package utils

import "testing"

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !Contains(&slice, "a") {
		t.Errorf("Expected slice to contain 'a'")
	}
	if Contains(&slice, "d") {
		t.Errorf("Expected slice to not contain 'd'")
	}
}

func TestChunK(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e", "f", "g"}
	chunks := Chunk(slice, 2)
	if len(chunks) != 4 {
		t.Errorf("Expected 4 chunks, got %d", len(chunks))
	}
	if chunks[0][0] != "a" || chunks[0][1] != "b" {
		t.Errorf("Expected first chunk to be [a b], got %v", chunks[0])
	}
	if chunks[1][0] != "c" || chunks[1][1] != "d" {
		t.Errorf("Expected second chunk to be [c d], got %v", chunks[1])
	}
	if chunks[2][0] != "e" || chunks[2][1] != "f" {
		t.Errorf("Expected third chunk to be [e f], got %v", chunks[2])
	}
	if chunks[3][0] != "g" {
		t.Errorf("Expected fourth chunk to be [g], got %v", chunks[3])
	}
}
