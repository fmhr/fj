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

func TraceErrorf(format string, err error) error {
	if !DebugModeEnabled() {
		return fmt.Errorf(format, err)
	}
	_, file, line, _ := runtime.Caller(1) // 1つ上のスタックフレームの情報を取得
	filename := filepath.Base(file)
	return fmt.Errorf("[%s:%d]: %s", filename, line, fmt.Errorf(format, err).Error())
}

func TraceError(msg string) error {
	if !DebugModeEnabled() {
		return fmt.Errorf(msg)
	}
	_, file, line, _ := runtime.Caller(1) // 1つ上のスタックフレームの情報を取得
	filename := filepath.Base(file)
	return fmt.Errorf("[%s:%d]: %s", filename, line, msg)
}
