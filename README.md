# DevTools

> **Unified, MCP-native developer toolkit** — bundles **75 operations** across 30+ tools behind a single CLI, a local web UI, and an MCP stdio server. Offline-first. One static binary. Plugin-extensible.

```
┌────────────────────────────────────────────────────────────────┐
│                       dev-tools (single binary)                 │
├────────────────────────────────────────────────────────────────┤
│   CLI                Web UI (embedded)         MCP server      │
│   ─────              ──────────────────        ──────────────  │
│   dev-tools uuid gen  http://localhost:port     stdio JSON-RPC  │
│   dev-tools run …     /api/v1/<tool>/<op>       tools/list,call │
└──────────────────┬───────────────┬─────────────────┬───────────┘
                   ▼               ▼                 ▼
           ┌────────────────────────────────────────────────┐
           │  Registry (single source of truth: 75 ops)     │
           └────────────────────────────────────────────────┘
                                 │
                                 ▼
        ┌──────────────────────────────────────────────────┐
        │  Engines (pure Go, 27 packages):                 │
        │  uuidx jsonfmt yamlfmt csvfmt smartdiff          │
        │  regextool cronx jwtx tzconv faker               │
        │  enc strx timex sqlfmt mdx idx colorx            │
        │  datax gitx devx netx httpx                      │
        │  cryptox totpx codefmt mathx ipx                 │
        └──────────────────────────────────────────────────┘
```

<img width="1703" height="1241" alt="image" src="https://github.com/user-attachments/assets/c6936054-a8f1-4c7b-a7b9-04180e720742" />


<img width="1702" height="1235" alt="image" src="https://github.com/user-attachments/assets/74a36959-c1fd-4a3d-a0c9-6794ee472207" />

---

## Table of contents

