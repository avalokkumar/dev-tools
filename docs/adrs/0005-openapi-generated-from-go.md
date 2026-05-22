# ADR 0005 — OpenAPI generated from Go

**Status:** Accepted
**Date:** 2026-05-09

## Context

The Web Surface and the MCP Surface both need precise schemas for the same Operations. Hand-writing them invites drift.

## Decision

Go structs are the source of truth. Pipeline:

1. Annotate `Options` and `Result` structs with `json` and `jsonschema` tags.
2. Generate JSON Schemas via `github.com/invopop/jsonschema` (used directly by MCP).
3. Generate an OpenAPI 3.1 document via the Registry's emitter.
4. Generate the TypeScript client from the OpenAPI doc via `openapi-typescript` (npm).

## Consequences

- Adding a new Operation = add struct + register with Registry. TS client and MCP schema regenerate.
- A `cross-Surface contract test` runs the same input through Web + MCP and asserts byte-equal `Result` JSON.
- Build adds one codegen step (`scripts/gen-openapi.sh`) before `pnpm build`.
