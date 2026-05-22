// Package idx generates IDs and URL-safe slugs.
//
// External API:
//
//	ULID(ULIDOptions) (ULIDResult, error)
//	Slugify(string, SlugOptions) (StringResult, error)
package idx

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/oklog/ulid/v2"

	"github.com/devforge/devforge/pkg/engine"
)

// ULIDOptions tunes ULID.
type ULIDOptions struct {
	// Count is the number of ULIDs to produce. Default 1; capped at 1024.
	Count int `json:"count,omitempty"`
	// Lowercase produces the canonical Crockford-base32 ULID in lowercase.
	Lowercase bool `json:"lowercase,omitempty"`
}

// ULIDResult is the success return.
type ULIDResult struct {
	Values      []string            `json:"values"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// ULID generates count ULIDs.
func ULID(opts ULIDOptions) (ULIDResult, error) {
	if opts.Count <= 0 {
		opts.Count = 1
	}
	if opts.Count > 1024 {
		return ULIDResult{Diagnostics: []engine.Diagnostic{{
			Code: "ID.ULID.COUNT_EXCEEDS_LIMIT",
			Message: fmt.Sprintf("count %d exceeds 1024", opts.Count),
			Severity: engine.SevError,
		}}}, nil
	}
	out := make([]string, 0, opts.Count)
	for i := 0; i < opts.Count; i++ {
		id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
		if err != nil {
			return ULIDResult{}, fmt.Errorf("idx: ulid: %w", err)
		}
		s := id.String()
		if opts.Lowercase {
			s = strings.ToLower(s)
		}
		out = append(out, s)
	}
	return ULIDResult{Values: out}, nil
}

// SlugOptions tunes Slugify.
type SlugOptions struct {
	// Sep is the joiner between words. Default "-".
	Sep string `json:"sep,omitempty"`
	// Lower forces output to lowercase. Default true.
	Lower *bool `json:"lower,omitempty"`
	// MaxLen truncates the output to at most N runes (after slugifying). 0 = no limit.
	MaxLen int `json:"maxLen,omitempty"`
}

// StringResult is the canonical string-out result.
type StringResult struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Slugify converts arbitrary text into a URL-safe slug.
func Slugify(input string, opts SlugOptions) (StringResult, error) {
	sep := opts.Sep
	if sep == "" {
		sep = "-"
	}
	lower := true
	if opts.Lower != nil {
		lower = *opts.Lower
	}

	// Replace any non-alphanumeric run with a single sep.
	var b strings.Builder
	prevSep := true
	for _, r := range input {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if lower {
				r = unicode.ToLower(r)
			}
			b.WriteRune(r)
			prevSep = false
		default:
			if !prevSep {
				b.WriteString(sep)
				prevSep = true
			}
		}
	}
	out := strings.Trim(b.String(), sep)
	if opts.MaxLen > 0 {
		// Truncate by rune count.
		runes := []rune(out)
		if len(runes) > opts.MaxLen {
			runes = runes[:opts.MaxLen]
			// Drop any trailing separator after truncation.
			out = strings.TrimRight(string(runes), sep)
		} else {
			out = string(runes)
		}
	}
	return StringResult{Output: out}, nil
}
