## MODIFIED

### SPEC SECTION Invariants

1. Public collection operations are safe for concurrent callers.
2. Empty lookup/pop operations return the element type's zero value with `false`.
3. Non-positive capacities normalize to a usable minimum where the type requires capacity.
4. Bounded queues reject pushes at capacity; ring buffers accept pushes and report the oldest evicted value.
5. LRU reads promote entries and insertion beyond capacity evicts the least recently used entry.
6. Priority queues always pop according to the caller's comparator.
7. Set algebra returns new sets without mutating either operand.
8. Pool size never exceeds its configured maximum and factory-created values are returned when empty.

### REQUIREMENT REQ-collections-001

All exported collection operations SHALL remain safe under concurrent use.

Acceptance Criteria
- The complete package test suite passes under Go's race detector.

### REQUIREMENT REQ-collections-002

Queue and ring-buffer operations SHALL preserve FIFO order while bounded queues reject overflow and ring buffers report oldest-value eviction.

Acceptance Criteria
- Queue and ring-buffer tests cover FIFO order, empty reads, bounded overflow rejection, and oldest-value eviction.

### REQUIREMENT REQ-collections-003

LRU operations SHALL promote accessed entries and evict the least recently used entry when insertion exceeds capacity.

Acceptance Criteria
- LRU tests demonstrate read promotion and eviction of the least recently used entry at capacity.

### REQUIREMENT REQ-collections-004

Priority queues SHALL order values by the supplied comparator and sets SHALL provide non-mutating union, intersection, and difference.

Acceptance Criteria
- Priority-queue tests pop values in comparator order.
- Set algebra tests preserve both operands while returning the expected union, intersection, and difference.

### REQUIREMENT REQ-collections-005

Pools SHALL honor their maximum retained size and factory behavior, and SyncMap SHALL preserve typed load/store/delete/range semantics.

Acceptance Criteria
- Pool tests cover maximum retained size, factory fallback, and prefill behavior.
- SyncMap tests cover typed load, store, delete, load-or-store, load-and-delete, and range behavior.

### REQUIREMENT REQ-collections-006

Empty retrieval operations SHALL return the element type's zero value and a false presence result.

Acceptance Criteria
- Empty queue, stack, heap, ring-buffer, pool, map, and cache retrieval tests return the documented zero value and `false` result where applicable.
