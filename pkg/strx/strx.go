// Package strx provides text/string utilities.
//
// External API:
//
//	Case(input, mode) (StringResult, error)
//	Diff(left, right, DiffOptions) (DiffResult, error)
//	Stats(input) (StatsResult, error)
//	SortUnique(input, SortOptions) (StringResult, error)
//	Replace(input, pattern, repl, ReplaceOptions) (StringResult, error)
package strx

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/devforge/devforge/pkg/engine"
)

// StringResult is the standard string-out result.
type StringResult struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// ---------- Case conversion ----------

// Case modes: camel, snake, kebab, pascal, constant, dot, train, header, lower, upper, title.
func Case(input, mode string) (StringResult, error) {
	parts := splitWords(input)
	switch strings.ToLower(mode) {
	case "camel":
		return StringResult{Output: joinCamel(parts, false)}, nil
	case "pascal":
		return StringResult{Output: joinCamel(parts, true)}, nil
	case "snake":
		return StringResult{Output: strings.ToLower(strings.Join(parts, "_"))}, nil
	case "kebab":
		return StringResult{Output: strings.ToLower(strings.Join(parts, "-"))}, nil
	case "constant", "screaming_snake":
		return StringResult{Output: strings.ToUpper(strings.Join(parts, "_"))}, nil
	case "dot":
		return StringResult{Output: strings.ToLower(strings.Join(parts, "."))}, nil
	case "train", "header":
		return StringResult{Output: joinTitle(parts, "-")}, nil
	case "title":
		return StringResult{Output: joinTitle(parts, " ")}, nil
	case "lower":
		return StringResult{Output: strings.ToLower(input)}, nil
	case "upper":
		return StringResult{Output: strings.ToUpper(input)}, nil
	default:
		return StringResult{Diagnostics: []engine.Diagnostic{{
			Code: "STR.CASE.UNKNOWN_MODE",
			Message: fmt.Sprintf("unknown case mode %q (try camel|snake|kebab|pascal|constant|dot|train|title|lower|upper)", mode),
			Severity: engine.SevError,
		}}}, nil
	}
}

// splitWords breaks an arbitrary identifier or phrase into lowercase words.
func splitWords(s string) []string {
	if s == "" {
		return nil
	}
	// Insert separators at case boundaries first.
	var b strings.Builder
	prev := ' '
	for i, r := range s {
		switch {
		case unicode.IsSpace(r), r == '_', r == '-', r == '.', r == '/':
			b.WriteRune(' ')
		case i > 0 && unicode.IsUpper(r) && unicode.IsLower(prev):
			b.WriteRune(' ')
			b.WriteRune(r)
		case i > 0 && unicode.IsLetter(r) && unicode.IsDigit(prev):
			b.WriteRune(' ')
			b.WriteRune(r)
		case i > 0 && unicode.IsDigit(r) && unicode.IsLetter(prev):
			b.WriteRune(' ')
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
		prev = r
	}
	out := []string{}
	for _, w := range strings.Fields(b.String()) {
		out = append(out, strings.ToLower(w))
	}
	return out
}

func joinCamel(parts []string, pascal bool) string {
	var b strings.Builder
	for i, p := range parts {
		if i == 0 && !pascal {
			b.WriteString(p)
			continue
		}
		b.WriteString(titleCase(p))
	}
	return b.String()
}

func joinTitle(parts []string, sep string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, titleCase(p))
	}
	return strings.Join(out, sep)
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// ---------- Line diff ----------

// DiffOptions tunes Diff.
type DiffOptions struct {
	IgnoreWhitespace bool `json:"ignoreWhitespace,omitempty"`
	IgnoreCase       bool `json:"ignoreCase,omitempty"`
}

