# ADR 0001 — Go core, TypeScript Web UI, Python adjuncts

**Status:** Accepted
**Date:** 2026-05-09

## Context

DevForge needs a single distributable binary, a rich web UI, and a future home for ML/NLP services (log analyzer, knowledge-base compiler).

## Decision

- **Core (CLI + MCP server + Engines):** Go. Single static binary, fast cold-start, strong concurrency, mature MCP SDK (`mcp-go`), excellent cross-compile.
- **Web UI:** TypeScript + React + Vite. Embedded into the Go binary via `embed.FS`.
- **Python adjuncts:** reserved slot under `adjuncts/`. Ship later via the plugin transport (stdio JSON-RPC) so they're language-agnostic and don't bloat the core binary.

## Consequences

- One artifact to ship, install, and version.
- Web UI dev experience uses Vite hot-reload via `--dev` flag in `webserver` (proxies `/` to `localhost:5173`).
- ML features stay out of the core; users opt-in by installing an adjunct.
