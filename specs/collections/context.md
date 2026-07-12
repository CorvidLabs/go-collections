---
spec: collections.spec.md
---

## Context

Applications repeatedly need a small common set of generic data structures
whose concurrency and capacity behavior is explicit and tested.

## Related Modules

- Go `container/heap`, `container/list`, and `sync` primitives.

## Design Decisions

- Prefer small purpose-built types over one configurable abstraction.
- Return zero-value/boolean pairs using standard Go lookup conventions.
