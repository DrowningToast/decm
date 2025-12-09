#!/bin/bash

RPC_URL="https://eth-sepolia.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5"
HOST_ADDRESS="0x7836f1b8B0FDf5Fb86A7617eF167EbeC23aa4e8E"
TRANSACTOR="0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A"

# IMPORTANT: Replace this with your actual EventCertificate contract address
# This is NOT the Event contract - it's the EventCertificate contract
# Find it in your database: certificate.EventCertificateAddress field
CERT_CONTRACT="$1"

if [ -z "$CERT_CONTRACT" ]; then
  echo "Usage: $0 <EventCertificateContractAddress>"
  echo ""
  echo "Example:"
  echo "  $0 0xYourCertificateContractAddressHere"
  echo ""
  echo "To find your certificate contract address:"
  echo "  1. Check your database: SELECT event_certificate_address FROM event_certificates WHERE id = '<your-certificate-id>';"
  echo "  2. Or check the logs: Look for 'contract_address' in claim_certificate.go logs"
  exit 1
fi

echo "Checking access control for Certificate Contract: $CERT_CONTRACT"
echo ""

# Get EventAccessManager
EAM=$(cast call $CERT_CONTRACT "EVENT_ACCESS_MANAGER()(address)" --rpc-url $RPC_URL)
echo "EventAccessManager: $EAM"

# Check host
IS_HOST=$(cast call $EAM "checkIsHostOrAdmin(address)(bool)" $HOST_ADDRESS --rpc-url $RPC_URL)
echo "Host ($HOST_ADDRESS) is host/admin: $IS_HOST"

# Get DecmAccessManager
DAM=$(cast call $EAM "DECM_ACCESS_MANAGER()(address)" --rpc-url $RPC_URL)
echo "DecmAccessManager: $DAM"

# Check transactor
IS_ALLOWED=$(cast call $DAM "checkIsAllowedMsgSender(address)(bool)" $TRANSACTOR --rpc-url $RPC_URL)
echo "Transactor ($TRANSACTOR) is allowed: $IS_ALLOWED"

echo ""
if [ "$IS_HOST" == "false" ] && [ "$IS_ALLOWED" == "false" ]; then
  echo "❌ ACCESS DENIED - Transaction will revert"
else
  echo "✅ ACCESS GRANTED"
fi
