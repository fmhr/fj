package fj

import (
	"encoding/json"
	"log"

	"github.com/elliotchance/orderedmap/v2"
)

type EncodableOrderedMap orderedmap.OrderedMap[string, any]

type EncodableOrderedMapItem struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (m *EncodableOrderedMap) MarshalJSON() ([]byte, error) {
	items := make([]EncodableOrderedMapItem, (*orderedmap.OrderedMap[string, any])(m).Len())
	var i int
	for el := (*orderedmap.OrderedMap[string, any])(m).Front(); el != nil; el = el.Next() {
		items[i] = EncodableOrderedMapItem{el.Key, el.Value}
		i++
	}
	return json.Marshal(items)
}

func (m *EncodableOrderedMap) UnmarshalJson(data []byte) error {
	var items []EncodableOrderedMapItem
	err := json.Unmarshal(data, &items)
	if err != nil {
		log.Println(string(data))
		return err
	}

	self := (*orderedmap.OrderedMap[string, any])(m)
	for _, item := range items {
		log.Println(item.Key, item.Value)
		self.Set(item.Key, item.Value)
	}
	return nil
}
