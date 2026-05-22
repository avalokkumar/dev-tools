package totpx

import (
	"testing"
)

// RFC 6238 Appendix B vectors for SHA-1 (8-digit). Secret = ASCII "12345678901234567890".
//
// T (sec)        Counter T   TOTP    Algorithm
// 59             1           94287082  SHA1
// 1111111109     0x023523EC  07081804  SHA1
// 1111111111     0x023523ED  14050471  SHA1
// 1234567890     0x0273EF07  89005924  SHA1
// 2000000000     0x03F940AA  69279037  SHA1
const sha1Secret = "12345678901234567890" // raw ASCII, 20 bytes

func TestRFC6238_SHA1_Vectors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		at   int64
		want string
	}{
		{59, "94287082"},
		{1111111109, "07081804"},
		{1111111111, "14050471"},
		{1234567890, "89005924"},
		{2000000000, "69279037"},
	}
	for _, c := range cases {
		r, _ := Generate(GenerateOptions{
			Secret: sha1Secret, SecretEncoding: "raw",
			Algorithm: "sha1", Digits: 8, PeriodSec: 30, At: c.at,
		})
		if r.Code != c.want {
			t.Fatalf("at=%d code=%s want=%s", c.at, r.Code, c.want)
		}
	}
}

func TestVerify_AcceptsCurrent(t *testing.T) {
	t.Parallel()
	at := int64(1700000000)
	gen, _ := Generate(GenerateOptions{Secret: sha1Secret, SecretEncoding: "raw", Digits: 6, At: at})
	r, _ := Verify(gen.Code, VerifyOptions{
		GenerateOptions: GenerateOptions{Secret: sha1Secret, SecretEncoding: "raw", Digits: 6, At: at},
	})
	if r.Valid == nil || !*r.Valid {
		t.Fatalf("expected valid; %+v", r.Diagnostics)
	}
}

func TestVerify_RejectsWrongCode(t *testing.T) {
	t.Parallel()
	at := int64(1700000000)
	r, _ := Verify("000000", VerifyOptions{
		GenerateOptions: GenerateOptions{Secret: sha1Secret, SecretEncoding: "raw", Digits: 6, At: at},
	})
	if r.Valid == nil || *r.Valid {
		t.Fatalf("expected invalid; %+v", r)
	}
}

func TestVerify_AcceptsPrevStepWithinWindow(t *testing.T) {
	t.Parallel()
	prev := int64(1700000000)
	now := prev + 30 // exactly one step later
	gen, _ := Generate(GenerateOptions{Secret: sha1Secret, SecretEncoding: "raw", Digits: 6, At: prev})
	r, _ := Verify(gen.Code, VerifyOptions{
		GenerateOptions: GenerateOptions{Secret: sha1Secret, SecretEncoding: "raw", Digits: 6, At: now},
		WindowSteps:     1,
	})
	if r.Valid == nil || !*r.Valid {
		t.Fatalf("expected valid (within window)")
	}
}

func TestGenerate_UnsupportedDigits(t *testing.T) {
	t.Parallel()
	r, _ := Generate(GenerateOptions{Secret: "JBSWY3DPEHPK3PXP", Digits: 7})
	if r.Diagnostics[0].Code != "TOTP.UNSUPPORTED_DIGITS" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestGenerate_BadSecret(t *testing.T) {
	t.Parallel()
	r, _ := Generate(GenerateOptions{Secret: "!!!not-base32!!!"})
	if r.Diagnostics[0].Code != "TOTP.SECRET" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}
