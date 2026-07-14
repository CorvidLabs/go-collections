---
change: CHG-0003-correct-the-go-collections-rollout-contract-and-verification-configuration-witho
artifact: context
---

# Context

The initial rollout documents existing collection behavior. Review identified one overbroad concurrency statement: set algebra acquires receiver and argument read locks in caller order, so opposite-direction operations can form a lock cycle when writers are pending. Workflow YAML was also counted as product source because the root source directory had no Go-only extension filter.
