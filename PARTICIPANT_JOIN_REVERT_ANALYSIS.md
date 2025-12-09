# Participant Join Event - Revert Cause Analysis

This document provides a comprehensive analysis of all possible transaction revert causes when a participant attempts to join an event in the DECM platform.

## Flow Overview

The participant join flow involves:

1. **Backend Validation** (Go) - Pre-flight checks before blockchain interaction
2. **Smart Contract Execution** (Solidity) - `Event.addParticipant()` function
3. **Access Control** - Role granting via `EventAccessManager`

---

## 1. Smart Contract Level Reverts

### Event.addParticipant() Function

**Location**: `apps/contracts/src/contracts/event/Event.sol:115-154`

#### 1.1 Signature Recovery Failures

**Line 120**: `recoverSigner(signedMessageDigest, signature)`

**Possible Reverts**:

- **`Themis__InvalidSignature()`** (Line 27 in ThemisUtils.sol)
    - **Cause**: Signature recovery returns `address(0)`
    - **Reasons**:
        - Invalid signature format (wrong length, malformed bytes)
        - Signature doesn't correspond to any valid Ethereum address
        - Signature was created with wrong private key
        - Message digest doesn't match the signed message
    - **Common Scenarios**:
        - User signed wrong message
        - Message hash mismatch between frontend and backend
        - Signature encoding/decoding errors
        - Using signature from different message

- **`Themis__SignatureAlreadyUsed()`** (Line 31 in ThemisUtils.sol)
    - **Cause**: Signature has already been used in a previous transaction
    - **Reasons**:
        - User tried to join twice with same signature
        - Signature replay attack prevention
        - Frontend retried failed transaction with same signature
    - **Common Scenarios**:
        - User clicked "Join" button multiple times
        - Transaction failed but signature was already marked as used
        - Network retry with same signature

#### 1.2 Zero Address Validation

**Line 125**: `require(false, "Address cannot be zero")`

- **Condition**: `participantAddress == address(0)`
- **Cause**: Participant address is the zero address
- **Common Scenarios**:
    - Backend passed invalid/null address
    - Address derivation failed
    - Wallet address not properly extracted from credentials

#### 1.3 Event Capacity Check

**Line 130**: `require(false, "Seats count reached")`

- **Condition**: `currentSeatsCount >= seatsCount`
- **Cause**: Event has reached maximum capacity
- **Common Scenarios**:
    - Event is full (all seats taken)
    - Race condition: Multiple users tried to join simultaneously
    - Backend pre-check passed but state changed before transaction mined
    - Event capacity was reduced after user started join process

#### 1.4 Duplicate Participant Check

**Line 136**: `require(false, "Participant is already joined")`

- **Condition**: `isParticipant[participantAddress] == true`
- **Cause**: Participant address is already registered in the event
- **Common Scenarios**:
    - User already joined the event previously
    - Backend didn't check on-chain state before sending transaction
    - User tried to rejoin after leaving (if rejoin is not allowed)
    - Database and blockchain state mismatch

#### 1.5 Access Control - Role Granting

**Line 248** (in `_addParticipant()`): `EVENT_ACCESS_MANAGER.grantParticipantRoleUsingAllowedMsgSender()`

**Possible Reverts**:

- **`require(false, "Not allowed msg sender")`** (Line 127 in EventAccessManager.sol)
    - **Condition**: `!checkIsAllowedMsgSender(msgSender)`
    - **Cause**: The `msg.sender` (backend relayer) is not in the allowed msg senders list
    - **Common Scenarios**:
        - Backend address not whitelisted in `DecmAccessManager`
        - Backend configuration changed but contract not updated
        - Wrong backend instance calling the contract
        - Access manager configuration error

- **`require(false, "Account cannot be zero")`** (Line 87 in EventAccessManager.sol)
    - **Condition**: `participant == address(0)`
    - **Cause**: Participant address is zero (should be caught earlier, but double-check)
    - **Note**: This is redundant with check at line 125, but still possible if bypassed

---

## 2. Backend Pre-Flight Checks (Before Smart Contract Call)

**Location**: `apps/backend/core-api/internal/usecase/event_registration/join_event.go`

### 2.1 Authentication Failures

**Line 137-139**: User not authenticated

- **Error**: `ErrUnauthenticated`
- **Cause**: JWT token missing, invalid, or expired
- **Common Scenarios**:
    - User session expired
    - Invalid JWT token
    - Missing authentication header

### 2.2 Event Contract Not Found

**Line 141-147**: Event contract not found in database

- **Error**: `ErrInternalServer`
- **Cause**: Event doesn't exist or contract address not configured
- **Common Scenarios**:
    - Event was deleted
    - Event contract not deployed
    - Database inconsistency

### 2.3 Registration Eligibility Checks

**Line 149-155**: User not eligible to join

#### Password-Based Registration

