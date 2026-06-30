//go:build windows

package cmd

import "os/exec"

// getMaxRSS は Windows では現状 0 を返します
func getMaxRSS(cmd *exec.Cmd) int64 {
	return 0
}
