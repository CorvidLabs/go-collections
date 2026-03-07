package collections

import (
	"container/heap"
	"sync"
)

// PriorityQueue is a thread-safe, generic min-heap priority queue.
// Items with lower priority values are dequeued first.
type PriorityQueue[T any] struct {
	mu   sync.Mutex
	heap pqHeap[T]
	less func(a, b T) bool
}

// NewPriorityQueue returns a new PriorityQueue ordered by the given less function.
// less(a, b) should return true when a has higher priority (should be dequeued before b).
func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{less: less}
	pq.heap.less = less
	heap.Init(&pq.heap)
	return pq
}

// Push adds an item to the queue.
func (pq *PriorityQueue[T]) Push(item T) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(&pq.heap, item)
}

// Pop removes and returns the highest-priority item.
// Returns the zero value and false if the queue is empty.
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.heap.Len() == 0 {
		var zero T
		return zero, false
	}
	item := heap.Pop(&pq.heap).(T)
	return item, true
}

// Peek returns the highest-priority item without removing it.
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.heap.Len() == 0 {
		var zero T
		return zero, false
	}
	return pq.heap.items[0], true
}

// Len returns the number of items in the queue.
func (pq *PriorityQueue[T]) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.heap.Len()
}

// IsEmpty returns true if the queue has no items.
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.Len() == 0
}

// Clear removes all items from the queue.
func (pq *PriorityQueue[T]) Clear() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.heap.items = nil
}

// Drain removes and returns all items in priority order (highest priority first).
func (pq *PriorityQueue[T]) Drain() []T {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	result := make([]T, 0, pq.heap.Len())
	for pq.heap.Len() > 0 {
		result = append(result, heap.Pop(&pq.heap).(T))
	}
	return result
}

// PushAll adds multiple items to the queue.
func (pq *PriorityQueue[T]) PushAll(items ...T) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	for _, item := range items {
		heap.Push(&pq.heap, item)
	}
}

// internal heap implementation for container/heap interface
type pqHeap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func (h pqHeap[T]) Len() int           { return len(h.items) }
func (h pqHeap[T]) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h pqHeap[T]) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *pqHeap[T]) Push(x any) {
	h.items = append(h.items, x.(T))
}

func (h *pqHeap[T]) Pop() any {
	old := h.items
	n := len(old)
	item := old[n-1]
	var zero T
	old[n-1] = zero // clear reference for GC
	h.items = old[:n-1]
	return item
}
