# Smart Contract Transaction Revert Analysis

This document provides a comprehensive analysis of all possible transaction revert conditions in the DECM smart contracts.

## Overview

The DECM platform consists of the following main contracts:

1. **Event.sol** - Core event management
2. **EventAccessManager.sol** - Role-based access control for events
3. **EventCertificate.sol** - ERC721 certificate NFTs
4. **EventTicket.sol** - ERC721 ticket NFTs
5. **DecmAccessManager.sol** - Global access control
6. **ThemisUtils.sol** - Signature verification utilities
7. **StringUtils.sol** - String parsing utilities

---

## 1. Event.sol Revert Conditions

### Constructor

- **Line 57**: `require(false, "Access manager cannot be zero address")`
    - **Condition**: `eventAccessManagerAddr == address(0)`
    - **Context**: Constructor parameter validation

- **Line 285**: `require(false, "Invalid event name")`
    - **Condition**: `bytes(_eventName).length == 0`
    - **Context**: Event name validation in `_validateEventName()`

### updateEvent()

- **Line 80**: `EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 285**: `require(false, "Invalid event name")`
    - **Condition**: `bytes(_eventName).length == 0`
    - **Context**: Event name validation

- **Line 87**: `require(false, "Cannot reduce seats count")`
    - **Condition**: `_seatsCount < seatsCount`
    - **Context**: Prevents reducing available seats

- **Line 27**: `revert Themis__InvalidSignature()` (from ThemisUtils)
    - **Condition**: Signature recovery returns `address(0)`
    - **Context**: Invalid signature in `recoverSigner()`

- **Line 31**: `revert Themis__SignatureAlreadyUsed()` (from ThemisUtils)
    - **Condition**: Signature has already been used
    - **Context**: Signature replay protection

### addParticipant()

- **Line 120**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 125**: `require(false, "Address cannot be zero")`
    - **Condition**: `participantAddress == address(0)`
    - **Context**: Zero address validation

- **Line 130**: `require(false, "Seats count reached")`
    - **Condition**: `currentSeatsCount >= seatsCount`
    - **Context**: Event capacity check

- **Line 136**: `require(false, "Participant is already joined")`
    - **Condition**: `isParticipant[participantAddress] == true`
    - **Context**: Duplicate participant check

- **Line 248**: `EVENT_ACCESS_MANAGER.grantParticipantRoleUsingAllowedMsgSender()` - Can revert if:
    - `require(false, "Not allowed msg sender")` - msg.sender not in allowed list
    - `require(false, "Account cannot be zero")` - participantAddress is zero (already checked above)

### leaveEvent()

- **Line 160**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 161**: `EVENT_ACCESS_MANAGER.requireParticipant(signer, msg.sender)`
    - **Condition**: Signer is not participant AND msg.sender is not allowed msg sender
    - **Error**: "Not participant or allowed msg sender"
    - **Context**: Access control check

- **Line 167**: `require(false, "Participant is not joined")`
    - **Condition**: `!isParticipant[participantAddress]`
    - **Context**: Participant existence check

### removeParticipant()

- **Line 188**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 189**: `EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 193**: `require(false, "Address cannot be zero")`
    - **Condition**: `participantAddress == address(0)`
    - **Context**: Zero address validation

- **Line 198**: `require(false, "Participant is not joined")`
    - **Condition**: `!isParticipant[participantAddress]`
    - **Context**: Participant existence check

### confirmEvent()

- **Line 212**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 213**: `EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 217**: `require(false, "Event is closed")`
    - **Condition**: `eventStatus == EventStatus.CLOSED`
    - **Context**: Event status validation

- **Line 221**: `require(false, "Event is inactive")`
    - **Condition**: `eventStatus == EventStatus.INACTIVE`
    - **Context**: Event status validation

---

## 2. EventAccessManager.sol Revert Conditions

### Constructor

- **Line 26**: `require(false, "Access manager cannot be zero address")`
    - **Condition**: `decmAccessManagerAddr == address(0)`
    - **Context**: Constructor parameter validation

- **Line 30**: `require(false, "Account cannot be zero")`
    - **Condition**: `hostAddress == address(0)`
    - **Context**: Host address validation

### grantIssuerRole()

- **Line 42**: `requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 45**: `require(false, "Account cannot be zero")`
    - **Condition**: `issuer == address(0)`
    - **Context**: Zero address validation

