package fj

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// normalRun は指定された設定とシードに基づいてコマンドを実行する
// normal モード用
func normalRun(cnf *Config, seed int) ([]byte, error) {
	if cnf.Cmd == "" {
		return nil, fmt.Errorf("config.Cmd is empty")
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	if _, err := os.Stat(inputfile); err != nil {
		return []byte{}, fmt.Errorf("input file [%s] does not exist", inputfile)
	}

	if err := checkOutputFolder(cnf.OutfilePath); err != nil {
		return nil, err
	}
	if !isExist(inputfile) {
		return nil, fmt.Errorf("input file [%s] does not exist", inputfile)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmdStr := fmt.Sprintf("%s < %s > %s", cnf.Cmd, inputfile, outputfile)
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmdStr := fmt.Sprintf("%s < %s > %s", cnf.Cmd, inputfile, outputfile)
		cmd = exec.Command("/bin/sh", "-c", cmdStr)
	}

	out, err := runCommandWithTimeout(cmd, cnf)
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmd, err)
	}
	return out, nil
}

func runCommandWithTimeout(cmd *exec.Cmd, cnf *Config) ([]byte, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("cmd.Start() failed with: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(cnf.TimeLimitMS) * time.Millisecond):
		if cmd.Process != nil {
			err := cmd.Process.Kill()
			if err != nil {
				return nil, fmt.Errorf("failed to kill process: %v", err)
			}
			return nil, fmt.Errorf("process killed as timeout reached")
		}
	case err := <-done:
		if err != nil {
			return stdoutBuf.Bytes(), fmt.Errorf("cmd.Wait() failed with: %v \n%v", err, stderrBuf.String())
		}
	}
	// エラーがなければ、標準出力を返す
	rtn := stdoutBuf.Bytes()
	rtn = append(rtn, stderrBuf.Bytes()...)
	return rtn, nil
}

func checkOutputFolder(dir string) error {
	if filepath.Clean(dir) != dir || filepath.IsAbs(dir) {
		return fmt.Errorf("invalid output folder: %s", dir)
	}

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

func isExist(file string) bool {
	if filepath.Clean(file) != file || filepath.IsAbs(file) {
		return false
	}

	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Print(err)
		return false
	}
	return true
}
