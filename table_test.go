package fj

import (
	"testing"
)

func TestSortBySeed(t *testing.T) {
	data := []map[string]float64{
		{"seed": 3, "value": 3.0},
		{"seed": 1, "value": 1.123},
		{"seed": 2, "value": 2.234},
	}

	sortBySeed(&data)

	if data[0]["seed"] != 1 || data[1]["seed"] != 2 || data[2]["seed"] != 3 {
		t.Fatalf("Data was not sorted correctly: %v", data)
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1.123, "1.123"},
		{3.0, "3"},
	}

	for _, test := range tests {
		result := formatFloat(test.input)
		if result != test.expected {
			t.Errorf("Expected %s for input %f, but got %s", test.expected, test.input, result)
		}
	}
}
