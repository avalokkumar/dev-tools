// Package uuid is the Adapter that exposes pkg/uuidx as Registry Operations.
//
// Adapters do translation only — no business logic. They receive raw JSON
// arguments, construct the Engine's Options struct, call the Engine, and
// return the JSON-encoded Result.
package uuid

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/pkg/uuidx"
)

// Operations returns the UUID Tool's Operations.
func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool:        "uuid",
			Op:          "generate",
			Description: "Generate one or more UUIDs (v4 or v7).",
			InputSchema: json.RawMessage(`{
  "type": "object",
  "properties": {
    "version": {"type": "integer", "enum": [4, 7], "default": 4, "description": "UUID version"},
    "count":   {"type": "integer", "minimum": 1, "maximum": 1024, "default": 1},
    "format":  {"type": "string", "enum": ["std", "compact", "urn"], "default": "std"}
  }
}`),
			Handler: handleGenerate,
		},
		{
			Tool:        "uuid",
			Op:          "hash",
			Description: "Compute one or more cryptographic digests of an input string.",
			InputSchema: json.RawMessage(`{
  "type": "object",
  "required": ["input"],
  "properties": {
    "input":    {"type": "string", "description": "UTF-8 input to hash"},
    "algos":    {"type": "array", "items": {"type": "string", "enum": ["md5", "sha1", "sha256", "sha512"]}, "default": ["sha256"]},
    "encoding": {"type": "string", "enum": ["hex", "base64"], "default": "hex"}
  }
}`),
			Handler: handleHash,
		},
	}
}

func handleGenerate(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var opts uuidx.GenerateOptions
	if len(args) > 0 {
		if err := json.Unmarshal(args, &opts); err != nil {
			return nil, err
		}
	}
	res, err := uuidx.Generate(opts)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type hashArgs struct {
	Input    string   `json:"input"`
	Algos    []string `json:"algos,omitempty"`
	Encoding string   `json:"encoding,omitempty"`
}

func handleHash(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a hashArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := uuidx.Hash([]byte(a.Input), uuidx.HashOptions{
		Algos:    a.Algos,
		Encoding: a.Encoding,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
