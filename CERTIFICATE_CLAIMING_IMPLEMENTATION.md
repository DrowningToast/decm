# Certificate Claiming Mechanism - Implementation Summary

## Overview

Implemented a complete certificate claiming mechanism on the frontend that allows participants to claim their certificates from the certificate detail page. The implementation follows the same pattern used in the event registration flow with password/PIN modal authentication.

## Files Created/Modified

### 1. New Hook: `useClaimCertificate.ts`

**Location:** `/apps/web/src/hooks/useClaimCertificate.ts`

**Purpose:** Custom React hook for claiming certificates with account password verification

**Features:**

- Uses React Query mutation for async certificate claiming
- Integrates with toast notifications for user feedback
- Automatically invalidates certificate and inbox queries after successful claim
- Includes comprehensive documentation for backend integration
- Mock implementation for development/testing

**API Structure (To be implemented on backend):**

```typescript
POST /api/v1/events/{event_id}/certificates/{certificate_id}/claim
Body: { account_password: string }
Response: {
  certificate_id: string,
  certificate_token_id: string,
  transaction_hash: string,
  claimed_at: string
}
```

### 2. Updated Component: `CertificateDetail.tsx`

**Location:** `/apps/web/src/components/pages/Participant/Certificates/CertificateDetail.tsx`

**Changes:**

- Added certificate claiming button UI
- Integrated password/PIN prompt modal using `usePasswordPrompt` hook
- Added loading states and disabled states during claiming
- Shows different UI based on claim status:
    - "Claim Certificate" button for unclaimed certificates
    - "Certificate Claimed" status badge for claimed certificates
- Added proper error handling with user feedback

**New Imports:**

```typescript
import { useState } from "react";
import { Award, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { usePasswordPrompt } from "@/hooks/usePassowordPrompt";
import { useClaimCertificate } from "@/hooks/useClaimCertificate";
```

### 3. Updated Translations: `en.json`

**Location:** `/apps/web/src/lib/i18n/locales/en.json`

**Added Keys:**

```json
{
    "participant.certificates": {
        "claimTitle": "Claim Certificate",
        "claimDescription": "Sign to claim your certificate on the blockchain",
        "claimButton": "Claim Certificate",
        "claiming": "Claiming...",
        "claimSuccess": "Certificate claimed successfully!",
        "claimError": "Failed to claim certificate. Please try again.",
        "claimNote": "Sign this certificate to add it to your blockchain credentials",
        "claimed": "Certificate Claimed"
    }
}
```

## User Flow

1. **Navigate to Certificate Detail Page**
    - User views a certificate from their certificates list or inbox
    - If certificate is not yet claimed (status !== "completed"), they see a "Claim Certificate" button

2. **Click Claim Certificate**
    - Button triggers `handleClaimCertificate` function
    - Loading state activates, button becomes disabled

3. **Password/PIN Verification**
    - Password/PIN modal pops up using the existing `usePasswordPrompt` hook
    - Same modal used in event registration flow
    - Shows signing details:
        - Contract Address (if available)
        - Transaction Type: "Certificate Claim"
        - Details: "Claiming certificate: {certificate name}"

4. **Password Verification & Claiming**
    - User enters their PIN (6 digits) or password
    - System verifies the password against stored credentials
    - If valid, proceeds to claim certificate
    - If invalid, shows error message

5. **Certificate Claimed**
    - Success toast notification appears
    - Button is replaced with "Certificate Claimed" badge
    - Certificate queries are invalidated and refetched
    - Inbox messages are updated to show "Claimed" status

## Technical Integration Points

### Password/PIN Modal

The implementation reuses the existing password prompt system:

```typescript
const accountPassword = await openPasswordPrompt({
    eventContractAddress: certificate.eventContractAddress || "",
    transactionType: "Certificate Claim",
    title: t("participant.certificates.claimTitle"),
    description: t("participant.certificates.claimDescription"),
    details: `Claiming certificate: ${certificate.name}`,
});
```

### Certificate Status Check

```typescript
const isCertificateClaimed = certificate?.status === "completed";
const canClaimCertificate = !isCertificateClaimed && !isClaiming && !isProcessing;
```

### Query Invalidation

After successful claim, these queries are automatically refetched:

```typescript
queryClient.invalidateQueries({ queryKey: QUERY_KEY.certificate.all });
queryClient.invalidateQueries({ queryKey: QUERY_KEY.inbox.all });
```

## Backend Requirements

### Required API Endpoint

**Endpoint:** `POST /api/v1/events/{event_id}/certificates/{certificate_id}/claim`

**Request Body:**

```go
type ClaimCertificateRequest struct {
    AccountPassword string `json:"account_password" binding:"required"`
}
```

**Response Body:**

```go
type ClaimCertificateResponse struct {
    CertificateID      string `json:"certificate_id"`
    CertificateTokenID string `json:"certificate_token_id"`
    TransactionHash    string `json:"transaction_hash"`
    ClaimedAt          string `json:"claimed_at"`
}
```

### Backend Logic Steps:

1. **Authenticate User**
    - Verify JWT token from cookie
    - Extract user's authentication credential ID

2. **Verify Password**
    - Verify `account_password` matches user's stored password
    - Use existing password verification logic from profile service

3. **Validate Certificate**
    - Check certificate exists and belongs to the event
    - Verify certificate is assigned to the authenticated user (`receiver_credential_id`)
    - Check certificate hasn't been claimed yet (`certificate_token_id` is null)
    - Verify user is a participant in the event

4. **Mint Certificate NFT**
    - Use system-managed wallet to mint certificate NFT
    - Call blockchain smart contract to mint token
    - Store transaction hash and token ID

5. **Update Database**
    - Update `event_certificates` table:
        - Set `certificate_token_id`
        - Set claimed timestamp
    - Update related inbox message to "claimed" status

