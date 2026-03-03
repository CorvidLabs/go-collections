package collections

import "sync"

// Queue is an unbounded, thread-safe FIFO queue.
// An optional maximum capacity can be set; zero means unlimited.
type Queue[T any] struct {
	mu      sync.Mutex
	items   []T
	maxSize int
}

// NewQueue returns a new Queue with no capacity limit.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// NewBoundedQueue returns a new Queue that rejects pushes when full.
// maxSize must be positive.
func NewBoundedQueue[T any](maxSize int) *Queue[T] {
	if maxSize <= 0 {
		maxSize = 1
	}
	return &Queue[T]{maxSize: maxSize}
}

// Push appends an item to the back of the queue.
// Returns false if the queue is at capacity.
func (q *Queue[T]) Push(item T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.maxSize > 0 && len(q.items) >= q.maxSize {
		return false
	}
	q.items = append(q.items, item)
	return true
}

// Pop removes and returns the item at the front of the queue.
// Returns the zero value and false if the queue is empty.
func (q *Queue[T]) Pop() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek returns the item at the front without removing it.
func (q *Queue[T]) Peek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// IsEmpty returns true if the queue has no items.
func (q *Queue[T]) IsEmpty() bool {
	return q.Len() == 0
}

// IsFull returns true if the queue is at capacity (bounded queues only).
// Always returns false for unbounded queues.
func (q *Queue[T]) IsFull() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.maxSize > 0 && len(q.items) >= q.maxSize
}

// Clear removes all items from the queue.
func (q *Queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = nil
}

// Drain removes and returns all items from the queue.
func (q *Queue[T]) Drain() []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	items := q.items
	q.items = nil
	return items
}
