#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");

// Configuration
const CONTRACTS_DIR = path.join(__dirname, "../apps/contracts");
const OUTPUT_DIR = path.join(CONTRACTS_DIR, "out");
const MAIN_CONTRACTS = [
    "Event.sol",
    "EventAccessManager.sol",
    "EventTicket.sol",
    "EventCertificate.sol",
    "DecmAccessManager.sol",
];

console.log("🔨 Building smart contracts...\n");

try {
    // Change to contracts directory and run forge build
    process.chdir(CONTRACTS_DIR);
    console.log("Running 'forge build'...");
    execSync("forge build", { stdio: "inherit" });
    
    console.log("\n✅ Contract build completed successfully!");
    
    // Verify that the output directory contains the expected files
    console.log("\n🔍 Verifying build artifacts...");
    
    let allContractsFound = true;
    MAIN_CONTRACTS.forEach((contractFile) => {
        const contractName = contractFile.replace(".sol", "");
        const contractDir = path.join(OUTPUT_DIR, contractFile);
        
        if (!fs.existsSync(contractDir)) {
            console.log(`❌ Directory not found for ${contractName} at ${contractDir}`);
            allContractsFound = false;
            return;
        }
        
        // Check for ABI and BIN files
        const abiFile = path.join(contractDir, `${contractName}.json`);
        if (!fs.existsSync(abiFile)) {
            console.log(`❌ ABI file not found for ${contractName} at ${abiFile}`);
            allContractsFound = false;
        } else {
            console.log(`✅ Found ABI for ${contractName}`);
        }
    });
    
    if (allContractsFound) {
        console.log("\n🎉 All contract artifacts verified successfully!");
        process.exit(0);
    } else {
        console.log("\n❌ Some contract artifacts are missing. Build may have failed.");
        process.exit(1);
    }
    
} catch (error) {
    console.error("\n❌ Contract build failed:");
    console.error(error.message);
    process.exit(1);
}
