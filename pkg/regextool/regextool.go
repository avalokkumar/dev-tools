// Package regextool tests regular expressions and produces plain-English
// explanations of pattern fragments.
//
// External API:
//
//	Test(pattern, input string, opts TestOptions) (TestResult, error)
//	Explain(pattern string, opts ExplainOptions) (ExplainResult, error)
//
// Only Go's RE2 flavor is supported in MVP. PCRE is reserved for Phase D.
package regextool

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// TestOptions tunes Test.
type TestOptions struct {
	// Flavor must be "" or "re2" in MVP.
	Flavor string `json:"flavor,omitempty"`
	// Flags is a set of single-letter modifiers: i (case-insensitive), m
	// (multi-line), s (dot matches newline). RE2 does not support 'g'.
	Flags string `json:"flags,omitempty"`
}

// Group is one capture group within a Match.
type Group struct {
	Name  string `json:"name,omitempty"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Value string `json:"value"`
}

// Match is one match of pattern against input.
type Match struct {
	Start  int     `json:"start"`
	End    int     `json:"end"`
	Value  string  `json:"value"`
	Groups []Group `json:"groups,omitempty"`
}

// TestResult is the success return.
type TestResult struct {
	Matches     []Match             `json:"matches"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Test compiles pattern and reports all matches in input.
func Test(pattern, input string, opts TestOptions) (TestResult, error) {
	if opts.Flavor != "" && opts.Flavor != "re2" {
		return TestResult{Diagnostics: []engine.Diagnostic{{
			Code: "REGEX.UNSUPPORTED_FLAVOR", Message: "only re2 is supported in MVP",
			Severity: engine.SevError,
		}}}, nil
	}
	flagPrefix := buildFlagPrefix(opts.Flags)
	re, err := regexp.Compile(flagPrefix + pattern)
	if err != nil {
		return TestResult{Diagnostics: []engine.Diagnostic{{
			Code: "REGEX.INVALID_PATTERN", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}

	matches := re.FindAllStringSubmatchIndex(input, -1)
	out := make([]Match, 0, len(matches))
	groupNames := re.SubexpNames()
	for _, idx := range matches {
		m := Match{Start: idx[0], End: idx[1], Value: input[idx[0]:idx[1]]}
		// indices come in pairs (start,end) per group, [0] is the full match
		for g := 1; g*2+1 < len(idx); g++ {
			s, e := idx[g*2], idx[g*2+1]
			if s < 0 || e < 0 {
				continue
			}
			m.Groups = append(m.Groups, Group{
				Name: groupNames[g], Start: s, End: e, Value: input[s:e],
			})
		}
		out = append(out, m)
	}
	return TestResult{Matches: out}, nil
}

func buildFlagPrefix(flags string) string {
	flags = strings.ToLower(flags)
	var b strings.Builder
	for _, r := range flags {
		switch r {
		case 'i', 'm', 's':
			// supported
		default:
			continue
		}
	}
	keep := strings.Map(func(r rune) rune {
		switch r {
		case 'i', 'm', 's':
			return r
		}
		return -1
	}, flags)
	if keep == "" {
		return ""
	}
	b.WriteString("(?")
	b.WriteString(keep)
	b.WriteString(")")
	return b.String()
}

// ExplainOptions reserves space for future tuning.
type ExplainOptions struct {
	Flavor string `json:"flavor,omitempty"`
}

// ExplainNode is one tokenised regex fragment with a plain-English label.
type ExplainNode struct {
	Token       string `json:"token"`
	Description string `json:"description"`
}

// ExplainResult is the success return.
type ExplainResult struct {
	Tree        []ExplainNode       `json:"tree"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Explain returns a token-by-token annotation of the pattern. This is a
// pragmatic best-effort explainer for the most common metacharacters; it does
// not fully parse the regex grammar.
func Explain(pattern string, _ ExplainOptions) (ExplainResult, error) {
	if pattern == "" {
		return ExplainResult{Diagnostics: []engine.Diagnostic{{
			Code: "REGEX.EMPTY_PATTERN", Message: "pattern is empty",
			Severity: engine.SevError,
		}}}, nil
	}
	if _, err := regexp.Compile(pattern); err != nil {
		return ExplainResult{Diagnostics: []engine.Diagnostic{{
			Code: "REGEX.INVALID_PATTERN", Message: err.Error(),
			Severity: engine.SevError,
		}}}, nil
	}
	out := make([]ExplainNode, 0, len(pattern))
	runes := []rune(pattern)
	i := 0
	for i < len(runes) {
		r := runes[i]
		var node ExplainNode
		switch r {
		case '^':
			node = ExplainNode{Token: "^", Description: "start of line/input"}
		case '$':
			node = ExplainNode{Token: "$", Description: "end of line/input"}
		case '.':
			node = ExplainNode{Token: ".", Description: "any single character (except newline)"}
		case '*':
			node = ExplainNode{Token: "*", Description: "zero or more of preceding"}
		case '+':
			node = ExplainNode{Token: "+", Description: "one or more of preceding"}
		case '?':
			node = ExplainNode{Token: "?", Description: "optional / zero or one"}
		case '|':
			node = ExplainNode{Token: "|", Description: "alternation (or)"}
		case '(':
			// Try to identify named groups.
			if i+2 < len(runes) && runes[i+1] == '?' && runes[i+2] == 'P' {
				node = ExplainNode{Token: "(?P<…>…)", Description: "named capturing group"}
			} else if i+1 < len(runes) && runes[i+1] == '?' {
				node = ExplainNode{Token: "(?…)", Description: "non-capturing or flagged group"}
			} else {
				node = ExplainNode{Token: "(", Description: "open capturing group"}
			}
		case ')':
			node = ExplainNode{Token: ")", Description: "close group"}
		case '[':
			// Capture entire char class up to closing ']'.
			j := i + 1
			for j < len(runes) && runes[j] != ']' {
				j++
			}
			if j < len(runes) {
				node = ExplainNode{Token: string(runes[i : j+1]), Description: "character class"}
				i = j
			} else {
				node = ExplainNode{Token: "[", Description: "open character class"}
			}
		case '\\':
			if i+1 < len(runes) {
				next := runes[i+1]
				node = ExplainNode{
					Token: "\\" + string(next), Description: explainEscape(next),
				}
				i++
			}
		default:
			node = ExplainNode{Token: string(r), Description: fmt.Sprintf("literal %q", r)}
		}
		out = append(out, node)
		i++
	}
	return ExplainResult{Tree: out}, nil
}

func explainEscape(r rune) string {
	switch r {
	case 'd':
		return "any digit (0-9)"
	case 'D':
		return "any non-digit"
	case 'w':
		return "word char (alpha/digit/underscore)"
	case 'W':
		return "non-word char"
	case 's':
		return "whitespace"
	case 'S':
		return "non-whitespace"
	case 'b':
		return "word boundary"
	case 'B':
		return "non-word-boundary"
	case 'n':
		return "newline"
	case 't':
		return "tab"
	case 'r':
		return "carriage return"
	}
	return fmt.Sprintf("escaped literal %q", r)
}
