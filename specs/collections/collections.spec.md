---
module: collections
version: 1
status: active
files:
  - doc.go
  - lru.go
  - pool.go
  - priorityqueue.go
  - queue.go
  - ringbuffer.go
  - set.go
  - syncmap.go

db_tables: []
depends_on: []
---

# Collections

## Purpose

Provides generic, dependency-free, concurrency-safe collection primitives for
Go programs: sets, queues, ring buffers, priority queues, LRU caches, reusable
object pools, and a typed wrapper around `sync.Map`.

## Public API

### Exported Functions

| Export | Description |
|--------|-------------|
| `NewLRU` | Create a fixed-capacity LRU cache |
| `NewPool` | Create a bounded reusable-value pool |
| `NewPriorityQueue` | Create a comparator-ordered priority queue |
| `NewQueue` | Create an unbounded FIFO queue |
| `NewBoundedQueue` | Create a bounded FIFO queue |
| `NewRingBuffer` | Create a fixed-capacity evicting FIFO |
| `NewSet` | Create an empty set |
| `NewSetFrom` | Create a set from values |
| `NewSyncMap` | Create a typed concurrent map |
| `LRU` | Least-recently-used cache type |
| `Get` | Retrieve and promote an LRU entry |
| `Put` | Insert a retained value |
| `Peek` | Inspect without removal |
| `Remove` | Remove an entry |
| `Contains` | Test membership |
| `Len` | Report current size |
| `Cap` | Report capacity |
| `Keys` | Return keys |
| `Clear` | Remove retained values |
| `Pool` | Reusable-value pool type |
| `Prefill` | Populate a pool from its factory |
| `PriorityQueue` | Comparator-ordered queue type |
| `Push` | Insert a value |
| `Pop` | Remove the next value |
| `IsEmpty` | Report whether empty |
| `Drain` | Remove and return all values |
| `PushAll` | Insert multiple priority values |
| `Less` | Delegate internal heap comparison |
| `Swap` | Swap internal heap positions |
| `Queue` | FIFO queue type |
| `IsFull` | Report whether capacity is reached |
| `RingBuffer` | Fixed-capacity evicting FIFO type |
| `PeekNewest` | Inspect the newest ring value |
| `ToSlice` | Return ring values in FIFO order |
| `Set` | Hash-set type |
| `Add` | Insert a set value |
| `Values` | Return values |
| `Union` | Return set union |
| `Intersection` | Return set intersection |
| `Difference` | Return set difference |
| `Range` | Iterate until stopped |
| `SyncMap` | Typed concurrent map type |
| `Store` | Store a typed map value |
| `Load` | Load a typed map value |
| `LoadOrStore` | Atomically load or insert |
| `LoadAndDelete` | Atomically load and remove |
| `Delete` | Delete a map entry |

### Structs & Enums

| Type | Description |
|------|-------------|
| `Set[T]` | Thread-safe hash set with algebra and early-stop iteration |
| `Queue[T]` | Thread-safe FIFO queue with optional capacity |
| `RingBuffer[T]` | Fixed-capacity FIFO that reports evictions |
| `PriorityQueue[T]` | Heap-backed queue ordered by a caller-supplied comparator |
| `LRU[K,V]` | Fixed-capacity least-recently-used cache |
| `Pool[T]` | Bounded reusable-value pool backed by a factory |
| `SyncMap[K,V]` | Typed generic facade over `sync.Map` |

### Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `Get` | Type-specific | Retrieve and promote an LRU entry |
| `Put` | Type-specific | Insert a cache or pool value |
| `Peek` | Type-specific | Inspect without removal or promotion |
| `Remove` | Type-specific | Remove an entry |
| `Contains` | Type-specific | Test membership |
| `Len` | Type-specific | Report current size |
| `Cap` | Type-specific | Report capacity |
| `Keys` | Type-specific | Return keys |
| `Clear` | Type-specific | Remove retained values |
| `Prefill` | Type-specific | Populate a pool from its factory |
| `Push` | Type-specific | Insert a queue, heap, or ring value |
| `Pop` | Type-specific | Remove the next ordered value |
| `IsEmpty` | Type-specific | Report whether empty |
| `Drain` | Type-specific | Remove and return all ordered values |
| `PushAll` | Type-specific | Insert multiple priority values |
| `Less` | Internal heap adapter | Delegate comparator ordering |
| `Swap` | Internal heap adapter | Swap heap positions |
| `IsFull` | Type-specific | Report whether capacity is reached |
| `PeekNewest` | Type-specific | Inspect newest ring-buffer value |
| `ToSlice` | Type-specific | Return ring values oldest to newest |
| `Add` | Type-specific | Insert a set value |
| `Values` | Type-specific | Return values |
| `Union` | Type-specific | Return set union |
| `Intersection` | Type-specific | Return set intersection |
| `Difference` | Type-specific | Return set difference |
| `Range` | Type-specific | Iterate until callback stops |
| `Store` | Type-specific | Store a typed map value |
| `Load` | Type-specific | Load a typed map value |
| `LoadOrStore` | Type-specific | Atomically load or insert |
| `LoadAndDelete` | Type-specific | Atomically load and remove |
| `Delete` | Type-specific | Delete a map entry |

## Invariants

1. Public collection operations are safe for concurrent callers.
2. Empty lookup/pop operations return the element type's zero value with `false`.
3. Non-positive capacities normalize to a usable minimum where the type requires capacity.
4. Bounded queues reject pushes at capacity; ring buffers accept pushes and report the oldest evicted value.
5. LRU reads promote entries and insertion beyond capacity evicts the least recently used entry.
6. Priority queues always pop according to the caller's comparator.
7. Set algebra returns new sets without mutating either operand.
8. Pool size never exceeds its configured maximum and factory-created values are returned when empty.

## Behavioral Examples

```
Given a ring buffer at capacity
When one more value is pushed
Then the oldest value is returned as evicted and remaining values preserve FIFO order
```

## Error Cases

| Error | When | Behavior |
|-------|------|----------|
| Empty collection | A lookup or removal has no value | Return zero value and false |
| Full bounded queue or pool | A caller adds beyond capacity | Reject the value without exceeding capacity |

## Dependencies

- Go standard library only

## Change Log

| Version | Date | Changes |
|---------|------|---------|
| 1 | 2026-07-12 | Initial active specification of existing collection behavior |
