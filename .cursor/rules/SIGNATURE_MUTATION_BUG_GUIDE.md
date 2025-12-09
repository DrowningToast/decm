# Signature Mutation Bug - Developer Guide

## The Bug

`cyptoutils.VerifySignatureByDigest()` **modifies the signature in-place**, changing the recovery ID (v) from Ethereum format (27/28) to go-ethereum format (0/1). Smart contracts expect v=27/28, so sending a mutated signature will cause `Event__InvalidSignature()` or `Themis__InvalidSignature()` revert errors.

## The Fix

### ❌ BROKEN Pattern

```go
// Wallet extension flow (user provides signature)
signature := []byte{...} // v=27 or 28

// Verify signature
isValid, err := cyptoutils.VerifySignatureByDigest(address, hash, signature)
// ⚠️ signature[64] is now 0 or 1 (MUTATED!)

// Send to contract
tx, err := contract.Method(transactor, address, message, signature)
// ❌ Contract rejects: expects v=27/28 but got v=0/1
```

### ✅ FIXED Pattern - Make a Copy First

```go
// Wallet extension flow (user provides signature)
signature := []byte{...} // v=27 or 28

// CRITICAL: Make a copy before verification
signatureCopy := make([]byte, len(signature))
copy(signatureCopy, signature)

// Verify the COPY
isValid, err := cyptoutils.VerifySignatureByDigest(address, hash, signatureCopy)
// signatureCopy[64] is mutated, but original signature is intact

// Send original to contract
tx, err := contract.Method(transactor, address, message, signature)
// ✅ Contract accepts: v=27/28 as expected
```

### ✅ ALTERNATIVE - Skip Verification for PIN Flow

```go
// PIN flow (server creates signature)
privateKey, address, err := cyptoutils.DecryptPrivateKey(encrypted, password)
signature, err := cyptoutils.Sign(hash.Bytes(), privateKey)
// signature already has v=27/28

// NO verification needed - we just created it!
// Send directly to contract
tx, err := contract.Method(transactor, address, message, signature)
// ✅ Works perfectly
```

## Functions Reference

### 🚨 Dangerous (Mutates Signature)

```go
// MODIFIES signature[64] from 27/28 → 0/1
cyptoutils.VerifySignatureByDigest(address, messageHash, signature)

// Also modifies the signature for same reason
cyptoutils.RecoverPublicKeyFromSignature(messageHash, signature)
```

**Rule:** Always make a copy before calling these if you need the original signature later.

### ✅ Safe Functions

```go
// Creates signature with v=27/28 (correct for contracts)
cyptoutils.Sign(digest, privateKey)

// Validates message format without mutating
cyptoutils.ValidateSignMessage(signMessage, address, contractAddr, deadline)

// Hashing functions (don't touch signature)
cyptoutils.HashEthereumMessage(message)
cyptoutils.GetSignMessage(address, contractAddr, deadline)
```

## Fixed Files

### 1. `join_event.go` - Both Flows Fixed

#### PIN Flow (`JoinEventWithPin`)

- **Line 187-189**: Creates signature with `cyptoutils.Sign()`
- **Line 192**: Sends signature directly to `joinEvent()` → contract ✅
- **Removed lines 293-303**: Deleted redundant verification that was mutating signature

#### Wallet Extension Flow (`JoinEventWithSignature`)

- **Line 226-229**: Makes a copy before verification ✅
- **Line 232**: Verifies the copy (mutates copy, not original)
- **Line 247**: Sends original signature to `joinEvent()` → contract ✅

### 2. `claim_certificate.go` - Fixed Defensively

Certificate claiming technically didn't have this bug (participant signature not used in contract), but added copy for consistency and safety:

#### Wallet Extension Flow (`ClaimCertificateWithSignature`)

- **Line 221-222**: Makes a copy before verification ✅
- **Line 225**: Verifies the copy (mutates copy, not original)
- **Line 242**: Recovers public key from original signature ✅
- **Line 254**: Passes original signature to `claimCertificate()` (not used in contract anyway)

**Why it was already safe:**

