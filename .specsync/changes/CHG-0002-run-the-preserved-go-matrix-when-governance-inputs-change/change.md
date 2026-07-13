---
id: CHG-0002-run-the-preserved-go-matrix-when-governance-inputs-change
state: implementing
type: bug_fix
base_commit: 043f70384b0cc269d651bd209a33dcae1714fc87
---

# Run the preserved Go matrix when governance inputs change

## Intent

Run the preserved Go matrix when governance inputs change

## Affected Canonical Specs

- None

## Acceptance Criteria

- Both push and pull_request path filters include .trust.toml; .specsync/**; fledge.toml; and .github/workflows/trust.yml while every Go job and matrix version remains unchanged; local strict SpecSync and native Trust verification pass before closing approval; hosted matrix results remain external

## No-spec Rationale

This correction changes only CI path-filter scheduling so the repository's preserved Go matrix validates governance-only pull-request heads; jobs, matrix versions, runtime behavior, and canonical product contracts remain unchanged.
