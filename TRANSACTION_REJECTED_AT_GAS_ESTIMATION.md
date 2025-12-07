# Transaction Rejected at Gas Estimation - Debugging Guide

## What's Happening

Your transaction is being **rejected BEFORE it's sent to the blockchain**. This happens at line 288 in `join_event.go`:

```go
tx, err := eventContractInstance.AddParticipant(transactor, participantAddr, messageHash.String(), signature)
```

## Why This Happens

When you call a contract function with go-ethereum, it automatically:

1. **Simulates the transaction** (eth_call)
2. **Estimates gas required** (eth_estimateGas)
3. **Only then submits** the transaction

If step 1 or 2 fails (transaction would revert), go-ethereum returns an error immediately without submitting the transaction. This is actually GOOD - it saves you gas fees!

## Enhanced Debug Output

Now when you run your code, you'll see detailed output like:

```
=== DEBUG AddParticipant ===
Participant Address: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
Contract Address: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Transactor Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Current Seats: 0
Max Seats: 100
========================

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ERROR: AddParticipant transaction REJECTED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Full Error: execution reverted: custom error 0x12345678

🔴 REVERT REASON: custom error 0x12345678

Possible causes:
  ❌ Smart contract rejected the call during gas estimation
     This means one of the contract's require/revert statements failed:
     - Event__SeatsCountReached: Event is full
     - Event__ParticipantIsAlreadyJoined: Already registered
     - Event__AddressCannotBeZero: Invalid address
     - Access control: Transactor doesn't have HOST/ADMIN role
     - Invalid signature: Signature verification failed

Transaction Details:
  Participant: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
  Contract: 0x5FbDB2315678afecb367f032d93F642f64180aa3
  Transactor: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
  Message Hash: 0xabcd1234...
  Signature: a1b2c3d4...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Common Causes & Solutions

### 1. Access Control Issue (MOST LIKELY)

**Problem**: The transactor wallet doesn't have HOST or ADMIN role.

Looking at your smart contract (`Event.sol` line 131-132):

```solidity
address signer = recoverSigner(signedMessageDigest, signature);
EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer);
```

The contract recovers the signer from the signature and checks if that signer has HOST/ADMIN role. **NOT the msg.sender!**

**Check:**

```bash
# In your blockchain console/script, verify:
# 1. What address does the signature recover to?
# 2. Does that address have HOST/ADMIN role?
```

**Solution:**

```solidity
// Grant HOST role to the signer address
eventAccessManager.grantHostRole(signerAddress, adminAddress);

// OR grant ADMIN role
eventAccessManager.grantAdminRole(signerAddress, adminAddress);
```

**Key Insight**: The signature must be signed by a wallet that has HOST/ADMIN role, not just any wallet.

### 2. Signature Format Issue

**Problem**: The message hash format doesn't match what the contract expects.

Your code passes (line 288):

```go
messageHash.String()  // This converts bytes32 to hex string
```

The contract expects (line 128):

```solidity
string memory signedMessageDigest
```

**Check**: Does `messageHash.String()` produce the format your contract's `recoverSigner` function expects?

**Debug**:

```go
// Add this before line 288:
println("Message Hash Type:", reflect.TypeOf(messageHash))
println("Message Hash String:", messageHash.String())
println("Message Hash Hex:", messageHash.Hex())
```

### 3. Already Joined (Unlikely but possible)

Your code already checks this at line 224-229, but there could be a race condition.

### 4. Event is Full (Unlikely)

Your code checks this at line 238-240, but the on-chain state might have changed.

### 5. Invalid Signature

The signature recovery might be failing in the contract's `recoverSigner` function.

**Check**:

- Signature length should be 65 bytes (shown in debug output)
- Signature format matches what `recoverSigner` expects (v, r, s)
- The message being signed matches exactly what the contract expects

## Step-by-Step Debugging

### Step 1: Run Your Code

Look at the console output, specifically:

- The "Full Error" line
- The "🔴 REVERT REASON" line
- The "Transaction Details" section

### Step 2: Test Signature Recovery

Create a test script to verify signature recovery:

```go
// Test if signature recovers to expected address
recovered := cyptoutils.RecoverSignerAddress(messageHash, signature)
println("Recovered Address:", recovered.Hex())
println("Expected (CurrentUser):", currentUser.WalletAddress)
```

### Step 3: Check Access Control

Create a script to check roles:

```solidity
// In your smart contract testing environment
bool isHost = eventAccessManager.hasRole(HOST_ROLE, signerAddress);
bool isAdmin = eventAccessManager.hasRole(ADMIN_ROLE, signerAddress);
console.log("Signer has HOST role:", isHost);
console.log("Signer has ADMIN role:", isAdmin);
```

### Step 4: Manual Contract Call

Test the contract call directly:

```bash
# Using cast (foundry) or similar tool
cast call $CONTRACT_ADDRESS \
  "addParticipant(address,string,bytes32,bytes)" \
  $PARTICIPANT_ADDRESS \
  "$MESSAGE_HASH" \
  $DIGEST_BYTES32 \
  $SIGNATURE_BYTES \
  --from $TRANSACTOR_ADDRESS
```

## Most Likely Root Cause

Based on your smart contract code, the **#1 most likely issue** is:

### The signer (recovered from signature) doesn't have HOST/ADMIN role

**Why this is confusing:**

- Your transactor wallet might be authorized
- But the contract checks the **recovered signer** from the signature
- These are TWO DIFFERENT addresses!

**The Flow:**

1. Your backend signs with user's wallet address
2. Backend submits transaction with admin wallet (transactor)
3. Contract receives transaction from admin wallet
4. Contract **ignores** who sent it (msg.sender)
5. Contract **recovers signer** from the signature
6. Contract checks if **recovered signer** has HOST/ADMIN role
7. If not → REVERT

**Solution:**
The recovered signer (likely the user's wallet or the admin wallet that created the signature) needs to be granted HOST or ADMIN role in the EventAccessManager.

## Quick Fix to Test

If you want to quickly test if this is the issue, temporarily modify your smart contract:

```solidity
// TEMPORARY TEST ONLY - DO NOT USE IN PRODUCTION
function addParticipant(
    address participantAddress,
    string memory signedMessageDigest,
    bytes memory signature
) external {
    address signer = recoverSigner(signedMessageDigest, signature);

    // Log for debugging
    emit Debug("Signer", signer);
    emit Debug("MsgSender", msg.sender);

    // Temporarily comment out access control
    // EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer);

    // ... rest of the function
}
```

If the transaction succeeds with access control commented out, you've confirmed the issue!

## Need More Help?

Share the console output with:

1. The "Full Error" message
2. The "Transaction Details" section
3. The "Transactor Address" and "Participant Address"

Then I can help you pinpoint the exact issue!
