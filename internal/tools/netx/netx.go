// Package netx is the Adapter that exposes pkg/netx as Operations.
// E6 ships url_parse + headers_analyze; E7 will add dns_lookup + http_request.
package netx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/netx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "url", Op: "parse",
			Description: "Decompose a URL into scheme/host/port/path/query/fragment + params.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleURLParse,
		},
		{Tool: "headers", Op: "analyze",
			Description: "Static security-header checklist (no network).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["headers"],
  "properties":{"headers":{"type":"object","additionalProperties":{"type":"string"}}}
}`),
			Handler: handleHeaders,
		},
		{Tool: "dns", Op: "lookup",
			Description: "Resolve DNS records (A/AAAA/MX/TXT/CNAME/NS/PTR). Private IPs blocked unless allowPrivate=true.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["host"],
  "properties":{"host":{"type":"string"},
                "type":{"type":"string","enum":["a","aaaa","mx","txt","cname","ns","ptr","all"],"default":"all"},
                "timeoutSec":{"type":"integer","minimum":1,"maximum":30,"default":3},
                "allowPrivate":{"type":"boolean","default":false}}
}`),
			Handler: handleDNS,
		},
	}
}

type strInput struct {
	Input string `json:"input"`
}

func handleURLParse(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strInput
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.URLParse(a.Input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type headerArgs struct {
	Headers map[string]string `json:"headers"`
}

func handleHeaders(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a headerArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.HeadersAnalyze(a.Headers)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type dnsArgs struct {
	Host         string `json:"host"`
	Type         string `json:"type,omitempty"`
	TimeoutSec   int    `json:"timeoutSec,omitempty"`
	AllowPrivate bool   `json:"allowPrivate,omitempty"`
}

func handleDNS(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a dnsArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.DNSLookup(a.Host, enginepkg.DNSLookupOptions{
		Type: a.Type, TimeoutSec: a.TimeoutSec, AllowPrivate: a.AllowPrivate,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
