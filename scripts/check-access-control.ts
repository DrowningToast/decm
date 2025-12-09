#!/usr/bin/env tsx

import { ethers } from "ethers";

// Configuration - FROM YOUR LOGS
const RPC_URL = process.env.BLOCKCHAIN_RPC_URL || "http://localhost:8545";
const RECOVERED_SIGNER = "0x9F6f8ef7c3CD068e96C6643754dbF57A49ee13aB";
const SYSTEM_TRANSACTOR = "0xf466e7cE6B06f9b3071557A790Bd45F051C1C60A";

// !!! REPLACE THIS WITH THE FULL 42-CHARACTER ADDRESS FROM DATABASE !!!
const CERTIFICATE_CONTRACT = "0x3B328a7049a374a5D728d575383d38486f1c727c";

// Certificate Contract ABI (minimal - to get EventAccessManager address)
const CERTIFICATE_ABI = ["function EVENT_ACCESS_MANAGER() external view returns (address)"];

// EventAccessManager ABI (minimal)
const ACCESS_MANAGER_ABI = [
    "function checkIsHostOrAdmin(address account) external view returns (bool)",
    "function checkIsAllowedMsgSender(address account) external view returns (bool)",
    "function getHosts() external view returns (address[])",
    "function getAllowedMsgSenders() external view returns (address[])",
];

async function checkAccessControl() {
    console.log("🔍 Checking Access Control...\n");

    const provider = new ethers.JsonRpcProvider(RPC_URL);

    // First, get EventAccessManager address from certificate contract
    console.log("📡 Fetching EventAccessManager address from certificate contract...");
    const certificateContract = new ethers.Contract(
        CERTIFICATE_CONTRACT,
        CERTIFICATE_ABI,
        provider,
    );

    const EVENT_ACCESS_MANAGER_ADDRESS = await certificateContract.EVENT_ACCESS_MANAGER();
    console.log(`✅ Found EventAccessManager at: ${EVENT_ACCESS_MANAGER_ADDRESS}\n`);

    const accessManager = new ethers.Contract(
        EVENT_ACCESS_MANAGER_ADDRESS,
        ACCESS_MANAGER_ABI,
        provider,
    );

    console.log("Configuration:");
    console.log(`  RPC: ${RPC_URL}`);
    console.log(`  Certificate Contract: ${CERTIFICATE_CONTRACT}`);
    console.log(`  EventAccessManager: ${EVENT_ACCESS_MANAGER_ADDRESS}`);
    console.log(`  Recovered Signer: ${RECOVERED_SIGNER}`);
    console.log(`  System Transactor: ${SYSTEM_TRANSACTOR}\n`);

    // Check 1: Is recovered signer a host/admin?
    console.log("━━━ CHECK 1: Is Recovered Signer a Host/Admin? ━━━");
    try {
        const isHost = await accessManager.checkIsHostOrAdmin(RECOVERED_SIGNER);
        console.log(`Result: ${isHost ? "✅ YES" : "❌ NO"}`);

        if (!isHost) {
            console.log("\n⚠️  Recovered signer is NOT registered as host/admin!");
            console.log("This is likely the issue.\n");
        }
    } catch (error) {
        console.error("❌ Error checking host status:", error);
    }

    // Check 2: Is system transactor allowed?
    console.log("\n━━━ CHECK 2: Is System Transactor Allowed? ━━━");
    try {
        const isAllowed = await accessManager.checkIsAllowedMsgSender(SYSTEM_TRANSACTOR);
        console.log(`Result: ${isAllowed ? "✅ YES" : "❌ NO"}`);

        if (!isAllowed) {
            console.log("\n⚠️  System transactor is NOT in allowed senders list!");
            console.log("This could be the issue if recovered signer is also not host/admin.\n");
        }
    } catch (error) {
        console.error("❌ Error checking allowed sender:", error);
    }

    // Check 3: List all hosts
    console.log("\n━━━ CHECK 3: All Registered Hosts ━━━");
    try {
        const hosts = await accessManager.getHosts();
        console.log(`Total hosts: ${hosts.length}`);
        hosts.forEach((host: string, i: number) => {
            const match = host.toLowerCase() === RECOVERED_SIGNER.toLowerCase();
            console.log(`  ${i + 1}. ${host} ${match ? "← MATCH!" : ""}`);
        });
    } catch (error) {
        console.error("❌ Error listing hosts:", error);
    }

    // Check 4: List all allowed senders
    console.log("\n━━━ CHECK 4: All Allowed Senders ━━━");
    try {
        const senders = await accessManager.getAllowedMsgSenders();
        console.log(`Total allowed senders: ${senders.length}`);
        senders.forEach((sender: string, i: number) => {
            const match = sender.toLowerCase() === SYSTEM_TRANSACTOR.toLowerCase();
            console.log(`  ${i + 1}. ${sender} ${match ? "← MATCH!" : ""}`);
        });
    } catch (error) {
        console.error("❌ Error listing allowed senders:", error);
    }

    // Final verdict
    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    console.log("VERDICT:");
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n");

    try {
        const isHost = await accessManager.checkIsHostOrAdmin(RECOVERED_SIGNER);
        const isAllowed = await accessManager.checkIsAllowedMsgSender(SYSTEM_TRANSACTOR);

        if (!isHost && !isAllowed) {
            console.log("❌ NEITHER condition is satisfied!");
            console.log("   - Recovered signer is NOT host/admin");
            console.log("   - System transactor is NOT allowed sender");
            console.log("\n💡 FIX: Add one of the following:");
            console.log(`   1. Register ${RECOVERED_SIGNER} as host in EventAccessManager`);
            console.log(`   2. Add ${SYSTEM_TRANSACTOR} to allowed senders in EventAccessManager`);
        } else if (isHost && isAllowed) {
            console.log("✅ BOTH conditions are satisfied - this should work!");
            console.log("   The issue might be something else.");
        } else if (isHost) {
            console.log("✅ Recovered signer IS host/admin - should work!");
        } else if (isAllowed) {
            console.log("✅ System transactor IS allowed sender - should work!");
        }
    } catch (error) {
        console.error("❌ Error determining verdict:", error);
    }
}

checkAccessControl()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error("Fatal error:", error);
        process.exit(1);
    });
