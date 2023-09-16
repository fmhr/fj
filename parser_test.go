package fj

import (
	"reflect"
	"testing"
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
		result, err := ExtractKeyValuePairs(test.input)
		if (err != nil) != test.err {
			t.Errorf("Expected error %v, but got %v", test.err, err)
		}
		for k, v := range test.output {
			if result[k] != v {
				t.Errorf("For key %s, expected %f but got %f", k, v, result[k])
			}
		}
	}
}
func TestExtractData(t *testing.T) {
	tests := []struct {
		input  string
		expect map[string]float64
		err    bool
	}{
		{
			input: `
Score = 122.441
Number of wrong answers = 0
Placement cost = 805758.046
Measurement cost = 10863.00
Measurement count = 9585
`,
			expect: map[string]float64{
				"Score":                   122.441,
				"Number of wrong answers": 0,
				"Placement cost":          805758.046,
				"Measurement cost":        10863.00,
				"Measurement count":       9585,
			},
			err: false,
		},
		{
			input:  `InvalidData = abc`,
			expect: nil,
			err:    true,
		},
	}

	for _, test := range tests {
		result, err := extractData(test.input)
		if (err != nil) != test.err {
			t.Errorf("For input %q, expected error? %v but got: %v\n", test.input, test.err, err)
		}
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("For input %q, expected %v but got %v\n", test.input, test.expect, result)
		}
	}
}
