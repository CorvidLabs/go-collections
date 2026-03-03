package collections

import "sync"

// Set is a thread-safe hash set backed by a map.
type Set[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

// NewSet returns a new empty Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// NewSetFrom returns a new Set initialized with the given values.
func NewSetFrom[T comparable](values ...T) *Set[T] {
	s := &Set[T]{items: make(map[T]struct{}, len(values))}
	for _, v := range values {
		s.items[v] = struct{}{}
	}
	return s
}

// Add inserts a value into the set. Returns true if the value was new.
func (s *Set[T]) Add(value T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[value]; exists {
		return false
	}
	s.items[value] = struct{}{}
	return true
}

// Remove deletes a value from the set. Returns true if the value was present.
func (s *Set[T]) Remove(value T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[value]; !exists {
		return false
	}
	delete(s.items, value)
	return true
}

// Contains returns true if the set contains the value.
func (s *Set[T]) Contains(value T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[value]
	return exists
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// IsEmpty returns true if the set has no elements.
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Values returns all elements as a slice. Order is not guaranteed.
func (s *Set[T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, 0, len(s.items))
	for v := range s.items {
		result = append(result, v)
	}
	return result
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
}

// Union returns a new set containing all elements from both sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()

	result := &Set[T]{items: make(map[T]struct{}, len(s.items)+len(other.items))}
	for v := range s.items {
		result.items[v] = struct{}{}
	}
	for v := range other.items {
		result.items[v] = struct{}{}
	}
	return result
}

// Intersection returns a new set containing elements present in both sets.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()

	// Iterate over the smaller set for efficiency.
	small, big := s.items, other.items
	if len(small) > len(big) {
		small, big = big, small
	}

	result := &Set[T]{items: make(map[T]struct{})}
	for v := range small {
		if _, ok := big[v]; ok {
			result.items[v] = struct{}{}
		}
	}
	return result
}

// Difference returns a new set containing elements in s that are not in other.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()

	result := &Set[T]{items: make(map[T]struct{})}
	for v := range s.items {
		if _, ok := other.items[v]; !ok {
			result.items[v] = struct{}{}
		}
	}
	return result
}

// Range calls f for each element. If f returns false, iteration stops.
func (s *Set[T]) Range(f func(value T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for v := range s.items {
		if !f(v) {
			return
		}
	}
}
