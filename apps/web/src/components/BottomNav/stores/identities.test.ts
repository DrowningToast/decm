import { describe, it, expect, beforeEach } from "vitest";
import { useSearchIdentitiesNavStore } from "./identities";

describe("useSearchIdentitiesNavStore", () => {
    beforeEach(() => {
        useSearchIdentitiesNavStore.setState({ searchQuery: "" });
    });

    it("initializes with empty search query", () => {
        expect(useSearchIdentitiesNavStore.getState().searchQuery).toBe("");
    });

    it("updates search query", () => {
        useSearchIdentitiesNavStore.getState().setSearchQuery("alice");
        expect(useSearchIdentitiesNavStore.getState().searchQuery).toBe("alice");
    });

    it("clears search query", () => {
        useSearchIdentitiesNavStore.getState().setSearchQuery("alice");
        useSearchIdentitiesNavStore.getState().setSearchQuery("");
        expect(useSearchIdentitiesNavStore.getState().searchQuery).toBe("");
    });
});
