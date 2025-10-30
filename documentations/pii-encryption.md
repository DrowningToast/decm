# PII Encryption Guide - DECM Platform

## Overview

This document provides comprehensive guidance on encrypting Personally Identifiable Information (PII) in the DECM platform. All PII encryption uses **AES-256-GCM** (Galois/Counter Mode) encryption at the **application layer in the repository only** (not database-level encryption).

## Critical Principle

**ALL PII MUST be encrypted at the application layer in the repository - NEVER at handler, use case, or database layers.**

This ensures:

- Single responsibility principle
- Consistent encryption/decryption logic in application code
- Clear security boundaries
- No dependency on database-specific encryption features
- Easier auditing and maintenance
- Database portability

## PII Fields Definition

### Authentication Fields

- `google_connector_ref` - Google OAuth reference
- `github_connector_ref` - GitHub OAuth reference
- `wallet_address` - Blockchain wallet address (if considered sensitive)

### Profile Fields

- `first_name` - User's first name
- `last_name` - User's last name
- `email` - Email address
- `phone_number` - Phone number
- `address` - Physical address
- `bio` - User biography
- `profile_picture_url` - Profile image URL (if storing file paths)
- `academic_institution` - School/university name
- `academic_email` - Institutional email address

## Encryption Architecture

### Components

1. **Encryption Algorithm**: AES-256-GCM
    - 256-bit key (32 bytes)
    - Galois/Counter Mode for authenticated encryption
    - Detects any tampering with encrypted data
    - **Encryption happens in Go application code, NOT in database**

2. **Key Management**: Environment variables
    - `PII_ENCRYPTION_KEY` - Primary encryption key
    - Never hardcoded in source code
    - 32 bytes (256 bits) minimum
    - Managed at application level

3. **Implementation Location**: `apps/backend/common/pgmapper/`
    - Provides high-level encryption/decryption utilities
    - Built on top of `apps/backend/common/encryptutils/`
    - Type-safe wrappers for pgtype conversions
    - **All encryption logic in Go code**

4. **Database Storage**: TEXT columns
    - PII stored encrypted in database (already encrypted by application)
    - Base64-encoded for safe transmission
    - Nonce included in ciphertext
    - **Database stores encrypted strings, no database-level encryption**

### Encryption Flow Diagram

```
User Input
    ↓
Handler (Parse/Validate)
    ↓
UseCase (Business Logic)
    ↓
Repository Layer (APPLICATION LAYER ENCRYPTION)
    ├─→ Encrypt PII fields in Go (before INSERT/UPDATE)
    ├─→ Execute SQL queries with encrypted data
    └─→ Decrypt PII fields in Go (after SELECT)
    ↓
UseCase
    ↓
Handler (Return Response)
    ↓
Client

Note: All encryption/decryption happens in Go application code.
Database only stores the encrypted strings.
```

## Implementation Patterns

### 1. CREATE Operation (Insert with Encryption)

**Pattern**:

```go
func (r *Repository) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, error) {
    // 1. Encrypt PII fields
    emailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("email encryption failed: %w", err)
    }

    firstNameEnc, err := pgmapper.EncryptStringPtrToPgText(profile.FirstName, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("first name encryption failed: %w", err)
    }

    // 2. Insert encrypted data using sqlc-generated query
    query, err := r.queries.CreateProfile(ctx, generated.CreateProfileParams{
        Email:     emailEnc,
        FirstName: firstNameEnc,
        // ... other fields
    })
    if err != nil {
        return pgerrutils.ParsePgError(err)
    }

    // 3. Decrypt for return to usecase
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("email decryption failed: %w", err)
    }

    firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("first name decryption failed: %w", err)
    }

    return &entity.Profile{
        Email:     emailDec,
        FirstName: firstNameDec,
        // ... other fields
    }, nil
}
```

**SQL Query** (in `packages/database/queries/profiles.sql`):

```sql
-- name: CreateProfile :one
INSERT INTO profiles (id, email, first_name, last_name, created_at, updated_at)
VALUES (sqlc.arg(id), sqlc.arg(email), sqlc.arg(first_name), sqlc.arg(last_name), NOW(), NOW())
RETURNING *;
```

