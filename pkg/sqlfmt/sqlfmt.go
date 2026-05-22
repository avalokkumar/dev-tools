// Package sqlfmt provides a hand-rolled SQL formatter and a lightweight
// validator. The formatter does not depend on a full parser; it tokenises
// keywords case-insensitively and re-indents around major clauses. This
// keeps the engine pure-Go with zero third-party deps and gives stable
// output on the most common query shapes (SELECT, INSERT, UPDATE, DELETE,
// CREATE/ALTER/DROP TABLE).
//
// External API:
//
//	Format(string, FormatOptions) (FormatResult, error)
//	Validate(string, ValidateOptions) (ValidateResult, error)
package sqlfmt

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/devforge/devforge/pkg/engine"
)

// FormatOptions tunes Format.
type FormatOptions struct {
	// Indent is the number of spaces per indent level. Default 2.
	Indent int `json:"indent,omitempty"`
	// Uppercase keywords. Default true (idiomatic SQL style).
	Uppercase *bool `json:"uppercase,omitempty"`
}

// FormatResult is the formatter's output.
type FormatResult struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// keywordSet is the set of SQL words we recognise for casing + line-breaks.
var keywordSet = map[string]struct{}{}

// majorClauses cause a newline + zero-indent line.
var majorClauses = []string{
	"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "HAVING",
	"LIMIT", "OFFSET", "UNION", "UNION ALL", "INTERSECT", "EXCEPT",
	"INSERT INTO", "VALUES", "UPDATE", "SET", "DELETE FROM",
	"CREATE TABLE", "CREATE INDEX", "CREATE VIEW",
	"ALTER TABLE", "DROP TABLE", "DROP INDEX",
	"LEFT JOIN", "RIGHT JOIN", "INNER JOIN", "FULL JOIN", "CROSS JOIN", "JOIN",
	"ON", "AND", "OR", "RETURNING",
}

// minorKeywords just get re-cased.
var minorKeywords = []string{
	"AS", "ASC", "DESC", "DISTINCT", "ALL", "IN", "NOT", "NULL",
	"IS", "BETWEEN", "LIKE", "ILIKE", "EXISTS", "CASE", "WHEN", "THEN",
	"ELSE", "END", "TRUE", "FALSE", "USING", "WITH", "RECURSIVE",
	"INT", "INTEGER", "BIGINT", "SMALLINT", "VARCHAR", "TEXT", "CHAR",
	"BOOLEAN", "DATE", "TIME", "TIMESTAMP", "SERIAL", "PRIMARY", "KEY",
	"FOREIGN", "REFERENCES", "DEFAULT", "CHECK", "UNIQUE", "INDEX",
	"COLUMN", "CONSTRAINT", "CASCADE", "RESTRICT",
}

func init() {
	for _, k := range majorClauses {
		keywordSet[k] = struct{}{}
	}
	for _, k := range minorKeywords {
		keywordSet[k] = struct{}{}
	}
}

// Format reformats SQL.
func Format(input string, opts FormatOptions) (FormatResult, error) {
	if opts.Indent <= 0 {
		opts.Indent = 2
	}
	upper := true
	if opts.Uppercase != nil {
		upper = *opts.Uppercase
	}
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return FormatResult{Diagnostics: []engine.Diagnostic{{
			Code: "SQL.EMPTY", Message: "input has no SQL tokens", Severity: engine.SevError,
		}}}, nil
	}
	out := render(tokens, opts.Indent, upper)
	return FormatResult{Output: out}, nil
}

// ValidateOptions tunes Validate.
type ValidateOptions struct{}

