/**
 * Diagnose the ACTUAL production error
 * Values from the error log provided
 */

const { ethers } = require("ethers");

// ACTUAL VALUES from the production error
const POLYGON_RPC = "https://polygon-mainnet.g.alchemy.com/v2/p-1v26f8urtKJAgsL3ry5";
const EVENT_CONTRACT = "0x73E452e141cf206216708a02257c530Eba8e40dA"; // From error
const PARTICIPANT_ADDRESS = "0x3F653E90D907581F1379474C815C1c4f7135F836"; // From error
const TRANSACTOR_ADDRESS = "0x3d21c0e5391A10B9c4DB7713a6a099bDe2b71C40"; // From error
const MESSAGE_HASH = "0x7cc9064cd2d6211535106f9b558375e55f5f37e1610665e18d8418f266b77737"; // From error
const SIGNATURE =
    "0x2bac07b22f1f738879e708e7810df35f641bf80c738d21bb28c8b47f6f452af44498fcf965e46354013e1a842aaa2bb8bca70d3bed4a0c1d3bb76b2944da89941c"; // From error

async function diagnose() {
    console.log("🔍 Diagnosing ACTUAL Production Error\n");
    console.log("=".repeat(80));
    console.log("Contract:", EVENT_CONTRACT);
    console.log("Participant:", PARTICIPANT_ADDRESS);
    console.log("Transactor:", TRANSACTOR_ADDRESS);
    console.log("=".repeat(80));

    const provider = new ethers.JsonRpcProvider(POLYGON_RPC);

    // Event contract ABI
    const eventAbi = [
        "function currentSeatsCount() view returns (uint256)",
        "function seatsCount() view returns (uint256)",
        "function getParticipants() view returns (address[])",
        "function EVENT_ACCESS_MANAGER() view returns (address)",
    ];

    const eventContract = new ethers.Contract(EVENT_CONTRACT, eventAbi, provider);

    // Check 1: Event capacity
    console.log("\n📋 Check 1: Event Capacity");
    try {
        const currentSeats = await eventContract.currentSeatsCount();
        const maxSeats = await eventContract.seatsCount();
        console.log("Current seats:", currentSeats.toString());
        console.log("Max seats:", maxSeats.toString());
        console.log("Event full:", currentSeats >= maxSeats ? "❌ YES - EVENT IS FULL!" : "✅ No");

        if (currentSeats >= maxSeats) {
            console.log("\n🔴 PROBLEM FOUND: Event is full!");
            return;
        }
    } catch (error) {
        console.log("❌ Error checking capacity:", error.message);
    }

    // Check 2: Participant already joined
    console.log("\n📋 Check 2: Participant Already Joined");
    try {
        const participants = await eventContract.getParticipants();
        console.log("Total participants:", participants.length);

        const alreadyJoined = participants.some(
            (addr) => addr.toLowerCase() === PARTICIPANT_ADDRESS.toLowerCase(),
        );
        console.log("Already joined:", alreadyJoined ? "❌ YES - ALREADY JOINED!" : "✅ No");

        if (alreadyJoined) {
            console.log("\n🔴 PROBLEM FOUND: Participant already joined!");
            return;
        }
    } catch (error) {
        console.log("❌ Error checking participants:", error.message);
    }

    // Check 3: Signature format
    console.log("\n📋 Check 3: Signature Analysis");
    console.log("Signature:", SIGNATURE);
    console.log("Length:", SIGNATURE.length, "chars");

    const sig = SIGNATURE.slice(2); // Remove 0x
    const r = "0x" + sig.substring(0, 64);
    const s = "0x" + sig.substring(64, 128);
    const v = parseInt(sig.substring(128, 130), 16);

    console.log("Components:");
    console.log("  r:", r);
    console.log("  s:", s);
    console.log("  v:", v, v === 27 || v === 28 ? "✅ Valid" : "❌ INVALID!");

    if (v !== 27 && v !== 28) {
        console.log("\n🔴 PROBLEM FOUND: Invalid signature v value!");
        return;
    }

    // Check 4: Signature already used (CRITICAL!)
    console.log("\n📋 Check 4: Signature Already Used (CRITICAL)");
    try {
        const eventAccessManagerAddr = await eventContract.EVENT_ACCESS_MANAGER();
        console.log("EventAccessManager:", eventAccessManagerAddr);

        // ThemisUtils is inherited by Event contract, so we can check directly
        const themisAbi = ["function usedSignatures(bytes) view returns (bool)"];

        // Try checking on the Event contract itself (since it inherits ThemisUtils)
        const eventContractWithThemis = new ethers.Contract(EVENT_CONTRACT, themisAbi, provider);
        const signatureBytes = ethers.getBytes(SIGNATURE);
        const wasUsed = await eventContractWithThemis.usedSignatures(signatureBytes);

        console.log("Signature already used:", wasUsed ? "🔴 YES - THIS IS THE PROBLEM!" : "✅ No");

        if (wasUsed) {
            console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
            console.log("🔴 ROOT CAUSE IDENTIFIED!");
            console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
            console.log("The signature has already been consumed by a previous transaction.");
            console.log("This happens when:");
            console.log(
                "  1. A previous transaction attempt succeeded (even if frontend didn't show it)",
            );
            console.log(
                "  2. A previous transaction used the same signature and failed AFTER consuming it",
            );
            console.log("  3. Someone else submitted a transaction with this signature");
            console.log("\n💡 SOLUTION:");
            console.log("  Generate a NEW signature with a fresh deadline block.");
            console.log("  Each signature can only be used ONCE, even if the transaction fails.");
            console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
            return;
        }
    } catch (error) {
        console.log("⚠️  Could not check signature usage:", error.message);
        console.log("   The contract might not expose usedSignatures mapping");
    }

    // Check 5: Try to verify signature recovery
    console.log("\n📋 Check 5: Signature Recovery Test");
    console.log("⚠️  Need the RAW sign message to verify signature");
    console.log('   Format should be: "participantAddr,contractAddr,deadlineBlock"');
    console.log('   Example: "0x3F65...,0x73E4...,12345678"');
    console.log("\n💡 The backend should log the RAW sign message used.");

    console.log("\n" + "=".repeat(80));
    console.log("📊 SUMMARY");
    console.log("If all checks passed, the most likely causes are:");
    console.log("  1. ❌ Signature doesn't match the message (wrong signer/message mismatch)");
    console.log("  2. ❌ Contract's ECDSA.recover is failing (returns address(0))");
    console.log("  3. ❌ Race condition between checks and transaction submission");
    console.log("=".repeat(80));
}

diagnose().catch(console.error);
