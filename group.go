package group

import (
	"errors"
	"sync"

	"github.com/TheChosenGay/memcache/lru"
)

type Getter interface {
	Get(key string) ([]byte, error)
}
type Group struct {
	name   string
	cache  *lru.Lru
	getter Getter
}

var (
	gmtx   sync.Mutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int, getter Getter) *Group {
	gmtx.Lock()
	defer gmtx.Unlock()
	g := &Group{
		name:   name,
		cache:  lru.NewLru(cacheBytes),
		getter: getter,
	}
	groups[name] = g
	return g
}

func (g *Group) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, nil
	}

	if value, exist := g.cache.Get(key); exist {
		return value.ByteSlice(), nil
	}

	if g.getter == nil {
		return nil, errors.New("nil getter")
	}
	// 如果找不到，就交给 getter获取了
	return g.getter.Get(key)
}
