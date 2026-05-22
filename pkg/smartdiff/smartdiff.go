// Package smartdiff produces semantic diffs for JSON, INI, and SQL inputs.
//
// External API:
//
//	Diff(left, right []byte, opts DiffOptions) (DiffResult, error)
//
// "Semantic" here means: structural diffs (per-key, per-section, per-statement)
// rather than line-by-line text diffs.
package smartdiff

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// DiffOptions tunes the diff.
type DiffOptions struct {
	// Mode is one of "json", "ini", "sql", or "auto".
	Mode string `json:"mode,omitempty"`
	// IgnoreOrder treats array order as immaterial (JSON only).
	IgnoreOrder bool `json:"ignoreOrder,omitempty"`
	// IgnoreWhitespace collapses runs of whitespace before comparison.
	IgnoreWhitespace bool `json:"ignoreWhitespace,omitempty"`
}

// Hunk is one differing element.
type Hunk struct {
	// Path locates the change. Format depends on Mode:
	//   json: dotted JSON-pointer-like path (e.g. "users.0.name")
	//   ini:  section + "." + key
	//   sql:  zero-based statement index
	Path string `json:"path"`
	// Op is "add", "remove", or "change".
	Op string `json:"op"`
	// Left and Right carry the diffing values when Op is "change".
	// They are nil for "add"/"remove" of the absent side.
	Left  any `json:"left,omitempty"`
	Right any `json:"right,omitempty"`
}

// Summary is high-level counts.
type Summary struct {
	Adds    int `json:"adds"`
	Removes int `json:"removes"`
	Changes int `json:"changes"`
}

