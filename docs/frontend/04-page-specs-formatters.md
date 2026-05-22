# 04 — Page Specs: Formatters

> Detailed page-level specifications for all Formatter tools:
> JSON, YAML, CSV, SQL, Code, and Markdown.

---

## Shared Formatter Layout

All formatter pages share a common **split-pane** layout:

```text
┌──────────────────────────────────────────────────────────┐
│  [icon] Tool Name                                        │
│  [Format]  [Validate]  (tabs if multi-op)                │
├─────────────────────────┬────────────────────────────────┤
│  Input                  │  Output                        │
│  ┌───────────────────┐  │  ┌──────────────────────────┐  │
│  │                   │  │  │                          │  │
│  │  CodeEditor       │  │  │  CodeEditor (read-only)  │  │
│  │  (editable)       │  │  │                          │  │
│  │                   │  │  │                          │  │
│  └───────────────────┘  │  └──────────────────────────┘  │
│                         │                                │
│  Options bar:           │  [ Copy ] [ Download ] [ Wrap ]│
│  [indent ▾] [sort ☐]   │                                │
│                         │                                │
│  [ ▶ Format  ⌘↵ ]      │                                │
├─────────────────────────┴────────────────────────────────┤
│  Diagnostics (if any)                                    │
│  ⚠️ [JSON.TRAILING_COMMA] Trailing comma at line 4       │
└──────────────────────────────────────────────────────────┘
```

- Input and output are side-by-side on `md+`, stacked vertically on `xs`.
- Input CodeEditor has syntax highlighting for the tool's language.
- Output CodeEditor is read-only with the same language highlighting.
- DiagnosticsPanel below renders validation warnings/errors.

---

## JSON Formatter (`/tools/json`)

- **Icon:** `Braces` (Lucide)
- **Operations:** `json_format`, `json_validate`
- **Tabs:** "Format" (default), "Validate"

### Tab: Format — Options

| Control | Type | Default | API Field |
| --- | --- | --- | --- |
| Indent | Select: 2/4/Tab/Compact | 2 | `indent` (0=compact) |
| Sort Keys | Toggle | off | `sortKeys` |
| Trailing Newline | Toggle | off | `trailingNewline` |

### Tab: Validate

- Input: JSON text.
- Output: "✅ Valid JSON" or list of diagnostics with line numbers.
- No options needed.

**Interactions:**

- Paste detection: if clipboard contains JSON, auto-populate input.
- Format on `Cmd+Enter`.
- Error positions: clicking a diagnostic scrolls the input editor to
  the relevant line.

---

## YAML Formatter (`/tools/yaml`)

- **Icon:** `FileText` (Lucide)
- **Operations:** `yaml_format`, `yaml_validate`, `yaml_convert`
- **Tabs:** "Format", "Validate", "Convert"

### Tab: Format — Options

| Control | Type | Default | API Field |
| --- | --- | --- | --- |
| Indent | Select: 2/4 | 2 | `indent` |

### Tab: Convert

Converts between YAML and JSON.

| Control | Type | Default | API Field |
| --- | --- | --- | --- |
| Target | Select: JSON or YAML | JSON | `to` |
| Indent | Select: 2 or 4 | 2 | `indent` |

- Input language auto-detects: if input starts with `{`, set input
  editor to JSON mode; otherwise YAML mode.
- Output editor switches language to match the target format.

---

## CSV Formatter (`/tools/csv`)

- **Icon:** `Table` (Lucide)
- **Operations:** `csv_format`, `csv_validate`
- **Tabs:** "Format", "Validate"

### Tab: Format — Options

| Control | Type | Default | API Field |
| --- | --- | --- | --- |
| Delimiter | Select: `,` or `\t` or `;` or pipe | `,` | `delimiter` |
| Has Header | Toggle | on | `header` |
| Align Columns | Toggle | off | `alignColumns` |

### Tab: Format — Output Enhancements

- Show a **preview table** below the raw output: parsed CSV rendered as
  an HTML table with zebra rows and sortable headers.
