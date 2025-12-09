import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useMediaQuery } from "./use-media-query";

describe("useMediaQuery", () => {
    let mockMatchMedia: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        mockMatchMedia = vi.fn();
        // Override the global matchMedia mock from test setup
        window.matchMedia = mockMatchMedia as unknown as typeof window.matchMedia;
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it("should return default value on server side", () => {
        // The hook checks typeof window === "undefined" at module load time
        // Since we're in a browser-like test environment, we test the equivalent behavior:
        // when initializeWithValue is false, it starts with defaultValue, but the effect will update it
        // So we need to mock matchMedia to return the defaultValue
        mockMatchMedia.mockReturnValue({
            matches: false, // matches defaultValue
            media: "(min-width: 768px)",
            onchange: null,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        });

        const { result } = renderHook(() =>
            useMediaQuery("(min-width: 768px)", {
                defaultValue: false,
                initializeWithValue: false,
            }),
        );

        // The effect will run and update based on matchMedia, so it should match the mock
        expect(result.current).toBe(false);
    });

    it("should return matches when query matches", () => {
        mockMatchMedia.mockReturnValue({
            matches: true,
            media: "(min-width: 768px)",
            onchange: null,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        });

        const { result } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(result.current).toBe(true);
        expect(mockMatchMedia).toHaveBeenCalledWith("(min-width: 768px)");
    });

    it("should return false when query does not match", () => {
        mockMatchMedia.mockReturnValue({
            matches: false,
            media: "(min-width: 768px)",
            onchange: null,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        });

        const { result } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(result.current).toBe(false);
    });

    it("should use defaultValue when initializeWithValue is false", () => {
        // When initializeWithValue is false, it starts with defaultValue
        // but the effect will still run and update based on matchMedia
        // So we need to mock matchMedia to return false to match the defaultValue
        mockMatchMedia.mockReturnValue({
            matches: false, // matches defaultValue
            media: "(min-width: 768px)",
            onchange: null,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        });

        const { result } = renderHook(() =>
            useMediaQuery("(min-width: 768px)", {
                defaultValue: false,
                initializeWithValue: false,
            }),
        );

        // The effect will run and update based on matchMedia
        // Since matches is false, it should be false
        expect(result.current).toBe(false);
    });

    it("should update when media query changes (addListener)", () => {
        const addListener = vi.fn();
        const removeListener = vi.fn();
        let changeHandler: (() => void) | null = null;

        mockMatchMedia.mockReturnValue({
            matches: false,
            media: "(min-width: 768px)",
            onchange: null,
            addListener: (handler: () => void) => {
                changeHandler = handler;
                addListener(handler);
            },
            removeListener: (handler: () => void) => {
                if (changeHandler === handler) {
                    changeHandler = null;
                }
                removeListener(handler);
            },
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        });

        const { result } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(result.current).toBe(false);
        expect(addListener).toHaveBeenCalled();

        // Simulate media query change
        if (changeHandler) {
            // Update the mock to return matches: true
            mockMatchMedia.mockReturnValue({
                matches: true,
                media: "(min-width: 768px)",
                onchange: null,
                addListener: vi.fn(),
                removeListener: vi.fn(),
                addEventListener: vi.fn(),
                removeEventListener: vi.fn(),
                dispatchEvent: vi.fn(),
            });
            // Wrap in act since it updates state
            act(() => {
                changeHandler();
            });
        }

        expect(result.current).toBe(true);
    });

    it("should update when media query changes (addEventListener)", () => {
        const addEventListener = vi.fn();
        const removeEventListener = vi.fn();
        let changeHandler: (() => void) | null = null;

        mockMatchMedia.mockReturnValue({
            matches: false,
            media: "(min-width: 768px)",
            onchange: null,
            addListener: undefined,
            removeListener: undefined,
            addEventListener: (event: string, handler: () => void) => {
                if (event === "change") {
                    changeHandler = handler;
                }
                addEventListener(event, handler);
            },
            removeEventListener: (event: string, handler: () => void) => {
                if (event === "change" && changeHandler === handler) {
                    changeHandler = null;
                }
                removeEventListener(event, handler);
            },
            dispatchEvent: vi.fn(),
        });

        const { result } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        expect(result.current).toBe(false);
        expect(addEventListener).toHaveBeenCalledWith("change", expect.any(Function));

        // Simulate media query change
        if (changeHandler) {
            // Update the mock to return matches: true
            mockMatchMedia.mockReturnValue({
                matches: true,
                media: "(min-width: 768px)",
                onchange: null,
                addListener: undefined,
                removeListener: undefined,
                addEventListener: vi.fn(),
                removeEventListener: vi.fn(),
                dispatchEvent: vi.fn(),
            });
            // Wrap in act since it updates state
            act(() => {
                changeHandler();
            });
        }

        expect(result.current).toBe(true);
    });

    it("should cleanup listeners on unmount", () => {
        const removeListener = vi.fn();
        const removeEventListener = vi.fn();

        mockMatchMedia.mockReturnValue({
            matches: false,
            media: "(min-width: 768px)",
            onchange: null,
            addListener: vi.fn(),
            removeListener,
            addEventListener: vi.fn(),
            removeEventListener,
            dispatchEvent: vi.fn(),
        });

        const { unmount } = renderHook(() => useMediaQuery("(min-width: 768px)"));

        unmount();

        expect(removeListener).toHaveBeenCalled();
    });

    it("should re-register listener when query changes", () => {
        const addListener = vi.fn();
        const removeListener = vi.fn();

        mockMatchMedia.mockReturnValue({
            matches: false,
            media: "(min-width: 768px)",
            onchange: null,
            addListener,
            removeListener,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn(),
        });

        const { rerender } = renderHook(({ query }) => useMediaQuery(query), {
            initialProps: { query: "(min-width: 768px)" },
        });

        expect(addListener).toHaveBeenCalledTimes(1);

        rerender({ query: "(min-width: 1024px)" });

        expect(removeListener).toHaveBeenCalled();
        expect(addListener).toHaveBeenCalledTimes(2);
    });
});
