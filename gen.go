package fj

import (
	"fmt"
	"os"
	"sync"
)

func Gen(cnf *Config, seed int) error {
	if cnf.GenPath == "" {
		return fmt.Errorf("GenPath is not set. please set GenPath: {} in fj/config.toml")
	}
	err := gen(cnf, seed)
	return err
}

var genMutex sync.Mutex

// seedを書き込んだd.txtをgenにわたすとin/0000.txtが生成される
// これをin2/{seed}.txtにリネームする
// config.InfilePathをin2/に変更する
func gen(cnf *Config, seed int) error {
	genMutex.Lock()
	defer genMutex.Unlock()
	// in2/がなければ作成
	path := "in2"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// genがあるか確認
	_, err := os.Stat(cnf.GenPath)
	if err != nil {
		return err
	}
	// seedを書き込んだ{seed}.txtを作成
	seedfile := "seed.txt"
	err = writeIntToFile(seed, seedfile)
	if err != nil {
		return err
	}
	// genを実行
	cmdStr := fmt.Sprintf("%s %s", cnf.GenPath, seedfile)
	cmd := createCommand(cmdStr)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run gen: %s", err)
	}
	// in/0000.txtをin/{seed}.txtにリネーム
	err = os.Rename("in/0000.txt", fmt.Sprintf("in2/%04d.txt", seed))
	if err != nil {
		return fmt.Errorf("failed to rename file: %s", err)
	}
	// cnf.InfilePathを更新
	cnf.InfilePath = "in2/"
	// (seed.txt)を削除
	return nil
}

func writeIntToFile(n int, filename string) error {
	data := fmt.Sprintf("%d", n)
	return os.WriteFile(filename, []byte(data), 0644)
}
