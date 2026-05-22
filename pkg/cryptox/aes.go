// Package cryptox is the crypto suite (AES-GCM, RSA, HMAC, password +
// strength). Stdlib + golang.org/x/crypto only — no third-party algorithms.
package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	"github.com/devforge/devforge/pkg/engine"
)

// AESEncryptOptions tunes AES encryption.
type AESEncryptOptions struct {
	// Key is the raw 32-byte key (hex- or base64-encoded). Empty when Passphrase is set.
	Key string `json:"key,omitempty"`
	// KeyEncoding selects how Key is interpreted: "hex" (default), "base64", or "raw".
	KeyEncoding string `json:"keyEncoding,omitempty"`
	// Passphrase derives a key via PBKDF2-SHA256 (salt is 16 random bytes; iters from PBKDF2Iters).
	Passphrase string `json:"passphrase,omitempty"`
	// PBKDF2Iters defaults to 200_000 when Passphrase is set.
	PBKDF2Iters int `json:"pbkdf2Iters,omitempty"`
	// OutputEncoding for the ciphertext bundle: "base64" (default) or "hex".
	OutputEncoding string `json:"outputEncoding,omitempty"`
}

// AESResult is the success return for AES ops.
type AESResult struct {
	Output      string              `json:"output"`
	KDF         string              `json:"kdf,omitempty"`
	Salt        string              `json:"salt,omitempty"`
	Iters       int                 `json:"iters,omitempty"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// AESEncrypt encrypts plaintext with AES-256-GCM. Output is the base64
// (or hex) encoding of: [salt(16) || nonce(12) || ciphertext+tag] when a
// passphrase is used, or [nonce(12) || ciphertext+tag] when a raw key is used.
// The format is documented in the package README; the engine writes it
// deterministically so AESDecrypt can reverse it without a side channel.
func AESEncrypt(plaintext []byte, opts AESEncryptOptions) (AESResult, error) {
	key, salt, iters, err := resolveKey(opts.Key, opts.KeyEncoding, opts.Passphrase, opts.PBKDF2Iters, nil)
	if err != nil {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.KEY", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.CIPHER", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return AESResult{}, fmt.Errorf("cryptox: gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return AESResult{}, fmt.Errorf("cryptox: nonce: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	var bundle []byte
	if salt != nil {
		bundle = append(bundle, salt...)
	}
	bundle = append(bundle, nonce...)
	bundle = append(bundle, ciphertext...)

	res := AESResult{Output: encodeBundle(bundle, opts.OutputEncoding)}
	if salt != nil {
		res.KDF = "PBKDF2-SHA256"
		res.Salt = hex.EncodeToString(salt)
		res.Iters = iters
	}
	return res, nil
}

// AESDecryptOptions mirrors AESEncryptOptions plus the input ciphertext.
type AESDecryptOptions struct {
	Key            string `json:"key,omitempty"`
	KeyEncoding    string `json:"keyEncoding,omitempty"`
	Passphrase     string `json:"passphrase,omitempty"`
	PBKDF2Iters    int    `json:"pbkdf2Iters,omitempty"`
	InputEncoding  string `json:"inputEncoding,omitempty"`
	OutputEncoding string `json:"outputEncoding,omitempty"` // "utf8" (default) or "base64" or "hex"
}

// AESDecrypt is the inverse of AESEncrypt.
func AESDecrypt(input string, opts AESDecryptOptions) (AESResult, error) {
	bundle, err := decodeBundle(input, opts.InputEncoding)
	if err != nil {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.INPUT", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	usingPass := opts.Passphrase != ""
	const saltLen = 16
	const nonceLen = 12
	if usingPass && len(bundle) < saltLen+nonceLen {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.SHORT", Message: "ciphertext bundle is too short", Severity: engine.SevError,
		}}}, nil
	}
	if !usingPass && len(bundle) < nonceLen {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.SHORT", Message: "ciphertext bundle is too short", Severity: engine.SevError,
		}}}, nil
	}
	var salt []byte
	if usingPass {
		salt = bundle[:saltLen]
		bundle = bundle[saltLen:]
	}
	nonce := bundle[:nonceLen]
	ct := bundle[nonceLen:]

	key, _, _, err := resolveKey(opts.Key, opts.KeyEncoding, opts.Passphrase, opts.PBKDF2Iters, salt)
	if err != nil {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.KEY", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.CIPHER", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return AESResult{}, fmt.Errorf("cryptox: gcm: %w", err)
	}
	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return AESResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.AES.AUTH", Message: "authentication failed (wrong key/passphrase or tampered ciphertext)",
			Severity: engine.SevError,
		}}}, nil
	}
	out := encodePlain(plaintext, opts.OutputEncoding)
	return AESResult{Output: out}, nil
}

func resolveKey(rawKey, keyEnc, passphrase string, iters int, salt []byte) ([]byte, []byte, int, error) {
	if passphrase != "" {
		if iters <= 0 {
			iters = 200_000
		}
		if salt == nil {
			salt = make([]byte, 16)
			if _, err := rand.Read(salt); err != nil {
				return nil, nil, 0, err
			}
		}
		key := pbkdf2.Key([]byte(passphrase), salt, iters, 32, sha256.New)
		return key, salt, iters, nil
	}
	if rawKey == "" {
		return nil, nil, 0, fmt.Errorf("either key or passphrase is required")
	}
	var b []byte
	var err error
	switch strings.ToLower(keyEnc) {
	case "", "hex":
		b, err = hex.DecodeString(rawKey)
	case "base64":
		b, err = base64.StdEncoding.DecodeString(rawKey)
	case "raw":
		b = []byte(rawKey)
	default:
		return nil, nil, 0, fmt.Errorf("unknown keyEncoding %q", keyEnc)
	}
	if err != nil {
		return nil, nil, 0, err
	}
	if len(b) != 32 {
		return nil, nil, 0, fmt.Errorf("AES-256 requires a 32-byte key (got %d)", len(b))
	}
	return b, nil, 0, nil
}

func encodeBundle(b []byte, enc string) string {
	switch strings.ToLower(enc) {
	case "hex":
		return hex.EncodeToString(b)
	default:
		return base64.StdEncoding.EncodeToString(b)
	}
}

func decodeBundle(s, enc string) ([]byte, error) {
	switch strings.ToLower(enc) {
	case "hex":
		return hex.DecodeString(s)
	default:
		return base64.StdEncoding.DecodeString(s)
	}
}

func encodePlain(b []byte, enc string) string {
	switch strings.ToLower(enc) {
	case "base64":
		return base64.StdEncoding.EncodeToString(b)
	case "hex":
		return hex.EncodeToString(b)
	default:
		return string(b)
	}
}
