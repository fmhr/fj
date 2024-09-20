package cmd

import (
	"context"
	"fmt"
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

	// OS-specific setup
	err := cmdCustom(cmd)

	if err != nil {
		return nil, false, fmt.Errorf("failed to set up command: %v", err)
	}

	cmd.WaitDelay = 5 * time.Second // Wait 5 seconds before sending SIGKILL　子プロセスとのIOの完了を待つ
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		cmd.Cancel() // Ensure the process is terminated on timeout
		return output, true, fmt.Errorf("command timed out after %d ms", timelimitMS)
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return output, false, fmt.Errorf("command failed with exit code %d: %v", exitErr.ExitCode(), err)
		}
		return output, false, fmt.Errorf("command execution failed: %v", err)
	}

	return output, false, nil
}
