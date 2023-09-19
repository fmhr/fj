package fj

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

func Gen(cnf *config, seed int) {
	if cnf.GenPath == "" {
		log.Fatal("GenPath is not set. please set GenPath: {} in fj_config.toml")
	}
	err := gen(cnf, seed)
	if err != nil {
		log.Println(err)
	}
}

var genMutex sync.Mutex

// seedを書き込んだd.txtをgenにすとin/0000.txtが生成される
func gen(cnf *config, seed int) error {
	genMutex.Lock()
	defer genMutex.Unlock()
	_, err := os.Stat(cnf.GenPath)
	if err != nil {
		return fmt.Errorf("gen not found: %s", cnf.GenPath)
	}
	// seedを書き込んだd.txtをgenによませるとin/0000.txtが生成される
	filename := fmt.Sprintf("%d.txt", seed)
	err = writeIntToFile(seed, filename)
	// d.txtを削除
	defer os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to write seed to file: %s", err)
	}
	// genを実行
	cmd := fmt.Sprintf("%s %s", cnf.GenPath, filename)
	err = exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		return fmt.Errorf("failed to run gen: %s", err)
	}
	// in/0000.txtをin/{seed}.txtにリネーム
	err = os.Rename("in/0000.txt", fmt.Sprintf("in/%04d.txt", seed))
	if err != nil {
		return fmt.Errorf("failed to rename file: %s", err)
	}
	// cnf.InfilePathを更新
	cnf.InfilePath = "in/" + fmt.Sprintf("%04d.txt", seed)
	return nil
}

func writeIntToFile(n int, filename string) error {
	data := fmt.Sprintf("%d", n)
	return os.WriteFile(filename, []byte(data), 0644)
}
