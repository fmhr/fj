//go:build windows
// +build windows

package fj

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func runCommandWithTimeout(cmdStrings []string, timelimitMS int) ([]byte, string, error) {
	if len(cmdStrings) == 0 {
		return nil, "", fmt.Errorf("cmdStrings must not be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timelimitMS)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdStrings[0], cmdStrings[1:]...)
	output, err := cmd.CombinedOutput()
	// タイムアウトの場合
	if ctx.Err() == context.DeadlineExceeded {
		return output, "TLE", nil
	}

	// タイムアウト以外のエラーの場合
	if err != nil {
		//log.Println("command:", cmd.String(), "output:", string(output))
		return output, "", fmt.Errorf("cmd.CombinedOutput() failed with: %v", err)
	}

	log.Println("pid=", cmd.ProcessState.Pid())

	return output, "Success", nil
}
