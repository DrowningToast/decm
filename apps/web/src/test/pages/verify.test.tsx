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

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
        i18n: { language: "en" },
    }),
}));

vi.mock("@/hooks/useOnChainCertificate", () => ({
    useOnChainCertificate: () => ({ data: null, isLoading: false, isError: false }),
}));

vi.mock("@/hooks/useCertificateShareImage", () => ({
    useCertificateShareImage: () => ({ imageUrl: undefined, isLoading: false }),
}));

vi.mock("@/lib/certificate/verifySignature", () => ({
    verifyProof: vi.fn().mockResolvedValue({ host: true, issuers: [] }),
}));

vi.mock("@/lib/certificate/verifyHash", () => ({
    verifyCertificateHash: vi.fn().mockReturnValue(true),
}));

// ---------------------------------------------------------------------------
// Import mocked modules so we can configure them per-test
// ---------------------------------------------------------------------------

import { AxiosError } from "axios";
import { coreApiClient } from "@/lib/api/api";
import { useSearchParams } from "react-router-dom";

function make403Error() {
    const err = new AxiosError("Forbidden");
    err.response = { status: 403 } as AxiosError["response"];
    return err;
}

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
    contract: {
        eventCertificateContractAddress: "0xcontract",
        certificateTokenId: "1",
    },
    data: {
        header: { "@context": [], id: "vc-1", issuanceDate: "", issuer: "", type: [] },
        data: {
            certificateId: "cert-1",
            certificateTitle: "My Cert",
            certificateSubtitle: "",
            certificateTokenId: "1",
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
            signMessage: "{}",
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
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue(make403Error());

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
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue(make403Error());

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
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue(make403Error());

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
            expect(screen.getAllByText("My Cert").length).toBeGreaterThan(0);
        });
    });

    it("shows certificate title when data loads successfully", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockResolvedValue(mockVcData as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getAllByText("My Cert").length).toBeGreaterThan(0);
        });
    });

    it("shows attendee data section when decryptedUserData is present", async () => {
        const dataWithAttendee = {
            ...mockVcData,
            decryptedUserData: {
                first_name: "Alice",
                last_name: "Smith",
                email: "alice@example.com",
                phone_number: "0812345678",
                academic_institution: "MIT",
                academic_email: null,
                address: null,
                bio: null,
            },
        };
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockResolvedValue(
            dataWithAttendee as never,
        );

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(
                screen.getByText("certificateVerify.table.attendeeDataHeading"),
            ).toBeInTheDocument();
        });
    });

    it("does not show attendee data section when decryptedUserData is absent", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockResolvedValue(mockVcData as never);

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        await waitFor(() => {
            expect(screen.getAllByText("My Cert").length).toBeGreaterThan(0);
        });

        expect(
            screen.queryByText("certificateVerify.table.attendeeDataHeading"),
        ).not.toBeInTheDocument();
    });

    it("shows verified badge on name when it matches decryptedCertificateData", async () => {
        const dataWithMatch = {
            ...mockVcData,
            decryptedUserData: {
                first_name: "Alice",
                last_name: "Smith",
                email: "alice@example.com",
                phone_number: null,
                academic_institution: null,
                academic_email: null,
                address: null,
                bio: null,
            },
            decryptedCertificateData: {
                name: "Alice Smith",
                academic_institution: "",
                certificate_title: "My Cert",
                certificate_subtitle: "",
            },
        };
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockResolvedValue(
            dataWithMatch as never,
        );

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        // Wait for attendee section heading, then open it
        await waitFor(() => {
            expect(
                screen.getByText("certificateVerify.table.attendeeDataHeading"),
            ).toBeInTheDocument();
        });
        fireEvent.click(screen.getByText("certificateVerify.table.attendeeDataHeading"));

        await waitFor(() => {
            expect(screen.getAllByText("Alice Smith").length).toBeGreaterThan(0);
        });
        expect(screen.getAllByText("certificateVerify.verified").length).toBeGreaterThan(0);
    });

    it("shows warning indicator on name when it does not match decryptedCertificateData", async () => {
        const dataWithMismatch = {
            ...mockVcData,
            decryptedUserData: {
                first_name: "Alice",
                last_name: "Smith",
                email: "alice@example.com",
                phone_number: null,
                academic_institution: null,
                academic_email: null,
                address: null,
                bio: null,
            },
            decryptedCertificateData: {
                name: "Bob Jones",
                academic_institution: "",
                certificate_title: "My Cert",
                certificate_subtitle: "",
            },
        };
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockResolvedValue(
            dataWithMismatch as never,
        );

        const Wrapper = makeWrapper();
        render(<VerifyPage />, { wrapper: Wrapper });

        // Wait for attendee section heading, then open it
        await waitFor(() => {
            expect(
                screen.getByText("certificateVerify.table.attendeeDataHeading"),
            ).toBeInTheDocument();
        });
        fireEvent.click(screen.getByText("certificateVerify.table.attendeeDataHeading"));

        await waitFor(() => {
            expect(screen.getAllByText("Alice Smith").length).toBeGreaterThan(0);
        });
        expect(screen.getByTestId("mock-alerttriangle")).toBeInTheDocument();
    });

    it("shows password required alert when unlock clicked with empty password", async () => {
        vi.mocked(coreApiClient.v1.getCertificateShareData).mockRejectedValue(make403Error());

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
