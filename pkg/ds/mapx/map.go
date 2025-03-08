package mapx

import "reflect"

type IMap[K comparable, V any] interface {
	Get(key K) (V, bool)
	GetOrDefault(key K, defaultValue V) V
	Set(key K, value V)
	Delete(key K)
	IsEqual(other IMap[K, V]) bool
}

type Map[K comparable, V any] struct {
	data map[K]V
}

func NewMapFromSource[K comparable, V any](source map[K]V) *Map[K, V] {
	return &Map[K, V]{
		data: source,
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *Map[K, V]) GetOrDefault(key K, defaultValue V) V {
	v, ok := m.Get(key)
	if !ok {
		return defaultValue
	}
	return v
}

func (m *Map[K, V]) Set(key K, value V) {
	m.data[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.data, key)
}

func (m *Map[K, V]) IsEqual(other IMap[K, V]) bool {
	if len(m.data) != len(other.(*Map[K, V]).data) {
		return false
	}
	for k, v1 := range m.data {
		if v2, ok := other.Get(k); !ok || !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	return true
}
