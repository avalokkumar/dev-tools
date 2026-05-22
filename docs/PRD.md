# Product Requirement Document — DevForge: Unified Agentic Developer Toolkit

**Status:** Draft v1.0
**Owner:** Product / Engineering
**Date:** 2026-05-09

---

## 1. Product Overview

### 1.1 Problem Statement
Modern developers waste significant time context-switching between dozens of disconnected single-purpose web utilities (regex testers, JSON formatters, JWT decoders, cron builders, mock servers, log parsers, etc.). Beyond fragmentation, five high-value workflows remain unsolved by existing tools:

1. **API contract drift** — specs diverge from running code, causing production outages.
2. **Log observability** — manual log triage scales poorly; existing platforms (ELK, Splunk) are heavyweight and reactive.
3. **Secrets / env management** — plaintext `.env` files leak credentials; enterprise SecretOps platforms have steep onboarding.
4. **Knowledge fragmentation** — documentation scattered across Slack, Notion, GitHub; LLM agents hallucinate without grounded context.
5. **Test brittleness** — Playwright suites break on UI changes; existing self-healers patch at runtime without updating source.

Crucially, none of these tools expose deterministic, MCP-native interfaces that let AI coding agents (Claude Code, Cursor) operate against them without hallucination.

### 1.2 Target Users
- **Primary:** Backend, frontend, full-stack developers in small-to-mid teams (2–50 engineers).
- **Secondary:** DevOps / SRE engineers, QA automation engineers, indie hackers.
- **Tertiary:** AI coding agents (as first-class consumers via MCP).

### 1.3 Core Objective
Deliver a **single, offline-first, MCP-enabled developer toolkit** that:
- Bundles high-frequency utilities (formatters, converters, generators) in one cohesive UI + CLI.
- Ships five flagship "agentic platforms" that solve the episodic-pain workflows above.
- Exposes every capability as a deterministic MCP server so AI agents can use it without hallucination.
- Uses GitHub identity for zero-friction onboarding.

### 1.4 Success Metrics
- 10K weekly active developers within 6 months of MVP launch.
- ≥40% of users invoke ≥2 distinct tool categories per week (composability proof).
- ≥30% of agentic-platform users (API Drift Guard, Healer, etc.) reach a paid tier.
- MCP-tool-call success rate ≥95% (deterministic interface integrity).

---

## 2. Feature Set

The product is organized into **two layers**: (A) Flagship Agentic Platforms (the wedge), (B) Utility Toolkit (the daily-driver glue).

---

### 🔹 Feature A1: API Contract Drift Guard

- **Description:** GitHub App + CLI + MCP server that compares OpenAPI/GraphQL specs against running code and live traffic. Fails CI on breaking drift; auto-generates client types; exposes the contract to AI agents.
- **Objective:** Prevent shadow APIs and schema drift from reaching production. Reduce LLM hallucination of nonexistent endpoints.
- **Scope:**
  - **In-Scope:** OpenAPI 3.1/3.2 + GraphQL; static AST diff; live-traffic capture via sidecar proxy; PR comment with breaking-change report; auto client-type generation (TS, Python); MCP server exposing `getEndpoint`, `listEndpoints`, `validateRequest`.
  - **Out-of-Scope (MVP):** gRPC / AsyncAPI; full Postman replacement; production traffic mirroring at scale.

### 🔹 Feature A2: Log Analyzer Lite

- **Description:** Lightweight, zero-config log parser + anomaly detector. Auto-classifies common formats (Nginx, Apache, syslog, JSON), surfaces anomalies via two-stage LLM filter, suggests remediations grounded in repo history.
- **Objective:** Give small teams Splunk-grade insight at <10% of cost and setup time.
- **Scope:**
  - **In-Scope:** File ingestion + tail mode; Drain3-based template extraction; anomaly scoring with context (deploy events, time-of-day); RAG over historical logs; CLI dashboard + minimal web UI.
  - **Out-of-Scope:** Full distributed tracing; multi-region log aggregation; SIEM-grade compliance features.

### 🔹 Feature A3: GitHub-Native Secrets Gateway

- **Description:** CLI tool that authenticates via GitHub OAuth, stores secrets remotely (encrypted at rest, AES-256-GCM), and injects them into process memory at runtime. Drop-in replacement for `dotenv` with team sync and drift alerts.
- **Objective:** Eliminate plaintext `.env` files without forcing teams onto Vault/Doppler.
- **Scope:**
  - **In-Scope:** `devforge run -- <cmd>` memory-only injection; per-environment secret sets; GitHub team-based ACLs; schema validation (Zod/Pydantic adapter); local drift alerts; MCP server with scoped secret read.
  - **Out-of-Scope (MVP):** Hardware-token auth; multi-cloud KMS; secret rotation automation.

