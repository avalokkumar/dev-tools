// Package strx is the Adapter that exposes pkg/strx as Registry Operations.
package strx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/strx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "str", Op: "case",
			Description: "Convert between identifier cases (camel/snake/kebab/...).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input","mode"],
  "properties":{"input":{"type":"string"},
                "mode":{"type":"string","enum":["camel","pascal","snake","kebab","constant","dot","train","header","title","lower","upper"]}}
}`),
			Handler: handleCase,
		},
		{
			Tool: "str", Op: "diff",
			Description: "Line-based unified diff between two text inputs.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["left","right"],
  "properties":{"left":{"type":"string"},"right":{"type":"string"},
                "ignoreWhitespace":{"type":"boolean"},"ignoreCase":{"type":"boolean"}}
}`),
			Handler: handleDiff,
		},
		{
			Tool: "str", Op: "stats",
			Description: "Count lines, words, characters, bytes, longest line.",
			InputSchema: json.RawMessage(`{"type":"object","required":["input"],"properties":{"input":{"type":"string"}}}`),
			Handler: handleStats,
		},
		{
			Tool: "str", Op: "sort_unique",
			Description: "Sort lines; optionally drop duplicates.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},"ignoreCase":{"type":"boolean"},
                "reverse":{"type":"boolean"},"unique":{"type":"boolean"}}
}`),
			Handler: handleSort,
		},
		{
			Tool: "str", Op: "replace",
			Description: "Find-and-replace; literal or regex; optional case-insensitive.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input","pattern","replacement"],
  "properties":{"input":{"type":"string"},"pattern":{"type":"string"},
                "replacement":{"type":"string"},"regex":{"type":"boolean"},
                "ignoreCase":{"type":"boolean"}}
}`),
			Handler: handleReplace,
		},
	}
}

type caseArgs struct {
	Input string `json:"input"`
	Mode  string `json:"mode"`
}

func handleCase(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a caseArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Case(a.Input, a.Mode)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type diffArgs struct {
	Left             string `json:"left"`
	Right            string `json:"right"`
	IgnoreWhitespace bool   `json:"ignoreWhitespace,omitempty"`
	IgnoreCase       bool   `json:"ignoreCase,omitempty"`
}

func handleDiff(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a diffArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Diff(a.Left, a.Right, enginepkg.DiffOptions{IgnoreWhitespace: a.IgnoreWhitespace, IgnoreCase: a.IgnoreCase})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type strInputOnly struct {
	Input string `json:"input"`
}

func handleStats(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strInputOnly
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Stats(a.Input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type sortArgs struct {
	Input      string `json:"input"`
	IgnoreCase bool   `json:"ignoreCase,omitempty"`
	Reverse    bool   `json:"reverse,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
}

func handleSort(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a sortArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.SortUnique(a.Input, enginepkg.SortOptions{IgnoreCase: a.IgnoreCase, Reverse: a.Reverse, Unique: a.Unique})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type replaceArgs struct {
	Input       string `json:"input"`
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
	Regex       bool   `json:"regex,omitempty"`
	IgnoreCase  bool   `json:"ignoreCase,omitempty"`
}

func handleReplace(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a replaceArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Replace(a.Input, a.Pattern, a.Replacement, enginepkg.ReplaceOptions{Regex: a.Regex, IgnoreCase: a.IgnoreCase})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
