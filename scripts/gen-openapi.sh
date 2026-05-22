#!/usr/bin/env bash
# gen-openapi.sh — regenerate the OpenAPI doc and TypeScript client.
#
# Phase A: validates the hand-written api/openapi.yaml and produces a TS client.
# Phase B+: replace this with `go run ./cmd/openapi-gen` that emits the spec
# from the live Registry.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SPEC="$ROOT/api/openapi.yaml"
OUT="$ROOT/packages/api-types/src/generated/devforge.ts"

if [ ! -f "$SPEC" ]; then
  echo "[gen-openapi] missing spec at $SPEC" >&2
  exit 1
fi

mkdir -p "$(dirname "$OUT")"

# Use openapi-typescript via pnpm dlx to avoid pinning a workspace dep until
# we wire it in for real.
pnpm dlx openapi-typescript "$SPEC" -o "$OUT"
echo "[gen-openapi] wrote $OUT"