// ValidateResult holds lint findings.
type ValidateResult struct {
	Valid       bool                `json:"valid"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Validate runs cheap syntactic + best-practice lints.
func Validate(input string, _ ValidateOptions) (ValidateResult, error) {
	res := ValidateResult{Valid: true}
	tokens := tokenize(input)
	if len(tokens) == 0 {
		res.Valid = false
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "SQL.EMPTY", Message: "no tokens", Severity: engine.SevError,
		})
		return res, nil
	}
	// Balanced parentheses.
	depth := 0
	for _, t := range tokens {
		if t.kind == tokParen {
			if t.value == "(" {
				depth++
			} else {
				depth--
				if depth < 0 {
					res.Valid = false
					res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
						Code: "SQL.UNBALANCED_PARENS", Message: "unmatched ')'", Severity: engine.SevError,
					})
					break
				}
			}
		}
	}
	if depth > 0 {
		res.Valid = false
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "SQL.UNBALANCED_PARENS",
			Message: fmt.Sprintf("%d unmatched '('", depth),
			Severity: engine.SevError,
		})
	}
	// Best-practice lints (warnings).
	upperJoined := strings.ToUpper(input)
	if strings.Contains(upperJoined, "SELECT *") {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "SQL.LINT.SELECT_STAR",
			Message: "SELECT * disallows future-proof column changes; prefer explicit columns",
			Severity: engine.SevWarn,
		})
	}
	if strings.Contains(upperJoined, "DELETE FROM") && !strings.Contains(upperJoined, "WHERE") {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "SQL.LINT.DELETE_NO_WHERE",
			Message: "DELETE without WHERE deletes every row",
			Severity: engine.SevWarn,
		})
	}
	if strings.Contains(upperJoined, "UPDATE") && strings.Contains(upperJoined, " SET ") &&
		!strings.Contains(upperJoined, "WHERE") {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "SQL.LINT.UPDATE_NO_WHERE",
			Message: "UPDATE without WHERE updates every row",
			Severity: engine.SevWarn,
		})
	}
	return res, nil
}

// ---- internal tokenizer ----

type tokKind int

const (
	tokWord tokKind = iota
	tokNumber
	tokString
	tokParen
	tokSymbol
	tokSemi
	tokComment
)

type token struct {
	kind  tokKind
	value string
}

func tokenize(in string) []token {
	var out []token
	r := []rune(in)
	i := 0
	for i < len(r) {
		c := r[i]
		switch {
		case unicode.IsSpace(c):
			i++
		case c == '\'' || c == '"':
			quote := c
			j := i + 1
			for j < len(r) && r[j] != quote {
				if r[j] == '\\' && j+1 < len(r) {
					j += 2
					continue
				}
				j++
			}
			if j < len(r) {
				j++
			}
			out = append(out, token{kind: tokString, value: string(r[i:j])})
			i = j
		case c == '(' || c == ')':
			out = append(out, token{kind: tokParen, value: string(c)})
			i++
		case c == ',' || c == ';' || c == '*' || c == '=' || c == '<' || c == '>' || c == '+' || c == '-' || c == '/':
			// Comment: -- or /* */
			if c == '-' && i+1 < len(r) && r[i+1] == '-' {
				j := i
				for j < len(r) && r[j] != '\n' {
					j++
				}
				out = append(out, token{kind: tokComment, value: string(r[i:j])})
				i = j
				continue
			}
			if c == '/' && i+1 < len(r) && r[i+1] == '*' {
				j := i + 2
				for j+1 < len(r) && !(r[j] == '*' && r[j+1] == '/') {
					j++
				}
				if j+1 < len(r) {
					j += 2
				}
				out = append(out, token{kind: tokComment, value: string(r[i:j])})
				i = j
				continue
			}
			kind := tokSymbol
			if c == ';' {
				kind = tokSemi
			}
			out = append(out, token{kind: kind, value: string(c)})
			i++
		case unicode.IsDigit(c):
			j := i
			for j < len(r) && (unicode.IsDigit(r[j]) || r[j] == '.') {
				j++
			}
			out = append(out, token{kind: tokNumber, value: string(r[i:j])})
			i = j
		default:
			j := i
			for j < len(r) && (unicode.IsLetter(r[j]) || unicode.IsDigit(r[j]) || r[j] == '_' || r[j] == '.') {
				j++
			}
			out = append(out, token{kind: tokWord, value: string(r[i:j])})
			i = j
		}
	}
	return out
}

// ---- internal renderer ----

func render(tokens []token, indent int, upper bool) string {
	var b strings.Builder
	pad := strings.Repeat(" ", indent)
	depth := 0
	atLineStart := true
	writeNewline := func() {
		b.WriteString("\n")
		for k := 0; k < depth; k++ {
			b.WriteString(pad)
		}
		atLineStart = true
	}
	writeSpace := func() {
		if !atLineStart {
			b.WriteString(" ")
		}
	}
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.kind {
		case tokWord:
			// Match multi-word major clauses (e.g., GROUP BY).
			matched := matchMajor(tokens, i)
			if matched > 0 {
				if b.Len() > 0 {
					writeNewline()
				}
				phrase := joinTokens(tokens[i : i+matched])
				if upper {
					phrase = strings.ToUpper(phrase)
				} else {
					phrase = strings.ToLower(phrase)
				}
				b.WriteString(phrase)
				i += matched - 1
				atLineStart = false
				continue
			}
			val := t.value
			if upper && isKeyword(val) {
				val = strings.ToUpper(val)
			} else if !upper && isKeyword(val) {
				val = strings.ToLower(val)
			}
			writeSpace()
			b.WriteString(val)
			atLineStart = false
		case tokNumber, tokString:
			writeSpace()
			b.WriteString(t.value)
			atLineStart = false
		case tokParen:
			if t.value == "(" {
				if !atLineStart {
					b.WriteString(" ")
				}
				b.WriteString("(")
				depth++
			} else {
				depth--
				if depth < 0 {
					depth = 0
				}
				b.WriteString(")")
			}
			atLineStart = false
		case tokSymbol:
			switch t.value {
			case ",":
				b.WriteString(",")
				writeNewline()
			default:
				writeSpace()
				b.WriteString(t.value)
			}
			atLineStart = false
		case tokSemi:
			b.WriteString(";")
			writeNewline()
		case tokComment:
			writeSpace()
			b.WriteString(t.value)
			writeNewline()
		}
	}
	return strings.TrimRight(b.String(), " \n")
}

func matchMajor(tokens []token, i int) int {
	// Try longest 2-word phrase first (e.g., GROUP BY, ORDER BY, INSERT INTO,
	// DELETE FROM, LEFT JOIN, UNION ALL, CREATE TABLE, ALTER TABLE).
	if i+1 < len(tokens) && tokens[i+1].kind == tokWord {
		two := strings.ToUpper(tokens[i].value + " " + tokens[i+1].value)
		for _, c := range majorClauses {
			if c == two {
				return 2
			}
		}
	}
	one := strings.ToUpper(tokens[i].value)
	for _, c := range majorClauses {
		if c == one {
			return 1
		}
	}
	return 0
}

func joinTokens(tt []token) string {
	parts := make([]string, len(tt))
	for i, t := range tt {
		parts[i] = t.value
	}
	return strings.Join(parts, " ")
}

func isKeyword(w string) bool {
	_, ok := keywordSet[strings.ToUpper(w)]
	return ok
}
