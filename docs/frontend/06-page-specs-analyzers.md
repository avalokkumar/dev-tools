# 06 — Page Specs: Analyzers

> Detailed page-level specifications for all Analyzer tools:
> Diff, Regex, Cron, JWT, String, URL, HTTP Headers, DNS, HTTP Client,
> and IP Calculator.

---

## Smart Diff (`/tools/diff`)

- **Icon:** `GitCompareArrows` (Lucide)
- **Operations:** `diff_compare`
- **Layout:** Dual-pane input + diff result view.

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  ↔️ Smart Diff                                           │
│  Semantic diff for JSON, INI, or SQL inputs              │
│                                                          │
│  Mode: [ auto ▾ ] (auto, json, ini, sql)                 │
│  ☐ Ignore order   ☐ Ignore whitespace                    │
├─────────────────────────┬────────────────────────────────┤
│  Left                   │  Right                         │
│  ┌───────────────────┐  │  ┌──────────────────────────┐  │
│  │ {"a": 1,          │  │  │ {"a": 99,               │  │
│  │  "b": 2}          │  │  │  "c": 3}                │  │
│  └───────────────────┘  │  └──────────────────────────┘  │
│                         │                                │
│        [ ▶ Compare  ⌘↵ ]                                │
├──────────────────────────────────────────────────────────┤
│  Diff Results         [ Side-by-side ▾ ] [ Unified ]     │
│  ┌────────────────────────────────────────────────────┐  │
│  │  ⊖ remove  .b  value: 2                  (coral)  │  │
│  │  ⊕ add     .c  value: 3                  (lime)   │  │
│  │  ⊙ change  .a  1 → 99                    (aqua)   │  │
│  └────────────────────────────────────────────────────┘  │
│  Summary: 1 added, 1 removed, 1 changed                 │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Each hunk row color-coded: additions in `--lemon-lime` bg, removals
  in `--vibrant-coral` bg, changes in `--sky-aqua` bg.
- Toggle between side-by-side and unified views.
- JSON paths clickable — scrolls input editors to the relevant position.

---

## Regex Tester (`/tools/regex`)

- **Icon:** `Regex` (Lucide)
- **Operations:** `regex_test`, `regex_explain`
- **Tabs:** "Test" (default), "Explain"

### Tab: Test

