package utils

import "testing"

func TestPercentOf(t *testing.T) {
	if PercentOf(10, 100) != 10 {
		t.Errorf("Expected 10%% of 100 to be 10, got %f", PercentOf(10, 100))
	}
	if PercentOf(10, 200) != 5 {
		t.Errorf("Expected 10%% of 200 to be 5, got %f", PercentOf(10, 200))
	}
}

func TestPercent(t *testing.T) {
	if Percent(10, 100) != 10 {
		t.Errorf("Expected 10%% of 100 to be 10, got %f", Percent(10, 100))
	}
	if Percent(10, 200) != 20 {
		t.Errorf("Expected 10%% of 200 to be 20, got %f", Percent(10, 200))
	}
}
