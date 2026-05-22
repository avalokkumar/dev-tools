package cryptox

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// ---------- AES ----------

func TestAES_Roundtrip_Passphrase(t *testing.T) {
	t.Parallel()
	enc, _ := AESEncrypt([]byte("secret payload"), AESEncryptOptions{Passphrase: "hunter2", PBKDF2Iters: 1000})
	if enc.Output == "" {
		t.Fatalf("empty output")
	}
	dec, _ := AESDecrypt(enc.Output, AESDecryptOptions{Passphrase: "hunter2", PBKDF2Iters: 1000})
	if dec.Output != "secret payload" {
		t.Fatalf("decrypt = %q, diags=%+v", dec.Output, dec.Diagnostics)
	}
}

func TestAES_Roundtrip_RawKey(t *testing.T) {
	t.Parallel()
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	keyHex := hex.EncodeToString(key)
	enc, _ := AESEncrypt([]byte("abc"), AESEncryptOptions{Key: keyHex})
	dec, _ := AESDecrypt(enc.Output, AESDecryptOptions{Key: keyHex})
	if dec.Output != "abc" {
		t.Fatalf("decrypt = %q", dec.Output)
	}
}

func TestAES_WrongPassphraseFails(t *testing.T) {
	t.Parallel()
	enc, _ := AESEncrypt([]byte("x"), AESEncryptOptions{Passphrase: "right", PBKDF2Iters: 1000})
	dec, _ := AESDecrypt(enc.Output, AESDecryptOptions{Passphrase: "wrong", PBKDF2Iters: 1000})
	if dec.Diagnostics[0].Code != "CRYPTO.AES.AUTH" {
		t.Fatalf("code = %q", dec.Diagnostics[0].Code)
	}
}

func TestAES_BadKeyLength(t *testing.T) {
	t.Parallel()
	_, _ = AESEncrypt([]byte("x"), AESEncryptOptions{Key: "aabb"})
	r, _ := AESEncrypt([]byte("x"), AESEncryptOptions{Key: "aabb"})
	if r.Diagnostics[0].Code != "CRYPTO.AES.KEY" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

// ---------- RSA ----------

func TestRSA_KeyGenPEMShapes(t *testing.T) {
	t.Parallel()
	r, err := RSAKeyGen(RSAKeyGenOptions{Bits: 2048})
	if err != nil {
		t.Fatalf("KeyGen: %v", err)
	}
	if !strings.Contains(r.PrivatePEM, "PRIVATE KEY") {
		t.Fatalf("private pem:\n%s", r.PrivatePEM)
	}
	if !strings.Contains(r.PublicPEM, "PUBLIC KEY") {
		t.Fatalf("public pem:\n%s", r.PublicPEM)
	}
	if r.Bits != 2048 {
		t.Fatalf("bits = %d", r.Bits)
	}
}

func TestRSA_UnsupportedBits(t *testing.T) {
	t.Parallel()
	r, _ := RSAKeyGen(RSAKeyGenOptions{Bits: 1024})
	if r.Diagnostics[0].Code != "CRYPTO.RSA.UNSUPPORTED_BITS" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

// ---------- HMAC RFC 4231 vectors ----------

func TestHMAC_SHA256_RFC4231_TestCase1(t *testing.T) {
	t.Parallel()
	// Test Case 1 from RFC 4231.
	keyHex := "0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b"
	want := "b0344c61d8db38535ca8afceaf0bf12b881dc200c9833da726e9376c2e32cff7"
	r, _ := HMAC([]byte("Hi There"), HMACOptions{
		Algorithm:   "sha256",
		Key:         keyHex,
		KeyEncoding: "hex",
	})
	if r.Output != want {
		t.Fatalf("got %s want %s", r.Output, want)
	}
}

func TestHMAC_SHA512_RFC4231_TestCase2(t *testing.T) {
	t.Parallel()
	// Test Case 2 (key="Jefe", data="what do ya want for nothing?").
	want := "164b7a7bfcf819e2e395fbe73b56e0a387bd64222e831fd610270cd7ea2505549758bf75c05a994a6d034f65f8f0e6fdcaeab1a34d4a6b4b636e070a38bce737"
	r, _ := HMAC([]byte("what do ya want for nothing?"), HMACOptions{
		Algorithm: "sha512",
		Key:       "Jefe",
	})
	if r.Output != want {
		t.Fatalf("got %s want %s", r.Output, want)
	}
}

func TestHMAC_UnsupportedAlgo(t *testing.T) {
	t.Parallel()
	r, _ := HMAC([]byte("x"), HMACOptions{Algorithm: "snorefish", Key: "k"})
	if r.Diagnostics[0].Code != "CRYPTO.HMAC.UNSUPPORTED_ALGO" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

// ---------- Password ----------

func TestPasswordHash_Bcrypt_RoundtripsWithCompare(t *testing.T) {
	t.Parallel()
	r, _ := PasswordHash("hunter2", PasswordHashOptions{Algorithm: "bcrypt", BcryptCost: 4})
	if !strings.HasPrefix(r.Hash, "$2") {
		t.Fatalf("not a bcrypt hash: %s", r.Hash)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(r.Hash), []byte("hunter2")); err != nil {
		t.Fatalf("compare: %v", err)
	}
}

func TestPasswordHash_Argon2id_PHCFormat(t *testing.T) {
	t.Parallel()
	r, _ := PasswordHash("hunter2", PasswordHashOptions{Algorithm: "argon2id"})
	if !strings.HasPrefix(r.Hash, "$argon2id$v=19$") {
		t.Fatalf("not a PHC argon2id hash: %s", r.Hash)
	}
}

func TestPasswordStrength_Tiers(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		in   string
		want int
	}{
		{"password", 0},
		{"abc", 1},
		{"abcdefgh", 1},
		{"Abcdefgh1!", 3},
		{"correcthorsebatterystaple9!ZZ", 4},
	} {
		r, _ := PasswordStrength(c.in)
		if r.Score != c.want {
			t.Fatalf("%q score = %d, want %d (entropy=%.1f)", c.in, r.Score, c.want, r.EntropyBits)
		}
	}
}