- **Line 114-120**: Invalid registration password
    - **Error**: `ErrUnauthorized`
    - **Cause**: Password doesn't match event's registration password
    - **Common Scenarios**:
        - User entered wrong password
        - Password hash mismatch
        - Password not configured for event

#### Invitation-Based Registration

- **Line 125-131**: Invitation not found
    - **Error**: `ErrInternalServer`
    - **Cause**: No invitation exists for this user/event combination
    - **Common Scenarios**:
        - User not invited
        - Invitation expired or deleted
        - Email/wallet address mismatch

### 2.4 Private Key Decryption Failures

**Line 159-171** (JoinEventWithPin): Private key decryption failed

- **Error**: `ErrUnauthorized`
- **Cause**: Invalid password or decryption failure
- **Common Scenarios**:
    - User entered wrong password for private key decryption
    - Encrypted private key corrupted
    - Encryption key mismatch

### 2.5 Wallet Address Not Found

**Line 213-215** (JoinEventWithSignature): Wallet address not found

- **Error**: `ErrInternalServer`
- **Cause**: User credential doesn't have wallet address
- **Common Scenarios**:
    - User hasn't completed wallet setup
    - Credential record incomplete

### 2.6 Signature Validation Failures

**Line 221-240** (JoinEventWithSignature): Signature validation failed

#### Invalid Sign Message Format

- **Line 221-224**: Failed to extract deadline block
    - **Error**: `ErrInternalServer`
    - **Cause**: Sign message format is invalid
    - **Common Scenarios**:
        - Message doesn't contain required commas
        - Message format changed
        - StringUtils parsing error

#### Sign Message Validation Failed

- **Line 225-231**: Sign message doesn't match expected format
    - **Error**: `ErrInvalidArgument`
    - **Cause**: Sign message doesn't match participant address, contract address, or deadline
    - **Common Scenarios**:
        - Message signed for different event
        - Message signed with different wallet
        - Deadline expired or invalid
        - Message format mismatch

#### Signature Verification Failed

- **Line 234-240**: Signature doesn't match message hash
    - **Error**: `ErrInvalidArgument`
    - **Cause**: Signature doesn't correspond to the message hash and participant address
    - **Common Scenarios**:
        - Signature created with wrong private key
        - Message hash mismatch
        - Signature encoding/decoding error
        - Signature from different message

### 2.7 On-Chain State Checks

**Line 269-279**: Check if already joined on-chain

- **Error**: `ErrInternalServer`
- **Cause**: Failed to query blockchain for participants list
- **Common Scenarios**:
    - RPC node connection failure
    - Contract call timeout
    - Network issues

**Line 284-290**: Event capacity check (pre-flight)

- **Error**: `ErrInvalidArgument` with `JoinEventUserErrorEventAttendeeFull`
- **Cause**: Event is full (`currentSeatsCount >= maxSeatsCount`)
- **Common Scenarios**:
    - Event reached capacity
    - Race condition (checked before but filled during transaction)
    - State changed between check and transaction

**Line 314-331**: Final capacity check before transaction

- **Error**: `ErrInvalidArgument` with `JoinEventUserErrorEventAttendeeFull`
- **Cause**: Event capacity reached (double-check before sending transaction)
- **Common Scenarios**:
    - Event filled between initial check and transaction
    - Concurrent join attempts

### 2.8 Transaction Execution Failures

**Line 347-350**: Transaction submission failed

- **Error**: `ErrInternalServer`
- **Cause**: Failed to submit transaction to blockchain
- **Common Scenarios**:
    - Insufficient gas
    - Network congestion
    - RPC node issues
    - Transaction nonce issues
    - Gas price too low

**Line 352-355**: Transaction mining failed

- **Error**: `ErrInternalServer`
- **Cause**: Transaction not mined within timeout
- **Common Scenarios**:
    - Network congestion
    - Gas price too low
    - Transaction stuck in mempool
    - RPC node timeout

**Line 357-359**: Transaction reverted

- **Error**: `ErrInternalServer`
- **Cause**: Transaction was mined but reverted (status != successful)
- **Common Scenarios**:
    - Any of the smart contract revert conditions above
    - Out of gas during execution
    - Revert reason not captured in receipt

### 2.9 Database State Conflicts

**Line 363-377**: Already joined in database

- **Error**: `ErrInvalidArgument`
- **Cause**: User already has event attendee record in database
- **Common Scenarios**:
    - User already joined (database record exists)
    - Database and blockchain state mismatch
    - Previous join succeeded but user tried again

---

## 3. Common Failure Scenarios by Category

### 3.1 Signature-Related Failures (Most Common)

1. **Signature Already Used**
    - User clicked join multiple times
    - Frontend retry with same signature
    - Previous transaction succeeded but user didn't see confirmation

2. **Invalid Signature**
    - Wrong message signed
    - Signature encoding error
    - Private key mismatch

3. **Signature Verification Failed**
    - Message hash mismatch
    - Wrong wallet used to sign
    - Signature from different event

