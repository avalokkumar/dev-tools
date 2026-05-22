# ADR 0003 — mcp-go SDK

**Status:** Accepted
**Date:** 2026-05-09

## Context

Need an MCP server SDK in Go. Candidates: `github.com/mark3labs/mcp-go`, `github.com/metoro-io/mcp-golang`.

## Decision

Use `github.com/mark3labs/mcp-go`. Most active, widest adoption, closest to spec.

Wrap it in a thin `internal/mcpserver` adapter (≤200 LOC) so the underlying SDK can be swapped if needed. The Registry is SDK-independent.

## Consequences

- MCP transport is stdio only for MVP. HTTP/SSE deferred.
- JSON Schemas are auto-derived from Go structs via `github.com/invopop/jsonschema` to avoid hand-written drift.
- Pre-1.0 SDK churn is contained behind the adapter.
