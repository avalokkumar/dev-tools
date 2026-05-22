# DevForge — Frontend Specification

> Comprehensive design and UX specifications for the DevForge Web Surface,
> covering all 75 operations across 31 tools, plus CLI reference and MCP
> configuration pages.

---

## Document Index

| Doc | Title | Contents |
| --- | --- | --- |
| [00](./00-design-system.md) | Design System | Color palette, typography, spacing, iconography, motion, a11y tokens |
| [01](./01-information-architecture.md) | Information Architecture | Sitemap, navigation model, routing table, Cmd+K palette, home dashboard |
| [02](./02-component-library.md) | Component Library | Tech stack, layout shells, input/output/action/feedback components, hooks |
| [03](./03-page-specs-generators.md) | Pages: Generators | UUID, Faker, ID (ULID/Slug), TOTP, Crypto Toolkit |
| [04](./04-page-specs-formatters.md) | Pages: Formatters | JSON, YAML, CSV, SQL, Code (Go/XML/HTML), Markdown |
| [05](./05-page-specs-converters.md) | Pages: Converters | Encoding, Data Transform, Color, Time, Timezone, Math/Unit |
| [06](./06-page-specs-analyzers.md) | Pages: Analyzers | Diff, Regex, Cron, JWT, String, URL, Headers, DNS, HTTP, IP |
| [07](./07-page-specs-devops.md) | Pages: DevOps | Git (Patch/Commit/Ignore), Dockerfile, Env, K8s |
| [08](./08-surfaces-cli-mcp.md) | Surfaces: CLI + MCP | CLI Reference page, MCP Config Wizard, Settings |
| [09](./09-interactions-and-ux.md) | Interactions & UX | Animations, feedback, keyboard shortcuts, a11y, responsive, errors |

---

## Color Theme (quick reference)

```css
--shadow-grey:   #212738;  /* dark bg, nav, code blocks         */
--vibrant-coral: #F97068;  /* primary action, CTA, errors       */
--lemon-lime:    #D1D646;  /* success, additions, highlights     */
--platinum:      #EDF2EF;  /* page bg, light surfaces            */
--sky-aqua:      #57C4E5;  /* links, info, focus rings, selected */
```

## Tech Stack (quick reference)

- **React 18** + React Router 6 + TypeScript 5
- **Tailwind CSS 3** for utility-first styling
- **shadcn/ui** headless primitives
- **CodeMirror 6** for syntax-highlighted editors
- **Lucide React** icons
- **cmdk** for command palette
- **Sonner** for toast notifications

## Coverage

These specs cover the MVP feature set:

- **31 tool pages** spanning 75 backend operations
- **5 tool categories:** Generators, Formatters, Converters, Analyzers, DevOps
- **3 surface pages:** Home Dashboard, CLI Reference, MCP Config Wizard
- **1 settings page:** Theme, telemetry, plugins, version info
- **Shared infrastructure:** Design tokens, component library, hooks, feedback patterns
