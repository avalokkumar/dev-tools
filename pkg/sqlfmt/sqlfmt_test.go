package sqlfmt

import (
	"strings"
	"testing"

	"github.com/devforge/devforge/pkg/engine"
)

func TestFormat_SelectIndentsClauses(t *testing.T) {
	t.Parallel()
	in := "select id, name from users where age > 18 order by name"
	r, err := Format(in, FormatOptions{})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	out := r.Output
	for _, want := range []string{"SELECT", "FROM", "WHERE", "ORDER BY"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "\nFROM") {
		t.Fatalf("FROM not on new line:\n%s", out)
	}
}

func TestFormat_LowercaseKeywords(t *testing.T) {
	t.Parallel()
	upper := false
	r, _ := Format("SELECT id FROM t", FormatOptions{Uppercase: &upper})
	if !strings.Contains(r.Output, "select") {
		t.Fatalf("expected lowercase select: %s", r.Output)
	}
}

func TestFormat_Empty(t *testing.T) {
	t.Parallel()
	r, _ := Format("   ", FormatOptions{})
	if r.Diagnostics[0].Code != "SQL.EMPTY" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestValidate_UnbalancedParens(t *testing.T) {
	t.Parallel()
	r, _ := Validate("SELECT (a, b FROM t", ValidateOptions{})
	if r.Valid {
		t.Fatalf("expected invalid")
	}
	if !hasCode(r.Diagnostics, "SQL.UNBALANCED_PARENS") {
		t.Fatalf("missing SQL.UNBALANCED_PARENS: %+v", r.Diagnostics)
	}
}

func TestValidate_LintsSelectStar(t *testing.T) {
	t.Parallel()
	r, _ := Validate("SELECT * FROM users", ValidateOptions{})
	if !hasCode(r.Diagnostics, "SQL.LINT.SELECT_STAR") {
		t.Fatalf("missing SELECT_STAR: %+v", r.Diagnostics)
	}
}

func TestValidate_DeleteWithoutWhere(t *testing.T) {
	t.Parallel()
	r, _ := Validate("DELETE FROM users", ValidateOptions{})
	if !hasCode(r.Diagnostics, "SQL.LINT.DELETE_NO_WHERE") {
		t.Fatalf("missing DELETE_NO_WHERE: %+v", r.Diagnostics)
	}
}

func TestValidate_UpdateWithoutWhere(t *testing.T) {
	t.Parallel()
	r, _ := Validate("UPDATE users SET active = false", ValidateOptions{})
	if !hasCode(r.Diagnostics, "SQL.LINT.UPDATE_NO_WHERE") {
		t.Fatalf("missing UPDATE_NO_WHERE: %+v", r.Diagnostics)
	}
}

func hasCode(d []engine.Diagnostic, code string) bool {
	for _, x := range d {
		if x.Code == code {
			return true
		}
	}
	return false
}
