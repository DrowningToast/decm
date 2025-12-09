# Certificate Claim - Debugging Guide

## 🐛 Current Issue

**Error**: `"failed to mint certificate NFT: execution reverted"`

**Location**: Smart contract transaction on blockchain

---

## 📊 Added Logging

I've added comprehensive logging throughout the certificate claiming process. Here's what to look for:

### 1. **Certificate Signature Data**

```
📝 Certificate signature data for contract verification
  - signed_message_digest: The digest that was signed
  - host_signature: The host's signature (hex string)
  - host_signature_length: Length of the signature string
  - sign_message: The original message that was signed
  - host_wallet_address: Expected host address
```

### 2. **Signature Decoding**

```
🔐 Decoded host signature
  - signature_bytes_length: Should be 65 bytes for ECDSA
  - signature_hex: The signature in hex format
```

### 3. **Signature Recovery Test**

```
✅ Signature recovery test
  - recovered_signer: Address recovered from the signature
  - expected_host_address: The host's wallet address
  - addresses_match: Whether they match (CRITICAL!)
```

**⚠️ CRITICAL CHECK**: If `addresses_match` is `false`, the contract will definitely revert because the recovered signer won't be recognized as a host/admin.

### 4. **Minting Parameters**

```
🚀 Attempting to mint certificate NFT
  - certificate_id: Certificate UUID
  - receiver_address: Participant's address
  - contract_address: Certificate contract address
  - transactor_address: System wallet address
  - user_id: User credential ID
  - issuer_id: Issuer credential ID
  - num_issuer_addresses: Number of issuers
  - num_issuer_proofs: Number of issuer signatures
  - certificate_title: Certificate title
  - certificate_subtitle: Certificate subtitle
  - data_hash: SHA256 hash of certificate data
```

### 5. **Transaction Status**

```
✅ Transaction submitted successfully
  - tx_hash: Transaction hash
```

```
⏳ Waiting for transaction to be mined...
```

```
📦 Transaction mined
  - tx_hash: Transaction hash
  - gas_used: Gas consumed
  - status: Transaction status (1 = success, 0 = reverted)
```

### 6. **Event Parsing**

```
🔍 Parsing transaction logs for CertificateMinted event
  - num_logs: Number of log entries
```

```
✅ Found CertificateMinted event
  - token_id: NFT token ID
  - log_index: Which log contained the event
```

---

## 🔍 Root Cause Analysis

Based on the smart contract code, the transaction reverts at this line:

```solidity
function mintNft(...) external nonReentrant {
    address signer = recoverSigner(signedMessageDigest, signature);
    requireHostOrAdmin(signer, msg.sender);  // <-- LIKELY FAILING HERE
    ...
}
```

The contract:

1. **Recovers the signer** from `signedMessageDigest` and `signature`
2. **Checks if the signer is a host or admin** in the EventAccessManager
3. **OR if msg.sender is an allowed sender**

### Possible Causes

#### 1. **Signature Mismatch** ❌

The `signedMessageDigest` and `hostSignatureBytes` don't correspond to each other.

**How to verify**:

- Check the logs for "Signature recovery test"
- If `addresses_match` is `false`, this is the problem
- The signature was not created by signing the digest

**Solution**:

- Verify that the certificate signature was created correctly when the issuer signed
- Check the `sign_certificate_requests` table to ensure signatures are valid

#### 2. **Host Not Registered in EventAccessManager** 🚫

The recovered signer address is not registered as a host in the EventAccessManager contract.

**How to verify**:

```bash
# Check if host is registered
cast call $EVENT_ACCESS_MANAGER_ADDRESS "checkIsHostOrAdmin(address)" $HOST_ADDRESS --rpc-url $RPC_URL
```

**Solution**:

- Register the host in the EventAccessManager contract
- Or add the system wallet address as an allowed message sender

#### 3. **System Wallet Not Allowed** 🔐

The `msg.sender` (system wallet) is not registered as an allowed sender.

**How to verify**:

```bash
# Check if system wallet is allowed
cast call $EVENT_ACCESS_MANAGER_ADDRESS "checkIsAllowedMsgSender(address)" $SYSTEM_WALLET_ADDRESS --rpc-url $RPC_URL
```

**Solution**:

- Add system wallet to allowed senders in EventAccessManager
- This is the recommended approach for backend minting

#### 4. **Wrong Certificate Contract Address** 📍

The certificate contract address in the database doesn't match the actual deployed contract.

**How to verify**:

- Check the `contract_address` in the minting logs
- Verify it matches the certificate contract from the event configuration

#### 5. **Signature Format Issue** 📝

The signature has an incorrect format (not 65 bytes, wrong encoding, etc.)

**How to verify**:

- Check `signature_bytes_length` in logs - should be 65
- Check if signature decoding succeeds

---

## 🛠️ Debugging Steps

### Step 1: Check the Logs

Run the claim operation and look for these critical log entries:

```bash
# Backend logs
tail -f backend.log | grep -E "📝|🔐|✅|❌|🚀"
```

Look for:

- ❌ **Red X marks** indicate errors
- ⚠️ **Warning signs** indicate potential issues
- ✅ **Check marks** with `addresses_match: false` = PROBLEM

### Step 2: Verify Signature Recovery

From the logs, copy:

