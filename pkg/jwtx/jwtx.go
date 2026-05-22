// Package jwtx decodes and verifies JSON Web Tokens.
//
// External API:
//
//	Decode(token string) (DecodeResult, error)
//	Verify(token string, opts VerifyOptions) (VerifyResult, error)
//
// Decode never validates a signature. Verify enforces algorithm allow-list,
// expiry (with optional leeway), and key material per algorithm.
package jwtx

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/devforge/devforge/pkg/engine"
)

// RawParts holds the three base64url segments separately.
type RawParts struct {
	Header    string `json:"header"`
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

// DecodeResult is the success return.
type DecodeResult struct {
	Header      map[string]any      `json:"header"`
	Payload     map[string]any      `json:"payload"`
	Signature   []byte              `json:"signature,omitempty"`
	Raw         RawParts            `json:"raw"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Decode parses (without verifying) a JWS Compact Serialization token.
func Decode(token string) (DecodeResult, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return DecodeResult{Diagnostics: []engine.Diagnostic{{
			Code: "JWT.MALFORMED", Message: "token must have 3 dot-separated parts",
			Severity: engine.SevError,
		}}}, nil
	}
	hdr, perr := decodeSegment(parts[0])
	if perr != nil {
		return DecodeResult{Diagnostics: []engine.Diagnostic{{
			Code: "JWT.HEADER.PARSE", Message: perr.Error(), Severity: engine.SevError,
		}}}, nil
	}
	pl, perr := decodeSegment(parts[1])
	if perr != nil {
		return DecodeResult{Diagnostics: []engine.Diagnostic{{
			Code: "JWT.PAYLOAD.PARSE", Message: perr.Error(), Severity: engine.SevError,
		}}}, nil
	}
	sig, _ := base64.RawURLEncoding.DecodeString(parts[2])
	return DecodeResult{
		Header:    hdr,
		Payload:   pl,
		Signature: sig,
		Raw:       RawParts{Header: parts[0], Payload: parts[1], Signature: parts[2]},
	}, nil
}

func decodeSegment(s string) (map[string]any, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		// allow non-padded standard base64 too
		b, err = base64.RawStdEncoding.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("base64: %w", err)
		}
	}
	var v map[string]any
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}
	return v, nil
}

// VerifyOptions configures Verify.
type VerifyOptions struct {
	// Key is the verification key material. For HMAC algorithms it is the
	// shared secret. For RS256/RS384/RS512 it is the PEM-encoded public key
	// (use KeyFormat="pem").
	Key []byte `json:"key,omitempty"`

	// KeyFormat is "hmac" (default), "pem", or "jwk" (jwk reserved for Phase D).
	KeyFormat string `json:"keyFormat,omitempty"`

	// ExpectedAlgs is the allow-list. Empty disables alg enforcement (not
	// recommended; "none" is always rejected).
	ExpectedAlgs []string `json:"expectedAlgs,omitempty"`

	// Leeway is added to expiry checks.
	Leeway time.Duration `json:"leeway,omitempty"`
}

// VerifyResult is the success return.
type VerifyResult struct {
	Valid       bool                `json:"valid"`
	ExpiresIn   *time.Duration      `json:"expiresIn,omitempty"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Verify checks signature, algorithm, and expiry.
func Verify(token string, opts VerifyOptions) (VerifyResult, error) {
	dec, err := Decode(token)
	if err != nil {
		return VerifyResult{}, err
	}
	if engine.HasError(dec.Diagnostics) {
		return VerifyResult{Valid: false, Diagnostics: dec.Diagnostics}, nil
	}
	algName, _ := dec.Header["alg"].(string)
	if strings.EqualFold(algName, "none") {
		return VerifyResult{Valid: false, Diagnostics: []engine.Diagnostic{{
			Code: "JWT.ALG_NONE", Message: "alg \"none\" is rejected",
			Severity: engine.SevError,
		}}}, nil
	}
	if len(opts.ExpectedAlgs) > 0 && !contains(opts.ExpectedAlgs, algName) {
		return VerifyResult{Valid: false, Diagnostics: []engine.Diagnostic{{
			Code: "JWT.ALG_NOT_ALLOWED",
			Message: fmt.Sprintf("alg %q not in allow-list %v", algName, opts.ExpectedAlgs),
			Severity: engine.SevError,
		}}}, nil
	}

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		switch strings.ToLower(opts.KeyFormat) {
		case "", "hmac":
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("token alg %q is not HMAC", t.Method.Alg())
			}
			return opts.Key, nil
		case "pem":
			block, _ := pem.Decode(opts.Key)
			if block == nil {
				return nil, fmt.Errorf("invalid PEM")
			}
			pub, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				// Try RSA-specific parser as fallback.
				pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
				if err != nil {
					return nil, fmt.Errorf("parse pubkey: %w", err)
				}
			}
			return pub, nil
		default:
			return nil, fmt.Errorf("unsupported keyFormat %q", opts.KeyFormat)
		}
	}, jwt.WithLeeway(opts.Leeway))

	if err != nil || parsed == nil || !parsed.Valid {
		// Distinguish expiry from signature errors via the library's error string.
		code := "JWT.INVALID"
		if err != nil && strings.Contains(strings.ToLower(err.Error()), "expired") {
			code = "JWT.EXPIRED"
		}
		return VerifyResult{Valid: false, Diagnostics: []engine.Diagnostic{{
			Code: code, Message: errString(err), Severity: engine.SevError,
		}}}, nil
	}

	res := VerifyResult{Valid: true}
	// Compute expires-in.
	if exp, ok := dec.Payload["exp"].(float64); ok {
		d := time.Until(time.Unix(int64(exp), 0))
		res.ExpiresIn = &d
	}
	return res, nil
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func errString(err error) string {
	if err == nil {
		return "invalid"
	}
	return err.Error()
}
