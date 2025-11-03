#!/usr/bin/env node

const fs = require("fs");
const path = require("path");

// Configuration
const OUTPUT_DIR = path.join(__dirname, "../apps/backend/contracts");

// Directories to create
const DIRECTORIES = [
    "event",
    "accessmanager", 
    "ticket",
    "certificate",
    "decm"
];

console.log("📁 Creating directory structure for contract bindings...");

// Create main contracts directory
if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
    console.log(`✅ Created main directory: ${OUTPUT_DIR}`);
}

// Create subdirectories
DIRECTORIES.forEach(dir => {
    const dirPath = path.join(OUTPUT_DIR, dir);
    if (!fs.existsSync(dirPath)) {
        fs.mkdirSync(dirPath, { recursive: true });
        console.log(`✅ Created directory: ${dirPath}`);
    } else {
        console.log(`ℹ️  Directory already exists: ${dirPath}`);
    }
});

console.log("\n🎉 Directory structure created successfully!");
