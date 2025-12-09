"use client";

import { useCallback, useRef } from "react";
import { useIsomorphicLayoutEffect } from "./use-isomorphic-layout-effect";

/**
 * Custom hook that creates a memoized event callback that's safe to call during rendering.
 *
 * This hook ensures that the callback always has access to the latest values
 * while maintaining a stable reference, preventing unnecessary re-renders.
 */
export function useEventCallback<Args extends unknown[], R>(
    fn: (...args: Args) => R,
): (...args: Args) => R;
export function useEventCallback<Args extends unknown[], R>(
    fn: ((...args: Args) => R) | undefined,
): ((...args: Args) => R) | undefined;
export function useEventCallback<Args extends unknown[], R>(
    fn: ((...args: Args) => R) | undefined,
): ((...args: Args) => R) | undefined {
    const ref = useRef<typeof fn>(() => {
        throw new Error("Cannot call an event handler while rendering.");
    });

    useIsomorphicLayoutEffect(() => {
        ref.current = fn;
    }, [fn]);

    // Always call useCallback at the top level, before any early returns
    const callback = useCallback(
        (...args: Args) => {
            if (ref.current === undefined) {
                return undefined as R;
            }
            return ref.current(...args);
        },
        [ref],
    ) as (...args: Args) => R;

    // If fn is undefined, return undefined
    if (fn === undefined) {
        return undefined;
    }

    return callback;
}
