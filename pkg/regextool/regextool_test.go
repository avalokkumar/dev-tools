package regextool

import (
	"strings"
	"testing"
)

// TestRegex_RE2_BasicMatch — C5: simple literal match.
func TestRegex_RE2_BasicMatch(t *testing.T) {
	t.Parallel()
	res, err := Test(`foo`, "foo bar foo", TestOptions{})
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if len(res.Matches) != 2 {
		t.Fatalf("got %d matches", len(res.Matches))
	}
	if res.Matches[0].Value != "foo" {
		t.Fatalf("first match = %q", res.Matches[0].Value)
	}
}

// TestRegex_NamedGroups — C5: named groups appear in Match.Groups.
func TestRegex_NamedGroups(t *testing.T) {
	t.Parallel()
	res, err := Test(`(?P<year>\d{4})-(?P<mon>\d{2})`, "release 2026-05", TestOptions{})
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if len(res.Matches) != 1 {
		t.Fatalf("got %d matches", len(res.Matches))
	}
	groups := res.Matches[0].Groups
	if len(groups) != 2 || groups[0].Name != "year" || groups[0].Value != "2026" {
		t.Fatalf("groups: %+v", groups)
	}
}

// TestRegex_Explain_Quantifiers — C5: explain catches +, ?, .
func TestRegex_Explain_Quantifiers(t *testing.T) {
	t.Parallel()
	res, err := Explain(`a+b?c.`, ExplainOptions{})
	if err != nil {
		t.Fatalf("Explain: %v", err)
	}
	want := []string{"+", "?", "."}
	got := strings.Join(tokens(res.Tree), " ")
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Fatalf("missing %q in tokens %q", w, got)
		}
	}
}

// TestRegex_InvalidPattern_DiagnosticCode — C5: bad regex emits diagnostic.
func TestRegex_InvalidPattern_DiagnosticCode(t *testing.T) {
	t.Parallel()
	res, err := Test(`(unclosed`, "x", TestOptions{})
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if res.Diagnostics[0].Code != "REGEX.INVALID_PATTERN" {
		t.Fatalf("code=%q", res.Diagnostics[0].Code)
	}
}

// TestRegex_CaseInsensitiveFlag — C5: flag 'i' makes match case-insensitive.
func TestRegex_CaseInsensitiveFlag(t *testing.T) {
	t.Parallel()
	res, err := Test(`hello`, "HELLO world", TestOptions{Flags: "i"})
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if len(res.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(res.Matches))
	}
}

func tokens(t []ExplainNode) []string {
	out := make([]string, len(t))
	for i, n := range t {
		out[i] = n.Token
	}
	return out
}
