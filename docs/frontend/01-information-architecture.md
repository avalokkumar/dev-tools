# 01 — Information Architecture

> Sitemap, navigation model, routing table, and URL structure for the
> DevForge Web Surface.

---

## Navigation Model

The app uses a **sidebar + main content** layout. The sidebar is the
primary navigation surface; it is always visible on `md+` screens and
collapses to an icon rail on `sm`. On `xs` it becomes a hamburger-triggered
overlay drawer.

```text
┌────────────────────────────────────────────────────────────┐
│  Header bar  [ ⚒ DevForge ]   [ 🔍 Cmd+K ]   [ ⚙ ]      │
├──────────┬─────────────────────────────────────────────────┤
│ Sidebar  │  Main Content                                   │
│          │                                                  │
│ Home     │  ┌────────────────────────────────────────────┐ │
│          │  │  Tool Page / Dashboard / Settings          │ │
│ ─────    │  │                                            │ │
│ Category │  │  Input panel     │  Output panel            │ │
│  Tool 1  │  │                  │                          │ │
│  Tool 2  │  │                  │                          │ │
│  Tool 3  │  │                  │                          │ │
│          │  └────────────────────────────────────────────┘ │
│ ─────    │                                                  │
│ Category │  ┌────────────────────────────────────────────┐ │
│  ...     │  │  Diagnostics / History bar (collapsible)   │ │
│          │  └────────────────────────────────────────────┘ │
│ ─────    │                                                  │
│ CLI Ref  │                                                  │
│ MCP Cfg  │                                                  │
│ Settings │                                                  │
└──────────┴─────────────────────────────────────────────────┘
```

### Sidebar Sections

1. **Home** — Tool dashboard grid (default landing page).
2. **Generators** — UUID, Faker, ID (ULID/Slug), TOTP, Crypto Keygen.
3. **Formatters** — JSON, YAML, CSV, SQL, Code (Go/XML/HTML), Markdown.
4. **Converters** — Encoding (Base64/URL/HTML/Hex), Data Transforms,
   Color, Time, Timezone, Math/Unit.
5. **Analyzers** — Diff, Regex, Cron, JWT, String Tools, URL Parser,
   HTTP Headers, DNS, HTTP Client, IP Calc.
6. **DevOps** — Git (Patch/Commit/Ignore), Dockerfile Lint, Env, K8s.
7. **Surfaces** — CLI Reference, MCP Config Wizard.
8. **Settings** — Theme toggle, telemetry opt-in, plugin manager.

---

## Sitemap & Routing Table

All tool routes follow the pattern `/tools/<tool>` or
`/tools/<tool>/<op>`. API calls go to `POST /api/v1/<tool>/<op>`.

### Top-level Routes

| Path | Page | Description |
| --- | --- | --- |
| `/` | Home Dashboard | Tool grid with search, categories, favourites |
| `/settings` | Settings | Theme, telemetry, plugins |
| `/cli` | CLI Reference | Searchable command tree, auto-generated from cobra |
| `/mcp` | MCP Config Wizard | Copy-paste config for Claude / Cursor / Desktop |

### Generator Tools

| Path | Page | Operations Used |
| --- | --- | --- |
| `/tools/uuid` | UUID Generator | `uuid_generate`, `uuid_hash` |
| `/tools/faker` | Data Faker | `faker_generate`, `faker_kinds` |
| `/tools/id` | ID Generator (ULID/Slug) | `id_ulid`, `id_slug` |
| `/tools/totp` | TOTP Generator | `totp_generate`, `totp_verify` |
| `/tools/crypto` | Crypto Toolkit | `crypto_aes_encrypt`, `crypto_aes_decrypt`, `crypto_rsa_keygen`, `crypto_hmac`, `crypto_password_hash`, `crypto_password_strength` |

### Formatter Tools

| Path | Page | Operations Used |
| --- | --- | --- |
| `/tools/json` | JSON Formatter | `json_format`, `json_validate` |
| `/tools/yaml` | YAML Formatter | `yaml_format`, `yaml_validate`, `yaml_convert` |
| `/tools/csv` | CSV Formatter | `csv_format`, `csv_validate` |
| `/tools/sql` | SQL Formatter | `sql_format`, `sql_validate` |
| `/tools/code` | Code Formatter | `code_fmt_go`, `code_fmt_xml`, `code_fmt_html` |
| `/tools/markdown` | Markdown Editor | `md_to_html`, `md_table_from_csv` |

### Converter Tools

| Path | Page | Operations Used |
| --- | --- | --- |
| `/tools/encoding` | Encoding Converter | `enc_base64_encode`, `enc_base64_decode`, `enc_url_encode`, `enc_url_decode`, `enc_html_encode`, `enc_html_decode`, `enc_hex_encode`, `enc_hex_decode` |
| `/tools/data` | Data Transformer | `data_json_to_csv`, `data_csv_to_json`, `data_json_to_xml`, `data_xml_to_json`, `data_flatten`, `data_unflatten`, `data_key_rename` |
| `/tools/color` | Color Converter | `color_convert` |
| `/tools/time` | Time Converter | `time_convert`, `time_relative`, `time_duration` |
| `/tools/timezone` | Timezone Converter | `tz_convert`, `tz_list` |
| `/tools/math` | Math & Unit Converter | `math_eval`, `math_unit` |

