#!/usr/bin/env node

/**
 * Database command runner with dynamic .env configuration
 * Runs migrate or psql commands with credentials from .env file
 */

const { spawn } = require("child_process");
const path = require("path");
const { getDatabaseConfig, displayConfig } = require("./db-env");

function printUsage() {
    console.log("Usage:");
    console.log("  node db-command.js migrate <command> [args...]");
    console.log("  node db-command.js psql [sql_command]");
    console.log("");
    console.log("Examples:");
    console.log("  node db-command.js migrate up");
    console.log("  node db-command.js migrate down");
    console.log("  node db-command.js migrate force 1");
    console.log("  node db-command.js migrate version");
    console.log('  node db-command.js psql "\\dt"');
    console.log("  node db-command.js psql");
}

async function runCommand(command, args, options = {}) {
    return new Promise((resolve, reject) => {
        console.log(`📋 Running: ${command} ${args.join(" ")}`);

        const childProcess = spawn(command, args, {
            stdio: "inherit",
            ...options,
        });

        childProcess.on("close", (code) => {
            if (code === 0) {
                resolve();
            } else {
                reject(new Error(`Command failed with exit code ${code}`));
            }
        });

        childProcess.on("error", (error) => {
            reject(error);
        });
    });
}

async function main() {
    const args = process.argv.slice(2);

    if (args.length < 1) {
        printUsage();
        process.exit(1);
    }

    // Load database configuration
    const { config, databaseUrl, isFromEnvFile } = getDatabaseConfig();
    displayConfig(config, isFromEnvFile);
    console.log(""); // Add spacing

    const tool = args[0];
    const toolArgs = args.slice(1);

    try {
        if (tool === "migrate") {
            if (toolArgs.length < 1) {
                console.error("❌ Migrate command required (up, down, force, version)");
                process.exit(1);
            }

            const migrateArgs = ["-path", "./migrations", "-database", databaseUrl, ...toolArgs];

            await runCommand("migrate", migrateArgs, {
                cwd: path.join(__dirname, "..", "packages", "database"),
                env: { ...process.env, DATABASE_URL: databaseUrl },
            });
        } else if (tool === "psql") {
            const psqlArgs = toolArgs.length > 0 ? [databaseUrl, "-c", ...toolArgs] : [databaseUrl];

            await runCommand("psql", psqlArgs, {
                env: { ...process.env, DATABASE_URL: databaseUrl },
            });
        } else {
            console.error(`❌ Unknown tool: ${tool}`);
            printUsage();
            process.exit(1);
        }

        console.log("✅ Command completed successfully!");
    } catch (error) {
        console.error("❌ Command failed:", error.message);
        console.error("\n💡 Troubleshooting:");
        console.error("   - Check your .env file configuration");
        console.error("   - Ensure database is running: pnpm db:start");
        console.error("   - Verify database credentials are correct");
        process.exit(1);
    }
}

// Handle Ctrl+C gracefully
process.on("SIGINT", () => {
    console.log("\n⚠️  Command cancelled");
    process.exit(130);
});

main().catch((error) => {
    console.error("❌ Unexpected error:", error);
    process.exit(1);
});
