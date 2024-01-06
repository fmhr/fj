package fj

import (
	"encoding/json"

	"github.com/elliotchance/orderedmap/v2"
)

type EncodableOrderedMap orderedmap.OrderedMap[string, any]

type EncodableOrderedMapItem struct {
	key   string
	value any
}

func (m *EncodableOrderedMap) MarshalJSON() ([]byte, error) {
	items := make([]EncodableOrderedMapItem, (*orderedmap.OrderedMap[string, any])(m).Len())
	for i, key := range (*orderedmap.OrderedMap[string, any])(m).Keys() {
		value, _ := (*orderedmap.OrderedMap[string, any])(m).Get(key)
		items[i] = EncodableOrderedMapItem{key, value}
	}
	return json.Marshal(items)
}

func (m *EncodableOrderedMap) UnmarshalJson(data []byte) error {
	var items []EncodableOrderedMapItem
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	self := (*orderedmap.OrderedMap[string, any])(m)
	for _, item := range items {
		self.Set(item.key, item.value)
	}
	return nil
}
