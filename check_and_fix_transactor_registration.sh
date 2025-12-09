#!/bin/bash

# Script to check and fix system transactor registration in DecmAccessManager
# This fixes the "execution reverted" errors in both join_event and claim_certificate

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "System Transactor Registration Checker"
echo "=========================================="
echo ""

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

# Configuration (from environment or defaults)
RPC_URL="${BLOCKCHAIN_RPC_URL:-${RPC_URL}}"
PRIVATE_KEY="${BLOCKCHAIN_PRIVATE_KEY:-${PRIVATE_KEY}}"
DAM_ADDRESS="${DECM_ACCESS_MANAGER_ADDRESS}"

if [ -z "$RPC_URL" ]; then
  echo -e "${RED}❌ Error: RPC_URL or BLOCKCHAIN_RPC_URL not set${NC}"
  exit 1
fi

if [ -z "$PRIVATE_KEY" ]; then
  echo -e "${RED}❌ Error: PRIVATE_KEY or BLOCKCHAIN_PRIVATE_KEY not set${NC}"
  exit 1
fi

if [ -z "$DAM_ADDRESS" ]; then
  echo -e "${RED}❌ Error: DECM_ACCESS_MANAGER_ADDRESS not set${NC}"
  echo "   Please set this environment variable to the DecmAccessManager contract address"
  exit 1
fi

# Derive transactor address from private key
echo "1. Deriving system transactor address from private key..."
TRANSACTOR=$(cast wallet address --private-key $PRIVATE_KEY 2>/dev/null || cast wallet address $PRIVATE_KEY)
echo "   Transactor address: $TRANSACTOR"
echo ""

# Check if transactor is allowed
echo "2. Checking if transactor is registered in DecmAccessManager..."
echo "   DecmAccessManager: $DAM_ADDRESS"
IS_ALLOWED=$(cast call $DAM_ADDRESS \
  "allowedMsgSenders(address)(bool)" \
  $TRANSACTOR \
  --rpc-url $RPC_URL 2>/dev/null || echo "false")

if [ "$IS_ALLOWED" == "true" ]; then
  echo -e "${GREEN}✅ Transactor IS registered: true${NC}"
  echo ""
  echo "The system transactor is properly configured. If you're still seeing errors,"
  echo "check other possible causes (signature replay, seats count, etc.)"
  exit 0
else
  echo -e "${RED}❌ Transactor is NOT registered: false${NC}"
  echo ""
fi

# Get DEFAULT_ADMIN_ROLE to check if we can register
echo "3. Checking if you have permission to register the transactor..."
echo "   (Checking if the transactor address has DEFAULT_ADMIN_ROLE...)"

# Check if transactor has DEFAULT_ADMIN_ROLE
TRANSACTOR_IS_ADMIN=$(cast call $DAM_ADDRESS \
  "hasRole(bytes32,address)(bool)" \
  $(cast sig "DEFAULT_ADMIN_ROLE()") \
  $TRANSACTOR \
  --rpc-url $RPC_URL 2>/dev/null || echo "false")

# Alternative: try to get DEFAULT_ADMIN_ROLE constant
DEFAULT_ADMIN_ROLE=$(cast call $DAM_ADDRESS \
  "DEFAULT_ADMIN_ROLE()(bytes32)" \
  --rpc-url $RPC_URL 2>/dev/null || echo "")

if [ -n "$DEFAULT_ADMIN_ROLE" ] && [ "$DEFAULT_ADMIN_ROLE" != "0x0000000000000000000000000000000000000000000000000000000000000000" ]; then
  TRANSACTOR_IS_ADMIN=$(cast call $DAM_ADDRESS \
    "hasRole(bytes32,address)(bool)" \
    $DEFAULT_ADMIN_ROLE \
    $TRANSACTOR \
    --rpc-url $RPC_URL 2>/dev/null || echo "false")
fi

echo ""
echo "=========================================="
echo "REGISTRATION REQUIRED"
echo "=========================================="
echo ""
echo "The system transactor ($TRANSACTOR) is NOT registered in DecmAccessManager."
echo "This is causing 'execution reverted' errors in:"
echo "  - join_event (addParticipant)"
echo "  - claim_certificate (mintNft)"
echo ""
echo -e "${YELLOW}To fix this, you need to register the transactor using one of these methods:${NC}"
echo ""

if [ "$TRANSACTOR_IS_ADMIN" == "true" ]; then
  echo -e "${GREEN}✅ Your transactor HAS admin role - you can register it directly!${NC}"
  echo ""
  echo "Run this command:"
  echo ""
  echo "  cast send $DAM_ADDRESS \\"
  echo "    \"addAllowedMsgSender(address)\" \\"
  echo "    $TRANSACTOR \\"
  echo "    --private-key $PRIVATE_KEY \\"
  echo "    --rpc-url $RPC_URL"
  echo ""
  read -p "Do you want to register it now? (y/n) " -n 1 -r
  echo ""
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Registering transactor..."
    cast send $DAM_ADDRESS \
      "addAllowedMsgSender(address)" \
      $TRANSACTOR \
      --private-key $PRIVATE_KEY \
      --rpc-url $RPC_URL
    echo ""
    echo -e "${GREEN}✅ Transactor registered successfully!${NC}"
    
    # Verify
    IS_ALLOWED_NOW=$(cast call $DAM_ADDRESS \
      "allowedMsgSenders(address)(bool)" \
      $TRANSACTOR \
      --rpc-url $RPC_URL)
    if [ "$IS_ALLOWED_NOW" == "true" ]; then
      echo -e "${GREEN}✅ Verified: Transactor is now registered${NC}"
    fi
  fi
else
  echo -e "${YELLOW}⚠️  Your transactor does NOT have admin role${NC}"
  echo ""
  echo "You need to ask someone with DEFAULT_ADMIN_ROLE to register it:"
  echo ""
  echo "  cast send $DAM_ADDRESS \\"
  echo "    \"addAllowedMsgSender(address)\" \\"
  echo "    $TRANSACTOR \\"
  echo "    --private-key <ADMIN_PRIVATE_KEY> \\"
  echo "    --rpc-url $RPC_URL"
  echo ""
  echo "Or use Foundry/Remix to call:"
  echo "  DecmAccessManager.addAllowedMsgSender($TRANSACTOR)"
  echo ""
  echo "After registration, both join_event and claim_certificate should work."
fi

echo ""
echo "=========================================="
echo "Why This Happened"
echo "=========================================="
echo "This issue likely occurred because:"
echo "  1. DecmAccessManager was redeployed and lost its configuration"
echo "  2. The system transactor was accidentally removed"
echo "  3. A new environment was set up without registering the transactor"
echo ""
echo "The transactor needs to be registered so that:"
echo "  - join_event can call grantParticipantRoleUsingAllowedMsgSender()"
echo "  - claim_certificate can pass access control checks"
echo ""
