package lru

import (
	"testing"
)

func TestNewCache(t *testing.T) {
	capacity := 10
	cache := NewCache[int, int](capacity)

	if cache.capacity != capacity {
		t.Errorf("Expected capacity %d, but got %d", capacity, cache.capacity)
	}
	if cache.list.Len() != 0 {
		t.Errorf("Expected empty list, but got length %d", cache.list.Len())
	}
	if len(cache.store) != 0 {
		t.Errorf("Expected empty store, but got length %d", len(cache.store))
	}
}

func TestGetAndPut(t *testing.T) {
	cache := NewCache[int, int](2)

	// Test Put
	cache.Put(1, 100)
	cache.Put(2, 200)

	// Test Get existing keys
	if val := cache.Get(1); val != 100 {
		t.Errorf("Expected 100, but got %d", val)
	}
	if val := cache.Get(2); val != 200 {
		t.Errorf("Expected 200, but got %d", val)
	}

	// Test Get non-existent key
	if val := cache.Get(3); val != 0 {
		t.Errorf("Expected 0, but got %d", val)
	}

	// Test LRU eviction
	cache.Put(3, 300)
	if val := cache.Get(1); val != 0 {
		t.Errorf("Expected 1 to be evicted, but got %d", val)
	}
}
