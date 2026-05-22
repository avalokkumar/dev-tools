# ADR 0007 — No GPL code from reference projects

**Status:** Accepted
**Date:** 2026-05-10

## Context

Phase E catalogued tool ideas from `it-tools` (CorentinTh, GPL-3.0) and
`DevToys`. Both are excellent inspiration sources but copying their code
would require DevForge to adopt the same license.

## Decision

DevForge uses these projects as **inspiration only**: tool names and
categories. Every engine is re-implemented in Go from scratch using:

- the language standard library, and/or
- permissively-licensed Go libraries (MIT, BSD, Apache-2.0, MPL-2.0)

PRs that copy GPL/AGPL/SSPL/CC-BY-SA code into the repo are rejected
regardless of how trivial the copy seems. Reviewers ask "did you copy code
from a GPL source?" on every PR introducing a new engine.

## Consequences

- DevForge ships under a permissive license (TBD — MIT or Apache-2.0).
- The repo can be vendored, embedded, or relicensed without negotiating
  copyleft compatibility.
- Every dependency added in Phase E was license-vetted; the running list
  is in `go.mod` and reviewable via `go mod why`.
