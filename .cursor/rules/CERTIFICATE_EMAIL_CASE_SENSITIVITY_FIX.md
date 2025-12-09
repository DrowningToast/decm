# Certificate Email Case Sensitivity Fix

## Problem Description

The `/api/v1/certificates/my-list-viewmodel` endpoint was not returning certificates for users who registered via Google OAuth, even when their email matched a certificate receiver email in the database.

### Root Cause

**Case-sensitive email comparison** in encrypted fields:

1. Certificate receivers were imported with emails in their original casing (e.g., `"User@Example.com"`)
2. Emails were encrypted using **deterministic AES-GCM** encryption
3. Users logging in with Google OAuth received a JWT with their email (e.g., `"user@example.com"`)
4. The query encrypted both emails and compared them directly
5. **Different casings produced different encrypted values**, causing no match

### Why This Happened

- Deterministic encryption (same plaintext → same ciphertext) allows for encrypted field comparison
- However, `"User@Example.com"` and `"user@example.com"` are different plaintexts
- Each produces a different ciphertext when encrypted
- PostgreSQL TEXT comparison is case-sensitive, so encrypted values don't match

## Solution Implemented

### 1. Normalize Emails to Lowercase

All emails are now normalized to lowercase **before encryption** to ensure consistent comparison:

#### Certificate Creation (`event_certificate.go`)

```go
// Normalize email to lowercase for case-insensitive comparison
normalizedEmail := params.ReceiverEmail
if params.ReceiverEmail != nil && *params.ReceiverEmail != "" {
    lowercaseEmail := strings.ToLower(*params.ReceiverEmail)
    normalizedEmail = &lowercaseEmail
}

receiverEmailEnc, err := pgmapper.EncryptStringPtrToPgText(normalizedEmail, r.piiEncryptionKey)
```

#### Certificate Queries (`event_certificate.go`)

```go
// GetClaimedCertificatesByCredentialID
// GetUnclaimedReadyCertificatesByCredentialID
normalizedEmail := email
if email != nil && *email != "" {
    lowercaseEmail := strings.ToLower(*email)
    normalizedEmail = &lowercaseEmail
}

encryptedEmail, err := pgmapper.EncryptStringPtrToPgText(normalizedEmail, r.piiEncryptionKey)
```

#### User Registration (`onboard.go`)

```go
// Normalize email to lowercase for case-insensitive comparison
lowercaseEmail := strings.ToLower(userInfo.Email)

credential = &entity.AuthenticationCredential{
    GoogleConnectorRef: &lowercaseEmail,
    // ...
}

// JWT token also uses normalized email
sessionToken, err := u.authService.CreateToken(auth.JwtPayload{
    Email: credential.GoogleConnectorRef,
    // ...
})
```

### 2. Files Modified

- `apps/backend/core-api/internal/repositories/postgres/event_certificate.go`
    - Added `strings` import
    - Normalized emails in `CreateEventCertificate`
    - Normalized emails in `GetClaimedCertificatesByCredentialID`
    - Normalized emails in `GetUnclaimedReadyCertificatesByCredentialID`

- `apps/backend/core-api/internal/usecase/onboard/onboard.go`
    - Added `strings` import
    - Normalized email in `RegisterWithGoogle` for both `GoogleConnectorRef` and JWT

## Migration Path

### For Existing Data

**Existing certificates with non-normalized emails will NOT automatically match** because:

1. They are already encrypted with the original casing
2. To fix them, the encrypted data would need to be:
    - Decrypted (requires encryption key)
    - Normalized to lowercase
    - Re-encrypted

### Recommended Actions

#### For Development/Testing

1. **Clear and re-import certificate receivers**:
    ```bash
    # Delete existing certificates for the event
    # Then re-import the receiver list
    ```

#### For Production

1. **Re-import certificate receiver lists** for affected events
    - The new import will use normalized lowercase emails
    - Existing certificates will be deleted and recreated (this is standard import behavior)

2. **For existing users**:
    - Users who registered before this fix will have non-normalized emails in their JWT
    - They need to log out and log back in to get a fresh JWT
    - Future logins will use normalized emails

3. **Monitor certificate claims**:
    - After the fix is deployed, new certificate imports will work correctly
    - Existing unclaimed certificates may need to be re-imported if users report issues

## Testing

### Test Case 1: New Certificate Import

1. Import certificate receivers with mixed-case email: `"User@Example.COM"`
2. Verify email is stored as lowercase in encrypted form (check behavior, not raw value)
3. Register user with Google OAuth using email: `"user@example.com"`
4. Call `/api/v1/certificates/my-list-viewmodel`
5. **Expected**: Certificate appears in `unclaimed_certificates`

### Test Case 2: Existing User Claims Certificate

1. Existing user with Google email: `"Test@Gmail.com"` (in JWT)
2. Import certificate for: `"test@gmail.com"`
3. User calls `/api/v1/certificates/my-list-viewmodel`
4. **Expected**: Certificate appears in `unclaimed_certificates`

### Test Case 3: Edge Cases

- Empty email: Should handle gracefully
- Null email: Should handle gracefully
- Special characters in email: Should normalize only letter casing

## Future Improvements

### Consider Database Migration Script (If Needed)

If re-importing is not feasible for production data:

```go
// Pseudo-code for migration script
func MigrateExistingEmails() {
    // 1. Get all certificates with receiver_email
    certificates := GetAllCertificates()

    // 2. For each certificate:
    for cert := range certificates {
        // Decrypt email
        email := DecryptEmail(cert.ReceiverEmail, encryptionKey)

        // Normalize to lowercase
        normalizedEmail := strings.ToLower(email)

        // Re-encrypt with normalized email
        encryptedEmail := EncryptEmail(normalizedEmail, encryptionKey)

        // Update certificate
        UpdateCertificateEmail(cert.ID, encryptedEmail)
    }
}
```

**Note**: This would require a one-time migration script with access to the PII encryption key.

## Related Documentation

- [Database Security Patterns](/.cursor/rules/database-security.mdc)
- [Encryption Utils](/.cursor/rules/encryption-utils.mdc)
- [Authentication Flows](/.cursor/rules/authentication-security.mdc)

## Summary

✅ **Fixed**: Email case sensitivity in certificate queries
✅ **Impact**: New certificates and users will work correctly
⚠️ **Action Required**: Re-import existing certificate receiver lists if needed
⚠️ **User Action**: Existing users should log out/in to get fresh JWT with normalized email

The fix ensures that email comparison is case-insensitive while maintaining data encryption security.
