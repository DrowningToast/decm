import { renderHook, waitFor, act } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useClaimCertificate } from "./useClaimCertificate";
import { certificateService } from "@/services/services";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import { toast } from "sonner";

// Mock the certificate service
vi.mock("@/services/services", () => ({
    certificateService: {
        claimCertificate: vi.fn(),
    },
}));

// Mock react-i18next
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string, fallback?: string) => fallback ?? key,
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
            queries: { retry: false },
            mutations: { retry: false },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

/** Synchronous-style response (old backend behavior or already-minted) */
const makeSuccessResponse = (overrides = {}) => ({
    certificateId: "cert-123",
    transactionHash: "0xabc",
    claimedAt: "2024-01-01T00:00:00Z",
    message: "Certificate claimed successfully",
    status: "success",
    ...overrides,
});

/**
 * Queued response shape — what the backend now returns after the async migration.
 * `transactionHash` will be empty because `certificate_token_id` is null until
 * the background worker mints the certificate on-chain.
 */
const makeQueuedResponse = () => ({
    certificateId: "cert-123",
    transactionHash: "", // certificate_token_id is null when status="queued"
    claimedAt: "2024-01-01T00:00:00Z",
    message:
        "Certificate claim has been queued for processing. The certificate will be minted shortly.",
    status: "queued",
    estimatedDeadline: "2024-01-02T00:00:00Z",
});

describe("useClaimCertificate", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("initialises with correct default values", () => {
        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        expect(result.current.isClaiming).toBe(false);
        expect(result.current.claimError).toBeNull();
        expect(result.current.claimCertificate).toBeDefined();
    });

    // ─── PIN / password flow ──────────────────────────────────────────────────

    it("claims certificate with PIN and shows success toast", async () => {
        vi.mocked(certificateService.claimCertificate).mockResolvedValue(makeSuccessResponse());

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.claimCertificate({
                certificateId: "cert-123",
                accountPassword: "secret123",
            });
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(false);
        });

        expect(certificateService.claimCertificate).toHaveBeenCalledWith({
            certificateId: "cert-123",
            accountPassword: "secret123",
        });
        expect(toast.success).toHaveBeenCalledWith("Certificate claimed successfully!");
        expect(toast.error).not.toHaveBeenCalled();
    });

    // ─── Wallet-signature flow ────────────────────────────────────────────────

    it("claims certificate with wallet signature and shows success toast", async () => {
        vi.mocked(certificateService.claimCertificate).mockResolvedValue(makeSuccessResponse());

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.claimCertificate({
                certificateId: "cert-123",
                signature: "0xsig",
                signMessage: "Please sign this",
            });
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(false);
        });

        expect(certificateService.claimCertificate).toHaveBeenCalledWith({
            certificateId: "cert-123",
            signature: "0xsig",
            signMessage: "Please sign this",
        });
        expect(toast.success).toHaveBeenCalledWith("Certificate claimed successfully!");
    });

    // ─── Async / offline-claim (queued) flow ─────────────────────────────────
    //
    // The backend now returns 202 Accepted with status="queued" instead of
    // minting immediately.  The frontend currently treats the 202 as a full
    // success because it only checks whether the service call resolved without
    // throwing.
    //
    // These tests document the **current** behaviour.  Once the frontend is
    // updated to handle the queued state (polling, pending UI, etc.) they will
    // need to be revised.

    it("shows queued toast when backend queues the claim (status=queued, no tokenId)", async () => {
        vi.mocked(certificateService.claimCertificate).mockResolvedValue(makeQueuedResponse());

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.claimCertificate({
                certificateId: "cert-123",
                accountPassword: "secret123",
            });
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(false);
        });

        expect(toast.success).toHaveBeenCalledWith(
            "Certificate claim submitted — minting in progress!",
        );
        expect(toast.error).not.toHaveBeenCalled();
    });

    it("returns empty transactionHash when claim is queued (not yet minted)", async () => {
        vi.mocked(certificateService.claimCertificate).mockResolvedValue(makeQueuedResponse());

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        let claimResult: Awaited<ReturnType<typeof result.current.claimCertificate>>;

        await act(async () => {
            claimResult = await result.current.claimCertificate({
                certificateId: "cert-123",
                accountPassword: "secret123",
            });
        });

        // transactionHash is empty because certificate_token_id is null in the queued response.
        expect(claimResult!.transactionHash).toBe("");
        expect(claimResult!.status).toBe("queued");
        expect(claimResult!.estimatedDeadline).toBe("2024-01-02T00:00:00Z");
    });

    // ─── isClaiming loading state ─────────────────────────────────────────────

    it("sets isClaiming to true while the mutation is in-flight", async () => {
        vi.mocked(certificateService.claimCertificate).mockImplementation(
            () => new Promise((resolve) => setTimeout(() => resolve(makeSuccessResponse()), 100)),
        );

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        act(() => {
            result.current.claimCertificate({
                certificateId: "cert-123",
                accountPassword: "secret123",
            });
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(true);
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(false);
        });
    });

    // ─── Error handling ───────────────────────────────────────────────────────

    it("shows error toast with structured backend message", async () => {
        const backendError = {
            error: { message: "Certificate has already been claimed" },
        };
        vi.mocked(certificateService.claimCertificate).mockRejectedValue(backendError);

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        await act(async () => {
            await result.current
                .claimCertificate({
                    certificateId: "cert-123",
                    accountPassword: "wrong-password",
                })
                .catch(() => {});
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(false);
        });

        expect(toast.error).toHaveBeenCalledWith("Certificate has already been claimed");
        expect(toast.success).not.toHaveBeenCalled();
    });

    it("shows error toast with top-level message property", async () => {
        const error = { message: "Invalid password" };
        vi.mocked(certificateService.claimCertificate).mockRejectedValue(error);

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        await act(async () => {
            await result.current
                .claimCertificate({
                    certificateId: "cert-123",
                    accountPassword: "wrong",
                })
                .catch(() => {});
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(false);
        });

        expect(toast.error).toHaveBeenCalledWith("Invalid password");
    });

    it("shows fallback error toast for non-object errors", async () => {
        vi.mocked(certificateService.claimCertificate).mockRejectedValue("unexpected string error");

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        await act(async () => {
            await result.current
                .claimCertificate({
                    certificateId: "cert-123",
                    accountPassword: "secret",
                })
                .catch(() => {});
        });

        await waitFor(() => {
            expect(result.current.isClaiming).toBe(false);
        });

        expect(toast.error).toHaveBeenCalledWith("Failed to claim certificate. Please try again.");
    });

    it("exposes claimError after a failed mutation", async () => {
        const error = new Error("Network error");
        vi.mocked(certificateService.claimCertificate).mockRejectedValue(error);

        const { result } = renderHook(() => useClaimCertificate(), {
            wrapper: createWrapper(),
        });

        await act(async () => {
            await result.current
                .claimCertificate({
                    certificateId: "cert-123",
                    accountPassword: "secret",
                })
                .catch(() => {});
        });

        await waitFor(() => {
            expect(result.current.claimError).toBeDefined();
        });
    });
});
