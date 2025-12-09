# Comparison: join_event vs claim_certificate

## Key Differences

### 1. Signature Creation

**join_event:**

- Participant signs their own message
- Message format: `{"walletAddress":"...","contractAddress":"...","deadlineBlock":...}`
- Signed with `HashEthereumMessage`
- Signature stored temporarily, passed directly to contract

**claim_certificate:**

- Host signs message about receivers
- Message format: `{"eventContractAddress":"...","receivers":["hash1","hash2"]}`
- Signed during `import_certificate_receivers`, stored in DB
- Retrieved from DB and passed to contract

### 2. Contract Function Calls

**join_event:**

```go
eventContractInstance.AddParticipant(transactor, *participantAddress, signMessage, signature)
```

- Parameters: participantAddress, signMessage (string), signature (bytes)
- No access control pre-flight check

**claim_certificate:**

```go
certificateContractInstance.MintNft(
    transactor,
    receiverAddress,
    userId, certificateId, issuerId,
    encryptedUserData, backendEncryptedUserData,
    issuerAddresses,
    signMessageStr,  // signedMessageDigest parameter
    hostSignatureBytes,  // signature parameter
    hostSignatureStr,  // hostSignature parameter
    hostPublicKey,
    signMessageStr,  // signMessage parameter
    userEncryptedProof, backendEncryptedProof,
    certificateTitle, certificateSubtitle,
    userDataHashStr,
    issuerProofs,
)
```

- More complex with many parameters
- Access control pre-flight check (currently failing)

### 3. Access Control

**join_event:**

- Contract's `addParticipant` calls `recoverSigner` then checks access control internally
- No pre-flight check

**claim_certificate:**

- Contract's `mintNft` calls `recoverSigner` then `requireHostOrAdmin(signer, msg.sender)`
- Pre-flight check added but failing with "execution reverted"

### 4. The Critical Issue

The access control pre-flight check in `claim_certificate` is reverting:

```
⚠️ Could not check if host is registered
error="execution reverted"
```

This suggests:

1. The `CheckIsHostOrAdmin` view function call is failing
2. This might be due to DecmAccessManager contract issues
3. OR the EventAccessManager contract itself has issues

However, **join_event works without any pre-flight checks**, so the contract access control must be working there.

### Potential Root Cause

The actual `mintNft` transaction is reverting, and our pre-flight check is also reverting. This suggests:

1. The signature recovery in the contract might be failing (returns address(0))
2. The access control check in the contract is failing
3. There's a contract setup issue (EventAccessManager/DecmAccessManager)

Since `join_event` works with the same pattern, the issue is likely specific to:

- The certificate contract setup
- The EventAccessManager registration for this certificate/event
- The signature format or message format for certificate minting

### Recommendation

Since the access control pre-flight check is also failing, we should:

1. Remove or make the pre-flight check non-blocking (already done - it's a WARN)
2. Focus on getting better error messages from the actual transaction
3. Compare the exact contract addresses and setup between working join_event and failing claim_certificate
