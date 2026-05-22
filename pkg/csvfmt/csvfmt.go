// Package csvfmt formats and validates CSV.
//
// External API:
//
//	Format([]byte, FormatOptions) (FormatResult, error)
//	Validate([]byte, ValidateOptions) (ValidateResult, error)
package csvfmt

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// FormatOptions tunes CSV reformat.
type FormatOptions struct {
	// Delimiter overrides the field separator. Default ','.
	Delimiter rune `json:"delimiter,omitempty"`
	// Header indicates the first record is a header row (informational).
	Header bool `json:"header,omitempty"`
	// AlignColumns pads each field with spaces so columns line up.
	AlignColumns bool `json:"alignColumns,omitempty"`
}

// FormatResult is the success return.
type FormatResult struct {
	Output      []byte              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Format normalises a CSV. With AlignColumns it pads to common column widths.
func Format(input []byte, opts FormatOptions) (FormatResult, error) {
	rows, diags, fatal := readAll(input, opts.Delimiter)
	if fatal != nil {
		return FormatResult{Diagnostics: diags}, nil
	}
	if opts.AlignColumns {
		return FormatResult{Output: alignColumns(rows, opts.Delimiter), Diagnostics: diags}, nil
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if opts.Delimiter != 0 {
		w.Comma = opts.Delimiter
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return FormatResult{}, fmt.Errorf("csvfmt: write: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return FormatResult{}, fmt.Errorf("csvfmt: flush: %w", err)
	}
	return FormatResult{Output: buf.Bytes(), Diagnostics: diags}, nil
}

func readAll(input []byte, delim rune) ([][]string, []engine.Diagnostic, error) {
	r := csv.NewReader(bytes.NewReader(input))
	if delim != 0 {
		r.Comma = delim
	}
	r.FieldsPerRecord = -1 // tolerate jagged rows; surface as diagnostic
	r.LazyQuotes = false
	var rows [][]string
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, []engine.Diagnostic{{
				Code:     "CSV.PARSE",
				Message:  err.Error(),
				Severity: engine.SevError,
			}}, err
		}
		rows = append(rows, row)
	}
	var diags []engine.Diagnostic
	if len(rows) > 1 {
		want := len(rows[0])
		for i, row := range rows[1:] {
			if len(row) != want {
				diags = append(diags, engine.Diagnostic{
					Code: "CSV.JAGGED_ROW",
					Message: fmt.Sprintf("row %d has %d fields, header has %d",
						i+2, len(row), want),
					Severity: engine.SevWarn,
				})
			}
		}
	}
	return rows, diags, nil
}

func alignColumns(rows [][]string, delim rune) []byte {
	if delim == 0 {
		delim = ','
	}
	if len(rows) == 0 {
		return nil
	}
	cols := 0
	for _, r := range rows {
		if len(r) > cols {
			cols = len(r)
		}
	}
	widths := make([]int, cols)
	for _, r := range rows {
		for i, f := range r {
			if len(f) > widths[i] {
				widths[i] = len(f)
			}
		}
	}
	var buf bytes.Buffer
	for _, r := range rows {
		for i, f := range r {
			if i > 0 {
				buf.WriteRune(delim)
				buf.WriteByte(' ')
			}
			pad := widths[i] - len(f)
			buf.WriteString(f)
			if i < len(r)-1 && pad > 0 {
				buf.WriteString(strings.Repeat(" ", pad))
			}
		}
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

// ValidateOptions tunes Validate.
type ValidateOptions struct {
	// ExpectedColumns, if non-empty, must equal the header row exactly.
	ExpectedColumns []string `json:"expectedColumns,omitempty"`
	// Delimiter overrides ','.
	Delimiter rune `json:"delimiter,omitempty"`
	// Strict requires every row to have exactly len(header) fields.
	Strict bool `json:"strict,omitempty"`
}

// ValidateResult is the success return.
type ValidateResult struct {
	Valid       bool                `json:"valid"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Validate checks parseability, optional header match, optional strict shape.
func Validate(input []byte, opts ValidateOptions) (ValidateResult, error) {
	rows, diags, fatal := readAll(input, opts.Delimiter)
	if fatal != nil {
		return ValidateResult{Valid: false, Diagnostics: diags}, nil
	}
	if len(opts.ExpectedColumns) > 0 {
		if len(rows) == 0 {
			diags = append(diags, engine.Diagnostic{
				Code: "CSV.HEADER_MISSING", Message: "input is empty", Severity: engine.SevError,
			})
		} else if !equalSlice(rows[0], opts.ExpectedColumns) {
			diags = append(diags, engine.Diagnostic{
				Code: "CSV.HEADER_MISMATCH",
				Message: fmt.Sprintf("header %v != expected %v",
					rows[0], opts.ExpectedColumns),
				Severity: engine.SevError,
			})
		}
	}
	if opts.Strict {
		// Promote any JAGGED_ROW warnings to errors.
		for i := range diags {
			if diags[i].Code == "CSV.JAGGED_ROW" {
				diags[i].Severity = engine.SevError
			}
		}
	}
	return ValidateResult{
		Valid:       !engine.HasError(diags),
		Diagnostics: diags,
	}, nil
}

func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
