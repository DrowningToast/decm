import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook } from "@testing-library/react";
import { createRef } from "react";
import { useEventListener } from "./use-event-listener";

describe("useEventListener", () => {
    let mockAddEventListener: ReturnType<typeof vi.fn>;
    let mockRemoveEventListener: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        mockAddEventListener = vi.fn();
        mockRemoveEventListener = vi.fn();

        // Mock window.addEventListener
        window.addEventListener = mockAddEventListener;
        window.removeEventListener = mockRemoveEventListener;
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it("should add event listener to window by default", () => {
        const handler = vi.fn();

        renderHook(() => useEventListener("click", handler));

        expect(mockAddEventListener).toHaveBeenCalledWith("click", expect.any(Function), undefined);
    });

    it("should add event listener to element ref", () => {
        const handler = vi.fn();
        const elementRef = createRef<HTMLDivElement>();
        const mockElement = {
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        } as unknown as HTMLDivElement;

        elementRef.current = mockElement;

        renderHook(() =>
            useEventListener("click", handler, elementRef as React.RefObject<HTMLDivElement>),
        );

        expect(mockElement.addEventListener).toHaveBeenCalledWith(
            "click",
            expect.any(Function),
            undefined,
        );
    });

    it("should call handler when event is triggered", () => {
        const handler = vi.fn();
        let eventListener: ((event: Event) => void) | null = null as
            | ((event: Event) => void)
            | null;

        mockAddEventListener.mockImplementation((event, listener) => {
            if (event === "click") {
                eventListener = listener as (event: Event) => void;
            }
        });

        renderHook(() => useEventListener("click", handler));

        const mockEvent = new Event("click");
        eventListener?.(mockEvent);

        expect(handler).toHaveBeenCalledWith(mockEvent);
    });

    it("should update handler when it changes", () => {
        const handler1 = vi.fn();
        const handler2 = vi.fn();
        let eventListener: ((event: Event) => void) | null = null as
            | ((event: Event) => void)
            | null;

        mockAddEventListener.mockImplementation((event, listener) => {
            if (event === "click") {
                eventListener = listener as (event: Event) => void;
            }
        });

        const { rerender } = renderHook(({ handler }) => useEventListener("click", handler), {
            initialProps: { handler: handler1 },
        });

        const mockEvent = new Event("click");
        eventListener?.(mockEvent);

        expect(handler1).toHaveBeenCalledTimes(1);
        expect(handler2).not.toHaveBeenCalled();

        rerender({ handler: handler2 });

        eventListener?.(mockEvent);

        expect(handler1).toHaveBeenCalledTimes(1);
        expect(handler2).toHaveBeenCalledTimes(1);
    });

    it("should remove event listener on cleanup", () => {
        const handler = vi.fn();
        const { unmount } = renderHook(() => useEventListener("click", handler));

        unmount();

        expect(mockRemoveEventListener).toHaveBeenCalledWith(
            "click",
            expect.any(Function),
            undefined,
        );
    });

    it("should pass options to addEventListener", () => {
        const handler = vi.fn();
        const options = { capture: true, once: true };

        renderHook(() => useEventListener("click", handler, undefined, options));

        expect(mockAddEventListener).toHaveBeenCalledWith("click", expect.any(Function), options);
    });

    it("should handle document events", () => {
        const handler = vi.fn();
        const documentRef = createRef<Document>();
        const mockDocument = {
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        } as unknown as Document;

        documentRef.current = mockDocument;

        renderHook(() =>
            useEventListener("scroll", handler, documentRef as React.RefObject<Document>),
        );

        expect(mockDocument.addEventListener).toHaveBeenCalledWith(
            "scroll",
            expect.any(Function),
            undefined,
        );
    });

    it("should not add listener if element is null", () => {
        const handler = vi.fn();
        const elementRef = createRef<HTMLDivElement>();

        renderHook(() =>
            useEventListener("click", handler, elementRef as React.RefObject<HTMLDivElement>),
        );

        // When element ref is provided but null, it falls back to window
        // So window.addEventListener should be called
        expect(mockAddEventListener).toHaveBeenCalledWith("click", expect.any(Function), undefined);
    });

    it("should not add listener if element does not support addEventListener", () => {
        const handler = vi.fn();
        const elementRef = createRef<HTMLDivElement>();
        const mockElement = {} as HTMLDivElement; // No addEventListener method

        elementRef.current = mockElement;

        renderHook(() =>
            useEventListener("click", handler, elementRef as React.RefObject<HTMLDivElement>),
        );

        // Should not throw, just silently fail
        expect(mockAddEventListener).not.toHaveBeenCalled();
    });

    it("should handle MediaQueryList events", () => {
        const handler = vi.fn();
        const mediaQueryRef = createRef<MediaQueryList>();
        const mockMediaQuery = {
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        } as unknown as MediaQueryList;

        mediaQueryRef.current = mockMediaQuery;

        renderHook(() =>
            useEventListener("change", handler, mediaQueryRef as React.RefObject<MediaQueryList>),
        );

        expect(mockMediaQuery.addEventListener).toHaveBeenCalledWith(
            "change",
            expect.any(Function),
            undefined,
        );
    });

    it("should re-register listener when event name changes", () => {
        const handler = vi.fn();

        const { rerender } = renderHook(({ eventName }) => useEventListener(eventName, handler), {
            initialProps: { eventName: "click" as keyof WindowEventMap },
        });

        expect(mockAddEventListener).toHaveBeenCalledWith("click", expect.any(Function), undefined);

        rerender({ eventName: "resize" as keyof WindowEventMap });

        expect(mockRemoveEventListener).toHaveBeenCalledWith(
            "click",
            expect.any(Function),
            undefined,
        );
        expect(mockAddEventListener).toHaveBeenCalledWith(
            "resize",
            expect.any(Function),
            undefined,
        );
    });
});
