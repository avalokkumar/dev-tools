// Package faker is the Adapter that exposes pkg/faker as Operations.
package faker

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/faker"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "faker", Op: "generate",
			Description: "Generate synthetic data rows from a field spec.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["spec"],
  "properties":{
    "spec":{"type":"object","required":["fields"],"properties":{"fields":{
      "type":"array","items":{"type":"object","required":["name","kind"],
      "properties":{"name":{"type":"string"},"kind":{"type":"string"},
                    "locale":{"type":"string"},"params":{"type":"object"}}}}}},
    "count":{"type":"integer","minimum":1,"maximum":10000,"default":10},
    "seed":{"type":"integer","default":0},
    "format":{"type":"string","enum":["json","csv","sql"],"default":"json"},
    "table":{"type":"string","default":"data"}}
}`),
			Handler: handleGenerate,
		},
		{
			Tool: "faker", Op: "kinds",
			Description: "List supported field kinds.",
			InputSchema: json.RawMessage(`{"type":"object"}`),
			Handler: func(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"kinds":   enginepkg.Kinds(),
					"locales": enginepkg.Locales(),
				})
			},
		},
	}
}

type genArgs struct {
	Spec   enginepkg.Spec `json:"spec"`
	Count  int            `json:"count,omitempty"`
	Seed   int64          `json:"seed,omitempty"`
	Format string         `json:"format,omitempty"`
	Table  string         `json:"table,omitempty"`
}

func handleGenerate(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a genArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Generate(a.Spec, enginepkg.GenerateOptions{
		Count: a.Count, Seed: a.Seed, Format: a.Format, Table: a.Table,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{
		"output":      string(res.Output),
		"diagnostics": res.Diagnostics,
	})
}
