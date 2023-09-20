package fj

import (
	"testing"
)

func TestCheckConfigFile(t *testing.T) {
	// 1. Cmdが空の場合
	t.Run("empty cmd", func(t *testing.T) {
		cnf := &Config{Cmd: ""}
		err := checkConfigFile(cnf)
		if err == nil {
			t.Errorf("expected error but got nil")
		}
	})

	// 2. Cmdに値が設定されている場合
	t.Run("valid cmd", func(t *testing.T) {
		cnf := &Config{Cmd: "sample command"}
		err := checkConfigFile(cnf)
		if err != nil {
			t.Errorf("expected nil but got error: %v", err)
		}
	})
}
