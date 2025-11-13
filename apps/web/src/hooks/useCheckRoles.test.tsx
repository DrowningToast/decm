import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useCheckRoles } from "./useCheckRoles";
import { authService } from "@/services/AuthService/AuthService";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";

// Mock the authService
vi.mock("@/services/AuthService/AuthService", () => ({
    authService: {
        checkRoles: vi.fn(),
    },
}));

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                retry: false,
            },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useCheckRoles", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should return hasRequiredRoles as true when user has host role and it is required", async () => {
        const mockResponse = { is_host: true, is_issuer: false };
        vi.mocked(authService.checkRoles).mockResolvedValue(mockResponse);

        const { result } = renderHook(
            () => useCheckRoles({ requireHost: true, requireIssuer: false }),
            { wrapper: createWrapper() },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.hasRequiredRoles).toBe(true);
        expect(result.current.isHost).toBe(true);
        expect(result.current.isIssuer).toBe(false);
    });

    it("should return hasRequiredRoles as true when user has issuer role and it is required", async () => {
        const mockResponse = { is_host: false, is_issuer: true };
        vi.mocked(authService.checkRoles).mockResolvedValue(mockResponse);

        const { result } = renderHook(
            () => useCheckRoles({ requireHost: false, requireIssuer: true }),
            { wrapper: createWrapper() },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.hasRequiredRoles).toBe(true);
        expect(result.current.isIssuer).toBe(true);
    });

    it("should return hasRequiredRoles as true when user has both roles", async () => {
        const mockResponse = { is_host: true, is_issuer: true };
        vi.mocked(authService.checkRoles).mockResolvedValue(mockResponse);

        const { result } = renderHook(
            () => useCheckRoles({ requireHost: true, requireIssuer: true }),
            { wrapper: createWrapper() },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.hasRequiredRoles).toBe(true);
        expect(result.current.isHost).toBe(true);
        expect(result.current.isIssuer).toBe(true);
    });

    it("should return hasRequiredRoles as false when user lacks required host role", async () => {
        const mockResponse = { is_host: false, is_issuer: true };
        vi.mocked(authService.checkRoles).mockResolvedValue(mockResponse);

        const { result } = renderHook(
            () => useCheckRoles({ requireHost: true, requireIssuer: false }),
            { wrapper: createWrapper() },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.hasRequiredRoles).toBe(false);
        expect(result.current.isHost).toBe(false);
    });

    it("should return hasRequiredRoles as false when user lacks required issuer role", async () => {
        const mockResponse = { is_host: true, is_issuer: false };
        vi.mocked(authService.checkRoles).mockResolvedValue(mockResponse);

        const { result } = renderHook(
            () => useCheckRoles({ requireHost: false, requireIssuer: true }),
            { wrapper: createWrapper() },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.hasRequiredRoles).toBe(false);
        expect(result.current.isIssuer).toBe(false);
    });

    it("should handle when both roles are not required", async () => {
        const mockResponse = { is_host: false, is_issuer: false };
        vi.mocked(authService.checkRoles).mockResolvedValue(mockResponse);

        const { result } = renderHook(
            () => useCheckRoles({ requireHost: false, requireIssuer: false }),
            { wrapper: createWrapper() },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.hasRequiredRoles).toBe(true);
    });

    it("should not fetch when enabled is false", async () => {
        const { result } = renderHook(() => useCheckRoles({ enabled: false }), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(authService.checkRoles).not.toHaveBeenCalled();
    });

    it("should handle errors gracefully", async () => {
        const error = new Error("Failed to check roles");
        vi.mocked(authService.checkRoles).mockRejectedValue(error);

        const { result } = renderHook(() => useCheckRoles({ requireHost: true }), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.isError).toBe(true);
        expect(result.current.error).toBeDefined();
        expect(result.current.hasRequiredRoles).toBe(false);
    });

    it("should call authService.checkRoles with correct parameters", async () => {
        vi.mocked(authService.checkRoles).mockResolvedValue({
            is_host: true,
            is_issuer: false,
        });

        renderHook(() => useCheckRoles({ requireHost: true, requireIssuer: false }), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(authService.checkRoles).toHaveBeenCalledWith({
                requireHost: true,
                requireIssuer: false,
            });
        });
    });
});
