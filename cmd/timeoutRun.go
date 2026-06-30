package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func runCommandWithTimeout(cmdStrings []string, timelimitMS int) ([]byte, bool, int64, error) {
	if timelimitMS == 0 {
		return nil, false, 0, fmt.Errorf("timelimitMS must not be zero")
	}
	if len(cmdStrings) == 0 {
		return nil, false, 0, fmt.Errorf("cmdStrings must not be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timelimitMS)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdStrings[0], cmdStrings[1:]...)

	// OS毎のコマンド修正
	err := cmdCustom(cmd)

	if err != nil {
		return nil, false, 0, fmt.Errorf("failed to set up command: %v", err)
	}

	// 標準出力と標準エラーの操作
	var outputBuf bytes.Buffer

	cmd.Stderr = &outputBuf

	// プロセス終了後にI/Oパイプがクローズされるのを待つ最大時間 (ゾンビプロセスによるハングアップ防止用)
	cmd.WaitDelay = 500 * time.Millisecond

	// コマンド実行
	err = cmd.Start()
	if err != nil {
		return nil, false, 0, fmt.Errorf("command execution failed: %v", err)
	}

	err = cmd.Wait()

	maxRSS := getMaxRSS(cmd)

	//output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return outputBuf.Bytes(), true, maxRSS, nil
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return outputBuf.Bytes(), false, maxRSS, fmt.Errorf("command failed with exit code %d: %v", exitErr.ExitCode(), err)
		}
		return outputBuf.Bytes(), false, maxRSS, fmt.Errorf("command execution failed: %v", err)
	}
	// 標準エラーの内容をログに出力 --stderrオプション
	if showStderr != nil && *showStderr {
		log.Println("stderr:\n", outputBuf.String())
	}

	return outputBuf.Bytes(), false, maxRSS, nil
}