### Analyzer Tools

| Path | Page | Operations Used |
| --- | --- | --- |
| `/tools/diff` | Smart Diff | `diff_compare` |
| `/tools/regex` | Regex Tester | `regex_test`, `regex_explain` |
| `/tools/cron` | Cron Builder | `cron_parse`, `cron_next` |
| `/tools/jwt` | JWT Debugger | `jwt_decode`, `jwt_verify` |
| `/tools/string` | String Tools | `str_case`, `str_diff`, `str_stats`, `str_sort_unique`, `str_replace` |
| `/tools/url` | URL Parser | `url_parse` |
| `/tools/headers` | HTTP Header Analyzer | `headers_analyze` |
| `/tools/dns` | DNS Lookup | `dns_lookup` |
| `/tools/http` | HTTP Client | `http_request` |
| `/tools/ip` | IP Calculator | `ip_calc` |

### DevOps Tools

| Path | Page | Operations Used |
| --- | --- | --- |
| `/tools/git` | Git Tools | `git_patch`, `git_commit_format`, `git_ignore_gen` |
| `/tools/dockerfile` | Dockerfile Linter | `dockerfile_lint` |
| `/tools/env` | Env File Tools | `env_parse`, `env_diff` |
| `/tools/k8s` | K8s Validator | `k8s_validate` |

---

## Command Palette (Cmd+K)

A global keyboard shortcut `Cmd+K` (or `Ctrl+K`) opens a command palette
overlay — fuzzy-searchable list of all 31 tools and 75 operations.

### Behavior

- Typing filters by tool name, operation name, or description.
- `Enter` navigates to the tool page.
- `Shift+Enter` opens the tool page with the selected operation pre-focused.
- Recent tools appear at the top as "Recently Used" section.
- Each result shows: icon + tool name + operation + keyboard shortcut
  (if assigned).

### Data Source

The palette fetches from `GET /api/v1/operations` on first open, caches
in memory. Each result maps to the routing table above.

---

## Home Dashboard (`/`)

The home page is a **searchable tool grid** — the primary discovery and
quick-access surface.

### Layout

```text
┌────────────────────────────────────────────────────────┐
│  Hero: "75 tools. One forge." + search bar             │
│  [Search tools...]                                     │
├────────────────────────────────────────────────────────┤
│  ★ Favourites (user-pinned, stored in localStorage)    │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐                  │
│  │ UUID │ │ JSON │ │ JWT  │ │ Diff │                   │
│  └──────┘ └──────┘ └──────┘ └──────┘                  │
├────────────────────────────────────────────────────────┤
│  [Generators] [Formatters] [Converters] [Analyzers]    │
│  [DevOps]     [All]                                    │
│                                                        │
│  ┌─────────────────┐  ┌─────────────────┐             │
│  │  🔑 UUID        │  │  {} JSON        │             │
│  │  Generate IDs   │  │  Format & lint  │             │
│  │  2 ops          │  │  2 ops          │             │
│  └─────────────────┘  └─────────────────┘             │
│  ┌─────────────────┐  ┌─────────────────┐             │
│  │  📊 CSV         │  │  🔍 Regex       │             │
│  │  Format & lint  │  │  Test & explain │             │
│  │  2 ops          │  │  2 ops          │             │
│  └─────────────────┘  └─────────────────┘             │
│  ...                                                   │
└────────────────────────────────────────────────────────┘
```

### Tool Card

Each card shows:

- **Icon** (from Lucide, colored per category).
- **Tool name** (e.g. "UUID Generator").
- **One-line description** (e.g. "Generate v4/v7 UUIDs and compute hashes").
- **Operation count** badge (e.g. "2 ops").
- **Favourite star** (toggle, persisted in `localStorage`).

Cards use `--bg-surface` with `--shadow-sm`; on hover they elevate to
`--shadow-md` with a `2px` left border in the category accent color.

### Category Filters

Horizontal chip row above the grid. Clicking a chip filters the grid.
"All" shows everything. Active chip uses `--sky-aqua` fill with white text.

---

## Breadcrumbs

Every tool page shows a breadcrumb trail:

```text
Home  >  Analyzers  >  Regex Tester
```

Breadcrumbs use `--text-secondary` with `--text-link` for the active
segment. Separator is a `/` or `>` glyph.

---

## URL Query State

Tool pages persist input state in the URL query string so users can share
links:

```text
/tools/uuid?version=7&count=5
/tools/diff?mode=json
/tools/cron?expr=*/5+*+*+*+*&tz=America/New_York
/tools/encoding?tab=base64&dir=encode
```

This allows bookmarking, sharing, and browser back/forward navigation
without losing inputs.

---

## API Integration Pattern

Every tool page follows the same data-flow pattern:

```text
User fills form  →  call POST /api/v1/<tool>/<op>  →  render result
                                                    →  show diagnostics
                                                    →  push query state
```

The shared `callOp<I,O>(tool, op, body)` function in
`web/src/api/client.ts` handles all calls. New tool pages only need to
define their form and render their result — no new API wiring required.
