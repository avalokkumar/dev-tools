package uuidx

import (
	"strings"
	"testing"
)

// TestGenerate_V4_Length36 — B1: a v4 UUID is 36 chars and parses.
func TestGenerate_V4_Length36(t *testing.T) {
	t.Parallel()
	res, err := Generate(GenerateOptions{Version: 4, Count: 1})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(res.Values) != 1 {
		t.Fatalf("len(Values) = %d, want 1", len(res.Values))
	}
	got := res.Values[0]
	if len(got) != 36 {
		t.Fatalf("len = %d, want 36; got %q", len(got), got)
	}
	// Position of dashes per RFC 4122.
	for _, i := range []int{8, 13, 18, 23} {
		if got[i] != '-' {
			t.Fatalf("expected dash at %d, got %q", i, got)
		}
	}
}

// TestGenerate_V4_CountThree — B1: Count=3 yields three distinct UUIDs.
func TestGenerate_V4_CountThree(t *testing.T) {
	t.Parallel()
	res, err := Generate(GenerateOptions{Version: 4, Count: 3})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(res.Values) != 3 {
		t.Fatalf("len = %d, want 3", len(res.Values))
	}
	seen := map[string]bool{}
	for _, v := range res.Values {
		if seen[v] {
			t.Fatalf("duplicate: %q", v)
		}
		seen[v] = true
	}
}

// TestGenerate_V7_Monotonic — B2: a sequence of v7 UUIDs increases lexically.
func TestGenerate_V7_Monotonic(t *testing.T) {
	t.Parallel()
	res, err := Generate(GenerateOptions{Version: 7, Count: 32})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	for i := 1; i < len(res.Values); i++ {
		if res.Values[i] <= res.Values[i-1] {
			t.Fatalf("not monotonic: [%d]=%s [%d]=%s",
				i-1, res.Values[i-1], i, res.Values[i])
		}
	}
}

// TestGenerate_RejectsUnknownVersion — B1: bad version returns Diagnostic.
func TestGenerate_RejectsUnknownVersion(t *testing.T) {
	t.Parallel()
	res, err := Generate(GenerateOptions{Version: 99, Count: 1})
	if err != nil {
		t.Fatalf("Generate err = %v, want nil (use Diagnostic)", err)
	}
	if len(res.Values) != 0 {
		t.Fatalf("want empty Values, got %v", res.Values)
	}
	if len(res.Diagnostics) == 0 {
		t.Fatalf("expected diagnostic")
	}
	if res.Diagnostics[0].Code != "UUID.UNSUPPORTED_VERSION" {
		t.Fatalf("code = %q", res.Diagnostics[0].Code)
	}
}

// TestGenerate_FormatCompact — B1: compact format strips dashes.
func TestGenerate_FormatCompact(t *testing.T) {
	t.Parallel()
	res, err := Generate(GenerateOptions{Version: 4, Count: 1, Format: "compact"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	got := res.Values[0]
	if len(got) != 32 {
		t.Fatalf("compact len = %d, want 32: %q", len(got), got)
	}
	if strings.Contains(got, "-") {
		t.Fatalf("compact contains dash: %q", got)
	}
}

// TestHash_Sha256_Hex — B3: sha256 hex matches the empty-string vector.
func TestHash_Sha256_Hex(t *testing.T) {
	t.Parallel()
	res, err := Hash(nil, HashOptions{Algos: []string{"sha256"}, Encoding: "hex"})
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	want := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got := res.Digests["sha256"]; got != want {
		t.Fatalf("sha256(\"\") = %q, want %q", got, want)
	}
}

// TestHash_MultiAlgo — B3: multiple algos all populate.
func TestHash_MultiAlgo(t *testing.T) {
	t.Parallel()
	res, err := Hash([]byte("devforge"), HashOptions{
		Algos:    []string{"md5", "sha1", "sha256", "sha512"},
		Encoding: "hex",
	})
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	for _, a := range []string{"md5", "sha1", "sha256", "sha512"} {
		if res.Digests[a] == "" {
			t.Fatalf("missing %q", a)
		}
	}
}

// TestHash_RejectsUnknownAlgo — B3: unknown algo emits diagnostic.
func TestHash_RejectsUnknownAlgo(t *testing.T) {
	t.Parallel()
	res, err := Hash(nil, HashOptions{Algos: []string{"snorefish"}})
	if err != nil {
		t.Fatalf("Hash err = %v", err)
	}
	if len(res.Diagnostics) == 0 || res.Diagnostics[0].Code != "HASH.UNSUPPORTED_ALGO" {
		t.Fatalf("expected HASH.UNSUPPORTED_ALGO, got %+v", res.Diagnostics)
	}
}
