# Certificate UserData Encryption Fix - Summary

## Overview

Fixed the certificate claiming process to properly structure UserData (Raw Certificate Data) as **CSV format** and implement dual encryption strategy as per the requirements.

## Key Changes

### 1. UserData Structure (PII - Raw Certificate Data)

**Previous (Incorrect):**

- Only 2 fields: `name`, `academicInstitution`
- Missing critical fields
- Incorrectly formatted as JSON

**Fixed (Correct):**

```
name,academic_institution,certificate_title,certificate_subtitle
```

**Format**: CSV (Comma-Separated Values), NOT JSON  
**Example**: `John Doe,MIT,Certificate of Achievement,Web3 Development Course`

All 4 fields are now included:

- `name` (TEXT) - Participant's name
- `academic_institution` (TEXT) - Academic institution
- `certificate_title` (TEXT) - Certificate title
- `certificate_subtitle` (TEXT) - Certificate subtitle

**Important**: This matches the pattern used in `import_certificate_receivers.go`

### 2. Dual Encryption Strategy

#### Strategy Overview

Two different encryption methods for different purposes:

1. **User-Encrypted Data (`encryptedUserData`)**
    - **Algorithm**: ECIES (Elliptic Curve Integrated Encryption Scheme)
    - **Key**: User's wallet PUBLIC KEY (secp256k1)
    - **Purpose**: Only the certificate owner can decrypt with their private key on client-side
    - **Decryption**: Client-side using user's private wallet key

2. **Backend-Encrypted Data (`backendEncryptedUserData`)**
    - **Algorithm**: AES-256-GCM (Deterministic for searchability)
    - **Key**: Backend PII_ENCRYPTION_KEY from environment variables
    - **Purpose**: Backend administrative access, auditing, and data management
    - **Decryption**: Backend-only using PII encryption key

#### Why Dual Encryption?

- **Privacy**: Users control their own data via ECIES encryption
- **Compliance**: Backend can decrypt for legal/administrative requirements
- **Security**: User data protected even if backend is compromised
- **Flexibility**: Different access levels for different purposes

### 3. Implementation Details

#### New Files Created

**`apps/backend/core-api/internal/usecase/cyptoutils/ecies.go`**
New utility for ECIES encryption/decryption:

- `EncryptWithPublicKeyBytes(plaintext, publicKey)` - Encrypt CSV with ECDSA public key
- `DecryptWithPrivateKey(ciphertext, privateKey)` - Decrypt CSV with ECDSA private key
- Uses `github.com/ethereum/go-ethereum/crypto/ecies` package
- Returns base64-encoded ciphertext

**`apps/backend/core-api/internal/usecase/cyptoutils/utils.go` (additions)**
New helper functions for wallet extension flow:

- `RecoverPublicKeyFromSignature(messageHash, signature)` - Recover public key from signature
- `PublicKeyToAddress(publicKey)` - Convert public key to Ethereum address

#### Modified Files

**`apps/backend/core-api/internal/usecase/event/claim_certificate.go`**

1. **Updated `claimCertificate` signature**:
    - Added `participantPublicKey *ecdsa.PublicKey` parameter
    - Required for ECIES encryption

2. **Fixed UserData CSV structure** (changed from JSON to CSV):
    - Now includes all 4 fields
    - CSV format: `name,academic_institution,certificate_title,certificate_subtitle`
    - Matches pattern from `import_certificate_receivers.go`

3. **Added hash computation**:

    ```go
    // Compute hash of the CSV data (for blockchain verification)
    userDataHash := cyptoutils.HashMessage(userDataCSV)
    userDataHashStr := hex.EncodeToString(userDataHash)
    ```

4. **Implemented dual encryption with CSV data**:

    ```go
    // Encrypt CSV data with user's public key (ECIES)
    encryptedUserData, err := cyptoutils.EncryptWithPublicKeyBytes(userDataCSV, participantPublicKey)

    // Encrypt CSV data with backend PII key (AES-GCM)
    backendEncryptedUserData, err := pgmapper.EncryptPII(userDataCSV, uc.cfg.PIIEncryptionKey)
    ```

5. **Updated `ClaimCertificateWithPin`** (PIN/Password Flow):
    - Derives public key from decrypted private key
    - Passes public key to `claimCertificate()`

6. **Updated `ClaimCertificateWithSignature`** (Wallet Extension Flow):
    - Recovers public key from signature
    - Verifies recovered address matches participant
    - Backend encrypts CSV with recovered public key
    - User receives encrypted data without sending any PII

#### Added Imports

- `crypto/ecdsa` - ECDSA public/private key types
- `apps/backend/common/pgmapper` - PII encryption utilities

### 4. Data Flow

#### Flow 1: Claiming Certificate with PIN (Mobile/Web App)

1. User provides password/PIN
2. Backend decrypts user's private key
3. Derive wallet address from private key (cryptographically consistent)
4. Derive public key from private key
5. Create UserData **CSV** (4 fields): `name,academic_institution,certificate_title,certificate_subtitle`
6. Compute hash of CSV data for blockchain verification
7. Encrypt CSV with user's public key (ECIES) → `encryptedUserData`
8. Encrypt CSV with PII key (AES-GCM) → `backendEncryptedUserData`
9. Mint NFT with both encrypted versions + hash
10. Store on blockchain

#### Flow 2: Claiming Certificate with Wallet Extension (MetaMask, etc.)