### 2. READ Operation (Select and Decrypt)

**Pattern**:

```go
func (r *Repository) GetProfile(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
    // 1. Query encrypted data
    query, err := r.queries.GetProfileByID(ctx, id)
    if err != nil {
        return nil, pgerrutils.ParsePgError(err)
    }

    // 2. Decrypt all PII fields
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("email decryption failed: %w", err)
    }

    firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("first name decryption failed: %w", err)
    }

    // 3. Return decrypted entity
    return &entity.Profile{
        ID:        query.ID,
        Email:     emailDec,
        FirstName: firstNameDec,
        // ... other fields
    }, nil
}
```

**SQL Query**:

```sql
-- name: GetProfileByID :one
SELECT id, email, first_name, last_name, created_at, updated_at
FROM profiles
WHERE id = sqlc.arg(id);
```

### 3. UPDATE Operation (Encrypt then Update)

**Pattern**:

```go
func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, updates entity.Profile) (*entity.Profile, error) {
    // 1. Encrypt PII fields to update
    emailEnc := pgtype.Text{}
    if updates.Email != nil {
        encrypted, err := pgmapper.EncryptStringPtrToPgText(updates.Email, r.piiEncryptionKey)
        if err != nil {
            return nil, fmt.Errorf("email encryption failed: %w", err)
        }
        emailEnc = encrypted
    }

    // 2. Update with encrypted data
    query, err := r.queries.UpdateProfile(ctx, generated.UpdateProfileParams{
        ID:    id,
        Email: emailEnc,
        // ... other encrypted fields
    })
    if err != nil {
        return nil, pgerrutils.ParsePgError(err)
    }

    // 3. Decrypt for return
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("email decryption failed: %w", err)
    }

    return &entity.Profile{
        ID:    query.ID,
        Email: emailDec,
        // ... other fields
    }, nil
}
```

**SQL Query**:

```sql
-- name: UpdateProfile :one
UPDATE profiles
SET email = COALESCE(sqlc.narg(email), email),
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;
```

### 4. SEARCH Operation (Deterministic Encryption)

**Pattern - Search by Email**:

```go
func (r *Repository) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error) {
    if email == "" {
        return nil, customerror.New(customerror.StatusBadRequest, "email required", nil)
    }

    // 1. Encrypt search term (deterministic - same plaintext = same ciphertext)
    encryptedEmail, err := pgmapper.EncryptPII(email, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("email encryption failed: %w", err)
    }

    // 2. Search with encrypted term
    query, err := r.queries.GetProfileByEmail(ctx, pgtype.Text{
        String: encryptedEmail,
        Valid:  true,
    })
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, customerror.New(customerror.StatusNotFound, "profile not found", nil)
        }
        return nil, pgerrutils.ParsePgError(err)
    }

    // 3. Decrypt result
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("email decryption failed: %w", err)
    }

    return &entity.Profile{
        ID:    query.ID,
        Email: emailDec,
        // ... other fields
    }, nil
}
```

**SQL Query**:

```sql
-- name: GetProfileByEmail :one
SELECT id, email, first_name, last_name, created_at, updated_at
FROM profiles
WHERE email = sqlc.arg(email);
```

### 5. DELETE Operation (No Decryption Needed)

**Pattern**:

```go
func (r *Repository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
    err := r.queries.DeleteProfile(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return customerror.New(customerror.StatusNotFound, "profile not found", nil)
        }
        return pgerrutils.ParsePgError(err)
    }
    return nil
}
```

**SQL Query**:

```sql
-- name: DeleteProfile :exec
DELETE FROM profiles WHERE id = sqlc.arg(id);
```

### 6. LIST/BATCH Operations (Decrypt All)

**Pattern**:

```go
func (r *Repository) ListProfiles(ctx context.Context, limit int) ([]entity.Profile, error) {
    // 1. Query encrypted data
    queries, err := r.queries.ListProfiles(ctx, int32(limit))
    if err != nil {
        return nil, pgerrutils.ParsePgError(err)
    }

    // 2. Decrypt all profiles
    profiles := make([]entity.Profile, len(queries))
    for i, q := range queries {
        emailDec, err := pgmapper.DecryptPgTextToStringPtr(q.Email, r.piiEncryptionKey)
        if err != nil {
            return nil, fmt.Errorf("email decryption failed for profile %d: %w", i, err)
        }

        profiles[i] = entity.Profile{
            ID:    q.ID,
            Email: emailDec,
            // ... other fields
        }
    }

    return profiles, nil
}
```

## Available Functions

### pgmapper Package Functions

```go
package pgmapper

// Encrypt string pointer to pgtype.Text
func EncryptStringPtrToPgText(field *string, encryptionKey string) (pgtype.Text, error)

// Decrypt pgtype.Text to string pointer
func DecryptPgTextToStringPtr(field pgtype.Text, encryptionKey string) (*string, error)

// Encrypt raw string (for search operations)
func EncryptPII(plaintext string, encryptionKey string) (string, error)

// Decrypt raw string (for search operations)
func DecryptPII(ciphertext string, encryptionKey string) (string, error)

// Type conversion utilities
func PgTextToStringPtr(pgText pgtype.Text) *string
func StringPtrToPgText(strPtr *string) pgtype.Text

func PgTimestampzToTimePtr(ts pgtype.Timestamptz) *time.Time
func TimePtrToPgTimestampz(t *time.Time) pgtype.Timestamptz

func PgUUIDToUUID(pgUUID pgtype.UUID) uuid.UUID
func UUIDToPgUUID(id uuid.UUID) pgtype.UUID
```

### encryptutils Package Functions

```go
package encryptutils

// Recommended: Simple AES-GCM with password
func EncryptAESGCM(plaintext string, password string) (string, error)
func DecryptAESGCM(ciphertext string, password string) (string, error)

// Advanced: AES with custom key (16, 24, or 32 bytes)
func EncryptAESWithKey(plaintext string, key []byte) (string, error)
func DecryptAESWithKey(ciphertext string, key []byte) (string, error)

// Utility functions
func GenerateAESKey() []byte
func KeyFromPassword(password string, salt []byte) []byte

// Deterministic encryption (same input = same output)
func EncryptDeterministicAES(plaintext string, password string) (string, error)
func DecryptDeterministicAES(ciphertext string, password string) (string, error)
```

## Database Schema

### PII Fields Definition

All PII fields MUST be stored as `TEXT` in the database:

```sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Standard fields (not encrypted)
    user_id UUID NOT NULL UNIQUE REFERENCES users(id),

    -- PII fields (stored encrypted as TEXT)
    email TEXT NOT NULL UNIQUE,
    first_name TEXT,
    last_name TEXT,
    phone_number TEXT,
    address TEXT,
    bio TEXT,
    profile_picture_url TEXT,
    academic_institution TEXT,
    academic_email TEXT,

    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_profiles_email ON profiles(email);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
```

**Important:** Email format validation must be performed in the application layer (e.g., using `validatorutils`) **before** encrypting the value and inserting or updating the database row. Database constraints cannot validate encrypted ciphertext against plaintext patterns.

### Searchable PII Fields

For fields you need to search on (email, username, etc.), use deterministic encryption:

- Same plaintext always encrypts to the same ciphertext
- Allows direct comparison in database
- Note: Deterministic encryption is less secure than random nonces - use carefully

## Environment Configuration

### Required Environment Variables

```bash
# .env file
PII_ENCRYPTION_KEY=your-256-bit-encryption-key-here
```

### Key Requirements

- **Length**: 32 bytes (256 bits) minimum
- **Format**: Can be any string (will be hashed to 256-bit key)
- **Security**: Generate with cryptographically secure random source
- **Rotation**: Plan key rotation strategy for production

### Generate Encryption Key

```bash
# Generate 32 random bytes (256-bit key) in base64
openssl rand -base64 32

# Result example:
# REDACTED_PII_ENCRYPTION_KEY
```

## Error Handling

### Common Encryption Errors

