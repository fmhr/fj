package fj

import (
	"os/exec"
	"runtime"
)

// OS別のコマンド実行
func createCommand(cmdStr string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", cmdStr)
	}
	return exec.Command("sh", "-c", cmdStr)
}
