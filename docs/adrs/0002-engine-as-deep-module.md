# ADR 0002 — Engines as deep modules

**Status:** Accepted
**Date:** 2026-05-09

## Context

Each utility has lots of internal complexity (parsers, locale tables, RFC compliance). All three Surfaces need to use it.

## Decision

Each `pkg/<tool>/` package exposes 1–3 exported functions with `Options` structs and `Result` structs. No interfaces, no factories, no abstraction layers. Internals stay opaque.

`error` return is reserved for catastrophic failure (engine bug, OOM). User-input issues are `Diagnostic` entries inside `Result`.

## Consequences

- Adapters are 1:1 mappings — easy to write, easy to review.
- Engines can be tested in pure isolation via golden tests.
- Refactoring internals does not break any Surface.
- Adding `context.Context` later (for cancellable long ops) is a non-breaking parameter addition.
