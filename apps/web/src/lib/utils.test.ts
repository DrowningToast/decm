import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { cn, delay } from "./utils";

describe("cn utility", () => {
    it("should merge class names correctly", () => {
        expect(cn("class1", "class2")).toBe("class1 class2");
    });

    it("should handle conditional classes", () => {
        const isFalse = false;
        const isTrue = true;
        expect(cn("base", isFalse && "conditional", "always")).toBe("base always");
        expect(cn("base", isTrue && "conditional", "always")).toBe("base conditional always");
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
