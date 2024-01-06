package fj

import (
	"testing"

	"github.com/elliotchance/orderedmap/v2"
)

func TestSortBySeed(t *testing.T) {
	data := make([]*orderedmap.OrderedMap[string, any], 0)
	data[0].Set("seed", 3)
	data[0].Set("value", 3.0)
	data[1].Set("seed", 1)
	data[1].Set("value", 1.123)
	data[2].Set("seed", 2)
	data[2].Set("value", 2.234)
	sortBySeed(&data)

	if v, _ := data[0].Get("seed"); v != 1 {
		t.Errorf("Expected 1 for data[0].seed, but got %d", v)
	}
	if v, _ := data[1].Get("seed"); v != 2 {
		t.Errorf("Expected 2 for data[1].seed, but got %d", v)
	}
	if v, _ := data[2].Get("seed"); v != 3 {
		t.Errorf("Expected 3 for data[2].seed, but got %d", v)
	}

}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1.123, "1.123"},
		{3.0, "3"},
	}

	for _, test := range tests {
		result := formatFloat(test.input)
		if result != test.expected {
			t.Errorf("Expected %s for input %f, but got %s", test.expected, test.input, result)
		}
	}
}