- Participant signature: Only for authentication + public key recovery (never sent to contract)
- Host signature: Used for blockchain transaction (fetched from DB as bytes, never mutated)

**Why we added the copy anyway:**

- Defensive programming: Prevents future bugs if code changes
- Consistency: Same pattern as `JoinEventWithSignature`
- No performance impact: Single allocation, minimal overhead

## Decision Matrix

| Flow Type                  | Signature Source                   | Need Verification? | Solution                              |
| -------------------------- | ---------------------------------- | ------------------ | ------------------------------------- |
| **PIN**                    | Server creates (`cyptoutils.Sign`) | ❌ No              | Send directly to contract             |
| **Wallet Extension**       | User provides                      | ✅ Yes             | Make copy, verify copy, send original |
| **Host Signature from DB** | Database                           | ❌ No\*            | Send directly to contract             |

\* Host signatures are pre-verified when created, no need to re-verify before contract call.

## Common Pitfalls

### ❌ Pitfall 1: "I need to verify before sending to contract"

**Wrong thinking:** "I should verify the signature is valid before wasting gas on a failed transaction."

**Why it's wrong:** The contract will verify the signature anyway. If you verify on the backend and mutate it, the contract will ALWAYS reject it.

**Correct approach:**

- **PIN flow:** Skip backend verification (you created it, you trust it)
- **Wallet flow:** Make a copy, verify the copy

### ❌ Pitfall 2: "I'll just adjust v back after verification"

```go
// Verify (mutates v)
isValid, err := cyptoutils.VerifySignatureByDigest(addr, hash, sig)

// Try to fix it
if sig[64] == 0 || sig[64] == 1 {
    sig[64] += 27 // ❌ Fragile! What if VerifySignatureByDigest changes?
}
```

**Why it's wrong:** Relies on internal implementation details that might change.

**Correct approach:** Make a copy before verification.

### ❌ Pitfall 3: "The function is called Verify, it shouldn't mutate"

**Wrong assumption:** "Functions with read-like names shouldn't have side effects."

**Reality:** Go allows functions to modify slices passed by value (they contain a pointer to the underlying array).

**Lesson:** Always check if a function mutates its parameters, especially in crypto code.

## Testing Checklist

When implementing a new blockchain transaction flow:

- [ ] Identify where signatures are created
- [ ] Identify where signatures are verified
- [ ] Identify where signatures are sent to contracts
- [ ] Ensure no mutation happens between creation/verification and contract call
- [ ] Add debug logging for `signature[64]` value (should be 27 or 28 before contract call)
- [ ] Test with actual blockchain (local devnet) to catch revert errors early

## Debug Tips

Add this to your code before sending to contract:

```go
if len(signature) == 65 {
    println("Signature V value:", signature[64], "(should be 27 or 28)")
    if signature[64] != 27 && signature[64] != 28 {
        println("⚠️ WARNING: Signature may have been corrupted!")
    }
}
```

If you see v=0 or v=1, the signature was mutated somewhere.

## Future Development

When creating new blockchain transaction flows:

1. **Identify signature usage pattern:**
    - Server-side signing (PIN) → Skip verification, send directly
    - Client-side signing (wallet) → Copy before verification
    - Pre-signed from DB → Send directly

2. **Document signature flow:**

    ```go
    // SIGNATURE FLOW: User wallet → Verify (copy) → Original to contract
    // or
    // SIGNATURE FLOW: Server creates → Direct to contract (no verify)
    ```

3. **Add safeguards:**
    ```go
    // Assert v value before contract call
    if len(signature) == 65 && signature[64] != 27 && signature[64] != 28 {
        return errors.New("signature corruption detected")
    }
    ```

## References

- **Ethereum Yellow Paper**: ECDSA signature format (r, s, v where v ∈ {27, 28})
- **go-ethereum crypto package**: Uses v ∈ {0, 1} internally
- **Solidity ECDSA.recover**: Expects v ∈ {27, 28}

## Related Issues

- Event registration transaction reverts: `Event__InvalidSignature()`
- Certificate claiming transaction reverts: `Themis__InvalidSignature()`
- Signature verification passes but contract rejects

All traced back to signature mutation bug.