### 🔹 Feature A4: Knowledge Base Compiler

- **Description:** CLI that "compiles" a `raw/` folder of unstructured docs (transcripts, PRDs, ADRs, READMEs) into a structured, interlinked Markdown wiki with citations. Implements the Karpathy Loop: ingest → compile → query → health-check.
- **Objective:** Self-maintaining team knowledge substrate that grounds AI agents and humans alike.
- **Scope:**
  - **In-Scope:** Incremental compilation (only re-process changed files); local SQLite vector store; `query` command with citations; health-check linter (contradictions, stale links, orphan pages); Marp slide export.
  - **Out-of-Scope:** Real-time collaborative editing; Notion-style WYSIWYG; cloud-hosted multi-tenant wiki.

### 🔹 Feature A5: Playwright Healer Dashboard

- **Description:** Wraps Playwright runs, detects locator failures, proposes fixes via accessibility-tree analysis, and opens PRs that update **source code** with stable `data-testid` selectors. Visual regression guard prevents silent functional drift.
- **Objective:** Replace runtime "self-healing" with auditable source-patching healing.
- **Scope:**
  - **In-Scope:** Playwright runner wrapper; failure classification (UI-change / timing / data-regression); `fix-proposal.md` artifact per failure; one-click PR creation; visual snapshot diff; healer transparency report.
  - **Out-of-Scope:** Cypress / Selenium support; mobile native testing; cloud-hosted browser farm.

---

### 🔹 Feature B1: Developer Utility Toolkit (Bundled Micro-Tools)

- **Description:** Cohesive bundle of high-frequency utilities accessible via web UI, CLI (`devforge <tool>`), and MCP. Offline-first; plugin-extensible.
- **Objective:** Eliminate context-switching to 30+ single-purpose websites.
- **Scope:**
  - **In-Scope (MVP set):**
    - JSON / YAML / CSV formatter + validator (with schema hints)
    - Smart diff viewer (semantic JSON / SQL / config diffs)
    - Regex tester + explainer
    - Cron expression builder + visualizer
    - JWT debugger (decode / verify / expiry alerts)
    - UUID / hash generator (multiple versions / formats)
    - Timezone converter (DST-aware)
    - Dev-focused unit converter (bytes, latency, throughput)
    - SQL formatter + linter
    - Markdown editor + live preview
    - Data faker (multi-locale)
    - cURL ↔ code-snippet ↔ Postman converter
    - Webhook tester (request bin + replay)
  - **Out-of-Scope (MVP):** Image / PDF compression; color palette tools; clipboard manager; bulk file rename (deferred to Phase 2).

### 🔹 Feature B2: Unified Shell — UI / CLI / MCP Surface

- **Description:** Three coordinated entry points sharing one core engine. Web UI for discovery; CLI for power users / scripting; MCP server for agents.
- **Objective:** Maximize adoption (web), retention (CLI), and AI-agent reach (MCP).
- **Scope:**
  - **In-Scope:** Single-binary CLI (Go or Rust); local web app served from CLI (`devforge ui`); MCP server (`devforge mcp`); shared core library; offline-first; GitHub OAuth sign-in for cloud-sync features.
  - **Out-of-Scope:** Desktop installer apps (Electron, etc.); mobile clients.

---

## 3. Functional Breakdown

### A1 — API Contract Drift Guard
- **Key functionalities:** spec parsing (OpenAPI/GraphQL), AST extraction from code, live-traffic capture, drift detection engine, PR commenter, client-type generator, MCP endpoint server.
- **User interactions:** install GitHub App → add `devforge.yml` → on PR, see drift report comment → optionally apply auto-fix; CLI: `devforge api diff`, `devforge api gen-types`.
- **Dependencies:** GitHub App permissions; language-specific AST parsers (TS, Python, Go MVP); sidecar proxy for live traffic.

### A2 — Log Analyzer Lite
- **Key functionalities:** format auto-detection, template extraction (Drain3), anomaly scoring, RAG over historical logs, remediation suggester.
- **User interactions:** `devforge logs tail <file|url>`; web dashboard for trends; alert hooks (webhook / Slack).
- **Dependencies:** Drain3 (Python); optional embedding model (local via Ollama or OpenAI); Redis for state.

### A3 — Secrets Gateway
- **Key functionalities:** OAuth flow, encrypted remote storage, in-memory injection, schema validator, drift watcher.
- **User interactions:** `devforge auth login`, `devforge secret set/get`, `devforge run -- <cmd>`; web dashboard for team visibility.
- **Dependencies:** GitHub OAuth; backend KV store (encrypted); language adapters for Pydantic / Zod.

