# Certificate Receiver List Change Impact Analysis

## Overview

This document describes what happens when the certificate receiver list is changed via the Import Certificate Receivers functionality.

## Import Flow

### API Endpoint

**POST** `/api/v1/events/{event_id}/certificates/import`

### Handler

**Location:** `apps/backend/core-api/internal/handler/event/import_certificate_receivers.go`

```48:79:apps/backend/core-api/internal/handler/event/import_certificate_receivers.go
func (h Handler) ImportCertificateReceivers(ctx *fiber.Ctx) error {
	requestBody := ImportCertificateReceiversRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	// Convert request to usecase format
	requests := make([]event.ImportCertificateReceiversRequest, 0, len(requestBody.Receivers))
	for _, receiver := range requestBody.Receivers {
		requests = append(requests, event.ImportCertificateReceiversRequest{
			FirstName:           receiver.FirstName,           // Already *string
			LastName:            receiver.LastName,            // Already *string
			AcademicInstitution: receiver.AcademicInstitution, // Already *string
			CertificateTitle:    receiver.CertificateTitle,    // Already *string
			CertificateSubtitle: receiver.CertificateSubtitle, // Already *string
			HostPin:             requestBody.HostPin,
		})
	}

	response, err := h.EventUc.ImportCertificateReceivers(ctx.UserContext(), requestBody.EventID, requests, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}
```

---

## UseCase Implementation

**Location:** `apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go`

### Step-by-Step Process

#### 1. **Authorization Check**

```46:55:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
func (uc *EventUsecase) ImportCertificateReceivers(ctx context.Context, eventID uuid.UUID, requests []ImportCertificateReceiversRequest, currentUser *auth.JwtClaims) (*ImportCertificateReceiversResponse, error) {
	// 1. Check if current user is authorized
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	if !credential.IsVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified organizer"))
	}
```

#### 2. **Event & Contract Validation**

```57:75:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	// 2. Check if event exists
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
	}

	// 3. Get eventContract from eventContracts table using eventID
	eventContract, err := uc.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if eventContract == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event contract not found"))
	}
```

#### 3. **🔴 CRITICAL: Reset All Issuer Signing Status**

```77:81:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	// Reset all event issuers' signing status
	err = uc.EventIssuerDataGateway.ResetAllEventIssuersSigningStatus(ctx, eventID)
	if err != nil {
		return nil, err
	}
```

**Impact:** All issuer signatures are invalidated when receiver list changes. Issuers must re-sign certificates.

#### 4. **Deploy or Use Existing Certificate Contract**

```88:121:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	eventCertificateAddressStr := ""

	if eventContract.CertificateContractAddress == nil {
		// 4. Deploy event certificate contract
		client, err := cyptoutils.GetEthereumClient()
		if err != nil {
			return nil, err
		}

		auth, err := cyptoutils.GetKeyedTransactor()
		if err != nil {
			return nil, err
		}

		eventCertificateAddress, tx, _, err := eventCertificateContract.DeployEventCertificate(
			auth,
			client,
			common.HexToAddress(eventContract.EventContractAddress),
			common.HexToAddress(eventContract.EventContractAddress),
		)
		if err != nil {
			return nil, err
		}

		// Wait for transaction to be mined
		_, err = bind.WaitMined(ctx, client, tx)
		if err != nil {
			return nil, err
		}

		eventCertificateAddressStr = eventCertificateAddress.Hex()
	} else {
		eventCertificateAddressStr = *eventContract.CertificateContractAddress
	}
```

**Note:** If a certificate contract already exists, it's reused. The same contract is used for all certificate imports for this event.

#### 5. **Update Event Contract**

