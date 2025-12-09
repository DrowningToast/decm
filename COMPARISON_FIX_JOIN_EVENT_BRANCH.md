# Comparison: Current Implementation vs fix/join-event Branch

## Critical Differences

### 1. **Participant Address Source** ⚠️ **CRITICAL**

**fix/join-event (Working):**

```go
// Uses JWT wallet address directly
participantAddr := common.HexToAddress(currentUser.WalletAddress)
tx, err := eventContractInstance.AddParticipant(transactor, participantAddr, signMessage, signature)
```

**Current Implementation:**

```go
// Gets address from decrypted private key
credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(...)
privateKey, participantAddress, err := cyptoutils.DecryptPrivateKey(...)
tx, err := eventContractInstance.AddParticipant(transactor, *participantAddress, signMessage, signature)
```

**Impact:** If `currentUser.WalletAddress` (from JWT) doesn't match the address derived from the private key, the signature verification will fail in the contract because:

- The message was signed with the private key (which has address `participantAddress`)
- But the contract receives `currentUser.WalletAddress` as the participant address
- The contract recovers the signer from the signature and compares it to the provided address
- **MISMATCH = REVERT!**

### 2. **Signing Method**

**fix/join-event (Working):**

```go
// Uses AuthUsecase helper
signature, _, err := uc.AuthUsecase.SecuredSignStringForManagedUser(ctx, currentUser, signMessage, password, &deadlineBlock)
```

**Current Implementation:**

```go
// Directly signs with private key
privateKey, participantAddress, err := cyptoutils.DecryptPrivateKey(...)
messageHash := cyptoutils.HashEthereumMessage(signMessage)
signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
```

**Impact:** Both should work, but the current implementation is more explicit about using the correct private key.

### 3. **joinEvent Function Signature**

**fix/join-event (Working):**

```go
func joinEvent(..., signature []byte, signMessage string, deadlineBlock *uint64)
```

**Current Implementation:**

```go
func joinEvent(..., signature []byte, signMessage string, participantAddress *common.Address)
```

**Impact:** The function signature changed - current version passes the participant address directly instead of deadlineBlock.

### 4. **Signature Verification Address**

**fix/join-event (Working):**

```go
// Verifies signature against JWT wallet address
isValidHash, err := cyptoutils.VerifySignatureByDigest(
    common.HexToAddress(currentUser.WalletAddress),
    messageHash,
    signature)
```

**Current Implementation:**

```go
// Verifies signature against address from private key
isValidHash, err := cyptoutils.VerifySignatureByDigest(
    *participantAddress,  // From decrypted private key
    messageHash,
    signature)
```

### 5. **Debug Logging**

**fix/join-event (Working):**

```go
// Uses println for simple debug output
println("=== DEBUG AddParticipant ===")
println("Participant Address:", participantAddr.Hex())
```

**Current Implementation:**

```go
// Uses structured slog logging with extensive pre-flight checks
slog.Info("🚀 Attempting to add participant to event", ...)
// Plus many pre-flight validation checks
```

### 6. **Error Handling**

**fix/join-event (Working):**

```go
// Simple error extraction with println
if strings.Contains(errStr, "execution reverted:") {
    parts := strings.SplitN(errStr, "execution reverted:", 2)
    revertReason := strings.TrimSpace(parts[1])
    println("🔴 REVERT REASON:", revertReason)
}
```

**Current Implementation:**

```go
// Comprehensive error extraction function
revertReason := extractRevertReasonFromError(err)
// Plus detailed logging with structured fields
```

### 7. **Pre-flight Checks**

**fix/join-event (Working):**

```go
// Basic contract accessibility check
participants, err := eventContractInstance.GetParticipants(callOpts)
if err != nil {
    println("⚠️  Warning: Cannot read contract state:", err.Error())
}
```

**Current Implementation:**

```go
// Extensive pre-flight checks:
// 1. Signature replay check
// 2. Signature validity (already verified in Go)
// 3. Seats availability
// 4. Participant already joined
// 5. Access control (transactor allowed)
```

## Root Cause Analysis

The **KEY DIFFERENCE** that could cause the revert:

### **Address Mismatch Issue**

In the working version:

- Message is signed with address from JWT: `currentUser.WalletAddress`
- Contract receives: `currentUser.WalletAddress`
- Signature is verified against: `currentUser.WalletAddress`
- ✅ **MATCH** → Works!

In the current version:

- Message is signed with address from private key: `participantAddress` (derived from decrypted private key)
- Contract receives: `*participantAddress` (from private key)
- Signature is verified against: `*participantAddress` (from private key)
- ✅ Should also match... BUT...

### **Potential Problem:**

If `currentUser.WalletAddress` (from JWT) **does not match** `participantAddress` (from private key):

- The signature is created with the private key's address
- But somewhere in the flow, `currentUser.WalletAddress` might still be used
- This could cause a mismatch

However, looking at the current code more carefully:

- Line 309: Verification uses `*participantAddress` ✓
- Line 459: Contract call uses `*participantAddress` ✓
- Line 279: Chain check uses `*participantAddress` ✓

So the current implementation should be correct IF the private key matches the user's actual wallet.

## Recommendation

**Check if the issue is:**

1. **Private key mismatch**: The encrypted private key in the database doesn't match the wallet address in JWT
2. **Signature replay**: The signature was already used (pre-flight check should catch this)
3. **Access control**: Transactor not allowed (pre-flight check shows `is_transactor_allowed=true` ✓)

Since your logs show:

```
is_transactor_allowed=true
```

The issue is likely:

- **Signature replay** (signature already used)
- **Address mismatch** between what's in JWT vs what's in the private key
- **Event is full** or **participant already joined**

The extensive pre-flight checks added should help identify which one!
