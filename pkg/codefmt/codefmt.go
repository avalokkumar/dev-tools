// Package codefmt provides pure-Go code formatters for Go, XML, HTML.
//
// External API:
//
//	FormatGo([]byte) (Result, error)
//	FormatXML([]byte, XMLOptions) (Result, error)
//	FormatHTML([]byte) (Result, error)
package codefmt

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"strings"

	"golang.org/x/net/html"

	"github.com/devforge/devforge/pkg/engine"
)

// Result holds the formatted output.
type Result struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// FormatGo runs go/format over the input (gofmt semantics).
func FormatGo(input []byte) (Result, error) {
	out, err := format.Source(input)
	if err != nil {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "CODE.GO.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	return Result{Output: string(out)}, nil
}

// XMLOptions tunes FormatXML.
type XMLOptions struct {
	Indent int `json:"indent,omitempty"`
}

// FormatXML re-emits XML with consistent indentation.
func FormatXML(input []byte, opts XMLOptions) (Result, error) {
	indent := opts.Indent
	if indent <= 0 {
		indent = 2
	}
	dec := xml.NewDecoder(bytes.NewReader(input))
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", strings.Repeat(" ", indent))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if cd, ok := tok.(xml.CharData); ok {
			s := strings.TrimSpace(string(cd))
			if s == "" {
				continue
			}
			cd = xml.CharData(s)
			tok = cd
		}
		if err := enc.EncodeToken(tok); err != nil {
			return Result{Diagnostics: []engine.Diagnostic{{
				Code: "CODE.XML.PARSE", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
	}
	if err := enc.Flush(); err != nil {
		return Result{}, fmt.Errorf("codefmt: xml flush: %w", err)
	}
	return Result{Output: buf.String() + "\n"}, nil
}

// FormatHTML re-renders an HTML fragment via golang.org/x/net/html.
func FormatHTML(input []byte) (Result, error) {
	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "CODE.HTML.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return Result{}, fmt.Errorf("codefmt: html render: %w", err)
	}
	return Result{Output: buf.String()}, nil
}
