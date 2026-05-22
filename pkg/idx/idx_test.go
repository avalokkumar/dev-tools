package idx

import (
	"strings"
	"testing"
)

func TestULID_LengthAndAlphabet(t *testing.T) {
	t.Parallel()
	r, err := ULID(ULIDOptions{Count: 5})
	if err != nil {
		t.Fatalf("ULID: %v", err)
	}
	if len(r.Values) != 5 {
		t.Fatalf("count = %d", len(r.Values))
	}
	for _, v := range r.Values {
		if len(v) != 26 {
			t.Fatalf("len = %d, want 26: %q", len(v), v)
		}
	}
}

func TestULID_Monotonic(t *testing.T) {
	t.Parallel()
	r, _ := ULID(ULIDOptions{Count: 32})
	// ULID strings sort lexically by timestamp; the first 10 chars encode time.
	for i := 1; i < len(r.Values); i++ {
		if r.Values[i][:10] < r.Values[i-1][:10] {
			t.Fatalf("non-monotonic timestamp: %s < %s", r.Values[i], r.Values[i-1])
		}
	}
}

func TestULID_Lowercase(t *testing.T) {
	t.Parallel()
	r, _ := ULID(ULIDOptions{Count: 1, Lowercase: true})
	if strings.ToLower(r.Values[0]) != r.Values[0] {
		t.Fatalf("not lowercase: %q", r.Values[0])
	}
}

func TestSlugify_Basic(t *testing.T) {
	t.Parallel()
	r, _ := Slugify("Hello, World!  Foo--Bar", SlugOptions{})
	if r.Output != "hello-world-foo-bar" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestSlugify_CustomSep(t *testing.T) {
	t.Parallel()
	r, _ := Slugify("a b c", SlugOptions{Sep: "_"})
	if r.Output != "a_b_c" {
		t.Fatalf("got %q", r.Output)
	}
}

func TestSlugify_MaxLen(t *testing.T) {
	t.Parallel()
	r, _ := Slugify("hello-world-foo-bar", SlugOptions{MaxLen: 11})
	if len([]rune(r.Output)) > 11 {
		t.Fatalf("len = %d: %q", len([]rune(r.Output)), r.Output)
	}
}

func TestSlugify_PreserveCase(t *testing.T) {
	t.Parallel()
	lower := false
	r, _ := Slugify("Hello World", SlugOptions{Lower: &lower})
	if r.Output != "Hello-World" {
		t.Fatalf("got %q", r.Output)
	}
}
