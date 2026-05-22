// Package smartdiff is the Adapter that exposes pkg/smartdiff as Operations.
package smartdiff

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/smartdiff"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "diff", Op: "compare",
			Description: "Diff for JSON (semantic, key-aware), INI (per-section/key), or SQL (lightweight statement-level split on ';'; not a full SQL parser — quoted ';' inside string literals or BEGIN/END blocks will be treated as a separator).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["left","right"],
  "properties":{"left":{"type":"string"},"right":{"type":"string"},
                "mode":{"type":"string","enum":["json","ini","sql","auto"],"default":"auto"},
                "ignoreOrder":{"type":"boolean"},"ignoreWhitespace":{"type":"boolean"}}
}`),
			Handler: handleDiff,
		},
	}
}

type diffArgs struct {
	Left             string `json:"left"`
	Right            string `json:"right"`
	Mode             string `json:"mode,omitempty"`
	IgnoreOrder      bool   `json:"ignoreOrder,omitempty"`
	IgnoreWhitespace bool   `json:"ignoreWhitespace,omitempty"`
}

func handleDiff(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a diffArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Diff([]byte(a.Left), []byte(a.Right), enginepkg.DiffOptions{
		Mode: a.Mode, IgnoreOrder: a.IgnoreOrder, IgnoreWhitespace: a.IgnoreWhitespace,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
