// Package tzconv is the Adapter that exposes pkg/tzconv as Operations.
package tzconv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/tzconv"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "tz", Op: "convert",
			Description: "Convert a wall time between IANA zones with DST awareness.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["time","fromTZ","toTZ"],
  "properties":{"time":{"type":"string","description":"RFC3339 or YYYY-MM-DDTHH:MM:SS"},
                "fromTZ":{"type":"string"},"toTZ":{"type":"string"}}
}`),
			Handler: handleConvert,
		},
		{
			Tool: "tz", Op: "list",
			Description: "List well-known IANA zones; optional substring filter.",
			InputSchema: json.RawMessage(`{
  "type":"object",
  "properties":{"filter":{"type":"string"}}
}`),
			Handler: handleList,
		},
	}
}

type convertArgs struct {
	Time   string `json:"time"`
	FromTZ string `json:"fromTZ"`
	ToTZ   string `json:"toTZ"`
}

func handleConvert(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a convertArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	t, err := parseFlexibleTime(a.Time)
	if err != nil {
		b, _ := json.Marshal(map[string]any{
			"diagnostics": []map[string]any{{
				"code": "TZ.CONVERT.BAD_TIME", "message": err.Error(), "severity": 2,
			}},
		})
		return b, nil
	}
	res, err := enginepkg.Convert(t, a.FromTZ, a.ToTZ)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type listArgs struct {
	Filter string `json:"filter,omitempty"`
}

func handleList(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a listArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	zones, err := enginepkg.ListZones(a.Filter)
	if err != nil {
		return nil, err
	}
	return json.Marshal(zones)
}

func parseFlexibleTime(s string) (time.Time, error) {
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &timeParseError{value: s}
}

type timeParseError struct{ value string }

func (e *timeParseError) Error() string {
	return "tzconv: cannot parse time " + e.value + " (try RFC3339)"
}
