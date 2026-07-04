package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func JsonOutput(datas []SliceMap) error {
	fileContent, err := json.MarshalIndent(datas, "", " ")
	if err != nil {
		log.Println("json marshal error")
		return err
	}
	err = EnsureDirExists("fj/data/")
	if err != nil {
		log.Println("create dir error")
		return err
	}
	now := time.Now()
	filename := fmt.Sprintf("fj/data/result_%s.json", now.Format("20060102_150405"))
	err = os.WriteFile(filename, fileContent, 0644)
	if err != nil {
		log.Println("json write error")
		return err
	}
	log.Println("success save json file:", filename)
	return nil
}
