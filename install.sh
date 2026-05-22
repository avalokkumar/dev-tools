#!/usr/bin/env bash
# shellcheck shell=bash
#
# DevForge installer — one-line, idempotent, checksum-verified.
#
#   curl -fsSL https://raw.githubusercontent.com/devforge/devforge/main/install.sh | bash
#
# Environment overrides (all optional):
#
#   DEVFORGE_VERSION   Pin a specific version, e.g. v0.1.0. Default: latest release.
#   DEVFORGE_REPO      GitHub <owner>/<repo>. Default: devforge/devforge.
#   DEVFORGE_BASE_URL  Override the release-asset base URL (useful for
#                      air-gapped mirrors or local smoke tests). When set, the
#                      installer fetches <BASE_URL>/<archive> directly and skips
#                      "latest" resolution; you must also set DEVFORGE_VERSION.
#   DEVFORGE_PREFIX    Install prefix. Default: /usr/local (then bin/devforge).
#                      Falls back to $HOME/.local automatically when unwritable.
#   DEVFORGE_INSTALL_DIR
#                      Skip prefix logic and install the binary into this dir.
#   DEVFORGE_NO_MODIFY_PATH=1
#                      Don't append the install dir to your shell rc files.
#   DEVFORGE_FORCE=1   Reinstall even if the same version is already present.
#
# Flags:
#
#   --version VER   Same as DEVFORGE_VERSION.
#   --prefix DIR    Same as DEVFORGE_PREFIX.
#   --dir DIR       Same as DEVFORGE_INSTALL_DIR.
#   --no-modify-path  Same as DEVFORGE_NO_MODIFY_PATH=1.
#   --force         Same as DEVFORGE_FORCE=1.
#   -h, --help      Print this help.
#
# Exits non-zero on any failure. Cleans up partial state on interrupt.

set -Eeuo pipefail

# ---------------------------------------------------------------------------
# Defaults & globals
# ---------------------------------------------------------------------------
REPO="${DEVFORGE_REPO:-devforge/devforge}"
VERSION="${DEVFORGE_VERSION:-}"
BASE_URL_OVERRIDE="${DEVFORGE_BASE_URL:-}"
PREFIX="${DEVFORGE_PREFIX:-/usr/local}"
INSTALL_DIR_OVERRIDE="${DEVFORGE_INSTALL_DIR:-}"
NO_MODIFY_PATH="${DEVFORGE_NO_MODIFY_PATH:-0}"
FORCE="${DEVFORGE_FORCE:-0}"
BINARY_NAME="devforge"

# Colors only on TTY.
if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  C_RESET=$'\033[0m'; C_BOLD=$'\033[1m'
  C_GREEN=$'\033[32m'; C_YELLOW=$'\033[33m'; C_RED=$'\033[31m'; C_CYAN=$'\033[36m'
else
  C_RESET=""; C_BOLD=""; C_GREEN=""; C_YELLOW=""; C_RED=""; C_CYAN=""
fi

# tmp_dir is set after we've parsed flags; trap cleans it up on any exit.
TMP_DIR=""
cleanup() { [ -n "${TMP_DIR}" ] && [ -d "${TMP_DIR}" ] && rm -rf "${TMP_DIR}"; }
trap cleanup EXIT
trap 'err "Interrupted."; exit 130' INT TERM

info()  { printf "%s%s%s\n" "${C_CYAN}" "$*" "${C_RESET}"; }
note()  { printf "%s%s%s\n" "${C_BOLD}" "$*" "${C_RESET}"; }
warn()  { printf "%swarning:%s %s\n" "${C_YELLOW}" "${C_RESET}" "$*" >&2; }
err()   { printf "%serror:%s %s\n"   "${C_RED}"    "${C_RESET}" "$*" >&2; }
ok()    { printf "%s✓%s %s\n"        "${C_GREEN}"  "${C_RESET}" "$*"; }

usage() {
  sed -n '2,30p' "$0" | sed 's/^# \{0,1\}//'
}

# ---------------------------------------------------------------------------
# Argument parsing
# ---------------------------------------------------------------------------
while [ $# -gt 0 ]; do
  case "$1" in
    --version)        VERSION="$2"; shift 2 ;;
    --version=*)      VERSION="${1#*=}"; shift ;;
    --prefix)         PREFIX="$2"; shift 2 ;;
    --prefix=*)       PREFIX="${1#*=}"; shift ;;
    --dir)            INSTALL_DIR_OVERRIDE="$2"; shift 2 ;;
    --dir=*)          INSTALL_DIR_OVERRIDE="${1#*=}"; shift ;;
    --no-modify-path) NO_MODIFY_PATH=1; shift ;;
    --force)          FORCE=1; shift ;;
    -h|--help)        usage; exit 0 ;;
    *) err "unknown argument: $1"; usage; exit 2 ;;
  esac
