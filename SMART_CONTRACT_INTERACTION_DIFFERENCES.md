# Smart Contract Interaction Differences: join_event vs claim_certificate

## Summary

| Aspect                     | join_event         | claim_certificate           |
| -------------------------- | ------------------ | --------------------------- |
| **Contract Function**      | `AddParticipant`   | `MintNft`                   |
| **Parameters**             | 3 params           | 19 params                   |
| **Access Control**         | ❌ None            | ✅ `requireHostOrAdmin`     |
| **Signature Replay Check** | ✅ In contract     | ✅ Pre-flight + In contract |
| **Reentrancy Protection**  | ❌ None            | ✅ `nonReentrant` modifier  |
| **Pre-flight Validation**  | Basic state checks | Extensive validation        |

---

## 1. Function Signature Comparison

### join_event: `AddParticipant`

```solidity
function addParticipant(
    address participantAddress,
    string memory signedMessageDigest,
    bytes memory signature
) external
```

**Go Call:**

```go
eventContractInstance.AddParticipant(
    transactor,
    *participantAddress,  // 1. Address to add
    signMessage,          // 2. Raw message string
    signature,            // 3. Participant's signature (bytes)
)
```

### claim_certificate: `MintNft`

```solidity
function mintNft(
    address receiverAddress,
    string memory userId,
    string memory certificateId,
    string memory issuerId,
    string memory encryptedUserData,
    string memory backendEncryptedUserData,
    address[] memory issuerAddresses,
    string memory signedMessageDigest,  // ← Line 84
    bytes memory signature,              // ← Line 85 (user's query)
    string memory hostSignature,
    string memory hostPublicKey,
    string memory signMessage,
    string memory userEncryptedProof,
    string memory backendEncryptedProof,
    string memory certificateTitle,
    string memory certificateSubtitle,
    string memory hash,
    CertificateVCStructs.IssuerProof[] memory issuerProofs
) external nonReentrant
```

**Go Call:**

```go
certificateContractInstance.MintNft(
    transactor,
    receiverAddress,           // 1. NFT receiver
    userId,                    // 2. User ID
    certificateId,             // 3. Certificate ID
    issuerId,                  // 4. Issuer ID
    encryptedUserData,         // 5. User data (encrypted)
    backendEncryptedUserData,  // 6. Backend data (encrypted)
    issuerAddresses,           // 7. Array of issuer addresses
    signMessageStr,            // 8. signedMessageDigest (original message)
    hostSignatureBytes,        // 9. signature (bytes) - HOST's signature
    hostSignatureStr,          // 10. hostSignature (string representation)
    hostPublicKey,             // 11. Host's public key
    signMessageStr,            // 12. signMessage (same as #8)
    userEncryptedProof,        // 13. User's encrypted proof
    backendEncryptedProof,     // 14. Backend's encrypted proof
    certificateTitle,          // 15. Certificate title
    certificateSubtitle,       // 16. Certificate subtitle
    userDataHashStr,           // 17. Hash of CSV data
    issuerProofs,              // 18. Array of issuer proofs
)
```

---

## 2. Critical Differences

### A. Signature Parameter (Line 85)

**join_event:**

- `signature` = **Participant's** signature (signed by the participant themselves)
- Used to verify participant authorized the join action
- No access control check on the signer

**claim_certificate:**

- `signature` = **Host's** signature (signed by the event host)
- Used to verify host authorized the certificate mint
- Access control checks: `requireHostOrAdmin(signer, msg.sender)`
- Signature must be from a registered host/admin OR transactor must be allowed sender

### B. Access Control

**join_event:**

```solidity
// NO access control check!
function addParticipant(...) external {
    address signer = recoverSigner(signedMessageDigest, signature);
    // No requireHostOrAdmin() call
    // Just validates participantAddress != 0 and seat count
}
```

**claim_certificate:**

```solidity
// HAS access control check!
function mintNft(...) external nonReentrant {
    address signer = recoverSigner(signedMessageDigest, signature);
    requireHostOrAdmin(signer, msg.sender);  // ← ACCESS CONTROL HERE
    // ...
}
```

