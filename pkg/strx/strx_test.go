package strx

import (
	"strings"
	"testing"
)

func TestCase_AllModes(t *testing.T) {
	t.Parallel()
	cases := []struct{ mode, in, want string }{
		{"camel", "hello world example", "helloWorldExample"},
		{"camel", "Hello-World_Example", "helloWorldExample"},
		{"pascal", "hello world", "HelloWorld"},
		{"snake", "HelloWorldExample", "hello_world_example"},
		{"kebab", "Hello World", "hello-world"},
		{"constant", "hello world", "HELLO_WORLD"},
		{"dot", "hello world", "hello.world"},
		{"train", "hello world", "Hello-World"},
		{"title", "hello world", "Hello World"},
		{"lower", "Hello World", "hello world"},
		{"upper", "Hello World", "HELLO WORLD"},
	}
	for _, c := range cases {
		r, _ := Case(c.in, c.mode)
		if r.Output != c.want {
			t.Fatalf("%s(%q) = %q, want %q", c.mode, c.in, r.Output, c.want)
		}
	}
}

func TestCase_UnknownMode(t *testing.T) {
	t.Parallel()
	r, _ := Case("x", "ZALGO")
	if len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "STR.CASE.UNKNOWN_MODE" {
		t.Fatalf("expected diag, got %+v", r.Diagnostics)
	}
}

func TestDiff_Basic(t *testing.T) {
	t.Parallel()
	r, _ := Diff("a\nb\nc\n", "a\nB\nc\n", DiffOptions{})
	if r.Summary.Adds != 1 || r.Summary.Removes != 1 {
		t.Fatalf("summary = %+v", r.Summary)
	}
}

func TestDiff_IgnoreCase(t *testing.T) {
	t.Parallel()
	r, _ := Diff("Hello\nWorld", "hello\nworld", DiffOptions{IgnoreCase: true})
	if r.Summary.Adds != 0 || r.Summary.Removes != 0 {
		t.Fatalf("expected no diff with IgnoreCase, got %+v", r.Summary)
	}
}

func TestStats_BasicCounts(t *testing.T) {
	t.Parallel()
	r, _ := Stats("alpha beta\ngamma\n")
	if r.Lines != 2 {
		t.Fatalf("Lines = %d", r.Lines)
	}
	if r.Words != 3 {
		t.Fatalf("Words = %d", r.Words)
	}
	if r.LongestLine != 10 {
		t.Fatalf("LongestLine = %d", r.LongestLine)
	}
}

func TestStats_Empty(t *testing.T) {
	t.Parallel()
	r, _ := Stats("")
	if r.Lines != 0 || r.Words != 0 || r.Bytes != 0 {
		t.Fatalf("empty mismatch: %+v", r)
	}
}

func TestSortUnique_Sort(t *testing.T) {
	t.Parallel()
	r, _ := SortUnique("c\na\nb\n", SortOptions{})
	if r.Output != "a\nb\nc" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestSortUnique_UniqueIgnoreCase(t *testing.T) {
	t.Parallel()
	r, _ := SortUnique("a\nA\nb", SortOptions{Unique: true, IgnoreCase: true})
	lines := strings.Split(r.Output, "\n")
	if len(lines) != 2 {
		t.Fatalf("got %v", lines)
	}
}

func TestReplace_Literal(t *testing.T) {
	t.Parallel()
	r, _ := Replace("foo bar foo", "foo", "FOO", ReplaceOptions{})
	if r.Output != "FOO bar FOO" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestReplace_RegexCaseInsensitive(t *testing.T) {
	t.Parallel()
	r, _ := Replace("Foo FOO foo", `f.o`, "X", ReplaceOptions{Regex: true, IgnoreCase: true})
	if r.Output != "X X X" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestReplace_InvalidRegex(t *testing.T) {
	t.Parallel()
	r, _ := Replace("x", "(open", "y", ReplaceOptions{Regex: true})
	if len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "STR.REPLACE.INVALID_REGEX" {
		t.Fatalf("expected diag, got %+v", r.Diagnostics)
	}
}
