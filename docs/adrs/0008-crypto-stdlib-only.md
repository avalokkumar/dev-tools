# ADR 0008 — Crypto: stdlib + golang.org/x/crypto only

**Status:** Accepted
**Date:** 2026-05-10

## Context

`pkg/cryptox` and `pkg/totpx` add user-facing crypto operations: AES-GCM,
RSA, HMAC, bcrypt, argon2id, TOTP. Crypto bugs are catastrophic; the
provenance of cipher implementations matters.

## Decision

The crypto suite uses ONLY:

- Go standard library (`crypto/aes`, `crypto/cipher`, `crypto/rsa`,
  `crypto/hmac`, `crypto/sha1`, `crypto/sha256`, `crypto/sha512`,
  `crypto/x509`, `encoding/pem`, `encoding/base32`, `encoding/base64`,
  `encoding/hex`).
- `golang.org/x/crypto` for `bcrypt`, `argon2`, and `pbkdf2`.

No third-party cipher implementations. No CGO, no FFI, no shell-out.

## Consequences

- Anyone can audit DevForge's crypto by reading two well-known modules.
- `go.sum` carries the exact `golang.org/x/crypto` revision used per
  release; reproducible builds catch rotation.
- Test vectors are checked in (RFC 4231 HMAC, RFC 6238 TOTP, NIST GCM
  examples, bcrypt comparison via the upstream library).
- Future algorithms (XChaCha20, age, post-quantum) require an ADR before
  adoption to keep the surface conservative.
