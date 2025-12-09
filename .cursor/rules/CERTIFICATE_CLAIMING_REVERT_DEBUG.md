# Certificate Claiming Transaction Revert - Debug Guide

## Overview

The `mintNft` function in `EventCertificate.sol` can revert for several reasons. This guide lists ALL possible causes and how to debug them.

## All Revert Scenarios

### 1. ❌ `Themis__InvalidSignature()`

**Contract Location:** `ThemisUtils.sol:27`
**Triggered When:** `ECDSA.recover()` returns `address(0)`

**Root Causes:**

- Signature has wrong v value (0 or 1 instead of 27 or 28)
- Signature was signed with different message than provided
- Signature bytes are corrupted
- Signature format is invalid (not 65 bytes)

**How to Debug:**
The updated code now checks this BEFORE sending to contract (line 743-756, 765-774).

Look for these log messages:

```
❌ CRITICAL: Invalid signature length!
❌ CRITICAL: Invalid signature V value!
❌ CRITICAL: Local signature recovery FAILED!
```

**Fix:**

- Re-import certificate receivers to generate fresh signatures
- Ensure import uses correct host private key
- Check database for signature corruption

---

### 2. ❌ `Themis__SignatureAlreadyUsed()`

**Contract Location:** `ThemisUtils.sol:31`
**Triggered When:** `usedSignatures[signature] == true`

**Root Causes:**

- Same signature was used in a previous successful transaction
- Certificate receivers were re-imported but old signature still in DB
- Retry of a previously successful claim

**How to Debug:**
Already checked at line 684:

```go
isSignatureUsed, err := certificateContractInstance.UsedSignatures(nil, hostSignatureBytes)
```

Look for this log:

```
❌ CRITICAL: Signature has already been used!
```

**Fix:**

- Re-import certificate receivers to generate NEW signatures
- Check if certificate was already claimed successfully in a previous transaction

---

### 3. ❌ `"Not host or admin or allowed msg sender"`

**Contract Location:** `EventCertificate.sol:58`
**Triggered When:** BOTH conditions are false:

- Recovered signer is NOT host/admin
- AND transactor (msg.sender) is NOT allowed sender

**Root Causes:**

