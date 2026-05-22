package cryptox

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// HMACOptions tunes HMAC.
type HMACOptions struct {
	// Algorithm: "sha1", "sha256" (default), "sha384", "sha512".
	Algorithm string `json:"algorithm,omitempty"`
	// Key is the HMAC secret.
	Key string `json:"key"`
	// KeyEncoding: "raw" (default), "hex", "base64".
	KeyEncoding string `json:"keyEncoding,omitempty"`
	// OutputEncoding: "hex" (default) or "base64".
	OutputEncoding string `json:"outputEncoding,omitempty"`
}

// HMACResult holds the digest.
type HMACResult struct {
	Output      string              `json:"output"`
	Algorithm   string              `json:"algorithm"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// HMAC computes an HMAC over input.
func HMAC(input []byte, opts HMACOptions) (HMACResult, error) {
	algo := strings.ToLower(opts.Algorithm)
	if algo == "" {
		algo = "sha256"
	}
	var newH func() hash.Hash
	switch algo {
	case "sha1":
		newH = sha1.New
	case "sha256":
		newH = sha256.New
	case "sha384":
		newH = sha512.New384
	case "sha512":
		newH = sha512.New
	default:
		return HMACResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.HMAC.UNSUPPORTED_ALGO",
			Message: fmt.Sprintf("algorithm %q not supported", opts.Algorithm),
			Severity: engine.SevError,
		}}}, nil
	}
	key, err := decodeKey(opts.Key, opts.KeyEncoding)
	if err != nil {
		return HMACResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.HMAC.KEY", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	mac := hmac.New(newH, key)
	_, _ = mac.Write(input)
	sum := mac.Sum(nil)
	out := encodePlain(sum, opts.OutputEncoding)
	if opts.OutputEncoding == "" {
		out = hex.EncodeToString(sum)
	}
	return HMACResult{Output: out, Algorithm: algo}, nil
}

func decodeKey(s, enc string) ([]byte, error) {
	switch strings.ToLower(enc) {
	case "", "raw":
		return []byte(s), nil
	case "hex":
		return hex.DecodeString(s)
	case "base64":
		return base64.StdEncoding.DecodeString(s)
	}
	return nil, fmt.Errorf("unknown keyEncoding %q", enc)
}
