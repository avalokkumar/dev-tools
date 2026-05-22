// Package totpx generates and verifies RFC 6238 TOTP codes.
//
// External API:
//
//	Generate(GenerateOptions) (Result, error)
//	Verify(code string, opts VerifyOptions) (Result, error)
//
// Default parameters: SHA-1, 30-second step, 6 digits, T0=0 — matching
// RFC 6238 §4 (canonical TOTP) and Google Authenticator behaviour.
package totpx

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
	"time"

	"github.com/devforge/devforge/pkg/engine"
)

// GenerateOptions tunes Generate.
type GenerateOptions struct {
	// Secret is the shared key.
	Secret string `json:"secret"`
	// SecretEncoding: "base32" (default), "hex", "raw".
	SecretEncoding string `json:"secretEncoding,omitempty"`
	// Algorithm: "sha1" (default), "sha256", "sha512".
	Algorithm string `json:"algorithm,omitempty"`
	// Digits is 6 (default) or 8.
	Digits int `json:"digits,omitempty"`
	// PeriodSec defaults to 30.
	PeriodSec int `json:"periodSec,omitempty"`
	// At lets callers compute a code at a specific Unix timestamp (seconds).
	// 0 means "now".
	At int64 `json:"at,omitempty"`
}

// Result holds the code and the generation context.
type Result struct {
	Code        string              `json:"code,omitempty"`
	Valid       *bool               `json:"valid,omitempty"`
	Counter     uint64              `json:"counter"`
	RemainingMS int64               `json:"remainingMs,omitempty"`
	Algorithm   string              `json:"algorithm"`
	Digits      int                 `json:"digits"`
	PeriodSec   int                 `json:"periodSec"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Generate computes the current TOTP code per the supplied options.
func Generate(opts GenerateOptions) (Result, error) {
	algo := pickAlgo(opts.Algorithm)
	if algo == nil {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "TOTP.UNSUPPORTED_ALGO", Message: fmt.Sprintf("algo %q not supported", opts.Algorithm),
			Severity: engine.SevError,
		}}}, nil
	}
	digits := opts.Digits
	if digits == 0 {
		digits = 6
	}
	if digits != 6 && digits != 8 {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "TOTP.UNSUPPORTED_DIGITS", Message: "digits must be 6 or 8",
			Severity: engine.SevError,
		}}}, nil
	}
	period := opts.PeriodSec
	if period == 0 {
		period = 30
	}
	if period <= 0 {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "TOTP.INVALID_PERIOD", Message: "periodSec must be positive",
			Severity: engine.SevError,
		}}}, nil
	}
	secret, err := decodeSecret(opts.Secret, opts.SecretEncoding)
	if err != nil {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "TOTP.SECRET", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	now := opts.At
	if now == 0 {
		now = time.Now().Unix()
	}
	counter := uint64(now / int64(period))
	code := hotp(secret, counter, digits, algo)
	remaining := int64(period)*1000 - (now*1000)%(int64(period)*1000)
	return Result{
		Code:        code,
		Counter:     counter,
		RemainingMS: remaining,
		Algorithm:   strings.ToLower(opts.Algorithm),
		Digits:      digits,
		PeriodSec:   period,
	}, nil
}

// VerifyOptions tunes Verify.
type VerifyOptions struct {
	GenerateOptions
	// WindowSteps is how many ±steps to accept around the current counter.
	// Default 1 (the current step plus one before/after).
	WindowSteps int `json:"windowSteps,omitempty"`
}

// Verify checks a code against the secret with a sliding window.
func Verify(code string, opts VerifyOptions) (Result, error) {
	algo := pickAlgo(opts.Algorithm)
	if algo == nil {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "TOTP.UNSUPPORTED_ALGO", Message: fmt.Sprintf("algo %q not supported", opts.Algorithm),
			Severity: engine.SevError,
		}}}, nil
	}
	digits := opts.Digits
	if digits == 0 {
		digits = 6
	}
	period := opts.PeriodSec
	if period == 0 {
		period = 30
	}
	window := opts.WindowSteps
	if window == 0 {
		window = 1
	}
	secret, err := decodeSecret(opts.Secret, opts.SecretEncoding)
	if err != nil {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "TOTP.SECRET", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	now := opts.At
	if now == 0 {
		now = time.Now().Unix()
	}
	current := uint64(now / int64(period))
	for i := -window; i <= window; i++ {
		c := uint64(int64(current) + int64(i))
		want := hotp(secret, c, digits, algo)
		if subtleEqual(want, strings.TrimSpace(code)) {
			t := true
			return Result{
				Valid:     &t,
				Counter:   c,
				Algorithm: strings.ToLower(opts.Algorithm),
				Digits:    digits,
				PeriodSec: period,
			}, nil
		}
	}
	f := false
	return Result{Valid: &f, Counter: current, Algorithm: strings.ToLower(opts.Algorithm), Digits: digits, PeriodSec: period}, nil
}

func subtleEqual(a, b string) bool {
	// Constant-time-ish compare; both strings are short fixed-width digits.
	if len(a) != len(b) {
		return false
	}
	var x byte
	for i := 0; i < len(a); i++ {
		x |= a[i] ^ b[i]
	}
	return x == 0
}

func hotp(secret []byte, counter uint64, digits int, newH func() hash.Hash) string {
	var ctr [8]byte
	binary.BigEndian.PutUint64(ctr[:], counter)
	mac := hmac.New(newH, secret)
	_, _ = mac.Write(ctr[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	bin := (uint32(sum[offset])&0x7f)<<24 |
		uint32(sum[offset+1])<<16 |
		uint32(sum[offset+2])<<8 |
		uint32(sum[offset+3])
	mod := uint32(1)
	for i := 0; i < digits; i++ {
		mod *= 10
	}
	code := bin % mod
	return fmt.Sprintf("%0*d", digits, code)
}

func pickAlgo(s string) func() hash.Hash {
	switch strings.ToLower(s) {
	case "", "sha1":
		return sha1.New
	case "sha256":
		return sha256.New
	case "sha512":
		return sha512.New
	}
	return nil
}

func decodeSecret(s, enc string) ([]byte, error) {
	switch strings.ToLower(enc) {
	case "", "base32":
		return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.ReplaceAll(s, " ", "")))
	case "hex":
		return hex.DecodeString(s)
	case "raw":
		return []byte(s), nil
	}
	return nil, fmt.Errorf("unknown secretEncoding %q", enc)
}
