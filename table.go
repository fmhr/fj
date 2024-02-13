package fj

import (
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/olekukonko/tablewriter"
)

// DisplayTable はデータをテーブル形式で表示する
func DisplayTable(data []*orderedmap.OrderedMap[string, any]) {
	if len(data) == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetAutoFormatHeaders(false)
	//	log.Println(data)
	headers := extractHeaders(data)
	table.SetHeader(headers)

	for _, rowMap := range data {
		row := make([]string, 0)
		// エラーでrowMap全体がnilの場合がある
		if rowMap == nil {
			continue
		}
		for _, key := range headers {
			value, ok := rowMap.Get(key)
			if !ok {
				// seed(key)がなんらかの理由でない場合はスキップ
				//log.Println("Error no value key:", key)
				//continue
				value = -1
			}
			row = append(row, formatFloat(value))
		}
		table.Append(row)
	}
	table.Render()
}

// extractHeaders はデータからヘッダーを抽出する
func extractHeaders(data []*orderedmap.OrderedMap[string, any]) []string {
	headers := append([]string(nil), data[0].Keys()...)
	// seedを先頭に移動
	seedIndex := slices.Index(headers, "seed")
	headers = append(headers[:seedIndex], headers[seedIndex+1:]...)
	headers = append([]string{"seed"}, headers...)

	return headers
}

func formatFloat(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		// 小数点以下が0の場合は整数に変換
		if v == float64(int(v)) {
			return strconv.Itoa(int(v))
		}
		return strconv.FormatFloat(v, 'f', 3, 64)
	case string:
		return v
	default:
		log.Fatal("invalid type")
	}
	return ""
}
