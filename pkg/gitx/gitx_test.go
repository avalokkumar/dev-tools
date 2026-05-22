package gitx

import (
	"strings"
	"testing"

	"github.com/devforge/devforge/pkg/engine"
)

func TestPatch_Basic(t *testing.T) {
	t.Parallel()
	left := "alpha\nbeta\ngamma\n"
	right := "alpha\nBETA\ngamma\n"
	r, _ := Patch(left, right, PatchOptions{LeftPath: "a.txt", RightPath: "b.txt"})
	for _, want := range []string{"--- a.txt", "+++ b.txt", "-beta", "+BETA"} {
		if !strings.Contains(r.Output, want) {
			t.Fatalf("missing %q in:\n%s", want, r.Output)
		}
	}
}

func TestPatch_Identical(t *testing.T) {
	t.Parallel()
	r, _ := Patch("same\n", "same\n", PatchOptions{})
	if r.Output != "" {
		t.Fatalf("expected empty patch for identical input, got %q", r.Output)
	}
}

func TestCommitFormat_Valid(t *testing.T) {
	t.Parallel()
	r, _ := CommitFormat("feat(auth): add OAuth login", CommitOptions{})
	if !r.Valid {
		t.Fatalf("expected valid: %+v", r.Diagnostics)
	}
	if r.Type != "feat" || r.Scope != "auth" {
		t.Fatalf("got %+v", r)
	}
}

func TestCommitFormat_Breaking(t *testing.T) {
	t.Parallel()
	r, _ := CommitFormat("feat!: drop legacy api\n\nBREAKING CHANGE: remove /v1", CommitOptions{})
	if !r.Breaking {
		t.Fatalf("expected breaking flag")
	}
}

func TestCommitFormat_InvalidHeader(t *testing.T) {
	t.Parallel()
	r, _ := CommitFormat("just a freeform message", CommitOptions{})
	if r.Valid {
		t.Fatalf("expected invalid")
	}
	if !contains(r.Diagnostics, "GIT.COMMIT.INVALID_HEADER") {
		t.Fatalf("missing diag, got %+v", r.Diagnostics)
	}
}

func TestCommitFormat_UnknownType(t *testing.T) {
	t.Parallel()
	r, _ := CommitFormat("snazzy: something", CommitOptions{})
	if !contains(r.Diagnostics, "GIT.COMMIT.UNKNOWN_TYPE") {
		t.Fatalf("missing UNKNOWN_TYPE warn: %+v", r.Diagnostics)
	}
}

func TestCommitFormat_HeaderTooLong(t *testing.T) {
	t.Parallel()
	long := "feat: " + strings.Repeat("x", 100)
	r, _ := CommitFormat(long, CommitOptions{})
	if !contains(r.Diagnostics, "GIT.COMMIT.HEADER_TOO_LONG") {
		t.Fatalf("missing TOO_LONG warn")
	}
}

func TestIgnoreGen_Single(t *testing.T) {
	t.Parallel()
	r, _ := IgnoreGen([]string{"go"})
	if !strings.Contains(r.Output, "*.exe") {
		t.Fatalf("missing go template content: %s", r.Output)
	}
}

func TestIgnoreGen_Combined(t *testing.T) {
	t.Parallel()
	r, _ := IgnoreGen([]string{"go", "macos"})
	if !strings.Contains(r.Output, ".DS_Store") || !strings.Contains(r.Output, "*.exe") {
		t.Fatalf("missing combined content: %s", r.Output)
	}
}

func TestIgnoreGen_UnknownTemplate(t *testing.T) {
	t.Parallel()
	r, _ := IgnoreGen([]string{"cobol"})
	if !contains(r.Diagnostics, "GIT.IGNORE.UNKNOWN_TEMPLATE") {
		t.Fatalf("missing diag, got %+v", r.Diagnostics)
	}
}

func contains(d []engine.Diagnostic, code string) bool {
	for _, x := range d {
		if x.Code == code {
			return true
		}
	}
	return false
}
