package fj

import (
	"os"
	"slices"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// DisplayTable はデータをテーブル形式で表示する
func DisplayTable(data []map[string]float64) {
	if len(data) == 0 {
		return
	}

	sortBySeed(&data)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)

	headers := extractHeaders(data)
	table.SetHeader(headers)

	for _, rowMap := range data {
		row := make([]string, 0, len(rowMap))
		for _, key := range headers {
			row = append(row, formatFloat(rowMap[key]))
		}
		table.Append(row)
	}
	table.Render()
}

// sortBySeed はデータをseedでソートする
func sortBySeed(data *[]map[string]float64) {
	sort.Slice(*data, func(i, j int) bool {
		return (*data)[i]["seed"] < (*data)[j]["seed"]
	})
}

// extractHeaders はデータからヘッダーを抽出する
func extractHeaders(data []map[string]float64) []string {
	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
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
func formatFloat(value float64) string {
	if value == float64(int(value)) {
		return strconv.Itoa(int(value))
	}
	return strconv.FormatFloat(value, 'f', 3, 64)
}
