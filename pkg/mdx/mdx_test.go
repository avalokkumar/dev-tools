package mdx

import (
	"strings"
	"testing"
)

func TestToHTML_BasicRender(t *testing.T) {
	t.Parallel()
	r, _ := ToHTML([]byte("# Title\n\nHello **world**."), ToHTMLOptions{})
	if !strings.Contains(r.Output, "<h1") {
		t.Fatalf("missing h1: %s", r.Output)
	}
	if !strings.Contains(r.Output, "<strong>world</strong>") {
		t.Fatalf("missing bold: %s", r.Output)
	}
}

func TestToHTML_StripsScripts(t *testing.T) {
	t.Parallel()
	in := []byte(`# Title

<script>alert('xss')</script>

inline <img src="x" onerror="alert(1)">`)
	r, _ := ToHTML(in, ToHTMLOptions{})
	if strings.Contains(r.Output, "<script") {
		t.Fatalf("scripts not stripped: %s", r.Output)
	}
	if strings.Contains(r.Output, "onerror") {
		t.Fatalf("onerror handler not stripped: %s", r.Output)
	}
}

func TestTableFromCSV_Basic(t *testing.T) {
	t.Parallel()
	in := []byte("name,age\nalok,42\nbob,1\n")
	r, _ := TableFromCSV(in, TableOptions{})
	expected := []string{"| name | age |", "| --- | --- ", "| alok | 42 |", "| bob | 1 |"}
	for _, w := range expected {
		if !strings.Contains(r.Output, w) {
			t.Fatalf("missing %q in:\n%s", w, r.Output)
		}
	}
}

func TestTableFromCSV_AlignmentMarkers(t *testing.T) {
	t.Parallel()
	in := []byte("a,b\n1,2\n")
	r, _ := TableFromCSV(in, TableOptions{Align: []string{"left", "right"}})
	if !strings.Contains(r.Output, ":---") || !strings.Contains(r.Output, "---:") {
		t.Fatalf("alignment markers missing: %s", r.Output)
	}
}

func TestTableFromCSV_EscapesPipes(t *testing.T) {
	t.Parallel()
	in := []byte("a,b\nx|y,z\n")
	r, _ := TableFromCSV(in, TableOptions{})
	if !strings.Contains(r.Output, `x\|y`) {
		t.Fatalf("pipe not escaped: %s", r.Output)
	}
}
