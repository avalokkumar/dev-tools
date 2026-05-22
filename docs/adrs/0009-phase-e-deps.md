# ADR 0009 — Phase E dependency vetting

**Status:** Accepted
**Date:** 2026-05-10

## Context

Phase E expansion added several new direct Go dependencies. Per the
locked rule (no third-party crypto, no GPL), each must be vetted.

## Decision — accepted dependencies

| Module | License | Purpose | Notes |
|---|---|---|---|
| `golang.org/x/crypto` | BSD-3-Clause | bcrypt, argon2id, PBKDF2 | Official Go subrepo |
| `golang.org/x/net` | BSD-3-Clause | `html` parser/renderer | Official Go subrepo |
| `github.com/gomarkdown/markdown` | BSD-2-Clause | Markdown → HTML | MIT-compatible |
| `github.com/microcosm-cc/bluemonday` | BSD-3-Clause | HTML sanitiser (UGC policy) | Industry standard |
| `github.com/oklog/ulid/v2` | Apache-2.0 | ULID generator | Hashicorp/sourcegraph stack |

Pre-existing (Phases A–D, restated for completeness):

| Module | License |
|---|---|
| `github.com/spf13/cobra` | Apache-2.0 |
| `github.com/go-chi/chi/v5` | MIT |
| `github.com/mark3labs/mcp-go` | MIT |
| `gopkg.in/yaml.v3` | MIT + Apache-2.0 |
| `github.com/golang-jwt/jwt/v5` | MIT |
| `github.com/robfig/cron/v3` | MIT |
| `github.com/google/uuid` | BSD-3-Clause |
| `github.com/brianvoe/gofakeit/v7` | MIT |
| `github.com/BurntSushi/toml` | MIT |

## Consequences

- All licenses are permissive (MIT/BSD/Apache-2.0). No copyleft.
- Adding a new dependency requires updating this ADR with license + reason.
- `go mod tidy` runs in CI to keep `go.sum` clean and reproducible.
