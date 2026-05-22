// Package jsonfmt is the Adapter that exposes pkg/jsonfmt as Registry Operations.
package jsonfmt

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/jsonfmt"
)

// Operations returns the JSON Tool's Operations.
func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool:        "json",
			Op:          "format",
			Description: "Pretty-print or compact JSON.",
			InputSchema: json.RawMessage(`{
  "type": "object",
  "required": ["input"],
  "properties": {
    "input":           {"type": "string"},
    "indent":          {"type": "integer", "minimum": 0, "maximum": 16, "default": 2},
    "sortKeys":        {"type": "boolean", "default": false},
    "trailingNewline": {"type": "boolean", "default": false}
  }
}`),
			Handler: handleFormat,
		},
		{
			Tool:        "json",
			Op:          "validate",
			Description: "Validate JSON syntax. When a JSON Schema is supplied, validate against it (draft 2020-12) via santhosh-tekuri/jsonschema/v6.",
			InputSchema: json.RawMessage(`{
  "type": "object",
  "required": ["input"],
  "properties": {
    "input":  {"type": "string"},
    "schema": {"type": "string", "description": "Optional JSON-Schema document."}
  }
}`),
			Handler: handleValidate,
		},
	}
}

type formatArgs struct {
	Input           string `json:"input"`
	Indent          int    `json:"indent,omitempty"`
	SortKeys        bool   `json:"sortKeys,omitempty"`
	TrailingNewline bool   `json:"trailingNewline,omitempty"`
}

func handleFormat(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a formatArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Format([]byte(a.Input), enginepkg.FormatOptions{
		Indent: a.Indent, SortKeys: a.SortKeys, TrailingNewline: a.TrailingNewline,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{
		"output":      string(res.Output),
		"diagnostics": res.Diagnostics,
	})
}

type validateArgs struct {
	Input  string `json:"input"`
	Schema string `json:"schema,omitempty"`
}

func handleValidate(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a validateArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Validate([]byte(a.Input), []byte(a.Schema))
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
