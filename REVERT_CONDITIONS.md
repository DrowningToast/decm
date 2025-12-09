# Blockchain Revert Conditions in EventCertificate.mintNft

## Question

**What would be the command/condition that triggers revert on the blockchain?**

## All Possible Revert Conditions

Based on the contract code, here are ALL the conditions that can cause `mintNft` to revert:

### 1. Reentrancy Protection

**Location:** Function modifier

```solidity
external nonReentrant
```

**Revert Condition:**

- If the function is called recursively during execution
- **Error Message:** `ReentrancyGuard: reentrant call`

---

### 2. Invalid Signature Recovery

**Location:** Line 96

```solidity
address signer = recoverSigner(signedMessageDigest, signature);
```

**Revert Conditions:**

#### a) ECDSA Recovery Fails

```solidity
address signer = ECDSA.recover(ethSignedMessageHash, signature);
if (signer == address(0)) {
    revert Themis__InvalidSignature();
}
```

**Revert When:**

- Signature is malformed
- Signature doesn't match the message hash
- Recovery returns zero address
- **Error:** `Themis__InvalidSignature()` (custom error)

#### b) Signature Already Used

```solidity
if (usedSignatures[signature]) {
    revert Themis__SignatureAlreadyUsed();
}
```

**Revert When:**

- Same signature bytes have been used in a previous transaction
- **Error:** `Themis__SignatureAlreadyUsed()` (custom error)

---

### 3. Access Control Failure

**Location:** Line 97

```solidity
requireHostOrAdmin(signer, msg.sender);
```

**Revert Condition:**

```solidity
bool isAllowedMsgSender = EVENT_ACCESS_MANAGER.checkIsAllowedMsgSender(msgSender);
bool isHostOrAdmin = EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(signer);
if (!isHostOrAdmin && !isAllowedMsgSender) {
    require(false, "Not host or admin or allowed msg sender");
}
```

**Revert When:**

- ❌ Signer is NOT a host/admin in EventAccessManager, **AND**
- ❌ msg.sender (transactor) is NOT an allowed sender in DecmAccessManager
- **Error:** `"Not host or admin or allowed msg sender"`

**This is the MOST LIKELY cause of your revert!**

---

### 4. ERC721 \_safeMint Revert

**Location:** Line 123

```solidity
_safeMint(receiverAddress, tokenId);
```

**Revert Conditions:**

#### a) Receiver is Zero Address

- OpenZeppelin ERC721 reverts if `to == address(0)`
- **Error:** `ERC721: mint to the zero address`

#### b) Token Already Exists

- If `tokenId` already exists (shouldn't happen with counter)
- **Error:** `ERC721: token already minted`

#### c) Receiver is Contract Without ERC721Receiver

If `receiverAddress` is a contract:

```solidity
function _checkOnERC721Received(address from, address to, uint256 tokenId, bytes memory data)
    private returns (bool)
{
    if (to.isContract()) {
        try IERC721Receiver(to).onERC721Received(...) returns (bytes4 retval) {
            return retval == IERC721Receiver.onERC721Received.selector;
        } catch (bytes memory reason) {
            if (reason.length == 0) {
                revert("ERC721: transfer to non ERC721Receiver implementer");
            } else {
                assembly {
                    revert(add(32, reason), mload(reason))
                }
            }
        }
    } else {
        return true;
    }
}
```

**Revert When:**

- Receiver is a contract
- Contract doesn't implement `IERC721Receiver`
- Contract's `onERC721Received` returns wrong selector
- Contract's `onERC721Received` reverts
- **Error:** `"ERC721: transfer to non ERC721Receiver implementer"` or custom error from hook

---

### 5. Gas Limit Exceeded

**Revert When:**

- Transaction runs out of gas
- Parameters are too large (long strings, large arrays)
- **Error:** `"out of gas"` or `"gas required exceeds allowance"`

---

## Most Likely Revert Reason (Based on Your Logs)

Based on your error logs showing `"execution reverted"` with no detailed reason, the most likely cause is:

### **Access Control Failure** (Line 97)

```solidity
requireHostOrAdmin(signer, msg.sender);
```

**Condition:**

```
NOT (signer is host/admin) AND NOT (msg.sender is allowed sender)
```

**This means:**

1. Host address `0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E` is NOT registered as host/admin in EventAccessManager
2. System transactor `0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A` is NOT registered as allowed sender in DecmAccessManager

---

## How to Check Which Condition Triggered

### Using `cast` (Foundry)

```bash
# Simulate the transaction to get revert reason
cast call <CONTRACT_ADDRESS> \
  "mintNft(address,string,string,string,string,string,address[],string,bytes,string,string,string,string,string,string,string,string,(string,string)[])" \
  <receiverAddress> \
  "<userId>" \
  "<certificateId>" \
  "<issuerId>" \
  "<encryptedUserData>" \
  "<backendEncryptedUserData>" \
  "[<issuerAddresses>]" \
  "<signMessageStr>" \
  "<hostSignatureBytes>" \
  "<hostSignatureStr>" \
  "<hostPublicKey>" \
  "<signMessageStr>" \
  "<userEncryptedProof>" \
  "<backendEncryptedProof>" \
  "<certificateTitle>" \
  "<certificateSubtitle>" \
  "<userDataHashStr>" \
  "[<issuerProofs>]" \
  --rpc-url <RPC_URL> \
  --from <TRANSACTOR_ADDRESS>
```

### Check Access Control Manually

```bash
# Check if host is registered
cast call <EVENT_ACCESS_MANAGER_ADDRESS> \
  "checkIsHostOrAdmin(address)(bool)" \
  <HOST_ADDRESS> \
  --rpc-url <RPC_URL>

# Check if transactor is allowed
cast call <DECM_ACCESS_MANAGER_ADDRESS> \
  "checkIsAllowedMsgSender(address)(bool)" \
  <TRANSACTOR_ADDRESS> \
  --rpc-url <RPC_URL>
```

---

## Summary Table

| Condition           | Error Type                                           | Most Likely?                        |
| ------------------- | ---------------------------------------------------- | ----------------------------------- |
| Reentrancy          | `ReentrancyGuard: reentrant call`                    | ❌ No                               |
| Invalid Signature   | `Themis__InvalidSignature()`                         | ❌ No (backend check passes)        |
| Signature Replay    | `Themis__SignatureAlreadyUsed()`                     | ❌ No (pre-flight check passes)     |
| **Access Control**  | `"Not host or admin or allowed msg sender"`          | ✅ **YES**                          |
| Zero Address        | `ERC721: mint to the zero address`                   | ❌ No (pre-flight check)            |
| Token Exists        | `ERC721: token already minted`                       | ❌ No (uses counter)                |
| ERC721 Hook Failure | `ERC721: transfer to non ERC721Receiver implementer` | ⚠️ Possible if receiver is contract |
| Out of Gas          | `out of gas`                                         | ❌ Unlikely                         |

---

## Fix

To resolve the access control revert, ensure ONE of these is true:

1. **Register Host in EventAccessManager:**

    ```solidity
    EVENT_ACCESS_MANAGER.grantRole(HOST_ROLE, <HOST_ADDRESS>);
    ```

2. **Register Transactor as Allowed Sender in DecmAccessManager:**
    ```solidity
    DECM_ACCESS_MANAGER.allowMsgSender(<TRANSACTOR_ADDRESS>);
    ```
