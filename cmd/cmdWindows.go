//go:build windows

package cmd

import "os/exec"

// cmdCustom は Windows 環境向けに SysProcAttr や Cancel のデフォルト処理を維持するためのスタブ関数です
func cmdCustom(cmd *exec.Cmd) error {
	return nil
}
