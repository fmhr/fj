package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fmhr/fj/cmd/setup"
)

// normalRun は指定された設定とシードに基づいてコマンドを実行する
// normal モード用
func normalRun(cnf *setup.Config, seed int) ([]byte, string, error) {
	cmd := setup.LanguageSets[cnf.Language].ExeCmd
	if cmd == "" {
		return nil, "", NewStackTraceError(fmt.Sprintf("error: LanguageSets[%s].ExecCmd must not be empty", cnf.Language))
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	if _, err := os.Stat(inputfile); err != nil {
		return nil, "", err
	}

	if err := checkOutputFolder(cnf.OutfilePath); err != nil {
		return nil, "", err
	}

	// コマンドを作成
	cmd += " " + setArgs(cnf.Args) // カスタム引数を追加
	cmdStr := fmt.Sprintf("%s < %s > %s", cmd, inputfile, outputfile)

	cmdStrings := createCommand(cmdStr)

	out, result, err := runCommandWithTimeout(cmdStrings, int(cnf.TimeLimitMS))
	if err != nil {
		//log.Println("Error: ", err, "\nout:", string(out))
		return out, result, fmt.Errorf("cmd.Run() for command [%q] failed with: %v out:%s", cmdStrings, err.Error(), out)
	}
	return out, result, nil
}

// checkOutputFolder は出力フォルダが存在するか確認し、存在しない場合は作成する
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
