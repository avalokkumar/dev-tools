// Package csvfmt is the Adapter that exposes pkg/csvfmt as Registry Operations.
package csvfmt

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/csvfmt"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "csv", Op: "format",
			Description: "Reformat CSV; optionally align columns.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"delimiter":{"type":"string","maxLength":1},
                "header":{"type":"boolean"},"alignColumns":{"type":"boolean"}}
}`),
			Handler: handleFormat,
		},
		{
			Tool: "csv", Op: "validate",
			Description: "Validate CSV; optionally enforce header / strict shape.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"delimiter":{"type":"string","maxLength":1},
                "expectedColumns":{"type":"array","items":{"type":"string"}},
                "strict":{"type":"boolean"}}
}`),
			Handler: handleValidate,
		},
	}
}

type fmtArgs struct {
	Input        string `json:"input"`
	Delimiter    string `json:"delimiter,omitempty"`
	Header       bool   `json:"header,omitempty"`
	AlignColumns bool   `json:"alignColumns,omitempty"`
}

func handleFormat(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a fmtArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Format([]byte(a.Input), enginepkg.FormatOptions{
		Delimiter: firstRune(a.Delimiter), Header: a.Header, AlignColumns: a.AlignColumns,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{"output": string(res.Output), "diagnostics": res.Diagnostics})
}

type valArgs struct {
	Input           string   `json:"input"`
	Delimiter       string   `json:"delimiter,omitempty"`
	ExpectedColumns []string `json:"expectedColumns,omitempty"`
	Strict          bool     `json:"strict,omitempty"`
}

func handleValidate(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a valArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Validate([]byte(a.Input), enginepkg.ValidateOptions{
		Delimiter: firstRune(a.Delimiter), ExpectedColumns: a.ExpectedColumns, Strict: a.Strict,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func firstRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}
