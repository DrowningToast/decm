#!/usr/bin/env node

/**
 * Script to run GitHub Actions workflows locally using act
 * Usage: node scripts/act-run.js [job-name] [options]
 */

const { spawn } = require("child_process");
const { execSync } = require("child_process");

// Colors for output
const colors = {
    reset: "\x1b[0m",
    red: "\x1b[31m",
    green: "\x1b[32m",
    yellow: "\x1b[33m",
};

const log = {
    error: (msg) => console.error(`${colors.red}${msg}${colors.reset}`),
    success: (msg) => console.log(`${colors.green}${msg}${colors.reset}`),
    warning: (msg) => console.log(`${colors.yellow}${msg}${colors.reset}`),
};

// Available jobs
const JOBS = [
    "docker-compose",
    "database-migrate",
    "database-generate",
    "backend-build",
    "backend-tests",
    "backend-openapi",
    "frontend-build",
    "frontend-tests",
    "contracts-build",
];

// Default event (simulate pull_request)
const EVENT = "pull_request";
const SECRET_FILE = ".env.test";

// Check if act is installed
function checkActInstalled() {
    try {
        execSync("which act", { stdio: "ignore" });
        return true;
    } catch (error) {
        log.error("❌ act is not installed");
        log.warning(
            "Install it with: brew install act (macOS) or visit https://github.com/nektos/act",
        );
        process.exit(1);
    }
}

// List available jobs
function listJobs() {
    log.success("Available jobs:");
    JOBS.forEach((job) => {
        console.log(`  - ${job}`);
    });
    console.log("");
    log.warning("Usage:");
    console.log("  node scripts/act-run.js [job-name] [options]");
    console.log("");
    log.warning("Examples:");
    console.log("  node scripts/act-run.js                    # Run all jobs");
    console.log("  node scripts/act-run.js docker-compose      # Run specific job");
    console.log("  node scripts/act-run.js backend-build -l    # List steps");
    console.log("  node scripts/act-run.js frontend-build -v  # Verbose output");
}

// Run act command
function runAct(jobName, extraArgs) {
    const args = [EVENT];

    // Explicitly specify workflow directory to avoid scanning other YAML files
    args.push("-W", ".github/workflows/");

    if (jobName) {
        args.push("--job", jobName);
    }

    args.push("--secret-file", SECRET_FILE);

    // Add any extra arguments
    if (extraArgs && extraArgs.length > 0) {
        args.push(...extraArgs);
    }

    return new Promise((resolve, reject) => {
        log.success(`🚀 Running act for job: ${jobName || "all"}`);
        console.log("");

        const actProcess = spawn("act", args, {
            stdio: "inherit",
            env: process.env,
        });

        actProcess.on("close", (code) => {
            if (code === 0) {
                resolve();
            } else {
                reject(new Error(`act exited with code ${code}`));
            }
        });

        actProcess.on("error", (error) => {
            reject(error);
        });
    });
}

// Main function
function main() {
    const args = process.argv.slice(2);
    const jobName = args[0];
    const extraArgs = args.slice(1);

    // Check if act is installed
    checkActInstalled();

    // Show help if requested or no job specified
    if (!jobName || jobName === "--help" || jobName === "-h" || jobName === "help") {
        listJobs();
        process.exit(0);
    }

    // Validate job name if provided
    if (jobName && !JOBS.includes(jobName)) {
        log.error(`❌ Unknown job: ${jobName}`);
        console.log("");
        listJobs();
        process.exit(1);
    }

    // Run act
    runAct(jobName, extraArgs).catch((error) => {
        log.error(`❌ Error running act: ${error.message}`);
        process.exit(1);
    });
}

// Run if called directly
if (require.main === module) {
    main();
}

module.exports = { runAct, listJobs, JOBS };
