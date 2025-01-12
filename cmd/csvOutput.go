package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elliotchance/orderedmap/v2"
)

// CsvOutput csvファイルに出力する
func CsvOutput(datas []*orderedmap.OrderedMap[string, string]) error {
	now := time.Now()
	filename := fmt.Sprintf("%s.csv", fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))
	filename = fmt.Sprintf("fj/result/%s", filename)
	if _, err := os.Stat("fj/result"); os.IsNotExist(err) {
		if err := os.Mkdir("fj/result", 0755); err != nil {
			return err
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// CSVファイルに書き込む
	writer := csv.NewWriter(file)
	defer writer.Flush()

	heagders := append([]string{}, datas[0].Keys()...)
	if err := writer.Write(heagders); err != nil {
		return err
	}

	for _, data := range datas {
		values := make([]string, 0)
		for _, key := range heagders {
			v, ok := data.Get(key)
			if !ok {
				values = append(values, "")
				continue
			}
			values = append(values, numToString(v))
		}
		if err := writer.Write(values); err != nil {
			return err
		}
	}
	log.Println("save csv file:", filename)
	return nil
}

func numToString(num any) string {
	switch v := num.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		if float64(int(v)) == v {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%f", v)
	case string:
		return string(v)
	case []byte:
		return string(v)
	default:
		return ""
	}
}
