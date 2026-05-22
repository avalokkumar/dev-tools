package netx

import (
	"fmt"
	"net/textproto"
	"sort"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// HeadersResult enumerates the security-header findings for an analysed
// header set.
type HeadersResult struct {
	Headers     map[string]string   `json:"headers"`
	Findings    []HeaderFinding     `json:"findings"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// HeaderFinding is one security-header observation.
type HeaderFinding struct {
	Header   string `json:"header"`
	Present  bool   `json:"present"`
	OK       bool   `json:"ok"`
	Note     string `json:"note"`
}

// HeadersAnalyze runs a static checklist over a request/response header map.
// Static-only — no network. Useful for "did I set my CSP / HSTS / X-Frame
// correctly" sanity checks.
func HeadersAnalyze(headers map[string]string) (HeadersResult, error) {
	canon := make(map[string]string, len(headers))
	for k, v := range headers {
		canon[textproto.CanonicalMIMEHeaderKey(k)] = v
	}
	res := HeadersResult{Headers: canon, Findings: []HeaderFinding{}}
	checklist := []struct {
		name string
		need string // recommended substring (case-insensitive); "" means presence-only
	}{
		{"Strict-Transport-Security", "max-age"},
		{"Content-Security-Policy", ""},
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", ""},
		{"Referrer-Policy", ""},
		{"Permissions-Policy", ""},
	}
	for _, c := range checklist {
		v, ok := canon[c.name]
		f := HeaderFinding{Header: c.name, Present: ok}
		switch {
		case !ok:
			f.OK = false
			f.Note = fmt.Sprintf("missing; recommend setting %s", c.name)
		case c.need != "" && !strings.Contains(strings.ToLower(v), c.need):
			f.OK = false
			f.Note = fmt.Sprintf("present but does not contain %q (got %q)", c.need, v)
		default:
			f.OK = true
			f.Note = "ok"
		}
		res.Findings = append(res.Findings, f)
	}
	// Stable order on Headers map output via sorted slice in Findings already.
	sort.Slice(res.Findings, func(i, j int) bool {
		return res.Findings[i].Header < res.Findings[j].Header
	})
	return res, nil
}
