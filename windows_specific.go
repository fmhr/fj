//go:build windows
// +build windows

package fj

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func runCommandWithTimeout(cmdStrings []string, cnf *Config) ([]byte, string, error) {
	var result string
	cmd := exec.Command(cmdStrings[0], cmdStrings[1:]...)
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
			err := cmd.Process.Kill()
			if err != nil {
				log.Println("Error: ", err, "command:", cmd.String())
				return nil, result, fmt.Errorf("cmd.Process.Kill() failed with: %v", err)
			}
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
