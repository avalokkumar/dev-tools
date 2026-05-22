package yamlfmt

import (
	"strings"
	"testing"
)

// TestYamlFormat_Indent2 — C2: format with indent 2 produces stable output.
func TestYamlFormat_Indent2(t *testing.T) {
	t.Parallel()
	in := []byte("a: 1\nnested:\n  b: 2\n  c:\n    - x\n    - y\n")
	res, err := Format(in, FormatOptions{Indent: 2})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	got := string(res.Output)
	if !strings.Contains(got, "a: 1") || !strings.Contains(got, "  b: 2") {
		t.Fatalf("unexpected output:\n%s", got)
	}
}

// TestYamlConvert_RoundtripJSON — C2: yaml → json → yaml produces equivalent doc.
func TestYamlConvert_RoundtripJSON(t *testing.T) {
	t.Parallel()
	yamlIn := []byte("name: alok\nage: 99\nhobbies:\n  - go\n  - kayaking\n")
	r1, err := Convert(yamlIn, ConvertOptions{To: "json"})
	if err != nil {
		t.Fatalf("yaml→json: %v", err)
	}
	if !strings.Contains(string(r1.Output), `"name": "alok"`) {
		t.Fatalf("missing name in JSON: %s", r1.Output)
	}
	r2, err := Convert(r1.Output, ConvertOptions{To: "yaml"})
	if err != nil {
		t.Fatalf("json→yaml: %v", err)
	}
	if !strings.Contains(string(r2.Output), "name: alok") {
		t.Fatalf("missing name in yaml: %s", r2.Output)
	}
	if !strings.Contains(string(r2.Output), "- kayaking") {
		t.Fatalf("missing list item: %s", r2.Output)
	}
}

// TestYamlValidate_Bad — C2: invalid YAML emits diagnostic.
func TestYamlValidate_Bad(t *testing.T) {
	t.Parallel()
	res, err := Validate([]byte("a:\n - 1\nb: : :"), nil)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	if res.Diagnostics[0].Code != "YAML.PARSE" {
		t.Fatalf("code = %q", res.Diagnostics[0].Code)
	}
}

// TestYamlConvert_UnknownTarget — C2: unknown target emits diagnostic.
func TestYamlConvert_UnknownTarget(t *testing.T) {
	t.Parallel()
	res, err := Convert([]byte("a: 1"), ConvertOptions{To: "xml"})
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if res.Diagnostics[0].Code != "YAML.CONVERT.UNKNOWN_TARGET" {
		t.Fatalf("code = %q", res.Diagnostics[0].Code)
	}
}
