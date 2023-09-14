package main

import (
	"fmt"
	"regexp"
	"strconv"
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
func extractScore(s string) (int, error) {
	re := regexp.MustCompile((`Score\s*=\s*(\d+)`))
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no score found in string: %s", s)
	}
	return strconv.Atoi(matches[1])
}
