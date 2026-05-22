# 03 — Page Specs: Generators

> Detailed page-level specifications for all Generator tools:
> UUID, Faker, ID, TOTP, and Crypto.

---

## UUID Generator (`/tools/uuid`)

### Overview

- **Category:** Generators
- **Icon:** `Key` (Lucide)
- **Operations:** `uuid_generate`, `uuid_hash`
- **Tabs:** "Generate" (default), "Hash"

### Tab: Generate

```text
┌──────────────────────────────────────────────────────────┐
│  🔑 UUID Generator                                       │
│  Generate v4 (random) or v7 (time-ordered) UUIDs        │
│                                                          │
│  [Generate]  [Hash]                                      │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  Version   [ v4 (random) ▾ ]                             │
│  Count     [ 5     ] (1–1024)                            │
│  Format    [ Standard ▾ ]  (std / compact / urn)         │
│                                                          │
│  [ ▶ Generate  ⌘↵ ]                                     │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  Results                                       [ Copy ]  │
│  ┌────────────────────────────────────────────────────┐  │
│  │ 01928d4a-7f3c-7abc-9d12-3e4f5a6b7c8d          📋 │  │
│  │ 01928d4a-7f3c-7abd-8e23-4f5a6b7c8d9e          📋 │  │
│  │ 01928d4a-7f3c-7abe-af34-5a6b7c8d9e0f          📋 │  │
│  │ 01928d4a-7f3c-7abf-b045-6b7c8d9e0f1a          📋 │  │
│  │ 01928d4a-7f3c-7ac0-c156-7c8d9e0f1a2b          📋 │  │
│  └────────────────────────────────────────────────────┘  │
│  5 UUIDs generated                                       │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Each UUID row has an inline copy button (clipboard icon) — click copies
  that single UUID. Toast: "UUID copied".
- "Copy" button copies all UUIDs newline-separated.
- Version dropdown: v4 (random), v7 (time-ordered).
- Format dropdown: Standard (`xxxxxxxx-xxxx-…`), Compact (no dashes),
  URN (`urn:uuid:…`).
- Count field: NumberStepper, min=1, max=1024, step=1.

### Tab: Hash

```text
│  Input     [________________________________]            │
│  Algorithms  ☑ SHA-256  ☑ MD5  ☐ SHA-1  ☐ SHA-512       │
│  Encoding    ( ● hex  ○ base64 )                         │
│                                                          │
│  [ ▶ Hash  ⌘↵ ]                                         │
│                                                          │
│  Results                                                 │
│  ┌────────────────────────────────────────────────────┐  │
│  │  SHA-256  e3b0c44298fc1c149afb…               📋  │  │
│  │  MD5      d41d8cd98f00b204e980…               📋  │  │
│  └────────────────────────────────────────────────────┘  │
```

**Interactions:**

- Checkboxes for algorithm selection (multi-select).
- Radio buttons for encoding (hex vs base64).
- Each result row: algorithm label + digest + copy button.

---

## Data Faker (`/tools/faker`)

### Overview

- **Category:** Generators
- **Icon:** `Drama` (Lucide)
- **Operations:** `faker_generate`, `faker_kinds`
- **Tabs:** "Generate" (default), "Field Kinds"

### Tab: Generate

```text
┌──────────────────────────────────────────────────────────┐
│  🎭 Data Faker                                           │
│  Generate realistic mock data in JSON, CSV, or SQL       │
│                                                          │
│  [Generate]  [Field Kinds]                               │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  Field Schema (JSON)                          [ Reset ]  │
│  ┌────────────────────────────────────────────────────┐  │
│  │ {                                                  │  │
│  │   "fields": [                                      │  │
│  │     { "name": "id", "kind": "uuid" },              │  │
│  │     { "name": "name", "kind": "name" },            │  │
│  │     { "name": "email", "kind": "email" },          │  │
│  │     { "name": "age", "kind": "int",                │  │
│  │       "params": { "min": 18, "max": 80 } }         │  │
│  │   ]                                                │  │
│  │ }                                                  │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Count   [ 10    ] (1–10,000)                            │
│  Format  [ JSON ▾ ]  (json / csv / sql)                  │
│  Table   [ data  ]  (SQL table name, only for sql fmt)   │
│  Seed    [ 0     ]  (0 = random)                         │
│                                                          │
│  [ ▶ Generate  ⌘↵ ]                                     │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  Output                    [ Copy ] [ Download ]         │
│  ┌────────────────────────────────────────────────────┐  │
│  │ [                                                  │  │
│  │   {"id":"…","name":"Alice","email":"…","age":34},  │  │
│  │   …                                               │  │
│  │ ]                                                  │  │
│  └────────────────────────────────────────────────────┘  │
│  10 rows generated                                       │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- CodeEditor for field schema (JSON syntax highlighting).
- Reset button restores default schema.
- Download button: `.json`, `.csv`, or `.sql` depending on format.
- Output rendered in CodeEditor (read-only) with appropriate language.