1. User signs a message with wallet extension (no PII sent by user)
2. Backend verifies signature
3. Backend recovers public key from signature
4. Backend verifies recovered address matches user's wallet
5. Create UserData **CSV** (4 fields): `name,academic_institution,certificate_title,certificate_subtitle`
6. Compute hash of CSV data for blockchain verification
7. Backend encrypts CSV with recovered public key (ECIES) → `encryptedUserData`
8. Backend encrypts CSV with PII key (AES-GCM) → `backendEncryptedUserData`
9. Mint NFT with both encrypted versions + hash
10. Store on blockchain
11. **User can decrypt `encryptedUserData` later using their wallet's private key**

#### Smart Contract Parameters

```solidity
function mintNft(
    address receiverAddress,           // User's wallet address
    string userId,                     // Credential ID
    string certificateId,              // Certificate UUID
    string issuerId,                   // Issuer credential ID
    string encryptedUserData,          // ECIES encrypted (user can decrypt)
    string backendEncryptedUserData,   // AES-GCM encrypted (backend can decrypt)
    // ... other params
)
```

### 5. Security Considerations

#### ECIES (User Encryption)

✅ **Advantages:**

- Only user can decrypt with their private key
- Blockchain stores encrypted data
- Client-side decryption possible
- No backend involvement for user access

⚠️ **Limitations:**

- Requires public key availability
- Slightly larger ciphertext size
- Cannot search encrypted data

#### AES-GCM (Backend Encryption)

✅ **Advantages:**

- Fast encryption/decryption
- Deterministic (searchable if needed)
- Backend can decrypt for admin purposes
- Smaller ciphertext size

⚠️ **Limitations:**

- Backend key compromise exposes all data
- Requires secure key management
- Must rotate keys periodically

### 6. Testing

#### Build Verification

```bash
cd /Users/supratouchsuwatno/Desktop/decm/apps/backend
go build -o /dev/null ./core-api/internal/usecase/event
# Exit code: 0 ✅ (Success)
```

#### Manual Testing Required

- [ ] Test PIN-based certificate claiming
- [ ] Verify UserData structure in minted NFT
- [ ] Decrypt `encryptedUserData` on client-side
- [ ] Decrypt `backendEncryptedUserData` on backend
- [ ] Verify all 4 fields are populated correctly
- [ ] Test with empty/null fields

### 7. Database Schema

**No database changes required** - the encrypted data is stored on-chain, not in the database.

The certificate table already has these fields for generating the UserData:

- `name` TEXT (already encrypted in DB)
- `academic_institution` TEXT (already encrypted in DB)
- `certificate_title` TEXT (already encrypted in DB)
- `certificate_subtitle` TEXT (already encrypted in DB)

These fields are decrypted from the DB, then re-encrypted with the dual-encryption strategy before minting.

### 8. Breaking Changes

#### API Changes

⚠️ **Function Signature Updated**:

**Before**:

```go
ClaimCertificateWithSignature(ctx, client, currentUser, certificateId, eligibilityProof, signature, signMessage)
```

**After** (wallet extension flow now fully supported):

```go
ClaimCertificateWithSignature(ctx, client, currentUser, certificateId, eligibilityProof, signature, signMessage)
// Same signature - but now works by recovering public key from signature
```

#### Key Changes

✅ **Wallet extension flow now fully functional**:

- User signs message with wallet extension (MetaMask, etc.)
- Backend recovers public key from signature
- Backend encrypts CSV with recovered public key
- No PII sent by user over the wire
- User can decrypt later with their wallet's private key

### 9. Environment Variables

Required in `.env`:

```bash
PII_ENCRYPTION_KEY=your-256-bit-encryption-key-here
```

Generate secure key:

```bash
openssl rand -base64 32
```

### 10. Future Enhancements

1. **Public Key Storage**
    - Store public keys in `authentication_credentials` table
    - Enable signature-based claiming without password
    - Improve UX for repeat claims

2. **Key Rotation**
    - Implement PII encryption key rotation
    - Re-encrypt existing backend data with new keys
    - Maintain multiple keys for decryption during transition

3. **Client-Side Decryption Library**
    - TypeScript utility for ECIES decryption
    - Browser wallet integration
    - Display decrypted certificate data in UI

4. **Encrypted Search**
    - Implement searchable encryption for backend data
    - Use deterministic encryption for specific fields
    - Balance privacy and functionality

## Summary

### What Was Fixed

✅ **Changed format from JSON to CSV** for UserData  
✅ Corrected UserData structure to include all 4 required fields  
✅ Implemented ECIES encryption for user-decryptable data  
✅ Implemented AES-GCM encryption for backend-accessible data  
✅ Added hash computation for blockchain verification  
✅ Added proper key management and encryption utilities  
✅ Updated function signatures to pass public keys  
✅ **Implemented wallet extension flow with public key recovery**  
✅ Added comprehensive error handling  
✅ Code compiles successfully

### Both Flows Now Supported

✅ **PIN Flow** - User provides password, backend derives keys  
✅ **Wallet Extension Flow** - User signs message, backend recovers public key

### Next Steps

1. Test the certificate claiming flow end-to-end
2. Implement client-side ECIES decryption in frontend
3. Add public key storage for better UX
4. Monitor blockchain transaction costs
5. Document decryption process for frontend developers

---

**Date**: December 8, 2025  
**Status**: ✅ Implementation Complete, Ready for Testing  
**Files Changed**: 3 (1 new, 2 modified)
