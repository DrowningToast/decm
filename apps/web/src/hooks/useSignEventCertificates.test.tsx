import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useSignEventCertificates } from "./useSignEventCertificates";
import { coreApiClient } from "@/lib/api/api";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import { toast } from "sonner";

// Mock the API client
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            signEventCertificates: vi.fn(),
        },
    },
}));

// Mock react-i18next
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string, options?: { count: number }) => {
            if (key === "issuer.sign.signingSuccess") {
                return `Successfully signed ${options?.count || 0} certificates`;
            }
            if (key === "issuer.sign.signingError") {
                return "Failed to sign certificates";
            }
            return key;
        },
    }),
}));

// Mock sonner toast
vi.mock("sonner", () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    },
}));

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                retry: false,
            },
            mutations: {
                retry: false,
            },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useSignEventCertificates", () => {
    const mockEventId = "event-123";
    const mockIssuerPin = "1234";
    const mockResponse = {
        certificates: [
            { id: "cert-1", event_id: mockEventId },
            { id: "cert-2", event_id: mockEventId },
        ],
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should initialize with correct default values", () => {
        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        expect(result.current.isSigning).toBe(false);
        expect(result.current.signingError).toBeNull();
        expect(result.current.signEventCertificates).toBeDefined();
    });

    it("should sign certificates successfully", async () => {
        vi.mocked(coreApiClient.v1.signEventCertificates).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        result.current.signEventCertificates({
            eventId: mockEventId,
            issuerPin: mockIssuerPin,
        });

        await waitFor(() => {
            expect(result.current.isSigning).toBe(false);
        });

        expect(toast.success).toHaveBeenCalledWith("Successfully signed 2 certificates");
    });

    it("should call API with correct parameters", async () => {
        vi.mocked(coreApiClient.v1.signEventCertificates).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        result.current.signEventCertificates({
            eventId: mockEventId,
            issuerPin: mockIssuerPin,
        });

        await waitFor(() => {
            expect(result.current.isSigning).toBe(false);
        });

        expect(coreApiClient.v1.signEventCertificates).toHaveBeenCalledWith(
            { eventId: mockEventId },
            { issuer_pin: mockIssuerPin },
        );
    });

    it("should handle signing errors", async () => {
        const error = new Error("Invalid PIN");
        vi.mocked(coreApiClient.v1.signEventCertificates).mockRejectedValue(error);

        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        result.current.signEventCertificates({
            eventId: mockEventId,
            issuerPin: mockIssuerPin,
        });

        await waitFor(() => {
            expect(result.current.isSigning).toBe(false);
        });

        expect(toast.error).toHaveBeenCalledWith("Failed to sign certificates");
        expect(result.current.signingError).toBeDefined();
    });

    it("should handle empty certificates response", async () => {
        const emptyResponse = {
            certificates: [],
        };
        vi.mocked(coreApiClient.v1.signEventCertificates).mockResolvedValue(emptyResponse);

        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        result.current.signEventCertificates({
            eventId: mockEventId,
            issuerPin: mockIssuerPin,
        });

        await waitFor(() => {
            expect(result.current.isSigning).toBe(false);
        });

        expect(toast.success).toHaveBeenCalledWith("Successfully signed 0 certificates");
    });

    it("should handle null certificates response", async () => {
        const nullResponse = {
            certificates: null,
        };
        vi.mocked(coreApiClient.v1.signEventCertificates).mockResolvedValue(nullResponse);

        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        result.current.signEventCertificates({
            eventId: mockEventId,
            issuerPin: mockIssuerPin,
        });

        await waitFor(() => {
            expect(result.current.isSigning).toBe(false);
        });

        expect(toast.success).toHaveBeenCalledWith("Successfully signed 0 certificates");
    });

    it("should set isSigning to true during mutation", async () => {
        vi.mocked(coreApiClient.v1.signEventCertificates).mockImplementation(
            () => new Promise((resolve) => setTimeout(() => resolve(mockResponse), 100)),
        );

        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        result.current.signEventCertificates({
            eventId: mockEventId,
            issuerPin: mockIssuerPin,
        });

        // Check that isSigning is true during mutation
        await waitFor(() => {
            expect(result.current.isSigning).toBe(true);
        });

        await waitFor(() => {
            expect(result.current.isSigning).toBe(false);
        });
    });

    it("should handle multiple signature requests sequentially", async () => {
        vi.mocked(coreApiClient.v1.signEventCertificates).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useSignEventCertificates(), {
            wrapper: createWrapper(),
        });

        // First request
        result.current.signEventCertificates({
            eventId: mockEventId,
            issuerPin: mockIssuerPin,
        });

        await waitFor(() => {
            expect(result.current.isSigning).toBe(false);
        });

        expect(coreApiClient.v1.signEventCertificates).toHaveBeenCalledTimes(1);

        // Second request
        result.current.signEventCertificates({
            eventId: "event-456",
            issuerPin: "5678",
        });

        await waitFor(() => {
            expect(coreApiClient.v1.signEventCertificates).toHaveBeenCalledTimes(2);
        });
    });
});
