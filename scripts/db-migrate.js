#!/usr/bin/env node

/**
 * Database migration script with error handling and recovery
 * Handles common migration issues like dirty state automatically
 */

const { spawn } = require("child_process");
const path = require("path");
const { getDatabaseConfig, displayConfig } = require("./db-env");

const MIGRATIONS_PATH = "./migrations";

console.log("🗄️  Running database migrations...");

// Load database configuration from .env
const { config, databaseUrl, isFromEnvFile } = getDatabaseConfig();
displayConfig(config, isFromEnvFile);
const DATABASE_URL = databaseUrl;

function runCommand(command, args, options = {}) {
	return new Promise((resolve, reject) => {
		console.log(`📋 Running: ${command} ${args.join(" ")}`);

		const childProcess = spawn(command, args, {
			stdio: "inherit",
			env: { ...process.env, DATABASE_URL },
			cwd: path.join(__dirname, "..", "packages", "database"),
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

async function getMigrationVersion() {
	return new Promise((resolve) => {
		const childProcess = spawn(
			"migrate",
			["-path", MIGRATIONS_PATH, "-database", DATABASE_URL, "version"],
			{
				stdio: ["ignore", "pipe", "ignore"],
				cwd: path.join(__dirname, "..", "packages", "database"),
			}
		);

		let output = "";
		childProcess.stdout.on("data", (data) => {
			output += data.toString();
		});

		childProcess.on("close", (code) => {
			if (code === 0) {
				const version = output.trim();
				resolve(version === "no migration" ? 0 : parseInt(version));
			} else {
				resolve(-1); // Error state
			}
		});
	});
}

async function forceMigrationVersion(version) {
	console.log(`🔧 Forcing migration version to ${version}...`);
	await runCommand("migrate", [
		"-path",
		MIGRATIONS_PATH,
		"-database",
		DATABASE_URL,
		"force",
		version.toString(),
	]);
}

async function runMigrations() {
	try {
		// First, try to run migrations normally
		console.log("⬆️  Attempting to run migrations...");
		await runCommand("migrate", [
			"-path",
			MIGRATIONS_PATH,
			"-database",
			DATABASE_URL,
			"up",
		]);
		console.log("✅ Migrations completed successfully!");
	} catch (error) {
		console.log("⚠️  Migration failed, checking for dirty state...");

		// Check if it's a dirty database error
		if (error.message.includes("Dirty database")) {
			console.log("🧹 Detected dirty database state, attempting to clean...");

			// Get the last successful migration version
			const version = await getMigrationVersion();

			if (version === -1) {
				console.log("🔄 Forcing to clean state...");
				await forceMigrationVersion(1); // Assume at least extensions migration worked
			} else {
				console.log(`🔄 Forcing to version ${version}...`);
				await forceMigrationVersion(version);
			}

			// Try migrations again
			console.log("⬆️  Retrying migrations after cleanup...");
			await runCommand("migrate", [
				"-path",
				MIGRATIONS_PATH,
				"-database",
				DATABASE_URL,
				"up",
			]);
			console.log("✅ Migrations completed successfully after cleanup!");
		} else {
			throw error;
		}
	}

	// Show final status
	const finalVersion = await getMigrationVersion();
	console.log(`📊 Current migration version: ${finalVersion}`);

	// Verify tables were created
	console.log("🔍 Verifying database tables...");
	await runCommand("psql", [DATABASE_URL, "-c", "\\dt"]);
}

// Handle Ctrl+C gracefully
process.on("SIGINT", () => {
	console.log("\n⚠️  Migration cancelled");
	process.exit(130);
});

runMigrations().catch((error) => {
	console.error("❌ Migration failed:", error.message);
	console.error("\n💡 Try these troubleshooting steps:");
	console.error("   pnpm db:migrate:force 0  # Reset to clean state");
	console.error("   pnpm db:migrate          # Try again");
	console.error("   pnpm db:status           # Check database tables");
	process.exit(1);
});
