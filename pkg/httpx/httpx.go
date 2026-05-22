// Package httpx is a curl-like HTTP request runner.
//
// External API:
//
//	Request(Req, Options) (Response, error)
//
// Security model: callers must opt in to private/loopback targets. The
// engine resolves the URL's hostname BEFORE dialing, refuses if any
// returned IP is private, and emits HTTP.PRIVATE_BLOCKED with the resolved
// IPs in diagnostics. This blocks SSRF when DevForge runs as an MCP server
// invoked by AI agents.
package httpx

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/devforge/devforge/internal/netguard"
	"github.com/devforge/devforge/pkg/engine"
)

// Req is the request shape.
type Req struct {
	Method  string            `json:"method,omitempty"`  // default GET
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// Options tunes the runner.
type Options struct {
	// AllowPrivate permits the request to resolve to a private/loopback IP.
	// Off by default; opt-in for power users on internal networks.
	AllowPrivate bool `json:"allowPrivate,omitempty"`
	// TimeoutSec caps the request. Default 30, max 120.
	TimeoutSec int `json:"timeoutSec,omitempty"`
	// FollowRedirects up to MaxRedirects (default 5; 0 disables redirects).
	FollowRedirects bool `json:"followRedirects,omitempty"`
	MaxRedirects    int  `json:"maxRedirects,omitempty"`
	// MaxResponseBytes caps how much of the response body is read into memory.
	// 0 uses the package default (8 MiB). Hard ceiling 64 MiB to keep AI-agent
	// callers from triggering OOM on hostile/large endpoints.
	MaxResponseBytes int64 `json:"maxResponseBytes,omitempty"`
}

// Response is the result.
type Response struct {
	Status      int                 `json:"status"`
	StatusText  string              `json:"statusText"`
	Headers     map[string]string   `json:"headers"`
	Body        string              `json:"body"`
	BodyBytes   int                 `json:"bodyBytes"`
	BodyTruncated bool              `json:"bodyTruncated"`
	DurationMS  int64               `json:"durationMs"`
	ResolvedIPs []string            `json:"resolvedIps,omitempty"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Default + hard-ceiling caps for the response body buffer.
const (
	defaultBodyBytes int64 = 8 << 20  // 8 MiB
	maxBodyCeiling   int64 = 64 << 20 // 64 MiB hard ceiling
)

// Request executes an HTTP call with the netguard rules applied.
func Request(req Req, opts Options) (Response, error) {
	if strings.TrimSpace(req.URL) == "" {
		return Response{Diagnostics: []engine.Diagnostic{{
			Code: "HTTP.EMPTY_URL", Message: "url is required", Severity: engine.SevError,
		}}}, nil
	}
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = "GET"
	}
	parsed, err := url.Parse(req.URL)
	if err != nil {
		return Response{Diagnostics: []engine.Diagnostic{{
			Code: "HTTP.INVALID_URL", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return Response{Diagnostics: []engine.Diagnostic{{
			Code: "HTTP.UNSUPPORTED_SCHEME",
			Message: fmt.Sprintf("scheme %q not allowed (use http or https)", parsed.Scheme),
			Severity: engine.SevError,
		}}}, nil
	}
	host := parsed.Hostname()
	guard, gerr := netguard.ResolveSafe(context.Background(), host, 3*time.Second)
	if gerr != nil {
		return Response{Diagnostics: []engine.Diagnostic{{
			Code: "HTTP.RESOLVE", Message: gerr.Error(), Severity: engine.SevError,
		}}}, nil
	}
	resolved := make([]string, 0, len(guard.IPs))
	for _, ip := range guard.IPs {
		resolved = append(resolved, ip.String())
	}
	if guard.AnyPrivate && !opts.AllowPrivate {
		return Response{
			ResolvedIPs: resolved,
			Diagnostics: []engine.Diagnostic{{
				Code: "HTTP.PRIVATE_BLOCKED",
				Message: fmt.Sprintf(
					"host %q resolves to private IPs %v; pass allowPrivate=true if you really mean to dial them",
					host, guard.Private),
				Severity: engine.SevError,
			}},
		}, nil
	}

	timeout := time.Duration(opts.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if timeout > 120*time.Second {
		timeout = 120 * time.Second
	}

	// Build a Transport with a custom Dial that re-checks the resolved IP
	// (defense in depth: protects against DNS-rebinding between resolve and dial).
	trans := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			h, port, splitErr := net.SplitHostPort(addr)
			if splitErr != nil {
				return nil, splitErr
			}
			ip := net.ParseIP(h)
			if ip == nil {
				// addr came as host:port; resolve again.
				inner, ierr := netguard.ResolveSafe(ctx, h, 3*time.Second)
				if ierr != nil {
					return nil, ierr
				}
				if inner.AnyPrivate && !opts.AllowPrivate {
					return nil, fmt.Errorf("netguard: %v resolves to private IPs", h)
				}
				if len(inner.IPs) == 0 {
					return nil, fmt.Errorf("no IPs for %s", h)
				}
				ip = inner.IPs[0]
			} else if netguard.IsPrivate(ip) && !opts.AllowPrivate {
				return nil, fmt.Errorf("netguard: %s is private", ip)
			}
			d := net.Dialer{Timeout: 10 * time.Second}
			return d.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		},
	}
	client := &http.Client{
		Transport: trans,
		Timeout:   timeout,
	}
	if !opts.FollowRedirects {
		client.CheckRedirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if opts.MaxRedirects > 0 {
		max := opts.MaxRedirects
		client.CheckRedirect = func(_ *http.Request, via []*http.Request) error {
			if len(via) >= max {
				return fmt.Errorf("stopped after %d redirects", max)
			}
			return nil
		}
	}

	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}
	httpReq, err := http.NewRequest(method, req.URL, bodyReader)
	if err != nil {
		return Response{Diagnostics: []engine.Diagnostic{{
			Code: "HTTP.BUILD", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	if httpReq.Header.Get("User-Agent") == "" {
		httpReq.Header.Set("User-Agent", "devforge-httpx/0.0.0-dev")
	}

	start := time.Now()
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return Response{
			ResolvedIPs: resolved,
			Diagnostics: []engine.Diagnostic{{
				Code: "HTTP.SEND", Message: err.Error(), Severity: engine.SevError,
			}},
		}, nil
	}
	defer httpResp.Body.Close()
	limit := opts.MaxResponseBytes
	if limit <= 0 {
		limit = defaultBodyBytes
	}
	if limit > maxBodyCeiling {
		limit = maxBodyCeiling
	}
	// Read one byte past the limit to detect truncation.
	bodyBytes, _ := io.ReadAll(io.LimitReader(httpResp.Body, limit+1))
	truncated := int64(len(bodyBytes)) > limit
	if truncated {
		bodyBytes = bodyBytes[:limit]
	}

	headers := make(map[string]string, len(httpResp.Header))
	for k, v := range httpResp.Header {
		headers[k] = strings.Join(v, ", ")
	}
	return Response{
		Status:        httpResp.StatusCode,
		StatusText:    http.StatusText(httpResp.StatusCode),
		Headers:       headers,
		Body:          string(bodyBytes),
		BodyBytes:     len(bodyBytes),
		BodyTruncated: truncated,
		DurationMS:    time.Since(start).Milliseconds(),
		ResolvedIPs: resolved,
	}, nil
}
