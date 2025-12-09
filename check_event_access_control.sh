#!/bin/bash

# Diagnostic script to check access control for Event contract's join_event functionality
# Usage: ./check_event_access_control.sh <EVENT_CONTRACT_ADDRESS>

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if EVENT_CONTRACT_ADDRESS is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Event contract address is required${NC}"
    echo "Usage: ./check_event_access_control.sh <EVENT_CONTRACT_ADDRESS>"
    exit 1
fi

EVENT_CONTRACT="$1"
RPC_URL="${BLOCKCHAIN_RPC_URL:-https://eth-sepolia.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5}"
SYSTEM_TRANSACTOR="0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A"

echo "Checking access control for Event Contract: $EVENT_CONTRACT"
echo "System Transactor: $SYSTEM_TRANSACTOR"
echo ""

# Get EventAccessManager address from Event contract
echo -e "${YELLOW}1. Getting EventAccessManager address...${NC}"
EVENT_ACCESS_MANAGER=$(cast call "$EVENT_CONTRACT" "EVENT_ACCESS_MANAGER()(address)" --rpc-url "$RPC_URL" 2>/dev/null || echo "")
if [ -z "$EVENT_ACCESS_MANAGER" ]; then
    echo -e "${RED}Error: Could not get EventAccessManager address${NC}"
    exit 1
fi
echo -e "${GREEN}EventAccessManager: $EVENT_ACCESS_MANAGER${NC}"
echo ""

# Get DecmAccessManager address from EventAccessManager
echo -e "${YELLOW}2. Getting DecmAccessManager address...${NC}"
DECM_ACCESS_MANAGER=$(cast call "$EVENT_ACCESS_MANAGER" "DECM_ACCESS_MANAGER()(address)" --rpc-url "$RPC_URL" 2>/dev/null || echo "")
if [ -z "$DECM_ACCESS_MANAGER" ]; then
    echo -e "${RED}Error: Could not get DecmAccessManager address${NC}"
    exit 1
fi
echo -e "${GREEN}DecmAccessManager: $DECM_ACCESS_MANAGER${NC}"
echo ""

# Check if system transactor is allowed msg sender
echo -e "${YELLOW}3. Checking if system transactor is allowed msg sender...${NC}"
IS_ALLOWED=$(cast call "$DECM_ACCESS_MANAGER" "checkIsAllowedMsgSender(address)(bool)" "$SYSTEM_TRANSACTOR" --rpc-url "$RPC_URL" 2>/dev/null || echo "false")

if [ "$IS_ALLOWED" = "true" ]; then
    echo -e "${GREEN}✅ System transactor ($SYSTEM_TRANSACTOR) IS allowed msg sender${NC}"
else
    echo -e "${RED}❌ System transactor ($SYSTEM_TRANSACTOR) is NOT allowed msg sender${NC}"
    echo ""
    echo -e "${YELLOW}To fix this, call allowMsgSender on DecmAccessManager:${NC}"
    echo "cast send $DECM_ACCESS_MANAGER \"allowMsgSender(address)\" $SYSTEM_TRANSACTOR --rpc-url $RPC_URL --private-key \$BLOCKCHAIN_PRIVATE_KEY"
fi
echo ""

# Summary
echo -e "${YELLOW}=== Summary ===${NC}"
echo "Event Contract: $EVENT_CONTRACT"
echo "EventAccessManager: $EVENT_ACCESS_MANAGER"
echo "DecmAccessManager: $DECM_ACCESS_MANAGER"
echo "System Transactor: $SYSTEM_TRANSACTOR"
echo "Is Allowed: $IS_ALLOWED"
echo ""

if [ "$IS_ALLOWED" = "true" ]; then
    echo -e "${GREEN}✅ ACCESS GRANTED - join_event should work${NC}"
else
    echo -e "${RED}❌ ACCESS DENIED - join_event will revert${NC}"
    echo -e "${YELLOW}Reason: Event._addParticipant calls grantParticipantRoleUsingAllowedMsgSender,${NC}"
    echo -e "${YELLOW}which requires the transactor (msg.sender) to be an allowed msg sender.${NC}"
fi
