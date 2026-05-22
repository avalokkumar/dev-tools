package cryptox

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/devforge/devforge/pkg/engine"
)

// RSAKeyGenOptions tunes RSAKeyGen.
type RSAKeyGenOptions struct {
	// Bits must be 2048, 3072, or 4096. Default 2048.
	Bits int `json:"bits,omitempty"`
}

// RSAKeyGenResult holds PEM-encoded private + public keys.
type RSAKeyGenResult struct {
	PrivatePEM  string              `json:"privatePem"`
	PublicPEM   string              `json:"publicPem"`
	Bits        int                 `json:"bits"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// RSAKeyGen generates a new RSA keypair and returns both halves as PEM.
func RSAKeyGen(opts RSAKeyGenOptions) (RSAKeyGenResult, error) {
	bits := opts.Bits
	if bits == 0 {
		bits = 2048
	}
	switch bits {
	case 2048, 3072, 4096:
	default:
		return RSAKeyGenResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRYPTO.RSA.UNSUPPORTED_BITS",
			Message: fmt.Sprintf("bits %d not supported (use 2048, 3072, or 4096)", bits),
			Severity: engine.SevError,
		}}}, nil
	}
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return RSAKeyGenResult{}, fmt.Errorf("cryptox: keygen: %w", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: mustMarshalPKCS8(priv),
	})
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return RSAKeyGenResult{}, fmt.Errorf("cryptox: pubkey: %w", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	return RSAKeyGenResult{
		PrivatePEM: string(privPEM),
		PublicPEM:  string(pubPEM),
		Bits:       bits,
	}, nil
}

func mustMarshalPKCS8(priv *rsa.PrivateKey) []byte {
	b, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		// PKCS#1 is always valid for RSA — fall back rather than panic.
		return x509.MarshalPKCS1PrivateKey(priv)
	}
	return b
}
