// Package mathx is the Adapter that exposes pkg/mathx as Operations.
package mathx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/mathx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "math", Op: "eval",
			Description: "Safely evaluate an arithmetic expression (no I/O, no shell).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["expression"],
  "properties":{"expression":{"type":"string"}}
}`),
			Handler: handleEval,
		},
		{Tool: "math", Op: "unit",
			Description: "Convert between units (bytes, time, throughput, temperature, length).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["value","from","to"],
  "properties":{"value":{"type":"number"},"from":{"type":"string"},"to":{"type":"string"}}
}`),
			Handler: handleUnit,
		},
	}
}

type evalArgs struct {
	Expression string `json:"expression"`
}

func handleEval(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a evalArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Eval(a.Expression)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type unitArgs struct {
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

func handleUnit(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a unitArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.UnitConvert(a.Value, a.From, a.To)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