### A4 — Knowledge Base Compiler
- **Key functionalities:** incremental compiler, embedding indexer, RAG querier, health-check linter, Marp exporter.
- **User interactions:** `devforge kb compile`, `devforge kb query "<question>"`, `devforge kb lint`.
- **Dependencies:** SQLite + vector extension (sqlite-vec); markdown parser; LLM API (BYO key).

### A5 — Playwright Healer
- **Key functionalities:** test runner wrapper, failure classifier, locator proposer, visual diff, PR writer.
- **User interactions:** `devforge pw run`, `devforge pw heal`; dashboard reviews; auto-PR.
- **Dependencies:** Playwright API; MCP browser tools; image-diff library; GitHub API.

### B1 — Utility Toolkit
- **Key functionalities:** per-tool engines (formatters, parsers, generators) sharing common UI shell + CLI command structure + MCP tool registration.
- **User interactions:** `devforge json fmt`, web UI tool grid, MCP tools (`devforge.regex_test`, etc.).
- **Dependencies:** per-tool libraries (e.g., faker-js, prettier core, croner).

### B2 — Unified Shell
- **Key functionalities:** plugin loader, command router, web server, MCP server, auth manager.
- **User interactions:** `devforge --help`, `devforge ui`, MCP client config snippet.
- **Dependencies:** Go/Rust CLI framework; embedded web assets; MCP SDK.

---

## 4. Tasks (Execution-Ready)

### A1 — API Contract Drift Guard
| Task | Description | Expected Outcome |
|---|---|---|
| A1-T1 | Build OpenAPI 3.1/3.2 + GraphQL parser core | Canonical AST for both spec types |
| A1-T2 | Build language AST extractors (TS, Python, Go) | Endpoint inventory from source |
| A1-T3 | Implement static drift engine (spec ↔ code) | Drift report JSON with breaking/non-breaking classification |
| A1-T4 | Build live-traffic sidecar proxy | Endpoint usage telemetry stream |
| A1-T5 | Build GitHub App + PR commenter | PR comment showing drift summary |
| A1-T6 | Build client-type generator (TS, Python) | `devforge api gen-types` outputs typed clients |
| A1-T7 | Expose MCP server (`getEndpoint`, `validateRequest`) | Agents query contract deterministically |
| A1-T8 | E2E demo with sample monorepo | Working PR flow end-to-end |

### A2 — Log Analyzer Lite
| Task | Description | Expected Outcome |
|---|---|---|
| A2-T1 | Implement Drain3-based template extractor | Templates from raw logs |
| A2-T2 | Build format auto-detector (Nginx/Apache/syslog/JSON) | Correct parsing without config |
| A2-T3 | Build two-stage anomaly filter (normal-pattern removal + LLM scoring) | Anomaly stream with scores |
| A2-T4 | Build RAG layer over historical logs | Context-aware anomaly explanations |
| A2-T5 | Build CLI dashboard (`devforge logs tail`) | Live terminal UI |
| A2-T6 | Build web dashboard (trends, frequency) | Browser view at `/logs` |
| A2-T7 | Webhook / Slack alert hooks | Outbound alerts on critical anomalies |

### A3 — Secrets Gateway
| Task | Description | Expected Outcome |
|---|---|---|
| A3-T1 | Build GitHub OAuth flow | `devforge auth login` works |
| A3-T2 | Build encrypted backend KV store (AES-256-GCM) | Secret CRUD API |
| A3-T3 | Build memory-only injection runner | `devforge run -- <cmd>` |
| A3-T4 | Build schema-validation adapters (Zod, Pydantic) | Startup validation of env |
| A3-T5 | Build drift watcher (local vs canonical) | Local diff alerts |
| A3-T6 | Expose MCP server for scoped secret read | Agents fetch secrets safely |
| A3-T7 | Build web dashboard for team secret management | UI for permissions / audit |

### A4 — Knowledge Base Compiler
| Task | Description | Expected Outcome |
|---|---|---|
| A4-T1 | Build markdown ingestor + change detector | Incremental file diff |
| A4-T2 | Build LLM compilation pipeline (raw → wiki) | Structured `.md` with backlinks |
| A4-T3 | Build SQLite + sqlite-vec index | Vector search over wiki |
| A4-T4 | Build `query` command with citations | Cited answers in CLI |
| A4-T5 | Build health-check linter | Contradiction / orphan / stale-link reports |
| A4-T6 | Build Marp slide exporter | `.md` → slides pipeline |

