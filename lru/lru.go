package lru

import (
	"container/list"
	"sync"
)

type Lru struct {
	maxBytes int
	curBytes int
	cache    *list.List
	keyMap   map[string]*list.Element
	mtx      sync.Mutex
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// 访问更新
// 怎么删除
// 添加
// 清空

func NewLru(maxBytes int) *Lru {
	return &Lru{
		maxBytes: maxBytes,
		cache:    &list.List{},
		keyMap:   make(map[string]*list.Element),
		mtx:      sync.Mutex{},
	}
}

func (lru *Lru) Get(key string) (string, Value) {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()

	if value := lru.keyMap[key]; value != nil {
		lru.cache.MoveToBack(value)
		key := value.Value.(*entry).key
		value := value.Value.(*entry).value
		return key, value
	}
	return "", nil
}

func (lru *Lru) Add(key string, value Value) {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()
	if e := lru.keyMap[key]; e != nil {
		lru.cache.MoveToBack(e)
	} else {
		entry := entry{key: key, value: value}
		element := &list.Element{
			Value: entry,
		}
		if !lru.couldInserEle(element) {
			lru.shrink(value.Len())
		}
		lru.cache.PushBack(entry)
	}
}

func (lru *Lru) Delete(key string) {
	lru.mtx.Lock()
	defer lru.mtx.Unlock()
	if element, exist := lru.has(key); exist {
		lru.delete(element)
	}
}

func (lru *Lru) shrink(size int) {
	if size > lru.bytes() {
		lru.clear()
		return
	}

	for size > 0 {
		size -= lru.cache.Front().Value.(*entry).value.Len()
		lru.delete(lru.cache.Front())
	}
}

func (lru *Lru) delete(e *list.Element) {
	if lru.cache.Len() == 0 {
		return
	}
	targetKey := e.Value.(*entry).key
	if value := lru.keyMap[targetKey]; value != nil {
		lru.cache.Remove(e)
		delete(lru.keyMap, targetKey)
		lru.curBytes -= value.Value.(*entry).value.Len()
	}
}

func (lru *Lru) clear() {
	lru.cache = &list.List{}
	lru.keyMap = make(map[string]*list.Element)
	lru.curBytes = 0
}

func (lru *Lru) couldInserEle(e *list.Element) bool {
	if lru.bytes()+e.Value.(*entry).value.Len() > lru.maxBytes {
		return false
	}
	return true
}

func (lru *Lru) bytes() int {
	return lru.curBytes
}

func (lru *Lru) has(key string) (*list.Element, bool) {
	if value := lru.keyMap[key]; value != nil {
		return value, true
	}
	return nil, false
}
