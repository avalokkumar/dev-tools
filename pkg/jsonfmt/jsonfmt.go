// Package jsonfmt formats and validates JSON.
//
// External API:
//
//	Format([]byte, FormatOptions) (FormatResult, error)
//	Validate([]byte, []byte /*JSON Schema, optional*/) (ValidateResult, error)
//
// All user-input issues surface as Diagnostic with stable codes.
package jsonfmt

import (
	"bytes"
	"encoding/json"
	"fmt"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/devforge/devforge/pkg/engine"
)

// FormatOptions tunes pretty-printing.
type FormatOptions struct {
	// Indent is the number of spaces per indent level. 0 means compact.
	Indent int `json:"indent,omitempty"`
	// SortKeys reorders object keys alphabetically before printing.
	SortKeys bool `json:"sortKeys,omitempty"`
	// TrailingNewline appends "\n" to the output.
	TrailingNewline bool `json:"trailingNewline,omitempty"`
}

// FormatResult is the success return.
type FormatResult struct {
	Output      []byte              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Format pretty-prints input. Compact when Indent==0.
func Format(input []byte, opts FormatOptions) (FormatResult, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return FormatResult{
			Diagnostics: []engine.Diagnostic{{
				Code:     "JSON.EMPTY",
				Message:  "input is empty",
				Severity: engine.SevError,
			}},
		}, nil
	}

	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		return FormatResult{
			Diagnostics: []engine.Diagnostic{{
				Code:     "JSON.PARSE",
				Message:  err.Error(),
				Severity: engine.SevError,
			}},
		}, nil
	}
	if opts.SortKeys {
		v = sortKeys(v)
	}

	var out []byte
	var err error
	if opts.Indent <= 0 {
		out, err = json.Marshal(v)
	} else {
		indent := bytes.Repeat([]byte(" "), opts.Indent)
		out, err = json.MarshalIndent(v, "", string(indent))
	}
	if err != nil {
		return FormatResult{}, fmt.Errorf("jsonfmt: marshal: %w", err)
	}
	if opts.TrailingNewline {
		out = append(out, '\n')
	}
	return FormatResult{Output: out}, nil
}

func sortKeys(v any) any {
	switch t := v.(type) {
	case map[string]any:
		// json.Marshal already sorts map[string]any keys alphabetically; the
		// recursion below is for nested values only.
		for k, vv := range t {
			t[k] = sortKeys(vv)
		}
		return t
	case []any:
		for i, vv := range t {
			t[i] = sortKeys(vv)
		}
		return t
	default:
		return v
	}
}

// ValidateResult is the success return for Validate.
type ValidateResult struct {
	Valid       bool                `json:"valid"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Validate parses input as JSON. When schema is non-empty it is compiled as
// a JSON-Schema (draft 2020-12 default) and the input is validated against
// it via santhosh-tekuri/jsonschema/v6. All issues surface as diagnostics
// with stable codes; the error return is reserved for catastrophic failure.
func Validate(input, schema []byte) (ValidateResult, error) {
	res := ValidateResult{Valid: true}
	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		res.Valid = false
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code:     "JSON.PARSE",
			Message:  err.Error(),
			Severity: engine.SevError,
		})
		return res, nil
	}
	if len(bytes.TrimSpace(schema)) == 0 {
		return res, nil
	}
	var schemaDoc any
	if err := json.Unmarshal(schema, &schemaDoc); err != nil {
		res.Valid = false
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code:     "JSON.SCHEMA.PARSE",
			Message:  "schema is not valid JSON: " + err.Error(),
			Severity: engine.SevError,
		})
		return res, nil
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", schemaDoc); err != nil {
		res.Valid = false
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "JSON.SCHEMA.LOAD", Message: err.Error(), Severity: engine.SevError,
		})
		return res, nil
	}
	sch, err := c.Compile("schema.json")
	if err != nil {
		res.Valid = false
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "JSON.SCHEMA.COMPILE", Message: err.Error(), Severity: engine.SevError,
		})
		return res, nil
	}
	if err := sch.Validate(v); err != nil {
		res.Valid = false
		// Best-effort detail: collect each leaf cause.
		if ve, ok := err.(*jsonschema.ValidationError); ok {
			for _, cause := range ve.Causes {
				res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
					Code: "JSON.SCHEMA.VIOLATION",
					Message: cause.Error(),
					Severity: engine.SevError,
				})
			}
			if len(ve.Causes) == 0 {
				res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
					Code: "JSON.SCHEMA.VIOLATION", Message: ve.Error(), Severity: engine.SevError,
				})
			}
		} else {
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code: "JSON.SCHEMA.VIOLATION", Message: err.Error(), Severity: engine.SevError,
			})
		}
	}
	return res, nil
}
