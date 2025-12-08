# Certificate Claiming - TODOs and Open Questions

## ✅ **COMPLETED**

### Implementation

- [x] Fixed UserData format from JSON to CSV
- [x] Implemented CSV format: `name,academic_institution,certificate_title,certificate_subtitle`
- [x] Implemented dual encryption (ECIES + AES-GCM)
- [x] Added hash computation with `hexutil.Encode` (consistent with import pattern)
- [x] Implemented PIN-based claiming flow
- [x] Implemented wallet extension flow with public key recovery
- [x] Added all contract parameters
- [x] Created comprehensive unit tests
- [x] Fixed hash encoding to use `hexutil.Encode()` with "0x" prefix
- [x] Added missing imports and fixed compilation errors
- [x] Tests compile successfully

---

## 📋 **TODOs IN CODE**

### 1. ⚠️ **CRITICAL: Smart Contract Integration** (Line 454-524)

**Location**: `claim_certificate.go:454-524`

**Status**: 🔴 **Blocked - Contract call commented out**

**What needs to be done**:

```go
// Currently commented out (lines 476-520):
certificateContractInstance, err := certificateContract.NewEventCertificate(...)
tx, err := certificateContractInstance.MintNft(...)
receipt, err := bind.WaitMined(ctx, client, tx)
```

**Tasks**:

1. Uncomment the contract call code
2. Test the minting transaction on testnet
3. Extract `tokenId` from `CertificateMinted` event
4. Update certificate record with `tokenId` and transaction hash
5. Handle blockchain errors (gas, revert, timeout)
6. Add retry logic for failed transactions

**Dependencies**:

- Smart contract must be deployed
- Backend must have access to system wallet private key
- Sufficient gas/ETH in system wallet

**Questions**:

- ❓ Should we use a system wallet or user's wallet for gas?
- ❓ Do we need transaction queue/nonce management?
- ❓ What's the gas limit strategy?
- ❓ Should we use EIP-1559 or legacy gas pricing?

---

### 2. ⚠️ **MEDIUM: VC Proof Structure** (Line 421-423)

**Location**: `claim_certificate.go:421-423`

**Status**: 🟡 **Placeholder - Empty JSON**

**Current**:

```go
userEncryptedProof := "{}"
backendEncryptedProof := "{}"
```

**What needs to be done**:
Build proper Verifiable Credential (VC) proof structure according to W3C VC Data Model.

**Expected Structure** (needs clarification):

```json
{
    "type": "EcdsaSecp256k1Signature2019",
    "created": "2024-12-08T12:00:00Z",
    "proofPurpose": "assertionMethod",
    "verificationMethod": "did:example:issuer#key-1",
    "jws": "eyJhbGc...",
    "cryptosuite": "ecdsa-2019"
}
```

**Questions**:

- ❓ **What specific VC proof format should be used?** (EcdsaSecp256k1Signature2019, JWS, other?)
- ❓ **Should proofs be encrypted before storing on-chain?**
- ❓ **Do we need DID (Decentralized Identifier) integration?**
- ❓ **What fields are absolutely required in the proof?**
- ❓ **Is this for compliance or functional verification?**

**Suggested approach**:

1. Define the exact VC proof schema
2. Create a `buildVCProof()` function
3. Sign the proof with issuer's private key
4. Encrypt both user and backend versions
5. Validate proof structure before minting

---

### 3. 🔵 **LOW: Certificate Password Validation** (Line 98-100)

**Location**: `claim_certificate.go:98-100`

**Status**: 🔵 **Optional feature**

**Current**:

```go
if params.CertificatePassword != nil {
    // TODO: Implement password validation logic if needed
    // This would require a certificate_password field in the database
    return true, nil
}
```

