import { renderHook, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { useFetchedCertificateSvg } from "./useFetchedCertificateSvg";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const EVENT_ID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890";
const BACKEND_API = "http://localhost:8080";
const TEMPLATE_URL = `${BACKEND_API}/api/v1/events/${EVENT_ID}/config/certificate/template`;

const VALID_SVG =
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1920 1080">' +
    '<text id="{{ name }}">NAME</text></svg>';

const XML_DECLARATION_SVG =
    '<?xml version="1.0" encoding="UTF-8"?>' +
    '<svg xmlns="http://www.w3.org/2000/svg"><text id="{{ name }}">NAME</text></svg>';

// Mock the env module so VITE_CORE_BACKEND_API is available in tests.
// vi.mock is hoisted before variable declarations, so a string literal must be
// used here — referencing BACKEND_API would cause a temporal dead zone error.
vi.mock("@/config/env", () => ({
    env: { VITE_CORE_BACKEND_API: "http://localhost:8080" },
}));

function mockFetchSuccess(text: string) {
    global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(text),
    } as unknown as Response);
}

function mockFetchFailure(message = "Network error") {
    global.fetch = vi.fn().mockRejectedValue(new Error(message));
}

function mockFetchHttpError(status: number) {
    global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status,
        text: () => Promise.resolve("error"),
    } as unknown as Response);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("useFetchedCertificateSvg", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    // --- Happy path ----------------------------------------------------------

    it("fetches SVG via the backend proxy and returns its text content", async () => {
        mockFetchSuccess(VALID_SVG);

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, false));

        await waitFor(() => {
            expect(result.current).toBe(VALID_SVG);
        });

        expect(global.fetch).toHaveBeenCalledWith(TEMPLATE_URL, { credentials: "include" });
    });

    it("accepts SVG content that starts with an XML declaration", async () => {
        mockFetchSuccess(XML_DECLARATION_SVG);

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, false));

        await waitFor(() => {
            expect(result.current).toBe(XML_DECLARATION_SVG);
        });
    });

    it("accepts SVG content with leading whitespace before the <svg tag", async () => {
        const paddedSvg = "   \n" + VALID_SVG;
        mockFetchSuccess(paddedSvg);

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, false));

        await waitFor(() => {
            expect(result.current).toBe(paddedSvg);
        });
    });

    // --- Fallback conditions -------------------------------------------------

    it("returns empty string when the event ID is undefined", async () => {
        global.fetch = vi.fn();

        const { result } = renderHook(() => useFetchedCertificateSvg(undefined, false));

        // Allow a tick for any erroneous effect to fire
        await new Promise((r) => setTimeout(r, 30));

        expect(result.current).toBe("");
        expect(global.fetch).not.toHaveBeenCalled();
    });

    it("returns empty string and skips the fetch when a new file has been uploaded", async () => {
        global.fetch = vi.fn();

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, true));

        await new Promise((r) => setTimeout(r, 30));

        expect(result.current).toBe("");
        expect(global.fetch).not.toHaveBeenCalled();
    });

    it("returns empty string when the fetch fails with a network error", async () => {
        mockFetchFailure();

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, false));

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalled();
        });

        expect(result.current).toBe("");
    });

    it("returns empty string when the server responds with an HTTP error", async () => {
        mockFetchHttpError(403);

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, false));

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalled();
        });

        expect(result.current).toBe("");
    });

    it("returns empty string when the response is not SVG content", async () => {
        mockFetchSuccess("<html><body>Not an SVG</body></html>");

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, false));

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalled();
        });

        expect(result.current).toBe("");
    });

    it("returns empty string for a plain-text response", async () => {
        mockFetchSuccess("Just some plain text");

        const { result } = renderHook(() => useFetchedCertificateSvg(EVENT_ID, false));

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalled();
        });

        expect(result.current).toBe("");
    });

    // --- Reactivity ----------------------------------------------------------

    it("re-fetches and updates when the event ID changes", async () => {
        const firstId = "aaaaaaaa-0000-0000-0000-000000000001";
        const secondId = "aaaaaaaa-0000-0000-0000-000000000002";
        const firstSvg = '<svg id="first"></svg>';
        const secondSvg = '<svg id="second"></svg>';

        global.fetch = vi
            .fn()
            .mockResolvedValueOnce({
                ok: true,
                text: () => Promise.resolve(firstSvg),
            } as unknown as Response)
            .mockResolvedValueOnce({
                ok: true,
                text: () => Promise.resolve(secondSvg),
            } as unknown as Response);

        const { result, rerender } = renderHook(
            ({ id }: { id: string }) => useFetchedCertificateSvg(id, false),
            { initialProps: { id: firstId } },
        );

        await waitFor(() => expect(result.current).toBe(firstSvg));

        rerender({ id: secondId });

        await waitFor(() => expect(result.current).toBe(secondSvg));

        expect(global.fetch).toHaveBeenCalledTimes(2);
        expect(global.fetch).toHaveBeenNthCalledWith(
            1,
            `${BACKEND_API}/api/v1/events/${firstId}/config/certificate/template`,
            { credentials: "include" },
        );
        expect(global.fetch).toHaveBeenNthCalledWith(
            2,
            `${BACKEND_API}/api/v1/events/${secondId}/config/certificate/template`,
            { credentials: "include" },
        );
    });

    it("clears the cached SVG and stops fetching when hasNewFileUploaded switches to true", async () => {
        mockFetchSuccess(VALID_SVG);

        const { result, rerender } = renderHook(
            ({ uploaded }: { uploaded: boolean }) => useFetchedCertificateSvg(EVENT_ID, uploaded),
            { initialProps: { uploaded: false } },
        );

        await waitFor(() => expect(result.current).toBe(VALID_SVG));

        // Simulate user selecting a new file
        global.fetch = vi.fn();
        rerender({ uploaded: true });

        await waitFor(() => expect(result.current).toBe(""));
        expect(global.fetch).not.toHaveBeenCalled();
    });

    it("does not apply a stale response after the event ID changes", async () => {
        const firstId = "aaaaaaaa-0000-0000-0000-000000000001";
        const secondId = "aaaaaaaa-0000-0000-0000-000000000002";
        const fastSvg = '<svg id="fast"></svg>';

        // First fetch is deliberately slow; second resolves immediately
        let resolveFirst!: (v: Response) => void;
        const slowResponse = new Promise<Response>((res) => (resolveFirst = res));

        global.fetch = vi
            .fn()
            .mockReturnValueOnce(slowResponse)
            .mockResolvedValueOnce({
                ok: true,
                text: () => Promise.resolve(fastSvg),
            } as unknown as Response);

        const { result, rerender } = renderHook(
            ({ id }: { id: string }) => useFetchedCertificateSvg(id, false),
            { initialProps: { id: firstId } },
        );

        // Switch event ID before the first fetch completes
        rerender({ id: secondId });

        await waitFor(() => expect(result.current).toBe(fastSvg));

        // Now resolve the stale first request — it should be ignored
        resolveFirst({
            ok: true,
            text: () => Promise.resolve('<svg id="stale"></svg>'),
        } as unknown as Response);

        await new Promise((r) => setTimeout(r, 30));
        expect(result.current).toBe(fastSvg);
    });
});