### Tab: Field Kinds

Fetches `faker_kinds` and displays a searchable reference table:

| Kind | Description | Params |
| --- | --- | --- |
| `uuid` | Random UUID | — |
| `name` | Full name | `locale?` |
| `email` | Email address | `domain?` |
| `int` | Random integer | `min`, `max` |
| `float` | Random float | `min`, `max`, `precision` |
| … | … | … |

---

## ID Generator (`/tools/id`)

### Overview

- **Category:** Generators
- **Icon:** `Tag` (Lucide)
- **Operations:** `id_ulid`, `id_slug`
- **Tabs:** "ULID", "Slug"

### Tab: ULID

```text
│  Count      [ 5 ]                                        │
│  Lowercase  [ ☑ ]                                        │
│                                                          │
│  [ ▶ Generate ]                                          │
│                                                          │
│  Results                                       [ Copy ]  │
│  01HX9YZ3K4M5N6P7Q8R9S0T1V2                        📋  │
│  01HX9YZ3K4M5N6P7Q8R9S0T1V3                        📋  │
```

### Tab: Slug

```text
│  Input      [ My Blog Post Title! ]                      │
│  Max Length  [ 60 ]                                       │
│  Locale     [ en ▾ ]                                      │
│                                                          │
│  [ ▶ Generate ]                                          │
│                                                          │
│  Result: my-blog-post-title                        📋   │
```

---

## TOTP Generator (`/tools/totp`)

### Overview

- **Category:** Generators
- **Icon:** `Hash` (Lucide)
- **Operations:** `totp_generate`, `totp_verify`
- **Tabs:** "Generate", "Verify"

### Tab: Generate

