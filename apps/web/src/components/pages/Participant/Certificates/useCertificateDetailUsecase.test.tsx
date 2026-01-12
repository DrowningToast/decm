import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useCertificateDetailUsecase } from "./useCertificateDetailUsecase";
import type { Certificate } from "@/services/CertificateService/mapper";

// Mock dependencies
const mockSetCertificateId = vi.fn();
const mockSetIsClaimed = vi.fn();

vi.mock("@/components/BottomNav/stores/certificates", () => ({
    useCertificateDetailNavStore: vi.fn(() => ({
        setCertificateId: mockSetCertificateId,
        setIsClaimed: mockSetIsClaimed,
    })),
}));

vi.mock("@/hooks/useMyCertificatesListViewModel", () => ({
    useMyCertificatesListViewModel: vi.fn(() => ({
        claimedCertificates: [],
        unclaimedCertificates: [],
        totalClaimed: 0,
        totalUnclaimed: 0,
        isLoading: true,
        isError: false,
        error: null,
    })),
}));

import { useMyCertificatesListViewModel } from "@/hooks/useMyCertificatesListViewModel";

describe("useCertificateDetailUsecase", () => {
    let queryClient: QueryClient;

    beforeEach(() => {
        queryClient = new QueryClient({
            defaultOptions: { queries: { retry: false } },
        });
        vi.clearAllMocks();
    });

    const wrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );

    const makeCertificate = (overrides: Partial<Certificate> = {}): Certificate => ({
        id: "cert-1",
        eventId: "evt-1",
        eventName: "Test Event",
        name: "John Doe",
        certificateTitle: "Completion",
        certificateSubtitle: "With honors",
        academicInstitution: "MIT",
        eventContractAddress: "0xEC",
        eventCertificateAddress: "0xECA",
        certificateTokenId: "tok-1",
        createdAt: "2024-06-15T00:00:00Z",
        broadcastedAt: undefined,
        estimatedDeadline: undefined,
        userClaimSignatureId: undefined,
        ...overrides,
    });

    it("returns loading state initially", () => {
        const { result } = renderHook(() => useCertificateDetailUsecase("cert-1"), {
            wrapper,
        });

        expect(result.current.isLoading).toBe(true);
        expect(result.current.certificate).toBeNull();
    });

    it("returns certificate detail after data loads", async () => {
        const cert = makeCertificate();
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [cert],
            unclaimedCertificates: [],
            totalClaimed: 1,
            totalUnclaimed: 0,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-1"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificate).not.toBeNull();
        expect(result.current.certificate!.id).toBe("cert-1");
        expect(result.current.certificate!.name).toBe("John Doe");
        expect(result.current.certificate!.event).toBe("Test Event");
        expect(result.current.certificate!.status).toBe("completed");
    });

    it("finds certificate from unclaimed list", async () => {
        const cert = makeCertificate({ id: "cert-2", certificateTokenId: undefined });
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [],
            unclaimedCertificates: [cert],
            totalClaimed: 0,
            totalUnclaimed: 1,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-2"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificate).not.toBeNull();
        expect(result.current.certificate!.status).toBe("pending");
    });

    it("maps status to queued when estimatedDeadline is set and certificateTokenId is absent", async () => {
        const futureDeadline = new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(); // 24h from now
        const cert = makeCertificate({
            id: "cert-3",
            certificateTokenId: undefined,
            estimatedDeadline: futureDeadline,
        });
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [],
            unclaimedCertificates: [cert],
            totalClaimed: 0,
            totalUnclaimed: 1,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-3"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificate!.status).toBe("queued");
        expect(result.current.certificate!.estimatedDeadline).toEqual(new Date(futureDeadline));
    });

    it("maps status to expired when estimatedDeadline is in the past", async () => {
        const pastDeadline = new Date(Date.now() - 60 * 60 * 1000).toISOString(); // 1 hour ago
        const cert = makeCertificate({
            id: "cert-exp",
            certificateTokenId: undefined,
            estimatedDeadline: pastDeadline,
        });
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [],
            unclaimedCertificates: [cert],
            totalClaimed: 0,
            totalUnclaimed: 1,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-exp"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificate!.status).toBe("expired");
    });

    it("omits estimatedDeadline from view model when status is expired", async () => {
        const pastDeadline = new Date(Date.now() - 60 * 60 * 1000).toISOString();
        const cert = makeCertificate({
            id: "cert-exp2",
            certificateTokenId: undefined,
            estimatedDeadline: pastDeadline,
        });
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [],
            unclaimedCertificates: [cert],
            totalClaimed: 0,
            totalUnclaimed: 1,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-exp2"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificate!.status).toBe("expired");
        expect(result.current.certificate!.estimatedDeadline).toBeUndefined();
    });

    it("formats date correctly", async () => {
        const cert = makeCertificate({ createdAt: "2024-06-15T00:00:00Z" });
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [cert],
            unclaimedCertificates: [],
            totalClaimed: 1,
            totalUnclaimed: 0,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-1"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // The formattedDate uses en-US locale with year, month short, day
        expect(result.current.formattedDate).toBeTruthy();
        expect(result.current.formattedDate).toContain("2024");
        expect(result.current.formattedDate).toContain("Jun");
    });

    it("handles certificate not found", async () => {
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [],
            unclaimedCertificates: [],
            totalClaimed: 0,
            totalUnclaimed: 0,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("nonexistent"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificate).toBeNull();
        expect(result.current.formattedDate).toBe("");
    });

    it("updates nav store on mount", async () => {
        const cert = makeCertificate();
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [cert],
            unclaimedCertificates: [],
            totalClaimed: 1,
            totalUnclaimed: 0,
            isLoading: false,
            isError: false,
            error: null,
        });

        renderHook(() => useCertificateDetailUsecase("cert-1"), { wrapper });

        await waitFor(() => {
            expect(mockSetCertificateId).toHaveBeenCalledWith("cert-1");
            expect(mockSetIsClaimed).toHaveBeenCalledWith(true);
        });
    });

    it("resets nav store on unmount", async () => {
        const cert = makeCertificate();
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [cert],
            unclaimedCertificates: [],
            totalClaimed: 1,
            totalUnclaimed: 0,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { unmount } = renderHook(() => useCertificateDetailUsecase("cert-1"), {
            wrapper,
        });

        await waitFor(() => {
            expect(mockSetCertificateId).toHaveBeenCalledWith("cert-1");
        });

        unmount();

        expect(mockSetCertificateId).toHaveBeenCalledWith(null);
        expect(mockSetIsClaimed).toHaveBeenCalledWith(false);
    });

    it("handles error state", async () => {
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [],
            unclaimedCertificates: [],
            totalClaimed: 0,
            totalUnclaimed: 0,
            isLoading: false,
            isError: true,
            error: new Error("Failed"),
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-1"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.isError).toBe(true);
    });

    it("defaults name and event when missing", async () => {
        const cert = makeCertificate({
            name: undefined,
            eventName: undefined,
        });
        vi.mocked(useMyCertificatesListViewModel).mockReturnValue({
            claimedCertificates: [cert],
            unclaimedCertificates: [],
            totalClaimed: 1,
            totalUnclaimed: 0,
            isLoading: false,
            isError: false,
            error: null,
        });

        const { result } = renderHook(() => useCertificateDetailUsecase("cert-1"), {
            wrapper,
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificate!.name).toBe("Untitled Certificate");
        expect(result.current.certificate!.event).toBe("Unknown Event");
    });
});