```go
// Encryption failed - invalid key
err := encryptutils.EncryptAESGCM(data, "invalid")
// -> Error: key must be 16, 24, or 32 bytes

// Decryption failed - wrong key
err := encryptutils.DecryptAESGCM(ciphertext, wrongKey)
// -> Error: ciphertext too short or authentication failed

// Data tampering detected
err := encryptutils.DecryptAESGCM(tamperedCiphertext, key)
// -> Error: cipher: message authentication failed
```

### Proper Error Handling

```go
// Encrypt
encrypted, err := pgmapper.EncryptStringPtrToPgText(field, key)
if err != nil {
    return fmt.Errorf("encryption failed: %w", err)
    // OR
    return customerror.New(customerror.StatusInternalError, "Encryption failed", err)
}

// Decrypt
decrypted, err := pgmapper.DecryptPgTextToStringPtr(field, key)
if err != nil {
    return fmt.Errorf("decryption failed: %w", err)
    // OR
    return customerror.New(customerror.StatusInternalError, "Data decryption failed", err)
}
```

## Testing Encryption

### Unit Test Example

```go
func TestProfileEncryption(t *testing.T) {
    // Setup
    encryptionKey := "test-key-32-bytes-long!!!!!!!!"
    plainEmail := "user@example.com"

    // Test encrypt → decrypt roundtrip
    encrypted, err := pgmapper.EncryptPII(plainEmail, encryptionKey)
    require.NoError(t, err)
    require.NotEqual(t, plainEmail, encrypted)

    decrypted, err := pgmapper.DecryptPII(encrypted, encryptionKey)
    require.NoError(t, err)
    require.Equal(t, plainEmail, decrypted)
}

func TestDeterministicEncryption(t *testing.T) {
    // Setup
    encryptionKey := "test-key-32-bytes-long!!!!!!!!"
    plainEmail := "user@example.com"

    // Same input should produce same output
    encrypted1, _ := encryptutils.EncryptDeterministicAES(plainEmail, encryptionKey)
    encrypted2, _ := encryptutils.EncryptDeterministicAES(plainEmail, encryptionKey)

    require.Equal(t, encrypted1, encrypted2)
}

func TestWrongKeyDecryption(t *testing.T) {
    // Setup
    key1 := "key-one-32-bytes-long!!!!!!!!!!"
    key2 := "key-two-32-bytes-long!!!!!!!!!!"
    plaintext := "sensitive data"

    // Encrypt with key1
    encrypted, _ := encryptutils.EncryptAESGCM(plaintext, key1)

    // Try to decrypt with key2 - should fail
    _, err := encryptutils.DecryptAESGCM(encrypted, key2)
    require.Error(t, err)
}
```

### Integration Test Example

```go
func TestRepositoryProfileEncryption(t *testing.T) {
    // Setup database and repository
    ctx := context.Background()
    repo := setupTestRepository(t)

    // Create profile with PII
    profile := &entity.Profile{
        Email:     stringPtr("test@example.com"),
        FirstName: stringPtr("John"),
        LastName:  stringPtr("Doe"),
    }

    // Create should encrypt
    created, err := repo.CreateProfile(ctx, *profile)
    require.NoError(t, err)
    require.NotNil(t, created)
    require.Equal(t, "test@example.com", *created.Email)

    // Retrieve should decrypt
    retrieved, err := repo.GetProfile(ctx, created.ID)
    require.NoError(t, err)
    require.Equal(t, "test@example.com", *retrieved.Email)
    require.Equal(t, "John", *retrieved.FirstName)

    // Verify encrypted in database
    var encryptedEmail string
    err = db.QueryRow("SELECT email FROM profiles WHERE id = $1", created.ID).Scan(&encryptedEmail)
    require.NoError(t, err)
    require.NotEqual(t, "test@example.com", encryptedEmail)
}
```

## Security Best Practices

### ✅ DO

- ✅ **Encrypt at repository layer** - single responsibility
- ✅ **Use pgmapper utilities** - tested, type-safe wrappers
- ✅ **Get key from environment** - never hardcode
- ✅ **Validate encryption key exists** - before repository initialization
- ✅ **Handle encryption errors** - don't silently ignore
- ✅ **Test encryption roundtrips** - verify encrypt/decrypt
- ✅ **Use AES-GCM** - authenticated encryption prevents tampering
- ✅ **Document PII fields** - mark clearly in code
- ✅ **Plan key rotation** - prepare strategy for production
- ✅ **Monitor encryption errors** - alert on decryption failures (possible tampering)

