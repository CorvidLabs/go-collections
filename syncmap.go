package collections

import "sync"

// SyncMap is a type-safe concurrent map. It wraps sync.Map with generics
// to provide compile-time type safety.
type SyncMap[K comparable, V any] struct {
	m sync.Map
}

// NewSyncMap returns a new empty SyncMap.
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// Store sets the value for a key.
func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Load returns the value stored for a key, or the zero value if absent.
// The ok result indicates whether the key was found.
func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	raw, ok := m.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return raw.(V), true
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	raw, loaded := m.m.LoadOrStore(key, value)
	return raw.(V), loaded
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	raw, loaded := m.m.LoadAndDelete(key)
	if !loaded {
		var zero V
		return zero, false
	}
	return raw.(V), true
}

// Delete removes the value for a key.
func (m *SyncMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Range calls f sequentially for each key-value pair in the map.
// If f returns false, Range stops the iteration.
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

// Len returns the number of entries in the map.
// Note: this requires a full scan and is O(n).
func (m *SyncMap[K, V]) Len() int {
	count := 0
	m.m.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

// Keys returns all keys in the map. Order is not guaranteed.
func (m *SyncMap[K, V]) Keys() []K {
	var keys []K
	m.m.Range(func(k, _ any) bool {
		keys = append(keys, k.(K))
		return true
	})
	return keys
}

// Values returns all values in the map. Order is not guaranteed.
func (m *SyncMap[K, V]) Values() []V {
	var values []V
	m.m.Range(func(_, v any) bool {
		values = append(values, v.(V))
		return true
	})
	return values
}
