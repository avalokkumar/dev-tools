// Package engine defines types shared by every utility engine under pkg/.
//
// Engines are pure: no I/O, no UI, no MCP knowledge. They report user-input
// problems via Diagnostic, never via panic. The exported error return is
// reserved for catastrophic failure (engine bug, OOM).
package engine

// Severity classifies a Diagnostic.
type Severity int

const (
	// SevInfo is purely advisory.
	SevInfo Severity = iota
	// SevWarn flags a likely-problematic input the engine handled best-effort.
	SevWarn
	// SevError indicates the engine could not produce a meaningful result.
	SevError
)

// Diagnostic carries a structured user-facing message with a stable Code.
// Callers should treat any Diagnostic with SevError as a failed Operation.
type Diagnostic struct {
	// Code is a stable dotted-path identifier, e.g. "JSON.PARSE.UNEXPECTED_TOKEN".
	// It is part of the public contract; do not rename without an ADR.
	Code string `json:"code"`

	Message string `json:"message"`

	Severity Severity `json:"severity"`

	// Span optionally locates the diagnostic in the input.
	Span *Span `json:"span,omitempty"`
}

// Span identifies a 1-based line/column range in a textual input.
type Span struct {
	StartLine int `json:"startLine"`
	StartCol  int `json:"startCol"`
	EndLine   int `json:"endLine"`
	EndCol    int `json:"endCol"`
}

// HasError reports whether any diagnostic has SevError severity.
func HasError(d []Diagnostic) bool {
	for i := range d {
		if d[i].Severity == SevError {
			return true
		}
	}
	return false
}