### revokeIssuerRole()

- **Line 52**: `requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 55**: `require(false, "Account cannot be zero")`
    - **Condition**: `issuer == address(0)`
    - **Context**: Zero address validation

### grantParticipantRole()

- **Line 62**: `requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 65**: `require(false, "Account cannot be zero")`
    - **Condition**: `participant == address(0)`
    - **Context**: Zero address validation

### revokeParticipantRole()

- **Line 73**: `requireHostOrAdminOrParticipant(signer, msg.sender)`
    - **Condition**: Signer is not host/admin/participant AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or participant or allowed msg sender"
    - **Context**: Access control check

- **Line 76**: `require(false, "Account cannot be zero")`
    - **Condition**: `participant == address(0)`
    - **Context**: Zero address validation

### grantParticipantRoleUsingAllowedMsgSender()

- **Line 84**: `requireAllowedMsgSender(msgSender)`
    - **Condition**: `!checkIsAllowedMsgSender(msgSender)`
    - **Error**: "Not allowed msg sender"
    - **Context**: Access control check (only allowed msg senders can call)

- **Line 87**: `require(false, "Account cannot be zero")`
    - **Condition**: `participant == address(0)`
    - **Context**: Zero address validation

### grantHostRole()

- **Line 95**: `requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 98**: `require(false, "Account cannot be zero")`
    - **Condition**: `host == address(0)`
    - **Context**: Zero address validation

### requireAllowedMsgSender()

- **Line 127**: `require(false, "Not allowed msg sender")`
    - **Condition**: `!checkIsAllowedMsgSender(addr)`
    - **Context**: Access control validation

### requireHostOrAdmin()

- **Line 135**: `require(false, "Not host or admin or allowed msg sender")`
    - **Condition**: `!isHostOrAdmin && !isAllowedMsgSender`
    - **Context**: Access control validation

### requireAdmin()

- **Line 143**: `require(false, "Not admin or allowed msg sender")`
    - **Condition**: `!isAdmin && !isAllowedMsgSender`
    - **Context**: Access control validation

### requireHostOrAdminOrParticipant()

- **Line 153**: `require(false, "Not host or admin or participant or allowed msg sender")`
    - **Condition**: `!hasHostRole && !hasAdminRole && !hasParticipantRole && !isAllowedMsgSender`
    - **Context**: Access control validation

### requireParticipant()

- **Line 161**: `require(false, "Not participant or allowed msg sender")`
    - **Condition**: `!hasParticipantRole && !isAllowedMsgSender`
    - **Context**: Access control validation

---

## 3. EventCertificate.sol Revert Conditions

### Constructor

- **No explicit reverts** - Constructor only sets state variables

### mintNft()

- **Line 95**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 96**: `requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 121**: `_safeMint(receiverAddress, tokenId)` - ERC721 `_safeMint()` can revert if:
    - Receiver is zero address (OpenZeppelin ERC721 check)
    - Receiver is a contract that doesn't implement `IERC721Receiver.onERC721Received()`
    - Token already exists (unlikely with counter, but possible if counter is manipulated)

- **Line 58**: `require(false, "Not host or admin or allowed msg sender")` (from requireHostOrAdmin)

### participantSignedCertificate()

- **Line 150**: `require(false, "Certificate not valid")`
    - **Condition**: `tokenIdToStatus[tokenId] != CertificateStatus.VALID`
    - **Context**: Certificate status validation

- **Line 157**: `require(false, "Not participant")`
    - **Condition**: `vc.data.receiverAddress != msg.sender`
    - **Context**: Only certificate owner can sign

### revokeCertificate()

- **Line 244**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 245**: `requireHostOrAdmin(signer, msg.sender)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

### getTokenData()

- **Line 185**: `require(false, "Token id out of bounds")`
    - **Condition**: `tokenId >= tokenCounter`
    - **Context**: Token existence check

### tokenURI()

- **Line 267**: `require(false, "Token id out of bounds")`
    - **Condition**: `tokenId >= tokenCounter`
    - **Context**: Token existence check

---

## 4. EventTicket.sol Revert Conditions

### Constructor

- **Line 69**: `require(false, "Access manager cannot be zero address")`
    - **Condition**: `eventAccessManagerAddr == address(0)`
    - **Context**: Constructor parameter validation

- **Line 72**: `require(false, "Event address cannot be zero address")`
    - **Condition**: `eventAddr == address(0)`
    - **Context**: Constructor parameter validation

### mintNft()

- **Line 91**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 92**: `requireHostOrAdmin(signer)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 95**: `require(false, "Invalid receiver")`
    - **Condition**: `receiverAddress == address(0)`
    - **Context**: Zero address validation

