package cmd

import (
	"fmt"
	"os"
	"runtime"
)

// OS別にコマンドを変更
func createCommand(cmdStr string) (cmdArgs []string) {
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/C", cmdStr}
	} else {
		cmdArgs = []string{"/bin/sh", "-c", cmdStr}
	}
	return cmdArgs
}

// EnsureDirExists はディレクトリが存在するか確認し、存在しない場合は作成します
func EnsureDirExists(dir string) error {
	if dir == "" {
		return nil
	}
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		} else {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dir)
	}
	return nil
}

