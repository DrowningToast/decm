# Certificate Mint Readiness - Final Implementation ✅

## Critical Requirement Correction

**Original Issue**: The initial implementation only checked if _at least one_ issuer had signed.

**Corrected**: The system now correctly validates that **ALL assigned issuers must have signed** before certificates can be minted.

## Key Changes Made

### 1. SQL Query Logic Enhancement

**New Query**: `AllIssuersHaveSigned`

```sql
SELECT
    CASE
        WHEN COUNT(*) = 0 THEN false  -- No issuers assigned
        WHEN COUNT(*) = COUNT(*) FILTER (WHERE is_signed = 1) THEN true  -- All signed
        ELSE false  -- Some unsigned
    END AS all_issuers_signed
FROM event_issuers
WHERE event_id = $1 AND deleted_at IS NULL;
```

This query:

- Returns `false` if no issuers are assigned
- Returns `true` ONLY if every issuer has `is_signed = 1`
- Returns `false` if even one issuer hasn't signed

### 2. Enhanced Response Structure

```typescript
{
  "is_ready": boolean,
  "has_certificate_config": boolean,
  "all_issuers_have_signed": boolean,  // ✅ Changed from has_signed_issuers
  "signed_issuers_count": number,
  "total_issuers_count": number,       // ✅ NEW - shows X/Y signed
  "has_certificate_contract": boolean,
  "certificate_contract_address": string | null,
  "missing_requirements": string[]
}
```

### 3. Detailed Error Messages

The system now provides specific feedback:

| Scenario                 | Message                                       |
| ------------------------ | --------------------------------------------- |
| No issuers assigned      | "No issuers have been assigned to this event" |
| Partial signing          | "Not all issuers have signed (1/3 signed)"    |
| No config                | "Certificate configuration is not set up"     |
| No contract              | "Event contracts are not deployed"            |
| Contract address missing | "Certificate contract address is not set"     |

## Complete SQL Queries Added

1. **AllIssuersHaveSigned** - Primary validation query
2. **GetTotalIssuersCount** - Count all assigned issuers
3. **GetSignedIssuersCount** - Count signed issuers (for reporting)
4. **HasSignedIssuers** - Quick check if any signed (helper)

## Implementation Files Updated

### Modified Files (8):

1. `packages/database/queries/event_issuers.sql` - Added 4 queries
2. `apps/backend/core-api/internal/datagateway/event/event_issuer.go` - Extended interface
3. `apps/backend/core-api/internal/repositories/postgres/event_issuer.go` - Implemented methods
4. `apps/backend/core-api/internal/usecase/eventconfig/event_config.go` - Added dependencies
5. `apps/backend/core-api/internal/usecase/eventconfig/check_certificate_mint_readiness.go` - Updated logic
6. `apps/backend/core-api/internal/handler/eventconfig/check_certificate_mint_readiness_response.go` - Updated response
7. `apps/backend/core-api/internal/handler/eventconfig/routes.go` - Added route
8. `apps/backend/core-api/cmd/main.go` - Updated DI

### Generated Files (4):

1. `packages/database/go/generated/event_issuers.sql.go` - sqlc generated
2. `apps/backend/core-api/docs/swagger.json` - OpenAPI spec
3. `apps/backend/core-api/docs/swagger.yaml` - OpenAPI spec
4. `apps/backend/core-api/docs/docs.go` - Go docs

## Validation Flow

```
┌─────────────────────────────────────────┐
│  GET /events/{id}/config/certificate/   │
│         mint-readiness                   │
└──────────────┬──────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│  1. Check Authorization                  │
│     - Verified Organizer OR              │
│     - Assigned Issuer                    │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│  2. Validate Certificate Config          │
│     ✓ Base template exists               │
│     ✓ Positions configured               │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│  3. Check ALL Issuers Signed ⚠️          │
│     Query: AllIssuersHaveSigned          │
│     - Count total issuers                │
│     - Count signed issuers               │
│     - Compare: signed == total           │
│     - Fail if total == 0                 │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│  4. Verify Contract Deployment           │
│     ✓ Contract record exists             │
│     ✓ Certificate address set            │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│  5. Return Comprehensive Status          │
│     - is_ready (boolean)                 │
│     - Individual checks                  │
│     - Issuer counts (X/Y)                │
│     - Missing requirements list          │
└──────────────────────────────────────────┘
```

