import { describe, it, expect } from "vitest";
import React from "react";
import { renderHook } from "@testing-library/react";
import { useIsomorphicLayoutEffect } from "./use-isomorphic-layout-effect";

describe("useIsomorphicLayoutEffect", () => {
    it("should export a function", () => {
        expect(typeof useIsomorphicLayoutEffect).toBe("function");
    });

    it("should be useLayoutEffect when window is defined", () => {
        // In a browser environment (and in vitest which simulates it)
        // useIsomorphicLayoutEffect should equal useLayoutEffect
        expect(useIsomorphicLayoutEffect).toBe(React.useLayoutEffect);
    });

    it("should work with the same signature as useLayoutEffect", () => {
        // The hook should be callable with the same arguments as useLayoutEffect
        const useTestHook = () => {
            useIsomorphicLayoutEffect(() => {
                // Effect body
            }, [1, 2, 3]);
        };

        // This test verifies the type compatibility - should not throw when used in a component
        expect(() => {
            renderHook(useTestHook);
        }).not.toThrow();
    });

    it("should be a valid React hook", () => {
        // Verify it's a function that can be used as a React hook
        expect(useIsomorphicLayoutEffect).toBeDefined();
        expect(typeof useIsomorphicLayoutEffect).toBe("function");

        // In test environment, it should be useLayoutEffect
        const isLayoutEffect = useIsomorphicLayoutEffect === React.useLayoutEffect;
        const isEffect = useIsomorphicLayoutEffect === React.useEffect;

        // One of them should be true
        expect(isLayoutEffect || isEffect).toBe(true);
    });

    it("should call the effect function", () => {
        let effectCalled = false;
        const useTestHook = () => {
            useIsomorphicLayoutEffect(() => {
                effectCalled = true;
            });
        };

        renderHook(useTestHook);

        expect(effectCalled).toBe(true);
    });

    it("should support dependency array", () => {
        let callCount = 0;
        const useTestHook = () => {
            useIsomorphicLayoutEffect(() => {
                callCount++;
            }, []);
        };

        renderHook(useTestHook);

        // Effect should be called once on mount
        expect(callCount).toBe(1);
    });
});
