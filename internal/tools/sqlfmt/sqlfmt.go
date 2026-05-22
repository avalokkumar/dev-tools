// Package sqlfmt is the Adapter that exposes pkg/sqlfmt as Operations.
package sqlfmt

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/sqlfmt"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "sql", Op: "format",
			Description: "Reformat SQL with consistent indentation and casing.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "indent":{"type":"integer","default":2},
                "uppercase":{"type":"boolean","default":true}}
}`),
			Handler: handleFormat,
		},
		{
			Tool: "sql", Op: "validate",
			Description: "Cheap SQL syntactic and best-practice lints (SELECT *, missing WHERE, etc.).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleValidate,
		},
	}
}

type fmtArgs struct {
	Input     string `json:"input"`
	Indent    int    `json:"indent,omitempty"`
	Uppercase *bool  `json:"uppercase,omitempty"`
}

func handleFormat(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a fmtArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Format(a.Input, enginepkg.FormatOptions{Indent: a.Indent, Uppercase: a.Uppercase})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type valArgs struct {
	Input string `json:"input"`
}

func handleValidate(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a valArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Validate(a.Input, enginepkg.ValidateOptions{})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