- **Line 115**: `_safeMint(receiverAddress, tokenId)` - ERC721 `_safeMint()` can revert if:
    - Receiver is zero address (already checked above)
    - Receiver is a contract that doesn't implement `IERC721Receiver.onERC721Received()`
    - Token already exists (unlikely with counter)

### bulkMintParticipantTickets()

- **Line 146**: `recoverSigner()` - Can revert with:
    - `Themis__InvalidSignature()` - Invalid signature
    - `Themis__SignatureAlreadyUsed()` - Signature already used

- **Line 147**: `requireHostOrAdmin(signer)`
    - **Condition**: Signer is not host/admin AND msg.sender is not allowed msg sender
    - **Error**: "Not host or admin or allowed msg sender"
    - **Context**: Access control check

- **Line 151**: `require(false, "Invalid receiver")`
    - **Condition**: `params[i].receiverAddress == address(0)`
    - **Context**: Zero address validation (inside loop)

- **Line 168**: `_safeMint(params[i].receiverAddress, tokenId)` - ERC721 `_safeMint()` can revert if:
    - Receiver is zero address (already checked above)
    - Receiver is a contract that doesn't implement `IERC721Receiver.onERC721Received()`
    - Token already exists (unlikely with counter)

### getTokenData()

- **Line 193**: `require(false, "Token id out of bounds")`
    - **Condition**: `tokenId >= tokenCounter`
    - **Context**: Token existence check

### tokenURI()

- **Line 234**: `require(false, "Token id out of bounds")`
    - **Condition**: `tokenId >= tokenId >= tokenCounter`
    - **Context**: Token existence check

---

## 5. DecmAccessManager.sol Revert Conditions

### Constructor

- **Line 24**: `require(false, "Admin cannot be zero address")`
    - **Condition**: `initialAdmins[i] == address(0)` (inside loop)
    - **Context**: Initial admin validation

### grantAdminRole()

- **Line 32**: `onlyRole(DEFAULT_ADMIN_ROLE)` - OpenZeppelin AccessControl modifier
    - **Condition**: Caller doesn't have DEFAULT_ADMIN_ROLE
    - **Error**: "AccessControl: account {account} is missing role {role}"
    - **Context**: Access control check

- **Line 34**: `require(false, "Admin cannot be zero address")`
    - **Condition**: `admin == address(0)`
    - **Context**: Zero address validation

### revokeAdminRole()

- **Line 42**: `onlyRole(DEFAULT_ADMIN_ROLE)` - OpenZeppelin AccessControl modifier
    - **Condition**: Caller doesn't have DEFAULT_ADMIN_ROLE
    - **Error**: "AccessControl: account {account} is missing role {role}"
    - **Context**: Access control check

- **Line 44**: `require(false, "Admin cannot be zero address")`
    - **Condition**: `admin == address(0)`
    - **Context**: Zero address validation

### addAllowedMsgSender()

- **Line 54**: `onlyRole(DEFAULT_ADMIN_ROLE)` - OpenZeppelin AccessControl modifier
    - **Condition**: Caller doesn't have DEFAULT_ADMIN_ROLE
    - **Error**: "AccessControl: account {account} is missing role {role}"
    - **Context**: Access control check

### removeAllowedMsgSender()

- **Line 58**: `onlyRole(DEFAULT_ADMIN_ROLE)` - OpenZeppelin AccessControl modifier
    - **Condition**: Caller doesn't have DEFAULT_ADMIN_ROLE
    - **Error**: "AccessControl: account {account} is missing role {role}"
    - **Context**: Access control check

---

## 6. ThemisUtils.sol Revert Conditions

### recoverSigner()

- **Line 27**: `revert Themis__InvalidSignature()`
    - **Condition**: `signer == address(0)` after ECDSA recovery
    - **Context**: Signature recovery failed (invalid signature format)

