#!/usr/bin/env node

import {
	readFileSync,
	writeFileSync,
	existsSync,
	mkdirSync,
	readdirSync,
	copyFileSync,
} from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const projectRoot = join(__dirname, "..");

console.log("🔄 Post-processing generated TypeScript API client...");

// Ensure src directory exists
const srcDir = join(projectRoot, "src");
if (!existsSync(srcDir)) {
	mkdirSync(srcDir, { recursive: true });
}

// Copy generated files from generated/src to our src directory
const generatedSrcDir = join(projectRoot, "generated", "src");
if (existsSync(generatedSrcDir)) {
	// Copy all TypeScript files from generated/src
	const files = readdirSync(generatedSrcDir);

	files.forEach((file) => {
		const srcPath = join(generatedSrcDir, file);
		const destPath = join(srcDir, file);

		if (file.endsWith(".ts")) {
			console.log(`📄 Copying ${file} to src/`);
			copyFileSync(srcPath, destPath);
		} else if (existsSync(srcPath) && readdirSync(srcPath).length > 0) {
			// Handle directories (like apis/)
			if (!existsSync(destPath)) {
				mkdirSync(destPath, { recursive: true });
			}

			const subFiles = readdirSync(srcPath);
			subFiles.forEach((subFile) => {
				if (subFile.endsWith(".ts")) {
					const subSrcPath = join(srcPath, subFile);
					const subDestPath = join(destPath, subFile);
					console.log(`📄 Copying ${file}/${subFile} to src/`);
					copyFileSync(subSrcPath, subDestPath);
				}
			});
		}
	});
}

// Create a barrel export file (index.ts)
const indexPath = join(srcDir, "index.ts");
let indexContent = `// DECM API Client
// This file is auto-generated from OpenAPI spec
// Do not edit manually

`;

// Look for API files in src directory
const srcFiles = readdirSync(srcDir);
const apiFiles = srcFiles
	.filter((file) => file.endsWith(".ts") && file !== "index.ts")
	.map((file) => file.replace(".ts", ""));

if (apiFiles.length > 0) {
	apiFiles.forEach((fileName) => {
		indexContent += `export * from './${fileName}';\n`;
	});
} else {
	indexContent += `// No API files generated yet
// Run 'bun gen-api' to generate TypeScript interfaces from backend OpenAPI spec
export const DECM_API_VERSION = '1.0.0';
`;
}

writeFileSync(indexPath, indexContent);

console.log("✅ Post-processing complete!");
console.log(`📦 Generated files available in src/`);
console.log(`🎯 Main export: src/index.ts`);