- "Copy as Markdown" button converts the CSV to a GFM table.

---

## SQL Formatter (`/tools/sql`)

- **Icon:** `Database` (Lucide)
- **Operations:** `sql_format`, `sql_validate`
- **Tabs:** "Format", "Validate"

### Tab: Format — Options

| Control | Type | Default | API Field |
| --- | --- | --- | --- |
| Uppercase Keywords | Toggle | on | `uppercase` |
| Indent Width | Select: 2/4 | 2 | `indentWidth` |

### Tab: Validate

- Input: SQL text.
- Output: "✅ Valid SQL" or diagnostics with error details.

---

## Code Formatter (`/tools/code`)

- **Icon:** `Terminal` (Lucide)
- **Operations:** `code_fmt_go`, `code_fmt_xml`, `code_fmt_html`
- **Tabs:** "Go", "XML", "HTML"

### Layout

Each tab is a split-pane formatter with language-specific highlighting:

| Tab | Input Language | Output Language | Options |
| --- | --- | --- | --- |
| Go | `go` | `go` | None (uses `gofmt` rules) |
| XML | `xml` | `xml` | `indent` (2/4) |
| HTML | `html` | `html` | None (uses standard rules) |

### Go Tab — Special Feature

- Shows "gofmt-clean" badge if output equals input (no changes needed).
- Diagnostics surface Go parse errors with line numbers.

---

## Markdown Editor (`/tools/markdown`)

- **Icon:** `PenTool` (Lucide)
- **Operations:** `md_to_html`, `md_table_from_csv`
- **Tabs:** "Preview", "CSV → Table"

### Tab: Preview

```text
┌─────────────────────────┬────────────────────────────────┐
│  Markdown Input          │  HTML Preview                  │
│  ┌───────────────────┐  │  ┌──────────────────────────┐  │
│  │ # Hello World     │  │  │  Hello World              │  │
│  │                   │  │  │  ═══════════              │  │
│  │ This is **bold**  │  │  │  This is bold             │  │
│  │ and _italic_.     │  │  │  and italic.              │  │
│  │                   │  │  │                          │  │
│  │ ```js             │  │  │  ┌─────────────────────┐ │  │
│  │ console.log("hi") │  │  │  │ console.log("hi")   │ │  │
│  │ ```               │  │  │  └─────────────────────┘ │  │
│  └───────────────────┘  │  └──────────────────────────┘  │
│                         │                                │
│  Options:               │  [ Copy HTML ] [ Download ]    │
│  ☑ GFM (GitHub Flavored)│                                │
│  ☐ Allow Unsafe HTML    │                                │
└─────────────────────────┴────────────────────────────────┘
```

**Interactions:**

- **Live preview:** HTML panel updates as user types (debounced 300ms).
- Left pane: CodeEditor in `markdown` mode.
- Right pane: rendered HTML in a sandboxed `<iframe>` or `dangerouslySetInnerHTML`
  with sanitized output (the API runs bluemonday by default).
- "Allow Unsafe HTML" toggle: warns with a coral banner "⚠️ Unsafe HTML
  is enabled — output may contain scripts" before enabling.
- GFM toggle: enables table/strikethrough/tasklist extensions.

### Tab: CSV → Table

```text
│  CSV Input                                               │
│  ┌────────────────────────────────────────────────────┐  │
│  │ name,age,city                                      │  │
│  │ Alice,30,NYC                                       │  │
│  │ Bob,25,LA                                          │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Delimiter   [ , ▾ ]                                     │
│  Alignment   [ None ▾ ] (none / left / center / right)   │
│                                                          │
│  [ ▶ Convert ]                                           │
│                                                          │
│  GFM Table Output                                 📋    │
│  ┌────────────────────────────────────────────────────┐  │
│  │ | name  | age | city |                             │  │
│  │ | ----- | --- | ---- |                             │  │
│  │ | Alice | 30  | NYC  |                             │  │
│  │ | Bob   | 25  | LA   |                             │  │
│  └────────────────────────────────────────────────────┘  │
```
