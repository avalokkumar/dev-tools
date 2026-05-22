// Package httpx is the Adapter that exposes pkg/httpx as Operations.
package httpx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/httpx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "http", Op: "request",
			Description: "curl-like HTTP runner; refuses to dial private IPs unless allowPrivate=true.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["url"],
  "properties":{"method":{"type":"string","default":"GET"},
                "url":{"type":"string"},
                "headers":{"type":"object","additionalProperties":{"type":"string"}},
                "body":{"type":"string"},
                "allowPrivate":{"type":"boolean","default":false},
                "followRedirects":{"type":"boolean","default":false},
                "maxRedirects":{"type":"integer","minimum":0,"maximum":20,"default":5},
                "timeoutSec":{"type":"integer","minimum":1,"maximum":120,"default":30},
                "maxResponseBytes":{"type":"integer","minimum":0,"maximum":67108864,"default":8388608,"description":"response body cap in bytes; hard ceiling 64 MiB"}}
}`),
			Handler: handleRequest,
		},
	}
}

type reqArgs struct {
	Method           string            `json:"method,omitempty"`
	URL              string            `json:"url"`
	Headers          map[string]string `json:"headers,omitempty"`
	Body             string            `json:"body,omitempty"`
	AllowPrivate     bool              `json:"allowPrivate,omitempty"`
	FollowRedirects  bool              `json:"followRedirects,omitempty"`
	MaxRedirects     int               `json:"maxRedirects,omitempty"`
	TimeoutSec       int               `json:"timeoutSec,omitempty"`
	MaxResponseBytes int64             `json:"maxResponseBytes,omitempty"`
}

func handleRequest(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a reqArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Request(
		enginepkg.Req{Method: a.Method, URL: a.URL, Headers: a.Headers, Body: a.Body},
		enginepkg.Options{
			AllowPrivate: a.AllowPrivate, FollowRedirects: a.FollowRedirects,
			MaxRedirects: a.MaxRedirects, TimeoutSec: a.TimeoutSec,
			MaxResponseBytes: a.MaxResponseBytes,
		},
	)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
