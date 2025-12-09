# Certificate Claim - Comprehensive Logging Added

## ✅ Fixed

All compilation errors have been resolved. The code now compiles successfully.

---

## 🔍 Logging Added

### 1. **Certificate Signature Verification**

**Location**: Before sending to smart contract

**Logs**:

```
📝 Certificate signature data for contract verification
  - signed_message_digest: The message digest that was signed
  - host_signature: Host's signature (with length)
  - sign_message: Original sign message
  - host_wallet_address: Expected host address

🔐 Decoded host signature
  - signature_bytes_length: Should be 65 for valid ECDSA
  - signature_hex: Full signature in hex

✅ Signature recovery test
  - recovered_signer: Address recovered from signature
  - expected_host_address: Host's wallet address
  - addresses_match: TRUE/FALSE (CRITICAL!)
```

**Critical Check**: If `addresses_match` is FALSE, the smart contract will revert because the recovered signer won't be recognized as a host.

---

### 2. **Minting Parameters**

**Location**: Before calling smart contract

**Logs**:

```
🚀 Attempting to mint certificate NFT
  - certificate_id: Certificate UUID
  - receiver_address: Participant's address
  - contract_address: Certificate contract address
  - transactor_address: System wallet doing the minting
  - user_id: User credential ID
  - issuer_id: Issuer credential ID
  - num_issuer_addresses: Number of issuers
  - num_issuer_proofs: Number of issuer signatures
  - certificate_title: Certificate title
  - certificate_subtitle: Certificate subtitle
  - data_hash: SHA256 hash of certificate data

Debug level:
  - encrypted_user_data_length
  - backend_encrypted_user_data_length
  - user_encrypted_proof_length
  - backend_encrypted_proof_length
  - sign_message_digest
  - host_signature_length
```

---

### 3. **Transaction Submission**

**Logs**:

```
✅ Transaction submitted successfully
  - tx_hash: Blockchain transaction hash
```

---

### 4. **Transaction Mining**

**Logs**:

```
⏳ Waiting for transaction to be mined...

📦 Transaction mined
  - tx_hash: Transaction hash
  - gas_used: Gas consumed
  - block_number: Block number
  - status: 1 (success) or 0 (reverted)
```

**On Revert**:

```
❌ Transaction reverted on-chain
  - tx_hash: Transaction hash
  - gas_used: Gas consumed
  - block_number: Block number
```

---

### 5. **Event Parsing**

**Logs**:

```
🔍 Parsing transaction logs for CertificateMinted event
  - num_logs: Number of log entries

✅ Found CertificateMinted event
  - token_id: NFT token ID
  - log_index: Which log entry contained the event
```

**On Failure**:

```
❌ Failed to find CertificateMinted event in transaction logs
```

---

## 🐛 Debugging the "execution reverted" Error

### Next Steps

1. **Run the claim operation** to trigger the new logging
2. **Check the logs** for the signature recovery test
3. **Look for the critical indicator**:

    ```
    ✅ Signature recovery test
      addresses_match: false  👈 THIS IS THE PROBLEM
    ```

4. **If addresses don't match**:
    - The signature was not created by the host
    - OR the signature doesn't correspond to the digest
    - OR the certificate signature creation process has a bug

5. **If addresses DO match**:
    - Check if the host is registered in EventAccessManager
    - Check if the system wallet is an allowed message sender
    - Verify the certificate contract address is correct

---

## 🔑 Most Likely Issue

Based on the smart contract code:

```solidity
function mintNft(...) external nonReentrant {
    address signer = recoverSigner(signedMessageDigest, signature);
    requireHostOrAdmin(signer, msg.sender);  // <-- LIKELY FAILING
    ...
}
```

The contract checks:

1. **Recover signer** from signature
2. **Verify signer is host/admin** OR msg.sender is allowed

**If signer recovery fails** → addresses won't match in our logs
**If signer is not host/admin** → contract reverts with "Not host or admin or allowed msg sender"

---

## 📊 How to Use the Logs

### Viewing Logs

```bash
# Real-time monitoring
tail -f backend.log | grep -E "📝|🔐|✅|❌|🚀|⏳|📦|🔍"

# Search for specific certificate claim
grep "certificate_id.*<certificate-id>" backend.log

# Extract just the signature verification
grep -A 20 "Certificate signature data" backend.log
```

### Analyzing the Output

**Success Pattern**:

```
📝 Certificate signature data...
🔐 Decoded host signature... (65 bytes)
✅ Signature recovery test... (addresses_match: true)
🚀 Attempting to mint...
✅ Transaction submitted... (tx_hash: 0x...)
⏳ Waiting for transaction...
📦 Transaction mined... (status: 1)
🔍 Parsing transaction logs...
✅ Found CertificateMinted event... (token_id: 123)
```

**Failure Pattern**:

```
📝 Certificate signature data...
🔐 Decoded host signature...
⚠️ Could not recover signer... OR
✅ Signature recovery test... (addresses_match: false) OR
❌ CRITICAL: Recovered signer does not match...
🚀 Attempting to mint...
✅ Transaction submitted... (tx_hash: 0x...)
⏳ Waiting for transaction...
📦 Transaction mined... (status: 0)  👈 REVERTED!
❌ Transaction reverted on-chain
```

---

## 🎯 Quick Diagnosis

Run this after getting the error:

```bash
# Get the full log output
grep -B 5 -A 20 "failed to mint certificate NFT" backend.log > claim_error.log

# Check signature recovery
grep "addresses_match" claim_error.log

# Check transaction details
grep "tx_hash" claim_error.log
```

Then:

- If `addresses_match: false` → Fix certificate signature creation
- If `addresses_match: true` → Check smart contract permissions
- Copy transaction hash and check on block explorer for detailed revert reason

---

## 📁 Related Files

- **Main File**: `apps/backend/core-api/internal/usecase/event/claim_certificate.go`
- **Debugging Guide**: `CERTIFICATE_CLAIM_DEBUGGING_GUIDE.md`
- **Smart Contract**: `apps/contracts/src/contracts/event/EventCertificate.sol`

---

## 🚀 Ready for Testing

The logging is now complete and ready to help diagnose the issue. Run the claim operation and analyze the logs to identify the root cause.

**Expected Result**: Clear indication of whether the issue is:

1. Signature mismatch (most likely)
2. Permission issue in smart contract
3. Other contract-related problem

---

**Status**: ✅ Logging complete, ready for testing
