# Installing DevForge

DevForge is a single static binary. Pick whichever channel you prefer — they
all install the **same artifact** built and signed by `goreleaser` from the
[Releases page](https://github.com/devforge/devforge/releases).

| Channel              | Best for                                | Auto-update |
|----------------------|-----------------------------------------|-------------|
| `install.sh`         | Most users (macOS, Linux)               | Re-run       |
| Homebrew             | macOS / Linuxbrew users                 | `brew upgrade` |
| Manual download      | Air-gapped or pinned environments       | Manual       |
| `go install`         | Contributors building from source       | Manual       |

> **Windows** — download the `.zip` from the Releases page, extract, and put
> `devforge.exe` somewhere on your `PATH`. A native winget / scoop manifest is
> on the roadmap.

---

## 1. One-line installer (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/devforge/devforge/main/install.sh | bash
```

What it does:

1. Detects your OS (`linux` / `darwin`) and architecture (`amd64` / `arm64`).
2. Resolves the latest release tag from GitHub.
3. Downloads `devforge_<version>_<os>_<arch>.tar.gz` over HTTPS.
4. Verifies the **SHA-256** against `checksums.txt` from the same release.
5. Installs the binary to `/usr/local/bin` if writable, otherwise to
   `~/.local/bin` (no sudo required).
6. Adds the install dir to your shell rc (`~/.zshrc`, `~/.bashrc`, or
   fish equivalent) only if it is not already on `PATH`.
7. Runs `devforge version` to verify the install.

The installer is **idempotent** — running it again on the same version is a
no-op; running it on a newer tag upgrades in place.

### Customising the install

Use environment variables or flags (flags win):

```bash
# Pin a specific release.
DEVFORGE_VERSION=v0.1.0 curl -fsSL .../install.sh | bash

# Install into a custom prefix (no sudo, no PATH edits).
curl -fsSL .../install.sh | bash -s -- --prefix "$HOME/.local" --no-modify-path

# Drop the binary directly into a folder you control.
curl -fsSL .../install.sh | bash -s -- --dir /opt/devforge/bin

# Force-reinstall the same version (e.g. after a bad install).
curl -fsSL .../install.sh | bash -s -- --force
```

| Variable                    | Flag                | Default                |
|-----------------------------|---------------------|------------------------|
| `DEVFORGE_VERSION`          | `--version`         | latest GitHub release  |
| `DEVFORGE_REPO`             | —                   | `devforge/devforge`    |
| `DEVFORGE_PREFIX`           | `--prefix`          | `/usr/local`           |
| `DEVFORGE_INSTALL_DIR`      | `--dir`             | `<prefix>/bin`         |
| `DEVFORGE_NO_MODIFY_PATH=1` | `--no-modify-path`  | (modifies PATH)        |
| `DEVFORGE_FORCE=1`          | `--force`           | (idempotent)           |

### Reviewing the script before running

Always a good habit:

```bash
curl -fsSL https://raw.githubusercontent.com/devforge/devforge/main/install.sh -o install.sh
less install.sh
bash install.sh
```

---

## 2. Homebrew (macOS + Linux)

Once the [`devforge/tap`](https://github.com/devforge/homebrew-tap) is published:

```bash
brew tap devforge/tap
brew install devforge
brew upgrade devforge   # later
```

Brew formula source lives in this repo at
[`packaging/homebrew/devforge.rb`](./homebrew/devforge.rb) and is the
template `goreleaser` publishes to the tap on every tagged release.

To dry-run from a local clone:

```bash
brew install --build-from-source ./packaging/homebrew/devforge.rb
```

---

## 3. Manual download

```bash
VERSION=0.1.0
OS=darwin                 # or linux
ARCH=arm64                # or amd64
URL="https://github.com/devforge/devforge/releases/download/v${VERSION}"

curl -fsSLO "${URL}/devforge_${VERSION}_${OS}_${ARCH}.tar.gz"
curl -fsSLO "${URL}/checksums.txt"

# 1) Verify (one of these will be available on your system).
shasum  -a 256 -c checksums.txt --ignore-missing  || \
sha256sum      -c checksums.txt --ignore-missing

# 2) Install.
tar -xzf "devforge_${VERSION}_${OS}_${ARCH}.tar.gz"
sudo install -m 0755 devforge /usr/local/bin/devforge

devforge version --json
```

---

## 4. Build from source

Requires Go 1.25+ and pnpm 10+ (only if you want to rebuild the embedded UI).

```bash
git clone https://github.com/devforge/devforge.git
cd devforge
./scripts/build-web.sh                # rebuild & embed the SPA
go build -o ./bin/devforge ./cmd/devforge
./bin/devforge version --json
```

---

## Post-install: verify everything works

```bash
devforge version --json     # build metadata
devforge run --list         # 75 operations
devforge ui                 # local web UI on an ephemeral port
devforge mcp                # MCP stdio server (Ctrl-D to exit)
```

### Wiring DevForge into an MCP-aware client

`devforge mcp` is a JSON-RPC stdio server compliant with MCP 2024-11-05.
Drop the snippet below into your client's config (Claude Code's `.mcp.json`,
Cursor's `~/.cursor/mcp.json`, or Claude Desktop's `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "devforge": {
      "command": "/usr/local/bin/devforge",
      "args": ["mcp"]
    }
  }
}
```

Then ask the agent things like:

* *"Generate three v7 UUIDs."*
* *"Decode this JWT and tell me when it expires."*
* *"Convert 2026-03-08 02:30 New York time to UTC and flag DST issues."*

---

## Uninstalling

```bash
# Remove the binary + PATH stanza.
curl -fsSL https://raw.githubusercontent.com/devforge/devforge/main/uninstall.sh | bash -s -- --yes

# Also remove ~/.devforge (telemetry log, plugins, cache).
curl -fsSL https://raw.githubusercontent.com/devforge/devforge/main/uninstall.sh | bash -s -- --yes --purge

# Or, if installed via Homebrew:
brew uninstall devforge
```

---

## Troubleshooting

**`devforge: command not found` after install**
The installer added the install directory to your shell rc but the current
shell has not reloaded yet. Open a new terminal, or run
`source ~/.zshrc` (or your shell's equivalent).

**`error: failed to download …`**
Check connectivity to `github.com`. If you are behind a corporate proxy,
export `HTTPS_PROXY=http://proxy:port` before running the installer.

**`error: checksum mismatch`**
The download was tampered with or truncated. Re-run; if the error persists,
report it on GitHub Issues with the `--version` you pinned.

**`Permission denied` on `/usr/local/bin`**
Re-run with `--prefix "$HOME/.local"` (no sudo needed) or with `sudo bash`.

**Architecture not supported**
DevForge ships pre-built binaries for `linux/{amd64,arm64}` and
`darwin/{amd64,arm64}`. For anything else, build from source (see § 4).

---

## Security notes

* Every published artifact has a deterministic SHA-256 in `checksums.txt`.
* `install.sh` only fetches over HTTPS (TLS 1.2+) and refuses to proceed on
  a checksum mismatch.
* The installer never runs the downloaded binary before verifying its hash.
* DevForge itself runs **localhost-only** by default; the `http_request`
  and `dns_lookup` operations refuse private-network destinations unless
  you explicitly opt in via `allowPrivate: true`.
* Telemetry is **off by default** and never auto-uploads. See `devforge
  telemetry --help` for details.
