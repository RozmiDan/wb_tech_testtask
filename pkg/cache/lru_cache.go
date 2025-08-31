package lru_cache

import (
	"sync"

	"iter"

	"github.com/RozmiDan/wb_tech_testtask/pkg/cache/linklist"
)

type LRU[K comparable, V any] interface {
	Put(key K, val V)
	Get(key K) V
	Size() int
	All() iter.Seq2[K, V]
}

type Node[K comparable, V any] struct {
	key   K
	value V
}

type LruCache[K comparable, V any] struct {
	list         *linklist.List[Node[K, V]]
	mp           map[K]*linklist.Node[Node[K, V]]
	defaultValue V
	capacity     int
	mu           sync.RWMutex
}

func NewLruCache[K comparable, V any](cap int, defVal V) *LruCache[K, V] {
	if cap <= 0 {
		panic("lru: capacity must be > 0")
	}
	return &LruCache[K, V]{
		list:         linklist.NewList[Node[K, V]](),
		mp:           make(map[K]*linklist.Node[Node[K, V]], cap),
		defaultValue: defVal,
		capacity:     cap,
	}
}

func (lru *LruCache[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		lru.mu.RLock()
		defer lru.mu.RUnlock()
		for it := range lru.list.All() {
			if !yield(it.key, it.value) {
				return
			}
		}
	}
}

func (lru *LruCache[K, V]) Put(key K, val V) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if n, ok := lru.mp[key]; ok {
		lru.list.PutNewValue(n, Node[K, V]{key: key, value: val})
		lru.list.MoveToFront(n)

		return
	}

	if lru.list.Size() >= lru.capacity {
		tail := lru.list.Back()
		if tail != nil {
			evicted := lru.list.Remove(tail)
			delete(lru.mp, evicted.key)
		}
	}

	lru.list.PushFront(Node[K, V]{key: key, value: val})
	lru.mp[key] = lru.list.Front()
}

func (lru *LruCache[K, V]) Get(key K) V {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if n, ok := lru.mp[key]; ok {
		lru.list.MoveToFront(n)
		return n.GetData().value
	}
	return lru.defaultValue
}

func (lru *LruCache[K, V]) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	return lru.list.Size()
}
