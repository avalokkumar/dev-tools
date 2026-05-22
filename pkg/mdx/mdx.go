// Package mdx provides Markdown utilities (render to sanitized HTML, build
// Markdown tables from CSV).
//
// External API:
//
//	ToHTML([]byte, ToHTMLOptions) (StringResult, error)
//	TableFromCSV([]byte, TableOptions) (StringResult, error)
//
// Security: ToHTML always sanitises with bluemonday's UGC policy. There is
// no opt-out — the engine refuses to emit unsanitised HTML.
package mdx

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"

	"github.com/devforge/devforge/pkg/engine"
)

// StringResult is the canonical string-out result.
type StringResult struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// ToHTMLOptions tunes ToHTML.
type ToHTMLOptions struct {
	// AddTOC inserts a table of contents.
	AddTOC bool `json:"addToc,omitempty"`
}

// ToHTML renders Markdown → sanitised HTML.
func ToHTML(input []byte, opts ToHTMLOptions) (StringResult, error) {
	exts := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(exts)
	flags := html.CommonFlags | html.HrefTargetBlank
	if opts.AddTOC {
		flags |= html.TOC
	}
	r := html.NewRenderer(html.RendererOptions{Flags: flags})
	rendered := markdown.ToHTML(input, p, r)
	clean := bluemonday.UGCPolicy().SanitizeBytes(rendered)
	return StringResult{Output: string(clean)}, nil
}

// TableOptions tunes TableFromCSV.
type TableOptions struct {
	// Delimiter overrides the CSV separator (default ',').
	Delimiter string `json:"delimiter,omitempty"`
	// HasHeader treats the first row as the header. Default true.
	HasHeader *bool `json:"hasHeader,omitempty"`
	// Align controls per-column alignment markers ("left"|"right"|"center"|"none").
	Align []string `json:"align,omitempty"`
}

// TableFromCSV converts a CSV blob into a Markdown table.
func TableFromCSV(input []byte, opts TableOptions) (StringResult, error) {
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
		return StringResult{Diagnostics: []engine.Diagnostic{{
			Code: "MD.TABLE.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	if len(rows) == 0 {
		return StringResult{Diagnostics: []engine.Diagnostic{{
			Code: "MD.TABLE.EMPTY", Message: "input has no rows", Severity: engine.SevError,
		}}}, nil
	}
	cols := len(rows[0])
	header := rows[0]
	body := rows
	if hasHeader {
		body = rows[1:]
	} else {
		header = make([]string, cols)
		for i := range header {
			header[i] = fmt.Sprintf("col%d", i+1)
		}
	}
	var b strings.Builder
	// Header row.
	b.WriteString("| ")
	b.WriteString(strings.Join(header, " | "))
	b.WriteString(" |\n")
	// Separator row with alignment markers.
	for i := 0; i < cols; i++ {
		mark := "---"
		if i < len(opts.Align) {
			switch strings.ToLower(opts.Align[i]) {
			case "left":
				mark = ":---"
			case "right":
				mark = "---:"
			case "center":
				mark = ":---:"
			}
		}
		b.WriteString("| ")
		b.WriteString(mark)
		b.WriteString(" ")
	}
	b.WriteString("|\n")
	// Body rows.
	for _, row := range body {
		// Pad/trim to cols.
		fixed := make([]string, cols)
		for i := 0; i < cols; i++ {
			if i < len(row) {
				fixed[i] = strings.ReplaceAll(row[i], "|", `\|`)
			}
		}
		b.WriteString("| ")
		b.WriteString(strings.Join(fixed, " | "))
		b.WriteString(" |\n")
	}
	return StringResult{Output: b.String()}, nil
}
