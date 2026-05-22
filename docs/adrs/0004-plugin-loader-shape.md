# ADR 0004 — Plugin loader (out-of-process JSON-RPC)

**Status:** Accepted
**Date:** 2026-05-09

## Context

Plugins must add Tools at runtime. Options:
1. Go's built-in `plugin` package (in-process, .so files). Brittle, cgo-heavy, Go-only, version-pinned.
2. WASM modules. Promising but immature ecosystem, limited stdlib access.
3. Out-of-process binaries over stdio JSON-RPC (HashiCorp `go-plugin` style).

## Decision

Out-of-process binaries. Each plugin is a separate executable with a `manifest.toml`, communicating with the core over stdio JSON-RPC.

## Consequences

- Plugins can be written in any language (Go now, Python adjuncts later).
- Plugin crashes do not crash the core.
- ~1–5 ms IPC overhead per call. Acceptable for utility latency budget.
- Plugin discovery: `$DEVFORGE_PLUGIN_DIR` (default `~/.devforge/plugins/`). Each subdirectory contains a binary + manifest.
- Sandboxing (landlock/seatbelt) and signing/trust deferred to post-MVP.