```123:133:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	// 5. Update eventContract.certificate_contract_address
	_, err = uc.EventContractDataGateway.UpdateEventContract(ctx, generated.UpdateEventContractParams{
		EventID:                      eventID,
		AccessManagerContractAddress: eventContract.AccessManagerContractAddress,
		EventContractAddress:         eventContract.EventContractAddress,
		TicketContractAddress:        pgmapper.StringPtrToPgText(eventContract.TicketContractAddress),
		CertificateContractAddress:   pgmapper.StringPtrToPgText(&eventCertificateAddressStr),
	})
	if err != nil {
		return nil, err
	}
```

#### 6. **🔴 CRITICAL: Create New Certificate Records**

```135:192:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	// 6. Save certificate data to event_certificates
	certificates := make([]*entity.EventCertificate, 0, len(requests))
	var certificateIDs []uuid.UUID

	for _, req := range requests {
		// Safely dereference pointer fields (use empty string if nil)
		firstName := ""
		if req.FirstName != nil {
			firstName = *req.FirstName
		}
		lastName := ""
		if req.LastName != nil {
			lastName = *req.LastName
		}
		academicInstitution := ""
		if req.AcademicInstitution != nil {
			academicInstitution = *req.AcademicInstitution
		}
		certificateTitle := ""
		if req.CertificateTitle != nil {
			certificateTitle = *req.CertificateTitle
		}
		certificateSubtitle := ""
		if req.CertificateSubtitle != nil {
			certificateSubtitle = *req.CertificateSubtitle
		}

		// Combine first and last name
		name := fmt.Sprintf("%s %s", firstName, lastName)

		// Create CSV value for hash (used for blockchain verification)
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, academicInstitution, certificateTitle, certificateSubtitle)

		// Hash the CSV value
		hash := cyptoutils.HashMessage(csvValue)
		encodedHash := hexutil.Encode(hash)

		// Create certificate - pass original pointers to preserve nil vs empty distinction
		certificate, err := uc.EventCertificateDataGateway.CreateEventCertificate(ctx, eventdatagateway.CreateEventCertificateParameters{
			EventID:                 eventID,
			ReceiverCredentialID:    nil, // Will be set when receiver claims certificate
			ReceiverEmail:           nil, // Will be set when receiver claims certificate
			Name:                    stringPtrIfNotEmpty(name),
			AcademicInstitution:     req.AcademicInstitution,
			CertificateTitle:        req.CertificateTitle,
			CertificateSubtitle:     req.CertificateSubtitle,
			EventContractAddress:    eventContract.EventContractAddress,
			EventCertificateAddress: &eventCertificateAddressStr,
			CertificateTokenID:      nil, // Will be set when minted,
			Digest:                  &encodedHash,
		})
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, certificate)
		certificateIDs = append(certificateIDs, certificate.Id)
	}
```

**⚠️ IMPORTANT:** New certificates are **ADDED** to the database. **Existing certificates are NOT deleted**.

#### 7. **Create Sign Messages**

```194:234:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	// 7. Create sign_message with hashes
	receivers := make([]string, 0, len(requests))
	for _, req := range requests {
		// Safely dereference pointer fields (use empty string if nil)
		firstName := ""
		if req.FirstName != nil {
			firstName = *req.FirstName
		}
		lastName := ""
		if req.LastName != nil {
			lastName = *req.LastName
		}
		academicInstitution := ""
		if req.AcademicInstitution != nil {
			academicInstitution = *req.AcademicInstitution
		}
		certificateTitle := ""
		if req.CertificateTitle != nil {
			certificateTitle = *req.CertificateTitle
		}
		certificateSubtitle := ""
		if req.CertificateSubtitle != nil {
			certificateSubtitle = *req.CertificateSubtitle
		}

		// Combine first and last name
		name := fmt.Sprintf("%s %s", firstName, lastName)

		// Create CSV value
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, academicInstitution, certificateTitle, certificateSubtitle)

		// Hash the CSV value
		hash := cyptoutils.HashMessage(csvValue)
		encodedHash := hexutil.Encode(hash)
		receivers = append(receivers, encodedHash)
	}

	signMessage := SignMessage{
		EventContractAddress: eventCertificateAddressStr,
		Receivers:            receivers,
	}
```

