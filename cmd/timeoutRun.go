package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func runCommandWithTimeout(cmdStrings []string, timelimitMS int) ([]byte, bool, error) {
	if timelimitMS == 0 {
		return nil, false, fmt.Errorf("timelimitMS must not be zero")
	}
	if len(cmdStrings) == 0 {
		return nil, false, fmt.Errorf("cmdStrings must not be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timelimitMS)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdStrings[0], cmdStrings[1:]...)

	// OS毎のコマンド修正
	err := cmdCustom(cmd)

	if err != nil {
		return nil, false, fmt.Errorf("failed to set up command: %v", err)
	}

	// 標準出力と標準エラーの操作
	var outputBuf bytes.Buffer
	multiOut := &outputBuf
	// --show-stderr オプションが指定されている場合はstderrを表示する
	var multiErr io.Writer
	if *showStderr {
		multiErr = io.MultiWriter(os.Stderr, &outputBuf)
	} else {
		multiErr = &outputBuf
	}

	cmd.Stdout = multiOut
	cmd.Stderr = multiErr

	cmd.WaitDelay = 5 * time.Second // Wait 5 seconds before sending SIGKILL　子プロセスとのIOの完了を待つ

	// コマンド実行
	err = cmd.Start()
	if err != nil {
		return nil, false, fmt.Errorf("command execution failed: %v", err)
	}

	err = cmd.Wait()

	//output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		cmd.Cancel() // Ensure the process is terminated on timeout
		return outputBuf.Bytes(), true, fmt.Errorf("command timed out after %d ms", timelimitMS)
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return outputBuf.Bytes(), false, fmt.Errorf("command failed with exit code %d: %v", exitErr.ExitCode(), err)
		}
		return outputBuf.Bytes(), false, fmt.Errorf("command execution failed: %v", err)
	}

	return outputBuf.Bytes(), false, nil
}
