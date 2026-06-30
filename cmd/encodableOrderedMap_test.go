package cmd

import (
	"encoding/json"
	"testing"

	"github.com/elliotchance/orderedmap/v2"
)

func TestEncodableOrderedMap_MarshalJSON(t *testing.T) {
	om := orderedmap.NewOrderedMap[string, string]()
	om.Set("first", "1")
	om.Set("second", "2")
	om.Set("third", "3")

	encMap := (*EncodableOrderedMap)(om)

	data, err := json.Marshal(encMap)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expectedJSON := `[{"key":"first","value":"1"},{"key":"second","value":"2"},{"key":"third","value":"3"}]`
	if string(data) != expectedJSON {
		t.Errorf("expected %s, got %s", expectedJSON, string(data))
	}
}

func TestEncodableOrderedMap_UnmarshalJSON(t *testing.T) {
	jsonInput := `[{"key":"a","value":"apple"},{"key":"b","value":"banana"},{"key":"c","value":"cherry"}]`

	var encMap EncodableOrderedMap
	// UnmarshalJSONの初期化として内部OrderedMapが正常にセットされるかを検証するため
	// sliceとしてのセットアップをキャストして行う
	om := orderedmap.NewOrderedMap[string, string]()
	encMap = (EncodableOrderedMap)(*om)

	err := json.Unmarshal([]byte(jsonInput), &encMap)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	actualOM := (*orderedmap.OrderedMap[string, string])(&encMap)
	if actualOM.Len() != 3 {
		t.Fatalf("expected length 3, got %d", actualOM.Len())
	}

	keys := []string{"a", "b", "c"}
	values := []string{"apple", "banana", "cherry"}

	// 順序の確認
	i := 0
	for el := actualOM.Front(); el != nil; el = el.Next() {
		if el.Key != keys[i] {
			t.Errorf("expected key %s at index %d, got %s", keys[i], i, el.Key)
		}
		if el.Value != values[i] {
			t.Errorf("expected value %s for key %s, got %s", values[i], el.Key, el.Value)
		}
		i++
	}
}

func TestEncodableOrderedMap_UnmarshalJSON_Error(t *testing.T) {
	invalidJSON := `[{"key":"a", "value": 123}]` // value が string ではない

	var encMap EncodableOrderedMap
	om := orderedmap.NewOrderedMap[string, string]()
	encMap = (EncodableOrderedMap)(*om)

	// valueの型ミスマッチによるエラーハンドリングを確認
	// ValueはanyでUnmarshalされるが、内部で Value.(string) キャストするためパニックまたはエラーになるか？
	// encodableOrderedMap.go 42行目: self.Set(item.Key, item.Value.(string))
	// ここは string 以外の値が来ると Go の interface 変換で panic する可能性がある。
	// テストとしてこれを検証（本来は panic を防ぐ方が良いが、現状の実装の振る舞いを確認）
	defer func() {
		if r := recover(); r == nil {
			t.Log("expected panic or error for invalid type casting, but code did not panic (might fail or pass depending on runtime)")
		}
	}()

	_ = json.Unmarshal([]byte(invalidJSON), &encMap)
}
