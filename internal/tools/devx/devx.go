// Package devx is the Adapter that exposes pkg/devx as Operations.
package devx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/devx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "dockerfile", Op: "lint",
			Description: "Run best-practice lints over a Dockerfile.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleDockerLint,
		},
		{Tool: "env", Op: "parse",
			Description: "Parse a .env file into a key/value map (reports duplicates and bad lines).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"allowExport":{"type":"boolean","default":false}}
}`),
			Handler: handleEnvParse,
		},
		{Tool: "env", Op: "diff",
			Description: "Diff two .env files (added / removed / changed keys).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["left","right"],
  "properties":{"left":{"type":"string"},"right":{"type":"string"}}
}`),
			Handler: handleEnvDiff,
		},
		{Tool: "k8s", Op: "validate",
			Description: "Structural validation of a Kubernetes YAML manifest (apiVersion/kind/name + parse).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleK8s,
		},
	}
}

type strInput struct {
	Input string `json:"input"`
}

func handleDockerLint(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strInput
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.LintDockerfile([]byte(a.Input))
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type envParseArgs struct {
	Input       string `json:"input"`
	AllowExport bool   `json:"allowExport,omitempty"`
}

func handleEnvParse(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a envParseArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.ParseEnv([]byte(a.Input), enginepkg.ParseEnvOptions{AllowExport: a.AllowExport})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type envDiffArgs struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

func handleEnvDiff(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a envDiffArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.DiffEnv([]byte(a.Left), []byte(a.Right))
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleK8s(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strInput
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.ValidateK8s([]byte(a.Input), enginepkg.K8sOptions{})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
