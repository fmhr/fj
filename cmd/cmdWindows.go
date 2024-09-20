//go:build windows

package cmd

import "os/exec"

func cmdCustom(cmd *exec.Cmd) error {
	cmd.Cancel = func() error {
		return cmd.Process.Kill()
	}
	return nil
}
