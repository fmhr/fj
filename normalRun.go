package fj

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// normalRun は指定された設定とシードに基づいてコマンドを実行する
// normal モード用
func normalRun(cnf *Config, seed int) ([]byte, string, error) {
	if cnf.Cmd == "" {
		return nil, "", NewStackTraceError("config.Cmd is empty")
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	if _, err := os.Stat(inputfile); err != nil {
		return nil, "", err
	}

	if err := checkOutputFolder(cnf.OutfilePath); err != nil {
		return nil, "", err
	}

	cmdStr := fmt.Sprintf("%s < %s > %s", cnf.Cmd, inputfile, outputfile)

	cmdStrings := createCommand(cmdStr)

	out, result, err := runCommandWithTimeout(cmdStrings, cnf)
	if err != nil {
		return out, result, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmdStrings, err)
	}
	return out, result, nil
}

func runCommandWithTimeout(cmdStrings []string, cnf *Config) ([]byte, string, error) {
	var result string
	cmd := exec.Command(cmdStrings[0], cmdStrings[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if err := cmd.Start(); err != nil {
		return nil, result, fmt.Errorf("cmd.Start() failed with: %v", err)
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
		close(errCh)
	}()
	var err error
	select {
	case <-time.After(time.Duration(cnf.TimeLimitMS) * time.Millisecond):
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // -pgidで子プロセスもkill
			result = "TLE"
		}
	case err := <-errCh:
		if err != nil {
			log.Println("Error: ", err, "command:", cmd.String())
			return nil, result, fmt.Errorf("cmd.Wait() failed with: %v", err)
		}
	}

	// タイムアウトして、-race をつけたときに、WARNING: DATA RACE がでるのは避けられない
	// https://github.com/golang/go/issues/22757
	rtn := make([]byte, 0, len(stdoutBuf.Bytes())+len(stderrBuf.Bytes()))
	rtn = append(rtn, stdoutBuf.Bytes()...)
	rtn = append(rtn, stderrBuf.Bytes()...)
	return rtn, result, err
}

func checkOutputFolder(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				return fmt.Errorf("failed to create output folder: %w", err)
			}
		} else {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("path is not directory: %s", dir)
	}
	return nil
}
