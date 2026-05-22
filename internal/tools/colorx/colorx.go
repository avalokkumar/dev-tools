// Package colorx is the Adapter that exposes pkg/colorx as Operations.
package colorx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/colorx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "color", Op: "convert",
			Description: "Convert color between hex / rgb / hsl / CSS-named.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "to":{"type":"string","enum":["hex","rgb","hsl",""],"default":""}}
}`),
			Handler: handleConvert,
		},
	}
}

type convertArgs struct {
	Input string `json:"input"`
	To    string `json:"to,omitempty"`
}

func handleConvert(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a convertArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Convert(a.Input, enginepkg.Options{To: a.To})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
