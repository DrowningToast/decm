# Signature Parameter Comparison: claim_certificate vs update_event

## Problem

`claim_certificate.go` transaction reverts with signature verification failure, while `update_event.go` works correctly.

## Key Differences

### update_event.go (WORKING) ✅

```go
// Line 160: Create fresh sign message
signMessage, err := cyptoutils.GetSignMessage(*hostAddress, eventContractAddress, calculatedDeadlineBlock)

// Line 165: Hash the message
messageHash := cyptoutils.HashEthereumMessage(signMessage)

// Line 166: Sign the hash
signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)

// Line 177-178: Pass to contract
instance.UpdateEvent(
    ...
    signMessage,  // ✅ Original message string
    signature,    // ✅ Signature bytes
)
```

**Flow:**

1. Creates message fresh
2. Hashes message
3. Signs hash
4. Passes original message + signature to contract
5. Contract hashes message again and recovers signer

### claim_certificate.go (BROKEN) ❌

```go
// Line 471-473: Retrieve from database
signedMessageDigest := *firstSignature.SignMessageDigest  // Hash (hex-encoded)
hostSignatureStr := firstSignature.HostSignature          // Signature (hex-encoded)
signMessageStr := *firstSignature.SignMessage             // Original message

// Line 484: Decode signature
hostSignatureBytes, err := hex.DecodeString(strings.TrimPrefix(hostSignatureStr, "0x"))

// Line 815-816: Pass to contract
certificateContractInstance.MintNft(
    ...
    signMessageStr,      // ✅ Original message string (should be correct)
    hostSignatureBytes,  // ✅ Signature bytes (should be correct)
    ...
)
```

**Flow:**

1. Retrieves stored message and signature from database
2. Decodes signature from hex
3. Passes original message + signature to contract
4. Contract hashes message again and recovers signer

## Contract Expectation

From `EventCertificate.sol:96`:

```solidity
address signer = recoverSigner(signedMessageDigest, signature);
```

From `ThemisUtils.sol:19-24`:

```solidity
function recoverSigner(string memory signedMessageDigest, bytes memory signature) public returns (address) {
    // Hash the message with the Ethereum signed message prefix
    bytes32 ethSignedMessageHash = MessageHashUtils.toEthSignedMessageHash(bytes(signedMessageDigest));

    // Recover the signer address
    address signer = ECDSA.recover(ethSignedMessageHash, signature);
    ...
}
```

**Contract expects:** Original message string (despite parameter name `signedMessageDigest`)

## How Signature Was Created

From `import_certificate_receivers.go:288-289`:

```go
signMessageDigest := cyptoutils.HashEthereumMessage(signMessageJSON)  // Hash the message
signature, err := cyptoutils.Sign(signMessageDigest.Bytes(), privateKey)  // Sign the hash
```

**Stored in database:**

- `SignMessage`: `signMessageJSON` (original message string)
- `SignMessageDigest`: hex-encoded hash
- `HostSignature`: hex-encoded signature bytes

## Analysis

Both flows should work identically:

1. Signature created by signing hash of original message
2. Contract receives original message and signature
3. Contract hashes message and recovers signer

**Possible issues:**

1. `signMessageStr` from database doesn't exactly match `signMessageJSON` used during signing
2. String encoding/whitespace differences
3. Message format changed between import and claim
4. Database corruption or modification

## Solution

The code appears correct, but we should:

1. Add validation to ensure `signMessageStr` matches expected format
2. Verify the signature locally before sending to contract (already done at line 759-795)
3. Consider regenerating the signature if verification fails locally

The local verification at line 759-795 should catch this issue before sending to contract. If local verification passes but contract fails, there may be a contract-side issue or the message format differs slightly.
