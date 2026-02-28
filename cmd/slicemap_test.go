package cmd

import (
	"testing"
)

func TestNewSliceMap(t *testing.T) {
	sm := NewSliceMap()
	if len(sm) != 0 {
		t.Errorf("NewSliceMap() should return empty SliceMap, got len=%d", len(sm))
	}
}

func TestSliceMapSet(t *testing.T) {
	sm := NewSliceMap()
	sm.Set("key1", "value1")

	if len(sm) != 1 {
		t.Fatalf("expected len=1, got len=%d", len(sm))
	}
	if sm[0].Key != "key1" || sm[0].Val != "value1" {
		t.Errorf("expected {key1, value1}, got {%s, %s}", sm[0].Key, sm[0].Val)
	}
}

func TestSliceMapSetOverwrite(t *testing.T) {
	sm := NewSliceMap()
	sm.Set("key1", "value1")
	sm.Set("key1", "value2")

	if len(sm) != 1 {
		t.Fatalf("Set with same key should overwrite, not append: expected len=1, got len=%d", len(sm))
	}
	val, ok := sm.Get("key1")
	if !ok || val != "value2" {
		t.Errorf("expected value2, got %s", val)
	}
}

func TestSliceMapGet(t *testing.T) {
	sm := NewSliceMap()
	sm.Set("a", "1")
	sm.Set("b", "2")

	val, ok := sm.Get("a")
	if !ok || val != "1" {
		t.Errorf("Get(a): expected (1, true), got (%s, %v)", val, ok)
	}

	val, ok = sm.Get("b")
	if !ok || val != "2" {
		t.Errorf("Get(b): expected (2, true), got (%s, %v)", val, ok)
	}

	val, ok = sm.Get("nonexistent")
	if ok || val != "" {
		t.Errorf("Get(nonexistent): expected ('', false), got (%s, %v)", val, ok)
	}
}

func TestSliceMapOrder(t *testing.T) {
	sm := NewSliceMap()
	keys := []string{"c", "a", "b"}
	for _, k := range keys {
		sm.Set(k, k+"_val")
	}

	for i, k := range keys {
		if sm[i].Key != k {
			t.Errorf("order not preserved: index %d expected key=%s, got key=%s", i, k, sm[i].Key)
		}
	}
}

func TestSliceMapOverwritePreservesOrder(t *testing.T) {
	sm := NewSliceMap()
	sm.Set("a", "1")
	sm.Set("b", "2")
	sm.Set("c", "3")
	sm.Set("a", "updated")

	expected := []string{"a", "b", "c"}
	for i, k := range expected {
		if sm[i].Key != k {
			t.Errorf("overwrite changed order: index %d expected key=%s, got key=%s", i, k, sm[i].Key)
		}
	}

	val, _ := sm.Get("a")
	if val != "updated" {
		t.Errorf("expected updated, got %s", val)
	}
}
