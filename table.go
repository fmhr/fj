package fj

import (
	"log"
	"os"
	"slices"
	"sort"
	"strconv"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/olekukonko/tablewriter"
)

// DisplayTable はデータをテーブル形式で表示する
func DisplayTable(data []*orderedmap.OrderedMap[string, any]) {
	if len(data) == 0 {
		return
	}

	sortBySeed(&data)

	table := tablewriter.NewWriter(os.Stderr)
	table.SetAutoFormatHeaders(false)
	//	log.Println(data)
	headers := extractHeaders(data)
	table.SetHeader(headers)

	for _, rowMap := range data {
		row := make([]string, 0)
		for _, key := range headers {
			value, _ := rowMap.Get(key)
			row = append(row, formatFloat(value))
		}
		table.Append(row)
	}
	table.Render()
}

// sortBySeed はデータをseedでソートする
func sortBySeed(data *[]*orderedmap.OrderedMap[string, any]) {
	sort.Slice(*data, func(i, j int) bool {
		iseed, _ := (*data)[i].Get("seed")
		jseed, _ := (*data)[j].Get("seed")
		switch iseed.(type) {
		case int:
			return iseed.(int) < jseed.(int)
		case float64:
			return iseed.(float64) < jseed.(float64)
		default:
			log.Fatal("invalid type")
		}
		return false
	})
}

// extractHeaders はデータからヘッダーを抽出する
func extractHeaders(data []*orderedmap.OrderedMap[string, any]) []string {
	headers := make([]string, 0)
	for _, key := range data[0].Keys() {
		headers = append(headers, key)
	}
	// seedを先頭に移動
	seedIndex := slices.Index(headers, "seed")
	headers = append(headers[:seedIndex], headers[seedIndex+1:]...)
	headers = append([]string{"seed"}, headers...)
	// seed以外をソート
	//sort.Strings(headers[1:])

	return headers
}

// formatFloat は小数点以下がh0の場合は整数に変換する
func formatFloat(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', 3, 64)
	default:
		log.Fatal("invalid type")
	}
	return ""
}
