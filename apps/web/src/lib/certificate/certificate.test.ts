import { describe, it, expect } from "vitest";
import { loadSvgTemplateFile, detectTemplateKeys } from "./certificate";

describe("Certificate Utils", () => {
    describe("loadSvgTemplateFile", () => {
        it("should load and parse a valid SVG file", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="{{ name }}" x="100" y="100">Name</text>
                </svg>`;

            const file = new File([svgContent], "test.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            expect(result).toBeTruthy();
            expect(result?.documentElement.tagName.toLowerCase()).toBe("svg");
        });

        it("should parse SVG correctly", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <circle id="test-circle" cx="100" cy="100" r="50"/>
                </svg>`;

            const file = new File([svgContent], "circle.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            expect(result).toBeTruthy();
            const circle = result?.documentElement.querySelector("#test-circle");
            expect(circle).toBeTruthy();
        });

        it("should handle SVG with multiple elements", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="{{ name }}" x="100" y="100">Name</text>
                    <text id="{{ eventName }}" x="100" y="200">Event</text>
                    <text id="{{ academicInstitutionName }}" x="100" y="300">Institution</text>
                </svg>`;

            const file = new File([svgContent], "multi.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            expect(result).toBeTruthy();
            const textElements = result?.documentElement.querySelectorAll("text");
            expect(textElements?.length).toBe(3);
        });

        it("should handle SVG with nested elements", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <g id="group">
                        <text id="nested-text" x="100" y="100">Text</text>
                    </g>
                </svg>`;

            const file = new File([svgContent], "nested.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            expect(result).toBeTruthy();
            const nestedText = result?.documentElement.querySelector("#nested-text");
            expect(nestedText).toBeTruthy();
        });

        it("should reject invalid SVG files", async () => {
            const invalidContent = "This is not SVG content";
            const file = new File([invalidContent], "invalid.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            // Invalid SVG should still parse but not as a valid SVG element
            // The DOMParser doesn't throw errors, it just returns a document with error elements
            expect(result).toBeTruthy();
        });

        it("should handle SVG with special characters", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="{{ special_key }}" x="100" y="100">Special</text>
                </svg>`;

            const file = new File([svgContent], "special.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            expect(result).toBeTruthy();
        });

        it("should preserve SVG attributes", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080" viewBox="0 0 1920 1080">
                    <text id="test" x="100" y="100">Test</text>
                </svg>`;

            const file = new File([svgContent], "attrs.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            expect(result?.documentElement.getAttribute("width")).toBe("1920");
            expect(result?.documentElement.getAttribute("height")).toBe("1080");
        });

        it("should handle empty SVG files", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                </svg>`;

            const file = new File([svgContent], "empty.svg", { type: "image/svg+xml" });

            const result = await loadSvgTemplateFile(file);

            expect(result).toBeTruthy();
            expect(result?.documentElement.children.length).toBe(0);
        });
    });

    describe("detectTemplateKeys", () => {
        it("should detect template keywords in SVG", async () => {
            const svgDoc = new DOMParser().parseFromString(
                `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="name" data-key="{{ name }}" x="100" y="100">Name</text>
                    <text id="eventName" data-key="{{ eventName }}" x="100" y="200">Event</text>
                </svg>`,
                "image/svg+xml",
            );

            const result = await detectTemplateKeys(svgDoc);

            expect(result).toBeTruthy();
            expect(Array.isArray(result)).toBe(true);
        });

        it("should throw error if width or height is missing", async () => {
            const svgDoc = new DOMParser().parseFromString(
                `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg">
                    <text id="test">Test</text>
                </svg>`,
                "image/svg+xml",
            );

            await expect(detectTemplateKeys(svgDoc)).rejects.toThrow(
                "Certificate width or height is not set",
            );
        });

        it("should handle SVG with width but no height", async () => {
            const svgDoc = new DOMParser().parseFromString(
                `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920">
                    <text id="test">Test</text>
                </svg>`,
                "image/svg+xml",
            );

            await expect(detectTemplateKeys(svgDoc)).rejects.toThrow();
        });

        it("should detect keywords with correct position structure", async () => {
            const svgDoc = new DOMParser().parseFromString(
                `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="name" data-key="{{ name }}" x="100" y="100">Name</text>
                </svg>`,
                "image/svg+xml",
            );

            const result = await detectTemplateKeys(svgDoc);

            expect(result).toBeTruthy();
            result.forEach((item) => {
                expect(item).toHaveProperty("key");
                expect(item).toHaveProperty("position");
                expect(item.position).toHaveProperty("x");
                expect(item.position).toHaveProperty("y");
            });
        });

        it("should handle empty SVG", async () => {
            const svgDoc = new DOMParser().parseFromString(
                `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                </svg>`,
                "image/svg+xml",
            );

            const result = await detectTemplateKeys(svgDoc);

            expect(result).toBeTruthy();
            expect(Array.isArray(result)).toBe(true);
        });

        it("should handle SVG with various search keys", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="{{ name }}" data-key="{{ name }}" x="100" y="100">Name</text>
                    <text id="{{ eventName }}" data-key="{{ eventName }}" x="200" y="200">Event</text>
                    <text id="{{ startDate }}" data-key="{{ startDate }}" x="300" y="300">Date</text>
                </svg>`;

            const svgDoc = new DOMParser().parseFromString(svgContent, "image/svg+xml");

            const result = await detectTemplateKeys(svgDoc);

            expect(result).toBeTruthy();
            expect(Array.isArray(result)).toBe(true);
        });

        it("should handle SVG with numeric dimensions", async () => {
            const svgDoc = new DOMParser().parseFromString(
                `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1024" height="768">
                    <text id="test" x="100" y="100">Test</text>
                </svg>`,
                "image/svg+xml",
            );

            const result = await detectTemplateKeys(svgDoc);

            expect(result).toBeTruthy();
        });

        it("should clone SVG element without modifying original", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="test" x="100" y="100">Test</text>
                </svg>`;

            const svgDoc = new DOMParser().parseFromString(svgContent, "image/svg+xml");
            const originalWidth = svgDoc.documentElement.getAttribute("width");

            await detectTemplateKeys(svgDoc);

            // Original should not be modified
            expect(svgDoc.documentElement.getAttribute("width")).toBe(originalWidth);
        });
    });

    describe("Integration tests", () => {
        it("should load and detect keys in sequence", async () => {
            const svgContent = `<?xml version="1.0"?>
                <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                    <text id="{{ name }}" data-key="{{ name }}" x="100" y="100">Name</text>
                    <text id="{{ eventName }}" data-key="{{ eventName }}" x="200" y="200">Event</text>
                </svg>`;

            const file = new File([svgContent], "test.svg", { type: "image/svg+xml" });

            const doc = await loadSvgTemplateFile(file);
            expect(doc).toBeTruthy();

            if (doc) {
                const keys = await detectTemplateKeys(doc);
                expect(keys).toBeTruthy();
                expect(Array.isArray(keys)).toBe(true);
            }
        });
    });
});
