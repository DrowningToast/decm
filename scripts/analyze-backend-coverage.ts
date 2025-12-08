#!/usr/bin/env tsx

/**
 * Backend Test Coverage Analyzer
 * Analyzes Go test coverage and generates comprehensive reports
 */

import { execSync } from "child_process";
import * as fs from "fs";
import * as path from "path";

interface PackageCoverage {
    package: string;
    coverage: number;
    status: "excellent" | "good" | "medium" | "low" | "none";
}

interface UncoveredFunction {
    package: string;
    function: string;
    coverage: number;
}

interface CoverageReport {
    summary: {
        totalPackages: number;
        averageCoverage: number;
        uncoveredFunctions: number;
    };
    packages: PackageCoverage[];
    criticalGaps: {
        category: string;
        packages: PackageCoverage[];
        uncoveredFunctions: UncoveredFunction[];
    }[];
}

class BackendCoverageAnalyzer {
    private backendPath = path.join(process.cwd(), "apps", "backend");
    private coverageFile = "/tmp/decm-backend-coverage.out";

    /**
     * Run Go tests with coverage
     */
    private runTests(): void {
        console.log("🧪 Running backend tests with coverage...\n");

        try {
            execSync(
                `cd ${this.backendPath} && go test ./... -cover -coverprofile=${this.coverageFile}`,
                { stdio: "inherit" },
            );
        } catch (error) {
            console.error("⚠️  Some tests failed, but continuing with coverage analysis...\n");
        }
    }

    /**
     * Parse coverage output for package-level coverage
     */
    private parsePackageCoverage(): PackageCoverage[] {
        const output = execSync(`cd ${this.backendPath} && go test ./... -cover 2>&1`, {
            encoding: "utf-8",
        });

        const packages: PackageCoverage[] = [];
        const lines = output.split("\n");

        for (const line of lines) {
            const match = line.match(/coverage:\s+([0-9.]+)%\s+of\s+statements/);
            const packageMatch = line.match(/^(ok|PASS|FAIL)\s+([^\s]+)/);

            if (match && packageMatch) {
                const coverage = parseFloat(match[1]);
                const packageName = packageMatch[2].replace("apps/backend/", "");

                packages.push({
                    package: packageName,
                    coverage,
                    status: this.getStatus(coverage),
                });
            }
        }

        return packages.sort((a, b) => a.coverage - b.coverage);
    }

