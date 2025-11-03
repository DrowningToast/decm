import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useInboxListUsecase } from "./useInboxListUsecase";
import type { ReactNode } from "react";

// Mock the zustand store
vi.mock("@/components/BottomNav/stores/notifications", () => ({
    useSearchNotificationNavStore: vi.fn(() => ({
        searchQuery: "",
    })),
}));

describe("useInboxListUsecase", () => {
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

    it("returns inbox items successfully", async () => {
        const { result } = renderHook(() => useInboxListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxItems).toBeDefined();
        expect(Array.isArray(result.current.inboxItems)).toBe(true);
    });

    it("returns mock data with correct structure", async () => {
        const { result } = renderHook(() => useInboxListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxItems.length).toBeGreaterThan(0);

        const item = result.current.inboxItems[0];
        expect(item).toHaveProperty("id");
        expect(item).toHaveProperty("title");
        expect(item).toHaveProperty("sender");
        expect(item).toHaveProperty("date");
        expect(item).toHaveProperty("status");
    });

    it("filters inbox items based on search query", async () => {
        const { useSearchNotificationNavStore } = await import(
            "@/components/BottomNav/stores/notifications"
        );

        // Mock with search query
        vi.mocked(useSearchNotificationNavStore).mockReturnValue({
            searchQuery: "certificate",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useInboxListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // All filtered items should contain "certificate" in their title
        result.current.inboxItems.forEach((item) => {
            expect(item.title.toLowerCase()).toContain("certificate");
        });
    });

    it("performs case-insensitive fuzzy search", async () => {
        const { useSearchNotificationNavStore } = await import(
            "@/components/BottomNav/stores/notifications"
        );

        // Mock with search query in different case
        vi.mocked(useSearchNotificationNavStore).mockReturnValue({
            searchQuery: "INVITATION",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useInboxListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxItems.length).toBeGreaterThan(0);
    });

    it("searches by sender field", async () => {
        const { useSearchNotificationNavStore } = await import(
            "@/components/BottomNav/stores/notifications"
        );

        // Mock with sender search
        vi.mocked(useSearchNotificationNavStore).mockReturnValue({
            searchQuery: "ToBeIT",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useInboxListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxItems.length).toBeGreaterThan(0);
    });

    it("returns empty array when no matches found", async () => {
        const { useSearchNotificationNavStore } = await import(
            "@/components/BottomNav/stores/notifications"
        );

        // Mock with non-existent search query
        vi.mocked(useSearchNotificationNavStore).mockReturnValue({
            searchQuery: "nonexistentquery12345",
            setSearchQuery: vi.fn(),
        });

        const { result } = renderHook(() => useInboxListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.inboxItems).toEqual([]);
    });

    it("includes all status types in mock data", async () => {
        const { result } = renderHook(() => useInboxListUsecase(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        const statuses = result.current.inboxItems.map((item) => item.status);
        const uniqueStatuses = new Set(statuses);

        // Should have multiple different statuses
        expect(uniqueStatuses.size).toBeGreaterThan(1);
        expect(Array.from(uniqueStatuses)).toEqual(
            expect.arrayContaining([
                expect.stringMatching(/pending|available|expired|action-required/),
            ]),
        );
    });
});
