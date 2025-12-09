import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { cn, delay, formatEthereumAddress, until } from "./utils";

describe("cn utility", () => {
    it("should merge class names correctly", () => {
        expect(cn("class1", "class2")).toBe("class1 class2");
    });

    it("should handle conditional classes", () => {
        const isFalsy = false;
        const isTruthy = true;
        expect(cn("base", isFalsy && "conditional", "always")).toBe("base always");
        expect(cn("base", isTruthy && "conditional", "always")).toBe("base conditional always");
    });

    it("should merge tailwind classes correctly", () => {
        expect(cn("px-2 py-1", "px-4")).toBe("py-1 px-4");
    });

    it("should handle arrays of classes", () => {
        expect(cn(["class1", "class2"], "class3")).toBe("class1 class2 class3");
    });

    it("should handle undefined and null values", () => {
        expect(cn("class1", undefined, "class2", null, "class3")).toBe("class1 class2 class3");
    });

    it("should handle empty strings", () => {
        expect(cn("class1", "", "class2")).toBe("class1 class2");
    });

    it("should merge duplicate classes correctly", () => {
        // cn doesn't deduplicate classes by default, it just merges them
        // The behavior depends on clsx/twMerge implementation
        const result = cn("class1 class2", "class1 class3");
        expect(result).toContain("class1");
        expect(result).toContain("class2");
        expect(result).toContain("class3");
    });

    it("should handle complex tailwind conflicts", () => {
        expect(cn("text-red-500 bg-blue-500", "text-green-500")).toBe("bg-blue-500 text-green-500");
    });

    it("should work with no arguments", () => {
        expect(cn()).toBe("");
    });

    it("should handle objects", () => {
        expect(cn({ class1: true, class2: false, class3: true })).toBe("class1 class3");
    });
});

