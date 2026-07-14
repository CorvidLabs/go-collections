---
change: CHG-0001-adopt-specsync-5-0-1-and-trust-1-0-0-governance-for-go-collections
artifact: context
---

# Context

go-collections has mature native Go validation but no previous SpecSync threshold or
verified SDD policy. This migration records the existing generic collections package behavior,
adds project-scoped agent guidance, and composes the native gate through Trust
without changing runtime semantics or public APIs.