#### 8. **Host Signs the Message**

```243:248:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	// 8. Host signs the message
	signMessageDigest := cyptoutils.HashMessage(signMessageJSON)
	signature, err := cyptoutils.Sign(signMessageDigest[:], privateKey)
	if err != nil {
		return nil, err
	}
```

#### 9. **Create Certificate Signatures for All Issuers**

```250:275:apps/backend/core-api/internal/usecase/event/import_certificate_receivers.go
	// 9. Create event_certificate_signatures for each certificate
	for _, certificateID := range certificateIDs {
		// Get event issuers for this event
		eventIssuers, err := uc.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
		if err != nil {
			return nil, err
		}

		// Create signature for each issuer
		for _, issuer := range eventIssuers {
			encodedSignMessageDigestStr := hexutil.Encode(signMessageDigest[:])
			encodedHostSignature := hexutil.Encode(signature)

			_, err := uc.EventCertificateSignatureDataGateway.CreateEventCertificateSignature(ctx, eventdatagateway.CreateEventCertificateSignatureParameters{
				EventCertificateID: certificateID,
				IssuerCredentialID: issuer.IssuerCredentialID,
				IssuerSignature:    nil, // Will be set when issuer signs
				HostSignature:      encodedHostSignature,
				SignMessage:        &signMessageJSON,
				SignMessageDigest:  &encodedSignMessageDigestStr,
			})
			if err != nil {
				return nil, err
			}
		}
	}
```

**Impact:** For each new certificate, signature records are created for ALL issuers (pending their signatures).

---

## Impact Summary

### ✅ What Happens

1. **New Certificates Are Created**
    - New certificate records are inserted into `event_certificates` table
    - Each certificate gets a unique digest hash
    - Certificates are created with `receiver_credential_id = NULL` (unclaimed)

2. **Issuer Signing Status Reset**
    - **ALL issuers' `is_signed` status is reset to 0**
    - Existing issuer signatures are invalidated
    - Issuers must re-sign certificates for the new receiver list

3. **Certificate Signatures Created**
    - For each new certificate, signature records are created for ALL issuers
    - Host signature is included in each record
    - Issuer signatures remain NULL until issuers sign

4. **Smart Contract Deployment**
    - If no certificate contract exists, a new one is deployed
    - If a contract already exists, it's reused for all certificates

5. **Frontend Cache Invalidation**
    ```typescript
    queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
    ```

### ⚠️ What Does NOT Happen

1. **Existing Certificates Are NOT Deleted**
    - Old certificate records remain in the database
    - This can lead to duplicate certificates for the same receiver
    - No cleanup or deduplication is performed

2. **Old Certificate Signatures Are NOT Deleted**
    - Old signature records remain in `event_certificate_signatures`
    - These become orphaned if old certificates are removed manually

3. **No Validation for Duplicates**
    - The system doesn't check if a receiver already has a certificate
    - Multiple certificates can be created for the same person

4. **No Transaction Rollback**
    - If the import fails partway, some certificates may already be created
    - No automatic cleanup of partial imports

---

## Potential Issues

### 🔴 Issue 1: Certificate Accumulation

**Problem:** Each import creates new certificates without deleting old ones.

**Example Scenario:**

- Import List A: [Alice, Bob]
- Import List B: [Alice, Charlie]
- Result: Database has 4 certificates (2 for Alice, 1 for Bob, 1 for Charlie)

**Impact:**

- Duplicate certificates for the same receiver
- Confusion about which certificate is valid
- Database bloat over time

**Recommendation:**

- Delete existing certificates before importing new ones, OR
- Check for existing certificates and update them instead of creating duplicates, OR
- Add a flag to indicate which certificate list is "active"

### 🔴 Issue 2: Orphaned Signatures

