#!/usr/bin/env bash
# shellcheck shell=bash
#
# DevForge uninstaller — removes the binary, the PATH stanza added by
# install.sh, and (optionally) the user data directory.
#
#   curl -fsSL https://raw.githubusercontent.com/devforge/devforge/main/uninstall.sh | bash
#
# Flags / env:
#
#   --purge | DEVFORGE_PURGE=1   Also delete ~/.devforge (telemetry log,
#                                plugins, cache).
#   --yes   | DEVFORGE_YES=1     Don't prompt for confirmation.
#   -h, --help                   Print this help.
#
# Safe to re-run: anything already gone is silently skipped.

set -Eeuo pipefail

PURGE="${DEVFORGE_PURGE:-0}"
ASSUME_YES="${DEVFORGE_YES:-0}"
BINARY_NAME="devforge"

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  C_RESET=$'\033[0m'; C_BOLD=$'\033[1m'
  C_GREEN=$'\033[32m'; C_YELLOW=$'\033[33m'; C_RED=$'\033[31m'; C_CYAN=$'\033[36m'
else
  C_RESET=""; C_BOLD=""; C_GREEN=""; C_YELLOW=""; C_RED=""; C_CYAN=""
fi

info() { printf "%s%s%s\n" "${C_CYAN}" "$*" "${C_RESET}"; }
warn() { printf "%swarning:%s %s\n" "${C_YELLOW}" "${C_RESET}" "$*" >&2; }
err()  { printf "%serror:%s %s\n"   "${C_RED}"    "${C_RESET}" "$*" >&2; }
ok()   { printf "%s✓%s %s\n"        "${C_GREEN}"  "${C_RESET}" "$*"; }

usage() { sed -n '2,18p' "$0" | sed 's/^# \{0,1\}//'; }

while [ $# -gt 0 ]; do
  case "$1" in
    --purge)    PURGE=1; shift ;;
    --yes|-y)   ASSUME_YES=1; shift ;;
    -h|--help)  usage; exit 0 ;;
    *) err "unknown argument: $1"; usage; exit 2 ;;
  esac
done

confirm() {
  [ "$ASSUME_YES" = "1" ] && return 0
  if [ ! -t 0 ]; then
    # No TTY (piped) — be cautious unless --yes given.
    err "non-interactive run: pass --yes to confirm uninstall."
    exit 1
  fi
  printf "%s%s%s [y/N] " "${C_BOLD}" "$1" "${C_RESET}"
  read -r reply
  case "$reply" in y|Y|yes|YES) return 0 ;; *) return 1 ;; esac
}

# ---------------------------------------------------------------------------
# Locate every devforge binary on PATH plus the well-known install prefixes.
# ---------------------------------------------------------------------------
candidates=()

# Anything currently on PATH (portable across bash 3.2 / 4+ / 5+).
saved_ifs="$IFS"
IFS=:
for d in $PATH; do
  [ -z "$d" ] && d="."
  [ -x "$d/$BINARY_NAME" ] && candidates+=("$d/$BINARY_NAME")
done
IFS="$saved_ifs"

# Common install dirs we (or other installers) might have written to.
for d in \
  "/usr/local/bin/${BINARY_NAME}" \
  "/opt/homebrew/bin/${BINARY_NAME}" \
  "${HOME}/.local/bin/${BINARY_NAME}" \
  "${HOME}/bin/${BINARY_NAME}" \
  "${HOME}/.devforge/bin/${BINARY_NAME}"
do
  [ -e "$d" ] && candidates+=("$d")
done

# Dedupe.
if [ "${#candidates[@]}" -gt 0 ]; then
  IFS=$'\n' read -r -d '' -a candidates < <(printf '%s\n' "${candidates[@]}" | awk '!seen[$0]++' && printf '\0')
fi

if [ "${#candidates[@]}" -eq 0 ]; then
  warn "No 'devforge' binary found on PATH or in standard install dirs."
else
  printf "%sFound the following installations:%s\n" "${C_BOLD}" "${C_RESET}"
  for c in "${candidates[@]}"; do printf "  • %s\n" "$c"; done
  printf '\n'
  if ! confirm "Remove the binaries listed above?"; then
    err "Cancelled."; exit 1
  fi
  for bin in "${candidates[@]}"; do
    if [ -w "$bin" ] || [ -w "$(dirname "$bin")" ]; then
      rm -f "$bin" && ok "Removed ${bin}"
    elif command -v sudo >/dev/null 2>&1; then
      warn "Removing ${bin} requires sudo…"
      sudo rm -f "$bin" && ok "Removed ${bin}"
    else
      err "Cannot remove ${bin} (no write permission, no sudo)."
    fi
  done
fi

# ---------------------------------------------------------------------------
# Remove the PATH stanza added by install.sh from common rc files.
# ---------------------------------------------------------------------------
remove_path_stanza() {
  local rc="$1"
  [ -f "$rc" ] || return 0
  grep -F "Added by DevForge installer" "$rc" >/dev/null 2>&1 || return 0

  # Use a portable sed that strips the marker comment plus the next line.
  local tmp; tmp="$(mktemp)"
  awk '
    /# Added by DevForge installer/ { skip = 2; next }
    skip > 0 { skip--; next }
    { print }
  ' "$rc" > "$tmp"
  mv "$tmp" "$rc"
  ok "Cleaned PATH stanza from ${rc}"
}

for rc in \
  "${HOME}/.zshrc" \
  "${HOME}/.bashrc" \
  "${HOME}/.bash_profile" \
  "${HOME}/.profile" \
  "${HOME}/.config/fish/config.fish"
do
  remove_path_stanza "$rc"
done

# ---------------------------------------------------------------------------
# Optional: purge user data.
# ---------------------------------------------------------------------------
if [ "$PURGE" = "1" ]; then
  data_dir="${HOME}/.devforge"
  if [ -d "$data_dir" ]; then
    if confirm "Also delete ${data_dir}?"; then
      rm -rf "$data_dir"
      ok "Removed ${data_dir}"
    else
      info "Kept ${data_dir}"
    fi
  fi
else
  if [ -d "${HOME}/.devforge" ]; then
    info "Kept ${HOME}/.devforge (telemetry log, plugins). Re-run with --purge to delete."
  fi
fi

ok "Uninstall complete."
