package cmd

import (
	"log"
	"runtime"
)

// OS別にコマンドを変更

func createCommand(cmdStr string) (cmdArgs []string) {
	var timeCmd string
	switch runtime.GOOS {
	case "windows":
		// Windows では　time コマンドが使えないのでそのまま実行
		cmdArgs = []string{"cmd", "/C", cmdStr}
		return cmdArgs
	case "linux":
		// Linux では　GNU time を使う
		timeCmd = "/usr/bin/time -v"
	case "darwin", "freebsd":
		// macOSとFreeBSD では　BSD time を使う
		timeCmd = "/usr/bin/time -l"
	default:
		log.Println("OS not supported")
	}

	if timeCmd != "" {
		cmdStr = timeCmd + " " + cmdStr
	}

	cmdArgs = []string{"/bin/sh", "-c", cmdStr}

	return cmdArgs
}
