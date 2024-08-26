package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fmhr/fj/cmd/setup"
)

// normalRun は指定された設定とシードに基づいてコマンドを実行する
// normal モード用
func normalRun(cnf *setup.Config, seed int) ([]byte, bool, error) {
	cmd := cnf.ExecuteCmd
	if cmd == "" {
		return nil, false, NewStackTraceError("error: Not found execute comand")
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	if _, err := os.Stat(inputfile); err != nil {
		return nil, false, err
	}

	if err := checkOutputFolder(cnf.OutfilePath); err != nil {
		return nil, false, err
	}

	// コマンドを作成
	cmd += " " + setArgs(cnf.Args) // カスタム引数を追加
	cmdStr := fmt.Sprintf("%s < %s > %s", cmd, inputfile, outputfile)

	cmdStrings := createCommand(cmdStr)
	out, timeout, err := runCommandWithTimeout(cmdStrings, int(cnf.TimeLimitMS))
	if err != nil || timeout {
		//log.Println("Error: ", err, "\nout:", string(out))
		log.Println("TIMEOUT")
		return out, timeout, fmt.Errorf("cmd.Run() for command [%q] failed with: %v out:%s", cmdStrings, err.Error(), out)
	}
	return out, timeout, nil
}

// checkOutputFolder は出力フォルダが存在するか確認し、存在しない場合は作成する
func checkOutputFolder(dir string) error {
	// dirが空の時
	if dir == "" {
		return nil
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
