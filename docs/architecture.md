# DevForge — Architecture

## One-paragraph summary

Each utility is a **deep module** in `pkg/<tool>/` exposing a tiny public
API (`Format`, `Generate`, `Decode`, etc.) that returns a structured
`Result` and `[]Diagnostic`. Three **Surfaces** (CLI, Web, MCP) implement
thin **Adapters** that translate their native request format into the
engine's `Options` struct. A single in-process **Registry** is the source
of truth consumed by both the Web server and the MCP server. **Plugins**
are out-of-process binaries that register additional Tools at startup over
stdio JSON-RPC. The Web UI's TypeScript client and the OpenAPI doc are
**generated** from Go struct tags, so the three Surfaces cannot drift.

## Layered view

```text
┌──────────────────────────────────────────────────────────────┐
│ Surfaces (user / agent entry points)                         │
│   CLI (cobra)    Web UI (React)    MCP server (stdio)        │
└─────┬───────────────────┬─────────────────────┬──────────────┘
      │                   │                     │
      ▼                   ▼                     ▼
┌──────────────────────────────────────────────────────────────┐
│ Adapters (translation only)                                  │
│   internal/cli/tools/*    internal/webserver/api.go          │
│   internal/mcpserver/tools_*.go                              │
└─────────────────────┬────────────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────────────┐
│ Registry (single source of truth)                            │
│   internal/mcpserver/registry.go                             │
│     - Tool list                                              │
│     - JSON Schema per Operation                              │
│     - Adapter dispatch                                       │
└─────────────────────┬────────────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────────────┐
│ Engines (pure Go, deep modules) — 27 packages                │
│   Phase A–D : pkg/uuidx, pkg/jsonfmt, pkg/yamlfmt,           │
│               pkg/csvfmt, pkg/smartdiff, pkg/regextool,      │
│               pkg/cronx, pkg/jwtx, pkg/tzconv, pkg/faker     │
│   Phase E   : pkg/enc, pkg/strx, pkg/timex, pkg/sqlfmt,      │
│               pkg/mdx, pkg/idx, pkg/colorx, pkg/datax,       │
│               pkg/gitx, pkg/devx, pkg/netx, pkg/httpx,       │
│               pkg/cryptox, pkg/totpx, pkg/codefmt,           │
│               pkg/mathx, pkg/ipx                             │
└───────────────────────────────────────────────────────────────┘
          ▲  pkg/httpx + pkg/netx route through internal/netguard
          │  (resolve → reject private IPs unless allowPrivate=true)
┌─────────┴───────────────────────────────────────────────────┐
│ Security gate — internal/netguard (ADR 0006)                │
└──────────────────────────────────────────────────────────────┘

         ┌──────────────────────────────────────────┐
         │ Plugins (out-of-process, JSON-RPC stdio) │
         │   plugins/example-hello/   (Go)          │
         │   plugins/example-python/  (Python)      │
         │   adjuncts/<future-python-services>/     │
         └─────────────────┬────────────────────────┘
                           │ register at startup
                           ▼
                       Registry
```

## Module boundaries (gray-box rules)

+ **Engines** know nothing about CLI, HTTP, MCP. Imports limited to
  stdlib + domain-specific libs.
+ **Adapters** know their Surface and the Engine's signature. Nothing
  else. They do not call other Adapters.
+ **Registry** knows Tools and Operations. It does not know Surface
  implementation details.
+ **Surfaces** know the Registry. They do not import Engines directly.
+ **Plugins** know only the JSON-RPC wire format defined in
  `internal/plugin/transport.go`.

## Data flow (sample: `devforge uuid gen --version 7`)

1. Cobra parses args → `internal/cli/tools/uuid.go` Adapter.
2. Adapter builds `uuidx.GenerateOptions{Version: 7, Count: 1}`.
3. Adapter calls `uuidx.Generate(opts)` → `(GenerateResult, error)`.
4. Adapter formats `Result.Values` as text (default) or JSON (`--json`).
5. Adapter writes to `os.Stdout`, exit 0 unless `Diagnostics` carry severity Error.

Same Engine call from Web Surface:

1. `POST /api/v1/uuid/generate` body → struct via OpenAPI-generated handler.
2. Handler calls `uuidx.Generate(opts)`.
3. Handler returns `Result` as JSON.

Same Engine call from MCP Surface:

1. MCP `tools/call` for `uuid_generate` with arguments → registered
   MCP-Tool handler.
2. Handler calls `uuidx.Generate(opts)`.
3. Handler returns structured content.

All three call paths invoke the **same** function with the **same** struct.
No drift possible.

## ADR index

+ [ADR 0001 — Go core, TS web UI, Python adjuncts](./adrs/0001-go-core-ts-ui.md)
+ [ADR 0002 — Engines as deep modules](./adrs/0002-engine-as-deep-module.md)
+ [ADR 0003 — mcp-go SDK](./adrs/0003-mcp-go-sdk.md)
+ [ADR 0004 — Plugin loader shape (out-of-process JSON-RPC)](./adrs/0004-plugin-loader-shape.md)
+ [ADR 0005 — OpenAPI generated from Go](./adrs/0005-openapi-generated-from-go.md)
+ [ADR 0006 — netguard private-IP guard](./adrs/0006-netguard.md)
+ [ADR 0007 — No GPL code from reference projects](./adrs/0007-no-gpl-code.md)
+ [ADR 0008 — Crypto: stdlib + golang.org/x/crypto only](./adrs/0008-crypto-stdlib-only.md)
+ [ADR 0009 — Phase E dependency vetting](./adrs/0009-phase-e-deps.md)

## Locked invariants

1. Engine signatures (see PRD/plan) are frozen before Phase C. Changes
   require an ADR.
2. No business logic in Adapters.
3. No `panic` in Engines for user-input issues. Use `Diagnostic`.
4. Registry is the only thing Web and MCP Surfaces share. Adapters
   never cross-call.
5. Plugin transport is stdio JSON-RPC only. No in-process Go plugins.
6. **Network-touching engines (`pkg/httpx`, `pkg/netx`) refuse private
   IPs unless `allowPrivate=true` is set explicitly. Enforced via
   `internal/netguard` (ADR 0006).**
7. **Crypto (`pkg/cryptox`, `pkg/totpx`) uses only the Go standard
   library and `golang.org/x/crypto`. No third-party cipher
   implementations (ADR 0008).**
8. **No GPL/AGPL/SSPL/CC-BY-SA code is copied into the repo. Tool
   ideas may be inspired by such projects, but every engine is
   re-implemented (ADR 0007).**

## Current size (after Phase E9)

+ **31 user-facing Tools** registered as `<tool>_<op>` names.
+ **75 built-in Operations** across 28 adapter packages under
  `internal/tools/` and 27 engine packages under `pkg/`.
+ One internal helper (`internal/netguard`) gates outbound network calls.
+ Reference plugins: `plugins/example-hello` (Go), `plugins/example-python`
  (Python). Both are exercised by `scripts/mcp-smoke.py`.

Deep dive into every Operation, every engine, every Surface invariant:
see `README.md`, `docs/terms.md`, and the ADRs under `docs/adrs/`.
