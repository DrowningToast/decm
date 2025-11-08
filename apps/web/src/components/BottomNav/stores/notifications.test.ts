import { describe, it, expect, beforeEach } from "vitest";
import { useSearchNotificationNavStore } from "./notifications";

describe("useSearchNotificationNavStore", () => {
    beforeEach(() => {
        // Reset store state before each test
        useSearchNotificationNavStore.setState({ searchQuery: "" });
    });

    it("initializes with empty search query", () => {
        const state = useSearchNotificationNavStore.getState();
        expect(state.searchQuery).toBe("");
    });

    it("updates search query", () => {
        const { setSearchQuery } = useSearchNotificationNavStore.getState();

        setSearchQuery("test query");

        const state = useSearchNotificationNavStore.getState();
        expect(state.searchQuery).toBe("test query");
    });

    it("updates search query multiple times", () => {
        const { setSearchQuery } = useSearchNotificationNavStore.getState();

        setSearchQuery("first");
        expect(useSearchNotificationNavStore.getState().searchQuery).toBe("first");

        setSearchQuery("second");
        expect(useSearchNotificationNavStore.getState().searchQuery).toBe("second");

        setSearchQuery("third");
        expect(useSearchNotificationNavStore.getState().searchQuery).toBe("third");
    });

    it("can clear search query by setting to empty string", () => {
        const { setSearchQuery } = useSearchNotificationNavStore.getState();

        setSearchQuery("some search");
        expect(useSearchNotificationNavStore.getState().searchQuery).toBe("some search");

        setSearchQuery("");
        expect(useSearchNotificationNavStore.getState().searchQuery).toBe("");
    });

    it("handles special characters in search query", () => {
        const { setSearchQuery } = useSearchNotificationNavStore.getState();

        const specialQuery = "test@#$%^&*()";
        setSearchQuery(specialQuery);

        expect(useSearchNotificationNavStore.getState().searchQuery).toBe(specialQuery);
    });

    it("handles unicode characters in search query", () => {
        const { setSearchQuery } = useSearchNotificationNavStore.getState();

        const unicodeQuery = "สวัสดี 你好 مرحبا";
        setSearchQuery(unicodeQuery);

        expect(useSearchNotificationNavStore.getState().searchQuery).toBe(unicodeQuery);
    });

    it("preserves whitespace in search query", () => {
        const { setSearchQuery } = useSearchNotificationNavStore.getState();

        const queryWithSpaces = "  test   query  ";
        setSearchQuery(queryWithSpaces);

        expect(useSearchNotificationNavStore.getState().searchQuery).toBe(queryWithSpaces);
    });

    it("allows subscribers to react to changes", () => {
        let callCount = 0;
        let lastValue = "";

        const unsubscribe = useSearchNotificationNavStore.subscribe((state) => {
            callCount++;
            lastValue = state.searchQuery;
        });

        const { setSearchQuery } = useSearchNotificationNavStore.getState();
        setSearchQuery("test");

        expect(callCount).toBeGreaterThan(0);
        expect(lastValue).toBe("test");

        unsubscribe();
    });
});
