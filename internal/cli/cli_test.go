package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/devforge/devforge/internal/version"
)

// TestRoot_HelpExitsZero — A2: root --help returns 0 and prints usage.
func TestRoot_HelpExitsZero(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"--help"}, Stdio{Out: &stdout, Err: &stderr})
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "devforge") {
		t.Fatalf("help output missing brand: %q", stdout.String())
	}
}

// TestVersionCmd_Json — A3: version --json returns parseable JSON.
func TestVersionCmd_Json(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	code := Execute(context.Background(), []string{"version", "--json"}, Stdio{Out: &stdout, Err: &bytes.Buffer{}})
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	var got version.Info
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("output is not JSON: %v; raw=%q", err, stdout.String())
	}
	if got.Version == "" {
		t.Fatalf("Version field empty: %+v", got)
	}
}

// TestRoot_UuidGen_E2E — B4: full root invocation reaches the uuid Adapter.
func TestRoot_UuidGen_E2E(t *testing.T) {
	t.Parallel()
	var out, errBuf bytes.Buffer
	code := Execute(context.Background(), []string{"uuid", "gen", "-n", "1"},
		Stdio{Out: &out, Err: &errBuf})
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errBuf.String())
	}
	if len(strings.TrimSpace(out.String())) != 36 {
		t.Fatalf("output not a UUID: %q", out.String())
	}
}

// TestVersionCmd_Text — A3: version (no flags) prints human-readable line.
func TestVersionCmd_Text(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	code := Execute(context.Background(), []string{"version"}, Stdio{Out: &stdout, Err: &bytes.Buffer{}})
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.HasPrefix(stdout.String(), "devforge ") {
		t.Fatalf("unexpected text output: %q", stdout.String())
	}
}