done

# ---------------------------------------------------------------------------
# Environment detection
# ---------------------------------------------------------------------------
require_cmd() {
  command -v "$1" >/dev/null 2>&1 || { err "missing required command: $1"; exit 1; }
}

# Need at least: uname, tar, mkdir, mv. Use curl OR wget for downloads.
require_cmd uname
require_cmd tar
require_cmd mkdir
require_cmd mv

if command -v curl >/dev/null 2>&1; then
  HTTP_CMD="curl"
elif command -v wget >/dev/null 2>&1; then
  HTTP_CMD="wget"
else
  err "neither 'curl' nor 'wget' is available — install one and re-run."
  exit 1
fi

# Pick a SHA-256 binary. macOS has shasum; most Linux has sha256sum.
if command -v sha256sum >/dev/null 2>&1; then
  SHA_CMD="sha256sum"
elif command -v shasum >/dev/null 2>&1; then
  SHA_CMD="shasum -a 256"
else
  warn "no sha256sum/shasum available — checksum verification will be skipped."
  SHA_CMD=""
fi

detect_os() {
  local os
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$os" in
    linux)  echo linux ;;
    darwin) echo darwin ;;
    msys*|mingw*|cygwin*)
      err "Windows is not supported by install.sh. Download the .zip from https://github.com/${REPO}/releases or use 'scoop install devforge' (when published)."
      exit 1 ;;
    *)
      err "unsupported OS: $os. Supported: linux, darwin."
      exit 1 ;;
  esac
}

detect_arch() {
  local arch
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) echo amd64 ;;
    aarch64|arm64) echo arm64 ;;
    *)
      err "unsupported architecture: $arch. Supported: amd64, arm64."
      exit 1 ;;
  esac
}

OS="$(detect_os)"
ARCH="$(detect_arch)"

# ---------------------------------------------------------------------------
# HTTP helpers (curl / wget abstraction)
# ---------------------------------------------------------------------------
http_get_to_file() {
  # $1 url, $2 dest path
  local url="$1" dest="$2"
  if [ "$HTTP_CMD" = "curl" ]; then
    if [ -n "$BASE_URL_OVERRIDE" ]; then
      curl -fsSL --retry 3 --retry-delay 1 --connect-timeout 15 -o "$dest" "$url"
    else
      curl --proto '=https' --tlsv1.2 -fsSL --retry 3 --retry-delay 1 \
           --connect-timeout 15 -o "$dest" "$url"
    fi
  else
    if [ -n "$BASE_URL_OVERRIDE" ]; then
      wget --tries=3 --timeout=30 -qO "$dest" "$url"
    else
      wget --https-only --tries=3 --timeout=30 -qO "$dest" "$url"
    fi
  fi
}

# ---------------------------------------------------------------------------
# Resolve version (latest if not pinned)
# ---------------------------------------------------------------------------
resolve_version() {
  if [ -n "$VERSION" ]; then
    # Normalise: accept "0.1.0" or "v0.1.0".
    case "$VERSION" in v*) ;; *) VERSION="v$VERSION" ;; esac
    echo "$VERSION"
    return
  fi
  if [ -n "$BASE_URL_OVERRIDE" ]; then
    err "DEVFORGE_BASE_URL is set but DEVFORGE_VERSION is empty. Pin a version when overriding the base URL."
    exit 1
  fi
  # Use the redirect from /releases/latest — works without auth and avoids
  # rate-limited /api/.
  local url="https://github.com/${REPO}/releases/latest"
  local resolved=""
  if [ "$HTTP_CMD" = "curl" ]; then
    resolved="$(curl -fsSLI -o /dev/null -w '%{url_effective}' --connect-timeout 15 "$url" || true)"
  else
    # wget prints redirect chain on stderr with -S.
    resolved="$(wget -S --max-redirect=10 --spider "$url" 2>&1 \
                | awk '/^Location: /{loc=$2} END{print loc}' || true)"
  fi
  if [ -z "$resolved" ]; then
    err "could not resolve latest release for ${REPO}. Set DEVFORGE_VERSION or pass --version."
    exit 1
  fi
  echo "${resolved##*/}"
}

