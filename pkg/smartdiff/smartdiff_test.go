package smartdiff

import (
	"strings"
	"testing"
)

// TestDiff_JSON_KeyOrderIgnored — C4: object key order does not matter.
func TestDiff_JSON_KeyOrderIgnored(t *testing.T) {
	t.Parallel()
	res, err := Diff([]byte(`{"a":1,"b":2}`), []byte(`{"b":2,"a":1}`), DiffOptions{Mode: "json"})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if len(res.Hunks) != 0 {
		t.Fatalf("expected empty diff, got %+v", res.Hunks)
	}
}

// TestDiff_JSON_ChangeAddRemove — C4: change/add/remove are emitted.
func TestDiff_JSON_ChangeAddRemove(t *testing.T) {
	t.Parallel()
	res, err := Diff([]byte(`{"a":1,"b":2}`), []byte(`{"a":99,"c":3}`), DiffOptions{Mode: "json"})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	ops := map[string]int{}
	for _, h := range res.Hunks {
		ops[h.Op]++
	}
	if ops["change"] != 1 || ops["remove"] != 1 || ops["add"] != 1 {
		t.Fatalf("got ops %+v in %+v", ops, res.Hunks)
	}
}

// TestDiff_JSON_IgnoreOrder — C4: array order option treats lists as sets.
func TestDiff_JSON_IgnoreOrder(t *testing.T) {
	t.Parallel()
	res, err := Diff([]byte(`[1,2,3]`), []byte(`[3,2,1]`), DiffOptions{Mode: "json", IgnoreOrder: true})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if len(res.Hunks) != 0 {
		t.Fatalf("expected empty, got %+v", res.Hunks)
	}
}

// TestDiff_INI_SectionRename — C4: an INI section rename surfaces add+remove.
func TestDiff_INI_SectionRename(t *testing.T) {
	t.Parallel()
	left := []byte("[old]\nfoo=1\n")
	right := []byte("[new]\nfoo=1\n")
	res, err := Diff(left, right, DiffOptions{Mode: "ini"})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	var add, rem int
	for _, h := range res.Hunks {
		if h.Op == "add" {
			add++
		}
		if h.Op == "remove" {
			rem++
		}
	}
	if add < 1 || rem < 1 {
		t.Fatalf("expected at least one add+remove: %+v", res.Hunks)
	}
}

// TestDiff_SQL_AlterTable_AddsColumn — C4: an extra ALTER statement appears as add.
func TestDiff_SQL_AlterTable_AddsColumn(t *testing.T) {
	t.Parallel()
	left := []byte("CREATE TABLE u (id INT);")
	right := []byte("CREATE TABLE u (id INT); ALTER TABLE u ADD COLUMN name TEXT;")
	res, err := Diff(left, right, DiffOptions{Mode: "sql"})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if len(res.Hunks) != 1 || res.Hunks[0].Op != "add" {
		t.Fatalf("unexpected hunks: %+v", res.Hunks)
	}
	if !strings.Contains(strings.ToUpper(toString(res.Hunks[0].Right)), "ALTER TABLE U") {
		t.Fatalf("unexpected right: %v", res.Hunks[0].Right)
	}
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// TestDiff_AutoMode_DetectsJSON — C4: omitting Mode picks json for braces.
func TestDiff_AutoMode_DetectsJSON(t *testing.T) {
	t.Parallel()
	res, err := Diff([]byte(`{"a":1}`), []byte(`{"a":2}`), DiffOptions{})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if res.Mode != "json" {
		t.Fatalf("Mode = %q", res.Mode)
	}
	if len(res.Hunks) != 1 || res.Hunks[0].Op != "change" {
		t.Fatalf("unexpected hunks: %+v", res.Hunks)
	}
}
