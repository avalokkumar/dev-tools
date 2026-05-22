// Package yamlfmt formats, validates, and converts YAML.
//
// External API:
//
//	Format([]byte, FormatOptions) (FormatResult, error)
//	Validate([]byte, []byte) (ValidateResult, error)
//	Convert([]byte, ConvertOptions) (ConvertResult, error)
package yamlfmt

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/devforge/devforge/pkg/engine"
)

// FormatOptions tunes pretty-printing.
type FormatOptions struct {
	// Indent is the number of spaces per indent level. Default 2.
	Indent int `json:"indent,omitempty"`
}

// FormatResult is the success return.
type FormatResult struct {
	Output      []byte              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Format reformats YAML with consistent indentation.
func Format(input []byte, opts FormatOptions) (FormatResult, error) {
	if opts.Indent <= 0 {
		opts.Indent = 2
	}
	var v any
	if err := yaml.Unmarshal(input, &v); err != nil {
		return FormatResult{
			Diagnostics: []engine.Diagnostic{{
				Code:     "YAML.PARSE",
				Message:  err.Error(),
				Severity: engine.SevError,
			}},
		}, nil
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(opts.Indent)
	if err := enc.Encode(v); err != nil {
		return FormatResult{}, fmt.Errorf("yamlfmt: encode: %w", err)
	}
	_ = enc.Close()
	return FormatResult{Output: buf.Bytes()}, nil
}

// ValidateResult is the success return for Validate.
type ValidateResult struct {
	Valid       bool                `json:"valid"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Validate checks YAML syntax. Schema validation is deferred to Phase D.
func Validate(input, _ []byte) (ValidateResult, error) {
	var v any
	if err := yaml.Unmarshal(input, &v); err != nil {
		return ValidateResult{
			Valid: false,
			Diagnostics: []engine.Diagnostic{{
				Code:     "YAML.PARSE",
				Message:  err.Error(),
				Severity: engine.SevError,
			}},
		}, nil
	}
	return ValidateResult{Valid: true}, nil
}

// ConvertOptions controls Convert.
type ConvertOptions struct {
	// To is one of "json" or "yaml".
	To string `json:"to"`
	// Indent applies to the output format. Default 2.
	Indent int `json:"indent,omitempty"`
}

// ConvertResult is the success return for Convert.
type ConvertResult struct {
	Output      []byte              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Convert transforms YAML ↔ JSON.
func Convert(input []byte, opts ConvertOptions) (ConvertResult, error) {
	if opts.Indent <= 0 {
		opts.Indent = 2
	}
	switch opts.To {
	case "json":
		var v any
		if err := yaml.Unmarshal(input, &v); err != nil {
			return ConvertResult{Diagnostics: []engine.Diagnostic{{
				Code: "YAML.PARSE", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
		v = normalizeForJSON(v)
		out, err := json.MarshalIndent(v, "", spaces(opts.Indent))
		if err != nil {
			return ConvertResult{}, fmt.Errorf("yamlfmt: json marshal: %w", err)
		}
		return ConvertResult{Output: out}, nil
	case "yaml":
		var v any
		if err := json.Unmarshal(input, &v); err != nil {
			return ConvertResult{Diagnostics: []engine.Diagnostic{{
				Code: "YAML.JSON_PARSE", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(opts.Indent)
		if err := enc.Encode(v); err != nil {
			return ConvertResult{}, fmt.Errorf("yamlfmt: yaml encode: %w", err)
		}
		_ = enc.Close()
		return ConvertResult{Output: buf.Bytes()}, nil
	default:
		return ConvertResult{Diagnostics: []engine.Diagnostic{{
			Code: "YAML.CONVERT.UNKNOWN_TARGET", Message: "to must be 'json' or 'yaml'",
			Severity: engine.SevError,
		}}}, nil
	}
}

// yaml.v3 decodes maps as map[string]interface{} when keys are strings, but
// nested maps may use map[interface{}]interface{}. JSON cannot encode the
// latter, so we recursively normalise to map[string]any.
func normalizeForJSON(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			t[k] = normalizeForJSON(vv)
		}
		return t
	case map[any]any:
		out := make(map[string]any, len(t))
		for k, vv := range t {
			out[fmt.Sprint(k)] = normalizeForJSON(vv)
		}
		return out
	case []any:
		for i, vv := range t {
			t[i] = normalizeForJSON(vv)
		}
		return t
	default:
		return v
	}
}

func spaces(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}
