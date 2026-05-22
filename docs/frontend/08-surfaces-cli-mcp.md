# 08 — Surface-Specific UX: CLI Reference & MCP Config

> Specifications for the non-tool pages that serve the CLI and MCP
> surfaces within the Web UI.

---

## CLI Reference Page (`/cli`)

### Purpose

A searchable, browsable command reference generated from the same cobra
tree that `devforge docs gen` uses. Developers land here to discover CLI
syntax without leaving the browser.

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  💻 CLI Reference                                        │
│                                                          │
│  [Search commands...]                                    │
│                                                          │
├──────────┬───────────────────────────────────────────────┤
│ Command  │  devforge uuid gen                            │
│ Tree     │                                               │
│          │  Generate one or more UUIDs (v4 or v7).       │
│ devforge │                                               │
│ ├ version│  Usage:                                       │
│ ├ mcp    │    devforge uuid gen [flags]                  │
│ ├ ui     │                                               │
│ ├ uuid   │  Flags:                                       │
│ │ ├ gen ◄│    --version int   UUID version (4 or 7)      │
│ │ └ hash │                    (default 4)                │
│ ├ json   │    --count int     Number of UUIDs to         │
│ │ ├ fmt  │                    generate (default 1)       │
│ │ └ val  │    --format string Output format: std,        │
│ ├ run    │                    compact, urn (default std) │
│ ├ update │    --json          Output as JSON              │
│ ├ telem  │                                               │
│ │ └ stat │  Examples:                                    │
│ └ docs   │    devforge uuid gen --version 7 --count 5    │
│   └ gen  │    devforge uuid gen --format urn             │
│          │                                               │
│          │  [ 📋 Copy command ]                          │
└──────────┴───────────────────────────────────────────────┘
```

### Layout Details

- **Left panel (240px):** Collapsible command tree. Each node shows the
  command name. Active command highlighted with `--sky-aqua` background.
  Tree uses indented lines with `├` / `└` connectors (CSS pseudo-elements).
- **Right panel (fluid):** Command detail view showing:
  - **Command path** (e.g. `devforge uuid gen`) as heading.
  - **Description** paragraph.
  - **Usage** line in monospace block.
  - **Flags table:** Name, type, default, description.
  - **Examples** in copyable code blocks.
  - **"Copy command"** button copies the full command template.
- **Search bar** at top filters the tree — fuzzy match on command name
  or description. Matching nodes expand their parent groups.

### Data Source

Two options (choose at implementation time):

1. **Static build:** Run `devforge docs gen --out web/src/cli-docs/`
   during `build-web.sh`, import the Markdown files as a JSON index.
2. **Dynamic:** Call a `/api/v1/cli/commands` endpoint that walks the
   cobra tree at runtime. (Requires a small backend addition.)

Recommendation: **Option 1** (static build) for MVP — no backend change
needed, content is always in sync with the binary.

### Generic Dispatcher Section

A dedicated section for `devforge run`:

```text
│  devforge run                                            │
│                                                          │
│  Generic dispatcher for any registered Operation.        │
│  Runs all 75 built-in operations + plugin operations.    │
│                                                          │
│  Usage:                                                  │
│    devforge run [operation_name] [flags]                  │
│    devforge run --list                                    │
│                                                          │
│  Input Methods:                                          │
│    --args '{"key":"value"}'     Inline JSON               │
│    --args-file input.json       Read from file            │
│    echo '{}' | devforge run op  Pipe from stdin           │
│                                                          │
│  Operation List: (searchable, sortable)                   │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Name              │ Description                    │  │
│  │  uuid_generate     │ Generate UUIDs (v4/v7)         │  │
│  │  json_format       │ Pretty-print or compact JSON   │  │
│  │  enc_base64_encode │ Base64-encode input             │  │
│  │  ...               │ ...                            │  │
│  └────────────────────────────────────────────────────┘  │
```

This table fetches from `GET /api/v1/operations` to show the live list
of all registered operations.

---

## MCP Config Wizard Page (`/mcp`)

### Purpose

One-click setup for connecting DevForge's MCP server to AI coding
clients (Claude Code, Cursor, Claude Desktop, custom agents). The
wizard generates the exact config snippet for the user's platform.

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  🤖 MCP Configuration                                   │
│                                                          │
│  Connect DevForge to your AI coding assistant.           │
│                                                          │
│  Step 1: Select your client                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────────────┐ │
│  │            │  │            │  │                    │ │
│  │  Claude    │  │  Cursor    │  │  Claude Desktop    │ │
│  │  Code      │  │            │  │                    │ │
│  │            │  │            │  │                    │ │
│  └──── ● ────┘  └────────────┘  └────────────────────┘ │
│                                                          │
│  Step 2: Binary path                                     │
│  [ /usr/local/bin/devforge          ]   [ Auto-detect ]  │
│                                                          │
│  Step 3: Plugin directory (optional)                     │
│  [ ~/.devforge/plugins              ]                    │
│                                                          │
│  Step 4: Copy config                                     │
│  ┌────────────────────────────────────────────────────┐  │
│  │ Config file: .mcp.json (project root)              │  │
│  │                                                    │  │
│  │ {                                                  │  │
│  │   "mcpServers": {                                  │  │
│  │     "devforge": {                                  │  │
│  │       "command": "/usr/local/bin/devforge",        │  │
│  │       "args": ["mcp"]                              │  │
│  │     }                                              │  │
│  │   }                                                │  │
│  │ }                                                  │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ 📋 Copy Config ]   [ ⬇ Download .mcp.json ]          │
│                                                          │
│  Step 5: Verify connection                               │
│  [ ▶ Run smoke test ]                                    │
│                                                          │
│  ✅ MCP server responded with 75 tools registered.       │
└──────────────────────────────────────────────────────────┘
```

