// Package timex is the Adapter that exposes pkg/timex as Registry Operations.
package timex

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/timex"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "time", Op: "convert",
			Description: "Convert between epoch (s/ms/us/ns) and RFC3339 dates.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "inputFormat":{"type":"string","enum":["auto","epoch_s","epoch_ms","epoch_us","epoch_ns","rfc3339","iso8601"],"default":"auto"},
                "tz":{"type":"string","default":"UTC"}}
}`),
			Handler: handleConvert,
		},
		{
			Tool: "time", Op: "relative",
			Description: "Human-readable relative phrase between two RFC3339 instants.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["from","to"],
  "properties":{"from":{"type":"string"},"to":{"type":"string"}}
}`),
			Handler: handleRelative,
		},
		{
			Tool: "time", Op: "duration",
			Description: "Parse a Go duration (or plain seconds) and break it down.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleDuration,
		},
	}
}

type convertArgs struct {
	Input       string `json:"input"`
	InputFormat string `json:"inputFormat,omitempty"`
	TZ          string `json:"tz,omitempty"`
}

func handleConvert(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a convertArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Convert(enginepkg.ConvertOptions{Input: a.Input, InputFormat: a.InputFormat, TZ: a.TZ})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type relArgs struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func handleRelative(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a relArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	from, err := parseTime(a.From)
	if err != nil {
		return diagJSON("TIME.RELATIVE.BAD_FROM", err.Error()), nil
	}
	to, err := parseTime(a.To)
	if err != nil {
		return diagJSON("TIME.RELATIVE.BAD_TO", err.Error()), nil
	}
	res, err := enginepkg.Relative(from, to)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func diagJSON(code, msg string) json.RawMessage {
	b, _ := json.Marshal(map[string]any{
		"phrase":  "",
		"seconds": 0,
		"diagnostics": []map[string]any{{
			"code": code, "message": msg, "severity": 2,
		}},
	})
	return b
}

type durArgs struct {
	Input string `json:"input"`
}

func handleDuration(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a durArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Duration(a.Input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func parseTime(s string) (time.Time, error) {
	for _, l := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(l, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &timeParseError{s}
}

type timeParseError struct{ value string }

func (e *timeParseError) Error() string { return "timex: cannot parse " + e.value }
