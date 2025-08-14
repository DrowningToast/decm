#!/usr/bin/env node

/**
 * Database environment configuration helper
 * Reads database credentials from .env file and constructs DATABASE_URL
 */

const fs = require("fs");
const path = require("path");

function loadEnvFile() {
	const envPath = path.join(__dirname, "..", ".env");
	const envExamplePath = path.join(__dirname, "..", ".env.example");

	let envContent = "";

	// Try to load .env file first, fallback to .env.example
	if (fs.existsSync(envPath)) {
		envContent = fs.readFileSync(envPath, "utf8");
		console.log("📄 Using .env file for database configuration");
	} else if (fs.existsSync(envExamplePath)) {
		envContent = fs.readFileSync(envExamplePath, "utf8");
		console.log("⚠️  .env file not found, using .env.example defaults");
		console.log(
			"💡 Copy .env.example to .env and configure your database credentials"
		);
	} else {
		throw new Error("Neither .env nor .env.example file found");
	}

	return envContent;
}

function parseEnvFile(content) {
	const env = {};
	const lines = content.split("\n");

	for (const line of lines) {
		const trimmed = line.trim();
		if (trimmed && !trimmed.startsWith("#")) {
			const [key, ...valueParts] = trimmed.split("=");
			if (key && valueParts.length > 0) {
				// Join back in case value contains '=' characters
				let value = valueParts.join("=").trim();

				// Remove surrounding quotes if present
				if (
					(value.startsWith('"') && value.endsWith('"')) ||
					(value.startsWith("'") && value.endsWith("'"))
				) {
					value = value.slice(1, -1);
				}

				env[key.trim()] = value;
			}
		}
	}

	return env;
}

function getDatabaseConfig() {
	try {
		const envContent = loadEnvFile();
		const env = parseEnvFile(envContent);

		// Use environment variables if available, otherwise use .env file values
		const config = {
			host: process.env.DB_HOST || env.DB_HOST || "localhost",
			port: process.env.DB_PORT || env.DB_PORT || "5432",
			user: process.env.DB_USER || env.DB_USER || "decm_user",
			password: process.env.DB_PASSWORD || env.DB_PASSWORD || "decm_password",
			database: process.env.DB_NAME || env.DB_NAME || "decm",
			sslMode: process.env.DB_SSL_MODE || env.DB_SSL_MODE || "disable",
		};

		// Construct PostgreSQL connection URL
		const databaseUrl = `postgres://${config.user}:${config.password}@${config.host}:${config.port}/${config.database}?sslmode=${config.sslMode}`;

		return {
			config,
			databaseUrl,
			isFromEnvFile: !process.env.DB_HOST, // True if reading from .env file
		};
	} catch (error) {
		console.error("❌ Error loading database configuration:", error.message);
		console.error("\n💡 Make sure you have either:");
		console.error("   1. A .env file with database configuration");
		console.error("   2. DB_* environment variables set");
		console.error("   3. Or copy .env.example to .env");
		process.exit(1);
	}
}

function displayConfig(config, isFromEnvFile) {
	console.log("🔧 Database Configuration:");
	console.log(`   Host: ${config.host}`);
	console.log(`   Port: ${config.port}`);
	console.log(`   User: ${config.user}`);
	console.log(`   Database: ${config.database}`);
	console.log(`   SSL Mode: ${config.sslMode}`);
	console.log(
		`   Source: ${isFromEnvFile ? ".env file" : "environment variables"}`
	);
}

// Export for use in other scripts
module.exports = {
	getDatabaseConfig,
	displayConfig,
};

// If run directly, display the configuration
if (require.main === module) {
	const { config, databaseUrl, isFromEnvFile } = getDatabaseConfig();
	displayConfig(config, isFromEnvFile);
	console.log("\n🔗 Database URL:");
	console.log(databaseUrl);
}