### Client Cards

Three selection cards, each with the client's logo/icon:

| Client | Config File Location | Notes |
| --- | --- | --- |
| Claude Code | `.mcp.json` (project root) or global settings | Project-local recommended |
| Cursor | `~/.cursor/mcp.json` or project-local | Same JSON shape |
| Claude Desktop | `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) | Full path shown |

Selecting a card updates the config file path label and any
platform-specific notes.

### Auto-detect Binary Path

"Auto-detect" button attempts to find the `devforge` binary:

- Checks common paths: `/usr/local/bin/devforge`,
  `$HOME/go/bin/devforge`, `./bin/devforge`.
- Falls back to prompting the user to enter the path manually.
- On the web UI (served by `devforge ui`), the binary path is already
  known — pre-fill with the path of the running process.

### Smoke Test

"Run smoke test" button sends a test request to `GET /api/v1/operations`
and verifies the response contains the expected 75+ operations. Shows:

- ✅ success count badge.
- ❌ error with details if the server is unreachable.

### Tool Discovery Table

Below the wizard, show a collapsible "Registered Tools" section listing
all 75 operations with their names and descriptions — fetched live from
the running server. This helps users verify what their AI agent will
have access to.

```text
│  Registered Tools (75)                    [ Expand All ] │
│  ┌────────────────────────────────────────────────────┐  │
│  │  ▶ Generators (12 ops)                             │  │
│  │  ▶ Formatters (14 ops)                             │  │
│  │  ▶ Converters (22 ops)                             │  │
│  │  ▶ Analyzers (20 ops)                              │  │
│  │  ▶ DevOps (7 ops)                                  │  │
│  └────────────────────────────────────────────────────┘  │
```

### Sample Prompts Section

A curated list of example prompts that users can copy-paste into their
AI client to test each tool category:

```text
│  Try these prompts with your AI:                         │
│                                                          │
│  "Generate three v7 UUIDs."                         📋  │
│  "Decode this JWT: eyJhbGci..."                     📋  │
│  "Format this JSON: {\"a\":1}"                      📋  │
│  "What's the strength of password 'hunter2'?"       📋  │
│  "Diff these two JSON objects..."                   📋  │
│  "Compute usable hosts in 10.0.0.0/24."            📋  │
```

---

## Settings Page (`/settings`)

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  ⚙️ Settings                                             │
│                                                          │
│  ┌──── Appearance ────────────────────────────────────┐  │
│  │  Theme:  ( ● Light  ○ Dark  ○ System )             │  │
│  │  Sidebar: ( ● Expanded  ○ Collapsed )              │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──── Telemetry ─────────────────────────────────────┐  │
│  │  ☐ Enable anonymous usage telemetry                │  │
│  │  Events are stored locally at                      │  │
│  │  ~/.devforge/events.jsonl. No data is uploaded.    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──── Plugins ───────────────────────────────────────┐  │
│  │  Plugin Directory: [ ~/.devforge/plugins ]         │  │
│  │                                                    │  │
│  │  Loaded Plugins:                                   │  │
│  │    example-hello   v1.0.0   ✅ Healthy   (2 ops)  │  │
│  │                                                    │  │
│  │  No plugins found? See the Plugin SDK docs.        │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──── About ─────────────────────────────────────────┐  │
│  │  DevForge v0.1.0                                   │  │
│  │  Commit: abc1234                                   │  │
│  │  Built: 2024-07-15                                 │  │
│  │  75 Operations registered                          │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

### Persistence

- **Theme:** `localStorage` key `devforge-theme` (values: `light`,
  `dark`, `system`). Applies a class on `<html>`.
- **Sidebar state:** `localStorage` key `devforge-sidebar`.
- **Favourites:** `localStorage` key `devforge-favourites` (JSON array
  of tool paths).
- **Telemetry:** Corresponds to the `DEVFORGE_TELEMETRY` env var /
  `~/.devforge/telemetry.enabled` file. The web UI shows the current
  state read-only (actual toggle requires CLI or file system change).
