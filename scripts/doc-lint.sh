#!/usr/bin/env bash
# doc-lint.sh — sanity-check the canonical project docs.
#
# Fails CI when:
#   - docs/terms.md or docs/architecture.md is missing
#   - terms.md is missing one of the Five Names (Tool, Operation, Engine,
#     Surface, Adapter), which would mean the ubiquitous-language contract
#     was broken
#   - any ADR file is malformed (no Status field)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
fail=0
note() { echo "[doc-lint] $*" >&2; }

for f in docs/terms.md docs/architecture.md README.md; do
  if [ ! -s "$ROOT/$f" ]; then
    note "MISSING $f"; fail=1
  fi
done

if [ -f "$ROOT/docs/terms.md" ]; then
  for term in Tool Operation Engine Surface Adapter; do
    if ! grep -q "\\*\\*$term\\*\\*" "$ROOT/docs/terms.md"; then
      note "terms.md missing core term: $term"; fail=1
    fi
  done
fi

shopt -s nullglob
for adr in "$ROOT"/docs/adrs/*.md; do
  if ! grep -q "^\\*\\*Status:" "$adr"; then
    note "ADR missing Status: $(basename "$adr")"; fail=1
  fi
done

if [ $fail -ne 0 ]; then
  note "FAIL"
  exit 1
fi
note "OK"