**What needs to be done**:
Add support for password-protected certificates (certificates that require a password to claim, not user's account password).

**Requirements**:

1. Add `certificate_password` field to `event_certificates` table
2. Hash passwords with bcrypt/argon2
3. Validate password during claiming
4. Consider use cases: private certificates, exclusive access

**Questions**:

- ❓ **Is this feature needed?** (Different from user account password)
- ❓ **When would certificates require passwords?**
- ❓ **Should passwords expire?**
- ❓ **Single password or unique per certificate?**

**Decision needed**: Is this a required feature or can we defer it?

---

## ❓ **OPEN QUESTIONS & UNCERTAINTIES**

### Security & Cryptography

1. **ECIES Key Management**
    - ❓ Should we cache public keys to avoid deriving from private key each time?
    - ❓ Do we need key rotation for ECIES encryption?
    - ❓ What happens if user loses their wallet private key? (Can't decrypt their data)

2. **Hash Verification**
    - ❓ Should we verify the computed hash matches the stored `certificate_digest` before minting?
    - ❓ Currently we compute hash but don't compare it to stored value - is this intentional?

3. **Signature Validation**
    - ❓ Should we validate ALL issuer signatures before minting, not just check count?
    - ❓ What if an issuer signature is invalid? Should we reject or retry?

### Smart Contract Integration

4. **Transaction Management**
    - ❓ Who pays for gas? System wallet or user's wallet?
    - ❓ Should we implement transaction queuing/batching?
    - ❓ How do we handle failed transactions? Retry? Refund?
    - ❓ What's the timeout for waiting for transaction confirmation?

5. **Contract Parameters**
    - ❓ The contract has 18 parameters - are all of them required? Can any be optional?
    - ❓ `signature` parameter (bytes) vs `hostSignature` (string) - why both formats?
    - ❓ Should `certificateTitle` and `certificateSubtitle` be duplicated (they're already in encrypted UserData)?

### Data Consistency

6. **CSV Format Edge Cases**
    - ❓ What if CSV data contains commas? Should we escape them?
    - ❓ What if data contains newlines or special characters?
    - ❓ Should we validate CSV data length/format before encryption?

7. **Hash Computation**
    - ❓ The import pattern uses `hexutil.Encode()` with "0x" prefix - does smart contract expect this?
    - ❓ Should hash be computed before or after trimming whitespace?

### Error Handling

8. **Recovery Scenarios**
    - ❓ What if minting succeeds but database update fails? How to recover?
    - ❓ Should we use database transactions to ensure atomicity?
    - ❓ What if user claims twice (race condition)?

9. **Validation**
    - ❓ Should we validate certificate hasn't been revoked again right before minting?
    - ❓ Should we check if certificate contract is still deployed/active?

### Testing

10. **Integration Tests**
    - ❓ Do we need integration tests with actual smart contract (testnet)?
    - ❓ Should we mock the blockchain or use ganache/hardhat?
    - ❓ How to test gas estimation and transaction failures?

---

## 🎯 **RECOMMENDED NEXT STEPS**

### Priority 1 (MUST DO)

1. **Decide on VC Proof format** → Get specification/requirements
2. **Implement smart contract calling** → Uncomment and test contract integration
3. **Add tokenId extraction** → Parse event logs for minted token ID
4. **Test end-to-end flow** → Full claiming flow on testnet

### Priority 2 (SHOULD DO)

5. **Add hash verification** → Compare computed hash with stored digest
6. **Implement transaction retry** → Handle failed blockchain transactions
7. **Add more unit tests** → Cover all error paths and edge cases
8. **Integration tests** → Test with real smart contract on testnet

### Priority 3 (NICE TO HAVE)

9. **Certificate password feature** → Only if required by product
10. **CSV escaping** → Handle special characters in certificate data
11. **Performance optimization** → Cache public keys, batch transactions

---

## 📊 **TESTING STATUS**

### Unit Tests Created

✅ **15 test cases covering**:

- `CheckClaimEligibility` - 7 test cases (all scenarios)
- `ClaimCertificateWithPin` - 3 test cases (error paths)
- `ClaimCertificateWithSignature` - 3 test cases (error paths)
- CSV format consistency - 1 test
- Hash encoding consistency - 1 test
- Dual encryption - 1 test
- Public key recovery - 1 test
- Nil pointer handling - 1 test
- Empty string handling - 1 test

### Tests Compile

✅ All tests compile successfully

### Tests Run

✅ Sample test runs successfully

### Coverage Gaps

⚠️ **Need more tests for**:

- Full successful claiming flow (requires mocking more dependencies)
- Smart contract interaction (requires testnet or mock)
- Transaction failure scenarios
- Race conditions
- Data validation edge cases

---

## 🔐 **SECURITY CONSIDERATIONS**

### Completed

✅ Dual encryption implemented
✅ PII encrypted with backend key
✅ User data encrypted with user's public key
✅ Hash computation for verification
✅ Signature verification for authentication

### Still Need to Address

⚠️ Input validation for CSV data
⚠️ SQL injection prevention (verify all queries use parameterized queries)
⚠️ Rate limiting for claiming attempts
⚠️ Audit logging for all certificate operations

---

## 📝 **DOCUMENTATION NEEDS**

1. **API Documentation**
    - Document both claiming flows (PIN vs Wallet Extension)
    - Document error codes and messages
    - Document rate limits and constraints

2. **Client-Side Documentation**
    - How to decrypt `encryptedUserData` with wallet private key
    - Example code for wallet extension integration
    - Error handling on client side

3. **Smart Contract Documentation**
    - Document all 18 parameters and their purpose
    - Document events emitted
    - Document gas costs

---

## 🎬 **CONCLUSION**

### What's Done ✅

- Core claiming logic implemented
- Dual encryption working
- CSV format correct
- Hash computation consistent
- Tests created and compiling
- Both flows (PIN + Wallet Extension) implemented

### What's Blocking 🔴

1. **Smart contract integration** - Need to uncomment and test
2. **VC Proof specification** - Need product/compliance requirements
3. **TokenId extraction** - Depends on contract integration

### What Needs Clarification ❓

See all questions marked with ❓ above - these require product/technical decisions

---

**Last Updated**: December 8, 2025  
**Status**: 🟡 **Implementation Complete, Integration Pending**  
**Next Milestone**: Smart Contract Integration + End-to-End Testing
