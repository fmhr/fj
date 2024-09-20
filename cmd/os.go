package cmd

import (
	"runtime"
)

// OS別にコマンドを変更

func createCommand(cmdStr string) (cmdArgs []string) {
	switch runtime.GOOS {
	case "windows":
		cmdArgs = []string{"cmd", "/C", cmdStr}
	case "linux", "darwin", "freebsd":
		cmdArgs = []string{"/bin/sh", "-c", cmdStr}
	default:
		cmdArgs = []string{"/bin/sh", "-c", cmdStr}
	}
	return cmdArgs
}
