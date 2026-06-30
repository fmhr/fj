package cmd

// OS別にコマンドを変更
func createCommand(cmdStr string) (cmdArgs []string) {
	cmdArgs = []string{"/bin/sh", "-c", cmdStr}
	return cmdArgs
}
