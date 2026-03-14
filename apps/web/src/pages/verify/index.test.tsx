import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter } from "react-router-dom";
import type { ReactNode } from "react";
import VerifyPage from "./index";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getCertificateShareStatus: vi.fn(),
            getCertificateShareData: vi.fn(),
            getCertificateShareDataWithPassword: vi.fn(),
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

    it('shows "Invalid Link" when no handle in query params', () => {
        vi.mocked(useSearchParams).mockReturnValue([new URLSearchParams(""), vi.fn()] as ReturnType<
            typeof useSearchParams
        >);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        expect(screen.getByText("Invalid Link")).toBeInTheDocument();
    });

    it("shows loading state while fetching", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockReturnValue(
            new Promise(() => {}),
        );

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        expect(screen.getByText("Loading certificate…")).toBeInTheDocument();
    });

    it('shows "Not Found" when API returns error', async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockRejectedValue(
            new Error("Network error"),
        );

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByText("Not Found")).toBeInTheDocument();
        });
    });

    it("shows password form when status is PASSWORD_LOCKED", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockResolvedValue({
            data: { status: "PASSWORD_LOCKED" },
        } as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByRole("button", { name: /unlock/i })).toBeInTheDocument();
        });

        expect(screen.getByPlaceholderText(/enter password/i)).toBeInTheDocument();
    });

    it('shows "Incorrect password" on failed unlock', async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockResolvedValue({
            data: { status: "PASSWORD_LOCKED" },
        } as never);

        vi.mocked(coreApiClient.v1.getCertificateShareDataWithPassword).mockRejectedValue(
            new Error("Wrong password"),
        );

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByPlaceholderText(/enter password/i)).toBeInTheDocument();
        });

        fireEvent.change(screen.getByPlaceholderText(/enter password/i), {
            target: { value: "wrongpass" },
        });
        fireEvent.click(screen.getByRole("button", { name: /unlock/i }));

        await waitFor(() => {
            expect(screen.getByText("Incorrect password. Please try again.")).toBeInTheDocument();
        });
    });

    it('shows "Certificate Verified" on successful unlock', async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockResolvedValue({
            data: { status: "PASSWORD_LOCKED" },
        } as never);

        vi.mocked(coreApiClient.v1.getCertificateShareDataWithPassword).mockResolvedValue({
            data: { header: {}, data: {}, proof: {} },
        } as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByPlaceholderText(/enter password/i)).toBeInTheDocument();
        });

        fireEvent.change(screen.getByPlaceholderText(/enter password/i), {
            target: { value: "correctpass" },
        });
        fireEvent.click(screen.getByRole("button", { name: /unlock/i }));

        await waitFor(() => {
            expect(
                screen.getByRole("heading", { name: /certificate verified/i }),
            ).toBeInTheDocument();
        });
    });

    it('shows "Certificate Pending" when status is VALID_BUT_PENDING', async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockResolvedValue({
            data: { status: "VALID_BUT_PENDING" },
        } as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByText("Certificate Pending")).toBeInTheDocument();
        });
    });

    it("shows certificate title when status is READY", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockResolvedValue({
            data: {
                status: "READY",
                certificate: {
                    id: "c1",
                    certificate_title: "My Cert",
                    event_contract_address: "0xEC",
                    created_at: "2024-01-01T00:00:00Z",
                },
            },
        } as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByText("My Cert")).toBeInTheDocument();
        });
    });

    it("Unlock button is disabled when password field is empty", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareStatus).mockResolvedValue({
            data: { status: "PASSWORD_LOCKED" },
        } as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getByRole("button", { name: /unlock/i })).toBeInTheDocument();
        });

        // Password field is empty — clicking Unlock should show a validation error
        fireEvent.click(screen.getByRole("button", { name: /unlock/i }));

        await waitFor(() => {
            expect(screen.getByRole("alert")).toBeInTheDocument();
        });

        expect(coreApiClient.v1.getCertificateShareDataWithPassword).not.toHaveBeenCalled();
    });
});
