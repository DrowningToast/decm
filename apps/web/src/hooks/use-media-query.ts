import { useState } from "react";

import { useIsomorphicLayoutEffect } from "./use-isomorphic-layout-effect";

interface UseMediaQueryOptions {
    defaultValue?: boolean;
    initializeWithValue?: boolean;
}

const IS_SERVER = typeof window === "undefined";

/**
 * SSR-aware React hook for media queries
 *
 * @param query - The media query string (e.g., "(max-width: 768px)")
 * @param options - Configuration options
 * @param options.defaultValue - Default value to use on server or when initializeWithValue is false (default: false)
 * @param options.initializeWithValue - Whether to initialize state with actual match value on client (default: true)
 *
 * @returns boolean indicating whether the media query matches the current viewport
 *
 * @description
 * This hook provides an SSR-safe way to use media queries in React applications.
 * - On the server, it returns the defaultValue
 * - On the client, it uses window.matchMedia to get the current match state
 * - Sets up a listener that updates the state when the media query match status changes
 * - Falls back to deprecated addListener/removeListener for Safari < 14 compatibility
 *
 * @example
 * const isMobile = useMediaQuery("(max-width: 768px)");
 * const isTablet = useMediaQuery("(min-width: 768px) and (max-width: 1024px)", { defaultValue: false });
 */
export function useMediaQuery(
    query: string,
    { defaultValue = false, initializeWithValue = true }: UseMediaQueryOptions = {},
): boolean {
    const getMatches = (query: string): boolean => {
        if (IS_SERVER) {
            return defaultValue;
        }
        return window.matchMedia(query).matches;
    };

    const [matches, setMatches] = useState<boolean>(() => {
        if (initializeWithValue) {
            return getMatches(query);
        }
        return defaultValue;
    });

    // Handles the change event of the media query.
    function handleChange() {
        setMatches(getMatches(query));
    }

    useIsomorphicLayoutEffect(() => {
        const matchMedia = window.matchMedia(query);

        // Triggered at the first client-side load and if query changes
        handleChange();

        // Use deprecated `addListener` and `removeListener` to support Safari < 14 (#135)
        if (matchMedia.addListener) {
            matchMedia.addListener(handleChange);
        } else {
            matchMedia.addEventListener("change", handleChange);
        }

        return () => {
            if (matchMedia.removeListener) {
                matchMedia.removeListener(handleChange);
            } else {
                matchMedia.removeEventListener("change", handleChange);
            }
        };
    }, [query]);

    return matches;
}

export type { UseMediaQueryOptions };
