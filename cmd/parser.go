package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	keyValueRe = regexp.MustCompile(`([\p{L}\p{N}\p{M}\p{S}_-]+)=([\d.]+)`)
	dataRe     = regexp.MustCompile(`([^\n]+) = ([\d.]+)`)
)

// ExtractKeyValuePairs はコマンドから出力を、キーと値のマップで返します。
func ExtractKeyValuePairs(msg string) (keys SliceMap, err error) {
	// 例: "score=100.0 time=1.0"
	keys = NewSliceMap()
	matches := keyValueRe.FindAllStringSubmatch(msg, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		keys.Set(key, value)
	}
	return keys, nil
}

// extractScore 公式toolのvisコマンドの出力からスコアを抽出します。
// 小数点を含む数字に対応しているけど、stringで返すので注意
func extractData(src string) (map[string]string, error) {
	if src == "" {
		log.Println("Error: source string is empty")
		return nil, fmt.Errorf("source string is empty")
	}
	// 例: "Score = 100"
	matches := dataRe.FindAllStringSubmatch(src, -1)
	data := make(map[string]string)
	for _, match := range matches {
		key := strings.TrimSpace(match[1])
		value := strings.TrimSpace(match[2])
		data[key] = value
	}
	if len(data) == 0 {
		log.Println("Error: no data found")
		log.Println("Source:", src)
		return nil, fmt.Errorf("no data found")
	}
	return data, nil
}
