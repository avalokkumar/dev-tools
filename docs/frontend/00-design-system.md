# 00 — Design System

> DevForge brand identity, color palette, typography, spacing, iconography,
> and shared design tokens for the Web Surface.

---

## Color Palette

### Primary Tokens

| Token | Hex | HSL | Role |
| --- | --- | --- | --- |
| `--shadow-grey` | `#212738` | `hsla(224, 26%, 17%, 1)` | Primary background, nav bars, code editor bg, dark panels |
| `--vibrant-coral` | `#F97068` | `hsla(3, 92%, 69%, 1)` | Primary action, CTA buttons, error states, destructive warnings |
| `--lemon-lime` | `#D1D646` | `hsla(62, 64%, 56%, 1)` | Success indicators, "safe" badges, diff additions, highlights |
| `--platinum` | `#EDF2EF` | `hsla(144, 16%, 94%, 1)` | Page background, card surfaces, text on dark backgrounds |
| `--sky-aqua` | `#57C4E5` | `hsla(194, 73%, 62%, 1)` | Links, secondary actions, info badges, selected states |

### Extended Semantic Tokens (derived)

| Token | Value | Usage |
| --- | --- | --- |
| `--bg-primary` | `var(--platinum)` | Main page background |
| `--bg-surface` | `#FFFFFF` | Card / panel surfaces |
| `--bg-elevated` | `#F5F8F6` | Hover states on cards |
| `--bg-code` | `var(--shadow-grey)` | Code editor, terminal output, monospace blocks |
| `--bg-sidebar` | `var(--shadow-grey)` | Sidebar navigation |
| `--bg-input` | `#FFFFFF` | Form input backgrounds |
| `--text-primary` | `var(--shadow-grey)` | Body text on light backgrounds |
| `--text-secondary` | `#5A6178` | Muted labels, descriptions |
| `--text-inverse` | `var(--platinum)` | Text on dark backgrounds |
| `--text-link` | `var(--sky-aqua)` | Clickable links |
| `--border-default` | `#D4DAE0` | Card borders, dividers |
| `--border-focus` | `var(--sky-aqua)` | Focus rings on inputs |
| `--accent-error` | `var(--vibrant-coral)` | Error text, validation messages |
| `--accent-success` | `var(--lemon-lime)` | Success text, badges |
| `--accent-warning` | `#F5A623` | Warnings (amber derived) |
| `--accent-info` | `var(--sky-aqua)` | Info banners, tooltips |

### Dark Mode (future)

Dark mode inverts `--bg-primary` to `#1A1E2E` (darker than shadow-grey),
`--bg-surface` to `var(--shadow-grey)`, and text tokens to platinum. All
semantic tokens already use CSS custom properties, so toggling is a
single class swap on `<html>`.

---

## Typography

### Font Stack

```css
--font-sans:  'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
--font-mono:  'JetBrains Mono', 'Fira Code', 'SF Mono', ui-monospace, monospace;
--font-brand: 'Space Grotesk', var(--font-sans);
```

- **Inter** — Primary UI font. Load weights 400, 500, 600, 700 from
  Google Fonts. Excellent screen legibility at small sizes.
- **JetBrains Mono** — All code editors, terminal output, monospaced
  fields (JSON, regex, cron, UUIDs). Ligatures enabled.
- **Space Grotesk** — Brand headings only (logo wordmark, hero titles).

### Scale (rem, base 16px)

| Token | Size | Weight | Line-height | Usage |
| --- | --- | --- | --- | --- |
| `--text-h1` | 2rem (32px) | 700 | 1.2 | Page titles |
| `--text-h2` | 1.5rem (24px) | 600 | 1.3 | Section headings |
| `--text-h3` | 1.25rem (20px) | 600 | 1.3 | Card titles, tool names |
| `--text-body` | 0.9375rem (15px) | 400 | 1.6 | Body copy |
| `--text-small` | 0.8125rem (13px) | 400 | 1.5 | Captions, labels, badges |
| `--text-mono` | 0.875rem (14px) | 400 | 1.5 | Code, monospaced content |
| `--text-code-lg` | 0.9375rem (15px) | 400 | 1.6 | Editor panes |

---

## Spacing & Layout

### Spacing Scale (4px grid)

```css
--space-0:   0;
--space-1:   0.25rem  (4px)
--space-2:   0.5rem   (8px)
--space-3:   0.75rem  (12px)
--space-4:   1rem     (16px)
--space-5:   1.25rem  (20px)
--space-6:   1.5rem   (24px)
--space-8:   2rem     (32px)
--space-10:  2.5rem   (40px)
--space-12:  3rem     (48px)
--space-16:  4rem     (64px)
```

### Border Radius

| Token | Value | Usage |
| --- | --- | --- |
| `--radius-sm` | `4px` | Small chips, inline badges |
| `--radius-md` | `8px` | Cards, inputs, buttons |
| `--radius-lg` | `12px` | Modals, popovers |
| `--radius-xl` | `16px` | Hero cards, large panels |
| `--radius-full` | `9999px` | Pills, avatars |

### Shadow Elevation

