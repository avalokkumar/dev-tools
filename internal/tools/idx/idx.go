// Package idx is the Adapter that exposes pkg/idx as Operations.
package idx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/idx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "id", Op: "ulid",
			Description: "Generate one or more ULIDs (Crockford-base32, time-ordered).",
			InputSchema: json.RawMessage(`{
  "type":"object",
  "properties":{"count":{"type":"integer","minimum":1,"maximum":1024,"default":1},
                "lowercase":{"type":"boolean","default":false}}
}`),
			Handler: handleULID,
		},
		{
			Tool: "id", Op: "slug",
			Description: "Slugify arbitrary text into a URL-safe identifier.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "sep":{"type":"string","default":"-"},
                "lower":{"type":"boolean","default":true},
                "maxLen":{"type":"integer","minimum":0,"default":0}}
}`),
			Handler: handleSlug,
		},
	}
}

type ulidArgs struct {
	Count     int  `json:"count,omitempty"`
	Lowercase bool `json:"lowercase,omitempty"`
}

func handleULID(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a ulidArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.ULID(enginepkg.ULIDOptions{Count: a.Count, Lowercase: a.Lowercase})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type slugArgs struct {
	Input  string `json:"input"`
	Sep    string `json:"sep,omitempty"`
	Lower  *bool  `json:"lower,omitempty"`
	MaxLen int    `json:"maxLen,omitempty"`
}

func handleSlug(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a slugArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Slugify(a.Input, enginepkg.SlugOptions{Sep: a.Sep, Lower: a.Lower, MaxLen: a.MaxLen})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
