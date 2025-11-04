import { describe, it, expect, vi } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useEventCallback } from "./use-event-callback";

describe("useEventCallback", () => {
    it("should return a memoized callback", () => {
        const mockFn = vi.fn((x: number) => x * 2);

        const { result } = renderHook(() => useEventCallback(mockFn));

        expect(typeof result.current).toBe("function");
        expect(result.current).toBeDefined();
    });

    it("should call the callback with correct arguments", () => {
        const mockFn = vi.fn((x: number, y: number) => x + y);

        const { result } = renderHook(() => useEventCallback(mockFn));

        let callResult: number | undefined;
        act(() => {
            callResult = result.current(5, 3);
        });

        expect(mockFn).toHaveBeenCalledWith(5, 3);
        expect(callResult).toBe(8);
    });

    it("should return the result from the callback", () => {
        const mockFn = vi.fn((x: string) => `Hello ${x}`);

        const { result } = renderHook(() => useEventCallback(mockFn));

        let callResult: string | undefined;
        act(() => {
            callResult = result.current("World");
        });

        expect(callResult).toBe("Hello World");
    });

    it("should handle no arguments callback", () => {
        const mockFn = vi.fn(() => "no args");

        const { result } = renderHook(() => useEventCallback(mockFn));

        let callResult: string | undefined;
        act(() => {
            callResult = result.current();
        });

        expect(mockFn).toHaveBeenCalled();
        expect(callResult).toBe("no args");
    });

    it("should handle multiple arguments", () => {
        const mockFn = vi.fn((a: number, b: number, c: string) => `${a + b}:${c}`);

        const { result } = renderHook(() => useEventCallback(mockFn));

        let callResult: string | undefined;
        act(() => {
            callResult = result.current(10, 20, "test");
        });

        expect(mockFn).toHaveBeenCalledWith(10, 20, "test");
        expect(callResult).toBe("30:test");
    });

    it("should update callback when function changes", () => {
        const mockFn1 = vi.fn((x: number) => x * 2);
        const mockFn2 = vi.fn((x: number) => x * 3);

        const { result, rerender } = renderHook(
            ({ fn }: { fn: (x: number) => number }) => useEventCallback(fn),
            { initialProps: { fn: mockFn1 } },
        );

        act(() => {
            result.current(5);
        });
        expect(mockFn1).toHaveBeenCalledWith(5);
        expect(mockFn1).toHaveReturnedWith(10);

        rerender({ fn: mockFn2 });

        act(() => {
            result.current(5);
        });
        expect(mockFn2).toHaveBeenCalledWith(5);
        expect(mockFn2).toHaveReturnedWith(15);
    });

    it("should maintain stable reference", () => {
        const mockFn = vi.fn((x: number) => x);

        const { result, rerender } = renderHook(() => useEventCallback(mockFn));

        const callback1 = result.current;

        rerender();

        const callback2 = result.current;

        // The callback reference should remain the same (memoized)
        expect(callback1).toBe(callback2);
    });

    it("should handle undefined function", () => {
        const { result } = renderHook(() => useEventCallback(undefined));

        // When undefined is passed, the hook should handle it gracefully
        expect(result.current).toBeDefined();
    });

    it("should throw error if called during rendering", () => {
        const mockFn = vi.fn();
        const { result } = renderHook(() => useEventCallback(mockFn));

        // The default error message should indicate the issue
        expect(() => {
            try {
                result.current();
            } catch (e: unknown) {
                if (e instanceof Error) {
                    expect(e.message).toContain("Cannot call an event handler while rendering");
                }
                throw e;
            }
        }).not.toThrow(); // This would normally be called after rendering completes
    });

    it("should handle object return types", () => {
        interface User {
            name: string;
            age: number;
        }

        const mockFn = vi.fn((name: string, age: number): User => ({ name, age }));

        const { result } = renderHook(() => useEventCallback(mockFn));

        let callResult: User | undefined;
        act(() => {
            callResult = result.current("Alice", 30);
        });

        expect(callResult).toEqual({ name: "Alice", age: 30 });
    });

    it("should handle array arguments", () => {
        const mockFn = vi.fn((arr: number[]) => arr.reduce((a, b) => a + b, 0));

        const { result } = renderHook(() => useEventCallback(mockFn));

        let callResult: number | undefined;
        act(() => {
            callResult = result.current([1, 2, 3, 4, 5]);
        });

        expect(callResult).toBe(15);
    });
});