- **Line 31**: `revert Themis__SignatureAlreadyUsed()`
    - **Condition**: `usedSignatures[signature] == true`
    - **Context**: Signature replay protection

---

## 7. StringUtils.sol Revert Conditions

### splitSignMessage()

- **Line 58**: `require(foundFirst && secondComma > 0, "Invalid Sign Message")`
    - **Condition**: Message doesn't contain at least 2 commas
    - **Context**: Sign message format validation

### toAddress() / toUint256()

- **Line 29/33**: `s.parseAddress()` / `s.parseUint()` - OpenZeppelin Strings library
    - **Condition**: String cannot be parsed as address/uint256
    - **Error**: OpenZeppelin library error (format may vary)
    - **Context**: String parsing validation

---

## 8. OpenZeppelin ERC721 Revert Conditions

### \_safeMint() (used in EventCertificate and EventTicket)

- **Zero Address Check**: Reverts if `to == address(0)`
    - **Error**: "ERC721: mint to the zero address"
    - **Note**: Already checked in EventTicket, but not explicitly in EventCertificate

- **Token Existence Check**: Reverts if token already exists
    - **Error**: "ERC721: token already minted"
    - **Note**: Unlikely with counter-based token IDs, but possible if counter is manipulated

- **ERC721Receiver Check**: Reverts if `to` is a contract that doesn't implement `IERC721Receiver.onERC721Received()`
    - **Error**: "ERC721: transfer to non ERC721Receiver implementer"
    - **Context**: Safe transfer to contracts

### AccessControl Modifiers

- **onlyRole()**: Reverts if caller doesn't have the required role
    - **Error**: "AccessControl: account {account} is missing role {role}"
    - **Context**: Used in DecmAccessManager

---

## Summary by Revert Category

### 1. Access Control Reverts

- **Most Common**: "Not host or admin or allowed msg sender"
- **Other**: "Not participant or allowed msg sender", "Not admin or allowed msg sender", "Not allowed msg sender"
- **OpenZeppelin**: "AccessControl: account {account} is missing role {role}"

### 2. Signature Reverts

- **Themis\_\_InvalidSignature()**: Invalid signature format
- **Themis\_\_SignatureAlreadyUsed()**: Signature replay protection

### 3. Zero Address Reverts

- "Access manager cannot be zero address"
- "Account cannot be zero"
- "Address cannot be zero"
- "Invalid receiver"
- "Admin cannot be zero address"
- "Event address cannot be zero address"
- ERC721: "mint to the zero address"

### 4. Validation Reverts

- "Invalid event name" (empty string)
- "Cannot reduce seats count"
- "Seats count reached"
- "Participant is already joined"
- "Participant is not joined"
- "Event is closed"
- "Event is inactive"
- "Certificate not valid"
- "Not participant" (certificate signing)
- "Token id out of bounds"
- "Invalid Sign Message" (StringUtils)

### 5. ERC721 Reverts

- "ERC721: mint to the zero address"
- "ERC721: token already minted"
- "ERC721: transfer to non ERC721Receiver implementer"

### 6. String Parsing Reverts

- OpenZeppelin Strings library errors for invalid address/uint256 parsing

---

## Recommendations

1. **Add explicit zero address checks** in EventCertificate.mintNft() before `_safeMint()`
2. **Consider using custom errors** instead of `require(false, "message")` for gas optimization
3. **Add event status checks** in mintNft functions to prevent minting for closed/inactive events
4. **Document signature format requirements** to prevent Themis\_\_InvalidSignature errors
5. **Add comprehensive error handling** in frontend/backend to catch and handle all revert conditions gracefully

---

## Testing Recommendations

1. Test all access control scenarios (host, admin, participant, allowed msg sender)
2. Test signature replay protection
3. Test zero address validations
4. Test capacity limits (seats count)
5. Test event status transitions
6. Test ERC721Receiver contract interactions
7. Test invalid signature formats
8. Test string parsing edge cases

---

**Last Updated**: Generated from contract analysis
**Contracts Analyzed**: Event.sol, EventAccessManager.sol, EventCertificate.sol, EventTicket.sol, DecmAccessManager.sol, ThemisUtils.sol, StringUtils.sol
