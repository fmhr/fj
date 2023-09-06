package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// ExtractKeyValuePairs は文字列を受け取り、キーと値のマップを返します。
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
