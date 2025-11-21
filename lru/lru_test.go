package lru

import (
	"testing"

	"github.com/TheChosenGay/memcache/byte_view"
)

type StringValue string

func (s StringValue) Len() int {
	return len(s)
}

func TestLruAdd(t *testing.T) {
	lru := NewLru(100)
	lru.Add("key1", byte_view.NewByteView([]byte("hello")))
	lru.Add("key2", byte_view.NewByteView([]byte("value2")))
	lru.Add("key2", byte_view.NewByteView([]byte("value2-updated")))
	if lru.Len() != 2 {
		t.Errorf("expected length 2, got %d", lru.Len())
	}
	if _, exist := lru.Get("key1"); !exist {
		t.Error("expected value1 for key1")
	}
}

func TestLruDelete(t *testing.T) {
	lru := NewLru(100)
	lru.Add("key1", byte_view.NewByteView([]byte("value1")))
	lru.Add("key2", byte_view.NewByteView([]byte("value2")))
	lru.Delete("key1")

	if lru.Len() != 1 {
		t.Errorf("expected length 1 after deletion, got %d", lru.Len())
	}
}
