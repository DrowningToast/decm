import { renderHook, act, waitFor } from "@testing-library/react";
import { useLocalStorage } from "./use-local-storage";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";

describe("useLocalStorage", () => {
    beforeEach(() => {
        // Clear localStorage before each test
        localStorage.clear();
        vi.clearAllMocks();
    });

    afterEach(() => {
        localStorage.clear();
    });

    it("should return initial value when localStorage is empty", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", "initial-value"));

        const [storedValue] = result.current;
        expect(storedValue).toBe("initial-value");
    });

    it("should return function initial value when localStorage is empty", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", () => "function-value"));

        const [storedValue] = result.current;
        expect(storedValue).toBe("function-value");
    });

    it("should persist value to localStorage", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", "initial"));

        act(() => {
            const [, setValue] = result.current;
            setValue("new-value");
        });

        expect(localStorage.getItem("test-key")).toBe(JSON.stringify("new-value"));
    });

    it("should update state when setValue is called", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", "initial"));

        act(() => {
            const [, setValue] = result.current;
            setValue("updated-value");
        });

        const [storedValue] = result.current;
        expect(storedValue).toBe("updated-value");
    });

    it("should handle setter function (like useState)", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", 0));

        act(() => {
            const [, setValue] = result.current;
            setValue((prev) => prev + 1);
        });

        expect(result.current[0]).toBe(1);

        act(() => {
            const [, setValue] = result.current;
            setValue((prev) => prev + 1);
        });

        expect(result.current[0]).toBe(2);
    });

    it("should serialize and deserialize objects", () => {
        interface TestObject {
            name: string;
            age: number;
        }

        const initialObj: TestObject = { name: "John", age: 30 };
        const { result } = renderHook(() => useLocalStorage("test-obj", initialObj));

        act(() => {
            const [, setValue] = result.current;
            setValue({ name: "Jane", age: 25 });
        });

        const [storedValue] = result.current;
        expect(storedValue).toEqual({ name: "Jane", age: 25 });
        expect(localStorage.getItem("test-obj")).toBe(JSON.stringify({ name: "Jane", age: 25 }));
    });

    it("should remove value from localStorage", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", "initial"));

        act(() => {
            const [, setValue] = result.current;
            setValue("stored-value");
        });

        expect(localStorage.getItem("test-key")).not.toBeNull();

        act(() => {
            const [, , removeValue] = result.current;
            removeValue();
        });

        expect(localStorage.getItem("test-key")).toBeNull();
        const [storedValue] = result.current;
        expect(storedValue).toBe("initial");
    });

    it("should restore initial value after removal", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", "initial-value"));

        act(() => {
            const [, setValue] = result.current;
            setValue("new-value");
        });

        expect(result.current[0]).toBe("new-value");

        act(() => {
            const [, , removeValue] = result.current;
            removeValue();
        });

        expect(result.current[0]).toBe("initial-value");
    });

    it("should use custom serializer", () => {
        const customSerializer = (value: string) => value.toUpperCase();
        const { result } = renderHook(() =>
            useLocalStorage("test-key", "initial", { serializer: customSerializer }),
        );

        act(() => {
            const [, setValue] = result.current;
            setValue("test");
        });

        expect(localStorage.getItem("test-key")).toBe("TEST");
    });

    it("should use custom deserializer", () => {
        localStorage.setItem("test-key", "STORED_VALUE");

        const customDeserializer = (value: string) => value.toLowerCase();
        const { result } = renderHook(() =>
            useLocalStorage("test-key", "initial", { deserializer: customDeserializer }),
        );

        const [storedValue] = result.current;
        expect(storedValue).toBe("stored_value");
    });

    it("should handle undefined as a special value", () => {
        localStorage.setItem("test-key", '"undefined"');

        const { result } = renderHook(() =>
            useLocalStorage<string | undefined>("test-key", "initial"),
        );

        const [storedValue] = result.current;
        expect(storedValue).toBeUndefined();
    });

    it("should handle invalid JSON gracefully", () => {
        localStorage.setItem("test-key", "invalid json");

        const { result } = renderHook(() => useLocalStorage("test-key", "fallback"));

        const [storedValue] = result.current;
        expect(storedValue).toBe("fallback");
    });

    it("should sync across hooks with same key", () => {
        const { result: result1 } = renderHook(() => useLocalStorage("shared-key", "initial"));
        const { result: result2 } = renderHook(() => useLocalStorage("shared-key", "initial"));

        expect(result1.current[0]).toBe(result2.current[0]);

        act(() => {
            const [, setValue] = result1.current;
            setValue("updated");
        });

        // Trigger a storage event to sync
        act(() => {
            window.dispatchEvent(new StorageEvent("local-storage", { key: "shared-key" }));
        });

        waitFor(() => {
            expect(result2.current[0]).toBe("updated");
        });
    });

    it("should not reinitialize when key changes but keeps watching key changes", () => {
        const { result, rerender } = renderHook(
            ({ key, value }: { key: string; value: string }) => useLocalStorage(key, value),
            { initialProps: { key: "key1", value: "initial1" } },
        );

        act(() => {
            const [, setValue] = result.current;
            setValue("stored-value");
        });

        expect(result.current[0]).toBe("stored-value");

        rerender({ key: "key2", value: "initial2" });

        expect(result.current[0]).toBe("initial2");
    });

    it("should handle initializeWithValue option", () => {
        localStorage.setItem("test-key", JSON.stringify("stored-value"));

        const { result } = renderHook(() =>
            useLocalStorage("test-key", "initial", { initializeWithValue: false }),
        );

        // Initially returns initialValue without reading from localStorage
        expect(result.current[0]).toBe("initial");

        // After effect runs, it should read from localStorage
        waitFor(() => {
            expect(result.current[0]).toBe("stored-value");
        });
    });

    it("should handle array values", () => {
        const { result } = renderHook(() => useLocalStorage<number[]>("test-array", [1, 2, 3]));

        act(() => {
            const [, setValue] = result.current;
            setValue([4, 5, 6]);
        });

        const [storedValue] = result.current;
        expect(storedValue).toEqual([4, 5, 6]);
    });

    it("should handle null values", () => {
        const { result } = renderHook(() => useLocalStorage<string | null>("test-key", null));

        act(() => {
            const [, setValue] = result.current;
            setValue("value");
        });

        expect(result.current[0]).toBe("value");

        act(() => {
            const [, setValue] = result.current;
            setValue(null);
        });

        expect(result.current[0]).toBeNull();
    });

    it("should dispatch custom storage event on setValue", () => {
        const dispatchEventSpy = vi.spyOn(window, "dispatchEvent");

        const { result } = renderHook(() => useLocalStorage("test-key", "initial"));

        act(() => {
            const [, setValue] = result.current;
            setValue("new-value");
        });

        expect(dispatchEventSpy).toHaveBeenCalled();
        const event = dispatchEventSpy.mock.calls.find(
            ([e]) => (e as StorageEvent).type === "local-storage",
        );
        expect(event).toBeDefined();
    });

    it("should handle removal with function initial value", () => {
        const { result } = renderHook(() => useLocalStorage("test-key", () => "default-value"));

        act(() => {
            const [, setValue] = result.current;
            setValue("stored-value");
        });

        act(() => {
            const [, , removeValue] = result.current;
            removeValue();
        });

        expect(result.current[0]).toBe("default-value");
    });
});
