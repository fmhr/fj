package fj

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func CsvOutput(datas []map[string]float64) {
	now := time.Now()
	filename := fmt.Sprintf("fj/data/result_%s.csv", fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(":", err)
	}
	defer file.Close()

	// CSVファイルに書き込む
	writer := csv.NewWriter(file)
	defer writer.Flush()

	heagders := make([]string, 0)
	for key := range datas[0] {
		heagders = append(heagders, key)
	}
	if err := writer.Write(heagders); err != nil {
		log.Fatal(":", err)
	}

	for _, data := range datas {
		values := make([]string, 0)
		for _, key := range heagders {
			values = append(values, fmt.Sprintf("%f", data[key]))
		}
		if err := writer.Write(values); err != nil {
			log.Fatal(":", err)
		}
	}
	log.Println("save csv file:", filename)
}
