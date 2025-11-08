import { describe, it, expect, beforeEach } from "vitest";
import { useSearchEventNavStore } from "./events";

describe("useSearchEventNavStore", () => {
    beforeEach(() => {
        // Reset store state before each test
        useSearchEventNavStore.setState({ searchQuery: "" });
    });

    it("initializes with empty search query", () => {
        const state = useSearchEventNavStore.getState();
        expect(state.searchQuery).toBe("");
    });

    it("updates search query", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();

        setSearchQuery("event query");

        const state = useSearchEventNavStore.getState();
        expect(state.searchQuery).toBe("event query");
    });

    it("updates search query multiple times", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();

        setSearchQuery("first");
        expect(useSearchEventNavStore.getState().searchQuery).toBe("first");

        setSearchQuery("second");
        expect(useSearchEventNavStore.getState().searchQuery).toBe("second");

        setSearchQuery("third");
        expect(useSearchEventNavStore.getState().searchQuery).toBe("third");
    });

    it("can clear search query by setting to empty string", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();

        setSearchQuery("some event");
        expect(useSearchEventNavStore.getState().searchQuery).toBe("some event");

        setSearchQuery("");
        expect(useSearchEventNavStore.getState().searchQuery).toBe("");
    });

    it("handles special characters in search query", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();

        const specialQuery = "event@2024!";
        setSearchQuery(specialQuery);

        expect(useSearchEventNavStore.getState().searchQuery).toBe(specialQuery);
    });

    it("handles unicode characters in search query", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();

        const unicodeQuery = "งานเทศกาล 活动 حدث";
        setSearchQuery(unicodeQuery);

        expect(useSearchEventNavStore.getState().searchQuery).toBe(unicodeQuery);
    });

    it("preserves whitespace in search query", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();

        const queryWithSpaces = "  event   name  ";
        setSearchQuery(queryWithSpaces);

        expect(useSearchEventNavStore.getState().searchQuery).toBe(queryWithSpaces);
    });

    it("handles long search queries", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();

        const longQuery = "a".repeat(1000);
        setSearchQuery(longQuery);

        expect(useSearchEventNavStore.getState().searchQuery).toBe(longQuery);
    });

    it("allows subscribers to react to changes", () => {
        let callCount = 0;
        let lastValue = "";

        const unsubscribe = useSearchEventNavStore.subscribe((state) => {
            callCount++;
            lastValue = state.searchQuery;
        });

        const { setSearchQuery } = useSearchEventNavStore.getState();
        setSearchQuery("test event");

        expect(callCount).toBeGreaterThan(0);
        expect(lastValue).toBe("test event");

        unsubscribe();
    });

    it("state persists across multiple accesses", () => {
        const { setSearchQuery } = useSearchEventNavStore.getState();
        setSearchQuery("persistent query");

        // Access state multiple times
        const state1 = useSearchEventNavStore.getState();
        const state2 = useSearchEventNavStore.getState();
        const state3 = useSearchEventNavStore.getState();

        expect(state1.searchQuery).toBe("persistent query");
        expect(state2.searchQuery).toBe("persistent query");
        expect(state3.searchQuery).toBe("persistent query");
    });
});
