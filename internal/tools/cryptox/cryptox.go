// Package cryptox is the Adapter that exposes pkg/cryptox as Operations.
package cryptox

import (
	"context"
	"encoding/json"

	"github.com/devforge/devforge/internal/mcpserver"
	enginepkg "github.com/devforge/devforge/pkg/cryptox"
)

func Operations() []mcpserver.Operation {
	return []mcpserver.Operation{
		{Tool: "crypto", Op: "aes_encrypt",
			Description: "AES-256-GCM encrypt with raw key or PBKDF2-derived passphrase.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["plaintext"],
  "properties":{"plaintext":{"type":"string"},
                "key":{"type":"string"},"keyEncoding":{"type":"string","enum":["hex","base64","raw"],"default":"hex"},
                "passphrase":{"type":"string"},"pbkdf2Iters":{"type":"integer","default":200000},
                "outputEncoding":{"type":"string","enum":["base64","hex"],"default":"base64"}}
}`),
			Handler: handleAESEncrypt,
		},
		{Tool: "crypto", Op: "aes_decrypt",
			Description: "AES-256-GCM decrypt; matches aes_encrypt bundle format.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input"],
  "properties":{"input":{"type":"string"},
                "key":{"type":"string"},"keyEncoding":{"type":"string","enum":["hex","base64","raw"],"default":"hex"},
                "passphrase":{"type":"string"},"pbkdf2Iters":{"type":"integer","default":200000},
                "inputEncoding":{"type":"string","enum":["base64","hex"],"default":"base64"},
                "outputEncoding":{"type":"string","enum":["utf8","base64","hex"],"default":"utf8"}}
}`),
			Handler: handleAESDecrypt,
		},
		{Tool: "crypto", Op: "rsa_keygen",
			Description: "Generate an RSA keypair (PEM, PKCS#8 + PKIX).",
			InputSchema: json.RawMessage(`{
  "type":"object",
  "properties":{"bits":{"type":"integer","enum":[2048,3072,4096],"default":2048}}
}`),
			Handler: handleRSAKeygen,
		},
		{Tool: "crypto", Op: "hmac",
			Description: "Compute HMAC over input.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["input","key"],
  "properties":{"input":{"type":"string"},"key":{"type":"string"},
                "algorithm":{"type":"string","enum":["sha1","sha256","sha384","sha512"],"default":"sha256"},
                "keyEncoding":{"type":"string","enum":["raw","hex","base64"],"default":"raw"},
                "outputEncoding":{"type":"string","enum":["hex","base64"],"default":"hex"}}
}`),
			Handler: handleHMAC,
		},
		{Tool: "crypto", Op: "password_hash",
			Description: "Hash a password with bcrypt or argon2id.",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["password"],
  "properties":{"password":{"type":"string"},
                "algorithm":{"type":"string","enum":["bcrypt","argon2id"],"default":"bcrypt"},
                "bcryptCost":{"type":"integer","minimum":4,"maximum":15,"default":12}}
}`),
			Handler: handlePasswordHash,
		},
		{Tool: "crypto", Op: "password_strength",
			Description: "Score a password (0=very weak … 4=very strong).",
			InputSchema: json.RawMessage(`{
  "type":"object","required":["password"],
  "properties":{"password":{"type":"string"}}
}`),
			Handler: handlePasswordStrength,
		},
	}
}

type aesEncArgs struct {
	Plaintext      string `json:"plaintext"`
	Key            string `json:"key,omitempty"`
	KeyEncoding    string `json:"keyEncoding,omitempty"`
	Passphrase     string `json:"passphrase,omitempty"`
	PBKDF2Iters    int    `json:"pbkdf2Iters,omitempty"`
	OutputEncoding string `json:"outputEncoding,omitempty"`
}

func handleAESEncrypt(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a aesEncArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.AESEncrypt([]byte(a.Plaintext), enginepkg.AESEncryptOptions{
		Key: a.Key, KeyEncoding: a.KeyEncoding,
		Passphrase: a.Passphrase, PBKDF2Iters: a.PBKDF2Iters,
		OutputEncoding: a.OutputEncoding,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type aesDecArgs struct {
	Input          string `json:"input"`
	Key            string `json:"key,omitempty"`
	KeyEncoding    string `json:"keyEncoding,omitempty"`
	Passphrase     string `json:"passphrase,omitempty"`
	PBKDF2Iters    int    `json:"pbkdf2Iters,omitempty"`
	InputEncoding  string `json:"inputEncoding,omitempty"`
	OutputEncoding string `json:"outputEncoding,omitempty"`
}

func handleAESDecrypt(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a aesDecArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.AESDecrypt(a.Input, enginepkg.AESDecryptOptions{
		Key: a.Key, KeyEncoding: a.KeyEncoding,
		Passphrase: a.Passphrase, PBKDF2Iters: a.PBKDF2Iters,
		InputEncoding: a.InputEncoding, OutputEncoding: a.OutputEncoding,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type rsaArgs struct {
	Bits int `json:"bits,omitempty"`
}

func handleRSAKeygen(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a rsaArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.RSAKeyGen(enginepkg.RSAKeyGenOptions{Bits: a.Bits})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type hmacArgs struct {
	Input          string `json:"input"`
	Key            string `json:"key"`
	Algorithm      string `json:"algorithm,omitempty"`
	KeyEncoding    string `json:"keyEncoding,omitempty"`
	OutputEncoding string `json:"outputEncoding,omitempty"`
}

func handleHMAC(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a hmacArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.HMAC([]byte(a.Input), enginepkg.HMACOptions{
		Algorithm: a.Algorithm, Key: a.Key,
		KeyEncoding: a.KeyEncoding, OutputEncoding: a.OutputEncoding,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type pwHashArgs struct {
	Password   string `json:"password"`
	Algorithm  string `json:"algorithm,omitempty"`
	BcryptCost int    `json:"bcryptCost,omitempty"`
}

func handlePasswordHash(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a pwHashArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.PasswordHash(a.Password, enginepkg.PasswordHashOptions{
		Algorithm: a.Algorithm, BcryptCost: a.BcryptCost,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

type pwStrengthArgs struct {
	Password string `json:"password"`
}

func handlePasswordStrength(_ context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a pwStrengthArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil { return nil, err }
	}
	res, err := enginepkg.PasswordStrength(a.Password)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}
