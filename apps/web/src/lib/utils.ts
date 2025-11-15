import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export const delay = (ms: number): Promise<NodeJS.Timeout> => {
    return new Promise((resolve) => {
        const timeoutId = setTimeout(() => resolve(timeoutId), ms);
    });
};

export function formatEthereumAddress(address: string): string {
    if (address.length <= 10) {
        return address;
    }
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
}

/**
 * Result type for the until function
 */
export interface UntilResult<T> {
    /** Whether the callback succeeded */
    isSuccess: boolean;
    /** Whether the callback failed after all retry attempts */
    isFailed: boolean;
    /** The response data from the callback (null if failed) */
    response: T | null;
    /** The last error encountered (only present if failed) */
    error?: Error;
}

/**
 * Options for the until function
 */
export interface UntilOptions {
    /** Maximum duration to retry in milliseconds */
    maxDurationMs: number;
    /** Delay between retry attempts in milliseconds */
    delayMs: number;
}

/**
 * Retries a callback function until it succeeds or the max duration is reached
 *
 * @param callback - Async or sync function to retry
 * @param options - Configuration for max duration and delay interval
 * @returns Object with isSuccess, isFailed, response, and optional error
 *
 * @example
 * const result = await until(
 *   async () => {
 *     const data = await fetchData();
 *     return data;
 *   },
 *   { maxDurationMs: 5000, delayMs: 500 }
 * );
 *
 * if (result.isSuccess) {
 *   console.log(result.response);
 * } else {
 *   console.error(result.error);
 * }
 */
export async function until<T>(
    callback: () => T | Promise<T>,
    options: UntilOptions,
): Promise<UntilResult<T>> {
    const startTime = Date.now();
    let lastError: Error | null = null;

    while (Date.now() - startTime < options.maxDurationMs) {
        try {
            const response = await callback();
            return {
                isSuccess: true,
                isFailed: false,
                response,
            };
        } catch (error) {
            lastError = error instanceof Error ? error : new Error(String(error));
            await delay(options.delayMs);
        }
    }

    return {
        isSuccess: false,
        isFailed: true,
        response: null,
        error: lastError || new Error("Timeout: max duration reached"),
    };
}
