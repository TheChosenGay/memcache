package consistanthash

import (
	"strconv"
	"testing"
)

func TestConsistantMapping(t *testing.T) {
	hashFunc := func(data []byte) int {
		key, _ := strconv.Atoi(string(data))
		return key
	}
	hash := NewConsistantHash(3, hashFunc)
	hash.AddKeys("2", "3", "4")

	tests := []struct {
		key    string
		expect int
	}{
		{"2", 2},
		{"3", 3},
		{"4", 4},
	}

	for _, tt := range tests {
		res, exist := hash.Get(tt.key)
		if !exist || res != tt.key {
			t.Errorf("GetNode(%q) = %q, want %q", tt.key, res, tt.key)
		}
	}

}

func TestConsistentRepeatable(t *testing.T) {
	hash := NewConsistantHash(3, nil)
	hash.AddKeys("NodeA", "NodeB", "NodeC")

	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	for _, k := range keys {
		first, ok1 := hash.Get(k)
		second, ok2 := hash.Get(k)
		if !ok1 || !ok2 {
			t.Fatalf("GetNode returned ok=false for key %q", k)
		}
		if first != second {
			t.Errorf("GetNode not repeatable for key %q: %q != %q", k, first, second)
		}
	}
}

func TestEmptyHashReturnsFalse(t *testing.T) {
	hash := NewConsistantHash(3, nil)
	if _, ok := hash.Get("anykey"); ok {
		t.Fatalf("expected GetNode to return ok=false when no nodes are added")
	}
}

func TestKeysCountMatchesReplicasAndNodes(t *testing.T) {
	replicas := 4
	nodes := []string{"A", "B", "C", "D"}
	hash := NewConsistantHash(replicas, nil)
	hash.AddKeys(nodes...)
	expected := replicas * len(nodes)
	if len(hash.keys) != expected {
		t.Fatalf("len(keys) = %d, want %d", len(hash.keys), expected)
	}
}

func TestDistributionAcrossNodes(t *testing.T) {
	hash := NewConsistantHash(5, nil)
	nodes := []string{"NodeA", "NodeB", "NodeC"}
	hash.AddKeys(nodes...)

	counts := make(map[string]int)
	const totalKeys = 200
	for i := 0; i < totalKeys; i++ {
		k := "key_" + strconv.Itoa(i)
		node, ok := hash.Get(k)
		if !ok {
			t.Fatalf("GetNode returned ok=false for key %q", k)
		}
		counts[node]++
	}

	// ensure each node got at least one key
	for _, n := range nodes {
		if counts[n] == 0 {
			t.Errorf("node %q received 0 keys, expected > 0", n)
		}
	}
}

func TestReplicaCountAffectsPlacement(t *testing.T) {
	nodes := []string{"NodeA", "NodeB", "NodeC"}
	hash1 := NewConsistantHash(1, nil)
	hash1.AddKeys(nodes...)
	hash5 := NewConsistantHash(5, nil)
	hash5.AddKeys(nodes...)

	// check mappings for a set of keys; expect at least one difference
	keys := []string{"k1", "k2", "k3", "k4", "k5", "k6", "k7"}
	differ := false
	for _, k := range keys {
		n1, _ := hash1.Get(k)
		n5, _ := hash5.Get(k)
		if n1 != n5 {
			differ = true
			break
		}
	}
	if !differ {
		t.Errorf("expected at least one key to map to different nodes when replicas differ")
	}
}
