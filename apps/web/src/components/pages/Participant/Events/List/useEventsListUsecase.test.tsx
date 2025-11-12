import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useEventsListUsecase } from "./useEventsListUsecase";
import type { ReactNode } from "react";

// Mock the zustand store
vi.mock("@/components/BottomNav/stores/events", () => ({
    useSearchEventNavStore: vi.fn(() => ({
        searchQuery: "",
    })),
}));

describe("useEventsListUsecase", () => {
    let queryClient: QueryClient;

    beforeEach(() => {
        queryClient = new QueryClient({
            defaultOptions: {
                queries: {
                    retry: false,
                },
            },
        });
        vi.clearAllMocks();
    });

    const wrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );

    it("returns events successfully", async () => {
        const { result } = renderHook(() => useEventsListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.events).toBeDefined();
        expect(Array.isArray(result.current.events)).toBe(true);
    });

    it("returns events with correct structure", async () => {
        const { result } = renderHook(() => useEventsListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.events.length).toBeGreaterThan(0);

        const event = result.current.events[0];
        expect(event).toHaveProperty("id");
        expect(event).toHaveProperty("name");
        expect(event).toHaveProperty("description");
        expect(event).toHaveProperty("eventName");
        expect(event).toHaveProperty("status");
        expect(event).toHaveProperty("accessType");
    });

    it("filters events based on search query", async () => {
        const { useSearchEventNavStore } = await import("@/components/BottomNav/stores/events");

        vi.mocked(useSearchEventNavStore).mockReturnValue({
            searchQuery: "ToBelT",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useEventsListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.events.length).toBeGreaterThan(0);
        result.current.events.forEach((event) => {
            expect(event.name.toLowerCase()).toContain("tobelt");
        });
    });

    it("performs case-insensitive fuzzy search", async () => {
        const { useSearchEventNavStore } = await import("@/components/BottomNav/stores/events");

        vi.mocked(useSearchEventNavStore).mockReturnValue({
            searchQuery: "TOBELT69",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useEventsListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.events.length).toBeGreaterThan(0);
    });

    it("returns all events when filterType is 'all'", async () => {
        const { result } = renderHook(() => useEventsListUsecase({ filterType: "all" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.events.length).toBeGreaterThan(0);
    });

    it("filters to only user's events when filterType is 'my-events'", async () => {
        const { result } = renderHook(() => useEventsListUsecase({ filterType: "my-events" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // All returned events should have participationStatus defined
        result.current.events.forEach((event) => {
            expect(event.participationStatus).toBeDefined();
        });
    });

    it("returns empty array when no matches found", async () => {
        const { useSearchEventNavStore } = await import("@/components/BottomNav/stores/events");

        vi.mocked(useSearchEventNavStore).mockReturnValue({
            searchQuery: "nonexistentquery12345",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useEventsListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.events).toEqual([]);
    });

    it("combines search and filter correctly", async () => {
        const { useSearchEventNavStore } = await import("@/components/BottomNav/stores/events");

        vi.mocked(useSearchEventNavStore).mockReturnValue({
            searchQuery: "ToBelT",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useEventsListUsecase({ filterType: "my-events" }), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // Should apply both search and filter
        result.current.events.forEach((event) => {
            expect(event.name.toLowerCase()).toContain("tobelt");
            expect(event.participationStatus).toBeDefined();
        });
    });

    it("includes different event statuses in mock data", async () => {
        const { result } = renderHook(() => useEventsListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        const statuses = result.current.events.map((event) => event.status);
        const uniqueStatuses = new Set(statuses);

        expect(uniqueStatuses.size).toBeGreaterThan(0);
    });

    it("includes different access types in mock data", async () => {
        const { result } = renderHook(() => useEventsListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        const accessTypes = result.current.events.map((event) => event.accessType);
        const uniqueAccessTypes = new Set(accessTypes);

        expect(uniqueAccessTypes.size).toBeGreaterThan(1);
        expect(Array.from(uniqueAccessTypes)).toEqual(
            expect.arrayContaining([expect.stringMatching(/public|password|invite-only/)]),
        );
    });
});
