---
change: CHG-0001-adopt-specsync-5-0-1-and-trust-1-0-0-governance-for-go-collections
artifact: testing
---

# Testing

- REQ-collections-001: run the complete collection suite under Go's race detector.
- REQ-collections-002: exercise FIFO order, bounded overflow, empty reads, and ring-buffer eviction.
- REQ-collections-003: exercise LRU read promotion and least-recently-used eviction.
- REQ-collections-004: exercise comparator ordering and non-mutating set algebra.
- REQ-collections-005: exercise pool capacity/factory behavior and typed SyncMap operations.
- REQ-collections-006: exercise zero-value and false-presence results for empty retrieval operations.
- Run specsync strict validation at committed threshold zero.
- Confirm all four agent integrations report installed.
- Run the Fledge verify lane and Trust doctor.
- Confirm formatting, vet, race-enabled uncached tests, and builds pass locally and in the hosted Trust lifecycle.
- Confirm the existing Go 1.22, 1.23, and 1.24 matrix remains green.
