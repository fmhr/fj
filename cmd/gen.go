package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/fmhr/fj/cmd/setup"
)

func Gen(cnf *setup.Config, seed int) error {
	if cnf.GenPath == "" {
		return fmt.Errorf("GenPath is not set. please set GenPath: {} to fj/config.toml")
	}
	return gen(cnf, seed)
}

var genMutex sync.Mutex

// seedを書き込んだd.txtをgenにわたすとin/0000.txtが生成される
// これをin2/{seed}.txtにリネームする
// config.InfilePathをin2/に変更する
func gen(cnf *setup.Config, seed int) error {
	genMutex.Lock()
	defer genMutex.Unlock()
	// in2/がなければ作成
	path := "in2"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("in2ディレクトリの作成に失敗: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("in2ディレクトリの確認に失敗: %w", err)
	}
	// genがあるか確認
	_, err := os.Stat(cnf.GenPath)
	if err != nil {
		return fmt.Errorf("genファイルの確認に失敗: %w", err)
	}
	// seedを書き込んだ{seed}.txtを作成
	seedfile := "seed.txt"
	err = writeIntToFile(seed, seedfile)
	if err != nil {
		return fmt.Errorf("シードファイルの作成に失敗: %w", err)
	}
	// genを実行
	cmdStr := fmt.Sprintf("%s %s", cnf.GenPath, seedfile)
	cmdStrings := createCommand(cmdStr)
	cmd := exec.Command(cmdStrings[0], cmdStrings[1:]...)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("genの実行に失敗: %w", err)
	}
	// in/0000.txtをin/{seed}.txtにリネーム
	err = os.Rename("in/0000.txt", fmt.Sprintf("in2/%04d.txt", seed))
	if err != nil {
		return fmt.Errorf("inファイルのリネームに失敗: %w", err)
	}
	// cnf.InfilePathを更新
	cnf.InfilePath = "in2/"
	// (seed.txt)を削除
	err = os.Remove("seed.txt")
	if err != nil {
		return fmt.Errorf("シードファイルの削除に失敗: %w", err)
	}
	return nil
}

func writeIntToFile(n int, filename string) error {
	data := fmt.Sprintf("%d", n)
	return os.WriteFile(filename, []byte(data), 0644)
}
