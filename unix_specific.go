//go:build linux || darwin
// +build linux darwin

package fj

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"
)

// runCommandWithTimeout は指定されたタイムアウトでコマンドを実行し、標準出力と標準エラー出力の結合された内容と結果文字列、およびエラーを返す
func runCommandWithTimeout(cmdStrings []string, timelimitMS int) ([]byte, string, error) {
	if len(cmdStrings) == 0 {
		return nil, "", fmt.Errorf("cmdStrings must not be empty")
	}

	cmd := exec.Command(cmdStrings[0], cmdStrings[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // -pgidで子プロセスもkillできるようにする

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return nil, "", fmt.Errorf("cmd.Start() failed with: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
		close(errCh)
	}()

	timedOut := false
	select {
	case <-time.After(time.Duration(timelimitMS) * time.Millisecond):
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // -pgidで子プロセスもkill
			//cmd.Process.Kill()
		}
		timedOut = true
	case err := <-errCh:
		if err != nil {
			rtn := make([]byte, 0, len(stdoutBuf.Bytes())+len(stderrBuf.Bytes()))
			rtn = append(rtn, stdoutBuf.Bytes()...)
			rtn = append(rtn, stderrBuf.Bytes()...)
			log.Println("Error: ", err, "command:", cmd.String())
			return rtn, "", fmt.Errorf("cmd.Wait() failed with: %v", err)
		}
	}
	output := append(stdoutBuf.Bytes(), stderrBuf.Bytes()...)

	if timedOut {
		return output, "TLE", nil
	}

	return output, "Sccess", nil
}
