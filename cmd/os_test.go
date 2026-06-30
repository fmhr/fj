package cmd

import (
	"reflect"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple command",
			input:    "echo hello",
			expected: []string{"/bin/sh", "-c", "echo hello"},
		},
		{
			name:     "Command with redirection",
			input:    "solver < input.txt > output.txt",
			expected: []string{"/bin/sh", "-c", "solver < input.txt > output.txt"},
		},
		{
			name:     "Empty command string",
			input:    "",
			expected: []string{"/bin/sh", "-c", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := createCommand(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
