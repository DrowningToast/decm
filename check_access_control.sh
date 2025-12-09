#!/bin/bash

# Your environment variables
RPC_URL="https://eth-sepolia.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5"

# Addresses - UPDATE THIS with your actual EventCertificate contract address
# Note: This is the CERTIFICATE contract (EventCertificate.sol), NOT the Event contract
# The address should be from certificate.EventCertificateAddress in your database
CERT_CONTRACT="0x1cb18A232e8955B77Ffb87771D343e7a5cCa8769"  # TODO: Replace with actual EventCertificate contract address
HOST_ADDRESS="0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E"
TRANSACTOR="0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A"

echo "=========================================="
echo "Checking Access Control for Certificate Mint"
echo "=========================================="
echo ""

# Step 1: Get EventAccessManager address from Certificate contract
echo "1. Getting EventAccessManager address..."
EAM_ADDRESS=$(cast call $CERT_CONTRACT \
  "EVENT_ACCESS_MANAGER()(address)" \
  --rpc-url $RPC_URL)

if [ -z "$EAM_ADDRESS" ] || [ "$EAM_ADDRESS" == "0x0000000000000000000000000000000000000000" ]; then
  echo "❌ Failed to get EventAccessManager address"
  exit 1
fi

echo "✅ EventAccessManager: $EAM_ADDRESS"
echo ""

# Step 2: Check if host is registered as host/admin
echo "2. Checking if host is registered..."
IS_HOST_OR_ADMIN=$(cast call $EAM_ADDRESS \
  "checkIsHostOrAdmin(address)(bool)" \
  $HOST_ADDRESS \
  --rpc-url $RPC_URL)

if [ "$IS_HOST_OR_ADMIN" == "true" ]; then
  echo "✅ Host IS registered as host/admin: true"
else
  echo "❌ Host is NOT registered as host/admin: false"
fi
echo ""

# Step 3: Get DecmAccessManager address
echo "3. Getting DecmAccessManager address..."
DAM_ADDRESS=$(cast call $EAM_ADDRESS \
  "DECM_ACCESS_MANAGER()(address)" \
  --rpc-url $RPC_URL)

if [ -z "$DAM_ADDRESS" ] || [ "$DAM_ADDRESS" == "0x0000000000000000000000000000000000000000" ]; then
  echo "❌ Failed to get DecmAccessManager address"
  exit 1
fi

echo "✅ DecmAccessManager: $DAM_ADDRESS"
echo ""

# Step 4: Check if transactor is allowed sender
echo "4. Checking if system transactor is allowed..."
IS_ALLOWED_SENDER=$(cast call $DAM_ADDRESS \
  "checkIsAllowedMsgSender(address)(bool)" \
  $TRANSACTOR \
  --rpc-url $RPC_URL)

if [ "$IS_ALLOWED_SENDER" == "true" ]; then
  echo "✅ Transactor IS allowed sender: true"
else
  echo "❌ Transactor is NOT allowed sender: false"
fi
echo ""

# Step 5: Summary
echo "=========================================="
echo "SUMMARY"
echo "=========================================="
echo "Host address: $HOST_ADDRESS"
echo "  └─ Is host/admin: $IS_HOST_OR_ADMIN"
echo ""
echo "Transactor address: $TRANSACTOR"
echo "  └─ Is allowed sender: $IS_ALLOWED_SENDER"
echo ""

if [ "$IS_HOST_OR_ADMIN" == "false" ] && [ "$IS_ALLOWED_SENDER" == "false" ]; then
  echo "❌ ACCESS DENIED: Neither condition is met!"
  echo "   The transaction will revert with: 'Not host or admin or allowed msg sender'"
  echo ""
  echo "FIX: You need to either:"
  echo "  1. Register host as host/admin in EventAccessManager, OR"
  echo "  2. Register transactor as allowed sender in DecmAccessManager"
else
  echo "✅ ACCESS GRANTED: At least one condition is met"
fi
