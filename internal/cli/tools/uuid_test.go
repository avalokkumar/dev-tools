package tools

import (
	"bytes"
	"strings"
	"testing"
)

// TestCLI_UuidGen_PrintsThree — B4: `uuid gen -n 3` prints three lines.
func TestCLI_UuidGen_PrintsThree(t *testing.T) {
	t.Parallel()
	cmd := NewUUIDCmd()
	cmd.SetArgs([]string{"gen", "-n", "3"})
	// cobra requires the persistent --json flag for our hasErrorDiag path; add it locally.
	cmd.PersistentFlags().Bool("json", false, "")

	var out, errBuf bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v; stderr=%s", err, errBuf.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3:\n%s", len(lines), out.String())
	}
	for _, l := range lines {
		if len(l) != 36 {
			t.Fatalf("line not 36 chars: %q", l)
		}
	}
}

// TestCLI_UuidGen_JsonFormat — B4: --json emits parseable JSON.
func TestCLI_UuidGen_JsonFormat(t *testing.T) {
	t.Parallel()
	cmd := NewUUIDCmd()
	cmd.PersistentFlags().Bool("json", false, "")
	cmd.SetArgs([]string{"gen", "-n", "1", "--json"})
	if err := cmd.PersistentFlags().Set("json", "true"); err != nil {
		t.Fatalf("set json: %v", err)
	}
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(out.String(), `"values"`) {
		t.Fatalf("missing 'values' in JSON: %s", out.String())
	}
}
