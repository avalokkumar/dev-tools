# 09 — Interactions & UX

> Animations, feedback patterns, accessibility requirements, responsive
> behavior, keyboard shortcuts, and error handling across the entire
> DevForge Web Surface.

---

## Micro-Interactions

### Button Press

- **Idle:** `--vibrant-coral` bg, white text, `--radius-md`.
- **Hover:** Darken bg 8% (`#E06058`), subtle scale `1.01`.
  Transition: `--duration-fast` with `--easing-default`.
- **Active (pressed):** Scale `0.98`, darken 12%.
- **Loading:** Button text replaced by a 16px spinner (white circle
  with coral arc). Button stays full width (no layout shift). Disabled
  state prevents double-click.
- **Success flash:** After a successful operation, the button briefly
  flashes `--lemon-lime` bg for 300ms before returning to coral.

### Copy-to-Clipboard

1. User clicks copy button (📋 icon).
2. Icon morphs to a checkmark (✓) with `--lemon-lime` color.
3. Toast slides in from bottom-right: "Copied to clipboard" with a
   `--lemon-lime` left border.
4. After 2 seconds, icon reverts to clipboard. Toast auto-dismisses.
5. If copy fails (e.g. Clipboard API unavailable): toast shows error
   in `--vibrant-coral`.

### Tab Switching

- Active tab underline slides horizontally to the new tab position.
  Uses CSS `transform: translateX()` with `--duration-normal`.
- Tab content fades in with `opacity: 0 → 1` over `--duration-normal`.
- No content jump — minimum height maintained during transition.

### Sidebar Navigation

- **Expand/collapse group:** Chevron rotates 90° with `--duration-fast`.
  Group items slide in/out with `max-height` transition.
- **Active item:** Left border slides in from `0 → 3px` in
  `--vibrant-coral`. Background fades to `rgba(87,196,229,0.1)`.
- **Sidebar collapse (icon rail):** Width transitions from `260px → 64px`
  over `--duration-normal`. Labels fade out at `50%` of the transition,
  icons remain centered.

### Command Palette (Cmd+K)

1. Overlay fades in with backdrop blur (`backdrop-filter: blur(4px)`).
2. Palette container scales from `0.95 → 1.0` with `--easing-spring`.
3. Results list renders immediately from cached operation data.
4. Typing filters results with instant highlight of matching characters.
5. `Escape` or clicking outside: reverse animation, `--duration-fast`.

### Page Transitions

- Route changes use a subtle fade: outgoing page `opacity: 1 → 0` over
  `100ms`, incoming page `opacity: 0 → 1` over `150ms`. No slide or
  scale — keep it fast and non-distracting.

---

## Feedback Patterns

### Operation Execution Flow

```text
1. User fills form
2. User clicks "Run" (or Cmd+Enter)
3. Button enters loading state (spinner)
4. Output panel shows <LoadingSkeleton> shimmer
5a. Success:
    - Spinner stops, button flashes green briefly
    - Output panel renders result with fade-in
    - Diagnostics panel updates (if any warnings)
    - Toast: "Operation completed" (only if > 1s elapsed)
5b. Error:
    - Spinner stops, button returns to idle
    - Error banner slides in above output: red border, error message
    - Diagnostics panel shows error-level items
    - Toast: "Operation failed" with error summary
```

### Diagnostic Severity Visual Language

| Severity | Icon | Color | Border | Background |
| --- | --- | --- | --- | --- |
| Error (2) | ❌ `XCircle` | `--vibrant-coral` | `2px left coral` | `rgba(249,112,104,0.08)` |
| Warning (1) | ⚠️ `AlertTriangle` | `--accent-warning` | `2px left amber` | `rgba(245,166,35,0.08)` |
| Info (0) | ℹ️ `Info` | `--sky-aqua` | `2px left aqua` | `rgba(87,196,229,0.08)` |

### Empty States

Every output panel has a meaningful empty state:

- **Before first run:** Ghost icon + "Fill in the form and press Run"
  in `--text-secondary`.
- **No results:** "No matches found" with suggestion text.
- **Error cleared:** Returns to the pre-run empty state.

### Toast Notifications

| Event | Type | Message | Duration |
| --- | --- | --- | --- |
| Copy success | Success | "Copied to clipboard" | 2s |
| Copy failure | Error | "Failed to copy — try selecting manually" | 4s |
| Download started | Info | "Downloading {filename}" | 2s |
| Operation success (slow) | Success | "{tool} completed in {ms}ms" | 3s |
| Operation error | Error | "{error message}" | 5s |
| Network error | Error | "Server unreachable — is devforge ui running?" | 5s |

---

## Keyboard Shortcuts

### Global

| Shortcut | Action |
| --- | --- |
| `Cmd+K` / `Ctrl+K` | Open command palette |
| `Cmd+/` / `Ctrl+/` | Toggle sidebar |
| `Escape` | Close modal / command palette / dropdown |

### Tool Pages

| Shortcut | Action |
| --- | --- |
| `Cmd+Enter` / `Ctrl+Enter` | Execute the current operation ("Run") |
| `Cmd+Shift+C` / `Ctrl+Shift+C` | Copy output to clipboard |
| `Cmd+Shift+D` / `Ctrl+Shift+D` | Download output |
| `Tab` | Move focus to next form field |
| `Shift+Tab` | Move focus to previous form field |

