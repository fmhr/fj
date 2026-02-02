package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fmhr/fj/cmd/setup"
)

// normalRun は指定された設定とシードに基づいてコマンドを実行する
// normal モード用
func normalRun(cnf *setup.Config, seed int) ([]byte, bool, error) {
	return normalRunWithTime(cnf, seed)
}

// normalRunWithTime は実行時間を計測しながらコマンドを実行する
func normalRunWithTime(cnf *setup.Config, seed int) ([]byte, bool, error) {
	cmd := cnf.ExecuteCmd
	if cmd == "" {
		return nil, false, NewStackTraceError("error: Not found execute comand")
	}
	inputfile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outputfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	if _, err := os.Stat(inputfile); err != nil {
		if os.IsNotExist(err) {
			return nil, false, fmt.Errorf("input file not found: %s", inputfile)
		}
		return nil, false, err
	}

	if err := checkOutputFolder(cnf.OutfilePath); err != nil {
		return nil, false, err
	}

	// コマンドを作成
	cmd += " " + setArgs(cnf.Args) // カスタム引数を追加
	cmdStr := fmt.Sprintf("%s < %s > %s", cmd, inputfile, outputfile)
	cmdStrings := createCommand(cmdStr)

	// 実行時間を計測
	startTime := time.Now()
	out, timeout, err := runCommandWithTimeout(cmdStrings, int(cnf.TimeLimitMS))
	elapsed := time.Since(startTime)

	if err != nil {
		log.Println("Error: ", err)
		if len(out) > 0 {
			log.Println("out:", string(out))
		}
		return out, timeout, fmt.Errorf("cmd.Run() for command [%q] failed with: %w out:%s", cmdStrings, err, out)
	}
	if timeout {
		log.Printf("処理が制限時間%dmsを超えたためタイムアウトしました", int(cnf.TimeLimitMS))
		return out, timeout, fmt.Errorf("TIMEOUT %dms", int(cnf.TimeLimitMS))
	}

	// 実行時間を出力に追加（秒単位）
	timeStr := fmt.Sprintf("time=%.3f", elapsed.Seconds())
	out = append(out, []byte("\n"+timeStr)...)

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
