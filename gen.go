package fj

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func Gen(cnf *Config, seed int) error {
	if cnf.GenPath == "" {
		log.Fatal("GenPath is not set. please set GenPath: {} in fj_config.toml")
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
	path := "in2"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %s", err)
		}
	} else if err != nil {
		log.Fatalf("error checking if directory exists: %s", err)
	}
	_, err := os.Stat(cnf.GenPath)
	if err != nil {
		return fmt.Errorf("gen not found: %s", cnf.GenPath)
	}
	// seedを書き込んだ{seed}.txtをgenによませるとin/0000.txtが生成される
	filename := fmt.Sprintf("%d.txt", seed)
	err = writeIntToFile(seed, filename)
	if err != nil {
		return fmt.Errorf("failed to write seed to file: %s", err)
	}
	// genを実行
	cmdStr := fmt.Sprintf("%s %s", cnf.GenPath, filename)
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
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to remove file: %s", err)
	}
	err = os.Remove("in/0000.txt")
	if err != nil {
		return fmt.Errorf("failed to remove file: %s", err)
	}
	return nil
}

func writeIntToFile(n int, filename string) error {
	data := fmt.Sprintf("%d", n)
	return os.WriteFile(filename, []byte(data), 0644)
}
