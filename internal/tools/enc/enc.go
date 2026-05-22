// Package enc is the Adapter that exposes pkg/enc as Registry Operations.
package enc

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/enc"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		// Base64
		{
			Tool: "enc", Op: "base64_encode",
			Description: "Base64-encode an input string (RFC 4648).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "urlSafe":{"type":"boolean","default":false},
                "noPadding":{"type":"boolean","default":false}}
}`),
			Handler: handleB64Encode,
		},
		{
			Tool: "enc", Op: "base64_decode",
			Description: "Base64-decode an input string (auto-detects alphabet/padding).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "urlSafe":{"type":"boolean","default":false},
                "noPadding":{"type":"boolean","default":false}}
}`),
			Handler: handleB64Decode,
		},
		// URL
		{
			Tool: "enc", Op: "url_encode",
			Description: "URL-encode (percent-encode) a string per RFC 3986.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "mode":{"type":"string","enum":["component","path"],"default":"component"}}
}`),
			Handler: handleURLEncode,
		},
		{
			Tool: "enc", Op: "url_decode",
			Description: "Reverse URL percent-encoding.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleURLDecode,
		},
		// HTML
		{
			Tool: "enc", Op: "html_encode",
			Description: "Replace HTML special characters with named/numeric entities.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleHTMLEncode,
		},
		{
			Tool: "enc", Op: "html_decode",
			Description: "Decode HTML entities back to plain text.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleHTMLDecode,
		},
		// Hex
		{
			Tool: "enc", Op: "hex_encode",
			Description: "Encode an input string to hex.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "upper":{"type":"boolean","default":false}}
}`),
			Handler: handleHexEncode,
		},
		{
			Tool: "enc", Op: "hex_decode",
			Description: "Decode a hex string (tolerates 0x prefix and case).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"}}
}`),
			Handler: handleHexDecode,
		},
	}
}

type b64Args struct {
	Input     string `json:"input"`
	URLSafe   bool   `json:"urlSafe,omitempty"`
	NoPadding bool   `json:"noPadding,omitempty"`
}

func handleB64Encode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a b64Args
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Base64Encode([]byte(a.Input), enginepkg.Base64Options{URLSafe: a.URLSafe, NoPadding: a.NoPadding})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleB64Decode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a b64Args
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Base64Decode(a.Input, enginepkg.Base64Options{URLSafe: a.URLSafe, NoPadding: a.NoPadding})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type urlArgs struct {
	Input string `json:"input"`
	Mode  string `json:"mode,omitempty"`
}

func handleURLEncode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a urlArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.URLEncode(a.Input, enginepkg.URLOptions{Mode: a.Mode})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type strArgs struct {
	Input string `json:"input"`
}

func handleURLDecode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.URLDecode(a.Input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleHTMLEncode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.HTMLEncode(a.Input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleHTMLDecode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.HTMLDecode(a.Input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type hexArgs struct {
	Input string `json:"input"`
	Upper bool   `json:"upper,omitempty"`
}

func handleHexEncode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a hexArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.HexEncode([]byte(a.Input), enginepkg.HexOptions{Upper: a.Upper})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func handleHexDecode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a strArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.HexDecode(a.Input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
