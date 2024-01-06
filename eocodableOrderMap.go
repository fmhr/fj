package fj

import (
	"encoding/json"
	"fmt"

	"github.com/elliotchance/orderedmap/v2"
)

type EncodableOrderedMap orderedmap.OrderedMap[string, any]

func (m *EncodableOrderedMap) MarshalJSON() ([]byte, error) {
	items := map[string]interface{}{}
	for el := (*orderedmap.OrderedMap[string, any])(m).Front(); el != nil; el = el.Next() {
		items[fmt.Sprintf("%v", el.Key)] = el.Value
	}

	return json.Marshal(items)
}

func (m *EncodableOrderedMap) UnmarshalJSON(data []byte) error {
	var items map[string]interface{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	self := (*orderedmap.OrderedMap[string, any])(m)
	for k, v := range items {
		self.Set(k, v)
	}

	return nil
}
