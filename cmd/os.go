package cmd

import "runtime"

// OS別にコマンドを変更
func createCommand(cmdStr string) (cmdArgs []string) {
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/C", cmdStr}
	} else {
		cmdArgs = []string{"/bin/sh", "-c", cmdStr}
	}
	return cmdArgs
}
