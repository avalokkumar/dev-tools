// Package gitx is the Adapter that exposes pkg/gitx as Operations.
package gitx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/gitx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "git", Op: "patch",
			Description: "Build a unified-diff patch between two text inputs.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["left","right"],
  "properties":{"left":{"type":"string"},"right":{"type":"string"},
                "leftPath":{"type":"string","default":"a"},"rightPath":{"type":"string","default":"b"},
                "context":{"type":"integer","default":3}}
}`),
			Handler: handlePatch,
		},
		{Tool: "git", Op: "commit_format",
			Description: "Validate a commit message against Conventional Commits v1.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleCommit,
		},
		{Tool: "git", Op: "ignore_gen",
			Description: "Generate a .gitignore body by combining curated templates.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["templates"],
  "properties":{"templates":{"type":"array","items":{"type":"string"}}}
}`),
			Handler: handleIgnore,
		},
	}
}

type patchArgs struct {
	Left      string `json:"left"`
	Right     string `json:"right"`
	LeftPath  string `json:"leftPath,omitempty"`
	RightPath string `json:"rightPath,omitempty"`
	Context   int    `json:"context,omitempty"`
}

func handlePatch(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a patchArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Patch(a.Left, a.Right, enginepkg.PatchOptions{
		LeftPath: a.LeftPath, RightPath: a.RightPath, Context: a.Context,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type commitArgs struct {
	Input string `json:"input"`
}

func handleCommit(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a commitArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.CommitFormat(a.Input, enginepkg.CommitOptions{})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type ignoreArgs struct {
	Templates []string `json:"templates"`
}

func handleIgnore(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a ignoreArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.IgnoreGen(a.Templates)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
