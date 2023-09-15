package fj

import (
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

func TestExtractScore(t *testing.T) {
	tests := []struct {
		input  string
		output int
		err    bool
	}{
		{
			input:  "Score = 100",
			output: 100,
			err:    false,
		},
		{
			input:  "Some other string",
			output: 0,
			err:    true,
		},
		{
			input:  "Score=50",
			output: 50,
			err:    false,
		},
	}

	for _, test := range tests {
		result, err := extractScore(test.input)
		if (err != nil) != test.err {
			t.Errorf("Expected error %v, but got %v", test.err, err)
		}
		if result != test.output {
			t.Errorf("Expected score %d, but got %d", test.output, result)
		}
	}
}
