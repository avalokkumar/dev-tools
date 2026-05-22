# 02 — Component Library

> Reusable UI components shared across all tool pages. Each component is
> specified with props, variants, and visual behavior.

---

## Tech Stack

| Layer | Choice | Rationale |
| --- | --- | --- |
| Framework | React 18 | Already in codebase, SPA embedded in binary |
| Routing | React Router 6 | Already in codebase |
| Styling | Tailwind CSS 3 | Utility-first, design-token integration via `tailwind.config` |
| Components | shadcn/ui | Headless primitives, copy-paste into `web/src/components/ui/` |
| Icons | Lucide React | MIT, consistent, 1000+ icons |
| Code Editor | CodeMirror 6 | Syntax highlighting for JSON/YAML/SQL/Go/XML/HTML/Markdown |
| Toasts | Sonner | Lightweight notification toasts |
| Command Palette | cmdk | `Cmd+K` fuzzy search, shadcn-compatible |

---

## Layout Components

### `<AppShell>`

Top-level layout wrapper.

```text
Props:
  children: ReactNode

Structure:
  <div class="flex h-screen">
    <Sidebar />
    <div class="flex-1 flex flex-col">
      <Header />
      <main class="flex-1 overflow-y-auto p-6 bg-[--bg-primary]">
        {children}
      </main>
    </div>
  </div>
```

### `<Sidebar>`

Left navigation rail.

- **Width:** `260px` expanded, `64px` collapsed (icon-only).
- **Background:** `--shadow-grey`.
- **Text:** `--platinum` for labels, `--sky-aqua` for active item.
- **Active indicator:** `3px` left border in `--vibrant-coral`.
- **Sections:** Collapsible groups (Generators, Formatters, etc.) with
  chevron toggles.
- **Footer:** Version badge + collapse toggle button.

```text
Props:
  collapsed: boolean
  onToggleCollapse: () => void
  activePath: string
```

### `<Header>`

Top bar spanning the main content area.

- **Left:** Breadcrumbs.
- **Center:** Cmd+K search trigger button (pill shape, muted text
  "Search tools… ⌘K").
- **Right:** Theme toggle (sun/moon), settings gear icon.
- **Height:** `56px`.
- **Border:** `1px solid --border-default` bottom.

### `<Breadcrumb>`

```text
Props:
  items: { label: string, href?: string }[]

Renders:
  Home  /  Analyzers  /  Regex Tester
         ↑ link        ↑ link       ↑ current (no link)
```

---

## Tool Page Layout Components

### `<ToolPage>`

Standard wrapper for every tool page.

```text
Props:
  title: string           // "UUID Generator"
  description: string     // "Generate v4/v7 UUIDs and compute hashes"
  icon: LucideIcon        // Key
  category: string        // "Generators"
  operations: string[]    // ["uuid_generate", "uuid_hash"]
  children: ReactNode

Structure:
  <div>
    <ToolPageHeader title icon description />
    <OperationTabs operations>
      {children}
    </OperationTabs>
  </div>
```

### `<OperationTabs>`

When a tool page has multiple operations (e.g. `json_format` and
`json_validate`), display them as tabs.

- **Style:** Underline tabs in `--sky-aqua` for active, `--text-secondary`
  for inactive.
- **URL sync:** Active tab is stored in `?tab=format` query param.

```text
Props:
  tabs: { id: string, label: string, content: ReactNode }[]
  defaultTab?: string
```

### `<SplitPane>`

Horizontal or vertical resizable split for input/output views.

- **Default:** 50/50 horizontal split.
- **Drag handle:** `4px` bar in `--border-default`, cursor `col-resize`.
- **Collapse:** Double-click handle collapses one side.

```text
Props:
  direction: "horizontal" | "vertical"
  defaultSplit?: number  // 0.0 - 1.0, default 0.5
  left: ReactNode
  right: ReactNode
```

---

## Input Components

### `<CodeEditor>`

CodeMirror 6 wrapper for syntax-highlighted input/output.

- **Background:** `--bg-code` (shadow-grey).
- **Text:** `--platinum`, monospace font.
- **Line numbers:** Muted `--text-secondary`.
- **Gutter:** `36px` wide.
- **Languages:** JSON, YAML, SQL, Go, XML, HTML, Markdown, plaintext.
- **Features:** Bracket matching, auto-indent, search (`Cmd+F`),
  line wrapping toggle.

```text
Props:
  value: string
  onChange: (value: string) => void
  language: "json" | "yaml" | "sql" | "go" | "xml" | "html" | "markdown" | "text"
  readOnly?: boolean
  placeholder?: string
  minHeight?: string     // default "200px"
```

### `<InputField>`

Standard text input with label, description, and validation.

```text
Props:
  label: string
  description?: string
  value: string
  onChange: (v: string) => void
  type?: "text" | "number" | "password"
  placeholder?: string
  error?: string
  required?: boolean
```

- Uses shadcn `Input` internally.
- Error text in `--accent-error` below the field.
- Focus ring: `--shadow-glow`.

### `<SelectField>`

Dropdown select with label.

```text
Props:
  label: string
  options: { value: string, label: string }[]
  value: string
  onChange: (v: string) => void
```

