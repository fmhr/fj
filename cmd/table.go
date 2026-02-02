package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/olekukonko/tablewriter"
)

// DisplayTable はデータをテーブル形式で表示する
func DisplayTable(data []*orderedmap.OrderedMap[string, string]) error {
	if len(data) == 0 {
		return nil
	}

	headers := extractHeaders(data)
	table := tablewriter.NewWriter(os.Stderr)
	//	log.Println(data)
	table.Header(toInterfaceSlice(headers)...)

	sort.Slice(data, func(i, j int) bool {
		seedI, _ := data[i].Get("seed")
		seedJ, _ := data[j].Get("seed")
		seedIint, _ := strconv.Atoi(seedI)
		seedJint, _ := strconv.Atoi(seedJ)
		return seedIint < seedJint
	})

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
				value = "-1"
			}
			v, err := formatFloat(value)
			if err != nil {
				log.Println("Error formatFloat:", err)
				return err
			}
			row = append(row, v)
		}
		table.Append(row)
	}
	table.Render()
	return nil
}

// extractHeaders はデータからヘッダーを抽出する
func extractHeaders(data []*orderedmap.OrderedMap[string, string]) []string {
	headers := append([]string(nil), data[0].Keys()...)
	// seedを先頭に移動
	seedIndex := slices.Index(headers, "seed")
	headers = append(headers[:seedIndex], headers[seedIndex+1:]...)
	headers = append([]string{"seed"}, headers...)

	return headers
}

func formatFloat(value any) (string, error) {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v), nil
	case float64:
		// 小数点以下が0の場合は整数に変換
		if v == float64(int(v)) {
			return strconv.Itoa(int(v)), nil
		}
		return strconv.FormatFloat(v, 'f', 3, 64), nil
	case string:
		return v, nil
	}
	log.Println("invalid type")
	return "", fmt.Errorf("invalid type")
}

// toInterfaceSlice converts a string slice to an interface slice
func toInterfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
