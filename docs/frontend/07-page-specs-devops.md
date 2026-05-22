# 07 — Page Specs: DevOps

> Detailed page-level specifications for all DevOps tools:
> Git (Patch/Commit/Ignore), Dockerfile Linter, Env, and K8s Validator.

---

## Git Tools (`/tools/git`)

- **Icon:** `GitBranch` (Lucide)
- **Operations:** `git_patch`, `git_commit_format`, `git_ignore_gen`
- **Tabs:** "Patch", "Commit Format", ".gitignore Generator"

### Tab: Patch — Unified Diff Generator

```text
┌──────────────────────────────────────────────────────────┐
│  🔗 Git Tools  >  Patch                                  │
│                                                          │
│  [Patch]  [Commit Format]  [.gitignore Generator]        │
├─────────────────────────┬────────────────────────────────┤
│  Left (original)        │  Right (modified)              │
│  ┌───────────────────┐  │  ┌──────────────────────────┐  │
│  │ func hello() {    │  │  │ func hello() {           │  │
│  │   fmt.Println("hi")│  │  │   fmt.Println("hello")  │  │
│  │ }                 │  │  │ }                        │  │
│  └───────────────────┘  │  └──────────────────────────┘  │
│                         │                                │
│  Left Path  [a/main.go] │  Right Path [b/main.go]       │
│  Context Lines  [ 3 ]   │                                │
│                         │                                │
│  [ ▶ Generate Patch ]   │                                │
├─────────────────────────┴────────────────────────────────┤
│  Unified Diff Output                        📋  ⬇       │
│  ┌────────────────────────────────────────────────────┐  │
│  │ --- a/main.go                                      │  │
│  │ +++ b/main.go                                      │  │
│  │ @@ -1,3 +1,3 @@                                    │  │
│  │  func hello() {                                    │  │
│  │ -  fmt.Println("hi")                               │  │
│  │ +  fmt.Println("hello")                            │  │
│  │  }                                                 │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

- Output uses standard diff coloring: red (`--vibrant-coral`) for
  removals, green (`--lemon-lime`) for additions.
- Download as `.patch` file.

### Tab: Commit Format — Conventional Commit Validator

```text
│  Commit Message:                                         │
│  ┌────────────────────────────────────────────────────┐  │
│  │ feat(auth): add OAuth2 login flow                  │  │
│  │                                                    │  │
│  │ Implements GitHub OAuth for user authentication.   │  │
│  │                                                    │  │
│  │ BREAKING CHANGE: removes session-cookie auth       │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Validate ]                                          │
│                                                          │
│  ┌──── Parsed ────────────────────────────────────────┐  │
│  │  Type:     feat                                    │  │
│  │  Scope:    auth                                    │  │
│  │  Subject:  add OAuth2 login flow                   │  │
│  │  Body:     Implements GitHub OAuth for...          │  │
│  │  Footer:   BREAKING CHANGE: removes session-…      │  │
│  │  Breaking: ✅ Yes                                   │  │
│  │  Valid:    ✅ Conventional Commit compliant          │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Diagnostics: (if any validation issues)                 │
```

- Parsed fields shown in a structured card.
- Valid/invalid badge: green check or red X.
- Diagnostics render specific issues (e.g. "subject too long",
  "missing type prefix").

### Tab: .gitignore Generator

```text
│  Templates (search or select):                           │
│  ┌────────────────────────────────────────────────────┐  │
│  │  [Search templates...]                             │  │
│  │                                                    │  │
│  │  ☑ Go        ☑ Node     ☐ Python                  │  │
│  │  ☐ Rust      ☐ Java     ☐ Ruby                    │  │
│  │  ☑ macOS     ☐ Linux    ☐ Windows                 │  │
│  │  ☐ JetBrains ☑ VS Code  ☐ Vim                     │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Generate ]                                          │
│                                                          │
│  .gitignore Output                              📋  ⬇  │
│  ┌────────────────────────────────────────────────────┐  │
│  │ # Go                                               │  │
│  │ *.exe                                              │  │
│  │ /bin/                                              │  │
│  │ /vendor/                                           │  │
│  │                                                    │  │
│  │ # Node                                             │  │
│  │ node_modules/                                      │  │
│  │ dist/                                              │  │
│  │ ...                                                │  │
│  └────────────────────────────────────────────────────┘  │
```

- Searchable checkbox grid of template names.
- Popular templates shown first (Go, Node, Python, Rust).
- Download as `.gitignore` file.

---

## Dockerfile Linter (`/tools/dockerfile`)

- **Icon:** `Container` (Lucide)
- **Operations:** `dockerfile_lint`

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  🐳 Dockerfile Linter                                    │
│                                                          │
│  Dockerfile Input                                        │
│  ┌────────────────────────────────────────────────────┐  │
│  │ FROM node:18-alpine                                │  │
│  │ RUN npm install -g pnpm                            │  │
│  │ COPY . /app                                        │  │
│  │ RUN cd /app && pnpm install                        │  │
│  │ EXPOSE 3000                                        │  │
│  │ CMD ["node", "server.js"]                          │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Lint ]                                              │
│                                                          │
│  Diagnostics:                                            │
│  ┌────────────────────────────────────────────────────┐  │
│  │  ⚠️ [DOCKER.LINT.PIN_VERSION] Line 1: Consider     │  │
│  │     pinning to a specific digest (node:18-alpine@…)│  │
│  │  ⚠️ [DOCKER.LINT.MERGE_RUN] Lines 2,4: Consider   │  │
│  │     combining RUN commands to reduce layers         │  │
│  │  ℹ️ [DOCKER.LINT.COPY_LAST] Line 3: COPY before    │  │
│  │     dependency install — cache invalidation risk    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Summary: 0 errors, 2 warnings, 1 info                   │
└──────────────────────────────────────────────────────────┘
```

