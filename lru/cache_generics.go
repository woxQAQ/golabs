package lru

import (
	"golabs/container/listx"
)

type keyValue[K comparable, V any] struct {
	key   K
	value V
}

type cache[K comparable, V any] struct {
	list     *listx.List[keyValue[K, V]]
	store    map[K]*listx.Element[keyValue[K, V]]
	capacity int
}

func NewCache[K comparable, V any](capacity int) cache[K, V] {
	return cache[K, V]{
		list:     listx.New[keyValue[K, V]](),
		capacity: capacity,
		store:    make(map[K]*listx.Element[keyValue[K, V]], capacity),
	}
}

func (c *cache[K, V]) Get(key K) V {
	if v, ok := c.store[key]; ok {
		c.list.MoveToFront(v)
		return v.Value.value
	}
	var v V
	return v
}

func (c *cache[K, V]) Put(key K, value V) {
	if v, ok := c.store[key]; ok {
		c.list.MoveToFront(v)
		v.Value.value = value
	} else {
		if c.list.Len() == c.capacity {
			last := c.list.Back()
			c.list.Remove(last)
			delete(c.store, last.Value.key)
		}
		elem := c.list.PushFront(keyValue[K, V]{key, value})
		c.store[key] = elem
	}
}