```text
┌──────────────────────────────────────────────────────────┐
│  🔢 TOTP Generator                                       │
│                                                          │
│  Secret          [JBSWY3DPEHPK3PXP          ]           │
│  Encoding        [ Base32 ▾ ]  (raw / hex / base32)      │
│  Algorithm       [ SHA-1 ▾ ]   (sha1 / sha256 / sha512)  │
│  Digits          ( ● 6   ○ 8 )                            │
│  Period (sec)    [ 30 ]                                    │
│  Time (unix)     [ auto (now) ]                           │
│                                                          │
│  [ ▶ Generate ]                                          │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │         ╔═══════════════╗                          │  │
│  │         ║   4  8  3  2  ║                          │  │
│  │         ║   9  1        ║                          │  │
│  │         ╚═══════════════╝                          │  │
│  │                                                    │  │
│  │  Code: 483291                              📋     │  │
│  │  Expires in: 18s  [████████░░░░░░░] 60%            │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

**Interactions:**

- Large-font TOTP code with animated countdown progress bar.
- Auto-refresh when `expiresIn` reaches 0 (call `totp_generate` again).
- Progress bar color: `--lemon-lime` when > 10s, `--vibrant-coral` when
  ≤ 10s (urgency).

### Tab: Verify

```text
│  Secret     [JBSWY3DPEHPK3PXP]                          │
│  Code       [483291            ]                          │
│  Skew       [ 1 ]   (allow ±N periods)                   │
│                                                          │
│  [ ▶ Verify ]                                            │
│                                                          │
│  Result:  ✅ Valid                                        │
```

- Valid: green check + `--lemon-lime` badge.
- Invalid: red X + `--vibrant-coral` badge.

---

## Crypto Toolkit (`/tools/crypto`)

### Overview

- **Category:** Generators / Security
- **Icon:** `ShieldCheck` (Lucide)
- **Operations:** 6 operations
- **Tabs:** "AES Encrypt", "AES Decrypt", "RSA Keygen", "HMAC",
  "Password Hash", "Password Strength"

### Tab: AES Encrypt

```text
│  Plaintext   [ CodeEditor — plaintext input ]            │
│  Mode        ( ● Passphrase  ○ Raw Key )                  │
│  Passphrase  [________________________]                   │
│  Key Size    [ 256 ▾ ]  (128 / 192 / 256)                │
│  PBKDF2 Itr  [ 100000 ]                                   │
│                                                          │
│  [ ▶ Encrypt ]                                           │
│                                                          │
│  Ciphertext (Base64)                             📋     │
│  ┌──────────────────────────────────────────────┐        │
│  │ SGVsbG8gV29ybGQ...                           │        │
│  └──────────────────────────────────────────────┘        │
```

### Tab: AES Decrypt

Mirror of encrypt — takes Base64 ciphertext, returns plaintext.

### Tab: RSA Keygen

```text
│  Key Size   ( ○ 2048  ● 3072  ○ 4096 )                  │
│                                                          │
│  [ ▶ Generate Key Pair ]                                 │
│                                                          │
│  Private Key (PEM)                              📋  ⬇   │
│  ┌──────────────────────────────────────────────┐        │
│  │ -----BEGIN RSA PRIVATE KEY-----               │        │
│  │ MIIEpAIBAAKCAQEA...                           │        │
│  │ -----END RSA PRIVATE KEY-----                 │        │
│  └──────────────────────────────────────────────┘        │
│                                                          │
│  Public Key (PEM)                               📋  ⬇   │
│  ┌──────────────────────────────────────────────┐        │
│  │ -----BEGIN PUBLIC KEY-----                    │        │
│  │ MIIBIjANBgkqhkiG...                          │        │
│  │ -----END PUBLIC KEY-----                      │        │
│  └──────────────────────────────────────────────┘        │
```

### Tab: HMAC

```text
│  Input          [ message to sign ]                      │
│  Key            [________________________]               │
│  Key Encoding   [ raw ▾ ]  (raw / hex / base64)           │
│  Algorithm      [ SHA-256 ▾ ]  (sha256 / sha384 / sha512) │
│                                                          │
│  [ ▶ Compute ]                                           │
│                                                          │
│  HMAC (hex): 5d41402abc4b2a76b971…                 📋   │
```

### Tab: Password Hash

```text
│  Password    [________________________]                   │
│  Algorithm   [ bcrypt ▾ ]  (bcrypt / argon2id)            │
│  Cost        [ 12 ]  (bcrypt: 4–31; argon2: mem/time)     │
│                                                          │
│  [ ▶ Hash ]                                              │
│                                                          │
│  Hash: $2a$12$LJ3m4ks...                           📋   │
```

### Tab: Password Strength

```text
│  Password    [________________________]    (live update)  │
│                                                          │
│  Score: ████████░░  4/4 — Very Strong                    │
│                                                          │
│  Crack Time: ~34 centuries                               │
│                                                          │
│  Feedback:                                               │
│    ✅ Good length (14+ chars)                            │
│    ✅ Contains uppercase, lowercase, digits, symbols      │
│    ⚠️ Consider avoiding common substitutions (@ for a)    │
```

**Interactions:**

- Score bar: gradient from `--vibrant-coral` (0) through
  `--accent-warning` (2) to `--lemon-lime` (4).
- Live update on keystroke (debounced 300ms).
- Reasons listed as inline diagnostics.
