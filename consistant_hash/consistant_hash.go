package consistanthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type ConsistantHash struct {
	replicas int
	keys     []int
	hashMap  map[int]string
	hashFunc func(data []byte) int
}

func NewConsistantHash(replicas int, hashFunc func(data []byte) int) *ConsistantHash {
	hash := &ConsistantHash{
		replicas: replicas,
		hashMap:  make(map[int]string),
		keys:     []int{},
	}
	if hashFunc == nil {
		hash.hashFunc = func(data []byte) int {
			return int(crc32.ChecksumIEEE(data))
		}
	} else {
		hash.hashFunc = hashFunc
	}
	return hash
}

// 初始化一致性hash环
func (hash *ConsistantHash) initCircle(nodes ...string) {
	if len(nodes) == 0 {
		return
	}
	if len(hash.keys) != 0 {
		hash.Clear()
	}

	for _, node := range nodes {
		for i := 0; i < hash.replicas; i++ {
			hashKey := hash.hashFunc([]byte(strconv.Itoa(i) + node))
			hash.keys = append(hash.keys, hashKey)
			hash.hashMap[hashKey] = node
		}
	}

	sort.Ints(hash.keys)
}

func (hash *ConsistantHash) Clear() {
	hash.keys = []int{}
	hash.hashMap = make(map[int]string)
}

func (hash *ConsistantHash) GetNode(key string) (string, bool) {
	if len(hash.keys) == 0 {
		return "", false
	}

	hashKey := hash.hashFunc([]byte(key))
	// 二分查找
	idx := sort.Search(len(hash.keys), func(idx int) bool {
		return hash.keys[idx] >= hashKey
	})

	if idx == len(hash.keys) {
		idx = 0
	}
	return hash.hashMap[hash.keys[idx]], true
}
