/**
 * Verify signature recovery to find the mismatch
 * Run with: node verify-signature.js "RAW_SIGN_MESSAGE_HERE"
 */

const { ethers } = require("ethers");

const PARTICIPANT_ADDRESS = "0x3F653E90D907581F1379474C815C1c4f7135F836";
const CONTRACT_ADDRESS = "0x73E452e141cf206216708a02257c530Eba8e40dA";
const SIGNATURE =
    "0x2bac07b22f1f738879e708e7810df35f641bf80c738d21bb28c8b47f6f452af44498fcf965e46354013e1a842aaa2bb8bca70d3bed4a0c1d3bb76b2944da89941c";

// Get RAW sign message from command line argument
const rawSignMessage = process.argv[2];

if (!rawSignMessage) {
    console.error("❌ ERROR: Please provide the RAW sign message from backend logs");
    console.error("\nUsage:");
    console.error(
        '  node verify-signature.js "0x3F653E90D907581F1379474C815C1c4f7135F836,0x73E452e141cf206216708a02257c530Eba8e40dA,12345678"',
    );
    console.error("\nThe RAW sign message should be printed in backend logs as:");
    console.error('  "Sign Message (RAW): ..."');
    process.exit(1);
}

console.log("🔍 Signature Verification Test\n");
console.log("=".repeat(80));
console.log("Expected Participant:", PARTICIPANT_ADDRESS);
console.log("Contract:", CONTRACT_ADDRESS);
console.log("Signature:", SIGNATURE);
console.log("RAW Sign Message:", rawSignMessage);
console.log("=".repeat(80));

// Parse sign message format: "participantAddr,contractAddr,deadlineBlock"
const parts = rawSignMessage.split(",");
if (parts.length !== 3) {
    console.error("\n❌ Invalid sign message format!");
    console.error('Expected format: "participantAddr,contractAddr,deadlineBlock"');
    console.error("Got:", rawSignMessage);
    process.exit(1);
}

console.log("\n📋 Sign Message Components:");
console.log("  Participant Address:", parts[0]);
console.log("  Contract Address:", parts[1]);
console.log("  Deadline Block:", parts[2]);

// Check if components match expected values
let hasError = false;

if (parts[0].toLowerCase() !== PARTICIPANT_ADDRESS.toLowerCase()) {
    console.log("  ⚠️  WARNING: Participant address mismatch!");
    console.log("     Expected:", PARTICIPANT_ADDRESS);
    console.log("     In message:", parts[0]);
    hasError = true;
}

if (parts[1].toLowerCase() !== CONTRACT_ADDRESS.toLowerCase()) {
    console.log("  ⚠️  WARNING: Contract address mismatch!");
    console.log("     Expected:", CONTRACT_ADDRESS);
    console.log("     In message:", parts[1]);
    hasError = true;
}

// Recreate the hash exactly as the smart contract does
console.log("\n📋 Signature Recovery:");

// Step 1: Convert message to bytes (as contract does with `bytes(signedMessageDigest)`)
const messageBytes = ethers.toUtf8Bytes(rawSignMessage);
console.log("Step 1 - Message as bytes (length):", messageBytes.length);

// Step 2: Hash with Ethereum signed message prefix (as contract does with `toEthSignedMessageHash`)
const messageHash = ethers.hashMessage(rawSignMessage);
console.log("Step 2 - Ethereum signed message hash:", messageHash);

// Step 3: Recover signer
try {
    const recoveredSigner = ethers.recoverAddress(messageHash, SIGNATURE);
    console.log("Step 3 - Recovered signer:", recoveredSigner);

    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    if (recoveredSigner === ethers.ZeroAddress) {
        console.log("🔴 PROBLEM FOUND: Recovered signer is ZERO address!");
        console.log('   This will cause "Invalid signature" revert');
        hasError = true;
    } else if (recoveredSigner.toLowerCase() === PARTICIPANT_ADDRESS.toLowerCase()) {
        console.log("✅ SUCCESS: Recovered signer matches participant!");
        console.log("   The signature is mathematically valid.");
        console.log("\n💡 If contract still reverts, check:");
        console.log("   1. Backend is passing the EXACT same rawSignMessage to contract");
        console.log("   2. No whitespace or encoding differences");
        console.log("   3. Transactor wallet has sufficient MATIC for gas");
    } else {
        console.log("🔴 PROBLEM FOUND: Recovered signer does NOT match participant!");
        console.log("   Expected:", PARTICIPANT_ADDRESS);
        console.log("   Recovered:", recoveredSigner);
        console.log("\n💡 This means:");
        console.log("   - The signature was created with a DIFFERENT private key");
        console.log("   - OR the message was different when signing");
        hasError = true;
    }
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
} catch (error) {
    console.log("🔴 ERROR during signature recovery:", error.message);
    hasError = true;
}

if (!hasError) {
    console.log("\n✅ All verification checks passed!");
    console.log("\nThe signature is valid. If it still fails on-chain, check:");
    console.log("  1. Ensure backend sends EXACT same message string to contract");
    console.log("  2. Check for character encoding issues (UTF-8 vs ASCII)");
    console.log("  3. Verify transactor has HOST or ADMIN role if required");
    console.log("  4. Confirm transactor wallet has gas funds");
}

console.log("\n" + "=".repeat(80));
