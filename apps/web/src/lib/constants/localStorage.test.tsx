import { describe, it, expect, beforeEach, afterEach } from "vitest";
import {
    LOCAL_STORAGE_KEYS,
    getLocalStorageItem,
    setLocalStorageItem,
    removeLocalStorageItem,
} from "./localStorage";

describe("localStorage utils", () => {
    const mockLocalStorage = (() => {
        let store: Record<string, string> = {};

        return {
            getItem: (key: string) => store[key] || null,
            setItem: (key: string, value: string) => {
                store[key] = value === undefined ? "undefined" : value.toString();
            },
            removeItem: (key: string) => {
                delete store[key];
            },
            clear: () => {
                store = {};
            },
        };
    })();

    beforeEach(() => {
        Object.defineProperty(window, "localStorage", {
            value: mockLocalStorage,
            writable: true,
        });
        mockLocalStorage.clear();
    });

    afterEach(() => {
        mockLocalStorage.clear();
    });

    describe("getLocalStorageItem", () => {
        it("should return undefined when key does not exist", () => {
            const result = getLocalStorageItem(LOCAL_STORAGE_KEYS.JWT);
            expect(result).toBeUndefined();
        });

        it("should return string value when key exists", () => {
            mockLocalStorage.setItem(LOCAL_STORAGE_KEYS.JWT, "test-jwt-token");
            const result = getLocalStorageItem(LOCAL_STORAGE_KEYS.JWT);
            expect(result).toBe("test-jwt-token");
        });

        it("should parse and return number for EXPIRES_IN", () => {
            mockLocalStorage.setItem(LOCAL_STORAGE_KEYS.EXPIRES_IN, "3600");
            const result = getLocalStorageItem(LOCAL_STORAGE_KEYS.EXPIRES_IN);
            expect(result).toBe(3600);
        });

        it("should parse JSON for EXPIRES_IN when value is JSON string", () => {
            mockLocalStorage.setItem(LOCAL_STORAGE_KEYS.EXPIRES_IN, JSON.stringify(7200));
            const result = getLocalStorageItem(LOCAL_STORAGE_KEYS.EXPIRES_IN);
            expect(result).toBe(7200);
        });

        it("should handle invalid JSON for EXPIRES_IN gracefully", () => {
            mockLocalStorage.setItem(LOCAL_STORAGE_KEYS.EXPIRES_IN, "invalid-json");
            const result = getLocalStorageItem(LOCAL_STORAGE_KEYS.EXPIRES_IN);
            expect(result).toBeUndefined();
        });

        it("should handle NaN for EXPIRES_IN", () => {
            mockLocalStorage.setItem(LOCAL_STORAGE_KEYS.EXPIRES_IN, "not-a-number");
            const result = getLocalStorageItem(LOCAL_STORAGE_KEYS.EXPIRES_IN);
            expect(result).toBeUndefined();
        });

        it("should return valid number string as number for EXPIRES_IN", () => {
            mockLocalStorage.setItem(LOCAL_STORAGE_KEYS.EXPIRES_IN, "12345");
            const result = getLocalStorageItem(LOCAL_STORAGE_KEYS.EXPIRES_IN);
            expect(result).toBe(12345);
        });
    });

    describe("setLocalStorageItem", () => {
        it("should set string value", () => {
            setLocalStorageItem(LOCAL_STORAGE_KEYS.JWT, "test-jwt");
            expect(mockLocalStorage.getItem(LOCAL_STORAGE_KEYS.JWT)).toBe("test-jwt");
        });

        it("should stringify number value", () => {
            setLocalStorageItem(LOCAL_STORAGE_KEYS.EXPIRES_IN, 3600);
            expect(mockLocalStorage.getItem(LOCAL_STORAGE_KEYS.EXPIRES_IN)).toBe("3600");
        });

        it("should stringify object value", () => {
            const obj = { test: "value" };
            setLocalStorageItem(
                LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT,
                JSON.stringify(obj),
            );
            expect(
                mockLocalStorage.getItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT),
            ).toBe(JSON.stringify(obj));
        });

        it("should handle undefined value", () => {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            setLocalStorageItem(LOCAL_STORAGE_KEYS.JWT, undefined as any);
            expect(mockLocalStorage.getItem(LOCAL_STORAGE_KEYS.JWT)).toBe("undefined");
        });
    });

    describe("removeLocalStorageItem", () => {
        it("should remove item from localStorage", () => {
            mockLocalStorage.setItem(LOCAL_STORAGE_KEYS.JWT, "test-jwt");
            expect(mockLocalStorage.getItem(LOCAL_STORAGE_KEYS.JWT)).toBe("test-jwt");

            removeLocalStorageItem(LOCAL_STORAGE_KEYS.JWT);
            expect(mockLocalStorage.getItem(LOCAL_STORAGE_KEYS.JWT)).toBeNull();
        });

        it("should not throw when removing non-existent item", () => {
            expect(() => {
                removeLocalStorageItem(LOCAL_STORAGE_KEYS.JWT);
            }).not.toThrow();
        });
    });

    describe("LOCAL_STORAGE_KEYS", () => {
        it("should have all expected keys", () => {
            expect(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT).toBe(
                "on_google_oauth_success_redirect",
            );
            expect(LOCAL_STORAGE_KEYS.JWT).toBe("jwt");
            expect(LOCAL_STORAGE_KEYS.ACCESS_TOKEN).toBe("accessToken");
            expect(LOCAL_STORAGE_KEYS.EXPIRES_IN).toBe("expiresIn");
            expect(LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE).toBe("authSignSignature");
        });
    });
});
