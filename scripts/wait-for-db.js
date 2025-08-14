#!/usr/bin/env node

/**
 * Wait for PostgreSQL database to be ready
 * This script polls the database until it's available or times out
 */

const { spawn } = require("child_process");
const fs = require("fs");
const { getDatabaseConfig, displayConfig } = require("./db-env");

const MAX_ATTEMPTS = 30;
const DELAY_MS = 2000;

console.log("🔍 Waiting for PostgreSQL database to be ready...");

// Load database configuration from .env
const { config, databaseUrl, isFromEnvFile } = getDatabaseConfig();
displayConfig(config, isFromEnvFile);
const DATABASE_URL = databaseUrl;

async function sleep(ms) {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

async function checkDatabase() {
	return new Promise((resolve) => {
		const process = spawn("psql", [DATABASE_URL, "-c", "SELECT 1;"], {
			stdio: ["ignore", "ignore", "ignore"],
		});

		process.on("close", (code) => {
			resolve(code === 0);
		});

		process.on("error", () => {
			resolve(false);
		});
	});
}

async function checkDockerContainer() {
	return new Promise((resolve) => {
		const process = spawn(
			"docker",
			["ps", "--filter", "name=decm-postgres", "--format", "{{.Status}}"],
			{
				stdio: ["ignore", "pipe", "ignore"],
			}
		);

		let output = "";
		process.stdout.on("data", (data) => {
			output += data.toString();
		});

		process.on("close", (code) => {
			const isHealthy = output.includes("healthy") || output.includes("Up");
			resolve(code === 0 && isHealthy);
		});

		process.on("error", () => {
			resolve(false);
		});
	});
}

async function waitForDatabase() {
	console.log(
		`⏳ Polling database every ${DELAY_MS / 1000}s (max ${MAX_ATTEMPTS} attempts)...`
	);

	for (let attempt = 1; attempt <= MAX_ATTEMPTS; attempt++) {
		console.log(`   Attempt ${attempt}/${MAX_ATTEMPTS}...`);

		// First check if Docker container is healthy
		const containerHealthy = await checkDockerContainer();
		if (!containerHealthy) {
			console.log("   🐳 Container not ready yet...");
			await sleep(DELAY_MS);
			continue;
		}

		// Then check if database accepts connections
		const dbReady = await checkDatabase();
		if (dbReady) {
			console.log("✅ Database is ready!");
			return;
		}

		console.log("   🔄 Database not ready yet...");
		await sleep(DELAY_MS);
	}

	console.error("❌ Database failed to become ready within timeout");
	console.error("💡 Try running: bun db:stop && bun db:start");
	process.exit(1);
}

// Handle Ctrl+C gracefully
process.on("SIGINT", () => {
	console.log("\n⚠️  Database wait cancelled");
	process.exit(130);
});

waitForDatabase().catch((error) => {
	console.error("❌ Error waiting for database:", error.message);
	process.exit(1);
});