- [Why DevTools](#why-dev-tools)
- [Quick start](#quick-start)
- [Three surfaces](#three-surfaces)
- [Tool catalog (75 operations)](#tool-catalog-75-operations)
- [CLI usage by tool](#cli-usage-by-tool)
- [Web UI usage](#web-ui-usage)
- [MCP usage](#mcp-usage-claude-code--cursor--claude-desktop)
- [Plugins](#plugins)
- [Operations & maintenance](#operations--maintenance)
- [Development](#development)
- [Project structure](#project-structure)
- [Releasing](#releasing)
- [Roadmap](#roadmap)

---

## Why DevTools

Modern developers waste hours context-switching between regex testers, JSON formatters, JWT decoders, cron builders, and 30+ other single-purpose websites. None of them expose deterministic interfaces that AI coding agents (Claude Code, Cursor) can call without hallucinating.

**DevTools fixes both problems**:

- **Bundled utilities** — one binary, one UI, one CLI for 30+ dev workflows.
- **MCP-native** — every operation is also exposed as a Model Context Protocol tool, so AI agents call your local DevTools instead of guessing endpoints.
- **Offline-first** — runs entirely on your machine. No accounts, no telemetry without opt-in, no network required.
- **Plugin SDK** — extend the registry with out-of-process plugins in any language.
- **Single source of truth** — the same `Registry` powers the CLI, Web, and MCP, so the three surfaces never drift.

---

## Quick start

### Prerequisites

- **Go 1.25+** (only required for building from source)
- **Node 20+ and pnpm 10+** (only required if rebuilding the web UI)

### Build from source

```bash
git clone https://github.com/avalokkumar/dev-tools.git
cd dev-tools

# Build the web UI and embed it into the binary
./scripts/build-web.sh

# Build the binary
go build -o ./bin/dev-tools ./cmd/dev-tools

# Verify
./bin/dev-tools version --json
# → {"version":"0.0.0-dev","commit":"none","date":"unknown"}
```

### Smoke test in 30 seconds

```bash
# Generate three v7 UUIDs
./bin/dev-tools uuid gen --version 7 -n 3

# List all 75 operations the binary exposes
./bin/dev-tools run --list

# Open the web UI on a free port
./bin/dev-tools ui
# → DevTools UI listening on http://127.0.0.1:54321

# Start the MCP server (stdio); Ctrl-D to exit
./bin/dev-tools mcp

# Inspect telemetry status (off by default)
./bin/dev-tools telemetry status

# Simulate self-update (no release channel wired in MVP)
./bin/dev-tools update --dry-run

# Generate Markdown reference docs for every CLI subcommand
./bin/dev-tools docs gen --out ./generated-cli-docs
```

---

## Three surfaces

Every operation is reachable three ways. **The same input produces a byte-equal JSON result on all three surfaces** — there is no drift.

| Surface | Entry point | Use it when |
|---|---|---|
| **CLI** | `dev-tools <tool> <op>` (dedicated) <br> `dev-tools run <name>` (generic) | Scripting, one-shot terminal work, piping with other shell tools |
| **Web UI** | `dev-tools ui` then open the printed URL | Interactive exploration, copy-paste workflows |
| **MCP** | `dev-tools mcp` over stdio | AI agent integration (Claude Code, Cursor, Claude Desktop) |

### CLI conventions

- Global flag `--json` makes any command emit machine-readable JSON.
- Operations registered with the Registry are addressable as `dev-tools run <tool>_<op>`.
- A few high-frequency tools have dedicated cobra subcommands too: `dev-tools uuid …`, `dev-tools json …`.
- Most subcommands accept `--in <file>` or stdin, and write results to stdout.

### HTTP conventions

- All routes are `POST /api/v1/<tool>/<op>` with a JSON body matching the operation's input schema.
- `GET /api/v1/operations` enumerates the registry (handy for clients).
- `GET /healthz` is the liveness probe.
- Server binds to `127.0.0.1:0` (ephemeral port) by default. Override with `--addr 127.0.0.1:8787`.

### MCP conventions

- Stdio transport only (HTTP/SSE deferred).
- Tool names are `<tool>_<op>` (e.g. `uuid_generate`, `jwt_verify`, `cron_next`).
- Each tool ships its JSON Schema as part of `tools/list`.
- Plugins discovered via `$DEVFORGE_PLUGIN_DIR` (default `~/.dev-tools/plugins`) appear alongside built-ins.

---

## Tool catalog (75 operations)

DevTools groups operations into **categories**. Every operation is reachable
through CLI (`dev-tools run <op>`), Web UI (`POST /api/v1/<tool>/<op>`), and
MCP (`tools/call`). Examples below show only the most common per category;
the full list is `dev-tools run --list` (75 lines) or `GET /api/v1/operations`.

### Identifiers + hashing

| Tool | Operation | What it does |
|---|---|---|
| **uuid** | `uuid_generate` | Generate v4 or v7 UUIDs (with format options) |
| **uuid** | `uuid_hash` | Compute MD5/SHA-1/SHA-256/SHA-512 digests |
| **id** | `id_ulid` | Generate Crockford-base32 ULIDs (time-ordered) |
| **id** | `id_slug` | Slugify arbitrary text (URL-safe identifier) |

### Encoding / decoding

| Tool | Operation | What it does |
|---|---|---|
| **enc** | `enc_base64_{encode,decode}` | Base64 (RFC 4648, URL-safe + no-padding variants) |
| **enc** | `enc_url_{encode,decode}` | URL percent-encoding (component / path modes) |
| **enc** | `enc_html_{encode,decode}` | HTML entity round-trip |
| **enc** | `enc_hex_{encode,decode}` | Hex (tolerates `0x` prefix on decode) |

### Crypto + auth

| Tool | Operation | What it does |
|---|---|---|
| **crypto** | `crypto_aes_{encrypt,decrypt}` | AES-256-GCM with raw key or PBKDF2 passphrase |
| **crypto** | `crypto_rsa_keygen` | RSA keypair (PEM, 2048/3072/4096) |
| **crypto** | `crypto_hmac` | HMAC-SHA1/256/384/512 (RFC 4231 vectors verified) |
| **crypto** | `crypto_password_hash` | bcrypt + argon2id |
| **crypto** | `crypto_password_strength` | Entropy-based 0–4 score + suggestions |
| **totp** | `totp_{generate,verify}` | RFC 6238 TOTP (vectors verified) |
| **jwt** | `jwt_{decode,verify}` | JWT parse + signature/alg/expiry check |

### Networking

| Tool | Operation | What it does |
|---|---|---|
| **http** | `http_request` | curl-like HTTP runner — **localhost-only by default**, opt-in `allowPrivate` |
| **dns** | `dns_lookup` | A/AAAA/MX/TXT/CNAME/NS/PTR — same private-IP guard |
| **url** | `url_parse` | Decompose URL + per-param query breakdown |
| **headers** | `headers_analyze` | Static security-header checklist |

### Data + format

| Tool | Operation | What it does |
|---|---|---|
| **json** | `json_{format,validate}` | Pretty/compact + schema check |
| **yaml** | `yaml_{format,validate,convert}` | Reformat / validate / YAML↔JSON |
| **csv** | `csv_{format,validate}` | Column align + header enforcement |
| **data** | `data_json_to_csv` / `data_csv_to_json` | Bidirectional conversion |
| **data** | `data_json_to_xml` / `data_xml_to_json` | Bidirectional conversion |
| **data** | `data_flatten` / `data_unflatten` | Dotted-path flatten / restore |
| **data** | `data_key_rename` | Pattern-based JSON key rename |
| **diff** | `diff_compare` | Semantic JSON / INI / SQL diff |

### Text + time

| Tool | Operation | What it does |
|---|---|---|
| **str** | `str_case` | Convert between camel/snake/kebab/pascal/constant/dot/train/title/lower/upper |
| **str** | `str_diff` | Line-based unified diff |
| **str** | `str_stats` | Lines / words / chars / bytes / longest line |
| **str** | `str_sort_unique` | Sort lines + optional uniq |
| **str** | `str_replace` | Literal or regex replace; case-insensitive option |
| **time** | `time_convert` | Epoch ↔ RFC3339 (auto-detects ms/us/ns) |
| **time** | `time_relative` | Human phrase between two instants |
| **time** | `time_duration` | Parse `1h30m` and break it down |
| **tz** | `tz_{convert,list}` | DST-aware wall-time conversion |
| **regex** | `regex_{test,explain}` | RE2 matches + plain-English explanation |
| **cron** | `cron_{parse,next}` | Cron decompose + next N runs |

### Code, markdown, SQL

| Tool | Operation | What it does |
|---|---|---|
| **code** | `code_fmt_go` / `code_fmt_xml` / `code_fmt_html` | Pure-Go formatters |
| **md** | `md_to_html` | Markdown → sanitised HTML (bluemonday UGC) |
| **md** | `md_table_from_csv` | Generate Markdown table from CSV |
| **sql** | `sql_format` | Indent + casing |
| **sql** | `sql_validate` | Lints (SELECT *, missing WHERE, unbalanced parens) |

### Git + DevOps

| Tool | Operation | What it does |
|---|---|---|
| **git** | `git_patch` | Build a unified-diff patch from two text inputs |
| **git** | `git_commit_format` | Conventional Commits v1 validator |
| **git** | `git_ignore_gen` | `.gitignore` from curated templates |
| **dockerfile** | `dockerfile_lint` | Best-practice rules (no :latest, USER, single CMD/ENTRYPOINT) |
| **env** | `env_parse` / `env_diff` | Validate `.env` syntax + add/remove/change diff |
| **k8s** | `k8s_validate` | Structural manifest check |

### Math, IDs, color, IP

| Tool | Operation | What it does |
|---|---|---|
| **math** | `math_eval` | Safe arithmetic expression evaluator |
| **math** | `math_unit` | Bytes / time / throughput / temperature / length conversion |
| **ip** | `ip_calc` | IPv4/IPv6 subnet math (network, broadcast, hosts, mask) |
| **color** | `color_convert` | hex ↔ rgb ↔ hsl + CSS named colors |
| **faker** | `faker_{generate,kinds}` | Synthetic data generator |

For the full enumeration: `dev-tools run --list` prints all 75 operations.

---

## CLI usage by tool

### UUID

```bash
# v4 (random) — default
dev-tools uuid gen

# Three v7 (time-ordered) UUIDs
dev-tools uuid gen --version 7 --count 3

# Compact format (no dashes)
dev-tools uuid gen --format compact

# JSON output
dev-tools uuid gen --json --count 2

# Hash a string with multiple algorithms
dev-tools uuid hash --input "hello" --algo sha256 --algo md5
```

### JSON

```bash
# Pretty-print stdin
echo '{"b":1,"a":[2,3]}' | dev-tools json fmt --indent 2

# Compact form
echo '{"b": 1}' | dev-tools json fmt --indent 0

# Sort keys + trailing newline
echo '{"b":1,"a":2}' | dev-tools json fmt --sort-keys --trailing-newline

# Validate file
dev-tools json validate --in config.json

# Validate against a schema
dev-tools json validate --in config.json --schema schema.json
```

### Generic dispatcher (every other tool)

The `run` subcommand can invoke any registered operation. Pass JSON arguments via `--args`, `--args-file`, or stdin.

```bash
# YAML → JSON
dev-tools run yaml_convert --args '{"input":"a: 1\nb: 2","to":"json"}'

# Reformat CSV with column alignment
dev-tools run csv_format --args '{"input":"a,b\n1,2","alignColumns":true}'

# Diff two JSON snippets
dev-tools run diff_compare --args '{"left":"{\"a\":1}","right":"{\"a\":2}","mode":"json"}'

# Regex match
dev-tools run regex_test --args '{"pattern":"\\d+","input":"build 42 release 7"}'

# Plain-English regex explanation
dev-tools run regex_explain --args '{"pattern":"^\\d{4}-\\d{2}$"}'

# Cron field breakdown
dev-tools run cron_parse --args '{"expression":"0 9 * * MON-FRI"}'

# Next 3 cron runs in NY time
dev-tools run cron_next --args '{"expression":"*/15 * * * *","n":3,"tz":"America/New_York"}'

# Decode a JWT
dev-tools run jwt_decode --args '{"token":"eyJhbGciOi..."}'

# Verify a JWT (HMAC)
dev-tools run jwt_verify --args '{
  "token":"eyJhbGciOi...",
  "key":"my-secret",
  "expectedAlgs":["HS256"]
}'

# Convert wall time across zones
dev-tools run tz_convert --args '{
  "time":"2026-05-09T12:00:00Z",
  "fromTZ":"America/New_York",
  "toTZ":"Asia/Tokyo"
}'

# List all registered ops
dev-tools run --list
```

---

## Web UI usage

```bash
# Bind to ephemeral port
dev-tools ui

# Bind to fixed port (handy for bookmarks)
dev-tools ui --addr 127.0.0.1:8787

# Open the printed URL in your browser. The home page links to all tools.
```

The web UI is a single-page React app embedded into the binary via `embed.FS` — no separate web server required. Each tool page calls the `/api/v1/<tool>/<op>` endpoint directly.

For UI development with hot reload, run Vite separately:

```bash
pnpm -C web dev   # starts Vite at http://localhost:5173 (proxies /api → :8787)
```

---

## MCP usage (Claude Code / Cursor / Claude Desktop)

DevTools runs as an MCP server over stdio, so any MCP-aware client can use it.

### Claude Code

Add to your project's `.mcp.json` or your global Claude Code settings:

```json
{
  "mcpServers": {
    "dev-tools": {
      "command": "/absolute/path/to/dev-tools",
      "args": ["mcp"]
    }
  }
}
```

### Cursor

Add to `~/.cursor/mcp.json` (or project-local equivalent):

```json
{
  "mcpServers": {
    "dev-tools": {
      "command": "/absolute/path/to/dev-tools",
      "args": ["mcp"]
    }
  }
}
```

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or the equivalent on Windows/Linux:

```json
{
  "mcpServers": {
    "dev-tools": {
      "command": "/absolute/path/to/dev-tools",
      "args": ["mcp"]
    }
  }
}
```

After registering, the agent sees all 75 DevTools tools. Try prompts like:

- *"Generate three v7 UUIDs."*
- *"Decode this JWT and tell me when it expires."*
- *"Show me a regex that matches a US ZIP code, then explain it."*
- *"Convert 2026-03-08 02:30 Eastern to UTC; flag any DST issue."*
- *"Diff these two JSON snippets and tell me what changed."*
- *"Encrypt this secret with AES-256-GCM using the passphrase 'hunter2'."*
- *"What's the strength score of the password 'correct-horse-battery'?"*
- *"Generate an RFC 6238 TOTP code from this base32 secret."*
- *"Lint this Dockerfile and tell me what best-practice rules it breaks."*
- *"Validate this Kubernetes YAML manifest."*
- *"Parse this .env file and show me any duplicate or malformed lines."*
- *"Compute usable hosts for the subnet 10.0.0.0/22."*
- *"Convert this JSON array of objects to CSV."*
- *"Slugify 'Hello World! (2026)' into a URL-safe identifier."*
- *"Format this SQL query with consistent indentation and uppercase keywords."*
- *"Generate 20 rows of fake user data with name, email, and age fields."*

### Verifying MCP manually

```bash
# Initialize handshake + list tools
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke","version":"0"}}}\n{"jsonrpc":"2.0","id":2,"method":"tools/list"}\n' \
  | dev-tools mcp
```

---

## Plugins

DevTools supports out-of-process plugins written in any language. Plugins communicate with the host over stdio JSON-RPC and contribute additional Tools to the Registry — they appear in MCP `tools/list`, the HTTP API, and `dev-tools run --list` exactly like built-ins.

### Plugin discovery

```bash
export DEVFORGE_PLUGIN_DIR=~/.dev-tools/plugins   # default if unset
```

Each subdirectory of `$DEVFORGE_PLUGIN_DIR` must contain:

- `manifest.toml` — declares name, SDK version, optional entrypoint
- the executable named `plugin` (or whatever `entrypoint` says)

### Authoring a plugin in Go (10 minutes)

```go
// plugins/my-plugin/main.go
package main

import (
 "context"
 "encoding/json"
 "fmt"

 "github.com/dev-tools/dev-tools/pkg/pluginsdk"
)

func main() {
 pluginsdk.Serve(pluginsdk.Plugin{
  Name:    "my-plugin",
  Version: "0.1.0",
  Operations: []pluginsdk.OpDecl{{
   Tool:        "greet",
   Op:          "say",
   Description: "Greet by name.",
   InputSchema: json.RawMessage(`{
              "type":"object",
              "properties":{"name":{"type":"string"}}
            }`),
   Handler: func(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
    var a struct{ Name string `json:"name"` }
    _ = json.Unmarshal(args, &a)
    if a.Name == "" { a.Name = "world" }
    return json.Marshal(map[string]string{
     "message": fmt.Sprintf("hello, %s", a.Name),
    })
   },
  }},
 })
}
```

```toml
# plugins/my-plugin/manifest.toml
sdk = 1
name = "my-plugin"
version = "0.1.0"
description = "Demo plugin"
entrypoint = "plugin"

[[tools]]
tool = "greet"
description = "Greets by name."
```

```bash
cd plugins/my-plugin
go build -o plugin .

# Install
mkdir -p ~/.dev-tools/plugins/my-plugin
cp plugin manifest.toml ~/.dev-tools/plugins/my-plugin/

# Verify it shows up
dev-tools run --list | grep greet_say
```

The plugin SDK exposes:

- `pluginsdk.Serve(Plugin)` — the canonical entry point.
- `pluginsdk.OpDecl{...}` — declares one Operation.
- Crash recovery + per-call timeout are handled by the host loader.

See `plugins/example-hello/` for a working reference.

### Authoring a plugin in another language

Any executable that speaks newline-delimited JSON-RPC 2.0 over stdio can be a plugin. The host calls these methods:

- `initialize` (params: `{sdk: 1}`) → `{plugin: <name>, operations: [...]}`
- `invoke` (params: `{name: "<tool>_<op>", args: {...}}`) → arbitrary JSON Result
- `shutdown` (no params) → `null`

A working Python reference plugin lives at `plugins/example-python/`. It uses
only the standard library (no external dependencies) and contributes a
`py_reverse` Operation. To try it locally:

```bash
mkdir -p ~/.dev-tools/plugins
cp -R plugins/example-python ~/.dev-tools/plugins/

# Confirm it shows up alongside built-ins
dev-tools run --list | grep py_

# Invoke it
dev-tools run py_reverse --args '{"input":"DevTools"}'
# → {"output": "egroFveD"}
```

The Go reference plugin is at `plugins/example-hello/`.

---

## Operations & maintenance

DevTools ships with a few subcommands for keeping the binary current and
auditable.

### `dev-tools update` — self-update

```bash
# Show what an update would do, without applying
dev-tools update --dry-run

# Plain run (no release channel wired in MVP — prints a friendly note)
dev-tools update
```

Internals: the update path is delegated to a pluggable `update.Source`. The
default is a stub that reports "no release channel configured". A real
GitHub-releases source plugs in via `update.SetSource(...)` once the project
publishes signed releases via `goreleaser` (see [Releasing](#releasing)).

### `dev-tools telemetry status` — opt-in usage events

Telemetry is **off by default**. No data leaves your machine without explicit
opt-in. Two opt-in mechanisms:

```bash
# Per-invocation
DEVFORGE_TELEMETRY=1 dev-tools telemetry status

# Persistent: a non-empty marker file
mkdir -p ~/.dev-tools && echo enabled > ~/.dev-tools/telemetry.enabled
```

When enabled, events are buffered in memory and flushed to
`~/.dev-tools/events.jsonl` as newline-delimited JSON. The package never
auto-uploads — exporting events to a remote endpoint will be added in a later
phase along with a clear opt-out path.

### `dev-tools docs gen` — auto-generated CLI reference

```bash
# Produces one Markdown file per subcommand
dev-tools docs gen --out ./generated-cli-docs
ls ./generated-cli-docs   # dev-tools.md, devforge_uuid.md, devforge_run.md, …
```

These files are produced via `cobra/doc.GenMarkdownTree` and stay in sync
with the code automatically — re-run `docs gen` after adding a new
subcommand. Suitable for publishing as static-site reference docs.

### `scripts/doc-lint.sh` — guard the canonical docs

```bash
./scripts/doc-lint.sh
```

Fails CI when:

- `docs/terms.md`, `docs/architecture.md`, or `README.md` is missing
- `terms.md` is missing one of the **Five Names** (Tool / Operation / Engine
  / Surface / Adapter), which would mean the ubiquitous-language contract
  was broken
- any ADR file under `docs/adrs/` is missing a `Status:` field

---

## Development

### Common commands

```bash
# Run all Go tests with race detector
go test -race ./...

# Build and embed the web UI
./scripts/build-web.sh

# Build the single static binary
go build -o ./bin/dev-tools ./cmd/dev-tools

# Type-check the web UI
pnpm -C web exec tsc -b --noEmit

# Run web E2E (Playwright)
pnpm -C web e2e

# Regenerate OpenAPI types for the web client
./scripts/gen-openapi.sh

# Lint canonical docs (terms.md, architecture.md, ADRs)
./scripts/doc-lint.sh

# Build the example Python plugin (no compile step — just make it executable)
chmod +x plugins/example-python/plugin plugins/example-python/plugin.py

# Build the example Go plugin
go build -o plugins/example-hello/plugin ./plugins/example-hello
```

### Adding a new utility

1. Create the engine at `pkg/<tool>/<tool>.go` with 1–3 public functions returning `Result` + `[]engine.Diagnostic`. Write tests **first** (`<tool>_test.go`).
2. Create the registry adapter at `internal/tools/<tool>/<tool>.go` exporting `Operations() []mcpserver.Operation`.
3. Append the adapter to the slice in `internal/tools/tools.go`.
4. (Optional) Add a dedicated cobra subcommand at `internal/cli/tools/<tool>.go` for ergonomic CLI access. The generic `dev-tools run <op>` works without this.
5. (Optional) Add a React page at `web/src/pages/<Tool>.tsx` and a route in `web/src/App.tsx`.

The new operation is automatically reachable via Web HTTP API and MCP `tools/call` — both consume the same Registry.

### Locked invariants

- **Engine signatures** are frozen. Changes require an ADR (`docs/adrs/`).
- **No business logic in adapters** — adapters translate, engines compute.
- **Engines never `panic`** on user-input issues. Use `engine.Diagnostic`.
- **Adapters never cross-call**. They go through the Registry or the Engine.
- **Plugins are out-of-process only**. No in-process Go `plugin` package.

See `docs/architecture.md` and `docs/terms.md`.

---

## Project structure

```
dev-tools/
├── cmd/dev-tools/                # Binary entrypoint (main.go)
├── internal/
│   ├── cli/                     # cobra command tree (root, run, ui, mcp,
│   │   │                        #   update, telemetry, docs, version)
│   │   └── tools/               # CLI adapters per Tool (uuid, json)
│   ├── mcpserver/               # MCP adapter + Registry (single source of truth)
│   ├── webserver/               # HTTP API + embedded SPA
│   ├── tools/                   # Registry adapters per Tool (28 subpackages)
│   ├── plugin/                  # Out-of-process plugin loader
│   ├── update/                  # Self-update planning (pluggable Source)
│   ├── telemetry/               # Opt-in event buffer
│   ├── version/                 # Build-time stamp
│   └── ...
├── pkg/                         # PUBLIC engines (deep modules, 27 packages)
│   ├── engine/                  # shared types: Diagnostic, Severity, Span
│   ├── pluginsdk/               # plugin SDK (used by plugin authors)
│   ├── uuidx/                   # UUID + hash
│   ├── jsonfmt/, yamlfmt/, csvfmt/, smartdiff/, regextool/
│   ├── cronx/, jwtx/, tzconv/, faker/
│   ├── enc/, strx/, timex/, sqlfmt/, mdx/, idx/, colorx/
│   ├── datax/, gitx/, devx/, netx/, httpx/
│   ├── cryptox/, totpx/, codefmt/, mathx/, ipx/
├── plugins/
│   ├── example-hello/           # reference Go plugin (in its own go.mod)
│   └── example-python/          # reference Python plugin (no deps)
├── adjuncts/                    # reserved slot for Phase-2 Python services
├── web/                         # Vite + React SPA (TypeScript)
├── packages/api-types/          # Generated TS types from OpenAPI
├── api/openapi.yaml             # API spec source of truth (Phase A seed)
├── scripts/                     # build-web.sh, gen-openapi.sh, doc-lint.sh
├── docs/
│   ├── PRD.md                   # product requirements
│   ├── architecture.md          # module boundaries
│   ├── terms.md                 # ubiquitous language
│   ├── adrs/0001..0005-*.md     # architecture decision records
│   └── …
├── .goreleaser.yaml             # cross-platform release recipe
└── .github/workflows/
    ├── ci.yml                   # build + test on Linux/macOS/Windows
    └── release.yml              # tagged-release pipeline (goreleaser)
```

### The Five Names

The codebase uses these terms exactly. See `docs/terms.md`.

- **Tool** — user-facing capability ("UUID generator").
- **Operation** — single named action of a Tool ("generate").
- **Engine** — pure Go package implementing a Tool's Operations.
- **Surface** — entry point: CLI, Web, or MCP.
- **Adapter** — glue translating a Surface request into an Engine call.

---

## Releasing

Releases are produced by [goreleaser](https://goreleaser.com) and triggered by
pushing a `v*` tag. The recipe lives at `.goreleaser.yaml`; the workflow at
`.github/workflows/release.yml`.

### What gets built

- Cross-platform binaries: `linux/amd64`, `linux/arm64`, `darwin/amd64`,
  `darwin/arm64`, `windows/amd64`.
- Embedded web UI (built fresh via `scripts/build-web.sh` in the
  before-hooks).
- Per-archive bundle of `LICENSE`, `README.md`, `docs/PRD.md`,
  `docs/architecture.md`, `docs/terms.md`.
- `checksums.txt` (SHA-256) alongside the archives.

Build-time `version`, `commit`, and `date` are stamped into the binary via
`-ldflags` so `dev-tools version --json` reports the real release.

### How to cut a release

```bash
# 1. Verify everything is green locally.
go test -race ./...
./scripts/doc-lint.sh
./scripts/build-web.sh

# 2. Tag and push.
git tag v0.1.0
git push origin v0.1.0
```

The `release.yml` workflow then runs goreleaser with `--clean` and uploads
the archives to a GitHub draft release. Promote the draft to "Latest" once
release notes look right.

### Local rehearsal (optional)

```bash
# Requires goreleaser installed locally.
goreleaser release --snapshot --clean --skip publish
```

Outputs land in `dist/` for inspection without touching GitHub.

---

## Roadmap

**Shipped (Phase A + B + C + D + E)**

- Three surfaces (CLI / Web / MCP) — same byte-equal `Result` JSON across all three
- **75 operations** across 30+ tools (Phase E added 55 new ops):
  encoding, crypto suite (AES-GCM, RSA, HMAC, bcrypt+argon2id, TOTP),
  HTTP client (with private-IP guard), DNS lookup, SQL formatter +
  validator, string + time + URL utilities, semantic data converters
  (JSON↔CSV↔XML, flatten/unflatten), markdown render + table builder,
  git utilities (patch / commit / ignore), Dockerfile linter, .env
  parser + diff, K8s YAML validator, code formatters (Go / XML / HTML),
  math evaluator, unit converter, ULID, slug, color, CIDR.
- Plugin SDK with crash + timeout resilience (Go + Python examples)
- Cross-platform CI matrix (Linux / macOS / Windows) — race-free
- `internal/netguard` blocks SSRF for `http_request` / `dns_lookup` by default
- `dev-tools update --dry-run` with pluggable release source
- Opt-in telemetry (off by default; `DEVFORGE_TELEMETRY=1` to enable)
- Auto-generated Markdown CLI docs (`dev-tools docs gen --out DIR`)
- goreleaser config + release workflow for tagged builds
- Doc-lint script enforcing terms.md / architecture.md invariants
- RFC vector tests for crypto (HMAC RFC 4231, TOTP RFC 6238, Base64 RFC 4648)

**Post-MVP flagship platforms** (deferred per PRD)

- A1 — API Contract Drift Guard
- A2 — Log Analyzer Lite
- A3 — GitHub-Native Secrets Gateway
- A4 — Knowledge Base Compiler
- A5 — Playwright Healer Dashboard

---

## License

TBD.

## Contributing

Issues and pull requests welcome. Read `docs/architecture.md` and `docs/terms.md` first — they explain the patterns the codebase relies on.
