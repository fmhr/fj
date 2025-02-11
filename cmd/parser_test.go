package cmd

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/elliotchance/orderedmap/v2"
)

func TestExtractKeyValuePairs(t *testing.T) {
	tests := []struct {
		input  string
		output map[string]float64
		err    bool
	}{
		{
			input:  "key1=5.0 key2=10.0",
			output: map[string]float64{"key1": 5.0, "key2": 10.0},
			err:    false,
		},
		{
			input:  "key1=abc key2=10.0",
			output: map[string]float64{"key2": 10.0},
			err:    false,
		},
		{
			input:  "key1=5.0",
			output: map[string]float64{"key1": 5.0},
			err:    false,
		},
	}

	for _, test := range tests {
		result := orderedmap.NewOrderedMap[string, string]()
		keys, err := ExtractKeyValuePairs(result, test.input)
		if (err != nil) != test.err {
			t.Errorf("Expected error %v, but got %v", test.err, err)
		}
		for key, value := range test.output {
			vstr, _ := result.Get(key)
			v, _ := strconv.ParseFloat(vstr, 64)
			if value != v {
				t.Error("Expected", value, "but got", v)
			}
		}
		if len(keys) == 0 {
			t.Errorf("no keys")
		}
	}
}

func TestExtractData(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var errNoData = errors.New("no data found")

	tests := []struct {
		input  string
		expect map[string]string
		err    error
	}{
		{
			input: `
Score = 122.441
Number of wrong answers = 0
Placement cost = 805758.046
Measurement cost = 10863.00
Measurement count = 9585
`,
			expect: map[string]string{
				"Score":                   "122.441",
				"Number of wrong answers": "0",
				"Placement cost":          "805758.046",
				"Measurement cost":        "10863.00",
				"Measurement count":       "9585",
			},
			err: nil,
		},
		{
			input:  `InvalidData = abc`,
			expect: nil,
			err:    errNoData,
		},
		{
			input:  `Score = 100`,
			expect: map[string]string{"Score": "100"},
			err:    nil,
		},
		{
			input:  `Score = 100.0`,
			expect: map[string]string{"Score": "100.0"},
			err:    nil,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("TestExtractData_%d", i), func(t *testing.T) {
			result, err := extractData(test.input)

			if (err == nil) != (test.err == nil) {
				t.Errorf("Expected error %v, but got %v", test.err, err)
			}
			if err != nil && err.Error() != test.err.Error() {
				t.Errorf("Expected error %v, but got %v", test.err, err)
			}
			if (err == nil && test.err == nil) && !reflect.DeepEqual(result, test.expect) {
				t.Errorf("Expected %v, but got %v", test.expect, result)
			}
		})
	}
}
