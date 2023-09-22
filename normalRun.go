package fj

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// normalRun は指定された設定とシードに基づいてコマンドを実行する
// normal モード用
func normalRun(cnf *Config, seed int) ([]byte, error) {
	if cnf.Cmd == "" {
		return []byte{}, fmt.Errorf("config.Cmd is empty")
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	if _, err := os.Stat(inputfile); err != nil {
		return []byte{}, fmt.Errorf("input file [%s] does not exist", inputfile)
	}
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.out", seed))

	err := checkOutputFolder(cnf.OutfilePath)
	if err != nil {
		log.Fatal(err)
	}

	cmdStr := cnf.Cmd + " < " + inputfile + " > " + outputfile
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmd, err)
	}
	return out, nil
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