## Example Responses

### ✅ Ready to Mint (All Requirements Met)

```json
{
    "is_ready": true,
    "has_certificate_config": true,
    "all_issuers_have_signed": true,
    "signed_issuers_count": 3,
    "total_issuers_count": 3,
    "has_certificate_contract": true,
    "certificate_contract_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7",
    "missing_requirements": []
}
```

### ❌ Not Ready (2 of 3 Issuers Signed)

```json
{
    "is_ready": false,
    "has_certificate_config": true,
    "all_issuers_have_signed": false,
    "signed_issuers_count": 2,
    "total_issuers_count": 3,
    "has_certificate_contract": true,
    "certificate_contract_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7",
    "missing_requirements": ["Not all issuers have signed (2/3 signed)"]
}
```

### ❌ Not Ready (No Issuers Assigned)

```json
{
    "is_ready": false,
    "has_certificate_config": true,
    "all_issuers_have_signed": false,
    "signed_issuers_count": 0,
    "total_issuers_count": 0,
    "has_certificate_contract": true,
    "certificate_contract_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7",
    "missing_requirements": ["No issuers have been assigned to this event"]
}
```

### ❌ Not Ready (Multiple Issues)

```json
{
    "is_ready": false,
    "has_certificate_config": false,
    "all_issuers_have_signed": false,
    "signed_issuers_count": 1,
    "total_issuers_count": 2,
    "has_certificate_contract": false,
    "certificate_contract_address": null,
    "missing_requirements": [
        "Certificate configuration is not set up",
        "Not all issuers have signed (1/2 signed)",
        "Event contracts are not deployed"
    ]
}
```

## Testing Checklist

### Unit Tests Scenarios:

- [ ] All issuers signed (3/3) → `is_ready: true`
- [ ] Partial signing (2/3) → `is_ready: false`
- [ ] No issuers assigned (0/0) → `is_ready: false`
- [ ] Single issuer signed (1/1) → `is_ready: true`
- [ ] Single issuer not signed (0/1) → `is_ready: false`
- [ ] No certificate config → `is_ready: false`
- [ ] No contract address → `is_ready: false`
- [ ] Soft-deleted issuers excluded from count
- [ ] Correct error messages for each scenario

### Integration Tests:

- [ ] Organizer can access any event
- [ ] Issuer can only access assigned events
- [ ] Non-issuer/non-organizer gets 403
- [ ] Invalid event ID returns 400
- [ ] Non-existent event returns 404

## Build Status

✅ **All checks passed:**

- Database generation successful
- Go compilation successful
- OpenAPI documentation generated
- No linter errors
- Type safety maintained

## API Endpoint

```
GET /api/v1/events/{event_id}/config/certificate/mint-readiness
Authorization: Bearer <jwt_token>

Response: 200 OK
Content-Type: application/json
```

## Usage in Frontend

```typescript
import { DefaultApi } from "@decm/api";

const api = new DefaultApi({
    basePath: "http://localhost:8080/api/v1",
});

// Check if certificate is ready to mint
const readiness = await api.checkCertificateMintReadiness(eventId);

if (readiness.is_ready) {
    // Enable "Mint Certificates" button
    console.log("✅ Ready to mint!");
} else {
    // Show missing requirements to user
    console.log("❌ Not ready:", readiness.missing_requirements);
    console.log(`Signed: ${readiness.signed_issuers_count}/${readiness.total_issuers_count}`);
}
```

## Security Considerations

1. **Authorization**: Only organizers or assigned issuers can check readiness
2. **Soft Deletes**: Deleted issuers are excluded from all counts
3. **Type Safety**: sqlc generates type-safe Go code
4. **Input Validation**: Event ID validated as valid UUID

## Future Enhancements

- [ ] Webhook notification when all issuers sign
- [ ] Batch readiness check for multiple events
- [ ] Estimated gas cost for minting
- [ ] Blockchain network health check
- [ ] Issuer reminder system for unsigned issuers

---

## Summary

✅ **Requirement Met**: The system now correctly enforces that **ALL assigned issuers must sign** before certificates can be minted.

The implementation provides comprehensive feedback about signing status with detailed counts (X/Y signed) and clear error messages to guide users through the minting preparation process.