// DiffResult is a unified-style line diff.
type DiffResult struct {
	Hunks       []DiffHunk          `json:"hunks"`
	Summary     DiffSummary         `json:"summary"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// DiffHunk is a single change.
type DiffHunk struct {
	Op       string `json:"op"`       // "add" | "remove" | "equal"
	LineNum  int    `json:"lineNum"`  // 1-based line number on its side
	Content  string `json:"content"`
}

// DiffSummary summarises hunk counts.
type DiffSummary struct {
	Adds    int `json:"adds"`
	Removes int `json:"removes"`
}

// Diff computes a Myers-style line diff using a simple LCS algorithm.
func Diff(left, right string, opts DiffOptions) (DiffResult, error) {
	a := strings.Split(left, "\n")
	b := strings.Split(right, "\n")
	norm := func(s string) string {
		if opts.IgnoreCase {
			s = strings.ToLower(s)
		}
		if opts.IgnoreWhitespace {
			s = strings.Join(strings.Fields(s), " ")
		}
		return s
	}
	// LCS table.
	la, lb := len(a), len(b)
	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
	}
	for i := la - 1; i >= 0; i-- {
		for j := lb - 1; j >= 0; j-- {
			if norm(a[i]) == norm(b[j]) {
				dp[i][j] = dp[i+1][j+1] + 1
			} else {
				dp[i][j] = max(dp[i+1][j], dp[i][j+1])
			}
		}
	}
	var hunks []DiffHunk
	i, j := 0, 0
	var summary DiffSummary
	for i < la && j < lb {
		if norm(a[i]) == norm(b[j]) {
			hunks = append(hunks, DiffHunk{Op: "equal", LineNum: i + 1, Content: a[i]})
			i++
			j++
		} else if dp[i+1][j] >= dp[i][j+1] {
			hunks = append(hunks, DiffHunk{Op: "remove", LineNum: i + 1, Content: a[i]})
			summary.Removes++
			i++
		} else {
			hunks = append(hunks, DiffHunk{Op: "add", LineNum: j + 1, Content: b[j]})
			summary.Adds++
			j++
		}
	}
	for ; i < la; i++ {
		hunks = append(hunks, DiffHunk{Op: "remove", LineNum: i + 1, Content: a[i]})
		summary.Removes++
	}
	for ; j < lb; j++ {
		hunks = append(hunks, DiffHunk{Op: "add", LineNum: j + 1, Content: b[j]})
		summary.Adds++
	}
	return DiffResult{Hunks: hunks, Summary: summary}, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ---------- Stats ----------

// StatsResult holds text statistics.
type StatsResult struct {
	Lines       int `json:"lines"`
	Words       int `json:"words"`
	Characters  int `json:"characters"`
	Bytes       int `json:"bytes"`
	LongestLine int `json:"longestLine"`
}

// Stats counts lines, words, characters, bytes, and the longest line length.
func Stats(input string) (StatsResult, error) {
	r := StatsResult{Bytes: len(input), Characters: utf8.RuneCountInString(input)}
	if input == "" {
		return r, nil
	}
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		r.Lines++
		r.Words += len(strings.Fields(line))
		if l := utf8.RuneCountInString(line); l > r.LongestLine {
			r.LongestLine = l
		}
	}
	// Trailing newline counts as a final empty line in some tools; we follow `wc -l`
	// which counts only newlines. Adjust: Lines = number of '\n' OR scan count.
	r.Lines = strings.Count(input, "\n")
	if !strings.HasSuffix(input, "\n") {
		r.Lines++ // last line without trailing newline still counts
	}
	return r, nil
}

// ---------- Sort + unique ----------

// SortOptions tunes SortUnique.
type SortOptions struct {
	IgnoreCase bool `json:"ignoreCase,omitempty"`
	Reverse    bool `json:"reverse,omitempty"`
	Unique     bool `json:"unique,omitempty"`
}

// SortUnique sorts lines and (optionally) drops duplicates.
func SortUnique(input string, opts SortOptions) (StringResult, error) {
	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	cmp := func(i, j int) bool {
		a, b := lines[i], lines[j]
		if opts.IgnoreCase {
			a, b = strings.ToLower(a), strings.ToLower(b)
		}
		if opts.Reverse {
			return a > b
		}
		return a < b
	}
	sort.SliceStable(lines, cmp)
	if opts.Unique {
		out := make([]string, 0, len(lines))
		seen := map[string]struct{}{}
		for _, l := range lines {
			k := l
			if opts.IgnoreCase {
				k = strings.ToLower(l)
			}
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			out = append(out, l)
		}
		lines = out
	}
	return StringResult{Output: strings.Join(lines, "\n")}, nil
}

// ---------- Find + replace ----------

// ReplaceOptions tunes Replace.
type ReplaceOptions struct {
	Regex      bool `json:"regex,omitempty"`
	IgnoreCase bool `json:"ignoreCase,omitempty"`
}

// Replace performs literal or regex find-and-replace.
func Replace(input, pattern, repl string, opts ReplaceOptions) (StringResult, error) {
	if !opts.Regex {
		if opts.IgnoreCase {
			lower := strings.ToLower(input)
			needle := strings.ToLower(pattern)
			var b strings.Builder
			i := 0
			for i < len(input) {
				idx := strings.Index(lower[i:], needle)
				if idx < 0 {
					b.WriteString(input[i:])
					break
				}
				b.WriteString(input[i : i+idx])
				b.WriteString(repl)
				i += idx + len(pattern)
			}
			return StringResult{Output: b.String()}, nil
		}
		return StringResult{Output: strings.ReplaceAll(input, pattern, repl)}, nil
	}
	flags := ""
	if opts.IgnoreCase {
		flags = "(?i)"
	}
	re, err := regexp.Compile(flags + pattern)
	if err != nil {
		return StringResult{Diagnostics: []engine.Diagnostic{{
			Code: "STR.REPLACE.INVALID_REGEX", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	return StringResult{Output: re.ReplaceAllString(input, repl)}, nil
}
