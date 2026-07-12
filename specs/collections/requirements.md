---
spec: collections.spec.md
---

## User Stories

- As a Go developer, I want generic concurrent collections with predictable
  capacity, ordering, and empty-value behavior and no third-party dependencies.

## Acceptance Criteria

### REQ-collections-001

All exported collection operations SHALL remain safe under concurrent use.

### REQ-collections-002

Queue and ring-buffer operations SHALL preserve FIFO order while bounded queues
reject overflow and ring buffers report oldest-value eviction.

### REQ-collections-003

LRU operations SHALL promote accessed entries and evict the least recently used
entry when insertion exceeds capacity.

### REQ-collections-004

Priority queues SHALL order values by the supplied comparator and sets SHALL
provide non-mutating union, intersection, and difference.

### REQ-collections-005

Pools SHALL honor their maximum retained size and factory behavior, and
SyncMap SHALL preserve typed load/store/delete/range semantics.

### REQ-collections-006

Empty retrieval operations SHALL return the element type's zero value and a
false presence result.

## Constraints

- Ordering is unspecified for hash-backed set and map enumeration.

## Out of Scope

- Persistence, distributed coordination, and lock-free guarantees beyond
  `sync.Map`.
