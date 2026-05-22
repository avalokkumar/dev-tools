# DevForge — Ubiquitous Language

Shared vocabulary used across code, docs, commits, PRs, and chat. Use these terms exactly. Do not invent synonyms.

## Core terms

| Term | Definition |
|------|------------|
| **Tool** | A user-facing capability surfaced by DevForge (e.g., "JWT debugger", "UUID generator"). One Tool may expose multiple Operations. |
| **Operation** | A single named action on a Tool (e.g., `decode`, `verify`, `generate`). Maps 1:1 to an exported engine function. |
| **Engine** | The pure Go package under `pkg/<tool>/` that implements a Tool's Operations. No I/O, no UI, no MCP knowledge. Pure functions. |
| **Surface** | An entry point that lets a user invoke an Operation. DevForge has exactly three Surfaces: **CLI**, **Web**, **MCP**. |
| **Adapter** | Glue code that translates a Surface's request format (cobra args / HTTP JSON / MCP `tools/call`) into an Engine call and back. Two homes: `internal/cli/tools/<tool>.go` for dedicated cobra subcommands, and `internal/tools/<tool>/` for Registry adapters consumed by the Web and MCP Surfaces. |
| **MCP-Tool** | The Model Context Protocol representation of an Operation. Has a JSON-Schema'd input and structured output. |
| **Registry** | Single in-process catalog of Tools/Operations. Source of truth consumed by both Web server and MCP server. Located at `internal/mcpserver/registry.go`. |
| **Plugin** | An out-of-process binary that adds Tools at runtime by registering with the Registry over stdio JSON-RPC. |
| **Diagnostic** | A structured user-facing message (info / warn / error) with a stable `Code`. Returned in `Result`, never via `panic`. |
| **Result** | The success-shaped return value of any Engine call. Always includes a `Diagnostics` slice for non-fatal issues. |
| **Spec** | A user-supplied schema (e.g., faker field list, JSON Schema) that parameterises an Operation. |
| **Walking Skeleton** | The Phase-A artifact: end-to-end binary with all three Surfaces wired but zero (or one demo) Tools registered. |
| **Golden Test** | A test that compares Engine output against a checked-in `testdata/golden/*.txt` fixture. Update via `-update` flag. |
| **Flavor** | A dialect variant of a domain language (cron unix/quartz/aws, regex re2/pcre). Encoded as an `Options` field. |
| **Adjunct** | A non-Go service (Phase 2+) like the Python log analyzer, communicating via the same Plugin transport. |

## Distinctions

- **Engine vs Adapter:** Engines do logic. Adapters do translation. Never put logic in an Adapter.
- **Diagnostic vs error:** `Diagnostic` = user-input problem (carry on). `error` (Go return) = engine bug or catastrophic failure.
- **Tool vs Operation:** "JWT" is a Tool. "decode" and "verify" are Operations of that Tool.
- **Plugin vs Adjunct:** Plugin = third-party binary registering Tools. Adjunct = first-party non-Go service (Python ML) using the plugin transport.