```text
┌──────────────────────────────────────────────────────────┐
│  .* Regex Tester                                         │
│                                                          │
│  Pattern  [ \b[A-Z][a-z]+\b          ]                   │
│  Flags    [ i ☐ ] [ m ☐ ] [ s ☐ ]                       │
│                                                          │
│  Test Input                                              │
│  ┌────────────────────────────────────────────────────┐  │
│  │ Hello World, this is a Test String for regex.      │  │
│  │ Another Line with Words.                           │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Test  ⌘↵ ]                                         │
│                                                          │
│  Matches (4)                                             │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Match 1: "Hello"   pos: 0–5                      │  │
│  │  Match 2: "World"   pos: 6–11                     │  │
│  │  Match 3: "Test"    pos: 22–26                    │  │
│  │  Match 4: "String"  pos: 27–33                    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Highlighted Input (matches shown in --sky-aqua bg):     │
│  [Hello] [World], this is a [Test] [String] for regex.   │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- **Live highlight** in the input textarea: matched portions get
  `--sky-aqua` background overlay (debounced 200ms).
- Match list shows position ranges; clicking scrolls to that match.
- Group captures shown as nested items under each match.
- Flag checkboxes: `i` (case-insensitive), `m` (multiline), `s` (dotall).

### Tab: Explain

```text
│  Pattern  [ \b[A-Z][a-z]+\b ]                           │
│                                                          │
│  [ ▶ Explain ]                                           │
│                                                          │
│  Token-by-token explanation:                             │
│  ┌────────────────────────────────────────────────────┐  │
│  │  \b      Word boundary                            │  │
│  │  [A-Z]   One uppercase letter (A through Z)       │  │
│  │  [a-z]+  One or more lowercase letters            │  │
│  │  \b      Word boundary                            │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Plain English: "Match a word starting with an           │
│  uppercase letter followed by one or more lowercase      │
│  letters."                                               │
```

---

## Cron Builder (`/tools/cron`)

- **Icon:** `Clock` (Lucide)
- **Operations:** `cron_parse`, `cron_next`
- **Layout:** Single-page with expression builder and timeline.

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  ⏰ Cron Builder                                         │
│                                                          │
│  Expression  [ */5 * * * *               ]               │
│                                                          │
│  Visual Builder:                                         │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Minute:   [*/5 ▾]   (every 5 minutes)            │  │
│  │  Hour:     [* ▾]     (every hour)                  │  │
│  │  Day:      [* ▾]     (every day)                   │  │
│  │  Month:    [* ▾]     (every month)                 │  │
│  │  Weekday:  [* ▾]     (every day of week)           │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Timezone   [ UTC ▾ ]  (searchable)                      │
│  Flavor     [ unix ▾ ] (unix, quartz, aws)               │
│                                                          │
│  [ ▶ Parse + Next 5  ⌘↵ ]                               │
│                                                          │
│  Description: "Every 5 minutes"                          │
│                                                          │
│  Next 5 Runs:                                            │
│  ┌────────────────────────────────────────────────────┐  │
│  │  1. 2024-07-15 14:35:00 UTC   (in 3 minutes)     │  │
│  │  2. 2024-07-15 14:40:00 UTC   (in 8 minutes)     │  │
│  │  3. 2024-07-15 14:45:00 UTC   (in 13 minutes)    │  │
│  │  4. 2024-07-15 14:50:00 UTC   (in 18 minutes)    │  │
│  │  5. 2024-07-15 14:55:00 UTC   (in 23 minutes)    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ┌─── Timeline (24h) ────────────────────────────────┐  │
│  │ ╠═══╬═══╬═══╬═══╬═══╬═══╬═══╬═══╬═══╬═══╬═══╬══ │  │
│  │ 0   2   4   6   8  10  12  14  16  18  20  22    │  │
│  │ dots mark each run occurrence                     │  │
│  └───────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- **Visual builder** dropdowns update the expression text field and
  vice versa (two-way sync).
- **24-hour timeline** visualization: dots on a horizontal bar showing
  when the cron fires during the day. Dense firing (every minute) shows
  a solid bar.
- Relative time labels ("in 3 minutes") update live.
- Diagnostics from `cron_parse` render below if expression is invalid.

---

## JWT Debugger (`/tools/jwt`)

- **Icon:** `Ticket` (Lucide)
- **Operations:** `jwt_decode`, `jwt_verify`
- **Tabs:** "Decode" (default), "Verify"

### Tab: Decode

```text
┌──────────────────────────────────────────────────────────┐
│  🎫 JWT Debugger                                         │
│                                                          │
│  Token (paste JWT)                                       │
│  ┌────────────────────────────────────────────────────┐  │
│  │ eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI │  │
│  │ xMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijo │  │
│  │ xNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6y │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Decode ]                                            │
│                                                          │
│  ┌──── Header ─────────────┐  ┌──── Payload ──────────┐ │
│  │ {                       │  │ {                      │ │
│  │   "alg": "HS256",      │  │   "sub": "123456789", │ │
│  │   "typ": "JWT"         │  │   "name": "John Doe", │ │
│  │ }                       │  │   "iat": 1516239022   │ │
│  └─────────────────────────┘  │ }                      │ │
│                               └────────────────────────┘ │
│                                                          │
│  ┌──── Token Status ─────────────────────────────────┐  │
│  │  Algorithm: HS256                                  │  │
│  │  Issued At: 2018-01-18 01:30:22 UTC               │  │
│  │  Expiration: ⚠️ Token has no "exp" claim           │  │
│  │  Signature: ⬛ (not verified — use Verify tab)     │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Header and Payload rendered in syntax-highlighted JSON panes.
- **Token status** panel shows parsed claims:
  - `iat` → human-readable date.
  - `exp` → expiry with countdown ("expires in 2h 15m" or "expired 3d ago").
  - Expiry warning: `--accent-warning` if < 1h, `--vibrant-coral` if expired.
- Color-coded JWT segments: header (coral), payload (aqua), signature
  (lime) — matches the tri-color JWT.io style.

### Tab: Verify