6. **Return Response**
    - Return certificate details with token ID and transaction hash

### Database Schema Considerations

The `entity_event_certificate` table already has the necessary fields:

```sql
certificate_token_id VARCHAR(255) NULL  -- NULL = unclaimed, populated = claimed
```

### OpenAPI Documentation Example

```go
// @Summary Claim certificate as participant
// @Description Mint certificate NFT on blockchain and mark as claimed. Requires valid account password.
// @Tags Certificates, Participant
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param certificate_id path string true "Certificate ID"
// @Param body body ClaimCertificateRequest true "Account password for verification"
// @Success 200 {object} ClaimCertificateResponse "Certificate claimed successfully"
// @Failure 400 {object} customerror.ErrResponse "Invalid request or missing parameters"
// @Failure 401 {object} customerror.ErrResponse "Invalid password or unauthorized"
// @Failure 404 {object} customerror.ErrResponse "Certificate or event not found"
// @Failure 409 {object} customerror.ErrResponse "Certificate already claimed"
// @Failure 500 {object} customerror.ErrResponse "Internal server error"
// @Router /api/v1/events/{event_id}/certificates/{certificate_id}/claim [post]
// @Security BearerAuth
func (h *CertificateHandler) ClaimCertificate(ctx *fiber.Ctx) error {
    // Implementation here
}
```

## Testing

### Manual Testing Checklist

1. **Unclaimed Certificate Flow**
    - [ ] Navigate to certificate detail page
    - [ ] Verify "Claim Certificate" button appears
    - [ ] Click button and verify password modal opens
    - [ ] Enter correct PIN/password
    - [ ] Verify success message appears
    - [ ] Verify button changes to "Certificate Claimed" badge

2. **Already Claimed Certificate**
    - [ ] Navigate to claimed certificate detail page
    - [ ] Verify "Certificate Claimed" badge appears
    - [ ] Verify no claim button is shown

3. **Error Handling**
    - [ ] Try claiming with wrong password
    - [ ] Verify error message appears
    - [ ] Verify can retry claim
    - [ ] Test network error scenarios

4. **Loading States**
    - [ ] Verify button shows loading spinner during claim
    - [ ] Verify button is disabled during claim
    - [ ] Verify loading text appears

### Integration Testing

Once backend is implemented:

```typescript
// Test file: useClaimCertificate.test.ts
describe("useClaimCertificate", () => {
    it("should claim certificate successfully", async () => {
        // Test implementation
    });

    it("should handle invalid password", async () => {
        // Test implementation
    });

    it("should handle already claimed certificate", async () => {
        // Test implementation
    });

    it("should invalidate queries after successful claim", async () => {
        // Test implementation
    });
});
```

## UI/UX Considerations

### Design Decisions

1. **Button Placement**
    - Claim button appears below certificate details
    - Centered for emphasis
    - Full width on mobile for easy tapping

2. **Loading States**
    - Button shows spinner and "Claiming..." text
    - Button is disabled during claim
    - Clear visual feedback

3. **Success State**
    - Green badge with checkmark icon
    - Clear "Certificate Claimed" message
    - Matches existing design patterns

4. **Helper Text**
    - Small gray text below button
    - Explains action: "Sign this certificate to add it to your blockchain credentials"
    - Educates users about blockchain aspect

### Accessibility

- Button has proper ARIA labels
- Loading states announced to screen readers
- High contrast colors for status badges
- Keyboard navigation supported

## Future Enhancements

### Potential Improvements

1. **Download Certificate**
    - Add download button after claiming
    - Generate certificate image with QR code
    - Include blockchain verification link

2. **Share Certificate**
    - Social media sharing
    - Generate shareable link
    - Copy certificate details

3. **Transaction History**
    - Show blockchain transaction details
    - Link to blockchain explorer
    - Transaction confirmation status

4. **Batch Claiming**
    - Claim multiple certificates at once
    - Show progress for multiple claims
    - Aggregate success/failure reporting

5. **Notification Integration**
    - Push notification when certificate is ready to claim
    - Email notification with claim link
    - Reminder for unclaimed certificates

## Related Files

### Frontend Dependencies

- `/apps/web/src/hooks/usePassowordPrompt.ts` - Password prompt hook
- `/apps/web/src/components/ui/password-pin-modal.tsx` - Password modal UI
- `/apps/web/src/components/ui/headless-password-pin-modal.tsx` - Headless modal
- `/apps/web/src/components/providers/SignPasswordModal/store.ts` - Modal state
- `/apps/web/src/lib/queryKeys.ts` - Query key definitions

### Backend Files (To be implemented)

- `/apps/backend/core-api/internal/handler/certificate/claim_certificate.go`
- `/apps/backend/core-api/internal/usecase/certificate/claim_certificate.go`
- `/apps/backend/core-api/internal/repository/certificate/update_token_id.go`

## Migration from Mock to Real API

When backend is ready:

1. Update `useClaimCertificate.ts`:

    ```typescript
    // Replace mock with:
    const response = await coreApiClient.v1.claimCertificate(
        { certificateId, eventId },
        { account_password: accountPassword },
    );
    return response;
    ```

2. Generate API client:

    ```bash
    pnpm gen-api:core
    ```

3. Test with real backend endpoint

4. Update error handling for backend-specific errors

5. Update success handling to use real token IDs

## Conclusion

The certificate claiming mechanism is fully implemented on the frontend and ready for backend integration. The implementation follows existing patterns in the codebase, uses proper TypeScript typing, includes comprehensive error handling, and provides a smooth user experience consistent with the rest of the application.

The mock implementation allows frontend development and testing to continue while the backend API is being developed. Once the backend endpoint is ready, integration should be straightforward with minimal code changes required.
