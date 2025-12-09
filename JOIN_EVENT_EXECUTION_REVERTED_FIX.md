# Fix: join_event Execution Reverted

## Problem

The `join_event` use case started failing with "execution reverted" errors even though it was working before. This is the **same root cause** as the `claim_certificate` issue.

## Root Cause

The `Event.addParticipant` function calls:

```solidity
EVENT_ACCESS_MANAGER.grantParticipantRoleUsingAllowedMsgSender(participantAddress, msgSender)
```

This function requires that the system transactor (`msg.sender`) is registered as an **allowed message sender** in `DecmAccessManager`. If it's not registered, the transaction reverts with:

```
"Not allowed msg sender"
```

## Why This Started Failing

Possible reasons it stopped working:

1. **DecmAccessManager was redeployed** and lost its configuration
2. **The system transactor was removed** from allowed senders
3. **New environment setup** didn't include the registration step
4. **Contract migration** reset the access control mappings

## The Exact Revert Point

In `Event.sol`, line 248:

```solidity
function _addParticipant(address participantAddress, address msgSender) private {
    // ...
    EVENT_ACCESS_MANAGER.grantParticipantRoleUsingAllowedMsgSender(participantAddress, msgSender);
}
```

Which calls `EventAccessManager.grantParticipantRoleUsingAllowedMsgSender()`:

```solidity
function grantParticipantRoleUsingAllowedMsgSender(address participant, address msgSender) public {
    requireAllowedMsgSender(msgSender);  // ← REVERTS HERE if not allowed
    // ...
}
```

Which checks `DecmAccessManager`:

```solidity
function requireAllowedMsgSender(address addr) public view {
    if (!checkIsAllowedMsgSender(addr)) {
        require(false, "Not allowed msg sender");  // ← THE REVERT
    }
}
```

## Solution

Register the system transactor in `DecmAccessManager`:

### Option 1: Use the Diagnostic Script

```bash
./check_and_fix_transactor_registration.sh
```

This script will:

1. Derive the transactor address from your private key
2. Check if it's registered
3. Automatically register it if you have admin permissions
4. Provide manual instructions if you don't have admin permissions

### Option 2: Manual Registration with Foundry

```bash
# Set your environment variables
export RPC_URL="https://eth-sepolia.g.alchemy.com/v2/..."
export PRIVATE_KEY="your_private_key"
export DECM_ACCESS_MANAGER_ADDRESS="0x..."

# Get transactor address
TRANSACTOR=$(cast wallet address --private-key $PRIVATE_KEY)

# Register it (requires admin private key)
cast send $DECM_ACCESS_MANAGER_ADDRESS \
  "addAllowedMsgSender(address)" \
  $TRANSACTOR \
  --private-key <ADMIN_PRIVATE_KEY> \
  --rpc-url $RPC_URL
```

### Option 3: Using Foundry Script

Create a script `scripts/RegisterTransactor.s.sol`:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {DecmAccessManager} from "../src/contracts/decm/DecmAccessManager.sol";

contract RegisterTransactor is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address damAddress = vm.envAddress("DECM_ACCESS_MANAGER_ADDRESS");
        address transactorAddress = vm.envAddress("TRANSACTOR_ADDRESS");

        vm.startBroadcast(deployerPrivateKey);

        DecmAccessManager dam = DecmAccessManager(damAddress);
        dam.addAllowedMsgSender(transactorAddress);

        vm.stopBroadcast();
    }
}
```

## Verification

After registration, verify it worked:

```bash
cast call $DECM_ACCESS_MANAGER_ADDRESS \
  "allowedMsgSenders(address)(bool)" \
  $TRANSACTOR \
  --rpc-url $RPC_URL
```

Should return: `true`

## Diagnostic Logging

I've added comprehensive logging to `join_event.go` that will:

1. **Pre-flight check**: Verify if the transactor is allowed before attempting the transaction
2. **Better error messages**: Extract and display the exact revert reason
3. **Actionable guidance**: Log suggestions for fixing the issue

The logs will show:

```
🔐 Access control pre-flight check
   is_transactor_allowed: false  ← If false, registration is needed
   note: "Register transactor in DecmAccessManager using addAllowedMsgSender() function"
```

## Other Possible Revert Reasons

While the most likely cause is the transactor not being registered, other possible reasons:

1. **Signature replay**: Signature already used (`Themis__SignatureAlreadyUsed`)
2. **Invalid signature**: Signature format issue (`Themis__InvalidSignature`)
3. **Seats count reached**: Event is full
4. **Already joined**: Participant already registered
5. **Zero address**: Invalid participant address

The diagnostic logging will help identify which one is the actual cause.

## Related Issues

This same issue affects:

- ✅ `join_event` - Fixed with diagnostic logging
- ✅ `claim_certificate` - Already has diagnostic logging

Both use cases require the system transactor to be registered in `DecmAccessManager`.

## Prevention

To prevent this in the future:

1. **Document the setup step**: Add to deployment scripts/documentation
2. **Automate registration**: Include it in contract deployment scripts
3. **Add health checks**: Periodically verify transactor is registered
4. **Environment validation**: Check on startup if transactor is registered
