//go:build linux || darwin
// +build linux darwin

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// runCommandWithTimeout は指定されたタイムアウトでコマンドを実行し、標準出力と標準エラー出力の結合された内容と結果文字列、およびエラーを返す
// ここのタイムアウトは強制終了で、問題のTLEとは異なる
func runCommandWithTimeout(cmdStrings []string, timelimitMS int) ([]byte, string, error) {
	if len(cmdStrings) == 0 {
		return nil, "", fmt.Errorf("cmdStrings must not be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timelimitMS)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdStrings[0], cmdStrings[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = 5 * time.Second
	output, err := cmd.CombinedOutput()
	// タイムアウトの場合
	if ctx.Err() == context.DeadlineExceeded {
		return output, "Timeout", nil
	}

	// タイムアウト以外のエラーの場合
	if err != nil {
		//log.Println("command:", cmd.String(), "output:", string(output))
		return output, "", fmt.Errorf("cmd.CombinedOutput() failed with: %v", err)
	}
	return output, "Success", nil
}