**Problem:** When old certificates remain, their signatures also remain.

**Impact:**

- Database clutter
- Potential confusion when querying signatures
- No clear way to know which signatures are valid

**Recommendation:**

- Delete signatures when certificates are deleted (cascade delete)
- Or delete old signatures when importing new receiver list

### 🔴 Issue 3: Issuer Re-signing Required

**Problem:** All issuer signatures are reset when receiver list changes, even for certificates that don't change.

**Example Scenario:**

- Issuers sign certificates for List A
- Host imports List B (which includes all receivers from List A plus new ones)
- All issuers must re-sign everything, even unchanged certificates

**Impact:**

- Unnecessary re-signing work
- Disruption to certificate issuance workflow

**Recommendation:**

- Only reset signatures for certificates that are new or changed
- Or provide option to preserve existing signatures

### 🔴 Issue 4: No Partial Import Handling

**Problem:** If import fails midway, some certificates are already created.

**Impact:**

- Inconsistent state
- Manual cleanup required

**Recommendation:**

- Use database transactions
- Implement rollback mechanism

### 🟡 Issue 5: Contract Reuse

**Problem:** Same smart contract is reused for all certificate imports.

**Impact:**

- All certificates share the same contract address
- Cannot distinguish between certificate "batches"
- May cause issues if contract needs to be reset

**Recommendation:**

- Consider deploying new contracts per import, OR
- Add versioning/tracking for certificate batches

---

## Frontend Impact

**Location:** `apps/web/src/hooks/events/useImportCertificates.ts`

```26:30:apps/web/src/hooks/events/useImportCertificates.ts
        onSuccess: () => {
            toast.success("Certificates imported successfully");
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
            navigate("/host/events/:eventId", { params: { eventId } });
        },
```

**Cache Invalidation:**

- All event queries are invalidated
- User is redirected to event details page
- Certificate list will show all certificates (old + new)

---

## Database Schema

### Tables Affected

1. **`event_certificates`**
    - New records created for each receiver
    - Old records remain untouched

2. **`event_certificate_signatures`**
    - New signature records created for each (certificate, issuer) pair
    - Old signature records remain

3. **`event_issuers`**
    - `is_signed` field reset to 0 for all issuers

4. **`event_contracts`**
    - `certificate_contract_address` updated (if new contract deployed)

---

## Recommendations

### Immediate Actions

1. **Add Certificate Cleanup**
    - Delete existing certificates before importing new ones
    - Or implement deduplication logic

2. **Cascade Delete Signatures**
    - Delete certificate signatures when certificates are deleted
    - Or clean up old signatures during import

3. **Add Validation**
    - Check for duplicate receivers before creating certificates
    - Warn users if duplicates are detected

4. **Transaction Safety**
    - Wrap import in database transaction
    - Implement rollback on failure

### Long-term Improvements

1. **Certificate Versioning**
    - Track certificate "batches" or "versions"
    - Allow users to see certificate history

2. **Incremental Updates**
    - Support adding/removing individual receivers
    - Preserve existing certificates that haven't changed

3. **Smart Signature Management**
    - Only reset signatures for certificates that changed
    - Preserve signatures for unchanged receivers

4. **User Warnings**
    - Warn users that existing certificates will be affected
    - Require confirmation before importing new list

---

## Testing Scenarios

1. **First Import**
    - Verify certificates are created
    - Verify signatures are created for all issuers
    - Verify contract is deployed

2. **Second Import (No Deletes)**
    - Verify old certificates remain
    - Verify new certificates are added
    - Verify all issuer signatures are reset

3. **Duplicate Receivers**
    - Verify duplicate certificates are created
    - Verify system handles duplicates gracefully

4. **Partial Failure**
    - Simulate failure during import
    - Verify system state is consistent

5. **Issuer Re-signing**
    - Verify issuers must re-sign after import
    - Verify signature status is reset correctly
