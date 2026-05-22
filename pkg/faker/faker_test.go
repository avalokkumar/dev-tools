package faker

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestFaker_Seeded_Deterministic — C9: same seed → same output.
func TestFaker_Seeded_Deterministic(t *testing.T) {
	t.Parallel()
	spec := Spec{Fields: []FieldSpec{
		{Name: "id", Kind: "uuid"},
		{Name: "name", Kind: "name"},
		{Name: "n", Kind: "int", Params: map[string]any{"min": 0, "max": 100}},
	}}
	a, err := Generate(spec, GenerateOptions{Count: 5, Seed: 42, Format: "json"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	b, err := Generate(spec, GenerateOptions{Count: 5, Seed: 42, Format: "json"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if string(a.Output) != string(b.Output) {
		t.Fatalf("seeded output differs:\nA=%s\nB=%s", a.Output, b.Output)
	}
}

// TestFaker_OutputFormat_SQL_Insert — C9: SQL format produces INSERT lines.
func TestFaker_OutputFormat_SQL_Insert(t *testing.T) {
	t.Parallel()
	spec := Spec{Fields: []FieldSpec{
		{Name: "id", Kind: "uuid"},
		{Name: "email", Kind: "email"},
	}}
	res, err := Generate(spec, GenerateOptions{Count: 2, Seed: 1, Format: "sql", Table: "users"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	out := string(res.Output)
	if !strings.Contains(out, "INSERT INTO users") {
		t.Fatalf("missing INSERT: %s", out)
	}
	if strings.Count(out, "\n") != 2 {
		t.Fatalf("expected 2 inserts, got: %s", out)
	}
}

// TestFaker_CSV — C9: CSV output has header + N rows.
func TestFaker_CSV(t *testing.T) {
	t.Parallel()
	spec := Spec{Fields: []FieldSpec{{Name: "id", Kind: "sequence"}, {Name: "k", Kind: "fixed", Params: map[string]any{"value": "K"}}}}
	res, err := Generate(spec, GenerateOptions{Count: 3, Format: "csv"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	got := strings.Split(strings.TrimRight(string(res.Output), "\n"), "\n")
	if len(got) != 4 { // header + 3 rows
		t.Fatalf("got %d lines: %v", len(got), got)
	}
	if got[0] != "id,k" {
		t.Fatalf("header = %q", got[0])
	}
}

// TestFaker_UnknownKind — C9: unknown kind surfaces diagnostic.
func TestFaker_UnknownKind(t *testing.T) {
	t.Parallel()
	res, err := Generate(Spec{Fields: []FieldSpec{{Name: "x", Kind: "spaceship"}}}, GenerateOptions{Count: 1})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Diagnostics[0].Code != "FAKER.UNKNOWN_KIND" {
		t.Fatalf("code = %q", res.Diagnostics[0].Code)
	}
}

// TestFaker_DefaultJSON — C9: default format is JSON parseable as []map.
func TestFaker_DefaultJSON(t *testing.T) {
	t.Parallel()
	spec := Spec{Fields: []FieldSpec{{Name: "n", Kind: "sequence"}}}
	res, err := Generate(spec, GenerateOptions{Count: 4})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	var rows []map[string]any
	if err := json.Unmarshal(res.Output, &rows); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, res.Output)
	}
	if len(rows) != 4 {
		t.Fatalf("got %d rows", len(rows))
	}
}
