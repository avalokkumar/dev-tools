// Package regextool is the Adapter that exposes pkg/regextool as Operations.
package regextool

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/regextool"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "regex", Op: "test",
			Description: "Run a regex against input and return all matches.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["pattern","input"],
  "properties":{"pattern":{"type":"string"},"input":{"type":"string"},
                "flavor":{"type":"string","enum":["re2"],"default":"re2"},
                "flags":{"type":"string","description":"i,m,s in any order"}}
}`),
			Handler: handleTest,
		},
		{
			Tool: "regex", Op: "explain",
			Description: "Token-by-token plain-English explanation of a pattern.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["pattern"],
  "properties":{"pattern":{"type":"string"}}
}`),
			Handler: handleExplain,
		},
	}
}

type testArgs struct {
	Pattern string `json:"pattern"`
	Input   string `json:"input"`
	Flavor  string `json:"flavor,omitempty"`
	Flags   string `json:"flags,omitempty"`
}

func handleTest(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a testArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Test(a.Pattern, a.Input, enginepkg.TestOptions{Flavor: a.Flavor, Flags: a.Flags})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type explainArgs struct {
	Pattern string `json:"pattern"`
}

func handleExplain(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a explainArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Explain(a.Pattern, enginepkg.ExplainOptions{})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
