# How to Check Signature Recovery and Access Control with Foundry

## Step 1: Recover Signer from Signature

The `EventCertificate` contract inherits `recoverSigner` from `ThemisUtils`. You can call it directly:

```bash
# Replace with your actual values
CERTIFICATE_CONTRACT="0x14000891290"  # From your logs
SIGN_MESSAGE='{"eventContractAddress":"0x1cb18A232e8955B77Ffb87771D343e7a5cCa8769","receivers":["0x148d532c97fb3f21940c9f6923ab7b6a7df0489091da9fcfd4925fe05bdc49af"]}'
HOST_SIGNATURE="0x06f43116545c56945f41b07afd6eb0924f9b4171fb86aa9c9d22ee55b95f7145730546167699ba05424c630928caf7be58a23eb28ca52fba1116f89d09e5fbc61b"
RPC_URL="http://localhost:8545"  # Or your RPC endpoint

# Call recoverSigner - NOTE: This will mark signature as used!
cast send $CERTIFICATE_CONTRACT \
  "recoverSigner(string,bytes)(address)" \
  "$SIGN_MESSAGE" \
  "$HOST_SIGNATURE" \
  --rpc-url $RPC_URL \
  --private-key $YOUR_PRIVATE_KEY

# OR use cast call if you want to simulate (but recoverSigner is not view, it's state-changing)
# So you might need to use a fork or test environment
```

**Warning:** `recoverSigner` is `public returns`, not `view`, because it marks the signature as used. This means:

- You can call it, but it will mark the signature as used (state-changing)
- Better to simulate in a test environment or use `cast call` on a forked chain

---

## Step 2: Check Access Control (Better Approach)

Instead of calling `recoverSigner` (which uses gas and marks signature as used), check access control directly:

### 2a. Get EventAccessManager Address

```bash
# Get the EventAccessManager address from EventCertificate
EVENT_ACCESS_MANAGER=$(cast call $CERTIFICATE_CONTRACT \
  "EVENT_ACCESS_MANAGER()(address)" \
  --rpc-url $RPC_URL)

echo "EventAccessManager: $EVENT_ACCESS_MANAGER"
```

### 2b. Check if Host Address is Registered as Host/Admin

```bash
HOST_ADDRESS="0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E"  # From your logs

# Check if host is registered
cast call $EVENT_ACCESS_MANAGER \
  "checkIsHostOrAdmin(address)(bool)" \
  $HOST_ADDRESS \
  --rpc-url $RPC_URL

# Should return: true or false
```

### 2c. Check if Host has HOST_ROLE

```bash
# Check if address has HOST_ROLE specifically
cast call $EVENT_ACCESS_MANAGER \
  "checkIsHost(address)(bool)" \
  $HOST_ADDRESS \
  --rpc-url $RPC_URL
```

### 2d. Check if System Transactor is Allowed Sender

```bash
# Get DecmAccessManager address from EventAccessManager
DECM_ACCESS_MANAGER=$(cast call $EVENT_ACCESS_MANAGER \
  "DECM_ACCESS_MANAGER()(address)" \
  --rpc-url $RPC_URL)

echo "DecmAccessManager: $DECM_ACCESS_MANAGER"

# Check if transactor is allowed
TRANSACTOR_ADDRESS="0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A"  # From your logs

cast call $DECM_ACCESS_MANAGER \
  "checkIsAllowedMsgSender(address)(bool)" \
  $TRANSACTOR_ADDRESS \
  --rpc-url $RPC_URL
```

---

## Step 3: Manual Signature Verification (Off-chain)

Since `recoverSigner` is state-changing, better to verify off-chain:

```bash
# Install cast if not already installed
# forge install foundry-rs/foundry

# Use cast to verify signature (this uses ECDSA recovery)
# Note: You need to hash the message with Ethereum prefix first

# The message that should be signed:
SIGN_MESSAGE='{"eventContractAddress":"0x1cb18A232e8955B77Ffb87771D343e7a5cCa8769","receivers":["0x148d532c97fb3f21940c9f6923ab7b6a7df0489091da9fcfd4925fe05bdc49af"]}'

# Cast can't directly verify, but you can use a Solidity script or test
```

---

## Step 4: Using Foundry Test Script

Create a test script to check everything:

```solidity
// test/CheckSignature.t.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {EventCertificate} from "../src/contracts/event/EventCertificate.sol";
import {EventAccessManager} from "../src/contracts/event/EventAccessManager.sol";

contract CheckSignatureTest is Test {
    function testCheckSignature() public {
        // Setup - replace with your addresses
        address certificateContract = 0x14000891290;
        address hostAddress = 0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E;
        address transactor = 0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A;

        bytes memory signature = hex"06f43116545c56945f41b07afd6eb0924f9b4171fb86aa9c9d22ee55b95f7145730546167699ba05424c630928caf7be58a23eb28ca52fba1116f89d09e5fbc61b";
        string memory signMessage = '{"eventContractAddress":"0x1cb18A232e8955B77Ffb87771D343e7a5cCa8769","receivers":["0x148d532c97fb3f21940c9f6923ab7b6a7df0489091da9fcfd4925fe05bdc49af"]}';

        EventCertificate cert = EventCertificate(certificateContract);

        // Recover signer (WARNING: This marks signature as used!)
        address recoveredSigner = cert.recoverSigner(signMessage, signature);
        console.log("Recovered signer:", recoveredSigner);
        console.log("Expected host:", hostAddress);
        console.log("Match:", recoveredSigner == hostAddress);

        // Check access control
        address eventAccessManager = cert.EVENT_ACCESS_MANAGER();
        EventAccessManager eam = EventAccessManager(eventAccessManager);

        bool isHostOrAdmin = eam.checkIsHostOrAdmin(recoveredSigner);
        bool isAllowedSender = eam.checkIsAllowedMsgSender(transactor);

        console.log("Is host/admin:", isHostOrAdmin);
        console.log("Is allowed sender:", isAllowedSender);
        console.log("Access granted:", isHostOrAdmin || isAllowedSender);
    }
}
```

Run with:

```bash
forge test --match-test testCheckSignature -vvv
```

---

## Step 5: Quick One-Liner Checks

```bash
# Complete check script
#!/bin/bash

CERT_CONTRACT="0x14000891290"
HOST_ADDR="0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E"
TRANSACTOR="0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A"
RPC="http://localhost:8545"

# Get EventAccessManager
EAM=$(cast call $CERT_CONTRACT "EVENT_ACCESS_MANAGER()(address)" --rpc-url $RPC)
echo "EventAccessManager: $EAM"

# Check host
echo -n "Host is registered: "
cast call $EAM "checkIsHostOrAdmin(address)(bool)" $HOST_ADDR --rpc-url $RPC

# Get DecmAccessManager
DAM=$(cast call $EAM "DECM_ACCESS_MANAGER()(address)" --rpc-url $RPC)
echo "DecmAccessManager: $DAM"

# Check transactor
echo -n "Transactor is allowed: "
cast call $DAM "checkIsAllowedMsgSender(address)(bool)" $TRANSACTOR --rpc-url $RPC
```

---

## Summary

**To check if signature is from host:**

1. **Recover signer:** Call `recoverSigner(signMessage, signature)` on EventCertificate (⚠️ marks signature as used)
2. **Compare addresses:** Check if recovered address matches expected host address
3. **Check access control:** Use `checkIsHostOrAdmin(recoveredAddress)` on EventAccessManager

**Best approach:** Use `cast call` to check access control directly without consuming the signature!