### ❌ DON'T

- ❌ **Never encrypt at handler layer** - violates separation of concerns
- ❌ **Never hardcode encryption keys** - security risk
- ❌ **Never skip encryption for "sensitive" fields** - all PII fields listed above
- ❌ **Never use weak algorithms** - use AES-GCM only
- ❌ **Never store keys in version control** - use environment variables
- ❌ **Never use user passwords as encryption keys** - use environment variable key
- ❌ **Never ignore decryption errors** - could indicate tampering or wrong key
- ❌ **Never mix encryption methods** - stick with pgmapper functions
- ❌ **Never encrypt in SQL** - all encryption in Go application code
- ❌ **Never use database-level encryption (pgcrypto)** - use application-layer encryption
- ❌ **Never reuse keys across different services** - use unique key per service/environment

## Production Considerations

### Key Management Strategy

1. **Key Generation**
    - Generate with cryptographically secure random source
    - 32 bytes (256 bits) minimum
    - Store in secure vault (HashiCorp Vault, AWS Secrets Manager, etc.)

2. **Key Rotation**
    - Plan rotation schedule (e.g., annually)
    - Maintain multiple keys during rotation period
    - Re-encrypt existing data with new key
    - Archive old keys for recovery

3. **Key Storage**
    - Use environment variables in containers
    - Inject from secrets manager at runtime
    - Never commit to version control
    - Restrict access to only necessary services

### Audit and Monitoring

- Log encryption/decryption failures
- Monitor for unusual decryption error patterns (possible tampering)
- Audit access to PII data
- Track key usage and rotation

### Performance Optimization

- Cache encryption key in memory during initialization
- Consider lazy decryption for bulk operations
- Profile encryption performance under load
- Use connection pooling for database operations

### Data Recovery

- Document encryption key locations
- Maintain secure key backups
- Plan recovery procedures for key loss
- Test recovery process regularly

## Migration Guide: Existing Data

### Encrypting Unencrypted PII

If migrating existing unencrypted data:

```go
// Migration function
func migrateUnencryptedProfiles(ctx context.Context, db *pgx.Conn, key string) error {
    // 1. Read unencrypted data
    rows, err := db.Query(ctx, "SELECT id, email, first_name FROM profiles")
    if err != nil {
        return err
    }

    // 2. Encrypt and update
    for rows.Next() {
        var id uuid.UUID
        var email, firstName string

        if err := rows.Scan(&id, &email, &firstName); err != nil {
            return err
        }

        // Encrypt
        encEmail, _ := pgmapper.EncryptPII(email, key)
        encName, _ := pgmapper.EncryptPII(firstName, key)

        // Update
        _, err := db.Exec(ctx,
            "UPDATE profiles SET email = $1, first_name = $2 WHERE id = $3",
            encEmail, encName, id,
        )
        if err != nil {
            return err
        }
    }

    return rows.Err()
}
```

## Troubleshooting

### Issue: "ciphertext too short"

- **Cause**: Invalid base64 decoding or corrupted data
- **Solution**: Verify encryption key matches, check database data integrity

### Issue: "cipher: message authentication failed"

- **Cause**: Data tampered with or wrong decryption key
- **Solution**: Verify encryption key, check data integrity, investigate access logs

### Issue: Decryption returns garbage

- **Cause**: Wrong encryption key or data corrupted
- **Solution**: Verify PII_ENCRYPTION_KEY environment variable

### Issue: Encryption performance slow

- **Cause**: High volume of PII fields or large values
- **Solution**: Profile code, optimize database queries, consider caching

## References

- [AES-GCM Wikipedia](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
- [Go crypto/cipher package](https://pkg.go.dev/crypto/cipher)
- [NIST Encryption Standards](https://csrc.nist.gov/publications/detail/sp/800-38d/final)
- [OWASP Encryption Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cryptographic_Storage_Cheat_Sheet.html)
