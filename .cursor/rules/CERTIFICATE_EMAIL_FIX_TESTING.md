# Certificate Email Case Sensitivity Fix - Testing Guide

## Quick Summary

Fixed the issue where `/api/v1/certificates/my-list-viewmodel` doesn't return certificates for users with matching Google OAuth emails due to case-sensitive email comparison in encrypted fields.

## What Was Changed

### Code Changes

1. **Email normalization to lowercase** before encryption in:
    - Certificate creation (`event_certificate.go::CreateEventCertificate`)
    - Certificate queries (`event_certificate.go::GetClaimedCertificatesByCredentialID`, `GetUnclaimedReadyCertificatesByCredentialID`)
    - User registration (`onboard.go::RegisterWithGoogle`)

2. **Files modified**:
    - `apps/backend/core-api/internal/repositories/postgres/event_certificate.go`
    - `apps/backend/core-api/internal/usecase/onboard/onboard.go`

### Build Status

✅ **Compilation successful** - No errors
✅ **Formatting applied** - Code follows Go standards

## Testing Steps

### Prerequisites

```bash
# Start the backend
pnpm dev:core

# Start the frontend (in another terminal)
pnpm dev
```

### Test Case 1: New Certificate Import (Recommended First Test)

1. **As Organizer**: Import certificate receivers with mixed-case email

    ```
    POST /api/v1/events/{event_id}/certificates/import
    {
      "receivers": [
        {
          "email": "Test.User@Example.COM",  // Note: Mixed case
          "first_name": "Test",
          "last_name": "User"
        }
      ]
    }
    ```

2. **As New User**: Register with Google OAuth using lowercase email
    - Use Google account: `test.user@example.com`
    - Complete registration flow

3. **Publish Certificate Config**:

    ```
    POST /api/v1/event-config/{event_id}/certificate/publish
    ```

4. **Verify Certificate Appears**:

    ```
    GET /api/v1/certificates/my-list-viewmodel
    ```

    **Expected Response**:

    ```json
    {
        "claimed_certificates": [],
        "unclaimed_certificates": [
            {
                "receiver_email": "<encrypted>",
                "event_id": "..."
                // ... certificate details
            }
        ],
        "total_claimed": 0,
        "total_unclaimed": 1
    }
    ```

### Test Case 2: Existing User with New Certificate

**Scenario**: User already registered with Google OAuth

1. **Import certificate** for existing user's email (any casing)
2. **Publish certificate config**
3. **User logs out and logs back in** (to get fresh JWT with normalized email)
4. **Call** `/api/v1/certificates/my-list-viewmodel`
5. **Expected**: Certificate appears in `unclaimed_certificates`

### Test Case 3: Claim Certificate

After Test Case 1:

1. **User claims certificate**:

    ```
    POST /api/v1/certificates/{certificate_id}/claim
    ```

2. **Verify it moves to claimed**:

    ```
    GET /api/v1/certificates/my-list-viewmodel
    ```

    **Expected**:

    ```json
    {
        "claimed_certificates": [
            /* certificate here */
        ],
        "unclaimed_certificates": [],
        "total_claimed": 1,
        "total_unclaimed": 0
    }
    ```

## Debugging

### If certificates still don't appear:

1. **Check JWT token contains email**:
    - Decode JWT token from browser cookies
    - Verify `email` field is present and lowercase
    - If not present, user needs to log out/in

2. **Check certificate config is published**:

    ```sql
    SELECT is_published FROM event_certificate_configs WHERE event_id = '<event_id>';
    ```

    - Should be `TRUE`

3. **Check certificate is unclaimed**:

    ```sql
    SELECT certificate_token_id FROM event_certificates WHERE receiver_email = '<encrypted_email>';
    ```

    - Should be `NULL`

4. **Check certificate is not revoked**:

    ```sql
    SELECT revoked_at FROM event_certificates WHERE id = '<certificate_id>';
    ```

    - Should be `NULL`

5. **Verify email normalization** (add temporary logging):
    ```go
    // In event_certificate.go::GetUnclaimedReadyCertificatesByCredentialID
    log.Printf("Searching for email: %v, normalized: %v", email, normalizedEmail)
    ```

## Known Limitations

### Existing Data

- **Certificates imported before this fix** will have non-normalized encrypted emails
- They will NOT match until the receiver list is re-imported
- **Solution**: Re-import certificate receiver lists for affected events

### Existing Users

- Users registered before this fix may have non-normalized emails in JWT
- **Solution**: Log out and log back in to get fresh JWT

## Rollout Plan

### Development/Staging

1. Deploy the fix
2. Clear all test certificates
3. Re-import test data
4. Run test cases above

### Production

1. **Deploy the fix** during low-traffic period
2. **Notify organizers** to re-import receiver lists if they have issues
3. **Monitor** certificate claim success rate
4. **Support**: If users report missing certificates:
    - Ask them to log out/in
    - If still not working, organizer should re-import receivers

## Success Criteria

✅ Mixed-case emails in import match lowercase emails in JWT
✅ Users can see their unclaimed certificates
✅ Users can claim certificates successfully
✅ No regression in existing functionality

## Related Files

- Fix Documentation: `CERTIFICATE_EMAIL_CASE_SENSITIVITY_FIX.md`
- Modified Code:
    - `apps/backend/core-api/internal/repositories/postgres/event_certificate.go`
    - `apps/backend/core-api/internal/usecase/onboard/onboard.go`

## Contact

If you encounter issues during testing, check:

1. Backend logs for any encryption/decryption errors
2. Database logs for query performance
3. Frontend console for API errors
