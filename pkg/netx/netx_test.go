package netx

import "testing"

func TestURLParse_Basic(t *testing.T) {
	t.Parallel()
	r, _ := URLParse("https://user:pw@example.com:8443/api/v1?b=2&a=1#frag")
	if r.Scheme != "https" {
		t.Fatalf("scheme = %q", r.Scheme)
	}
	if !r.IsHTTPS {
		t.Fatalf("IsHTTPS false")
	}
	if r.Hostname != "example.com" || r.Port != "8443" {
		t.Fatalf("host/port = %q/%q", r.Hostname, r.Port)
	}
	if r.User != "user" {
		t.Fatalf("user = %q", r.User)
	}
	if len(r.Params) != 2 || r.Params[0].Key != "a" {
		t.Fatalf("params = %+v", r.Params)
	}
	if r.Fragment != "frag" {
		t.Fatalf("frag = %q", r.Fragment)
	}
}

func TestURLParse_Empty(t *testing.T) {
	t.Parallel()
	r, _ := URLParse("")
	if r.Diagnostics[0].Code != "URL.PARSE.EMPTY" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestURLParse_Invalid(t *testing.T) {
	t.Parallel()
	r, _ := URLParse("http://[bad")
	if r.Diagnostics[0].Code != "URL.PARSE.INVALID" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestHeadersAnalyze_AllPresent(t *testing.T) {
	t.Parallel()
	in := map[string]string{
		"Strict-Transport-Security": "max-age=31536000",
		"Content-Security-Policy":   "default-src 'self'",
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"Referrer-Policy":           "no-referrer",
		"Permissions-Policy":        "geolocation=()",
	}
	r, _ := HeadersAnalyze(in)
	for _, f := range r.Findings {
		if !f.OK {
			t.Fatalf("expected ok for %s: %s", f.Header, f.Note)
		}
	}
}

func TestHeadersAnalyze_MissingFlagsFail(t *testing.T) {
	t.Parallel()
	r, _ := HeadersAnalyze(map[string]string{})
	missing := 0
	for _, f := range r.Findings {
		if !f.OK {
			missing++
		}
	}
	if missing != len(r.Findings) {
		t.Fatalf("expected all missing, got %d/%d", missing, len(r.Findings))
	}
}

func TestHeadersAnalyze_WrongValueFlags(t *testing.T) {
	t.Parallel()
	r, _ := HeadersAnalyze(map[string]string{
		"X-Content-Type-Options": "wrong",
	})
	for _, f := range r.Findings {
		if f.Header == "X-Content-Type-Options" {
			if f.OK {
				t.Fatalf("expected not ok")
			}
		}
	}
}
