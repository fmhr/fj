package cmd

import (
	"testing"
)

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.123", "1.123"},
		{"3.0", "3"},
		{"abc", "abc"},
	}

	for _, test := range tests {
		result := formatFloat(test.input)
		if result != test.expected {
			t.Errorf("Expected %s for input %s, but got %s", test.expected, test.input, result)
		}
	}
}
