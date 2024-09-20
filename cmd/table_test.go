package cmd

import (
	"testing"
)

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1.123, "1.123"},
		{3.0, "3"},
	}

	for _, test := range tests {
		result, err := formatFloat(test.input)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if result != test.expected {
			t.Errorf("Expected %s for input %f, but got %s", test.expected, test.input, result)
		}
	}
}
