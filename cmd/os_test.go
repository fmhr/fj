package cmd

import (
	"reflect"
	"runtime"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	var shell []string
	if runtime.GOOS == "windows" {
		shell = []string{"cmd", "/C"}
	} else {
		shell = []string{"/bin/sh", "-c"}
	}

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Simple command",
			input: "echo hello",
		},
		{
			name:  "Command with redirection",
			input: "solver < input.txt > output.txt",
		},
		{
			name:  "Empty command string",
			input: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := append([]string{}, shell...)
			expected = append(expected, tt.input)

			actual := createCommand(tt.input)
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		})
	}
}
