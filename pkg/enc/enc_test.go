package enc

import (
	"strings"
	"testing"
	"testing/quick"
)

// ---------- Base64 RFC 4648 §10 vectors ----------

func TestBase64Encode_RFC4648_Vectors(t *testing.T) {
	t.Parallel()
	cases := []struct{ in, want string }{
		{"", ""},
		{"f", "Zg=="},
		{"fo", "Zm8="},
		{"foo", "Zm9v"},
		{"foob", "Zm9vYg=="},
		{"fooba", "Zm9vYmE="},
		{"foobar", "Zm9vYmFy"},
	}
	for _, c := range cases {
		r, err := Base64Encode([]byte(c.in), Base64Options{})
		if err != nil {
			t.Fatalf("Encode(%q): %v", c.in, err)
		}
		if r.Output != c.want {
			t.Fatalf("Encode(%q) = %q, want %q", c.in, r.Output, c.want)
		}
	}
}

func TestBase64Decode_RFC4648_Vectors(t *testing.T) {
	t.Parallel()
	cases := []struct{ in, want string }{
		{"", ""},
		{"Zg==", "f"},
		{"Zm9v", "foo"},
		{"Zm9vYmFy", "foobar"},
	}
	for _, c := range cases {
		r, err := Base64Decode(c.in, Base64Options{})
		if err != nil {
			t.Fatalf("Decode(%q): %v", c.in, err)
		}
		if string(r.Bytes) != c.want {
			t.Fatalf("Decode(%q) = %q, want %q", c.in, r.Bytes, c.want)
		}
	}
}

func TestBase64_URLSafe_NoPadding(t *testing.T) {
	t.Parallel()
	r, _ := Base64Encode([]byte{0xff, 0xfe, 0xfd}, Base64Options{URLSafe: true, NoPadding: true})
	if r.Output != "__79" {
		t.Fatalf("urlsafe = %q, want %q", r.Output, "__79")
	}
	d, _ := Base64Decode("__79", Base64Options{URLSafe: true, NoPadding: true})
	if string(d.Bytes) != "\xff\xfe\xfd" {
		t.Fatalf("decode mismatch: % x", d.Bytes)
	}
}

func TestBase64Decode_InvalidEmitsDiagnostic(t *testing.T) {
	t.Parallel()
	r, _ := Base64Decode("not!base64$", Base64Options{})
	if len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "ENC.BASE64.INVALID" {
		t.Fatalf("expected diagnostic, got %+v", r)
	}
}

// Property: encode → decode is identity for arbitrary bytes.
func TestBase64_RoundTrip_Property(t *testing.T) {
	t.Parallel()
	f := func(b []byte) bool {
		enc, _ := Base64Encode(b, Base64Options{})
		dec, _ := Base64Decode(enc.Output, Base64Options{})
		return string(dec.Bytes) == string(b)
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 200}); err != nil {
		t.Fatal(err)
	}
}

// ---------- URL ----------

func TestURLEncode_Component(t *testing.T) {
	t.Parallel()
	r, _ := URLEncode("a b&c=d", URLOptions{})
	if r.Output != "a+b%26c%3Dd" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestURLEncode_Path(t *testing.T) {
	t.Parallel()
	r, _ := URLEncode("a b/c", URLOptions{Mode: "path"})
	if r.Output != "a%20b%2Fc" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestURLDecode_RoundTrip(t *testing.T) {
	t.Parallel()
	f := func(s string) bool {
		// QueryEscape produces ASCII; restrict to printable to keep round-trip well-defined.
		clean := strings.Map(func(r rune) rune {
			if r < 0x20 || r > 0x7e {
				return -1
			}
			return r
		}, s)
		enc, _ := URLEncode(clean, URLOptions{})
		dec, _ := URLDecode(enc.Output)
		return dec.Output == clean
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 200}); err != nil {
		t.Fatal(err)
	}
}

// ---------- HTML ----------

func TestHTMLEncode_Decode(t *testing.T) {
	t.Parallel()
	in := `<a href="x?b=1&c=2">A&B</a>`
	enc, _ := HTMLEncode(in)
	if !strings.Contains(enc.Output, "&lt;") || !strings.Contains(enc.Output, "&amp;") {
		t.Fatalf("not entity-encoded: %s", enc.Output)
	}
	dec, _ := HTMLDecode(enc.Output)
	if dec.Output != in {
		t.Fatalf("round-trip lost: %q", dec.Output)
	}
}

// ---------- Hex ----------

func TestHexEncode_Lower(t *testing.T) {
	t.Parallel()
	r, _ := HexEncode([]byte{0x00, 0x0a, 0xff}, HexOptions{})
	if r.Output != "000aff" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestHexEncode_Upper(t *testing.T) {
	t.Parallel()
	r, _ := HexEncode([]byte{0xab, 0xcd}, HexOptions{Upper: true})
	if r.Output != "ABCD" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestHexDecode_With0xPrefix(t *testing.T) {
	t.Parallel()
	d, _ := HexDecode("0xDEADBEEF")
	if string(d.Bytes) != "\xde\xad\xbe\xef" {
		t.Fatalf("got % x", d.Bytes)
	}
}

func TestHexDecode_Invalid(t *testing.T) {
	t.Parallel()
	d, _ := HexDecode("zzz")
	if len(d.Diagnostics) == 0 || d.Diagnostics[0].Code != "ENC.HEX.INVALID" {
		t.Fatalf("expected diag, got %+v", d.Diagnostics)
	}
}

// Property: hex encode → decode identity.
func TestHex_RoundTrip_Property(t *testing.T) {
	t.Parallel()
	f := func(b []byte) bool {
		enc, _ := HexEncode(b, HexOptions{})
		dec, _ := HexDecode(enc.Output)
		return string(dec.Bytes) == string(b)
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 200}); err != nil {
		t.Fatal(err)
	}
}
