package collections

import "sync"

// RingBuffer is a fixed-size circular buffer. When full, new writes
// overwrite the oldest entries.
type RingBuffer[T any] struct {
	mu    sync.Mutex
	items []T
	head  int // read position (oldest item)
	tail  int // write position (next slot)
	count int
	cap   int
}

// NewRingBuffer returns a new RingBuffer with the given capacity.
// Capacity must be at least 1.
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	if capacity <= 0 {
		capacity = 1
	}
	return &RingBuffer[T]{
		items: make([]T, capacity),
		cap:   capacity,
	}
}

// Push adds an item to the buffer. If full, the oldest item is overwritten
// and returned along with true. Otherwise the zero value and false are returned.
func (rb *RingBuffer[T]) Push(item T) (evicted T, wasEvicted bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.count == rb.cap {
		evicted = rb.items[rb.head]
		wasEvicted = true
		rb.head = (rb.head + 1) % rb.cap
		rb.count--
	}

	rb.items[rb.tail] = item
	rb.tail = (rb.tail + 1) % rb.cap
	rb.count++
	return evicted, wasEvicted
}

// Pop removes and returns the oldest item. Returns the zero value and false
// if the buffer is empty.
func (rb *RingBuffer[T]) Pop() (T, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.count == 0 {
		var zero T
		return zero, false
	}

	item := rb.items[rb.head]
	var zero T
	rb.items[rb.head] = zero // clear reference for GC
	rb.head = (rb.head + 1) % rb.cap
	rb.count--
	return item, true
}

// Peek returns the oldest item without removing it.
func (rb *RingBuffer[T]) Peek() (T, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.count == 0 {
		var zero T
		return zero, false
	}
	return rb.items[rb.head], true
}

// PeekNewest returns the most recently added item without removing it.
func (rb *RingBuffer[T]) PeekNewest() (T, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.count == 0 {
		var zero T
		return zero, false
	}
	idx := (rb.tail - 1 + rb.cap) % rb.cap
	return rb.items[idx], true
}

// Len returns the number of items currently in the buffer.
func (rb *RingBuffer[T]) Len() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.count
}

// Cap returns the capacity of the buffer.
func (rb *RingBuffer[T]) Cap() int {
	return rb.cap
}

// IsFull returns true when the buffer is at capacity.
func (rb *RingBuffer[T]) IsFull() bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.count == rb.cap
}

// IsEmpty returns true when the buffer has no items.
func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.Len() == 0
}

// Clear removes all items from the buffer.
func (rb *RingBuffer[T]) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	var zero T
	for i := range rb.items {
		rb.items[i] = zero
	}
	rb.head = 0
	rb.tail = 0
	rb.count = 0
}

// ToSlice returns all items in order from oldest to newest.
func (rb *RingBuffer[T]) ToSlice() []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	result := make([]T, rb.count)
	for i := 0; i < rb.count; i++ {
		result[i] = rb.items[(rb.head+i)%rb.cap]
	}
	return result
}
