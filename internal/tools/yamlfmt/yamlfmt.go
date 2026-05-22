// Package yamlfmt is the Adapter that exposes pkg/yamlfmt as Registry Operations.
package yamlfmt

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/yamlfmt"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool:        "yaml",
			Op:          "format",
			Description: "Reformat YAML with consistent indentation.",
			InputSchema: schemaFormat,
			Handler:     handleFormat,
		},
		{
			Tool:        "yaml",
			Op:          "validate",
			Description: "Validate YAML syntax.",
			InputSchema: schemaValidate,
			Handler:     handleValidate,
		},
		{
			Tool:        "yaml",
			Op:          "convert",
			Description: "Convert between YAML and JSON.",
			InputSchema: schemaConvert,
			Handler:     handleConvert,
		},
	}
}

var (
	schemaFormat = json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"indent":{"type":"integer","default":2}}
}`)
	schemaValidate = json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`)
	schemaConvert = json.RawMessage(`{
  "type":"object","required":["input","to"],
  "properties":{"input":{"type":"string"},"to":{"type":"string","enum":["json","yaml"]},"indent":{"type":"integer","default":2}}
}`)
)

type fmtArgs struct {
	Input  string `json:"input"`
	Indent int    `json:"indent,omitempty"`
}

func handleFormat(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a fmtArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Format([]byte(a.Input), enginepkg.FormatOptions{Indent: a.Indent})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{"output": string(res.Output), "diagnostics": res.Diagnostics})
}

type valArgs struct {
	Input string `json:"input"`
}

func handleValidate(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a valArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Validate([]byte(a.Input), nil)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type convArgs struct {
	Input  string `json:"input"`
	To     string `json:"to"`
	Indent int    `json:"indent,omitempty"`
}

func handleConvert(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a convArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Convert([]byte(a.Input), enginepkg.ConvertOptions{To: a.To, Indent: a.Indent})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{"output": string(res.Output), "diagnostics": res.Diagnostics})
}