    /**
     * Parse function-level coverage from coverage profile
     */
    private parseUncoveredFunctions(): UncoveredFunction[] {
        if (!fs.existsSync(this.coverageFile)) {
            return [];
        }

        const output = execSync(`go tool cover -func=${this.coverageFile}`, { encoding: "utf-8" });

        const uncovered: UncoveredFunction[] = [];
        const lines = output.split("\n");

        for (const line of lines) {
            const match = line.match(/^(.+):(\d+):\s+(\S+)\s+([0-9.]+)%/);
            if (match) {
                const [, filePath, , functionName, coverage] = match;
                const coverageNum = parseFloat(coverage);

                if (coverageNum === 0) {
                    // Extract package from file path
                    const packageMatch = filePath.match(/apps\/backend\/(.+?)\//);
                    const packageName = packageMatch ? packageMatch[1] : path.dirname(filePath);

                    uncovered.push({
                        package: packageName,
                        function: functionName,
                        coverage: coverageNum,
                    });
                }
            }
        }

        return uncovered;
    }

    /**
     * Get status based on coverage percentage
     */
    private getStatus(coverage: number): PackageCoverage["status"] {
        if (coverage === 0) return "none";
        if (coverage < 30) return "low";
        if (coverage < 50) return "medium";
        if (coverage < 80) return "good";
        return "excellent";
    }

    /**
     * Get emoji for status
     */
    private getStatusEmoji(status: PackageCoverage["status"]): string {
        const emojis = {
            excellent: "✅",
            good: "🟢",
            medium: "🟡",
            low: "🟠",
            none: "🔴",
        };
        return emojis[status];
    }

    /**
     * Categorize packages by type
     */
    private categorizePackages(packages: PackageCoverage[]): CoverageReport["criticalGaps"] {
        const categories = [
            {
                category: "Usecase Layer",
                pattern: /^core-api\/internal\/usecase/,
                priority: "high",
            },
            {
                category: "Common/Utils",
                pattern: /^common\//,
                priority: "critical",
            },
            {
                category: "Services",
                pattern: /^services\//,
                priority: "medium",
            },
            {
                category: "Middleware",
                pattern: /^core-api\/internal\/middleware/,
                priority: "high",
            },
            {
                category: "Handlers",
                pattern: /^core-api\/internal\/handler/,
                priority: "low",
            },
        ];

        const uncoveredFunctions = this.parseUncoveredFunctions();

        return categories.map(({ category, pattern }) => {
            const categoryPackages = packages.filter((pkg) => pattern.test(pkg.package));
            const categoryFunctions = uncoveredFunctions.filter((fn) => pattern.test(fn.package));

            return {
                category,
                packages: categoryPackages,
                uncoveredFunctions: categoryFunctions,
            };
        });
    }

    /**
     * Generate detailed report
     */
    private generateReport(report: CoverageReport): void {
        console.log("\n📊 BACKEND TEST COVERAGE REPORT\n");
        console.log("═".repeat(80));

        // Summary
        console.log("\n## Summary\n");
        console.log(`Total Packages: ${report.summary.totalPackages}`);
        console.log(`Average Coverage: ${report.summary.averageCoverage.toFixed(1)}%`);
        console.log(`Uncovered Functions: ${report.summary.uncoveredFunctions}`);

        // Critical gaps by category
        console.log("\n## Coverage by Category\n");

        for (const gap of report.criticalGaps) {
            if (gap.packages.length === 0) continue;

            const avgCoverage =
                gap.packages.reduce((sum, p) => sum + p.coverage, 0) / gap.packages.length;
            console.log(`\n### ${gap.category} (${avgCoverage.toFixed(1)}% avg)\n`);

            // Show top 10 packages sorted by coverage
            const topPackages = gap.packages.sort((a, b) => a.coverage - b.coverage).slice(0, 10);

            for (const pkg of topPackages) {
                const emoji = this.getStatusEmoji(pkg.status);
                const pkgName = pkg.package.replace(/^core-api\/internal\/usecase\//, "usecase/");
                console.log(`  ${emoji} ${pkgName.padEnd(40)} ${pkg.coverage.toFixed(1)}%`);
            }

            // Show uncovered functions for critical packages
            if (gap.uncoveredFunctions.length > 0 && avgCoverage < 50) {
                console.log("\n  Critical Uncovered Functions:");
                const uniqueFunctions = new Map<string, string[]>();

                gap.uncoveredFunctions.slice(0, 20).forEach((fn) => {
                    const pkg = fn.package.split("/").pop() || fn.package;
                    if (!uniqueFunctions.has(pkg)) {
                        uniqueFunctions.set(pkg, []);
                    }
                    uniqueFunctions.get(pkg)!.push(fn.function);
                });

                for (const [pkg, functions] of uniqueFunctions) {
                    console.log(`    ${pkg}:`);
                    functions.slice(0, 5).forEach((fn) => {
                        console.log(`      - ${fn}()`);
                    });
                    if (functions.length > 5) {
                        console.log(`      ... and ${functions.length - 5} more`);
                    }
                }
            }
        }

        // Critical priorities
        console.log("\n## 🚨 Top Priorities (0% Coverage)\n");

        const criticalPackages = report.packages
            .filter(
                (pkg) =>
                    pkg.coverage === 0 &&
                    (pkg.package.includes("encryptutils") ||
                        pkg.package.includes("hashutils") ||
                        pkg.package.includes("pgmapper") ||
                        pkg.package.includes("usecase/event") ||
                        pkg.package.includes("usecase/eventconfig")),
            )
            .slice(0, 10);

        for (const pkg of criticalPackages) {
            console.log(`  🔴 ${pkg.package}`);
        }

        // Well-covered areas
        console.log("\n## ✅ Well-Covered Areas\n");

        const wellCovered = report.packages
            .filter((pkg) => pkg.coverage >= 80)
            .sort((a, b) => b.coverage - a.coverage)
            .slice(0, 10);

        for (const pkg of wellCovered) {
            const emoji = this.getStatusEmoji(pkg.status);
            console.log(`  ${emoji} ${pkg.package.padEnd(50)} ${pkg.coverage.toFixed(1)}%`);
        }

        console.log("\n" + "═".repeat(80) + "\n");
    }

    /**
     * Save report to JSON file
     */
    private saveJsonReport(report: CoverageReport): void {
        const outputPath = path.join(process.cwd(), "coverage-report.json");
        fs.writeFileSync(outputPath, JSON.stringify(report, null, 2));
        console.log(`📄 Full report saved to: ${outputPath}\n`);
    }

    /**
     * Main analysis method
     */
    public analyze(): void {
        // Run tests
        this.runTests();

        // Parse coverage
        const packages = this.parsePackageCoverage();
        const uncoveredFunctions = this.parseUncoveredFunctions();

        // Calculate summary
        const totalPackages = packages.length;
        const averageCoverage = packages.reduce((sum, p) => sum + p.coverage, 0) / totalPackages;

        // Generate report
        const report: CoverageReport = {
            summary: {
                totalPackages,
                averageCoverage,
                uncoveredFunctions: uncoveredFunctions.length,
            },
            packages,
            criticalGaps: this.categorizePackages(packages),
        };

        // Output report
        this.generateReport(report);
        this.saveJsonReport(report);
    }
}

// Run analyzer
const analyzer = new BackendCoverageAnalyzer();
analyzer.analyze();
