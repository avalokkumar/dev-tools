// Package datax provides JSON↔CSV↔XML transforms plus flatten/unflatten/key
// rename for JSON.
//
// External API:
//
//	JSONToCSV([]byte, JSONToCSVOptions) (BytesResult, error)
//	CSVToJSON([]byte, CSVToJSONOptions) (BytesResult, error)
//	JSONToXML([]byte, JSONToXMLOptions) (BytesResult, error)
//	XMLToJSON([]byte, XMLToJSONOptions) (BytesResult, error)
//	Flatten([]byte, FlattenOptions) (BytesResult, error)
//	Unflatten([]byte, FlattenOptions) (BytesResult, error)
//	KeyRename([]byte, []KeyRenameRule) (BytesResult, error)
package datax

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// BytesResult holds the transformed output (UTF-8 string).
type BytesResult struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// ---------- JSON ↔ CSV ----------

// JSONToCSVOptions tunes JSONToCSV.
type JSONToCSVOptions struct {
	Delimiter string `json:"delimiter,omitempty"`
}

// JSONToCSV converts an array of objects to CSV. The set of columns is the
// union of all object keys (alphabetically sorted for stable output).
func JSONToCSV(input []byte, opts JSONToCSVOptions) (BytesResult, error) {
	var rows []map[string]any
	if err := json.Unmarshal(input, &rows); err != nil {
		return BytesResult{Diagnostics: []engine.Diagnostic{{
			Code: "DATA.JSON.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	delim := ','
	if opts.Delimiter != "" {
		for _, r := range opts.Delimiter {
			delim = r
			break
		}
	}
	cols := unionKeys(rows)
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = delim
	_ = w.Write(cols)
	for _, row := range rows {
		fields := make([]string, len(cols))
		for i, k := range cols {
			fields[i] = stringify(row[k])
		}
		_ = w.Write(fields)
	}
	w.Flush()
	return BytesResult{Output: buf.String()}, nil
}

func unionKeys(rows []map[string]any) []string {
	seen := map[string]struct{}{}
	for _, r := range rows {
		for k := range r {
			seen[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func stringify(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		// Trim trailing .0 so integer-valued floats serialise cleanly.
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// CSVToJSONOptions tunes CSVToJSON.
type CSVToJSONOptions struct {
	Delimiter string `json:"delimiter,omitempty"`
	HasHeader *bool  `json:"hasHeader,omitempty"`
}

// CSVToJSON converts a CSV input to an array of objects.
func CSVToJSON(input []byte, opts CSVToJSONOptions) (BytesResult, error) {
	delim := ','
	if opts.Delimiter != "" {
		for _, r := range opts.Delimiter {
			delim = r
			break
		}
	}
	hasHeader := true
	if opts.HasHeader != nil {
		hasHeader = *opts.HasHeader
	}
	r := csv.NewReader(bytes.NewReader(input))
	r.Comma = delim
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil {
		return BytesResult{Diagnostics: []engine.Diagnostic{{
			Code: "DATA.CSV.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	if len(rows) == 0 {
		return BytesResult{Output: "[]"}, nil
	}
	var header []string
	body := rows
	if hasHeader {
		header = rows[0]
		body = rows[1:]
	} else {
		header = make([]string, len(rows[0]))
		for i := range header {
			header[i] = fmt.Sprintf("col%d", i+1)
		}
	}
	out := make([]map[string]string, 0, len(body))
	for _, row := range body {
		obj := make(map[string]string, len(header))
		for i, k := range header {
			if i < len(row) {
				obj[k] = row[i]
			}
		}
		out = append(out, obj)
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return BytesResult{}, fmt.Errorf("datax: json: %w", err)
	}
	return BytesResult{Output: string(b)}, nil
}

// ---------- JSON ↔ XML ----------

// JSONToXMLOptions tunes JSONToXML.
type JSONToXMLOptions struct {
	// Root is the wrapper element name. Default "root".
	Root string `json:"root,omitempty"`
	// Indent is the per-level indent. Default 2.
	Indent int `json:"indent,omitempty"`
}

// JSONToXML converts JSON to XML. Maps become elements, arrays become
// repeated <item>, scalars become text content.
func JSONToXML(input []byte, opts JSONToXMLOptions) (BytesResult, error) {
	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		return BytesResult{Diagnostics: []engine.Diagnostic{{
			Code: "DATA.JSON.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	root := opts.Root
	if root == "" {
		root = "root"
	}
	indent := opts.Indent
	if indent <= 0 {
		indent = 2
	}
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	writeXML(&b, root, v, 0, indent)
	return BytesResult{Output: b.String()}, nil
}

func writeXML(b *strings.Builder, name string, v any, depth, indent int) {
	pad := strings.Repeat(" ", depth*indent)
	switch t := v.(type) {
	case map[string]any:
		fmt.Fprintf(b, "%s<%s>\n", pad, name)
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			writeXML(b, sanitiseTag(k), t[k], depth+1, indent)
		}
		fmt.Fprintf(b, "%s</%s>\n", pad, name)
	case []any:
		fmt.Fprintf(b, "%s<%s>\n", pad, name)
		for _, x := range t {
			writeXML(b, "item", x, depth+1, indent)
		}
		fmt.Fprintf(b, "%s</%s>\n", pad, name)
	default:
		fmt.Fprintf(b, "%s<%s>%s</%s>\n", pad, name, xmlEscape(stringify(v)), name)
	}
}

func sanitiseTag(s string) string {
	var b strings.Builder
	for i, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' ||
			(i > 0 && (r >= '0' && r <= '9'))
		if !ok {
			b.WriteRune('_')
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if out == "" {
		return "field"
	}
	return out
}

func xmlEscape(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// XMLToJSONOptions tunes XMLToJSON.
type XMLToJSONOptions struct {
	Indent int `json:"indent,omitempty"`
}

// XMLToJSON parses XML and re-emits it as JSON. Repeated child tags become arrays.
func XMLToJSON(input []byte, opts XMLToJSONOptions) (BytesResult, error) {
	tree, err := parseXMLTree(input)
	if err != nil {
		return BytesResult{Diagnostics: []engine.Diagnostic{{
			Code: "DATA.XML.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	indent := opts.Indent
	if indent <= 0 {
		indent = 2
	}
	pad := strings.Repeat(" ", indent)
	out, err := json.MarshalIndent(tree, "", pad)
	if err != nil {
		return BytesResult{}, fmt.Errorf("datax: json: %w", err)
	}
	return BytesResult{Output: string(out)}, nil
}

func parseXMLTree(input []byte) (any, error) {
	dec := xml.NewDecoder(bytes.NewReader(input))
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if start, ok := tok.(xml.StartElement); ok {
			obj := map[string]any{}
			obj[start.Name.Local] = readElement(dec, start)
			return obj, nil
		}
	}
}

func readElement(dec *xml.Decoder, start xml.StartElement) any {
	obj := map[string]any{}
	var text strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			child := readElement(dec, t)
			name := t.Name.Local
			if existing, ok := obj[name]; ok {
				if arr, isArr := existing.([]any); isArr {
					obj[name] = append(arr, child)
				} else {
					obj[name] = []any{existing, child}
				}
			} else {
				obj[name] = child
			}
		case xml.CharData:
			text.WriteString(string(t))
		case xml.EndElement:
			if t.Name == start.Name {
				if len(obj) == 0 {
					s := strings.TrimSpace(text.String())
					if s == "" {
						return nil
					}
					return s
				}
				s := strings.TrimSpace(text.String())
				if s != "" {
					obj["#text"] = s
				}
				return obj
			}
		}
	}
	return obj
}

// ---------- Flatten / Unflatten ----------

// FlattenOptions tunes Flatten/Unflatten.
type FlattenOptions struct {
	// Sep is the dotted-path separator. Default ".".
	Sep string `json:"sep,omitempty"`
}

// Flatten turns nested JSON into a single-level map keyed by dotted paths.
func Flatten(input []byte, opts FlattenOptions) (BytesResult, error) {
	sep := opts.Sep
	if sep == "" {
		sep = "."
	}
	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		return BytesResult{Diagnostics: []engine.Diagnostic{{
			Code: "DATA.JSON.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	flat := map[string]any{}
	flattenInto("", v, flat, sep)
	out, err := json.MarshalIndent(flat, "", "  ")
	if err != nil {
		return BytesResult{}, fmt.Errorf("datax: json: %w", err)
	}
	return BytesResult{Output: string(out)}, nil
}

func flattenInto(prefix string, v any, out map[string]any, sep string) {
	switch t := v.(type) {
	case map[string]any:
		if len(t) == 0 && prefix != "" {
			out[prefix] = map[string]any{}
			return
		}
		for k, vv := range t {
			next := k
			if prefix != "" {
				next = prefix + sep + k
			}
			flattenInto(next, vv, out, sep)
		}
	case []any:
		if len(t) == 0 && prefix != "" {
			out[prefix] = []any{}
			return
		}
		for i, vv := range t {
			next := fmt.Sprintf("%s%s%d", prefix, sep, i)
			if prefix == "" {
				next = strconv.Itoa(i)
			}
			flattenInto(next, vv, out, sep)
		}
	default:
		out[prefix] = v
	}
}

// Unflatten is the inverse of Flatten.
func Unflatten(input []byte, opts FlattenOptions) (BytesResult, error) {
	sep := opts.Sep
	if sep == "" {
		sep = "."
	}
	var flat map[string]any
	if err := json.Unmarshal(input, &flat); err != nil {
		return BytesResult{Diagnostics: []engine.Diagnostic{{
			Code: "DATA.JSON.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	var root any = map[string]any{}
	for path, v := range flat {
		root = setPath(root, strings.Split(path, sep), v)
	}
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return BytesResult{}, fmt.Errorf("datax: json: %w", err)
	}
	return BytesResult{Output: string(out)}, nil
}

func setPath(root any, parts []string, v any) any {
	if len(parts) == 0 {
		return v
	}
	head := parts[0]
	rest := parts[1:]
	if idx, err := strconv.Atoi(head); err == nil {
		// Array index.
		arr, ok := root.([]any)
		if !ok {
			arr = []any{}
		}
		for len(arr) <= idx {
			arr = append(arr, nil)
		}
		arr[idx] = setPath(arr[idx], rest, v)
		return arr
	}
	m, ok := root.(map[string]any)
	if !ok {
		m = map[string]any{}
	}
	m[head] = setPath(m[head], rest, v)
	return m
}

// ---------- Key rename ----------

// KeyRenameRule maps source paths/keys to new keys.
type KeyRenameRule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// KeyRename rewrites top-level (and nested object) keys per rules. Supports
// only exact key match for MVP; pattern-based rules are deferred.
func KeyRename(input []byte, rules []KeyRenameRule) (BytesResult, error) {
	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		return BytesResult{Diagnostics: []engine.Diagnostic{{
			Code: "DATA.JSON.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	renamed := renameKeys(v, rules)
	out, err := json.MarshalIndent(renamed, "", "  ")
	if err != nil {
		return BytesResult{}, fmt.Errorf("datax: json: %w", err)
	}
	return BytesResult{Output: string(out)}, nil
}

func renameKeys(v any, rules []KeyRenameRule) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, vv := range t {
			newK := k
			for _, r := range rules {
				if r.From == k {
					newK = r.To
					break
				}
			}
			out[newK] = renameKeys(vv, rules)
		}
		return out
	case []any:
		for i, vv := range t {
			t[i] = renameKeys(vv, rules)
		}
		return t
	}
	return v
}
