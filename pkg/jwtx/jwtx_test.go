package jwtx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const hmacSecret = "test-secret-do-not-use"

func makeHS256Token(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(hmacSecret))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

// TestJwt_Decode_HS256_Header — C7: header alg parses as HS256.
func TestJwt_Decode_HS256_Header(t *testing.T) {
	t.Parallel()
	tok := makeHS256Token(t, jwt.MapClaims{"sub": "alok"})
	res, err := Decode(tok)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if res.Header["alg"] != "HS256" {
		t.Fatalf("alg = %v", res.Header["alg"])
	}
	if res.Payload["sub"] != "alok" {
		t.Fatalf("sub = %v", res.Payload["sub"])
	}
}

// TestJwt_Verify_RejectsAlgNone — C7: alg "none" is always rejected.
func TestJwt_Verify_RejectsAlgNone(t *testing.T) {
	t.Parallel()
	// Build a manual {alg:none} token without signing.
	hdr := `{"alg":"none","typ":"JWT"}`
	payload := `{"sub":"x"}`
	token := b64(hdr) + "." + b64(payload) + "."
	res, err := Verify(token, VerifyOptions{Key: []byte(hmacSecret)})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	if res.Diagnostics[0].Code != "JWT.ALG_NONE" {
		t.Fatalf("code=%q", res.Diagnostics[0].Code)
	}
}

// TestJwt_Verify_HS256_Valid — C7: a properly signed HS256 token verifies.
func TestJwt_Verify_HS256_Valid(t *testing.T) {
	t.Parallel()
	tok := makeHS256Token(t, jwt.MapClaims{"sub": "alok", "exp": time.Now().Add(time.Hour).Unix()})
	res, err := Verify(tok, VerifyOptions{Key: []byte(hmacSecret), ExpectedAlgs: []string{"HS256"}})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !res.Valid {
		t.Fatalf("expected valid: %+v", res.Diagnostics)
	}
	if res.ExpiresIn == nil || *res.ExpiresIn <= 0 {
		t.Fatalf("ExpiresIn = %v", res.ExpiresIn)
	}
}

// TestJwt_Verify_ExpiredReturnsDiagnostic — C7: expired token surfaces JWT.EXPIRED.
func TestJwt_Verify_ExpiredReturnsDiagnostic(t *testing.T) {
	t.Parallel()
	tok := makeHS256Token(t, jwt.MapClaims{"sub": "alok", "exp": time.Now().Add(-time.Hour).Unix()})
	res, err := Verify(tok, VerifyOptions{Key: []byte(hmacSecret), ExpectedAlgs: []string{"HS256"}})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	if res.Diagnostics[0].Code != "JWT.EXPIRED" {
		t.Fatalf("code=%q", res.Diagnostics[0].Code)
	}
}

// TestJwt_Verify_RS256_PEM — C7: RS256 verifies with a PEM public key.
func TestJwt_Verify_RS256_PEM(t *testing.T) {
	t.Parallel()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("genkey: %v", err)
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"sub": "alok"})
	signed, err := tok.SignedString(priv)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	pubBytes, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	res, err := Verify(signed, VerifyOptions{
		Key: pemBytes, KeyFormat: "pem", ExpectedAlgs: []string{"RS256"},
	})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !res.Valid {
		t.Fatalf("expected valid: %+v", res.Diagnostics)
	}
}

// TestJwt_Verify_AlgNotAllowed — C7: HS512 token blocked when allow-list = HS256.
func TestJwt_Verify_AlgNotAllowed(t *testing.T) {
	t.Parallel()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{"sub": "alok"})
	signed, _ := tok.SignedString([]byte(hmacSecret))
	res, err := Verify(signed, VerifyOptions{Key: []byte(hmacSecret), ExpectedAlgs: []string{"HS256"}})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	if !strings.Contains(res.Diagnostics[0].Code, "JWT.ALG_NOT_ALLOWED") {
		t.Fatalf("code=%q", res.Diagnostics[0].Code)
	}
}

func b64(s string) string {
	// Use raw url-encoding to match decoder.
	return strings.TrimRight(strings.NewReplacer("+", "-", "/", "_").Replace(
		base64Encode([]byte(s))), "=")
}

func base64Encode(b []byte) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	out := make([]byte, 0, ((len(b)+2)/3)*4)
	for i := 0; i < len(b); i += 3 {
		var x uint32
		var n int
		for j := 0; j < 3 && i+j < len(b); j++ {
			x = (x << 8) | uint32(b[i+j])
			n++
		}
		x <<= uint32(8 * (3 - n))
		out = append(out, charset[(x>>18)&0x3f], charset[(x>>12)&0x3f])
		if n >= 2 {
			out = append(out, charset[(x>>6)&0x3f])
		} else {
			out = append(out, '=')
		}
		if n == 3 {
			out = append(out, charset[x&0x3f])
		} else {
			out = append(out, '=')
		}
	}
	return string(out)
}
