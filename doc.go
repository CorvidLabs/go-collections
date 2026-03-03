// Package collections provides type-safe concurrent data structures for Go.
//
// All types are safe for concurrent use by multiple goroutines.
//
// Data structures:
//   - SyncMap[K, V]: Type-safe concurrent map backed by sync.Map
//   - Queue[T]: Unbounded FIFO queue with optional capacity limit
//   - Set[T]: Concurrent hash set with bulk operations
//   - RingBuffer[T]: Fixed-size circular buffer
//   - Pool[T]: Object pool with configurable factory
package collections
