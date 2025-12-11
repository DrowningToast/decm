/**
 * Diagnose why AddParticipant is failing
 * Run with: node diagnose-join-failure.js
 */

const { ethers } = require("ethers");

const POLYGON_RPC = "https://polygon-mainnet.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5";
const EVENT_CONTRACT = "0xF214925eF0874fcE9F9EDaf3F7E2992E210cA97a";
const PARTICIPANT_ADDRESS = "0x3F653E90D907581F1379474C815C1c4f7135F836";
const MESSAGE_HASH = "0xf576c365a08aa6b1686cf46a6383f3bb065274d8bd9be9db29a5f6ab44b630e1";
const SIGNATURE =
    "0x5f61be9523911b1e0d9dd791d3cfb416ff7e0573d7b48d23185d31959d9e421a3e588af9c9e2047d673327362b852adab5dbe36ede4664f10de6fd06946452181b";

async function diagnoseJoinFailure() {
    console.log("🔍 Diagnosing AddParticipant Transaction Failure\n");
    console.log("=".repeat(80));

    const provider = new ethers.JsonRpcProvider(POLYGON_RPC);

    // Event contract ABI (minimal)
    const eventAbi = [
        "function currentSeatsCount() view returns (uint256)",
        "function seatsCount() view returns (uint256)",
        "function getParticipants() view returns (address[])",
        "function EVENT_ACCESS_MANAGER() view returns (address)",
    ];

    const eventContract = new ethers.Contract(EVENT_CONTRACT, eventAbi, provider);

    // Check 1: Event capacity
    console.log("\n📋 Check 1: Event Capacity");
    const currentSeats = await eventContract.currentSeatsCount();
    const maxSeats = await eventContract.seatsCount();

    console.log("Current seats:", currentSeats.toString());
    console.log("Max seats:", maxSeats.toString());
    console.log("Event full:", currentSeats >= maxSeats ? "❌ YES - EVENT IS FULL!" : "✅ No");

    if (currentSeats >= maxSeats) {
        console.log("\n🔴 PROBLEM FOUND: Event has reached maximum capacity!");
        console.log('   This will cause: "Seats count reached" revert');
        return;
    }

    // Check 2: Get all participants and check if already joined
    console.log("\n📋 Check 2: Current Participants List");
    const participants = await eventContract.getParticipants();
    console.log("Total participants:", participants.length);

    // Check 3: Participant already joined
    console.log("\n📋 Check 3: Participant Already Joined");
    console.log("Participant:", PARTICIPANT_ADDRESS);

    const alreadyJoined = participants.some(
        (addr) => addr.toLowerCase() === PARTICIPANT_ADDRESS.toLowerCase(),
    );
    console.log("Already joined:", alreadyJoined ? "❌ YES - ALREADY A PARTICIPANT!" : "✅ No");

    if (alreadyJoined) {
        console.log("\n🔴 PROBLEM FOUND: Participant already joined this event!");
        console.log('   This will cause: "Participant is already joined" revert');
        return;
    }

    console.log("Total participants:", participants.length);

    if (participants.length > 0) {
        console.log("First 5 participants:");
        participants.slice(0, 5).forEach((addr, i) => {
            console.log(`   ${i + 1}. ${addr}`);
        });
    }

    // Check 4: Signature analysis
    console.log("\n📋 Check 4: Signature Analysis");
    console.log("Message Hash:", MESSAGE_HASH);
    console.log("Signature:", SIGNATURE);
    console.log("Signature length:", SIGNATURE.length, "chars (130 expected for 65 bytes)");

    // Extract v, r, s from signature
    const sig = SIGNATURE.slice(2); // Remove 0x
    const r = "0x" + sig.substring(0, 64);
    const s = "0x" + sig.substring(64, 128);
    const v = parseInt(sig.substring(128, 130), 16);

    console.log("Signature components:");
    console.log("  r:", r);
    console.log("  s:", s);
    console.log("  v:", v, v === 27 || v === 28 ? "✅ Valid (27 or 28)" : "❌ Invalid!");

    if (v !== 27 && v !== 28) {
        console.log("\n🔴 PROBLEM FOUND: Signature has invalid recovery ID (v)");
        console.log("   Expected: 27 or 28");
        console.log("   Got:", v);
        console.log("   This will cause signature verification to fail!");
        return;
    }

    // Check 5: Try to recover signer from signature
    console.log("\n📋 Check 5: Signature Recovery Test");

    try {
        // We need the original message, not the hash
        // The contract does: MessageHashUtils.toEthSignedMessageHash(bytes(signedMessageDigest))
        // So we need to figure out what the original message was

        console.log("⚠️  Cannot fully verify without the original sign message (not hash)");
        console.log("   The contract expects the RAW message string, not the hash");
        console.log(
            '   Example: "0x3F653E90D907581F1379474C815C1c4f7135F836,0xF214925eF0874fcE9F9EDaf3F7E2992E210cA97a,12345678"',
        );
    } catch (error) {
        console.log("❌ Error during signature recovery:", error.message);
    }

    // Check 6: ThemisUtils signature tracking
    console.log("\n📋 Check 6: Signature Already Used Check");
    console.log("⚠️  ThemisUtils tracks used signatures in a mapping");
    console.log("   If this signature was used in a previous transaction attempt,");
    console.log('   it will revert with "Signature already used"');
    console.log("\n   To check this, we would need to:");
    console.log("   1. Get EventAccessManager address from Event contract");
    console.log("   2. Call usedSignatures(signature) on ThemisUtils");

    // Get EventAccessManager
    const eventAccessManagerAddr = await eventContract.EVENT_ACCESS_MANAGER();
    console.log("\n   EventAccessManager address:", eventAccessManagerAddr);

    // Try to check if signature was used
    const themisAbi = ["function usedSignatures(bytes) view returns (bool)"];

    try {
        const themisUtils = new ethers.Contract(eventAccessManagerAddr, themisAbi, provider);
        const signatureBytes = ethers.getBytes(SIGNATURE);
        const wasUsed = await themisUtils.usedSignatures(signatureBytes);

        console.log(
            "   Signature already used:",
            wasUsed ? "❌ YES - SIGNATURE CONSUMED!" : "✅ No",
        );

        if (wasUsed) {
            console.log("\n🔴 PROBLEM FOUND: Signature has already been used!");
            console.log('   This will cause: "Signature already used" revert');
            console.log("\n💡 Solution: Generate a new signature with a fresh deadline block");
            return;
        }
    } catch (error) {
        console.log(
            "   ⚠️  Could not check signature usage (method may not exist on this contract)",
        );
    }

    console.log("\n" + "=".repeat(80));
    console.log("\n📊 Summary:");
    console.log("✅ Event not full");
    console.log("✅ Participant not already joined");
    console.log("✅ Signature format looks valid");
    console.log("\n🤔 If all checks pass, the issue might be:");
    console.log("   1. Message format mismatch (backend sends hash, contract expects raw message)");
    console.log("   2. Signature doesn't match the message");
    console.log("   3. Race condition (state changed between backend check and transaction)");
    console.log("\n💡 Next step: Review backend code - ensure passing RAW sign message, not hash!");
}

diagnoseJoinFailure().catch(console.error);
