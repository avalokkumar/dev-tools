# 05 — Page Specs: Converters

> Detailed page-level specifications for all Converter tools:
> Encoding, Data Transform, Color, Time, Timezone, and Math/Unit.

---

## Encoding Converter (`/tools/encoding`)

- **Icon:** `Binary` (Lucide)
- **Operations:** 8 operations (4 encode/decode pairs)
- **Layout:** Tab-bar for codec, direction toggle per tab.

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  🔤 Encoding Converter                                   │
│                                                          │
│  [Base64]  [URL]  [HTML]  [Hex]                          │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  Direction:  ( ● Encode   ○ Decode )                     │
│                                                          │
│  Base64 Options:                                         │
│    ☐ URL-Safe alphabet                                   │
│    ☐ No padding                                          │
│                                                          │
│  Input                                                   │
│  ┌────────────────────────────────────────────────────┐  │
│  │ Hello, World!                                      │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Encode  ⌘↵ ]                                       │
│                                                          │
│  Output                                           📋    │
│  ┌────────────────────────────────────────────────────┐  │
│  │ SGVsbG8sIFdvcmxkIQ==                               │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

### Per-Tab Options

| Tab | Encode Options | Decode Options |
| --- | --- | --- |
| Base64 | `urlSafe`, `noPadding` | `urlSafe` |
| URL | `path` (path-encode vs query-encode) | — |
| HTML | — | — |
| Hex | `uppercase` | — |

**Interactions:**

- Switching direction flips input ↔ output content.
- "Swap" button (↔ icon) swaps input and output values and toggles
  direction — e.g. encode result becomes decode input.
- Live encoding on keypress (debounced 200ms) for small inputs (< 10KB).

---

## Data Transformer (`/tools/data`)

