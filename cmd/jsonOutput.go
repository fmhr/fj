package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elliotchance/orderedmap/v2"
)

func JsonOutput(datas []*orderedmap.OrderedMap[string, any]) error {
	fileContent, err := json.MarshalIndent(datas, "", " ")
	if err != nil {
		log.Println("json marshal error")
		return err
	}
	err = createDirIfNotExist("fj/data/")
	if err != nil {
		log.Println("create dir error")
		return err
	}
	now := time.Now()
	filename := fmt.Sprintf("fj/data/result_%s.json", fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))
	err = os.WriteFile(filename, fileContent, 0644)
	if err != nil {
		log.Println("json write error")
		return err
	}
	log.Println("success save json file:", filename)
	return nil
}
