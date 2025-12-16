import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useMarkInboxMessageAsRead } from "./useMarkInboxMessageAsRead";
import { QUERY_KEY } from "@/lib/queryKeys";

// Mock the InboxService
vi.mock("@/services/InboxService/InboxService", () => ({
    defaultInboxService: {
        markAsRead: vi.fn(),
    },
}));

import { defaultInboxService } from "@/services/InboxService/InboxService";

describe("useMarkInboxMessageAsRead", () => {
    let queryClient: QueryClient;

    beforeEach(() => {
        vi.clearAllMocks();
        queryClient = new QueryClient({
            defaultOptions: {
                queries: {
                    retry: false,
                },
                mutations: {
                    retry: false,
                },
            },
        });
    });

    const wrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );

    it("should mark message as read successfully", async () => {
        // Arrange
        const messageId = "message-123";
        vi.mocked(defaultInboxService.markAsRead).mockResolvedValue(undefined);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert
        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        expect(defaultInboxService.markAsRead).toHaveBeenCalledWith(messageId);
        expect(defaultInboxService.markAsRead).toHaveBeenCalledTimes(1);
    });

    it("should invalidate inbox list query on success", async () => {
        // Arrange
        const messageId = "message-456";
        vi.mocked(defaultInboxService.markAsRead).mockResolvedValue(undefined);

        const invalidateQueriesSpy = vi.spyOn(queryClient, "invalidateQueries");

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert
        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        expect(invalidateQueriesSpy).toHaveBeenCalledWith({
            queryKey: QUERY_KEY.inbox.list(),
        });
    });

    it("should invalidate specific message query on success", async () => {
        // Arrange
        const messageId = "message-789";
        vi.mocked(defaultInboxService.markAsRead).mockResolvedValue(undefined);

        const invalidateQueriesSpy = vi.spyOn(queryClient, "invalidateQueries");

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert
        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        expect(invalidateQueriesSpy).toHaveBeenCalledWith({
            queryKey: QUERY_KEY.inbox.byId(messageId),
        });
    });

    it("should handle error when marking as read fails", async () => {
        // Arrange
        const messageId = "message-error";
        const error = new Error("Failed to mark as read");
        vi.mocked(defaultInboxService.markAsRead).mockRejectedValue(error);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert
        await waitFor(() => {
            expect(result.current.isError).toBe(true);
        });

        expect(result.current.error).toBeDefined();
        expect(defaultInboxService.markAsRead).toHaveBeenCalledWith(messageId);
    });

    it("should have correct initial state", () => {
        // Arrange & Act
        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Assert
        expect(result.current.isIdle).toBe(true);
        expect(result.current.isPending).toBe(false);
        expect(result.current.isSuccess).toBe(false);
        expect(result.current.isError).toBe(false);
    });

    it("should transition through loading state", async () => {
        // Arrange
        const messageId = "message-loading";
        let resolveMarkAsRead: () => void;
        const markAsReadPromise = new Promise<void>((resolve) => {
            resolveMarkAsRead = resolve;
        });

        vi.mocked(defaultInboxService.markAsRead).mockReturnValue(markAsReadPromise);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert - should be loading
        await waitFor(() => {
            expect(result.current.isPending).toBe(true);
        });

        // Resolve the promise
        resolveMarkAsRead!();

        // Assert - should be success
        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });
    });

    it("should allow mutating multiple messages sequentially", async () => {
        // Arrange
        const messageId1 = "message-1";
        const messageId2 = "message-2";
        vi.mocked(defaultInboxService.markAsRead).mockResolvedValue(undefined);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act - First mutation
        result.current.mutate(messageId1);

        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        // Act - Second mutation
        result.current.mutate(messageId2);

        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        // Assert
        expect(defaultInboxService.markAsRead).toHaveBeenCalledWith(messageId1);
        expect(defaultInboxService.markAsRead).toHaveBeenCalledWith(messageId2);
        expect(defaultInboxService.markAsRead).toHaveBeenCalledTimes(2);
    });

    it("should support mutateAsync for promise-based usage", async () => {
        // Arrange
        const messageId = "message-async";
        vi.mocked(defaultInboxService.markAsRead).mockResolvedValue(undefined);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        const promise = result.current.mutateAsync(messageId);

        // Assert
        await expect(promise).resolves.toBeUndefined();
        expect(defaultInboxService.markAsRead).toHaveBeenCalledWith(messageId);
    });

    it("should reject mutateAsync on error", async () => {
        // Arrange
        const messageId = "message-async-error";
        const error = new Error("Async error");
        vi.mocked(defaultInboxService.markAsRead).mockRejectedValue(error);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        const promise = result.current.mutateAsync(messageId);

        // Assert
        await expect(promise).rejects.toThrow("Async error");
    });

    it("should reset mutation state", async () => {
        // Arrange
        const messageId = "message-reset";
        vi.mocked(defaultInboxService.markAsRead).mockResolvedValue(undefined);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act - Mutate
        result.current.mutate(messageId);

        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        // Act - Reset
        result.current.reset();

        // Assert
        await waitFor(() => {
            expect(result.current.isIdle).toBe(true);
            expect(result.current.isSuccess).toBe(false);
        });
    });

    it("should not invalidate queries if mutation fails", async () => {
        // Arrange
        const messageId = "message-fail";
        vi.mocked(defaultInboxService.markAsRead).mockRejectedValue(new Error("Failed"));

        const invalidateQueriesSpy = vi.spyOn(queryClient, "invalidateQueries");

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert
        await waitFor(() => {
            expect(result.current.isError).toBe(true);
        });

        // onSuccess should not have been called, so no invalidations
        expect(invalidateQueriesSpy).not.toHaveBeenCalled();
    });

    it("should handle network timeout error", async () => {
        // Arrange
        const messageId = "message-timeout";
        const timeoutError = new Error("Network timeout");
        vi.mocked(defaultInboxService.markAsRead).mockRejectedValue(timeoutError);

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert
        await waitFor(() => {
            expect(result.current.isError).toBe(true);
        });

        expect(result.current.error).toEqual(timeoutError);
    });

    it("should invalidate both queries in correct order", async () => {
        // Arrange
        const messageId = "message-order";
        vi.mocked(defaultInboxService.markAsRead).mockResolvedValue(undefined);

        const invalidateCalls: string[] = [];
        vi.spyOn(queryClient, "invalidateQueries").mockImplementation((options: any) => {
            invalidateCalls.push(JSON.stringify(options.queryKey));
            return Promise.resolve();
        });

        const { result } = renderHook(() => useMarkInboxMessageAsRead(), { wrapper });

        // Act
        result.current.mutate(messageId);

        // Assert
        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        expect(invalidateCalls).toHaveLength(2);
        expect(invalidateCalls[0]).toContain(messageId); // Specific message query first
        expect(invalidateCalls[1]).toContain("inbox"); // Then list query
    });
});