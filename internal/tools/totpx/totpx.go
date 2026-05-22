// Package totpx is the Adapter that exposes pkg/totpx as Operations.
package totpx

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/totpx"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "totp", Op: "generate",
			Description: "Generate an RFC 6238 TOTP code from a base32 secret.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["secret"],
  "properties":{"secret":{"type":"string"},
                "secretEncoding":{"type":"string","enum":["base32","hex","raw"],"default":"base32"},
                "algorithm":{"type":"string","enum":["sha1","sha256","sha512"],"default":"sha1"},
                "digits":{"type":"integer","enum":[6,8],"default":6},
                "periodSec":{"type":"integer","default":30},
                "at":{"type":"integer","description":"Unix seconds; 0 = now"}}
}`),
			Handler: handleGenerate,
		},
		{Tool: "totp", Op: "verify",
			Description: "Verify an RFC 6238 TOTP code with a sliding window.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["code","secret"],
  "properties":{"code":{"type":"string"},"secret":{"type":"string"},
                "secretEncoding":{"type":"string","enum":["base32","hex","raw"],"default":"base32"},
                "algorithm":{"type":"string","enum":["sha1","sha256","sha512"],"default":"sha1"},
                "digits":{"type":"integer","enum":[6,8],"default":6},
                "periodSec":{"type":"integer","default":30},
                "windowSteps":{"type":"integer","default":1},
                "at":{"type":"integer"}}
}`),
			Handler: handleVerify,
		},
	}
}

type genArgs struct {
	Secret         string `json:"secret"`
	SecretEncoding string `json:"secretEncoding,omitempty"`
	Algorithm      string `json:"algorithm,omitempty"`
	Digits         int    `json:"digits,omitempty"`
	PeriodSec      int    `json:"periodSec,omitempty"`
	At             int64  `json:"at,omitempty"`
}

func handleGenerate(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a genArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Generate(enginepkg.GenerateOptions{
		Secret: a.Secret, SecretEncoding: a.SecretEncoding,
		Algorithm: a.Algorithm, Digits: a.Digits, PeriodSec: a.PeriodSec, At: a.At,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type verArgs struct {
	Code           string `json:"code"`
	Secret         string `json:"secret"`
	SecretEncoding string `json:"secretEncoding,omitempty"`
	Algorithm      string `json:"algorithm,omitempty"`
	Digits         int    `json:"digits,omitempty"`
	PeriodSec      int    `json:"periodSec,omitempty"`
	WindowSteps    int    `json:"windowSteps,omitempty"`
	At             int64  `json:"at,omitempty"`
}

func handleVerify(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a verArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.Verify(a.Code, enginepkg.VerifyOptions{
		GenerateOptions: enginepkg.GenerateOptions{
			Secret: a.Secret, SecretEncoding: a.SecretEncoding,
			Algorithm: a.Algorithm, Digits: a.Digits, PeriodSec: a.PeriodSec, At: a.At,
		},
		WindowSteps: a.WindowSteps,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