| Token | Value | Usage |
| --- | --- | --- |
| `--shadow-sm` | `0 1px 2px rgba(33,39,56,0.06)` | Cards at rest |
| `--shadow-md` | `0 4px 12px rgba(33,39,56,0.10)` | Hover cards, dropdowns |
| `--shadow-lg` | `0 8px 24px rgba(33,39,56,0.14)` | Modals, popovers |
| `--shadow-glow` | `0 0 0 3px rgba(87,196,229,0.3)` | Focus ring glow |

### Layout Grid

- **Max content width:** `1440px` (centered with `auto` margins).
- **Sidebar width:** `260px` collapsed to `64px` icon-only on mobile.
- **Main content area:** fluid, minimum `640px` before sidebar collapses.
- **Gutter:** `--space-6` (24px) between columns.
- **Card grid:** CSS Grid `repeat(auto-fill, minmax(320px, 1fr))` for
  the home page tool grid.

---

## Iconography

### Icon Library

Use **Lucide React** (`lucide-react`) — MIT-licensed, consistent 24×24
stroke icons. Every tool category gets a canonical icon:

| Category | Icon | Lucide Name |
| --- | --- | --- |
| Generators | ⚙️ | `Cog` / `Wand2` |
| Formatters | 📐 | `AlignLeft` / `Code2` |
| Converters | 🔄 | `ArrowRightLeft` / `Repeat` |
| Analyzers | 🔍 | `Search` / `ScanSearch` |
| Security | 🔒 | `Lock` / `Shield` |
| DevOps | 🚀 | `Rocket` / `GitBranch` |
| Network | 🌐 | `Globe` / `Wifi` |

### Per-Tool Icons

| Tool | Icon | Lucide Name |
| --- | --- | --- |
| UUID | 🔑 | `Key` |
| JSON | `{}` | `Braces` |
| YAML | 📄 | `FileText` |
| CSV | 📊 | `Table` |
| Diff | ↔️ | `GitCompareArrows` |
| Regex | `.*` | `Regex` |
| Cron | ⏰ | `Clock` |
| JWT | 🎫 | `Ticket` |
| Timezone | 🌍 | `Globe2` |
| Faker | 🎭 | `Drama` |
| Encoding | 🔤 | `Binary` |
| String | ✂️ | `Scissors` |
| Time | ⏱️ | `Timer` |
| SQL | 🗃️ | `Database` |
| Markdown | ✍️ | `PenTool` |
| ID (ULID/Slug) | 🏷️ | `Tag` |
| Color | 🎨 | `Palette` |
| Data Transform | 🔀 | `Shuffle` |
| Git | 🔗 | `GitBranch` |
| Dockerfile | 🐳 | `Container` |
| Env | 📋 | `ClipboardList` |
| K8s | ☸️ | `Boxes` |
| URL | 🔗 | `Link` |
| HTTP Headers | 📨 | `Mail` |
| DNS | 🌐 | `Globe` |
| HTTP | 📡 | `Radio` |
| Crypto | 🔐 | `ShieldCheck` |
| TOTP | 🔢 | `Hash` |
| Code Format | 💻 | `Terminal` |
| Math | ➗ | `Calculator` |
| IP | 📍 | `MapPin` |

### Logo

The DevForge wordmark uses **Space Grotesk Bold** in `--shadow-grey`
with an anvil/forge icon rendered in a two-color mark (`--vibrant-coral`
for the hammer, `--sky-aqua` for the anvil). The SVG is embedded in the
header at `32px` height. Favicon is a `16×16` / `32×32` simplified anvil
in coral on transparent.

---

## Motion & Animation

| Token | Value | Usage |
| --- | --- | --- |
| `--duration-fast` | `120ms` | Hover color changes, button states |
| `--duration-normal` | `200ms` | Panel open/close, tab switches |
| `--duration-slow` | `350ms` | Page transitions, modal entrances |
| `--easing-default` | `cubic-bezier(0.4, 0, 0.2, 1)` | Standard easing |
| `--easing-spring` | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Bouncy micro-interactions |

---

## Responsive Breakpoints

| Name | Width | Behavior |
| --- | --- | --- |
| `xs` | `< 640px` | Sidebar hidden, single column, full-width cards |
| `sm` | `640–768px` | Sidebar collapsed to icon rail (64px) |
| `md` | `768–1024px` | Sidebar open, 2-column tool grid |
| `lg` | `1024–1440px` | Full layout, 3-column tool grid |
| `xl` | `> 1440px` | Content centered at max-width |

---

## Accessibility

- All interactive elements must have ≥ 4.5:1 contrast ratio (WCAG AA).
- `--shadow-grey` on `--platinum` = 11.2:1 ✅
- `--vibrant-coral` on `#FFFFFF` = 3.8:1 — use as accent only, never
  for small body text. For error text on white, darken to `#E05550`.
- `--sky-aqua` on `#FFFFFF` = 3.0:1 — darken link text to `#2A9BBF`
  for AA compliance on small text.
- Focus rings use `--shadow-glow` (3px sky-aqua glow) for visibility.
- All icons carry `aria-label` or are paired with visible text.
- Keyboard navigation: `Tab` through sidebar, then tool grid, then
  main content. `Escape` closes modals and dropdowns.
