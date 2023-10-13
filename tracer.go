package fj

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func DebugModeEnabled() bool {
	return true
	// return os.Getenv("DEBUG") == "true"
}

func ErrorTrace(msg string, err error) error {
	if !DebugModeEnabled() {
		return err
	}
	_, file, line, _ := runtime.Caller(1) // 1つ上のスタックフレームの情報を取得
	filename := filepath.Base(file)
	return fmt.Errorf("[%s:%d]%s: %w", filename, line, msg, err)
}