# ---------------------------------------------------------------------------
# Choose install directory
# ---------------------------------------------------------------------------
choose_install_dir() {
  if [ -n "$INSTALL_DIR_OVERRIDE" ]; then
    echo "$INSTALL_DIR_OVERRIDE"
    return
  fi
  local primary="${PREFIX%/}/bin"
  local fallback="${HOME}/.local/bin"

  # If the user asked for a custom prefix, honour it without falling back.
  if [ "${PREFIX}" != "/usr/local" ] && [ -n "${DEVFORGE_PREFIX:-}" ]; then
    echo "$primary"
    return
  fi

  # Default: prefer /usr/local/bin if writable now or via existing sudo session,
  # otherwise fall back to ~/.local/bin without prompting.
  if mkdir -p "$primary" 2>/dev/null && [ -w "$primary" ]; then
    echo "$primary"
  else
    mkdir -p "$fallback"
    echo "$fallback"
  fi
}

# ---------------------------------------------------------------------------
# PATH wiring
# ---------------------------------------------------------------------------
shell_rc_for() {
  # Print the rc file we should append PATH to, based on the user's shell.
  local sh="${SHELL##*/}"
  case "$sh" in
    zsh)  echo "${ZDOTDIR:-$HOME}/.zshrc" ;;
    bash)
      # Prefer .bashrc on Linux, .bash_profile on macOS (login shells).
      if [ "$OS" = "darwin" ] && [ -f "$HOME/.bash_profile" ]; then
        echo "$HOME/.bash_profile"
      else
        echo "$HOME/.bashrc"
      fi
      ;;
    fish) echo "$HOME/.config/fish/config.fish" ;;
    *)    echo "" ;;
  esac
}

ensure_in_path() {
  local dir="$1"
  case ":${PATH}:" in *":${dir}:"*) return 0 ;; esac

  if [ "$NO_MODIFY_PATH" = "1" ]; then
    warn "${dir} is not on your PATH. Add it manually, or re-run without --no-modify-path."
    return 0
  fi

  local rc; rc="$(shell_rc_for)"
  if [ -z "$rc" ]; then
    warn "could not detect a shell rc file (SHELL=${SHELL:-unset}). Add ${dir} to PATH manually."
    return 0
  fi
  mkdir -p "$(dirname "$rc")"
  touch "$rc"

  local marker="# Added by DevForge installer"
  if grep -F "$marker" "$rc" >/dev/null 2>&1; then
    return 0
  fi

  {
    printf '\n%s\n' "$marker"
    case "$rc" in
      */config.fish) printf 'fish_add_path -gP "%s"\n' "$dir" ;;
      *)             printf 'export PATH="%s:$PATH"\n' "$dir" ;;
    esac
  } >> "$rc"

  ok "Added ${dir} to PATH in ${rc}"
  warn "Open a new shell or run: source \"${rc}\""
}

# ---------------------------------------------------------------------------
# Existing install detection (idempotency / upgrade)
# ---------------------------------------------------------------------------
existing_version_at() {
  # Print the version reported by an existing devforge binary at $1, if any.
  local bin="$1"
  [ -x "$bin" ] || return 1
  "$bin" version --json 2>/dev/null \
    | sed -n 's/.*"version"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
    | head -n1
}