### A5 — Playwright Healer
| Task | Description | Expected Outcome |
|---|---|---|
| A5-T1 | Build Playwright runner wrapper | Captures failures + traces |
| A5-T2 | Build accessibility-tree locator proposer | Stable selector candidates |
| A5-T3 | Build failure classifier (UI / timing / data) | Tagged failures |
| A5-T4 | Build visual snapshot differ | Diff artifacts |
| A5-T5 | Build `fix-proposal.md` generator | Reviewable patch proposal |
| A5-T6 | Build PR auto-creator (GitHub API) | One-click source patches |
| A5-T7 | Build healer transparency dashboard | Per-failure audit view |

### B1 — Utility Toolkit
| Task | Description | Expected Outcome |
|---|---|---|
| B1-T1 | JSON / YAML / CSV formatter + validator | CLI + web + MCP |
| B1-T2 | Smart diff viewer (semantic) | JSON / SQL / config diffs |
| B1-T3 | Regex tester + explainer | Match preview + plain-English explanation |
| B1-T4 | Cron builder + visualizer | Next-run preview |
| B1-T5 | JWT debugger | Decode / verify / expiry alerts |
| B1-T6 | UUID / hash generator | Multiple versions / formats |
| B1-T7 | Timezone + dev-unit converter | DST-aware; bytes/latency/throughput |
| B1-T8 | SQL formatter + linter | Pretty + lint warnings |
| B1-T9 | Markdown editor + preview | Web UI editor |
| B1-T10 | Data faker (multi-locale) | Realistic mock data |
| B1-T11 | cURL ↔ snippet ↔ Postman | Bidirectional conversion |
| B1-T12 | Webhook tester | Request bin + replay |

### B2 — Unified Shell
| Task | Description | Expected Outcome |
|---|---|---|
| B2-T1 | Build single-binary CLI scaffold (Go) | `devforge` binary |
| B2-T2 | Build plugin loader / command router | Extensible architecture |
| B2-T3 | Embed web UI assets + serve | `devforge ui` opens browser |
| B2-T4 | Build MCP server scaffold | All tools registered as MCP tools |
| B2-T5 | Build GitHub OAuth manager | Shared auth across features |
| B2-T6 | Build auto-update + version channel | One-command updates |
| B2-T7 | Build telemetry (opt-in, anonymized) | Usage insight without PII |

---

## 5. Prioritization

### MVP (Phase 1 — first release)
**Goal:** Ship the wedge product with one flagship + a strong utility bundle, all MCP-enabled.

- **B2 — Unified Shell** (CLI + UI + MCP scaffolding) — non-negotiable foundation.
- **B1 — Utility Toolkit** subset: JSON/YAML/CSV formatter, regex tester, JWT debugger, cron builder, UUID/hash, timezone converter, smart diff, data faker. (8 of 12 tools.)
- **A3 — GitHub-Native Secrets Gateway** — highest market-gap-to-complexity ratio (Tier 1, medium build).
- **A1 — API Contract Drift Guard (lite)** — static-only drift (no live traffic yet); TS + Python only.

### Phase 2 (next 3–6 months)
- **A1 full** — add live-traffic sidecar, GraphQL, Go support.
- **A5 — Playwright Healer Dashboard.**
- **A4 — Knowledge Base Compiler.**
- **B1 remainder** — SQL formatter, markdown editor, cURL converter, webhook tester.

### Phase 3 (6–12 months)
- **A2 — Log Analyzer Lite** (full anomaly + RAG).
- Plugin marketplace for community-built utilities.
- Team-tier billing + RBAC.
- Additional language support across A1 and A3 (Ruby, Rust, Java).
- Optional cloud-sync layer for utility state (snippets, configs).

---

## 6. Out-of-Scope (Product-wide)

- Mobile native apps.
- Hosted SaaS multi-tenant version of every tool (cloud is opt-in, only for sync / team features).
- Replacement of full IDEs, full APM platforms, or full SIEM.
- Non-MCP LLM integrations (custom plugins for ChatGPT, Gemini, etc.) — MCP-first only.

---

## 7. Key Risks & Mitigations

| Risk | Mitigation |
|---|---|
| Scope creep across 5 flagships + 12 utilities | Strict MVP gating — only A3 + A1-lite + 8 utilities ship in v1 |
| MCP ecosystem still maturing | Build MCP layer thin; CLI remains canonical |
| Auth / secrets storage = high security surface | External pen-test before A3 GA; AES-256-GCM + per-org keys |
| LLM costs for A2 / A4 | BYO-key model; aggressive incremental compilation |
| Differentiation vs Postman / Doppler / Splunk | Win on offline-first + MCP-native + GitHub-native onboarding |