### `<ToggleField>`

Boolean toggle (switch) with label.

```text
Props:
  label: string
  description?: string
  checked: boolean
  onChange: (v: boolean) => void
```

### `<TextAreaField>`

Multi-line text input (for non-code inputs like raw text, CSV).

```text
Props:
  label: string
  value: string
  onChange: (v: string) => void
  rows?: number
  monospace?: boolean
  placeholder?: string
```

### `<NumberStepper>`

Numeric input with +/- buttons and min/max bounds.

```text
Props:
  label: string
  value: number
  onChange: (v: number) => void
  min?: number
  max?: number
  step?: number
```

---

## Output Components

### `<OutputPanel>`

Displays operation results with copy-to-clipboard.

```text
Props:
  value: string
  language?: string       // for syntax highlighting
  label?: string          // e.g. "Output"
  loading?: boolean
  emptyMessage?: string   // "Run the tool to see results"

Features:
  - Copy button (top-right corner) — copies raw value.
  - Download button — saves as file with appropriate extension.
  - Word-wrap toggle.
  - Line count badge.
```

### `<DiagnosticsPanel>`

Renders `engine.Diagnostic[]` from any operation response.

```text
Props:
  diagnostics: { code: string, message: string, severity: 0|1|2 }[]

Rendering:
  severity 0 (Info)    → --accent-info    → ℹ️ icon
  severity 1 (Warning) → --accent-warning → ⚠️ icon
  severity 2 (Error)   → --accent-error   → ❌ icon

Each diagnostic: pill with icon + code + message.
Collapsible if > 3 items.
```

### `<ResultTable>`

Tabular result display (used by Faker, IP Calc, DNS, etc.).

```text
Props:
  columns: { key: string, label: string }[]
  rows: Record<string, unknown>[]
  onCopyRow?: (row) => void

Features:
  - Sortable columns (click header).
  - Zebra striping with --bg-elevated on odd rows.
  - Copy row as JSON.
```

### `<DiffView>`

Side-by-side or unified diff viewer (for `diff_compare`, `str_diff`,
`git_patch`, `env_diff`).

```text
Props:
  hunks: { path: string, op: string, left?: any, right?: any }[]
  mode: "side-by-side" | "unified"

Rendering:
  - Additions: --lemon-lime bg at 10% opacity + left border.
  - Removals: --vibrant-coral bg at 10% opacity + left border.
  - Changes: split color (coral left, lime right).
  - Path labels above each hunk.
```

---

## Action Components

### `<RunButton>`

Primary action button to execute an operation.

```text
Props:
  onClick: () => void
  loading?: boolean
  label?: string         // default "Run"
  shortcut?: string      // e.g. "⌘↵" displayed as badge

Style:
  bg: --vibrant-coral
  text: white
  hover: darken 10%
  loading: spinner replaces label
  border-radius: --radius-md
  padding: --space-3 --space-6
  font-weight: 600
```

### `<CopyButton>`

One-click copy to clipboard with feedback animation.

```text
Props:
  value: string
  label?: string        // default "Copy"

Behavior:
  Click → copy to clipboard → icon morphs from Copy to Check
  → revert after 2s. Toast "Copied to clipboard".
```

### `<ToolCardGrid>`

The home page tool card grid.

```text
Props:
  tools: ToolCard[]
  filter?: string          // category filter
  search?: string          // text search

ToolCard:
  {
    name: string
    description: string
    icon: LucideIcon
    category: string
    path: string
    opCount: number
    favourite: boolean
  }
```

---

## Feedback Components

### `<Toast>`

Uses Sonner. Positioned bottom-right.

- **Success:** `--lemon-lime` left border.
- **Error:** `--vibrant-coral` left border.
- **Info:** `--sky-aqua` left border.
- **Duration:** 3 seconds, dismissible on click.

### `<LoadingSkeleton>`

Placeholder shimmer for loading states. Uses `--bg-elevated` with a
sliding gradient animation.

### `<EmptyState>`

Shown when no results yet.

```text
Props:
  icon: LucideIcon
  title: string           // "No results yet"
  description: string     // "Fill in the form and click Run"
```

Centered in the output panel, muted icon at `48px`, text in
`--text-secondary`.

---

## Shared Hooks

### `useOperation<I, O>(tool, op)`

Encapsulates the API call lifecycle.

```typescript
const { run, result, loading, error, diagnostics } =
  useOperation<Input, Output>("uuid", "generate");
```

- `run(input: I)` — calls `callOp`, sets loading/result/error.
- `result: O | null` — last successful response.
- `loading: boolean` — spinner state.
- `error: string | null` — HTTP or parse error.
- `diagnostics: Diagnostic[]` — extracted from response.

### `useQueryState(key, defaultValue)`

Syncs a form field with the URL query string for shareable links.

```typescript
const [version, setVersion] = useQueryState("version", "4");
```

### `useFavourites()`

Manages the favourite tool list in `localStorage`.

```typescript
const { favourites, toggle, isFavourite } = useFavourites();
```

### `useCommandPalette()`

Provides the global `Cmd+K` registration and data.

```typescript
const { open, close, isOpen, search, results } = useCommandPalette();
```
