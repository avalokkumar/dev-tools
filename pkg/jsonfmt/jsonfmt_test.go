package jsonfmt

import (
	"strings"
	"testing"
)

// TestJsonFormat_PrettyIndent2 — C1: 2-space indent matches expected layout.
func TestJsonFormat_PrettyIndent2(t *testing.T) {
	t.Parallel()
	res, err := Format([]byte(`{"b":1,"a":[2,3]}`), FormatOptions{Indent: 2})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	got := string(res.Output)
	want := "{\n  \"a\": [\n    2,\n    3\n  ],\n  \"b\": 1\n}"
	if got != want {
		t.Fatalf("output mismatch:\nwant=%q\ngot =%q", want, got)
	}
}

// TestJsonFormat_Compact — C1: Indent=0 → compact JSON.
func TestJsonFormat_Compact(t *testing.T) {
	t.Parallel()
	res, err := Format([]byte("  {\"a\":  1}  "), FormatOptions{})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if string(res.Output) != `{"a":1}` {
		t.Fatalf("got %q", res.Output)
	}
}

// TestJsonValidate_RejectsTrailingComma — C1: trailing commas are rejected.
func TestJsonValidate_RejectsTrailingComma(t *testing.T) {
	t.Parallel()
	res, err := Validate([]byte(`{"a":1,}`), nil)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	if len(res.Diagnostics) == 0 || res.Diagnostics[0].Code != "JSON.PARSE" {
		t.Fatalf("expected JSON.PARSE, got %+v", res.Diagnostics)
	}
}

// TestJsonValidate_SchemaPass — H-2: input that matches schema passes.
func TestJsonValidate_SchemaPass(t *testing.T) {
	t.Parallel()
	schema := []byte(`{"type":"object","required":["name"],"properties":{"name":{"type":"string"},"age":{"type":"integer","minimum":0}}}`)
	res, err := Validate([]byte(`{"name":"alok","age":42}`), schema)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !res.Valid {
		t.Fatalf("expected valid, diags=%+v", res.Diagnostics)
	}
}

// TestJsonValidate_SchemaFail — H-2: violations surface as VIOLATION diags.
func TestJsonValidate_SchemaFail(t *testing.T) {
	t.Parallel()
	schema := []byte(`{"type":"object","required":["name"],"properties":{"name":{"type":"string"},"age":{"type":"integer","minimum":0}}}`)
	res, _ := Validate([]byte(`{"name":"alok","age":-5}`), schema)
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	found := false
	for _, d := range res.Diagnostics {
		if d.Code == "JSON.SCHEMA.VIOLATION" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected JSON.SCHEMA.VIOLATION, got %+v", res.Diagnostics)
	}
}

// TestJsonValidate_SchemaMissingRequired — H-2: missing required surfaces violation.
func TestJsonValidate_SchemaMissingRequired(t *testing.T) {
	t.Parallel()
	schema := []byte(`{"type":"object","required":["name"]}`)
	res, _ := Validate([]byte(`{}`), schema)
	if res.Valid {
		t.Fatalf("expected invalid")
	}
}

// TestJsonValidate_BadSchema — H-2: invalid schema JSON surfaces JSON.SCHEMA.PARSE.
func TestJsonValidate_BadSchema(t *testing.T) {
	t.Parallel()
	res, _ := Validate([]byte(`{"a":1}`), []byte(`{not-json`))
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	if res.Diagnostics[0].Code != "JSON.SCHEMA.PARSE" {
		t.Fatalf("code = %q", res.Diagnostics[0].Code)
	}
}

// TestJsonFormat_EmptyInput — C1: empty input emits diagnostic.
func TestJsonFormat_EmptyInput(t *testing.T) {
	t.Parallel()
	res, err := Format(nil, FormatOptions{})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if !strings.Contains(res.Diagnostics[0].Code, "JSON.EMPTY") {
		t.Fatalf("expected JSON.EMPTY, got %+v", res.Diagnostics)
	}
}

// TestJsonFormat_TrailingNewline — C1: TrailingNewline option appends \n.
func TestJsonFormat_TrailingNewline(t *testing.T) {
	t.Parallel()
	res, err := Format([]byte(`{}`), FormatOptions{TrailingNewline: true})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if !strings.HasSuffix(string(res.Output), "\n") {
		t.Fatalf("missing newline: %q", res.Output)
	}
}
