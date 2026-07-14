## MODIFIED

### SPEC SECTION Invariants

1. Single-collection operations are safe for concurrent callers; callers externally serialize opposite-direction set-algebra operations over the same pair when concurrent mutations may be pending.
2. Empty lookup/pop operations return the element type's zero value with `false`.
3. Non-positive capacities normalize to a usable minimum where the type requires capacity.
4. Bounded queues reject pushes at capacity; ring buffers accept pushes and report the oldest evicted value.
5. LRU reads promote entries and insertion beyond capacity evicts the least recently used entry.
6. Priority queues always pop according to the caller's comparator.
7. Set algebra returns new sets without mutating either operand.
8. Pool size never exceeds its configured maximum and factory-created values are returned when empty.

### REQUIREMENT REQ-collections-001

Exported single-collection operations SHALL remain safe under concurrent use. Callers SHALL externally serialize opposite-direction `Union`, `Intersection`, or `Difference` calls over the same set pair when concurrent mutations may be pending.

Acceptance Criteria
- The complete package test suite passes under Go's race detector.
- The canonical contract does not claim deadlock-free lock ordering across two independently locked sets.
