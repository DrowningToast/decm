# PII Encryption - Quick Reference

## ⚡ 10-Second Summary

- **ALL PII encrypted at: Application layer (Repository) ONLY**
- **Algorithm: AES-256-GCM in Go code**
- **Key from: Environment variable `PII_ENCRYPTION_KEY`**
- **Encrypt before: INSERT/UPDATE (in Go code)**
- **Decrypt after: SELECT (in Go code)**
- **Use: `pgmapper` package**
- **NO database-level encryption (pgcrypto)**

## 🔐 PII Fields (Must Encrypt)

```
Email, First/Last Name, Phone, Address, Bio, Academic Fields
```

## 💾 Implementation Checklist

- [ ] PII fields stored as `TEXT` in database schema
- [ ] Encrypt using `pgmapper.EncryptStringPtrToPgText()` in Go before INSERT/UPDATE
- [ ] Decrypt using `pgmapper.DecryptPgTextToStringPtr()` in Go after SELECT
- [ ] Search by encrypting search term with `pgmapper.EncryptPII()` in Go
- [ ] NO database-level encryption - all encryption in Go application code
- [ ] `PII_ENCRYPTION_KEY` from environment only
- [ ] Test encrypt/decrypt roundtrip

## 📝 Copy-Paste Patterns

### CREATE

```go
emailEnc, _ := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
result, _ := r.queries.CreateProfile(ctx, generated.CreateProfileParams{Email: emailEnc})
emailDec, _ := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
```

### READ

```go
result, _ := r.queries.GetProfileByID(ctx, id)
emailDec, _ := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
```

### SEARCH

```go
encEmail, _ := pgmapper.EncryptPII(email, r.piiEncryptionKey)
result, _ := r.queries.GetProfileByEmail(ctx, pgtype.Text{String: encEmail, Valid: true})
```

### DELETE

```go
r.queries.DeleteProfile(ctx, id)  // No encryption needed
```

## 🛠️ Functions

| Function                     | Usage                        |
| ---------------------------- | ---------------------------- |
| `EncryptStringPtrToPgText()` | Encrypt string → pgtype.Text |
| `DecryptPgTextToStringPtr()` | Decrypt pgtype.Text → string |
| `EncryptPII()`               | Encrypt raw string (search)  |
| `DecryptPII()`               | Decrypt raw string (search)  |

## 🚫 Common Mistakes

```go
// ❌ WRONG - Encrypting in handler
handler: emailEnc := pgmapper.EncryptPII(req.Email, key)

// ✅ RIGHT - Encrypting in repository
repository: emailEnc := pgmapper.EncryptStringPtrToPgText(profile.Email, key)

// ❌ WRONG - Hardcoded key
encryptionKey := "my-secret-key"

// ✅ RIGHT - Environment variable
encryptionKey := r.piiEncryptionKey  // from config

// ❌ WRONG - Skipping encryption for "optional" fields
if profile.Email != nil {
    // MUST STILL ENCRYPT
}

// ✅ RIGHT - Always encrypt PII fields
encrypted, _ := pgmapper.EncryptStringPtrToPgText(profile.Email, key)
```

## 🔧 Setup

```bash
# 1. Generate encryption key
openssl rand -base64 32

# 2. Add to .env
PII_ENCRYPTION_KEY=<your-generated-key>

# 3. Initialize repository
repo := &Repository{
    piiEncryptionKey: config.PiiEncryptionKey,
}
```

## 📊 Encryption Flow

```
Handler (receive request)
    ↓
UseCase (business logic)
    ↓
Repository (APPLICATION LAYER)
├─→ Encrypt PII in Go (before INSERT/UPDATE)
├─→ Execute SQL with encrypted data
└─→ Decrypt PII in Go (after SELECT)
    ↓
UseCase (process decrypted data)
    ↓
Handler (return response)

Note: All encryption in Go code, NOT in database
```

## 🧪 Test Template

```go
func TestEncryption(t *testing.T) {
    key := os.Getenv("PII_ENCRYPTION_KEY")
    plaintext := "test@example.com"

    // Encrypt
    encrypted, err := pgmapper.EncryptPII(plaintext, key)
    assert.NoError(t, err)
    assert.NotEqual(t, plaintext, encrypted)

    // Decrypt
    decrypted, err := pgmapper.DecryptPII(encrypted, key)
    assert.NoError(t, err)
    assert.Equal(t, plaintext, decrypted)
}
```

## 🆘 Troubleshooting

| Error                           | Cause                  | Fix                           |
| ------------------------------- | ---------------------- | ----------------------------- |
| "ciphertext too short"          | Bad data/key           | Verify `PII_ENCRYPTION_KEY`   |
| "message authentication failed" | Wrong key or tampering | Check environment variable    |
| Returns garbage                 | Wrong key              | Verify encryption key matches |

## 📚 Full Documentation

See `documentations/pii-encryption.md` for:

- Complete implementation patterns (CREATE, READ, UPDATE, DELETE, SEARCH, BATCH)
- Database schema examples
- Error handling details
- Testing strategies
- Production considerations
- Migration guides

## 🎯 Key Rules

1. **Application Layer ONLY** - Never encrypt at handler/usecase/database
2. **Always Encrypt PII** - No exceptions for optional fields
3. **Use pgmapper** - Don't implement your own encryption
4. **Environment Variable Key** - Never hardcode keys
5. **Handle Errors** - Don't silently ignore encryption failures
6. **Test Roundtrips** - Verify encrypt/decrypt works
7. **Document PII** - Mark fields clearly in code comments
8. **No pgcrypto** - All encryption in Go application code

## 📦 Related Files

- `.cursorules` - Comprehensive project rules
- `documentations/pii-encryption.md` - Full guide
- `apps/backend/common/pgmapper/` - Implementation
- `apps/backend/common/encryptutils/` - Low-level functions
- `.env.example` - Configuration template
