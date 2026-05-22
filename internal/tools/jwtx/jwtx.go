// Package jwtx is the Adapter that exposes pkg/jwtx as Operations.
package jwtx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/jwtx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{
			Tool: "jwt", Op: "decode",
			Description: "Parse a JWT (no signature check).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["token"],
  "properties":{"token":{"type":"string"}}
}`),
			Handler: handleDecode,
		},
		{
			Tool: "jwt", Op: "verify",
			Description: "Verify a JWT signature, allow-list alg, and check expiry.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["token"],
  "properties":{"token":{"type":"string"},
                "key":{"type":"string"},
                "keyFormat":{"type":"string","enum":["hmac","pem"],"default":"hmac"},
                "expectedAlgs":{"type":"array","items":{"type":"string"}},
                "leewaySeconds":{"type":"integer","default":0}}
}`),
			Handler: handleVerify,
		},
	}
}

type decodeArgs struct {
	Token string `json:"token"`
}

func handleDecode(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a decodeArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Decode(a.Token)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type verifyArgs struct {
	Token         string   `json:"token"`
	Key           string   `json:"key,omitempty"`
	KeyFormat     string   `json:"keyFormat,omitempty"`
	ExpectedAlgs  []string `json:"expectedAlgs,omitempty"`
	LeewaySeconds int      `json:"leewaySeconds,omitempty"`
}

func handleVerify(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a verifyArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	res, err := enginepkg.Verify(a.Token, enginepkg.VerifyOptions{
		Key:          []byte(a.Key),
		KeyFormat:    a.KeyFormat,
		ExpectedAlgs: a.ExpectedAlgs,
		Leeway:       time.Duration(a.LeewaySeconds) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
