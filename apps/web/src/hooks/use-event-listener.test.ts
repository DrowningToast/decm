import { describe, it, expect, vi, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useEventListener } from "./use-event-listener";
import { createRef } from "react";

describe("useEventListener", () => {
    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should add event listener to window when no element ref provided", () => {
        const mockHandler = vi.fn();
        const addListenerSpy = vi.spyOn(window, "addEventListener");
        const removeListenerSpy = vi.spyOn(window, "removeEventListener");

        const { unmount } = renderHook(() => useEventListener("click", mockHandler));

        expect(addListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), undefined);

        unmount();

        expect(removeListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), undefined);

        addListenerSpy.mockRestore();
        removeListenerSpy.mockRestore();
    });

    it("should add event listener to provided element ref", () => {
        const mockHandler = vi.fn();
        const mockDiv = document.createElement("div");
        const elementRef = createRef<HTMLDivElement>();
        Object.defineProperty(elementRef, "current", {
            value: mockDiv,
            writable: false,
        });

        const addListenerSpy = vi.spyOn(mockDiv, "addEventListener");
        const removeListenerSpy = vi.spyOn(mockDiv, "removeEventListener");

        const { unmount } = renderHook(() => useEventListener("click", mockHandler, elementRef));

        expect(addListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), undefined);

        unmount();

        expect(removeListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), undefined);

        addListenerSpy.mockRestore();
        removeListenerSpy.mockRestore();
    });

    it("should call handler when event is triggered", () => {
        const mockHandler = vi.fn();

        renderHook(() => useEventListener("resize", mockHandler));

        act(() => {
            window.dispatchEvent(new Event("resize"));
        });

        expect(mockHandler).toHaveBeenCalled();
    });

    it("should update handler when dependencies change", () => {
        const mockHandler1 = vi.fn();
        const mockHandler2 = vi.fn();

        const { rerender } = renderHook(
            ({ handler }: { handler: () => void }) => useEventListener("click", handler),
            { initialProps: { handler: mockHandler1 } },
        );

        act(() => {
            window.dispatchEvent(new Event("click"));
        });

        expect(mockHandler1).toHaveBeenCalledTimes(1);

        rerender({ handler: mockHandler2 });

        act(() => {
            window.dispatchEvent(new Event("click"));
        });

        expect(mockHandler1).toHaveBeenCalledTimes(1); // Should not increase
        expect(mockHandler2).toHaveBeenCalledTimes(1);
    });

    it("should handle event listener options", () => {
        const mockHandler = vi.fn();
        const mockDiv = document.createElement("div");
        const elementRef = createRef<HTMLDivElement>();
        Object.defineProperty(elementRef, "current", {
            value: mockDiv,
            writable: false,
        });

        const addListenerSpy = vi.spyOn(mockDiv, "addEventListener");

        const options = { capture: true, passive: true };
        renderHook(() => useEventListener("click", mockHandler, elementRef, options));

        expect(addListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), options);

        addListenerSpy.mockRestore();
    });

    it("should handle boolean options", () => {
        const mockHandler = vi.fn();
        const addListenerSpy = vi.spyOn(window, "addEventListener");

        renderHook(() => useEventListener("click", mockHandler, undefined, true));

        expect(addListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), true);

        addListenerSpy.mockRestore();
    });

    it("should handle null element ref gracefully", () => {
        const mockHandler = vi.fn();
        const elementRef = createRef<HTMLDivElement>();

        const { unmount } = renderHook(() => useEventListener("click", mockHandler, elementRef));

        // Should not throw when element is null
        expect(() => unmount()).not.toThrow();
    });

    it("should not add listener if target doesn't support addEventListener", () => {
        const mockHandler = vi.fn();
        const invalidRef = createRef<any>();
        (invalidRef as any).current = {}; // Object without addEventListener

        const { unmount } = renderHook(() => useEventListener("click", mockHandler, invalidRef));

        // Should not throw
        expect(() => unmount()).not.toThrow();
    });

    it("should remove previous listener when event name changes", () => {
        const mockHandler = vi.fn();
        const addListenerSpy = vi.spyOn(window, "addEventListener");
        const removeListenerSpy = vi.spyOn(window, "removeEventListener");

        const { rerender } = renderHook(
            ({ eventName }: { eventName: string }) =>
                useEventListener(eventName as keyof WindowEventMap, mockHandler),
            { initialProps: { eventName: "click" } },
        );

        expect(addListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), undefined);

        rerender({ eventName: "scroll" });

        expect(removeListenerSpy).toHaveBeenCalledWith("click", expect.any(Function), undefined);
        expect(addListenerSpy).toHaveBeenCalledWith("scroll", expect.any(Function), undefined);

        addListenerSpy.mockRestore();
        removeListenerSpy.mockRestore();
    });

    it("should handle multiple event listeners", () => {
        const mockHandler1 = vi.fn();
        const mockHandler2 = vi.fn();

        renderHook(() => useEventListener("click", mockHandler1));
        renderHook(() => useEventListener("scroll", mockHandler2));

        act(() => {
            window.dispatchEvent(new Event("click"));
            window.dispatchEvent(new Event("scroll"));
        });

        expect(mockHandler1).toHaveBeenCalled();
        expect(mockHandler2).toHaveBeenCalled();
    });

    it("should work with media query list events", () => {
        const mockHandler = vi.fn();
        const mediaQuery = window.matchMedia("(max-width: 600px)");
        const mqlRef = createRef<MediaQueryList>();
        Object.defineProperty(mqlRef, "current", {
            value: mediaQuery,
            writable: false,
        });

        const addListenerSpy = vi.spyOn(mediaQuery, "addEventListener");

        renderHook(() => useEventListener("change", mockHandler, mqlRef));

        expect(addListenerSpy).toHaveBeenCalledWith("change", expect.any(Function), undefined);

        addListenerSpy.mockRestore();
    });
});
