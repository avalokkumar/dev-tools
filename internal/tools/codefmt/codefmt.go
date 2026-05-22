// Package codefmt is the Adapter that exposes pkg/codefmt as Operations.
package codefmt

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/codefmt"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "code", Op: "fmt_go",
			Description: "Format Go source via the standard go/format package.",
			InputSchema: schemaInputOnly,
			Handler:     handleGo,
		},
		{Tool: "code", Op: "fmt_xml",
			Description: "Reformat XML with consistent indentation.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"indent":{"type":"integer","default":2}}
}`),
			Handler: handleXML,
		},
		{Tool: "code", Op: "fmt_html",
			Description: "Re-render HTML through golang.org/x/net/html (canonicalises markup).",
			InputSchema: schemaInputOnly,
			Handler:     handleHTML,
		},
	}
}

var schemaInputOnly = json.RawMessage(`{"type":"object","required":["input"],"properties":{"input":{"type":"string"}}}`)

type strInput struct {
	Input string `json:"input"`
}

type xmlArgs struct {
	Input  string `json:"input"`
	Indent int    `json:"indent,omitempty"`
}

func handleGo(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strInput
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.FormatGo([]byte(a.Input))
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleXML(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a xmlArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.FormatXML([]byte(a.Input), enginepkg.XMLOptions{Indent: a.Indent})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleHTML(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strInput
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.FormatHTML([]byte(a.Input))
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
