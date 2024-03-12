package fj

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

// ExtractKeyValuePairs はコマンドから出力を、キーと値のマップで返します。
func ExtractKeyValuePairs(m *orderedmap.OrderedMap[string, any], msg string) error {
	re := regexp.MustCompile(`(\w+)=([\d.]+)`)
	matches := re.FindAllStringSubmatch(msg, -1)

	for _, match := range matches {
		key := match[1]
		value, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			log.Println("Error: ", err, "key:", key, "from:", match[2])
			return fmt.Errorf("failed to convert %s to number: %s", match[2], err)
		}
		m.Set(key, value)
	}
	return nil
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
			log.Println("Error: ", err, "key:", key, "from:", match[2])
			return nil, fmt.Errorf("failed to convert %s to number: %s", match[2], err)
		}
		data[key] = value
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return data, nil
}