// DiffResult is the success return.
type DiffResult struct {
	Mode        string              `json:"mode"`
	Hunks       []Hunk              `json:"hunks"`
	Summary     Summary             `json:"summary"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Diff compares two inputs.
func Diff(left, right []byte, opts DiffOptions) (DiffResult, error) {
	mode := opts.Mode
	if mode == "" || mode == "auto" {
		mode = detectMode(left, right)
	}
	switch mode {
	case "json":
		return diffJSON(left, right, opts)
	case "ini":
		return diffINI(left, right, opts)
	case "sql":
		return diffSQL(left, right, opts)
	default:
		return DiffResult{
			Mode: mode,
			Diagnostics: []engine.Diagnostic{{
				Code:     "DIFF.UNSUPPORTED_MODE",
				Message:  fmt.Sprintf("mode %q not supported", mode),
				Severity: engine.SevError,
			}},
		}, nil
	}
}

func detectMode(l, r []byte) string {
	t := bytes.TrimSpace(l)
	if len(t) == 0 {
		t = bytes.TrimSpace(r)
	}
	if len(t) == 0 {
		return "json"
	}
	if t[0] == '{' || t[0] == '[' {
		return "json"
	}
	if t[0] == '[' || strings.HasPrefix(string(t), "[") {
		return "ini"
	}
	if bytes.Contains(bytes.ToUpper(t), []byte("SELECT ")) ||
		bytes.Contains(bytes.ToUpper(t), []byte("INSERT ")) ||
		bytes.Contains(bytes.ToUpper(t), []byte("CREATE ")) ||
		bytes.Contains(bytes.ToUpper(t), []byte("ALTER ")) {
		return "sql"
	}
	return "ini"
}

// ---- JSON ----

func diffJSON(left, right []byte, opts DiffOptions) (DiffResult, error) {
	var lv, rv any
	if err := json.Unmarshal(left, &lv); err != nil {
		return DiffResult{Mode: "json", Diagnostics: []engine.Diagnostic{{
			Code: "DIFF.JSON.PARSE_LEFT", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	if err := json.Unmarshal(right, &rv); err != nil {
		return DiffResult{Mode: "json", Diagnostics: []engine.Diagnostic{{
			Code: "DIFF.JSON.PARSE_RIGHT", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	var hunks []Hunk
	walkJSON("", lv, rv, opts.IgnoreOrder, &hunks)
	return DiffResult{Mode: "json", Hunks: hunks, Summary: summarise(hunks)}, nil
}

func walkJSON(path string, l, r any, ignoreOrder bool, hunks *[]Hunk) {
	if reflect.DeepEqual(l, r) {
		return
	}
	switch lt := l.(type) {
	case map[string]any:
		rt, ok := r.(map[string]any)
		if !ok {
			*hunks = append(*hunks, Hunk{Path: path, Op: "change", Left: l, Right: r})
			return
		}
		keys := unionKeys(lt, rt)
		for _, k := range keys {
			lv, lok := lt[k]
			rv, rok := rt[k]
			subPath := joinPath(path, k)
			switch {
			case lok && !rok:
				*hunks = append(*hunks, Hunk{Path: subPath, Op: "remove", Left: lv})
			case !lok && rok:
				*hunks = append(*hunks, Hunk{Path: subPath, Op: "add", Right: rv})
			default:
				walkJSON(subPath, lv, rv, ignoreOrder, hunks)
			}
		}
	case []any:
		rt, ok := r.([]any)
		if !ok {
			*hunks = append(*hunks, Hunk{Path: path, Op: "change", Left: l, Right: r})
			return
		}
		if ignoreOrder {
			diffArraysSet(path, lt, rt, hunks)
			return
		}
		max := len(lt)
		if len(rt) > max {
			max = len(rt)
		}
		for i := 0; i < max; i++ {
			subPath := joinPath(path, fmt.Sprint(i))
			switch {
			case i >= len(lt):
				*hunks = append(*hunks, Hunk{Path: subPath, Op: "add", Right: rt[i]})
			case i >= len(rt):
				*hunks = append(*hunks, Hunk{Path: subPath, Op: "remove", Left: lt[i]})
			default:
				walkJSON(subPath, lt[i], rt[i], ignoreOrder, hunks)
			}
		}
	default:
		*hunks = append(*hunks, Hunk{Path: path, Op: "change", Left: l, Right: r})
	}
}

func diffArraysSet(path string, l, r []any, hunks *[]Hunk) {
	used := make([]bool, len(r))
	for _, lv := range l {
		matched := false
		for i, rv := range r {
			if !used[i] && reflect.DeepEqual(lv, rv) {
				used[i] = true
				matched = true
				break
			}
		}
		if !matched {
			*hunks = append(*hunks, Hunk{Path: path, Op: "remove", Left: lv})
		}
	}
	for i, rv := range r {
		if !used[i] {
			*hunks = append(*hunks, Hunk{Path: path, Op: "add", Right: rv})
		}
	}
}

func unionKeys(a, b map[string]any) []string {
	seen := map[string]struct{}{}
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}

// ---- INI ----

func diffINI(left, right []byte, opts DiffOptions) (DiffResult, error) {
	lm := parseINI(left, opts.IgnoreWhitespace)
	rm := parseINI(right, opts.IgnoreWhitespace)
	var hunks []Hunk
	sections := unionINISections(lm, rm)
	for _, sec := range sections {
		lkv := lm[sec]
		rkv := rm[sec]
		keys := unionINIKeys(lkv, rkv)
		for _, k := range keys {
			lv, lok := lkv[k]
			rv, rok := rkv[k]
			path := sec + "." + k
			switch {
			case lok && !rok:
				hunks = append(hunks, Hunk{Path: path, Op: "remove", Left: lv})
			case !lok && rok:
				hunks = append(hunks, Hunk{Path: path, Op: "add", Right: rv})
			case lv != rv:
				hunks = append(hunks, Hunk{Path: path, Op: "change", Left: lv, Right: rv})
			}
		}
	}
	return DiffResult{Mode: "ini", Hunks: hunks, Summary: summarise(hunks)}, nil
}

func parseINI(b []byte, collapseWS bool) map[string]map[string]string {
	out := map[string]map[string]string{}
	cur := ""
	out[cur] = map[string]string{}
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			cur = strings.TrimSpace(line[1 : len(line)-1])
			if _, ok := out[cur]; !ok {
				out[cur] = map[string]string{}
			}
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		if collapseWS {
			val = strings.Join(strings.Fields(val), " ")
		}
		out[cur][key] = val
	}
	return out
}

func unionINISections(a, b map[string]map[string]string) []string {
	seen := map[string]struct{}{}
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func unionINIKeys(a, b map[string]string) []string {
	seen := map[string]struct{}{}
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ---- SQL ----

// Statement-level diff: split on ';' (naive but adequate for MVP), normalise
// whitespace+case, then compare per index.
func diffSQL(left, right []byte, opts DiffOptions) (DiffResult, error) {
	ls := splitSQL(string(left), opts.IgnoreWhitespace)
	rs := splitSQL(string(right), opts.IgnoreWhitespace)
	var hunks []Hunk
	max := len(ls)
	if len(rs) > max {
		max = len(rs)
	}
	for i := 0; i < max; i++ {
		switch {
		case i >= len(ls):
			hunks = append(hunks, Hunk{Path: fmt.Sprint(i), Op: "add", Right: rs[i]})
		case i >= len(rs):
			hunks = append(hunks, Hunk{Path: fmt.Sprint(i), Op: "remove", Left: ls[i]})
		case ls[i] != rs[i]:
			hunks = append(hunks, Hunk{Path: fmt.Sprint(i), Op: "change", Left: ls[i], Right: rs[i]})
		}
	}
	return DiffResult{Mode: "sql", Hunks: hunks, Summary: summarise(hunks)}, nil
}

func splitSQL(s string, collapseWS bool) []string {
	parts := strings.Split(s, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		p = strings.ToUpper(p)
		if collapseWS {
			p = strings.Join(strings.Fields(p), " ")
		}
		out = append(out, p)
	}
	return out
}

func summarise(h []Hunk) Summary {
	var s Summary
	for _, x := range h {
		switch x.Op {
		case "add":
			s.Adds++
		case "remove":
			s.Removes++
		case "change":
			s.Changes++
		}
	}
	return s
}
