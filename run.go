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

	//inputfile := fmt.Sprintf("%s%04d.txt", cnf.InfilePath, seed)
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	if _, err := os.Stat(inputfile); err != nil {
		return []byte{}, fmt.Errorf("input file [%s] does not exist", inputfile)
	}
	//outputfile := fmt.Sprintf("%s%04d.out", cnf.OutfilePath, seed)
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.out", seed))
	cmdStr := cnf.Cmd + " < " + inputfile + " > " + outputfile
	cmd := exec.Command("sh", "-c", cmdStr)
	log.Println(cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmd, err)
	}
	log.Println(string(out))
	return out, nil
}