### C. Pre-flight Validation

**join_event:**

- ✅ Checks signature validity (backend)
- ✅ Checks current/max seat count
- ✅ Checks if already joined
- ✅ Simple contract accessibility check (`GetParticipants`)

**claim_certificate:**

- ✅ Checks signature validity (backend)
- ✅ Checks signature replay (`UsedSignatures` pre-flight)
- ✅ Checks receiver address != zero
- ✅ Extensive error diagnostics on revert
- ❌ Access control pre-flight removed (was causing issues)

### D. Signature Replay Protection

**join_event:**

- ✅ Contract's `recoverSigner` marks signature as used
- ❌ No pre-flight check

**claim_certificate:**

- ✅ Contract's `recoverSigner` marks signature as used
- ✅ **Pre-flight check** via `UsedSignatures` mapping
- Fails fast if signature already used

### E. Reentrancy Protection

**join_event:**

- ❌ No `nonReentrant` modifier
- Simpler function, less attack surface

**claim_certificate:**

- ✅ `nonReentrant` modifier
- More complex function, multiple state changes

### F. Message Format

**join_event:**

```json
{
    "walletAddress": "0x...",
    "contractAddress": "0x...",
    "deadlineBlock": 12345
}
```

**claim_certificate:**

```json
{
  "eventContractAddress": "0x...",
  "receivers": ["0xhash1", "0xhash2", ...]
}
```

### G. Signature Creation Context

**join_event:**

- Signature created **on-demand** when user joins
- Participant signs their own authorization
- Signature passed directly to contract call

**claim_certificate:**

- Signature created **during import** (`import_certificate_receivers`)
- Host signs message containing all receiver hashes
- Signature **stored in database**, retrieved later during claim

---

## 3. Contract Function Flow Comparison

### join_event Flow:

1. User provides signature + message
2. Backend verifies signature locally
3. Check state (seat count, already joined)
4. Call `AddParticipant(participantAddress, signMessage, signature)`
5. Contract: `recoverSigner()` → validate → add participant
6. **NO access control enforcement**

### claim_certificate Flow:

1. User provides their signature + message (for authentication)
2. Backend retrieves **host's signature** from database
3. Backend verifies both signatures locally
4. Pre-flight: Check signature replay, zero address
5. Call `MintNft(...19 parameters...)`
6. Contract: `recoverSigner()` → **`requireHostOrAdmin()`** → mint NFT
7. **Access control enforced** (host/admin or allowed sender)

---

## 4. Error Handling

**join_event:**

```go
tx, err := eventContractInstance.AddParticipant(...)
if err != nil {
    return errors.Wrapf(err, "failed to add participant...")
}
// Simple error wrapping
```

**claim_certificate:**

```go
tx, err := certificateContractInstance.MintNft(...)
if err != nil {
    revertReason := extractRevertReasonFromError(err)
    // Extensive diagnostics:
    // - Possible error types
    // - Pre-flight check results
    // - Detected error patterns
    // - Full error message
}
// Comprehensive error analysis
```

---

## 5. Key Insight: Why claim_certificate Has More Security

1. **Access Control**: Only authorized hosts/admins can mint certificates
2. **Reentrancy Protection**: Prevents reentrancy attacks
3. **Signature Replay**: Pre-flight check prevents wasting gas
4. **Complex State**: Manages NFT state, certificate data, proofs
5. **Multi-party Verification**: Host signature + issuer proofs

**join_event is simpler** because:

- Anyone with valid signature can join
- No access control needed (self-authorization)
- Simpler state (just participant list)

---

## 6. The Critical Issue: Access Control

The main difference causing failures is:

**join_event works because:**

- No access control check → No dependency on EventAccessManager setup

**claim_certificate fails because:**

- Access control check requires:
    - Host registered in EventAccessManager, **OR**
    - System transactor registered in DecmAccessManager as allowed sender
- If neither condition met → transaction reverts

This is why the pre-flight access control check was reverting - the view function call to `CheckIsHostOrAdmin` was also failing, suggesting a contract setup or configuration issue.
