package lru

import "testing"

type StringValue string

func (s StringValue) Len() int {
	return len(s)
}

func TestLruAdd(t *testing.T) {
	lru := NewLru(100)
	lru.Add("key1", StringValue("value1"))
	lru.Add("key2", StringValue("value2"))
	lru.Add("key2", StringValue("value2-updated"))
	if lru.Len() != 2 {
		t.Errorf("expected length 2, got %d", lru.Len())
	}
	if _, value := lru.Get("key1"); value == nil || value != StringValue("value1") {
		t.Error("expected value1 for key1")
	}
}

func TestLruDelete(t *testing.T) {
	lru := NewLru(100)
	lru.Add("key1", StringValue("value1"))
	lru.Add("key2", StringValue("value2"))
	lru.Delete("key1")

	if lru.Len() != 1 {
		t.Errorf("expected length 1 after deletion, got %d", lru.Len())
	}
}
