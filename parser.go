package fj

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ExtractKeyValuePairs はコマンドから出力を、キーと値のマップで返します。
func ExtractKeyValuePairs(msg string) (map[string]float64, error) {
	re := regexp.MustCompile(`(\w+)=([\d.]+)`)
	matches := re.FindAllStringSubmatch(msg, -1)

	data := make(map[string]float64)

	for _, match := range matches {
		key := match[1]
		value, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to number: %s", match[2], err)
		}
		data[key] = value
	}
	return data, nil
}

// extractScore 公式toolのvisコマンドの出力からスコアを抽出します。
func extractData(src string) (map[string]float64, error) {
	re := regexp.MustCompile(`([^\n]+) = ([\d.]+)`)
	matches := re.FindAllStringSubmatch(src, -1)
	data := make(map[string]float64)
	for _, match := range matches {
		key := strings.TrimSpace(match[1])
		value, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to number: %s", match[2], err)
		}
		data[key] = value
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("failed to convert %s to number: %s", src, "no data")
	}
	return data, nil
}