- Signature recovers to wrong address (doesn't match host wallet)
- Host was removed from EventAccessManager
- System transactor not added to allowed senders
- Wrong private key used during import

**How to Debug:**
The updated code now checks this BEFORE sending (line 779-795).

Look for this log:

```
❌ CRITICAL: Local signature recovery returned WRONG ADDRESS!
```

And check access control separately:

```solidity
// Manual check if needed
bool isHost = eventAccessManager.checkIsHostOrAdmin(recoveredAddress)
bool isAllowedSender = eventAccessManager.checkIsAllowedMsgSender(transactor.from)
// One must be true
```

**Possible Causes Detail:**

```
1. Wrong private key used to sign during import
   → Host's credential.EncryptedPrivateKey doesn't match their wallet_address

2. Sign message was modified after signing
   → Database corruption or encoding issues

3. Signature was corrupted
   → Hex encoding/decoding issues

4. Host wallet address in database doesn't match signing key
   → Mismatched credentials
```

**Fix:**

- Verify host credential: ensure `GetAddressFromPrivateKey(decrypted) == credential.WalletAddress`
- Check EventAccessManager: ensure host is still registered
- Re-import with correct host credentials

---

### 4. ❌ `"ERC721: mint to the zero address"`

**Contract Location:** OpenZeppelin ERC721
**Triggered When:** `receiverAddress == 0x0000000000000000000000000000000000000000`

**How to Debug:**
Already checked at line 715-720:

```go
if receiverAddress == (common.Address{}) {
    return error
}
```

**Fix:**

- Ensure participant has valid wallet address
- Check certificate.ReceiverCredentialId exists and has valid wallet

---

### 5. ❌ `"ERC721InvalidReceiver"` / `"ERC721: transfer to non ERC721Receiver implementer"`

**Contract Location:** OpenZeppelin ERC721 `_safeMint`
**Triggered When:** Receiver is a smart contract that:

- Doesn't implement `onERC721Received()` function
- OR implements it but returns wrong magic value
- OR reverts in the function

**How to Debug:**
The updated code now checks this (line 723-730):

```go
code, err := client.CodeAt(ctx, receiverAddress, nil)
if len(code) > 0 {
    // Receiver is a contract
}
```

Look for this log:

```
⚠️ Receiver is a smart contract - must implement ERC721Receiver interface
```

**Fix:**

- If receiver is supposed to be EOA (regular wallet), something went wrong
- If receiver is supposed to be contract, implement ERC721Receiver interface
- Use `_mint()` instead of `_safeMint()` (requires contract change)

---

### 6. ❌ Out of Gas

**Triggered When:** Transaction runs out of gas

**Root Causes:**

- Parameters too large (very long strings)
- Many receivers in sign message (large JSON)
- Complex contract execution

**How to Debug:**
Look at gas used in error message. If it's very close to gas limit, this is the issue.

**Fix:**

- Increase gas limit for transaction
- Reduce parameter sizes if possible
- Batch fewer certificates per import

---

### 7. ❌ Reentrancy Guard

**Contract Location:** `EventCertificate.sol:95` has `nonReentrant` modifier

**Triggered When:** Attempting to call mintNft while a previous call is still executing

**Root Causes:**

- Malicious receiver contract trying to re-enter
- Should be impossible in normal operation

**Fix:**

- This is a security feature, no fix needed

---

## Message Format Issues

### The Sign Message Must Match EXACTLY

**During Import (creates signature):**

```go
signMessage := SignMessage{
    EventContractAddress: "0x...",
    Receivers: ["hash1", "hash2", ...] // ALL receivers in batch
}
signMessageJSON := json.Marshal(signMessage)
// Result: {"eventContractAddress":"0x...","receivers":["hash1","hash2",...]}

hash := HashEthereumMessage(signMessageJSON) // Adds Ethereum prefix
signature := Sign(hash, privateKey) // Signs the hash
```

**During Claim (verifies signature):**

```go
signMessageStr := *certificate.SignMessage // Must be EXACT same JSON
hostSignatureBytes := decodeHex(certificate.HostSignature)

// Send to contract
MintNft(..., signMessageStr, hostSignatureBytes, ...)
```

**Contract Verification:**

```solidity
bytes32 hash = MessageHashUtils.toEthSignedMessageHash(bytes(signMessageStr));
address signer = ECDSA.recover(hash, signature);
// Must match host address
```

### Common Format Mistakes

❌ **Wrong:** Using different JSON format
❌ **Wrong:** Using hash instead of original message
❌ **Wrong:** Different field order in JSON
❌ **Wrong:** Modified receivers array
✅ **Correct:** Exact same string stored during import

---

## Debugging Workflow

### Step 1: Run Certificate Claim

The updated code will now catch errors BEFORE sending transaction.

### Step 2: Check Logs

Look for these critical error messages:

```
❌ CRITICAL: Invalid signature length!
→ Signature is corrupted, not 65 bytes

❌ CRITICAL: Invalid signature V value!
→ Signature has v=0 or v=1, should be 27 or 28
→ Signature was mutated or created incorrectly

❌ CRITICAL: Local signature recovery FAILED!
→ Signature is invalid for the given message
→ Message or signature was modified

❌ CRITICAL: Local signature recovery returned WRONG ADDRESS!
→ Signature is valid but signed by wrong private key
→ Host wallet address doesn't match signing key
```

### Step 3: Verify Database State

```sql
-- Get certificate signature
SELECT
    ecs.id,
    ecs.host_signature,
    ecs.sign_message,
    LENGTH(ecs.host_signature) as sig_length,
    ec.receiver_credential_id,
    ac.wallet_address as receiver_wallet
FROM event_certificate_signatures ecs
JOIN event_certificates ec ON ecs.event_certificate_id = ec.id
JOIN authentication_credentials ac ON ec.receiver_credential_id = ac.id
WHERE ec.id = '<certificate_id>';

-- Verify host credentials
SELECT
    ac.id,
    ac.wallet_address,
    LENGTH(ac.encrypted_private_key) as pk_length
FROM authentication_credentials ac
JOIN events e ON e.owner_credential_id = ac.id
WHERE e.id = '<event_id>';
```

### Step 4: Verify Signature Manually

```go
// In a test or debug script
signMessage := "<value from database>"
hostSignature := "<hex from database>"
hostWalletAddress := "<from database>"

// Decode signature
sigBytes := hexutil.MustDecode(hostSignature)
fmt.Println("Signature length:", len(sigBytes))
fmt.Println("Signature V value:", sigBytes[64])

// Hash message
hash := cyptoutils.HashEthereumMessage(signMessage)

// Recover signer
pubKey, err := cyptoutils.RecoverPublicKeyFromSignature(hash, sigBytes)
recovered := cyptoutils.PublicKeyToAddress(pubKey)

fmt.Println("Recovered:", recovered.Hex())
fmt.Println("Expected:", hostWalletAddress)
fmt.Println("Match:", recovered.Hex() == hostWalletAddress)
```

---

## Quick Fix Checklist

If certificate claiming fails with "execution reverted":

- [ ] Check signature V value is 27 or 28 (now auto-checked)
- [ ] Verify signature recovers to correct host address (now auto-checked)
- [ ] Verify signature hasn't been used before (checked at line 684)
- [ ] Verify host is still registered in EventAccessManager
- [ ] Verify receiver address is not zero (checked at line 715)
- [ ] Check if receiver is contract requiring ERC721Receiver (now checked)
- [ ] Try re-importing certificate receivers to generate fresh signatures
- [ ] Verify host credential's encrypted private key matches their wallet address

---

## Prevention

To avoid these issues:

1. **During Import:**
    - Use correct host credentials
    - Verify host wallet address matches their private key
    - Don't modify signatures after creation
    - Store complete, unmodified sign message

2. **Database Integrity:**
    - Don't manually edit signatures
    - Don't manually edit sign messages
    - Keep certificates and signatures in sync

3. **Testing:**
    - Test certificate claiming immediately after import
    - Verify signatures before storing
    - Log all signature creation and verification steps

---

## If All Checks Pass But Still Fails

If local verification passes but contract still reverts:

1. **Check blockchain state:**
    - Is host still registered in EventAccessManager?
    - Is system transactor address allowed?
    - Was certificate already minted? (check tokenCounter)

2. **Check contract version:**
    - Is the deployed contract the expected version?
    - Was contract upgraded/redeployed?

3. **Check network:**
    - Correct RPC endpoint?
    - Correct chain ID?
    - Contract address correct?

4. **Last resort:**
    - Re-deploy certificate contract
    - Re-import all certificate receivers
    - Generate completely fresh signatures
