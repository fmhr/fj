package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

// ExtractKeyValuePairs はコマンドから出力を、キーと値のマップで返します。
func ExtractKeyValuePairs(m *orderedmap.OrderedMap[string, string], msg string) (keys []string, err error) {
	// 例: "score=100.0 time=1.0"
	re := regexp.MustCompile(`(\w+)=([\d.]+)`)
	matches := re.FindAllStringSubmatch(msg, -1)
	for _, match := range matches {
		key := match[1]
		//value, err := strconv.ParseFloat(match[2], 64)
		value := match[2]
		m.Set(key, value)
		keys = append(keys, key)
	}
	return keys, nil
}

// extractScore 公式toolのvisコマンドの出力からスコアを抽出します。
// 小数点を含む数字に対応	しているけど、stringで返すので注意
func extractData(src string) (map[string]string, error) {
	// 例: "Score = 100"
	re := regexp.MustCompile(`([^\n]+) = ([\d.]+)`)
	matches := re.FindAllStringSubmatch(src, -1)
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
