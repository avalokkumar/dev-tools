package cryptox

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"

	"github.com/devforge/devforge/pkg/engine"
)

// PasswordHashOptions tunes PasswordHash.
type PasswordHashOptions struct {
	// Algorithm: "bcrypt" (default) or "argon2id".
	Algorithm string `json:"algorithm,omitempty"`
	// BcryptCost overrides the default 12 (range 4..31).
	BcryptCost int `json:"bcryptCost,omitempty"`
	// Argon2: time / memoryKiB / threads / keyLen
	Argon2Time      uint32 `json:"argon2Time,omitempty"`
	Argon2MemoryKiB uint32 `json:"argon2MemoryKiB,omitempty"`
	Argon2Threads   uint8  `json:"argon2Threads,omitempty"`
	Argon2KeyLen    uint32 `json:"argon2KeyLen,omitempty"`
}

// PasswordHashResult holds the encoded hash string.
type PasswordHashResult struct {
	Hash        string              `json:"hash"`
	Algorithm   string              `json:"algorithm"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// PasswordHash hashes a plaintext password. The output is the standard
// crypt-style encoded string (bcrypt) or argon2id PHC string.
func PasswordHash(password string, opts PasswordHashOptions) (PasswordHashResult, error) {
	algo := strings.ToLower(opts.Algorithm)
	if algo == "" {
		algo = "bcrypt"
	}
	switch algo {
	case "bcrypt":
		cost := opts.BcryptCost
		if cost == 0 {
			cost = 12
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
		if err != nil {
			return PasswordHashResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRYPTO.PASSWORD.HASH", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
		return PasswordHashResult{Hash: string(hashed), Algorithm: "bcrypt"}, nil
	case "argon2id":
		t := opts.Argon2Time
		if t == 0 {
			t = 3
		}
		m := opts.Argon2MemoryKiB
		if m == 0 {
			m = 64 * 1024 // 64 MiB
		}
		p := opts.Argon2Threads
		if p == 0 {
			p = 4
		}
		k := opts.Argon2KeyLen
		if k == 0 {
			k = 32
		}
		salt := make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			return PasswordHashResult{}, fmt.Errorf("cryptox: salt: %w", err)
		}
		dk := argon2.IDKey([]byte(password), salt, t, m, p, k)
		out := fmt.Sprintf(
			"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
			m, t, p,
			base64.RawStdEncoding.EncodeToString(salt),
			base64.RawStdEncoding.EncodeToString(dk),
		)
		return PasswordHashResult{Hash: out, Algorithm: "argon2id"}, nil
	default:
		return PasswordHashResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.PASSWORD.UNKNOWN_ALGO",
			Message: fmt.Sprintf("algorithm %q not supported", opts.Algorithm),
			Severity: engine.SevError,
		}}}, nil
	}
}

// PasswordStrengthResult ranks a password.
type PasswordStrengthResult struct {
	// Score 0 (very weak) … 4 (strong).
	Score        int      `json:"score"`
	EntropyBits  float64  `json:"entropyBits"`
	Length       int      `json:"length"`
	HasLower     bool     `json:"hasLower"`
	HasUpper     bool     `json:"hasUpper"`
	HasDigit     bool     `json:"hasDigit"`
	HasSymbol    bool     `json:"hasSymbol"`
	IsCommon     bool     `json:"isCommon"`
	Suggestions  []string `json:"suggestions"`
	Verdict      string   `json:"verdict"`
}

// PasswordStrength scores a password using a simple Shannon-style entropy
// approximation plus heuristics. Not a substitute for zxcvbn but adequate
// for "show users a green/yellow/red bar" use cases.
func PasswordStrength(password string) (PasswordStrengthResult, error) {
	res := PasswordStrengthResult{Length: len(password), Suggestions: []string{}}
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			res.HasLower = true
		case unicode.IsUpper(r):
			res.HasUpper = true
		case unicode.IsDigit(r):
			res.HasDigit = true
		default:
			res.HasSymbol = true
		}
	}
	pool := 0
	if res.HasLower {
		pool += 26
	}
	if res.HasUpper {
		pool += 26
	}
	if res.HasDigit {
		pool += 10
	}
	if res.HasSymbol {
		pool += 32
	}
	if pool > 0 {
		res.EntropyBits = float64(res.Length) * math.Log2(float64(pool))
	}
	if _, ok := commonPasswords[strings.ToLower(password)]; ok {
		res.IsCommon = true
	}
	switch {
	case res.IsCommon:
		res.Score = 0
	case res.EntropyBits < 40:
		res.Score = 1
	case res.EntropyBits < 55:
		res.Score = 2
	case res.EntropyBits < 75:
		res.Score = 3
	default:
		res.Score = 4
	}
	if res.Length < 12 {
		res.Suggestions = append(res.Suggestions, "use 12 or more characters")
	}
	if !res.HasUpper {
		res.Suggestions = append(res.Suggestions, "add uppercase letters")
	}
	if !res.HasDigit {
		res.Suggestions = append(res.Suggestions, "add digits")
	}
	if !res.HasSymbol {
		res.Suggestions = append(res.Suggestions, "add symbols")
	}
	if res.IsCommon {
		res.Suggestions = append(res.Suggestions, "this password appears in well-known leak lists; pick something else")
	}
	res.Verdict = []string{"very weak", "weak", "fair", "strong", "very strong"}[res.Score]
	return res, nil
}

// commonPasswords is a tiny built-in list. A full SecLists dump would bloat
// the binary; this catches the obvious cases and clearly documents the
// limitation.
var commonPasswords = map[string]struct{}{
	"password":    {},
	"123456":      {},
	"123456789":   {},
	"qwerty":      {},
	"abc123":      {},
	"letmein":     {},
	"admin":       {},
	"welcome":     {},
	"iloveyou":    {},
	"monkey":      {},
	"dragon":      {},
	"sunshine":    {},
	"princess":    {},
	"passw0rd":    {},
	"password1":   {},
	"qwerty123":   {},
	"trustno1":    {},
	"changeme":    {},
}
