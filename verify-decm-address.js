/**
 * Comprehensive signature diagnosis
 * This will try different sign message formats to find the issue
 */

const { ethers } = require("ethers");

const PARTICIPANT_ADDRESS = "0x3F653E90D907581F1379474C815C1c4f7135F836";
const CONTRACT_ADDRESS = "0x73E452e141cf206216708a02257c530Eba8e40dA";
const SIGNATURE =
    "0x2bac07b22f1f738879e708e7810df35f641bf80c738d21bb28c8b47f6f452af44498fcf965e46354013e1a842aaa2bb8bca70d3bed4a0c1d3bb76b2944da89941c";
const MESSAGE_HASH = "0x7cc9064cd2d6211535106f9b558375e55f5f37e1610665e18d8418f266b77737";

console.log("🔍 Comprehensive Signature Diagnosis\n");
console.log("=".repeat(80));
console.log("Participant:", PARTICIPANT_ADDRESS);
console.log("Contract:", CONTRACT_ADDRESS);
console.log("Signature:", SIGNATURE);
console.log("Message Hash:", MESSAGE_HASH);
console.log("=".repeat(80));

async function tryDirectRecovery() {
    console.log("\n📋 Test: Direct recovery from signature");

    try {
        const recoveredFromHash = ethers.recoverAddress(MESSAGE_HASH, SIGNATURE);
        console.log("Recovered signer from message hash:", recoveredFromHash);

        if (recoveredFromHash.toLowerCase() === PARTICIPANT_ADDRESS.toLowerCase()) {
            console.log("✅ Signer matches participant!");
        } else {
            console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
            console.log("🔴 ROOT CAUSE FOUND!");
            console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
            console.log("❌ Signer does NOT match participant!");
            console.log("   Expected signer:", PARTICIPANT_ADDRESS);
            console.log("   Actual signer:", recoveredFromHash);
            console.log("\n💡 This means:");
            console.log("   The signature was signed by private key for:", recoveredFromHash);
            console.log("   But trying to add participant:", PARTICIPANT_ADDRESS);
            console.log("\n🔍 POSSIBLE CAUSES:");
            console.log("   1. User entered wrong password when joining");
            console.log("   2. Database has wrong encrypted_private_key for this user");
            console.log("   3. Sign message was created with wrong participant address");
            console.log("\n🛠️  NEXT STEPS:");
            console.log("   1. Query database for wallet", recoveredFromHash);
            console.log("   2. Check if this wallet belongs to a different user");
            console.log("   3. Verify the password decryption logic");
            console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
        }
    } catch (error) {
        console.log("❌ Failed to recover signer:", error.message);
    }
}

tryDirectRecovery().catch(console.error);
