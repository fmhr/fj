package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

// ExtractKeyValuePairs はコマンドから出力を、キーと値のマップで返します。
func ExtractKeyValuePairs(m *orderedmap.OrderedMap[string, any], msg string) (keys []string, err error) {
	// 例: "score=100.0 time=1.0"
	re := regexp.MustCompile(`(\w+)=([\d.]+)`)
	matches := re.FindAllStringSubmatch(msg, -1)
	for _, match := range matches {
		key := match[1]
		value, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			log.Println("Error: ", err, "key:", key, "from:", match[2])
			return keys, fmt.Errorf("failed to convert %s to number: %s", match[2], err)
		}
		m.Set(key, value)
		keys = append(keys, key)
	}
	return keys, nil
}

// extractScore 公式toolのvisコマンドの出力からスコアを抽出します。
func extractData(src string) (map[string]float64, error) {
	// 例: "Score = 100"
	re := regexp.MustCompile(`([^\n]+) = ([\d.]+)`)
	matches := re.FindAllStringSubmatch(src, -1)
	data := make(map[string]float64)
	for _, match := range matches {
		key := strings.TrimSpace(match[1])
		value, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			log.Println("Error: ", err, "key:", key, "from:", match[2])
			return nil, fmt.Errorf("failed to convert %s to number: %s", match[2], err)
		}
		data[key] = value
	}
	if len(data) == 0 {
		log.Println("Error: no data found")
		log.Println("Source:", src)
		return nil, fmt.Errorf("no data found")
	}
	return data, nil
}
