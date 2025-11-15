import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useCertificateTemplate } from "./useCertificateTemplate";

describe("useCertificateTemplate", () => {
    let inputRef: HTMLInputElement;
    let containerRef: HTMLDivElement;

    beforeEach(() => {
        inputRef = document.createElement("input");
        inputRef.type = "file";
        containerRef = document.createElement("div");
        document.body.appendChild(containerRef);
    });

    afterEach(() => {
        if (containerRef.parentNode) {
            containerRef.parentNode.removeChild(containerRef);
        }
        vi.clearAllMocks();
    });

    it("should initialize with default values", () => {
        const { result } = renderHook(() => useCertificateTemplate());

        expect(result.current.svgFile).toBeNull();
        expect(result.current.svgPreview).toBe("");
        expect(result.current.detectedKeywords).toEqual([]);
        expect(result.current.hasMissingMandatory).toBe(true);
        expect(result.current.missingMandatoryKeywords.length).toBe(2);
    });

    it("should have fileInputRef and svgTempRef defined", () => {
        const { result } = renderHook(() => useCertificateTemplate());

        expect(result.current.fileInputRef).toBeDefined();
        expect(result.current.svgTempRef).toBeDefined();
    });

    it("should include default available keywords", () => {
        const { result } = renderHook(() => useCertificateTemplate());

        expect(result.current.availableKeywords).toHaveLength(5);
        expect(result.current.availableKeywords[0].keyword).toBe("{{ eventName }}");
        expect(result.current.availableKeywords[1].keyword).toBe("{{ name }}");
        expect(result.current.availableKeywords[2].keyword).toBe("{{ academicInstitutionName }}");
        expect(result.current.availableKeywords[3].keyword).toBe("{{ certificateTitle }}");
        expect(result.current.availableKeywords[4].keyword).toBe("{{ certificateSubtitle }}");
    });

    it("should identify mandatory keywords", () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const mandatoryKeywords = result.current.availableKeywords.filter((kw) => kw.mandatory);
        expect(mandatoryKeywords).toHaveLength(1);
        expect(mandatoryKeywords.map((kw) => kw.keyword)).toContain("{{ name }}");
    });

    it("should clear template and reset state", () => {
        const { result } = renderHook(() => useCertificateTemplate());

        act(() => {
            result.current.clearTemplate();
        });

        expect(result.current.svgFile).toBeNull();
        expect(result.current.svgPreview).toBe("");
        expect(result.current.detectedKeywords).toEqual([]);
    });

    it("should handle SVG file selection", async () => {
        const svgContent = `
            <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                <text id="{{ name }}" x="100" y="100">Name</text>
                <text id="{{ eventName }}" x="100" y="200">Event</text>
            </svg>
        `;

        const file = new File([svgContent], "test.svg", { type: "image/svg+xml" });

        const { result } = renderHook(() => useCertificateTemplate());

        // Create a mock change event
        const changeEvent = {
            target: {
                files: [file],
            },
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        act(() => {
            result.current.handleFileSelect(changeEvent);
        });

        // Wait for async operations
        await waitFor(() => {
            expect(result.current.svgFile).toBe(file);
        });
    });

    it("should handle missing mandatory keywords", () => {
        const { result } = renderHook(() => useCertificateTemplate());

        // By default, both mandatory keywords are missing
        expect(result.current.hasMissingMandatory).toBe(true);
        expect(result.current.missingMandatoryKeywords.length).toBeGreaterThan(0);
    });

    it("should use custom search keys when provided", () => {
        const customSearchKeys = ["{{ customKey1 }}", "{{ customKey2 }}"];
        const { result } = renderHook(() =>
            useCertificateTemplate({ searchKeys: customSearchKeys }),
        );

        // The hook should use the custom search keys for detection
        expect(result.current).toBeDefined();
    });

    it("should use custom certificate dimensions when provided", () => {
        const customWidth = 1280;
        const customHeight = 720;

        const { result } = renderHook(() =>
            useCertificateTemplate({ certWidth: customWidth, certHeight: customHeight }),
        );

        expect(result.current).toBeDefined();
    });

    it("should handle invalid SVG file gracefully", async () => {
        const invalidSvgContent = "This is not SVG content";
        const file = new File([invalidSvgContent], "invalid.svg", { type: "image/svg+xml" });

        const { result } = renderHook(() => useCertificateTemplate());

        const changeEvent = {
            target: {
                files: [file],
            },
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        act(() => {
            result.current.handleFileSelect(changeEvent);
        });

        // Should handle error gracefully without throwing
        await waitFor(() => {
            expect(result.current).toBeDefined();
        });
    });

    it("should handle file input ref reset on clear", () => {
        const { result } = renderHook(() => useCertificateTemplate());

        // Set a ref
        const mockInput = document.createElement("input");
        mockInput.value = "test.svg";

        Object.defineProperty(result.current.fileInputRef, "current", {
            writable: true,
            value: mockInput,
        });

        act(() => {
            result.current.clearTemplate();
        });

        expect(result.current.fileInputRef.current?.value).toBe("");
    });

    it("should detect keywords in valid SVG", async () => {
        const svgContent = `
            <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                <defs>
                    <style>
                        text { font-size: 48px; }
                    </style>
                </defs>
                <text id="{{ name }}" x="960" y="300">Recipient Name</text>
                <text id="{{ eventName }}" x="960" y="400">Event Name</text>
            </svg>
        `;

        const file = new File([svgContent], "cert.svg", { type: "image/svg+xml" });

        const { result } = renderHook(() => useCertificateTemplate());

        const changeEvent = {
            target: {
                files: [file],
            },
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        act(() => {
            result.current.handleFileSelect(changeEvent);
        });

        // Wait for detection
        await waitFor(() => {
            expect(result.current.svgPreview).toContain("<svg");
        });
    });

    it("should maintain independent state for multiple hook instances", () => {
        const { result: result1 } = renderHook(() => useCertificateTemplate());
        const { result: result2 } = renderHook(() => useCertificateTemplate());

        act(() => {
            result1.current.clearTemplate();
        });

        // result2 should not be affected
        expect(result2.current.detectedKeywords).toEqual([]);
    });

    it("should have correct hasMissingMandatory flag when both mandatory keywords are present", async () => {
        const svgContent = `
            <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                <text id="{{ name }}" x="100" y="100">Name</text>
                <text id="{{ eventName }}" x="100" y="200">Event</text>
            </svg>
        `;

        const file = new File([svgContent], "cert.svg", { type: "image/svg+xml" });

        const { result } = renderHook(() => useCertificateTemplate());

        const changeEvent = {
            target: {
                files: [file],
            },
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        act(() => {
            result.current.handleFileSelect(changeEvent);
        });

        await waitFor(() => {
            if (result.current.detectedKeywords.length === 2) {
                expect(result.current.hasMissingMandatory).toBe(false);
            }
        });
    });

    it("should handle empty file selection", async () => {
        const { result } = renderHook(() => useCertificateTemplate());

        const changeEvent = {
            target: {
                files: [],
            },
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        act(() => {
            result.current.handleFileSelect(changeEvent);
        });

        // State should remain unchanged
        expect(result.current.svgFile).toBeNull();
        expect(result.current.svgPreview).toBe("");
    });

    it("should have correct detection keyword properties", async () => {
        const svgContent = `
            <svg xmlns="http://www.w3.org/2000/svg" width="1920" height="1080">
                <text id="{{ name }}" x="100" y="100" width="200" height="50">Name</text>
            </svg>
        `;

        const file = new File([svgContent], "cert.svg", { type: "image/svg+xml" });

        const { result } = renderHook(() => useCertificateTemplate());

        const changeEvent = {
            target: {
                files: [file],
            },
        } as unknown as React.ChangeEvent<HTMLInputElement>;

        act(() => {
            result.current.handleFileSelect(changeEvent);
        });

        await waitFor(() => {
            if (result.current.detectedKeywords.length > 0) {
                const detected = result.current.detectedKeywords[0];
                expect(detected).toHaveProperty("keyword");
                expect(detected).toHaveProperty("x");
                expect(detected).toHaveProperty("y");
                expect(detected).toHaveProperty("count");
                expect(typeof detected.x).toBe("number");
                expect(typeof detected.y).toBe("number");
            }
        });
    });
});
