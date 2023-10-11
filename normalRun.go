package fj

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// normalRun は指定された設定とシードに基づいてコマンドを実行する
// normal モード用
func normalRun(cnf *Config, seed int) ([]byte, error) {
	if cnf.Cmd == "" {
		return nil, fmt.Errorf("config.Cmd is empty")
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.out", seed))

	if _, err := os.Stat(inputfile); err != nil {
		return []byte{}, fmt.Errorf("input file [%s] does not exist", inputfile)
	}

	if err := checkOutputFolder(cnf.OutfilePath); err != nil {
		return nil, err
	}
	if !isExist(inputfile) {
		return nil, fmt.Errorf("input file [%s] does not exist", inputfile)
	}

	cmdArgs := []string{"-c", cnf.Cmd, "<", inputfile, ">", outputfile}
	cmd := exec.Command("/bin/sh", cmdArgs...)

	out, err := runCommandWithTimeout(cmd, cnf)
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmd, err)
	}
	return out, nil
}

func runCommandWithTimeout(cmd *exec.Cmd, cnf *Config) ([]byte, error) {
	timeout := time.Duration(cnf.TimeLimitMS) * time.Millisecond
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		return out.Bytes(), fmt.Errorf("cmd.Start() failed with: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		if cmd.Process != nil {
			err := cmd.Process.Kill()
			if err != nil {
				return out.Bytes(), fmt.Errorf("failed to kill process: %v", err)
			}
			return out.Bytes(), fmt.Errorf("process killed as timeout reached")
		}
	case err := <-done:
		if err != nil {
			return out.Bytes(), fmt.Errorf("cmd.Wait() failed with: %v stderr:%v", err, out.String())
		}
	}

	return out.Bytes(), nil
}

func checkOutputFolder(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func isExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