### Code Editor (CodeMirror)

| Shortcut | Action |
| --- | --- |
| `Cmd+F` / `Ctrl+F` | Find in editor |
| `Cmd+H` / `Ctrl+H` | Find and replace |
| `Cmd+A` / `Ctrl+A` | Select all |
| `Cmd+Z` / `Ctrl+Z` | Undo |
| `Cmd+Shift+Z` / `Ctrl+Shift+Z` | Redo |

A keyboard shortcut reference overlay is available via `Cmd+?` or the
`?` button in the header.

---

## Accessibility (a11y)

### WCAG 2.1 AA Compliance

- **Contrast ratios:** All text meets 4.5:1 minimum. Large text (≥18px
  or ≥14px bold) meets 3:1. See `00-design-system.md` for color
  adjustments.
- **Focus indicators:** Every interactive element shows a visible focus
  ring (`--shadow-glow`: 3px sky-aqua glow). Never `outline: none`
  without a replacement.
- **Focus order:** Logical tab order: sidebar → header → main content
  → output panel → diagnostics. No focus traps except modals (which
  trap focus intentionally until dismissed).

### Screen Reader Support

- All icons have `aria-label` attributes or are marked `aria-hidden`
  when paired with visible text.
- Tool pages use `<main>`, `<nav>`, `<section>`, `<header>` landmarks.
- Operation results announced with `aria-live="polite"` regions so
  screen readers announce new output.
- Error messages use `role="alert"` for immediate announcement.
- Form inputs linked to labels via `htmlFor`/`id` pairs.

### Keyboard Navigation

- All actions reachable via keyboard (no mouse-only interactions).
- Dropdown menus navigable with `ArrowUp`/`ArrowDown`, select with
  `Enter`, close with `Escape`.
- Tab groups use `role="tablist"`, `role="tab"`, `role="tabpanel"`
  with `aria-selected` and `arrow-key` navigation.
- Command palette results navigable with arrows, `Enter` to select.

### Motion Preferences

- All animations respect `prefers-reduced-motion: reduce`. When
  active:
  - Transitions reduce to `0ms` duration (instant state changes).
  - No sliding, scaling, or spring animations.
  - Loading spinners remain (essential feedback).
  - Toasts still appear but without slide animation.

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## Responsive Behavior

### Breakpoint Behaviors

| Element | `xs` (< 640px) | `sm` (640–768px) | `md` (768–1024px) | `lg+` (1024px+) |
| --- | --- | --- | --- | --- |
| Sidebar | Hidden (hamburger) | Icon rail (64px) | Full (260px) | Full (260px) |
| Tool grid | 1 column | 2 columns | 2 columns | 3 columns |
| Split pane | Stacked vertical | Stacked vertical | Side-by-side | Side-by-side |
| Header search | Icon only | Icon only | Pill button | Pill button |
| Code editor | Full width, 150px min-h | Full width | 50% split | 50% split |
| Breadcrumbs | Hidden | Abbreviated | Full | Full |

### Mobile-Specific Adaptations

- **Touch targets:** All buttons and interactive elements minimum
  `44×44px` tap area.
- **Hamburger menu:** Tapping opens the sidebar as a full-height overlay
  drawer with `backdrop-filter: blur(4px)`. Swipe-right to close.
- **Stacked layout:** Input and output sections stack vertically. A
  "Jump to output" floating button appears after running an operation
  (since output may be below the fold).
- **Bottom action bar:** On `xs`, the "Run" button moves to a sticky
  bottom bar so it's always reachable without scrolling.

---

## Error Handling

### Network Errors

When the API is unreachable (server stopped, network issue):

- A persistent banner appears at the top of the page:
  `⚠️ Cannot reach DevForge server. Is "devforge ui" running?`
- Banner uses `--accent-warning` background.
- All "Run" buttons disabled while disconnected.
- Auto-retry every 5 seconds; banner dismisses when connection restored.

### Validation Errors

- Client-side validation (empty required fields, out-of-range numbers)
  shows inline errors below the field in `--accent-error`.
- Form submission is blocked until all client-side validations pass.
- Server-side validation errors (from engine diagnostics) render in the
  `DiagnosticsPanel`.

### Rate Limiting & Timeouts

- Operations that exceed the 30s timeout show:
  `⏱️ Operation timed out after 30 seconds. Try reducing input size.`
- No client-side rate limiting for MVP (the backend handles it).

---

## Performance Considerations

- **Code splitting:** Each tool page is lazy-loaded via
  `React.lazy()` + `Suspense`. Only the home page and the active tool
  page are loaded initially.
- **Operation list caching:** `GET /api/v1/operations` response cached
  in memory for the session. Invalidated on page reload.
- **Debouncing:** All live-update features (regex highlighting, color
  preview, password strength) debounced at 200–300ms.
- **Large output virtualization:** For outputs > 1000 lines (e.g. Faker
  with 10,000 rows), use a virtualized list (`react-window`) to prevent
  DOM bloat.
- **Web Worker for diffing:** Heavy operations like `str_diff` on large
  texts run the API call in the main thread but rendering uses
  `requestIdleCallback` to avoid UI blocking.
