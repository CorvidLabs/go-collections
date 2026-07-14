---
spec: collections.spec.md
---

## Test Plan

### Unit Tests

- Cover constructors, empty values, normal mutation, capacity boundaries,
  ordering, eviction, clearing, draining, and set algebra for every type.
- Exercise concurrent operations under the race detector.

### Integration Tests

- Run the complete package with `go test -race -count=1 ./...` across the
  supported Go-version matrix.
