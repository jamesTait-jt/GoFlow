package store

type InMemoryKVStore[K comparable, V any] struct {
	data map[K]V
}

func (kv *InMemoryKVStore[K, V]) Put(k K, v V) {
	kv.data[k] = v
}

func (kv *InMemoryKVStore[K, V]) Get(k K) (V, bool) {
	v, ok := kv.data[k]

	return v, ok
}