describe("delay utility", () => {
    beforeEach(() => {
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it("should delay execution by specified milliseconds", async () => {
        const promise = delay(1000);

        vi.advanceTimersByTime(999);
        expect(promise).toBeInstanceOf(Promise);

        vi.advanceTimersByTime(1);
        await expect(promise).resolves.toBeUndefined();
    });

    it("should resolve immediately with 0ms delay", async () => {
        const promise = delay(0);

        vi.advanceTimersByTime(0);
        await expect(promise).resolves.toBeUndefined();
    });

    it("should handle multiple delays", async () => {
        const promise1 = delay(100);
        const promise2 = delay(200);
        const promise3 = delay(300);

        vi.advanceTimersByTime(100);
        await expect(promise1).resolves.toBeUndefined();

        vi.advanceTimersByTime(100);
        await expect(promise2).resolves.toBeUndefined();

        vi.advanceTimersByTime(100);
        await expect(promise3).resolves.toBeUndefined();
    });

    it("should work in async functions", async () => {
        const callback = vi.fn();

        const asyncFunction = async () => {
            callback("start");
            await delay(500);
            callback("end");
        };

        const promise = asyncFunction();

        expect(callback).toHaveBeenCalledWith("start");
        expect(callback).not.toHaveBeenCalledWith("end");

        vi.advanceTimersByTime(500);
        await promise;

        expect(callback).toHaveBeenCalledWith("end");
        expect(callback).toHaveBeenCalledTimes(2);
    });
});

describe("formatEthereumAddress utility", () => {
    it("should format ethereum address correctly", () => {
        const address = "0x1234567890123456789012345678901234567890";
        const result = formatEthereumAddress(address);
        expect(result).toBe("0x1234...7890");
    });

    it("should return short address as-is when length is <= 10", () => {
        const shortAddress = "0x12345678";
        const result = formatEthereumAddress(shortAddress);
        expect(result).toBe("0x12345678");
    });

    it("should handle exactly 10 character address", () => {
        const address = "0x12345678";
        const result = formatEthereumAddress(address);
        expect(result).toBe("0x12345678");
    });

    it("should handle 11 character address", () => {
        const address = "0x123456789";
        const result = formatEthereumAddress(address);
        expect(result).toBe("0x1234...6789");
    });

    it("should format standard ethereum address (42 chars)", () => {
        const address = "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb";
        const result = formatEthereumAddress(address);
        expect(result).toBe("0x742d...0bEb");
    });

    it("should handle very long addresses", () => {
        const address = "0x1234567890123456789012345678901234567890123456789012345678901234567890";
        const result = formatEthereumAddress(address);
        expect(result).toBe("0x1234...7890");
    });
});

describe("until utility", () => {
    it("should return success when callback succeeds on first attempt", async () => {
        const callback = vi.fn().mockResolvedValueOnce({ data: "success" });

        const result = await until(callback, { maxDurationMs: 5000, delayMs: 100 });

        expect(result.isSuccess).toBe(true);
        expect(result.isFailed).toBe(false);
        expect(result.response).toEqual({ data: "success" });
        expect(result.error).toBeUndefined();
        expect(callback).toHaveBeenCalledTimes(1);
    });

    it("should support synchronous callbacks", async () => {
        const callback = vi.fn().mockReturnValueOnce("sync success");

        const result = await until(callback, { maxDurationMs: 5000, delayMs: 100 });

        expect(result.isSuccess).toBe(true);
        expect(result.response).toBe("sync success");
        expect(callback).toHaveBeenCalledTimes(1);
    });

    it("should handle zero max duration", async () => {
        const callback = vi.fn().mockResolvedValueOnce("success");

        const result = await until(callback, { maxDurationMs: 0, delayMs: 100 });

        expect(result.isSuccess).toBe(false);
        expect(result.isFailed).toBe(true);
        expect(result.error?.message).toBe("Timeout: max duration reached");
    });

    it("should handle return correct types on success", async () => {
        const callback = vi.fn().mockResolvedValueOnce({ success: true, value: 42 });

        const result = await until(callback, { maxDurationMs: 5000, delayMs: 100 });

        expect(result).toHaveProperty("isSuccess");
        expect(result).toHaveProperty("isFailed");
        expect(result).toHaveProperty("response");
        expect(result.isSuccess).toBe(true);
        expect(result.response?.success).toBe(true);
    });

    it("should pass through response without error property on success", async () => {
        const callback = vi.fn().mockResolvedValueOnce({ data: "test" });

        const result = await until(callback, { maxDurationMs: 5000, delayMs: 100 });

        expect(result.isSuccess).toBe(true);
        expect(result.isFailed).toBe(false);
        expect(result.response).toEqual({ data: "test" });
        // On success, error is undefined (not included in the object)
        expect(result.error).toBeUndefined();
    });

    it("should retry on failure until max duration", async () => {
        vi.useFakeTimers();
        const callback = vi
            .fn()
            .mockRejectedValueOnce(new Error("First attempt failed"))
            .mockRejectedValueOnce(new Error("Second attempt failed"))
            .mockResolvedValueOnce("Success on third attempt");

        const promise = until(callback, { maxDurationMs: 1000, delayMs: 100 });

        // Advance timers to allow retries
        await vi.advanceTimersByTimeAsync(300);

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.response).toBe("Success on third attempt");
        expect(callback).toHaveBeenCalledTimes(3);
        vi.useRealTimers();
    });

    it("should fail after max duration is reached", async () => {
        vi.useFakeTimers();
        const error = new Error("Always fails");
        const callback = vi.fn().mockRejectedValue(error);

        const promise = until(callback, { maxDurationMs: 500, delayMs: 100 });

        // Advance timers past max duration
        await vi.advanceTimersByTimeAsync(600);

        const result = await promise;

        expect(result.isSuccess).toBe(false);
        expect(result.isFailed).toBe(true);
        expect(result.response).toBeNull();
        expect(result.error).toBe(error);
        expect(callback).toHaveBeenCalled();
        vi.useRealTimers();
    });

    it("should handle non-Error exceptions", async () => {
        vi.useFakeTimers();
        const callback = vi.fn().mockRejectedValue("String error");

        const promise = until(callback, { maxDurationMs: 100, delayMs: 50 });

        await vi.advanceTimersByTimeAsync(150);

        const result = await promise;

        expect(result.isFailed).toBe(true);
        expect(result.error).toBeInstanceOf(Error);
        expect(result.error?.message).toBe("String error");
        vi.useRealTimers();
    });
});
