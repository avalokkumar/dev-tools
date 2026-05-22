package datax

import (
	"strings"
	"testing"
)

func TestJSONToCSV_Basic(t *testing.T) {
	t.Parallel()
	in := []byte(`[{"name":"alok","age":42},{"name":"bob","age":1}]`)
	r, _ := JSONToCSV(in, JSONToCSVOptions{})
	if !strings.Contains(r.Output, "age,name") || !strings.Contains(r.Output, "42,alok") {
		t.Fatalf("unexpected:\n%s", r.Output)
	}
}

func TestCSVToJSON_Basic(t *testing.T) {
	t.Parallel()
	in := []byte("name,age\nalok,42\nbob,1\n")
	r, _ := CSVToJSON(in, CSVToJSONOptions{})
	if !strings.Contains(r.Output, `"name": "alok"`) {
		t.Fatalf("unexpected:\n%s", r.Output)
	}
}

func TestJSON_CSV_Roundtrip(t *testing.T) {
	t.Parallel()
	src := []byte(`[{"a":"1","b":"2"},{"a":"3","b":"4"}]`)
	csvR, _ := JSONToCSV(src, JSONToCSVOptions{})
	back, _ := CSVToJSON([]byte(csvR.Output), CSVToJSONOptions{})
	if !strings.Contains(back.Output, `"a": "1"`) || !strings.Contains(back.Output, `"b": "4"`) {
		t.Fatalf("round-trip lost data:\n%s", back.Output)
	}
}

func TestJSONToXML_Basic(t *testing.T) {
	t.Parallel()
	r, _ := JSONToXML([]byte(`{"name":"alok","tags":["a","b"]}`), JSONToXMLOptions{Root: "user"})
	if !strings.Contains(r.Output, "<user>") {
		t.Fatalf("missing root: %s", r.Output)
	}
	if !strings.Contains(r.Output, "<name>alok</name>") {
		t.Fatalf("missing element: %s", r.Output)
	}
	if !strings.Contains(r.Output, "<item>a</item>") {
		t.Fatalf("missing array item: %s", r.Output)
	}
}

func TestXMLToJSON_Basic(t *testing.T) {
	t.Parallel()
	in := []byte(`<root><name>alok</name><age>42</age></root>`)
	r, _ := XMLToJSON(in, XMLToJSONOptions{})
	if !strings.Contains(r.Output, `"name": "alok"`) {
		t.Fatalf("unexpected:\n%s", r.Output)
	}
}

func TestFlatten_Basic(t *testing.T) {
	t.Parallel()
	r, _ := Flatten([]byte(`{"a":{"b":{"c":1}},"d":[10,20]}`), FlattenOptions{})
	out := r.Output
	if !strings.Contains(out, `"a.b.c": 1`) {
		t.Fatalf("missing dotted path: %s", out)
	}
	if !strings.Contains(out, `"d.0": 10`) {
		t.Fatalf("missing array index: %s", out)
	}
}

func TestFlatten_Unflatten_Roundtrip(t *testing.T) {
	t.Parallel()
	src := []byte(`{"a":{"b":{"c":1}},"d":[{"e":"x"}]}`)
	flat, _ := Flatten(src, FlattenOptions{})
	un, _ := Unflatten([]byte(flat.Output), FlattenOptions{})
	if !strings.Contains(un.Output, `"c": 1`) {
		t.Fatalf("round-trip drop: %s", un.Output)
	}
	if !strings.Contains(un.Output, `"e": "x"`) {
		t.Fatalf("array-of-object path lost: %s", un.Output)
	}
}

func TestKeyRename(t *testing.T) {
	t.Parallel()
	in := []byte(`{"old_key":1,"nested":{"old_key":2}}`)
	r, _ := KeyRename(in, []KeyRenameRule{{From: "old_key", To: "newKey"}})
	if !strings.Contains(r.Output, `"newKey"`) {
		t.Fatalf("rename failed:\n%s", r.Output)
	}
	if strings.Contains(r.Output, `"old_key"`) {
		t.Fatalf("old key still present:\n%s", r.Output)
	}
}

func TestJSONToCSV_BadJSON(t *testing.T) {
	t.Parallel()
	r, _ := JSONToCSV([]byte("{not-json"), JSONToCSVOptions{})
	if r.Diagnostics[0].Code != "DATA.JSON.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}