- **Icon:** `Shuffle` (Lucide)
- **Operations:** 7 operations
- **Layout:** Dropdown selector for transform type + split pane.

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  🔀 Data Transformer                                     │
│                                                          │
│  Transform:  [ JSON → CSV ▾ ]                            │
│                                                          │
│  Options:                                                │
│    json_to_csv: delimiter [,] flattenSeparator [.]       │
│    csv_to_json: header [☑] typedValues [☐]               │
│    json_to_xml: root [root] indent [2]                   │
│    xml_to_json: attrPrefix [@] textKey [#text]           │
│    flatten:     separator [.]                            │
│    unflatten:   separator [.]                            │
│    key_rename:  rules editor (see below)                 │
│                                                          │
├─────────────────────────┬────────────────────────────────┤
│  Input                  │  Output                        │
│  ┌───────────────────┐  │  ┌──────────────────────────┐  │
│  │ {"users":[        │  │  │ name,age                 │  │
│  │  {"name":"A",     │  │  │ A,30                     │  │
│  │   "age":30}]}     │  │  │                          │  │
│  └───────────────────┘  │  └──────────────────────────┘  │
│                         │                                │
│  [ ▶ Transform  ⌘↵ ]   │  [ Copy ] [ Download ]         │
└─────────────────────────┴────────────────────────────────┘
```

### Transform Picker

Dropdown lists all 7 operations as human-readable labels:

| Label | Operation |
| --- | --- |
| JSON → CSV | `data_json_to_csv` |
| CSV → JSON | `data_csv_to_json` |
| JSON → XML | `data_json_to_xml` |
| XML → JSON | `data_xml_to_json` |
| Flatten JSON | `data_flatten` |
| Unflatten JSON | `data_unflatten` |
| Rename Keys | `data_key_rename` |

### Key Rename — Rules Editor

For `data_key_rename`, show a dynamic rules list:

```text
│  Rules:                                                  │
│  ┌────────────────────────────────────────────────────┐  │
│  │  From: [old_name   ]  To: [new_name   ]  ☐ Regex  │  │
│  │  From: [email_addr ]  To: [email      ]  ☐ Regex  │  │
│  │  [ + Add Rule ]                                    │  │
│  └────────────────────────────────────────────────────┘  │
```

Each rule row: two text inputs + regex toggle + remove (×) button.
"Add Rule" appends a blank row.

---

## Color Converter (`/tools/color`)

- **Icon:** `Palette` (Lucide)
- **Operations:** `color_convert`
- **Layout:** Single-page with visual preview.

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  🎨 Color Converter                                      │
│                                                          │
│  Input   [ #F97068         ]    To: [ All ▾ ]            │
│                                                          │
│  [ ▶ Convert ]                                           │
│                                                          │
│  ┌─────────────────────────────────────────────────────┐ │
│  │                                                     │ │
│  │            ┌────────────────────┐                   │ │
│  │            │                    │                   │ │
│  │            │   Color Preview    │                   │ │
│  │            │   #F97068          │                   │ │
│  │            │                    │                   │ │
│  │            └────────────────────┘                   │ │
│  │                                                     │ │
│  │  HEX:  #F97068                               📋   │ │
│  │  RGB:  rgb(249, 112, 104)                     📋   │ │
│  │  HSL:  hsl(3, 92%, 69%)                       📋   │ │
│  │  R: 249  G: 112  B: 104                            │ │
│  │  H: 3°   S: 92%  L: 69%                            │ │
│  │                                                     │ │
│  │  ┌──── Contrast Check ─────────────────────────┐   │ │
│  │  │  On white (#FFF): 3.8:1  ⚠️ AA-large only   │   │ │
│  │  │  On black (#000): 5.5:1  ✅ AA               │   │ │
│  │  └─────────────────────────────────────────────┘   │ │
│  └─────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Large color swatch preview (120×120px rounded square).
- Each format row has a copy button.
- **Contrast checker** shows WCAG ratio against white and black.
- Input accepts any format: `#hex`, `rgb(…)`, `hsl(…)`.
- Live update on input change (debounced 200ms).

---

## Time Converter (`/tools/time`)

- **Icon:** `Timer` (Lucide)
- **Operations:** `time_convert`, `time_relative`, `time_duration`
- **Tabs:** "Convert", "Relative", "Duration"

### Tab: Convert

```text
│  Input      [ 1700000000                 ]               │
│  From       [ Unix Epoch ▾ ]                              │
│  To         [ ISO 8601 ▾ ]                                │
│  Timezone   [ UTC ▾ ]                                     │
│                                                          │
│  [ ▶ Convert ]                                           │
│                                                          │
│  Output: 2023-11-14T22:13:20Z                      📋   │
│  Parsed: Tuesday, November 14, 2023 10:13:20 PM UTC      │
```

From/To dropdowns: `epoch` (unix seconds), `epoch_ms` (milliseconds),
`iso` (ISO 8601), `rfc2822`, `rfc3339`, `human`.

### Tab: Relative

```text
│  From   [ 2024-01-01T00:00:00Z ]                        │
│  To     [ 2024-07-04T12:00:00Z ]                        │
│                                                          │
│  [ ▶ Calculate ]                                         │
│                                                          │
│  "6 months, 3 days, 12 hours ago"                        │
│  Total seconds: 16,005,600                               │
```

### Tab: Duration

```text
│  Input   [ 2h30m15s ]   or  [ PT2H30M15S ]              │
│                                                          │
│  [ ▶ Parse ]                                             │
│                                                          │
│  Hours: 2   Minutes: 30   Seconds: 15                    │
│  Total seconds: 9,015                                    │
```

---

## Timezone Converter (`/tools/timezone`)

- **Icon:** `Globe2` (Lucide)
- **Operations:** `tz_convert`, `tz_list`
- **Tabs:** "Convert" (default), "Zone List"

### Tab: Convert

```text
┌──────────────────────────────────────────────────────────┐
│  🌍 Timezone Converter                                   │
│                                                          │
│  Time       [ 2026-03-08T02:30:00 ]                     │
│  From Zone  [ America/New_York ▾ ]   (searchable)        │
│  To Zone    [ Europe/London ▾ ]      (searchable)        │
│                                                          │
│  [ ▶ Convert ]                                           │
│                                                          │
│  Result: 2026-03-08T07:30:00Z (Europe/London)      📋   │
│                                                          │
│  ⚠️ DST Note: America/New_York springs forward on        │
│  2026-03-08 at 02:00. The input time 02:30 falls in the │
│  gap and may be ambiguous.                               │
│                                                          │
│  ┌─── World Clock Preview ─────────────────────────┐    │
│  │  🇺🇸 New York    02:30 AM  EST (UTC-5)          │    │
│  │  🇬🇧 London       07:30 AM  GMT (UTC+0)          │    │
│  │  🇯🇵 Tokyo        04:30 PM  JST (UTC+9)          │    │
│  │  🇦🇺 Sydney       06:30 PM  AEDT (UTC+11)        │    │
│  └──────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Zone selectors are searchable comboboxes (type to filter).
- DST warnings render as `--accent-warning` banners.
- **World clock preview** shows the same instant in 4 popular zones
  (configurable in settings).

### Tab: Zone List

Fetches `tz_list` and renders a searchable table:

| Zone | UTC Offset | Abbreviation |
| --- | --- | --- |
| America/New_York | -05:00 | EST |
| Europe/London | +00:00 | GMT |
| … | … | … |

Filter input at top narrows results. Click a zone to auto-fill the
converter "From Zone" field.

---

## Math & Unit Converter (`/tools/math`)

- **Icon:** `Calculator` (Lucide)
- **Operations:** `math_eval`, `math_unit`
- **Tabs:** "Calculator", "Unit Converter"

### Tab: Calculator

```text
│  Expression   [ (2 + 3) * 4 / 2         ]               │
│                                                          │
│  [ ▶ Evaluate ]                                          │
│                                                          │
│  Result: 10                                        📋   │
```

- Supports: `+`, `-`, `*`, `/`, `%`, `^`, parentheses, math functions.
- Live evaluation on Enter key.
- History sidebar (last 10 calculations) stored in `sessionStorage`.

### Tab: Unit Converter

```text
│  Value   [ 1     ]                                       │
│  From    [ GiB ▾ ]                                        │
│  To      [ MiB ▾ ]                                        │
│                                                          │
│  [ ▶ Convert ]                                           │
│                                                          │
│  1 GiB = 1,024 MiB                                📋   │
```

Unit categories (grouped in dropdown):

- **Data:** B, KB, KiB, MB, MiB, GB, GiB, TB, TiB
- **Time:** ns, µs, ms, s, min, h, d
- **Throughput:** bps, Kbps, Mbps, Gbps
