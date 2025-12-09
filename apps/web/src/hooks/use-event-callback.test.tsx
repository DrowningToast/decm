import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useEventCallback } from "./use-event-callback";

describe("useEventCallback", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should return a stable callback reference", () => {
        const fn = vi.fn();
        const { result, rerender } = renderHook(() => useEventCallback(fn));

        const firstCallback = result.current;

        rerender();

        expect(result.current).toBe(firstCallback);
    });

    it("should call the latest function with correct arguments", () => {
        const fn1 = vi.fn((a: number, b: string) => `${a}-${b}`);
        const { result, rerender } = renderHook(({ fn }) => useEventCallback(fn), {
            initialProps: { fn: fn1 },
        });

        act(() => {
            const result1 = result.current(1, "test");
            expect(result1).toBe("1-test");
            expect(fn1).toHaveBeenCalledWith(1, "test");
        });

        const fn2 = vi.fn((a: number, b: string) => `${b}-${a}`);
        rerender({ fn: fn2 });

        act(() => {
            const result2 = result.current(2, "hello");
            expect(result2).toBe("hello-2");
            expect(fn2).toHaveBeenCalledWith(2, "hello");
            expect(fn1).toHaveBeenCalledTimes(1); // Should not be called again
        });
    });

    it("should handle undefined function", () => {
        const { result } = renderHook(() => useEventCallback(undefined));

        expect(result.current).toBeUndefined();
    });

    it("should handle function that returns undefined", () => {
        const fn = vi.fn(() => undefined);
        const { result } = renderHook(() => useEventCallback(fn));

        act(() => {
            const returnValue = result.current();
            expect(returnValue).toBeUndefined();
            expect(fn).toHaveBeenCalled();
        });
    });

    it("should handle async functions", async () => {
        const fn = vi.fn(async (value: string) => {
            await new Promise((resolve) => setTimeout(resolve, 10));
            return `async-${value}`;
        });

        const { result } = renderHook(() => useEventCallback(fn));

        let returnValue: string | undefined;
        await act(async () => {
            returnValue = await result.current("test");
        });

        expect(returnValue).toBe("async-test");
        expect(fn).toHaveBeenCalledWith("test");
    });

    it("should maintain stable reference when function changes", () => {
        const fn1 = vi.fn();
        const fn2 = vi.fn();
        const { result, rerender } = renderHook(({ fn }) => useEventCallback(fn), {
            initialProps: { fn: fn1 },
        });

        const callback1 = result.current;

        rerender({ fn: fn2 });

        // Reference should remain stable
        expect(result.current).toBe(callback1);

        // But should call the new function
        act(() => {
            result.current();
        });

        expect(fn2).toHaveBeenCalled();
        expect(fn1).not.toHaveBeenCalled();
    });

    it("should handle multiple arguments", () => {
        const fn = vi.fn((a: number, b: string, c: boolean) => `${a}-${b}-${c}`);
        const { result } = renderHook(() => useEventCallback(fn));

        act(() => {
            const returnValue = result.current(1, "test", true);
            expect(returnValue).toBe("1-test-true");
            expect(fn).toHaveBeenCalledWith(1, "test", true);
        });
    });

    it("should handle no arguments", () => {
        const fn = vi.fn(() => "no-args");
        const { result } = renderHook(() => useEventCallback(fn));

        act(() => {
            const returnValue = result.current();
            expect(returnValue).toBe("no-args");
            expect(fn).toHaveBeenCalledWith();
        });
    });
});
