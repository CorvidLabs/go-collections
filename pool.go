package collections

import "sync"

// Pool is a thread-safe object pool with a configurable factory function.
// Unlike sync.Pool, objects are not subject to garbage collection and the
// pool tracks exact size.
type Pool[T any] struct {
	mu      sync.Mutex
	items   []T
	factory func() T
	maxSize int
}

// NewPool returns a new Pool that creates objects using the given factory.
// maxSize limits the number of idle objects stored; zero means unlimited.
func NewPool[T any](factory func() T, maxSize int) *Pool[T] {
	if maxSize < 0 {
		maxSize = 0
	}
	return &Pool[T]{
		factory: factory,
		maxSize: maxSize,
	}
}

// Get returns an object from the pool, or creates one using the factory
// if the pool is empty.
func (p *Pool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.items) == 0 {
		return p.factory()
	}
	last := len(p.items) - 1
	item := p.items[last]
	p.items = p.items[:last]
	return item
}

// Put returns an object to the pool. If the pool is at capacity, the object
// is discarded. Returns true if the object was stored.
func (p *Pool[T]) Put(item T) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.maxSize > 0 && len(p.items) >= p.maxSize {
		return false
	}
	p.items = append(p.items, item)
	return true
}

// Len returns the number of idle objects in the pool.
func (p *Pool[T]) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.items)
}

// Clear discards all idle objects from the pool.
func (p *Pool[T]) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.items = nil
}

// Prefill adds n objects to the pool using the factory, up to maxSize.
func (p *Pool[T]) Prefill(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := 0; i < n; i++ {
		if p.maxSize > 0 && len(p.items) >= p.maxSize {
			break
		}
		p.items = append(p.items, p.factory())
	}
}
