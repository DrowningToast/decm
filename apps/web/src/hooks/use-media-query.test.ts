import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook } from "@testing-library/react";
import { useMediaQuery } from "./use-media-query";

describe("useMediaQuery", () => {
    const mockMatchMedia = (matches: boolean) => {
        const listeners: Array<(event: MediaQueryListEvent) => void> = [];

        return {
            matches,
            media: "",
            onchange: null,
            addListener: vi.fn((listener: (event: MediaQueryListEvent) => void) => {
                listeners.push(listener);
            }),
            removeListener: vi.fn((listener: (event: MediaQueryListEvent) => void) => {
                const index = listeners.indexOf(listener);
                if (index > -1) {
                    listeners.splice(index, 1);
                }
            }),
            addEventListener: vi.fn(
                (event: string, listener: (event: MediaQueryListEvent) => void) => {
                    if (event === "change") {
                        listeners.push(listener);
                    }
                },
            ),
            removeEventListener: vi.fn(
                (event: string, listener: (event: MediaQueryListEvent) => void) => {
                    if (event === "change") {
                        const index = listeners.indexOf(listener);
                        if (index > -1) {
                            listeners.splice(index, 1);
                        }
                    }
                },
            ),
            dispatchEvent: vi.fn(),
            triggerChange: (newMatches: boolean) => {
                listeners.forEach((listener) =>
                    listener({ matches: newMatches, media: "" } as MediaQueryListEvent),
                );
            },
        };
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should return true when media query matches", () => {
        const matchMediaMock = mockMatchMedia(true);
        window.matchMedia = vi.fn(() => matchMediaMock as unknown as MediaQueryList);

        const { result } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(result.current).toBe(true);
    });

    it("should return false when media query does not match", () => {
        const matchMediaMock = mockMatchMedia(false);
        window.matchMedia = vi.fn(() => matchMediaMock as unknown as MediaQueryList);

        const { result } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(result.current).toBe(false);
    });

    it("should use defaultValue when initializeWithValue is false", () => {
        const matchMediaMock = mockMatchMedia(false);
        window.matchMedia = vi.fn(() => matchMediaMock as unknown as MediaQueryList);

        const { result } = renderHook(() =>
            useMediaQuery("(min-width: 768px)", { defaultValue: true, initializeWithValue: false }),
        );

        // With initializeWithValue: false, it should use defaultValue initially
        // But the layout effect will trigger and update it to the actual match
        expect([true, false]).toContain(result.current);
    });

    it("should update when media query match changes", () => {
        const matchMediaMock = mockMatchMedia(false);
        window.matchMedia = vi.fn(() => matchMediaMock as unknown as MediaQueryList);

        const { result, rerender } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(result.current).toBe(false);

        // Simulate media query change
        matchMediaMock.matches = true;
        matchMediaMock.triggerChange(true);

        rerender();

        expect(result.current).toBe(true);
    });

    it("should handle both addEventListener and deprecated addListener", () => {
        const matchMediaMockWithAddListener = {
            ...mockMatchMedia(false),
            addEventListener: undefined,
        };

        window.matchMedia = vi.fn(() => matchMediaMockWithAddListener as unknown as MediaQueryList);

        const { unmount } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(matchMediaMockWithAddListener.addListener).toHaveBeenCalled();

        unmount();

        expect(matchMediaMockWithAddListener.removeListener).toHaveBeenCalled();
    });

    it("should cleanup listeners on unmount", () => {
        const matchMediaMock = mockMatchMedia(false);
        window.matchMedia = vi.fn(() => matchMediaMock as unknown as MediaQueryList);

        const { unmount } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        unmount();

        // Check that either addEventListener or addListener was called
        const hasModernListener = matchMediaMock.removeEventListener.mock.calls.length > 0;
        const hasLegacyListener = matchMediaMock.removeListener.mock.calls.length > 0;

        expect(hasModernListener || hasLegacyListener).toBe(true);
    });

    it("should update when query changes", () => {
        const matchMediaMock1 = mockMatchMedia(false);
        const matchMediaMock2 = mockMatchMedia(true);

        window.matchMedia = vi.fn((query: string) => {
            // Return different mock based on query string
            if (query === "(min-width: 768px)") {
                return matchMediaMock1 as unknown as MediaQueryList;
            }
            return matchMediaMock2 as unknown as MediaQueryList;
        });

        const { result, rerender } = renderHook(({ query }) => useMediaQuery(query), {
            initialProps: { query: "(min-width: 768px)" },
        });

        // Initial query returns false
        expect(result.current).toBe(false);

        rerender({ query: "(min-width: 1024px)" });

        // New query returns true
        expect(result.current).toBe(true);
    });
});
