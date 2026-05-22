package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequest_PrivateBlockedByDefault(t *testing.T) {
	t.Parallel()
	r, _ := Request(Req{URL: "http://localhost:9/x"}, Options{})
	if r.Status != 0 {
		t.Fatalf("expected blocked (no Status), got %d", r.Status)
	}
	if len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "HTTP.PRIVATE_BLOCKED" {
		t.Fatalf("expected HTTP.PRIVATE_BLOCKED, got %+v", r.Diagnostics)
	}
}

func TestRequest_PrivateAllowedWhenOptIn(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	r, _ := Request(Req{URL: srv.URL}, Options{AllowPrivate: true})
	if r.Status != 200 {
		t.Fatalf("status = %d; diags=%+v", r.Status, r.Diagnostics)
	}
	if !strings.Contains(r.Body, "hello") {
		t.Fatalf("body = %q", r.Body)
	}
	if r.BodyBytes != 5 {
		t.Fatalf("BodyBytes = %d", r.BodyBytes)
	}
}

func TestRequest_UnsupportedScheme(t *testing.T) {
	t.Parallel()
	r, _ := Request(Req{URL: "file:///etc/passwd"}, Options{})
	if r.Diagnostics[0].Code != "HTTP.UNSUPPORTED_SCHEME" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestRequest_EmptyURL(t *testing.T) {
	t.Parallel()
	r, _ := Request(Req{URL: ""}, Options{})
	if r.Diagnostics[0].Code != "HTTP.EMPTY_URL" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestRequest_MethodAndHeaders(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s", r.Method)
		}
		if r.Header.Get("X-Test") != "yes" {
			t.Errorf("X-Test = %q", r.Header.Get("X-Test"))
		}
		w.WriteHeader(204)
	}))
	defer srv.Close()
	r, _ := Request(Req{Method: "POST", URL: srv.URL, Headers: map[string]string{"X-Test": "yes"}}, Options{AllowPrivate: true})
	if r.Status != 204 {
		t.Fatalf("status = %d; diags=%+v", r.Status, r.Diagnostics)
	}
}

// TestRequest_MaxResponseBytesTruncates — H-5: small cap → BodyTruncated true.
func TestRequest_MaxResponseBytesTruncates(t *testing.T) {
	t.Parallel()
	big := strings.Repeat("x", 5000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(big))
	}))
	defer srv.Close()
	r, _ := Request(Req{URL: srv.URL}, Options{AllowPrivate: true, MaxResponseBytes: 1024})
	if r.Status != 200 {
		t.Fatalf("status = %d", r.Status)
	}
	if !r.BodyTruncated {
		t.Fatalf("expected BodyTruncated=true")
	}
	if r.BodyBytes != 1024 {
		t.Fatalf("BodyBytes = %d, want 1024", r.BodyBytes)
	}
}

// TestRequest_DefaultBodyLimit — H-5: response under default cap NOT truncated.
func TestRequest_DefaultBodyLimit(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("small"))
	}))
	defer srv.Close()
	r, _ := Request(Req{URL: srv.URL}, Options{AllowPrivate: true})
	if r.BodyTruncated {
		t.Fatalf("expected not truncated")
	}
	if r.Body != "small" {
		t.Fatalf("body = %q", r.Body)
	}
}
