// Package mdx is the Adapter that exposes pkg/mdx as Operations.
package mdx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/mdx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "md", Op: "to_html",
			Description: "Render Markdown to sanitised HTML (UGC policy enforced).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"addToc":{"type":"boolean","default":false}}
}`),
			Handler: handleToHTML,
		},
		{
			Tool: "md", Op: "table_from_csv",
			Description: "Build a Markdown table from CSV input.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "delimiter":{"type":"string","maxLength":1},
                "hasHeader":{"type":"boolean","default":true},
                "align":{"type":"array","items":{"type":"string","enum":["left","right","center","none"]}}}
}`),
			Handler: handleTable,
		},
	}
}

type htmlArgs struct {
	Input  string `json:"input"`
	AddTOC bool   `json:"addToc,omitempty"`
}

func handleToHTML(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a htmlArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.ToHTML([]byte(a.Input), enginepkg.ToHTMLOptions{AddTOC: a.AddTOC})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type tableArgs struct {
	Input     string   `json:"input"`
	Delimiter string   `json:"delimiter,omitempty"`
	HasHeader *bool    `json:"hasHeader,omitempty"`
	Align     []string `json:"align,omitempty"`
}

func handleTable(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a tableArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.TableFromCSV([]byte(a.Input), enginepkg.TableOptions{
		Delimiter: a.Delimiter, HasHeader: a.HasHeader, Align: a.Align,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
