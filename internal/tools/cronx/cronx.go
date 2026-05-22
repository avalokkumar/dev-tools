// Package cronx is the Adapter that exposes pkg/cronx as Operations.
package cronx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/cronx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "cron", Op: "parse",
			Description: "Validate a cron expression and break it into fields.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["expression"],
  "properties":{"expression":{"type":"string"},
                "flavor":{"type":"string","enum":["unix","quartz","aws"],"default":"unix"}}
}`),
			Handler: handleParse,
		},
		{
			Tool: "cron", Op: "next",
			Description: "Compute the next N run times of a cron expression.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["expression"],
  "properties":{"expression":{"type":"string"},
                "n":{"type":"integer","minimum":1,"maximum":1000,"default":5},
                "from":{"type":"string","description":"RFC3339; default now"},
                "tz":{"type":"string","description":"IANA zone; default UTC"}}
}`),
			Handler: handleNext,
		},
	}
}

type parseArgs struct {
	Expression string `json:"expression"`
	Flavor     string `json:"flavor,omitempty"`
}

func handleParse(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a parseArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Parse(a.Expression, enginepkg.ParseOptions{Flavor: a.Flavor})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type nextArgs struct {
	Expression string `json:"expression"`
	N          int    `json:"n,omitempty"`
	From       string `json:"from,omitempty"`
	TZ         string `json:"tz,omitempty"`
}

func handleNext(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a nextArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	from := time.Now().UTC()
	if a.From != "" {
		t, err := time.Parse(time.RFC3339, a.From)
		if err != nil {
			return diagJSON("CRON.NEXT.BAD_FROM",
				"from must be RFC3339: "+err.Error()), nil
		}
		from = t
	}
	loc := time.UTC
	if a.TZ != "" {
		l, err := time.LoadLocation(a.TZ)
		if err != nil {
			return diagJSON("CRON.NEXT.BAD_TZ",
				"unknown IANA zone "+a.TZ+": "+err.Error()), nil
		}
		loc = l
	}
	res, err := enginepkg.NextRuns(a.Expression, from, a.N, loc)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

// diagJSON returns a Result-shaped error response so user-input parse
// failures stay in the engine diagnostic channel rather than transport errors.
func diagJSON(code, msg string) json.RawMessage {
	b, _ := json.Marshal(map[string]any{
		"runs": []any{},
		"diagnostics": []map[string]any{{
			"code": code, "message": msg, "severity": 2,
		}},
	})
	return b
}
