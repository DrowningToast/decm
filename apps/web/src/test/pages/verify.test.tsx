import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter } from "react-router-dom";
import type { ReactNode } from "react";
import VerifyPage from "@/pages/verify/index";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getCertificateShareData: vi.fn(),
        },
    },
}));

vi.mock("react-router-dom", async () => {
    const actual = await vi.importActual("react-router-dom");
    return { ...actual, useSearchParams: vi.fn() };
});

vi.mock("@/components/layouts/navigations/PublicNavbar", () => ({
    PublicNavbar: () => <div data-testid="public-navbar" />,
}));

// ---------------------------------------------------------------------------
// Import mocked modules so we can configure them per-test
// ---------------------------------------------------------------------------

import { coreApiClient } from "@/lib/api/api";
import { useSearchParams } from "react-router-dom";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function makeWrapper() {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false } },
    });

    return function Wrapper({ children }: { children: ReactNode }) {
        return (
            <QueryClientProvider client={queryClient}>
                <MemoryRouter>{children}</MemoryRouter>
            </QueryClientProvider>
        );
    };
}

const mockVcData = {
    data: {
        header: { "@context": [], id: "vc-1", issuanceDate: "", issuer: "", type: [] },
        data: {
            certificateId: "cert-1",
            certificateTitle: "My Cert",
            certificateSubtitle: "",
            certificateTokenId: "",
            eventName: "Event",
            eventDescription: "",
            issuedAt: "",
            issuerAddresses: "",
            issuerId: "",
            receiverAddress: "",
            status: "VALID",
            userId: "",
            backendEncryptedUserData: "",
            encryptedUserData: "",
        },
        proof: {
            hash: "",
            encryptedByBackendRawData: "",
            encryptedByUserRawData: "",
            signMessage: "",
            host: { publicKey: "", signature: "" },
            issuers: [],
        },
    },
};

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("VerifyPage", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Default: provide a handle so most tests don't need to set this up
        vi.mocked(useSearchParams).mockReturnValue([
            new URLSearchParams("handle=abc123"),
            vi.fn(),
        ] as ReturnType<typeof useSearchParams>);
    });

    it("shows loading state while fetching", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockReturnValue(new Promise(() => {}));

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        // The spinner should be present during loading
        const spinner = document.querySelector(".animate-spin");
        expect(spinner).toBeInTheDocument();
    });

    it("shows error state when API returns a non-403 error", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue(
            new Error("Network error"),
        );

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByText(/certificateVerify\.notFound/i)).toBeInTheDocument();
        });
    });

    it("shows password form when API returns 403", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue({ status: 403 });

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByRole("button", { name: /unlock/i })).toBeInTheDocument();
        });

        expect(
            screen.getByPlaceholderText(/certificateVerify\.passwordPlaceholder/i),
        ).toBeInTheDocument();
    });

    it('shows "Incorrect password" on failed unlock', async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue({ status: 403 });

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(
                screen.getByPlaceholderText(/certificateVerify\.passwordPlaceholder/i),
            ).toBeInTheDocument();
        });

        // Now mock the second call (unlock attempt) to also reject
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue(
            new Error("Wrong password"),
        );

        fireEvent.change(screen.getByPlaceholderText(/certificateVerify\.passwordPlaceholder/i), {
            target: { value: "wrongpass" },
        });
        fireEvent.click(screen.getByRole("button", { name: /unlock/i }));

        await waitFor(() => {
            expect(screen.getByRole("alert")).toBeInTheDocument();
        });
    });

    it("shows certificate title on successful unlock", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue({ status: 403 });

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(
                screen.getByPlaceholderText(/certificateVerify\.passwordPlaceholder/i),
            ).toBeInTheDocument();
        });

        vi.mocked(coreApiClient.v1.getCertificateShareData).mockResolvedValue(mockVcData as never);

        fireEvent.change(screen.getByPlaceholderText(/certificateVerify\.passwordPlaceholder/i), {
            target: { value: "correctpass" },
        });
        fireEvent.click(screen.getByRole("button", { name: /unlock/i }));

        await waitFor(() => {
            expect(screen.getByText("My Cert")).toBeInTheDocument();
        });
    });

    it("shows certificate title when data loads successfully", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockResolvedValue(mockVcData as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByText("My Cert")).toBeInTheDocument();
        });
    });

    it("shows password required alert when unlock clicked with empty password", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue({ status: 403 });

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByRole("button", { name: /unlock/i })).toBeInTheDocument();
        });

        fireEvent.click(screen.getByRole("button", { name: /unlock/i }));

        await waitFor(() => {
            expect(screen.getByRole("alert")).toBeInTheDocument();
        });

        expect(coreApiClient.v1.getCertificateShareData).toHaveBeenCalledTimes(1); // only the initial query
    });
});
