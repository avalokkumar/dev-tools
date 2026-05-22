// Package uuidx generates UUIDs and cryptographic digests.
//
// External API:
//
//	Generate(GenerateOptions) (GenerateResult, error)
//	Hash([]byte, HashOptions) (HashResult, error)
//
// All user-input issues surface as Diagnostic; the error return is reserved
// for catastrophic failure (entropy source unavailable, etc.).
package uuidx

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/google/uuid"

	"github.com/devforge/devforge/pkg/engine"
)

// GenerateOptions configures UUID generation.
type GenerateOptions struct {
	// Version is one of 4 or 7. Default 4 if zero.
	Version int `json:"version,omitempty"`
	// Count is the number of UUIDs to produce. Default 1; capped at 1024.
	Count int `json:"count,omitempty"`
	// Format is one of "std" (default), "compact" (no dashes), or "urn".
	Format string `json:"format,omitempty"`
}

// GenerateResult is the success return.
type GenerateResult struct {
	Values      []string             `json:"values"`
	Diagnostics []engine.Diagnostic  `json:"diagnostics,omitempty"`
}

const maxCount = 1024

// Generate produces UUIDs per the given options.
func Generate(opts GenerateOptions) (GenerateResult, error) {
	if opts.Version == 0 {
		opts.Version = 4
	}
	if opts.Count <= 0 {
		opts.Count = 1
	}
	if opts.Count > maxCount {
		return GenerateResult{
			Diagnostics: []engine.Diagnostic{{
				Code:     "UUID.COUNT_EXCEEDS_LIMIT",
				Message:  fmt.Sprintf("count %d exceeds limit %d", opts.Count, maxCount),
				Severity: engine.SevError,
			}},
		}, nil
	}
	if opts.Version != 4 && opts.Version != 7 {
		return GenerateResult{
			Diagnostics: []engine.Diagnostic{{
				Code:     "UUID.UNSUPPORTED_VERSION",
				Message:  fmt.Sprintf("version %d is not supported (use 4 or 7)", opts.Version),
				Severity: engine.SevError,
			}},
		}, nil
	}

	values := make([]string, 0, opts.Count)
	for i := 0; i < opts.Count; i++ {
		var (
			id  uuid.UUID
			err error
		)
		switch opts.Version {
		case 4:
			id, err = uuid.NewRandom()
		case 7:
			id, err = uuid.NewV7()
		}
		if err != nil {
			return GenerateResult{}, fmt.Errorf("uuidx: generate v%d: %w", opts.Version, err)
		}
		values = append(values, formatUUID(id, opts.Format))
	}
	return GenerateResult{Values: values}, nil
}

func formatUUID(id uuid.UUID, format string) string {
	switch strings.ToLower(format) {
	case "compact":
		return strings.ReplaceAll(id.String(), "-", "")
	case "urn":
		return id.URN()
	default:
		return id.String()
	}
}

// HashOptions configures Hash.
type HashOptions struct {
	// Algos is the list of algorithms; defaults to ["sha256"].
	Algos []string `json:"algos,omitempty"`
	// Encoding is one of "hex" (default) or "base64".
	Encoding string `json:"encoding,omitempty"`
}

// HashResult is the success return.
type HashResult struct {
	Digests     map[string]string    `json:"digests"`
	Diagnostics []engine.Diagnostic  `json:"diagnostics,omitempty"`
}

// Hash computes one or more digests of input.
func Hash(input []byte, opts HashOptions) (HashResult, error) {
	if len(opts.Algos) == 0 {
		opts.Algos = []string{"sha256"}
	}
	if opts.Encoding == "" {
		opts.Encoding = "hex"
	}
	res := HashResult{Digests: map[string]string{}}
	for _, a := range opts.Algos {
		h, ok := newHash(a)
		if !ok {
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code:     "HASH.UNSUPPORTED_ALGO",
				Message:  fmt.Sprintf("algorithm %q is not supported", a),
				Severity: engine.SevError,
			})
			continue
		}
		_, _ = h.Write(input)
		sum := h.Sum(nil)
		res.Digests[strings.ToLower(a)] = encode(sum, opts.Encoding)
	}
	return res, nil
}

func newHash(algo string) (hash.Hash, bool) {
	switch strings.ToLower(algo) {
	case "md5":
		return md5.New(), true
	case "sha1":
		return sha1.New(), true
	case "sha256":
		return sha256.New(), true
	case "sha512":
		return sha512.New(), true
	}
	return nil, false
}

func encode(b []byte, enc string) string {
	switch strings.ToLower(enc) {
	case "base64":
		return base64.StdEncoding.EncodeToString(b)
	default:
		return hex.EncodeToString(b)
	}
}
