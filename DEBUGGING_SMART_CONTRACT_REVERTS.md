# Debugging Smart Contract Execution Reverts

## Overview

When you encounter "execution reverted" errors from smart contract calls, this guide will help you identify and fix the root cause.

## What I've Added

### 1. Enhanced Error Logging in `join_event.go`

Added comprehensive debug logging that shows:

- **Participant Address**: The wallet address being added
- **Message Hash**: The hash being used for signature verification
- **Signature Length**: Validates signature is correct size (65 bytes expected)
- **Contract Address**: The event contract being called
- **Transactor Address**: The admin wallet making the call
- **Current/Max Seats**: Shows capacity status before adding participant
- **Transaction Hash**: For blockchain explorer lookup
- **Gas Used**: Helps identify out-of-gas issues
- **Revert Reason**: Decoded error message from the smart contract

### 2. New Utility Functions in `cyptoutils/utils.go`

#### `GetRevertReason(ctx, client, tx, receipt)`

Attempts to extract the human-readable revert reason by replaying the failed transaction.

#### `DecodeRevertReason(errData)`

Decodes custom error selectors from the Event contract:

- `Event__InvalidEventName()`
- `Event__CannotReduceSeatsCount()`
- `Event__SeatsCountReached()` - Event is full
- `Event__ParticipantIsNotJoined()`
- `Event__ParticipantIsAlreadyJoined()` - Duplicate registration
- `Event__AddressCannotBeZero()` - Invalid address
- `Event__InvalidSignature()`

## How to Debug

### Step 1: Run Your Code and Check Console Output

When you trigger the join event flow, you'll now see detailed output like:

```
=== DEBUG AddParticipant ===
Participant Address: 0x1234...5678
Message Hash: 0xabcd...ef01
Signature Length: 65
Contract Address: 0x9876...5432
Transactor Address: 0x5555...6666
Current Seats: 5
Max Seats: 100
========================
Transaction submitted: 0x7890...1234
Transaction mined. Status: 0 Gas Used: 45823
Revert Reason: Event__ParticipantIsAlreadyJoined()
```

### Step 2: Identify the Root Cause

Common issues and their solutions:

#### A. `Event__SeatsCountReached()`

**Cause**: Event is at full capacity  
**Solution**:

- Check `Current Seats` vs `Max Seats` in debug output
- Increase event capacity if needed
- Remove inactive participants

#### B. `Event__ParticipantIsAlreadyJoined()`

**Cause**: User already registered for this event  
**Solution**:

- Check the `hasJoinedOnChain` logic (line 223-229)
- This is actually handled in your code, so shouldn't reach the contract
- May indicate cache/sync issue between DB and blockchain

#### C. Access Control Errors (requireHostOrAdmin failed)

**Cause**: The transactor wallet doesn't have HOST or ADMIN role  
**Solution**:

- Verify `Transactor Address` in debug output
- Check EventAccessManager contract for role assignments
- Ensure `GetKeyedTransactor()` returns the correct admin wallet
- Verify the admin wallet has been granted HOST role on the EventAccessManager

#### D. Signature Verification Failed

**Cause**: The signature doesn't match the message hash  
**Solution**:

- Check `Signature Length` (should be 65 bytes)
- Verify message hash format matches what contract expects
- Ensure signature recovery produces the expected address
- The signature should be signed by a wallet with HOST/ADMIN role

#### E. `Event__AddressCannotBeZero()`

**Cause**: Invalid participant address  
**Solution**: Check `Participant Address` in debug output isn't 0x0000...

### Step 3: Use Blockchain Explorer

Copy the transaction hash from the console output and paste it into your blockchain explorer:

- **Hardhat/Local**: Check console logs
- **Testnet**: Use Etherscan/Blockscout
- **Look for**: Revert reason, emitted events, internal transactions

### Step 4: Check Smart Contract State

You can add calls to check contract state before the transaction:

```go
// Check if the transactor has the right role
// (You may need to add bindings for EventAccessManager)

// Check event status
status, _ := eventContractInstance.EventStatus(&bind.CallOpts{Context: ctx})
println("Event Status:", status) // 0=ACTIVE, 1=INACTIVE, 2=CLOSED
```

## Common Debugging Commands

### Check Transaction Receipt

```go
receipt, err := client.TransactionReceipt(ctx, txHash)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %d\n", receipt.Status) // 1 = success, 0 = failed
fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
```

### Simulate Transaction Before Sending

```go
// Use CallContract to simulate without sending
msg := ethereum.CallMsg{
    From:  transactor.From,
    To:    &contractAddress,
    Data:  callData,
    Gas:   1000000,
}
result, err := client.CallContract(ctx, msg, nil)
if err != nil {
    println("Simulation failed:", err.Error())
    // This will show the revert reason without wasting gas
}
```

## Most Likely Issue in Your Case

Based on the smart contract code, the most common causes are:

### 1. **Access Control** (Most Likely)

The transactor (admin wallet) doesn't have HOST or ADMIN role on the EventAccessManager contract.

**How to fix:**

```solidity
// On the EventAccessManager contract, ensure:
eventAccessManager.grantHostRole(adminWallet, adminWallet);
// Or
eventAccessManager.grantAdminRole(adminWallet, adminWallet);
```

### 2. **Signature Format Mismatch**

The contract expects `string memory signedMessageDigest` but you're passing `messageHash.String()`.

**Check:**

- Look at line 268 in your code: `messageHash.String()`
- The contract uses this to recover the signer
- Ensure this matches the format the contract expects

### 3. **Already Joined**

Your code checks `hasJoinedOnChain` (line 224-229), but there might be a race condition or sync issue.

## Testing the Fix

1. Run your code and look at the console output
2. Check which specific error is occurring
3. Verify the transactor has proper permissions
4. Ensure the signature format matches contract expectations
5. Test with a blockchain explorer for detailed transaction info

## Removing Debug Logging (Later)

Once you've fixed the issue, you can:

1. Keep the error decoding (it's useful!)
2. Remove or comment out the `println` statements
3. Consider using proper logging library instead of `println`

```go
// Replace println with proper logging
log.Printf("DEBUG: Participant=%s, Contract=%s", participantAddr.Hex(), contractAddr)
```

## Additional Resources

- **Ethereum Error Codes**: https://docs.soliditylang.org/en/latest/control-structures.html#error-handling
- **Go-Ethereum Documentation**: https://geth.ethereum.org/docs
- **Custom Errors**: Solidity 0.8.4+ custom errors are gas-efficient but require decoding

---

**Pro Tip**: The most useful information for debugging is in the **Revert Reason** line. Everything else helps you understand the context of why that revert happened.
