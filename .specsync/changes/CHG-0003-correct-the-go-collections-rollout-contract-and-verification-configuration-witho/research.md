---
change: CHG-0003-correct-the-go-collections-rollout-contract-and-verification-configuration-witho
artifact: research
---

# Research

`Set.Union`, `Set.Intersection`, and `Set.Difference` acquire two `sync.RWMutex` read locks in receiver-then-argument order. The implementation has no deterministic cross-set ordering. All other collection operations lock one collection or use `sync.Map`. SpecSync 5.0.1 supports a Go-only `source_extensions` filter, which excludes governance YAML without weakening Go coverage.
