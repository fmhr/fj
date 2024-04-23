//go:build linux || darwin
// +build linux darwin

package fj

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// runCommandWithTimeout は指定されたタイムアウトでコマンドを実行し、標準出力と標準エラー出力の結合された内容と結果文字列、およびエラーを返す
func runCommandWithTimeout(cmdStrings []string, timelimitMS int) ([]byte, string, error) {
	if len(cmdStrings) == 0 {
		return nil, "", fmt.Errorf("cmdStrings must not be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timelimitMS)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdStrings[0], cmdStrings[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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
	child, err := findChildPIDs(cmd.ProcessState.Pid())
	if err != nil {
		log.Println("Error: ", err, "command:", cmd.String())
		return output, "", fmt.Errorf("FindChildPIDs failed with: %v", err)
	}
	log.Println("child=", child)

	return output, "Success", nil
}

// FindChildPIDs は指定された親PIDを持つすべてのプロセスのPIDを返します。
// Linux
func FindChildPIDs(parentPID int) ([]int, error) {
	var childPIDs []int
	procPath := "/proc"

	// /proc 下のすべてのディレクトリを読み取る
	entries, err := os.ReadDir(procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pidDir := entry.Name()
			// ディレクトリ名が数字であることを確認（プロセスIDディレクトリ）
			if _, err := os.Stat(filepath.Join(procPath, pidDir)); os.IsNotExist(err) || !strings.HasPrefix(pidDir, "self") {
				statusFilePath := filepath.Join(procPath, pidDir, "status")
				file, err := os.Open(statusFilePath)
				if err != nil {
					continue // ファイルが開けない場合はスキップ
				}
				defer file.Close()

				// ファイルからPPidを探す
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "PPid:") {
						fields := strings.Fields(line)
						if len(fields) >= 2 && fields[1] == fmt.Sprintf("%d", parentPID) {
							// PPidが一致するプロセスIDを保存
							pid, err := strconv.Atoi(pidDir)
							if err == nil {
								childPIDs = append(childPIDs, pid)
							}
						}
					}
				}
				if err := scanner.Err(); err != nil {
					fmt.Printf("reading status file failed: %v\n", err)
				}
			}
		}
	}

	return childPIDs, nil
}

// macOS, BSD
func findChildPIDs(parentPID int) ([]string, error) {
	cmd := exec.Command("ps", "-ax", "-o", "pid,ppid")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	var childPIDs []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, fmt.Sprintf("%d", parentPID)) {
			parts := strings.Fields(line)
			if len(parts) > 1 && parts[1] == fmt.Sprintf("%d", parentPID) {
				childPIDs = append(childPIDs, parts[0])
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command failed: %v", err)
	}

	return childPIDs, nil
}
