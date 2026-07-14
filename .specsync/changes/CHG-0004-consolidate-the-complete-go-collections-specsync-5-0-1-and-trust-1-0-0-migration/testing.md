---
change: CHG-0004-consolidate-the-complete-go-collections-specsync-5-0-1-and-trust-1-0-0-migration
artifact: testing
---

# Testing

- Run `fledge lanes run verify` for formatting, vet, the race-enabled suite, and build.
- Run released `specsync check --strict --require-coverage 100 --force`.
- Confirm all four agent integrations.
- Run `fledge trust doctor` and `fledge trust verify`.
- Require exact-head hosted Trust, Go 1.22 through 1.24, and CodeQL.
