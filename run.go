package fj

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Run は指定された設定とシードに基づいてコマンドを実行する
func Run(cnf *Config, seed int) ([]byte, error) {
	if cnf.Cmd == "" {
		return []byte{}, fmt.Errorf("config.Cmd is empty")
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	if _, err := os.Stat(inputfile); err != nil {
		return []byte{}, fmt.Errorf("input file [%s] does not exist", inputfile)
	}
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.out", seed))
	cmdStr := cnf.Cmd + " < " + inputfile + " > " + outputfile
	cmd := exec.Command("sh", "-c", cmdStr)
	log.Println(cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmd, err)
	}
	return out, nil
}