- `signed_message_digest`
- `host_signature`
- `recovered_signer`
- `expected_host_address`

Check if `recovered_signer` == `expected_host_address`

**If NOT EQUAL** → The signature is invalid or doesn't match the digest

### Step 3: Check Smart Contract Permissions

```bash
# Set variables
EVENT_ACCESS_MANAGER_ADDRESS="<from contract deployment>"
HOST_ADDRESS="<from logs: expected_host_address>"
SYSTEM_WALLET_ADDRESS="<from logs: transactor_address>"
RPC_URL="<your RPC endpoint>"

# Check if host is registered
cast call $EVENT_ACCESS_MANAGER_ADDRESS \
  "checkIsHostOrAdmin(address)" \
  $HOST_ADDRESS \
  --rpc-url $RPC_URL

# Check if system wallet is allowed
cast call $EVENT_ACCESS_MANAGER_ADDRESS \
  "checkIsAllowedMsgSender(address)" \
  $SYSTEM_WALLET_ADDRESS \
  --rpc-url $RPC_URL
```

### Step 4: Test Signature Locally

Create a test script to verify signature recovery:

```go
package main

import (
    "fmt"
    "apps/backend/core-api/internal/usecase/cyptoutils"
)

func main() {
    signedMessageDigest := "<from logs>"
    hostSignature := "<from logs>"
    expectedAddress := "<from logs>"

    recovered, err := cyptoutils.GetAddressFromSignature(signedMessageDigest, hostSignature)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Recovered: %s\n", recovered.Hex())
    fmt.Printf("Expected: %s\n", expectedAddress)
    fmt.Printf("Match: %v\n", recovered.Hex() == expectedAddress)
}
```

### Step 5: Check Certificate Signature Creation

Look at the `event_certificate_signatures` table:

```sql
SELECT
  id,
  event_id,
  event_certificate_id,
  issuer_credential_id,
  host_signature,
  sign_message,
  sign_message_digest,
  is_signed,
  created_at
FROM event_certificate_signatures
WHERE event_certificate_id = '<certificate_id>';
```

Verify:

- `host_signature` is not null
- `sign_message_digest` is not null
- `is_signed` is 1 (true)

---

## 🎯 Quick Fixes

### Fix 1: Add System Wallet as Allowed Sender (RECOMMENDED)

This allows the backend to mint certificates on behalf of hosts:

```solidity
// On EventAccessManager contract
function addAllowedMsgSender(address sender) external onlyOwner {
    allowedMsgSenders[sender] = true;
}
```

Call this with the system wallet address from the logs.

### Fix 2: Use Host's Wallet to Sign

If the system wallet approach doesn't work, the backend needs to:

1. Get the host's private key (decrypted)
2. Create a fresh signature specifically for minting
3. Use that signature in the contract call

### Fix 3: Verify Certificate Signature Process

Check the certificate signing flow to ensure:

1. Host signs the certificate configuration correctly
2. Signature is stored properly in the database
3. Digest matches the signed message

---

## 📝 Log Example

Here's what successful logs should look like:

```
📝 Certificate signature data for contract verification
  signed_message_digest: "I want to sign certificate config abc-123 with deadline 12345678"
  host_signature: "0x1234abcd..."
  host_signature_length: 132
  sign_message: "I want to sign certificate config abc-123 with deadline 12345678"
  host_wallet_address: "0xHostAddress..."

🔐 Decoded host signature
  signature_bytes_length: 65
  signature_hex: "0x1234abcd..."

✅ Signature recovery test
  recovered_signer: "0xHostAddress..."
  expected_host_address: "0xHostAddress..."
  addresses_match: true  👈 MUST BE TRUE!

🚀 Attempting to mint certificate NFT
  certificate_id: "abc-123"
  receiver_address: "0xParticipantAddress..."
  contract_address: "0xCertificateContract..."
  transactor_address: "0xSystemWallet..."
  ...

✅ Transaction submitted successfully
  tx_hash: "0xTransactionHash..."

⏳ Waiting for transaction to be mined...

📦 Transaction mined
  tx_hash: "0xTransactionHash..."
  gas_used: 250000
  status: 1  👈 1 = SUCCESS!

🔍 Parsing transaction logs for CertificateMinted event
  num_logs: 3

✅ Found CertificateMinted event
  token_id: "123"
  log_index: 2
```

---

## 🚨 Error Indicators

Watch for these in the logs:

```
❌ Failed to decode host signature
❌ CRITICAL: Recovered signer does not match host wallet address!
❌ Failed to mint certificate NFT
❌ Transaction mining failed
❌ Transaction reverted on-chain
```

Any of these indicate the specific failure point.

---

## 🔄 Next Steps

1. **Run the claim operation** with the new logging
2. **Collect the logs** and analyze them
3. **Identify which check fails**:
    - Signature recovery?
    - Host permissions?
    - System wallet permissions?
4. **Apply the appropriate fix** based on the failure
5. **Test again** and verify success

---

## 💡 Tips

- **Enable debug logging** for maximum detail
- **Copy full log output** when reporting issues
- **Check blockchain explorer** for the transaction details
- **Verify contract deployment** is correct
- **Test with a known-good certificate** first

---

**Status**: Logging added, awaiting test results to diagnose root cause.