```text
│  Token      [paste JWT here]                             │
│  Secret/Key [________________________]                   │
│  Key Format  ( ● HMAC secret  ○ PEM public key )         │
│  Allowed Algorithms  ☑ HS256  ☐ HS384  ☐ HS512  ☐ RS256 │
│  Leeway (sec) [ 0 ]                                      │
│                                                          │
│  [ ▶ Verify ]                                            │
│                                                          │
│  Result: ✅ Signature valid, token not expired            │
│  — or —                                                  │
│  Result: ❌ Signature invalid                             │
│  Result: ⚠️ Token expired 2 hours ago                    │
```

---

## String Tools (`/tools/string`)

- **Icon:** `Scissors` (Lucide)
- **Operations:** `str_case`, `str_diff`, `str_stats`, `str_sort_unique`,
  `str_replace`
- **Tabs:** "Case", "Diff", "Stats", "Sort", "Replace"

### Tab: Case

```text
│  Input   [ Hello World Example ]                         │
│  Mode    [ snake_case ▾ ]                                │
│                                                          │
│  Output: hello_world_example                       📋   │
```

Mode options: lower, upper, title, camel, pascal, snake, kebab, constant,
dot, path. **Live update** on input change.

### Tab: Diff

String-level diff (not semantic like `diff_compare`):

```text
│  Left    [ old text here ]                               │
│  Right   [ new text here ]                               │
│  ☐ Ignore whitespace   ☐ Ignore case                    │
│                                                          │
│  [ ▶ Diff ]                                              │
│  Hunks displayed in DiffView component.                  │
```

### Tab: Stats

```text
│  Input (paste or type text)                              │
│  ┌────────────────────────────────────────────────────┐  │
│  │ Lorem ipsum dolor sit amet...                      │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──── Statistics ────────────────────────────────────┐  │
│  │  Characters: 445      Bytes: 445                   │  │
│  │  Words: 68             Lines: 5                     │  │
│  │  Runes: 445           Unique Words: 52             │  │
│  └────────────────────────────────────────────────────┘  │
```

Live update on input change (no button needed).

### Tab: Sort

```text
│  Input (one item per line)                               │
│  ☐ Reverse   ☐ Numeric sort   ☐ Unique only             │
│                                                          │
│  [ ▶ Sort ]                                              │
│                                                          │
│  Output (sorted lines)                             📋   │
```

### Tab: Replace

```text
│  Input     [ The quick brown fox ]                       │
│  Pattern   [ fox                 ]                       │
│  Replace   [ cat                 ]                       │
│  ☐ Regex   ☐ Ignore case                                │
│                                                          │
│  [ ▶ Replace ]                                           │
│                                                          │
│  Output: The quick brown cat                       📋   │
│  Replacements: 1                                         │
```

---

## URL Parser (`/tools/url`)

- **Icon:** `Link` (Lucide)
- **Operations:** `url_parse`

### Wireframe

```text
│  URL   [ https://example.com:8080/path?q=1&r=2#frag ]   │
│                                                          │
│  [ ▶ Parse ]  (or live parse on input)                   │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Scheme:    https                            📋   │  │
│  │  Hostname:  example.com                      📋   │  │
│  │  Port:      8080                              📋   │  │
│  │  Path:      /path                             📋   │  │
│  │  Query:     q=1&r=2                           📋   │  │
│  │  Fragment:  frag                              📋   │  │
│  │                                                    │  │
│  │  Query Parameters:                                 │  │
│  │    q = 1                                           │  │
│  │    r = 2                                           │  │
│  └────────────────────────────────────────────────────┘  │
```

---

## HTTP Header Analyzer (`/tools/headers`)

- **Icon:** `Mail` (Lucide)
- **Operations:** `headers_analyze`

### Wireframe

```text
│  Headers (key: value, one per line)                      │
│  ┌────────────────────────────────────────────────────┐  │
│  │ Content-Type: text/html; charset=utf-8             │  │
│  │ X-Frame-Options: DENY                              │  │
│  │ Strict-Transport-Security: max-age=31536000        │  │
│  │ Cache-Control: no-store                            │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Analyze ]                                           │
│                                                          │
│  Findings:                                               │
│  ┌────────────────────────────────────────────────────┐  │
│  │  ✅ X-Frame-Options: DENY — clickjacking protected │  │
│  │  ✅ HSTS: max-age=31536000 — HTTPS enforced        │  │
│  │  ⚠️ Missing: Content-Security-Policy               │  │
│  │  ⚠️ Missing: X-Content-Type-Options                │  │
│  └────────────────────────────────────────────────────┘  │
```

