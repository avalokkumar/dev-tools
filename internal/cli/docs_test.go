package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestDocsGen_OutputsAllSubcommands — D3: `devforge docs gen --out DIR`
// produces one .md per subcommand including the major tools.
func TestDocsGen_OutputsAllSubcommands(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(),
		[]string{"docs", "gen", "--out", dir},
		Stdio{Out: &stdout, Err: &stderr})
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{
		"devforge.md",
		"devforge_uuid.md",
		"devforge_json.md",
		"devforge_mcp.md",
		"devforge_ui.md",
		"devforge_run.md",
		"devforge_update.md",
	} {
		p := filepath.Join(dir, want)
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing %s: %v", p, err)
		}
	}
}
