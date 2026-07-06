//go:build !windows

package cmd

import (
	"os/exec"
	"runtime"
	"syscall"
)

// getMaxRSS はプロセスの最大メモリ使用量 (Max RSS) をキロバイト (KB) 単位で返します
func getMaxRSS(cmd *exec.Cmd) int64 {
	if cmd.ProcessState == nil {
		return 0
	}
	sysUsage := cmd.ProcessState.SysUsage()
	if sysUsage == nil {
		return 0
	}
	if rusage, ok := sysUsage.(*syscall.Rusage); ok {
		maxRSS := rusage.Maxrss
		// macOS (darwin) では Maxrss の単位がバイト (Byte) なので、キロバイト (KB) に変換する
		if runtime.GOOS == "darwin" {
			maxRSS /= 1024
		}
		return maxRSS // KB
	}
	return 0
}