Color-coded findings: green checks for present headers, warning badges
for missing security headers.

---

## DNS Lookup (`/tools/dns`)

- **Icon:** `Globe` (Lucide)
- **Operations:** `dns_lookup`

### Wireframe

```text
│  Host   [ example.com           ]                        │
│  Type   [ A ▾ ]  (A, AAAA, CNAME, MX, NS, TXT, PTR)    │
│  ☐ Allow Private IPs                                     │
│                                                          │
│  [ ▶ Lookup ]                                            │
│                                                          │
│  Records:                                                │
│  ┌────────────────────────────────────────────────────┐  │
│  │  A    93.184.216.34           TTL: 86400           │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Diagnostics:                                            │
│  (empty — or DNS.PRIVATE_BLOCKED if private IP found)    │
```

- "Allow Private IPs" toggle is OFF by default with a warning tooltip:
  "⚠️ Enabling this allows lookups that resolve to private/internal IPs."

---

## HTTP Client (`/tools/http`)

- **Icon:** `Radio` (Lucide)
- **Operations:** `http_request`

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  📡 HTTP Client                                          │
│                                                          │
│  Method  [ GET ▾ ]   URL [ https://api.example.com/data]│
│                                                          │
│  Headers:                                                │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Accept        : application/json           [ × ] │  │
│  │  Authorization : Bearer ...                  [ × ] │  │
│  │  [ + Add Header ]                                  │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Body (for POST/PUT/PATCH):                              │
│  ┌────────────────────────────────────────────────────┐  │
│  │ { "key": "value" }                                 │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Options:                                                │
│  ☐ Follow Redirects (max: [5])    Timeout: [30] sec     │
│  ☐ Allow Private IPs   ⚠️                               │
│                                                          │
│  [ ▶ Send Request  ⌘↵ ]                                 │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  Response                                                │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Status: 200 OK             Duration: 145ms        │  │
│  │                                                    │  │
│  │  Headers:                                          │  │
│  │    Content-Type: application/json                   │  │
│  │    Content-Length: 234                               │  │
│  │                                                    │  │
│  │  Body:                                    📋  ⬇   │  │
│  │  {                                                 │  │
│  │    "data": { ... }                                 │  │
│  │  }                                                 │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Diagnostics: (HTTP.PRIVATE_BLOCKED if blocked)          │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Method dropdown: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS.
- Dynamic header rows: key-value pairs with add/remove.
- Body editor: CodeEditor with JSON highlighting (only shown for methods
  with body).
- Response body: syntax-highlighted based on Content-Type.
- Status badge: green (2xx), yellow (3xx), red (4xx/5xx).
- Duration badge shows response time.
- "Allow Private IPs" toggle with prominent coral warning icon.

---

## IP Calculator (`/tools/ip`)

- **Icon:** `MapPin` (Lucide)
- **Operations:** `ip_calc`

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  📍 IP Calculator                                        │
│                                                          │
│  CIDR   [ 192.168.1.0/24        ]                       │
│  Max Host List  [ 256 ]  (0 = no listing)                │
│                                                          │
│  [ ▶ Calculate ]                                         │
│                                                          │
│  ┌──── Subnet Info ──────────────────────────────────┐  │
│  │  Network:      192.168.1.0                        │  │
│  │  Broadcast:    192.168.1.255                      │  │
│  │  First Host:   192.168.1.1                        │  │
│  │  Last Host:    192.168.1.254                      │  │
│  │  Netmask:      255.255.255.0                      │  │
│  │  Wildcard:     0.0.0.255                          │  │
│  │  Prefix:       /24                                │  │
│  │  Usable Hosts: 254                                │  │
│  └───────────────────────────────────────────────────┘  │
│                                                          │
│  Host List (254 hosts)                          📋      │
│  ┌────────────────────────────────────────────────────┐  │
│  │  192.168.1.1                                      │  │
│  │  192.168.1.2                                      │  │
│  │  ...                                              │  │
│  │  192.168.1.254                                    │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Subnet info rendered as a clean key-value grid.
- Host list in a scrollable monospace list with copy-all button.
- Visual subnet map (optional enhancement): a visual block showing
  the IP range as a filled bar.
