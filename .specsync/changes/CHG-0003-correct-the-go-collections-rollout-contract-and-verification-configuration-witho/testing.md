---
change: CHG-0003-correct-the-go-collections-rollout-contract-and-verification-configuration-witho
artifact: testing
---

# Testing

- REQ-collections-001: run the complete collection suite under Go's race detector and review the canonical exclusion of cross-set lock-order guarantees against `set.go`.
- Run formatting, vet, race-enabled tests, and build through `fledge lanes run verify`.
- Run strict SpecSync 5.0.1 at 100 percent file and LOC coverage.
- Confirm all four agent integrations.
- Run Trust doctor and verification.
- Require hosted Trust, Go 1.22 through 1.24, and CodeQL at the exact head.