# ---------------------------------------------------------------------------
# Main install flow
# ---------------------------------------------------------------------------
main() {
  note "DevForge installer"
  info "  repo:    ${REPO}"
  info "  os:      ${OS}"
  info "  arch:    ${ARCH}"

  VERSION="$(resolve_version)"
  info "  version: ${VERSION}"

  local install_dir; install_dir="$(choose_install_dir)"
  local target="${install_dir}/${BINARY_NAME}"
  info "  target:  ${target}"
  printf '\n'

  # Idempotency check.
  if [ -x "$target" ] && [ "$FORCE" != "1" ]; then
    local cur; cur="$(existing_version_at "$target" || true)"
    local pinned="${VERSION#v}"
    if [ -n "$cur" ] && [ "v$cur" = "$VERSION" ]; then
      ok "DevForge ${VERSION} is already installed at ${target} — nothing to do."
      ok "(Re-run with --force to reinstall.)"
      print_usage_footer "$install_dir"
      return 0
    fi
    if [ -n "$cur" ]; then
      info "Upgrading ${BINARY_NAME} from v${cur} → ${VERSION}…"
    fi
    : "${pinned}" # silence unused
  fi

  # Build asset URLs. Archive name strips the leading "v" from the version.
  local ver_no_v="${VERSION#v}"
  local archive="${BINARY_NAME}_${ver_no_v}_${OS}_${ARCH}.tar.gz"
  local base
  if [ -n "$BASE_URL_OVERRIDE" ]; then
    base="${BASE_URL_OVERRIDE%/}"
  else
    base="https://github.com/${REPO}/releases/download/${VERSION}"
  fi
  local archive_url="${base}/${archive}"
  local checksums_url="${base}/checksums.txt"

  TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t devforge)"

  info "Downloading ${archive}…"
  if ! http_get_to_file "$archive_url" "${TMP_DIR}/${archive}"; then
    err "failed to download ${archive_url}"
    err "is the version correct? See https://github.com/${REPO}/releases"
    exit 1
  fi

  if [ -n "$SHA_CMD" ]; then
    info "Verifying SHA-256…"
    if ! http_get_to_file "$checksums_url" "${TMP_DIR}/checksums.txt"; then
      warn "could not download checksums.txt — skipping verification."
    else
      local expected actual
      expected="$(grep " ${archive}\$" "${TMP_DIR}/checksums.txt" | awk '{print $1}')"
      if [ -z "$expected" ]; then
        err "no checksum entry for ${archive} in checksums.txt"
        exit 1
      fi
      actual="$( ( cd "$TMP_DIR" && $SHA_CMD "$archive" ) | awk '{print $1}')"
      if [ "$expected" != "$actual" ]; then
        err "checksum mismatch for ${archive}"
        err "  expected: ${expected}"
        err "  actual:   ${actual}"
        exit 1
      fi
      ok "Checksum verified."
    fi
  fi

  info "Extracting…"
  ( cd "$TMP_DIR" && tar -xzf "$archive" )
  local src="${TMP_DIR}/${BINARY_NAME}"
  if [ ! -x "$src" ]; then
    # Some archives nest the binary under a folder.
    src="$(find "$TMP_DIR" -type f -name "$BINARY_NAME" -perm -u+x | head -n1 || true)"
  fi
  if [ ! -x "$src" ]; then
    err "extracted archive does not contain a '${BINARY_NAME}' binary"
    exit 1
  fi

  info "Installing to ${target}…"
  mkdir -p "$install_dir"
  if ! mv -f "$src" "$target" 2>/dev/null; then
    if command -v sudo >/dev/null 2>&1; then
      warn "writing to ${install_dir} requires elevated privileges. Re-running with sudo…"
      sudo mkdir -p "$install_dir"
      sudo mv -f "$src" "$target"
      sudo chmod 0755 "$target"
    else
      err "no write permission for ${install_dir} and 'sudo' is not available."
      err "Re-run with --prefix=\$HOME/.local or --dir=<writable-path>."
      exit 1
    fi
  else
    chmod 0755 "$target"
  fi
  ok "Installed."

  ensure_in_path "$install_dir"

  # Verify the freshly-installed binary actually runs.
  info "Verifying installation…"
  if ! "$target" version --json >/dev/null 2>&1; then
    # Fall back to plain `version` (older builds may not support --json).
    if ! "$target" version >/dev/null 2>&1; then
      err "the installed binary did not run successfully."
      exit 1
    fi
  fi
  local installed
  installed="$("$target" version --json 2>/dev/null \
              | sed -n 's/.*"version"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
              | head -n1)"
  ok "DevForge v${installed:-unknown} is ready."

  print_usage_footer "$install_dir"
}

print_usage_footer() {
  local install_dir="$1"
  cat <<EOF

${C_BOLD}Get started${C_RESET}
  ${C_GREEN}devforge run --list${C_RESET}        list all 75 operations
  ${C_GREEN}devforge ui${C_RESET}                start the local web UI
  ${C_GREEN}devforge mcp${C_RESET}               start the MCP stdio server
  ${C_GREEN}devforge --help${C_RESET}            full CLI reference

${C_BOLD}MCP integration${C_RESET}
  Add this to Claude Code / Cursor / Claude Desktop:

    "mcpServers": {
      "devforge": {
        "command": "${install_dir}/${BINARY_NAME}",
        "args": ["mcp"]
      }
    }

${C_BOLD}Uninstall${C_RESET}
  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/uninstall.sh | bash

EOF
}

main "$@"
