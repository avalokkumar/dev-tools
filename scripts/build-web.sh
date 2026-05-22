#!/usr/bin/env bash
# build-web.sh — build the Vite SPA and place it where Go can embed it.
#
# Why: //go:embed must read from the package directory or below, so we
# build into web/dist (Vite default) and copy to internal/webserver/dist.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEB="$ROOT/web"
EMBED="$ROOT/internal/webserver/dist"

if [ ! -f "$WEB/package.json" ]; then
  echo "[build-web] no web/package.json yet — keeping placeholder dist." >&2
  exit 0
fi

(cd "$WEB" && pnpm install --frozen-lockfile && pnpm build)

rm -rf "$EMBED"
mkdir -p "$EMBED"
cp -R "$WEB/dist/." "$EMBED/"
echo "[build-web] embedded SPA from $WEB/dist into $EMBED"
