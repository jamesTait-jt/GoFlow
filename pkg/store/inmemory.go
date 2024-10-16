package store

import "sync"

type InMemoryKVStore[K comparable, V any] struct {
	data map[K]V
	mu   sync.Mutex
}

func NewInMemoryKVStore[K comparable, V any]() *InMemoryKVStore[K, V] {
	return &InMemoryKVStore[K, V]{
		data: make(map[K]V),
	}
}

func (kv *InMemoryKVStore[K, V]) Put(k K, v V) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[k] = v
}

func (kv *InMemoryKVStore[K, V]) Get(k K) (V, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	v, ok := kv.data[k]

	return v, ok
}
