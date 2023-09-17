package fj

import (
	"os"
	"slices"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func DisplayTabel(data []map[string]float64) {
	sortBySeed(&data)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	headers := make([]string, 0, len(data[0]))
	if len(data) > 0 {
		for key := range data[0] {
			headers = append(headers, key)
		}
		seedIndex := slices.Index(headers, "seed")
		if seedIndex != 0 {
			headers[0], headers[seedIndex] = headers[seedIndex], headers[0]
		}
		sort.Strings(headers[1:])
		table.SetHeader(headers)
	}

	for _, rowMap := range data {
		row := make([]string, 0, len(rowMap))
		for _, key := range headers {
			row = append(row, formatFloat(rowMap[key]))
		}
		table.Append(row)
	}
	table.Render()
}

func sortBySeed(data *[]map[string]float64) {
	sort.Slice(*data, func(i, j int) bool {
		return (*data)[i]["seed"] < (*data)[j]["seed"]
	})
}

// 小数点以下が０のときintとして表示させる
func formatFloat(value float64) string {
	if value == float64(int(value)) {
		return strconv.Itoa(int(value))
	}
	return strconv.FormatFloat(value, 'f', 3, 64)
}
