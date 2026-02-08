package cmd

import (
	"log"
	"os/exec"
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
		// GNU time
		timeCmd = "/usr/bin/time -v"
	case "darwin", "freebsd":
		// macOS, FreeBSDではGNU timeがあればそれを使う
		if commandExist("gtime") {
			timeCmd = `gtime -f "time:%E memory=%M KB nCPU:%P "`
		} else {
			// なければ BSD time
			timeCmd = "/usr/bin/time -l"
		}
	default:
		// 未対応OS
		log.Println("OS not supported")
	}

	if timeCmd != "" {
		cmdStr = timeCmd + " " + cmdStr
	}

	cmdArgs = []string{"/bin/sh", "-c", cmdStr}

	return cmdArgs
}

// commandExist はコマンドが存在するか確認する
func commandExist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