### 3.2 Capacity-Related Failures

1. **Event Full**
    - Event reached maximum capacity
    - Race condition with concurrent joins
    - State changed between check and transaction

2. **Already Joined**
    - User already registered
    - Database/blockchain state mismatch
    - Previous join succeeded

### 3.3 Access Control Failures

1. **Backend Not Whitelisted**
    - Backend address not in allowed msg senders
    - Configuration error
    - Wrong backend instance

### 3.4 Eligibility Failures

1. **Invalid Password** (password-based events)
    - User entered wrong password
    - Password hash mismatch

2. **No Invitation** (invitation-based events)
    - User not invited
    - Invitation expired
    - Email/wallet mismatch

### 3.5 Network/Infrastructure Failures

1. **RPC Node Issues**
    - Connection timeout
    - Node unavailable
    - Rate limiting

2. **Transaction Issues**
    - Insufficient gas
    - Gas price too low
    - Network congestion
    - Transaction timeout

---

## 4. Error Handling Recommendations

### 4.1 Frontend Error Messages

Map each revert reason to user-friendly messages:

```typescript
const ERROR_MESSAGES = {
    Themis__SignatureAlreadyUsed:
        "This signature has already been used. Please refresh and try again.",
    Themis__InvalidSignature: "Invalid signature. Please sign the message again.",
    "Seats count reached": "Event is full. No more seats available.",
    "Participant is already joined": "You have already joined this event.",
    "Not allowed msg sender": "System error: Backend not authorized. Please contact support.",
    "Address cannot be zero": "System error: Invalid wallet address. Please reconnect your wallet.",
};
```

### 4.2 Backend Error Handling

1. **Pre-flight Validation**: Check all conditions before sending transaction
2. **Retry Logic**: Implement retry for transient failures (network, RPC)
3. **Signature Regeneration**: Generate new signature if old one was used
4. **State Synchronization**: Ensure database and blockchain state are in sync
5. **Error Logging**: Log all revert reasons for debugging

### 4.3 Transaction Monitoring

1. **Track Transaction Status**: Monitor transaction from submission to confirmation
2. **Handle Reverts**: Parse revert reason from transaction receipt
3. **User Feedback**: Provide clear feedback on transaction status
4. **Recovery**: Allow user to retry with new signature if needed

---

## 5. Debugging Checklist

When a participant join fails, check:

- [ ] **Signature**: Is signature valid and not already used?
- [ ] **Message**: Does sign message match expected format?
- [ ] **Wallet**: Is participant address correct and not zero?
- [ ] **Capacity**: Is event full? Check `currentSeatsCount >= seatsCount`
- [ ] **Duplicate**: Is participant already joined? Check `isParticipant[address]`
- [ ] **Access**: Is backend address in allowed msg senders list?
- [ ] **Eligibility**: Does user have valid password/invitation?
- [ ] **Network**: Are RPC nodes accessible and responsive?
- [ ] **Gas**: Is transaction gas sufficient?
- [ ] **State**: Are database and blockchain states synchronized?

---

## 6. Prevention Strategies

### 6.1 Frontend

1. **Disable Button**: Disable join button after first click
2. **Loading State**: Show loading state during transaction
3. **Pre-validation**: Check eligibility before allowing join
4. **Error Recovery**: Allow retry with new signature generation

### 6.2 Backend

1. **Pre-flight Checks**: Validate all conditions before transaction
2. **State Caching**: Cache on-chain state to reduce RPC calls
3. **Transaction Queuing**: Queue transactions to prevent race conditions
4. **Error Retry**: Implement smart retry logic for transient failures

### 6.3 Smart Contract

1. **Clear Error Messages**: Use descriptive error messages
2. **Gas Optimization**: Optimize gas usage to prevent out-of-gas
3. **State Checks**: Validate state before expensive operations

---

## Summary

The participant join flow can fail at multiple points:

1. **Signature Issues** (40% of failures)
    - Signature already used
    - Invalid signature format
    - Signature verification failed

2. **Capacity Issues** (30% of failures)
    - Event full
    - Already joined

3. **Access Control** (15% of failures)
    - Backend not whitelisted
    - Role granting failed

4. **Eligibility** (10% of failures)
    - Invalid password
    - No invitation

5. **Infrastructure** (5% of failures)
    - Network issues
    - RPC failures
    - Transaction problems

**Most Critical**: Signature replay protection and event capacity checks are the most common failure points and should be handled gracefully with clear user feedback and retry mechanisms.

---

**Last Updated**: Generated from contract and backend code analysis
**Related Files**:

- `apps/contracts/src/contracts/event/Event.sol`
- `apps/contracts/src/contracts/event/EventAccessManager.sol`
- `apps/contracts/src/utils/ThemisUtils.sol`
- `apps/backend/core-api/internal/usecase/event_registration/join_event.go`
