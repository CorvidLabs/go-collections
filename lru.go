package collections

import (
	"container/list"
	"sync"
)

// LRU is a thread-safe, generic Least Recently Used cache with a fixed capacity.
// When the cache is full and a new key is inserted, the least recently accessed
// entry is evicted.
type LRU[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[K]*list.Element
	order    *list.List // front = most recent, back = least recent
}

type lruEntry[K comparable, V any] struct {
	key   K
	value V
}

// NewLRU returns a new LRU cache with the given capacity.
// Capacity must be at least 1.
func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	if capacity <= 0 {
		capacity = 1
	}
	return &LRU[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element, capacity),
		order:    list.New(),
	}
}

// Get retrieves a value by key and marks it as recently used.
// Returns the zero value and false if the key is not present.
func (c *LRU[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.order.MoveToFront(elem)
	return elem.Value.(*lruEntry[K, V]).value, true
}

// Put inserts or updates a key-value pair. If the cache is at capacity,
// the least recently used entry is evicted and returned along with true.
func (c *LRU[K, V]) Put(key K, value V) (evictedKey K, evictedValue V, evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing key.
	if elem, ok := c.items[key]; ok {
		entry := elem.Value.(*lruEntry[K, V])
		entry.value = value
		c.order.MoveToFront(elem)
		return evictedKey, evictedValue, false
	}

	// Evict if at capacity.
	if c.order.Len() >= c.capacity {
		back := c.order.Back()
		entry := back.Value.(*lruEntry[K, V])
		evictedKey = entry.key
		evictedValue = entry.value
		evicted = true
		delete(c.items, entry.key)
		c.order.Remove(back)
	}

	// Insert new entry.
	entry := &lruEntry[K, V]{key: key, value: value}
	elem := c.order.PushFront(entry)
	c.items[key] = elem
	return evictedKey, evictedValue, evicted
}

// Peek retrieves a value by key without updating its recency.
func (c *LRU[K, V]) Peek(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	return elem.Value.(*lruEntry[K, V]).value, true
}

// Remove deletes a key from the cache. Returns the value and true if the key
// was present.
func (c *LRU[K, V]) Remove(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	entry := elem.Value.(*lruEntry[K, V])
	delete(c.items, key)
	c.order.Remove(elem)
	return entry.value, true
}

// Contains returns true if the key is in the cache. Does not update recency.
func (c *LRU[K, V]) Contains(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.items[key]
	return ok
}

// Len returns the number of entries in the cache.
func (c *LRU[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

// Cap returns the cache capacity.
func (c *LRU[K, V]) Cap() int {
	return c.capacity
}

// Keys returns all keys in order from most to least recently used.
func (c *LRU[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]K, 0, c.order.Len())
	for e := c.order.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*lruEntry[K, V]).key)
	}
	return keys
}

// Clear removes all entries from the cache.
func (c *LRU[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[K]*list.Element, c.capacity)
	c.order.Init()
}
