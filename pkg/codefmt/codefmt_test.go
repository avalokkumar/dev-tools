package codefmt

import (
	"strings"
	"testing"
)

func TestFormatGo_Reformats(t *testing.T) {
	t.Parallel()
	in := []byte("package x\nfunc f( a int)int{return a+1}\n")
	r, _ := FormatGo(in)
	if !strings.Contains(r.Output, "func f(a int) int") {
		t.Fatalf("not reformatted:\n%s", r.Output)
	}
}

func TestFormatGo_BadSyntax(t *testing.T) {
	t.Parallel()
	r, _ := FormatGo([]byte("package x\nfunc {\n"))
	if r.Diagnostics[0].Code != "CODE.GO.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestFormatXML_Indents(t *testing.T) {
	t.Parallel()
	in := []byte(`<root><a>1</a><b><c>2</c></b></root>`)
	r, _ := FormatXML(in, XMLOptions{Indent: 2})
	if !strings.Contains(r.Output, "\n  <a>1</a>") {
		t.Fatalf("not indented:\n%s", r.Output)
	}
}

func TestFormatHTML_RoundTrip(t *testing.T) {
	t.Parallel()
	in := []byte(`<p>hi <b>there</b></p>`)
	r, _ := FormatHTML(in)
	if !strings.Contains(r.Output, "<p>hi <b>there</b></p>") {
		t.Fatalf("missing fragment in:\n%s", r.Output)
	}
}
