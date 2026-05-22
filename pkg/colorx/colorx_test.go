package colorx

import "testing"

func TestConvert_HexLong(t *testing.T) {
	t.Parallel()
	r, _ := Convert("#aabbcc", Options{})
	if r.R != 0xaa || r.G != 0xbb || r.B != 0xcc {
		t.Fatalf("RGB = %d,%d,%d", r.R, r.G, r.B)
	}
	if r.Hex != "#aabbcc" {
		t.Fatalf("hex = %s", r.Hex)
	}
}

func TestConvert_HexShort(t *testing.T) {
	t.Parallel()
	r, _ := Convert("#abc", Options{})
	if r.R != 0xaa || r.G != 0xbb || r.B != 0xcc {
		t.Fatalf("RGB = %d,%d,%d", r.R, r.G, r.B)
	}
}

func TestConvert_RGB(t *testing.T) {
	t.Parallel()
	r, _ := Convert("rgb(255, 0, 0)", Options{})
	if r.Hex != "#ff0000" {
		t.Fatalf("hex = %s", r.Hex)
	}
}

func TestConvert_HSL_Roundtrip(t *testing.T) {
	t.Parallel()
	r1, _ := Convert("#ff8800", Options{})
	// Convert back from HSL string.
	r2, _ := Convert(r1.HSL, Options{})
	if abs(r2.R-0xff) > 1 || abs(r2.G-0x88) > 1 || abs(r2.B-0x00) > 1 {
		t.Fatalf("round-trip drift: %s -> %s -> rgb(%d,%d,%d)", r1.Hex, r1.HSL, r2.R, r2.G, r2.B)
	}
}

func TestConvert_NamedColor(t *testing.T) {
	t.Parallel()
	r, _ := Convert("red", Options{})
	if r.Hex != "#ff0000" {
		t.Fatalf("got %s", r.Hex)
	}
}

func TestConvert_OutputField(t *testing.T) {
	t.Parallel()
	r, _ := Convert("#ff0000", Options{To: "hsl"})
	if r.Output != r.HSL {
		t.Fatalf("output = %s, hsl = %s", r.Output, r.HSL)
	}
}

func TestConvert_Invalid(t *testing.T) {
	t.Parallel()
	r, _ := Convert("not-a-color", Options{})
	if r.Diagnostics[0].Code != "COLOR.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
