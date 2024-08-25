package cmd

import (
	"testing"
)

func TestRunCommandWithTimeout(t *testing.T) {
	t.Run("Test successful command execution", func(t *testing.T) {
		cmd := []string{"echo", "hello"}
		output, timeout, err := runCommandWithTimeout(cmd, 5000)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if timeout {
			t.Errorf("Expected result timeout is'false', got %v", timeout)
		}
		expectedOutput := "hello\n"
		if string(output) != expectedOutput {
			t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
		}
	})

	t.Run("Test empty command slice", func(t *testing.T) {
		_, _, err := runCommandWithTimeout([]string{}, 1000)
		if err == nil {
			t.Error("Expected a:wn error for empty command slice, but got none")
		}
	})

	t.Run("Test command timeout", func(t *testing.T) {
		cmd := []string{"sleep", "2"}
		_, result, err := runCommandWithTimeout(cmd, 500) // 500 ms
		if err != nil {
			t.Errorf("Did not expect error: %v", err)
		}
		if result == false {
			t.Errorf("Expected result 'Timeout' for timeout, got %v", result)
		}
	})

	t.Run("Test failing command", func(t *testing.T) {
		cmd := []string{"ls", "--fake-option"}
		_, _, err := runCommandWithTimeout(cmd, 5000)
		if err == nil {
			t.Error("Expected an error for failing command, but got none")
		}
	})
}
