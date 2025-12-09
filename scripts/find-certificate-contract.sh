#!/bin/bash

echo "🔍 Finding Certificate Contract Address..."
echo ""

# From the certificate ID in your logs
CERT_ID="158cc7cd-0451-4fa6-bdba-a85a1801a2df"

cd /Users/supratouchsuwatno/Desktop/decm/apps/backend

# Query the database
echo "Querying database for certificate contract address..."
echo ""

# You'll need to run this SQL query
echo "Please run this SQL query:"
echo ""
echo "SELECT"
echo "  ec.id as certificate_id,"
echo "  ec.event_certificate_address,"
echo "  e.id as event_id,"
echo "  e.name as event_name"
echo "FROM event_certificates ec"
echo "JOIN events e ON ec.event_id = e.id"
echo "WHERE ec.id = '$CERT_ID';"
echo ""
echo "Or use: pnpm db:console"