- Diagnostics inline in the code editor (wavy underlines on relevant
  lines) + listed below.
- Click a diagnostic to scroll to the relevant line in the editor.
- Summary badge: errors (coral), warnings (amber), info (aqua).

---

## Env File Tools (`/tools/env`)

- **Icon:** `ClipboardList` (Lucide)
- **Operations:** `env_parse`, `env_diff`
- **Tabs:** "Parse", "Diff"

### Tab: Parse

```text
│  .env Input                                              │
│  ┌────────────────────────────────────────────────────┐  │
│  │ DATABASE_URL=postgres://localhost:5432/mydb         │  │
│  │ API_KEY=sk-1234567890                               │  │
│  │ DEBUG=true                                          │  │
│  │ # This is a comment                                 │  │
│  │ API_KEY=override-value                              │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ☐ Allow export prefix                                   │
│                                                          │
│  [ ▶ Parse ]                                             │
│                                                          │
│  Parsed Values:                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  DATABASE_URL = postgres://localhost:5432/mydb 📋  │  │
│  │  API_KEY      = override-value                📋  │  │
│  │  DEBUG        = true                          📋  │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ⚠️ Duplicates: API_KEY (defined on lines 2 and 5)       │
```

### Tab: Diff

```text
│  Left (.env.example)    │  Right (.env.local)            │
│  ┌──────────────────┐   │  ┌──────────────────────────┐  │
│  │ DB_HOST=localhost │   │  │ DB_HOST=prod.db.com      │  │
│  │ API_KEY=          │   │  │ API_KEY=sk-real-key      │  │
│  │ DEBUG=false       │   │  │ NEW_VAR=hello            │  │
│  └──────────────────┘   │  └──────────────────────────┘  │
│                                                          │
│  [ ▶ Diff ]                                              │
│                                                          │
│  ⊕ Added:   NEW_VAR = hello                              │
│  ⊙ Changed: DB_HOST  localhost → prod.db.com             │
│  ⊙ Changed: API_KEY  (empty) → sk-real-key               │
│  ⊖ Removed: DEBUG = false                                │
```

Uses the same `DiffView` component with add/remove/change styling.

---

## K8s Validator (`/tools/k8s`)

- **Icon:** `Boxes` (Lucide)
- **Operations:** `k8s_validate`

### Wireframe

```text
┌──────────────────────────────────────────────────────────┐
│  ☸️ Kubernetes Manifest Validator                        │
│                                                          │
│  YAML Input                                              │
│  ┌────────────────────────────────────────────────────┐  │
│  │ apiVersion: apps/v1                                │  │
│  │ kind: Deployment                                   │  │
│  │ metadata:                                          │  │
│  │   name: my-app                                     │  │
│  │ spec:                                              │  │
│  │   replicas: 3                                      │  │
│  │   selector:                                        │  │
│  │     matchLabels:                                   │  │
│  │       app: my-app                                  │  │
│  │   template:                                        │  │
│  │     metadata:                                      │  │
│  │       labels:                                      │  │
│  │         app: my-app                                │  │
│  │     spec:                                          │  │
│  │       containers:                                  │  │
│  │       - name: app                                  │  │
│  │         image: my-app:latest                       │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [ ▶ Validate ]                                          │
│                                                          │
│  ┌──── Result ────────────────────────────────────────┐  │
│  │  ✅ Valid Kubernetes manifest                       │  │
│  │  Kind: Deployment    API: apps/v1                  │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Diagnostics: (if any)                                   │
│  ⚠️ [K8S.LATEST_TAG] "latest" tag is not recommended     │
│     for production deployments                           │
└──────────────────────────────────────────────────────────┘
```

- Input uses CodeEditor in YAML mode.
- Result shows Kind + API version extracted from the manifest.
- Diagnostics surface best-practice warnings.
