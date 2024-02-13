package fj

import (
	"runtime"
)

// OS別のコマンド実行
func createCommand(cmdStr string) (strs []string) {
	if runtime.GOOS == "windows" {
		strs = []string{"cmd", "/C", cmdStr}
	} else {
		strs = []string{"/bin/sh", "-c", cmdStr}
	}
	return strs
}
