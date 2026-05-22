// Package enc provides encoder/decoder Operations for Base64, URL, HTML, Hex.
//
// External API:
//
//	Base64Encode([]byte, Base64Options) (Base64Result, error)
//	Base64Decode(string, Base64Options) (DecodeResult, error)
//	URLEncode(string, URLOptions) (StringResult, error)
//	URLDecode(string) (StringResult, error)
//	HTMLEncode(string) (StringResult, error)
//	HTMLDecode(string) (StringResult, error)
//	HexEncode([]byte, HexOptions) (StringResult, error)
//	HexDecode(string) (DecodeResult, error)
//
// All input issues surface as Diagnostic; error reserved for catastrophic failure.
package enc

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// ---------- Base64 ----------

// Base64Options tunes Base64 encode/decode.
type Base64Options struct {
	// URLSafe selects RFC 4648 §5 URL-safe alphabet ("-_" instead of "+/").
	URLSafe bool `json:"urlSafe,omitempty"`
	// NoPadding strips/accepts the trailing '=' padding.
	NoPadding bool `json:"noPadding,omitempty"`
}

// Base64Result holds encoded output.
type Base64Result struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// DecodeResult is shared by decoders that produce raw bytes.
type DecodeResult struct {
	Output      string              `json:"output"`      // best-effort UTF-8 view
	Bytes       []byte              `json:"bytes"`       // raw decoded bytes
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Base64Encode encodes input per options.
func Base64Encode(input []byte, opts Base64Options) (Base64Result, error) {
	enc := pickBase64(opts)
	return Base64Result{Output: enc.EncodeToString(input)}, nil
}

// Base64Decode decodes per options. Padding is auto-detected if input ends in '=';
// callers can force no-padding via NoPadding=true.
func Base64Decode(input string, opts Base64Options) (DecodeResult, error) {
	// Auto-tolerate either alphabet by trying URL first when standard fails (and vice versa).
	candidates := []*base64.Encoding{pickBase64(opts)}
	// Add fallbacks so users don't have to guess.
	if opts.URLSafe {
		candidates = append(candidates, base64.StdEncoding, base64.RawStdEncoding)
	} else {
		candidates = append(candidates, base64.URLEncoding, base64.RawURLEncoding)
	}
	candidates = append(candidates, base64.RawStdEncoding, base64.RawURLEncoding)

	var lastErr error
	for _, c := range candidates {
		b, err := c.DecodeString(input)
		if err == nil {
			return DecodeResult{Output: string(b), Bytes: b}, nil
		}
		lastErr = err
	}
	return DecodeResult{
		Diagnostics: []engine.Diagnostic{{
			Code:     "ENC.BASE64.INVALID",
			Message:  fmt.Sprintf("not valid base64: %v", lastErr),
			Severity: engine.SevError,
		}},
	}, nil
}

func pickBase64(opts Base64Options) *base64.Encoding {
	switch {
	case opts.URLSafe && opts.NoPadding:
		return base64.RawURLEncoding
	case opts.URLSafe:
		return base64.URLEncoding
	case opts.NoPadding:
		return base64.RawStdEncoding
	default:
		return base64.StdEncoding
	}
}

// ---------- URL ----------

// URLOptions tunes URL encode.
type URLOptions struct {
	// Mode is "component" (default) or "path".
	Mode string `json:"mode,omitempty"`
}

// StringResult is the canonical output for string-out engines.
type StringResult struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// URLEncode encodes a string per RFC 3986 with the chosen escape policy.
func URLEncode(input string, opts URLOptions) (StringResult, error) {
	switch opts.Mode {
	case "", "component":
		return StringResult{Output: url.QueryEscape(input)}, nil
	case "path":
		return StringResult{Output: url.PathEscape(input)}, nil
	default:
		return StringResult{Diagnostics: []engine.Diagnostic{{
			Code: "ENC.URL.UNKNOWN_MODE", Message: "mode must be component or path",
			Severity: engine.SevError,
		}}}, nil
	}
}

// URLDecode reverses URL percent-encoding.
func URLDecode(input string) (StringResult, error) {
	out, err := url.QueryUnescape(input)
	if err != nil {
		return StringResult{Diagnostics: []engine.Diagnostic{{
			Code: "ENC.URL.INVALID", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	return StringResult{Output: out}, nil
}

// ---------- HTML ----------

// HTMLEncode escapes "<", ">", "&", "\"", "'" to named/numeric entities.
func HTMLEncode(input string) (StringResult, error) {
	return StringResult{Output: html.EscapeString(input)}, nil
}

// HTMLDecode reverses HTML entity encoding.
func HTMLDecode(input string) (StringResult, error) {
	return StringResult{Output: html.UnescapeString(input)}, nil
}

// ---------- Hex ----------

// HexOptions tunes hex encode.
type HexOptions struct {
	// Upper produces uppercase hex digits.
	Upper bool `json:"upper,omitempty"`
}

// HexEncode encodes bytes as a hex string.
func HexEncode(input []byte, opts HexOptions) (StringResult, error) {
	out := hex.EncodeToString(input)
	if opts.Upper {
		out = strings.ToUpper(out)
	}
	return StringResult{Output: out}, nil
}

// HexDecode parses a hex string. Tolerates upper, lower, and "0x" prefix.
func HexDecode(input string) (DecodeResult, error) {
	in := strings.TrimSpace(input)
	in = strings.TrimPrefix(in, "0x")
	in = strings.TrimPrefix(in, "0X")
	b, err := hex.DecodeString(in)
	if err != nil {
		return DecodeResult{Diagnostics: []engine.Diagnostic{{
			Code: "ENC.HEX.INVALID", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	return DecodeResult{Output: string(b), Bytes: b}, nil
}
