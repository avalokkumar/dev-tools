// Package netx provides network-adjacent utilities. This file: URL parser
// + HTTP header analyser. dns_lookup and http_request live in sibling files
// and depend on internal/netguard for safety.
package netx

import (
	"net/url"
	"sort"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// URLParseResult is the structured breakdown of a parsed URL.
type URLParseResult struct {
	Scheme      string              `json:"scheme"`
	User        string              `json:"user,omitempty"`
	Host        string              `json:"host"`
	Hostname    string              `json:"hostname"`
	Port        string              `json:"port,omitempty"`
	Path        string              `json:"path"`
	Query       string              `json:"query,omitempty"`
	Fragment    string              `json:"fragment,omitempty"`
	Params      []QueryParam        `json:"params"`
	IsAbsolute  bool                `json:"isAbsolute"`
	IsHTTPS     bool                `json:"isHttps"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// QueryParam is one decoded query parameter.
type QueryParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// URLParse splits a URL into its parts and decodes query parameters.
func URLParse(input string) (URLParseResult, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return URLParseResult{Diagnostics: []engine.Diagnostic{{
			Code: "URL.PARSE.EMPTY", Message: "input is empty", Severity: engine.SevError,
		}}}, nil
	}
	u, err := url.Parse(input)
	if err != nil {
		return URLParseResult{Diagnostics: []engine.Diagnostic{{
			Code: "URL.PARSE.INVALID", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	res := URLParseResult{
		Scheme:     u.Scheme,
		Host:       u.Host,
		Hostname:   u.Hostname(),
		Port:       u.Port(),
		Path:       u.Path,
		Query:      u.RawQuery,
		Fragment:   u.Fragment,
		IsAbsolute: u.IsAbs(),
		IsHTTPS:    u.Scheme == "https",
	}
	if u.User != nil {
		res.User = u.User.Username()
	}
	values := u.Query()
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range values[k] {
			res.Params = append(res.Params, QueryParam{Key: k, Value: v})
		}
	}
	if res.Params == nil {
		res.Params = []QueryParam{}
	}
	return res, nil
}
