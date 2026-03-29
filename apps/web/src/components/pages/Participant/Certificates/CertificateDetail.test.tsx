import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { CertificateDetail } from "./CertificateDetail";

// ─── Mocks ───────────────────────────────────────────────────────────────────

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string | Record<string, unknown>) =>
            typeof fallback === "string" ? fallback : _key,
        i18n: { language: "en" },
    }),
}));

vi.mock("react-router-dom", () => ({
    Link: ({ children, to }: React.PropsWithChildren<{ to: string }>) => (
        <a href={to}>{children}</a>
    ),
}));

const mockSetOnClickShareable = vi.fn();
const mockSetImageUrl = vi.fn();
const mockSetOnChangePublish = vi.fn();
const mockSetIsPublished = vi.fn();
const mockSetShareableHandle = vi.fn();
const mockSetIsPasswordDialogOpen = vi.fn();

let mockIsPasswordDialogOpen = false;

vi.mock("@/components/BottomNav/stores/certificates", () => ({
    useCertificateDetailNavStore: vi.fn(() => ({
        setOnClickShareable: mockSetOnClickShareable,
        setImageUrl: mockSetImageUrl,
    })),
    useCertificateDetailsSharedNavStore: vi.fn(() => ({
        setOnChangePublish: mockSetOnChangePublish,
        setIsPublished: mockSetIsPublished,
        setIsPasswordProtected: vi.fn(),
        setShareableHandle: mockSetShareableHandle,
        isPasswordDialogOpen: mockIsPasswordDialogOpen,
        setIsPasswordDialogOpen: mockSetIsPasswordDialogOpen,
    })),
}));

vi.mock("./useCertificateDetailUsecase", () => ({
    useCertificateDetailUsecase: vi.fn(),
}));

vi.mock("./useCertificateCreateShareLinkUsecase", () => ({
    useCertificateCreateShareLinkUsecase: vi.fn(() => ({
        createCertificateShareLink: vi.fn(),
    })),
}));

vi.mock("./useCertificateUpdateShareVisibilityUsecase", () => ({
    useCertificateUpdateShareVisibilityUsecase: vi.fn(() => ({
        updateShareVisibility: vi.fn(),
    })),
}));

vi.mock("@/hooks/usePassowordPrompt", () => ({
    usePasswordPrompt: vi.fn(() => ({ mutateAsync: vi.fn() })),
}));

vi.mock("@/hooks/useClaimCertificate", () => ({
    useClaimCertificate: vi.fn(() => ({ claimCertificate: vi.fn(), isClaiming: false })),
}));

vi.mock("@/hooks/useCertificateImage", () => ({
    useCertificateImage: vi.fn(() => ({
        imageUrl: null,
        imageData: null,
        isLoading: false,
        error: null,
    })),
}));

vi.mock("@/hooks/events/useEvent", () => ({
    useEvent: vi.fn(() => ({ event: null, isLoadingEvent: false, isLoadingEventError: false })),
}));

// ─── BYOK / Auth / Services mocks ────────────────────────────────────────────

const mockGetClaimCertificateSignMessage = vi.hoisted(() => vi.fn());
const mockWalletSignMessage = vi.hoisted(() => vi.fn());

vi.mock("@/context/AuthContext", () => ({
    useAuth: vi.fn(() => ({
        user: { solutionStatus: "SYSTEM_MANAGED" },
        isFetching: false,
        isAuthenticated: true,
        refetch: vi.fn(),
    })),
}));

vi.mock("@/services/services", () => ({
    certificateService: {
        getClaimCertificateSignMessage: mockGetClaimCertificateSignMessage,
    },
    walletSigningService: {
        signMessage: mockWalletSignMessage,
    },
}));

vi.mock("@/services/WalletSigningService/WalletSigningService", () => ({
    WalletNotConnectedError: class WalletNotConnectedError extends Error {
        constructor() {
            super("Wallet is not connected");
            this.name = "WalletNotConnectedError";
        }
    },
    WalletSigningRejectedError: class WalletSigningRejectedError extends Error {
        constructor() {
            super("User rejected the signing request");
            this.name = "WalletSigningRejectedError";
        }
    },
}));

vi.mock("sonner", () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    },
}));

vi.mock("@/components/BottomNav/BottomNav", () => ({
    BottomNav: ({ variant }: { variant: string }) => (
        <div data-testid="bottom-nav" data-variant={variant} />
    ),
}));

vi.mock("@/components/common/EthscanLink", () => ({
    EthExplorerLink: ({ children }: React.PropsWithChildren) => <span>{children}</span>,
}));

vi.mock("./CertificatePasswordDialog", () => ({
    CertificatePasswordDialog: ({ open }: { open: boolean; onOpenChange: (v: boolean) => void }) =>
        open ? <div data-testid="password-dialog" /> : null,
}));

vi.mock("./CertificateShareModal", () => ({
    CertificateShareModal: () => <div data-testid="share-modal" />,
}));

// ─── Import mocked hooks so we can control return values per test ─────────────

import { useCertificateDetailUsecase } from "./useCertificateDetailUsecase";
import { useAuth } from "@/context/AuthContext";
import { useClaimCertificate } from "@/hooks/useClaimCertificate";
import { usePasswordPrompt } from "@/hooks/usePassowordPrompt";
import { toast } from "sonner";
import { useCertificateImage } from "@/hooks/useCertificateImage";
import { useEvent } from "@/hooks/events/useEvent";
import { useCertificateDetailsSharedNavStore } from "@/components/BottomNav/stores/certificates";

// ─── Helpers ─────────────────────────────────────────────────────────────────

const baseCertificate = {
    id: "cert-1",
    eventId: "evt-1",
    event: "My Event",
    name: "Jane Doe",
    status: "completed" as "completed" | "queued" | "expired" | "aborted" | "pending",
    issuedAt: "2024-06-15T00:00:00Z",
    certificateTitle: "Certificate of Completion",
    certificateSubtitle: undefined as string | undefined,
    academicInstitution: undefined as string | undefined,
    certificateContractAddress: undefined as string | undefined,
    eventContractAddress: undefined as string | undefined,
    verifiableCredentialUrl: undefined as string | undefined,
    estimatedDeadline: undefined as Date | undefined,
    shareable: undefined as
        | {
              id: string;
              handle: string;
              active: boolean;
              hasPassword: boolean;
              eventCertificateId: string;
              createdAt: string;
              updatedAt: string;
          }
        | undefined,
};

function mockUsecase(
    overrides: Partial<typeof baseCertificate> | null,
    opts?: {
        isLoading?: boolean;
        isError?: boolean;
    },
) {
    vi.mocked(useCertificateDetailUsecase).mockReturnValue({
        certificate: overrides !== null ? { ...baseCertificate, ...overrides } : null,
        formattedDate: "Jun 15, 2024",
        isLoading: opts?.isLoading ?? false,
        isError: opts?.isError ?? false,
    });
}

// ─── Tests ───────────────────────────────────────────────────────────────────

describe("CertificateDetail", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockIsPasswordDialogOpen = false;
        vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
            setOnChangePublish: mockSetOnChangePublish,
            setIsPublished: mockSetIsPublished,
            setIsPasswordProtected: vi.fn(),
            setShareableHandle: mockSetShareableHandle,
            isPasswordDialogOpen: false,
            setIsPasswordDialogOpen: mockSetIsPasswordDialogOpen,
        } as ReturnType<typeof useCertificateDetailsSharedNavStore>);
        vi.mocked(useCertificateImage).mockReturnValue({
            imageUrl: null,
            imageData: null,
            isLoading: false,
            error: null,
        });
        vi.mocked(useEvent).mockReturnValue({
            event: undefined,
            isLoadingEvent: false,
            isLoadingEventError: false,
        } as ReturnType<typeof useEvent>);
    });

    // ── Page states ──────────────────────────────────────────────────────────

    it("renders spinner while loading", () => {
        mockUsecase({}, { isLoading: true });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(document.querySelector(".animate-spin")).toBeInTheDocument();
    });

    it("renders error message on error", () => {
        mockUsecase({}, { isError: true });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("An error occurred")).toBeInTheDocument();
    });

    it("renders not-found message when certificate is null", () => {
        mockUsecase(null);
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Certificate not found")).toBeInTheDocument();
    });

    // ── Core content ─────────────────────────────────────────────────────────

    it("renders certificate title", () => {
        mockUsecase({});
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Certificate of Completion")).toBeInTheDocument();
    });

    it("renders participant name", () => {
        mockUsecase({});
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Jane Doe")).toBeInTheDocument();
    });

    it("renders event link", () => {
        mockUsecase({});
        render(<CertificateDetail certificateId="cert-1" />);
        const link = screen.getByRole("link", { name: "My Event" });
        expect(link).toHaveAttribute("href", "/app/events/evt-1");
    });

    it("renders issued date for claimed certificate", () => {
        mockUsecase({ status: "completed" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Certificate issued on")).toBeInTheDocument();
    });

    it("renders available-from label for pending certificate", () => {
        mockUsecase({ status: "pending" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Available to claim from")).toBeInTheDocument();
    });

    it("renders formatted date", () => {
        mockUsecase({});
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Jun 15, 2024")).toBeInTheDocument();
    });

    // ── Conditional content ───────────────────────────────────────────────────

    it("renders academic institution when present", () => {
        mockUsecase({ academicInstitution: "MIT" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Academic Institution")).toBeInTheDocument();
        expect(screen.getByText("MIT")).toBeInTheDocument();
    });

    it("does not render academic institution section when absent", () => {
        mockUsecase({ academicInstitution: undefined });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.queryByText("Academic Institution")).not.toBeInTheDocument();
    });

    it("renders certificate contract address when present", () => {
        mockUsecase({ certificateContractAddress: "0xABC" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("0xABC")).toBeInTheDocument();
    });

    it("renders event contract address when present", () => {
        mockUsecase({ eventContractAddress: "0xDEF" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("0xDEF")).toBeInTheDocument();
    });

    it("renders verifiable credential link when url is present", () => {
        mockUsecase({ verifiableCredentialUrl: "https://vc.example.com" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(
            screen.getByRole("link", {
                name: /participant\.certificates\.detail\.downloadAttributes/i,
            }),
        ).toHaveAttribute("href", "https://vc.example.com");
    });

    // ── Status banners ────────────────────────────────────────────────────────

    it("renders queued banner when status is queued", () => {
        mockUsecase({ status: "queued" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Certificate minting in progress...")).toBeInTheDocument();
    });

    it("renders expired banner when status is expired", () => {
        mockUsecase({ status: "expired" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Certificate claim expired")).toBeInTheDocument();
    });

    it("renders aborted banner when status is aborted", () => {
        mockUsecase({ status: "aborted" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Certificate claim could not be processed")).toBeInTheDocument();
    });

    it("renders claimed banner when certificate is claimed", () => {
        mockUsecase({ status: "completed" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Certificate Claimed")).toBeInTheDocument();
    });

    // ── Claim button ──────────────────────────────────────────────────────────

    it("shows claim button for pending certificate when event is joined", () => {
        mockUsecase({ status: "pending" });
        vi.mocked(useEvent).mockReturnValue({
            event: { isJoined: true },
            isLoadingEvent: false,
            isLoadingEventError: false,
        } as ReturnType<typeof useEvent>);

        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByRole("button", { name: /claim certificate/i })).toBeInTheDocument();
    });

    it("hides claim button when certificate is claimed", () => {
        mockUsecase({ status: "completed" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.queryByRole("button", { name: /claim/i })).not.toBeInTheDocument();
    });

    it("hides claim button when status is queued", () => {
        mockUsecase({ status: "queued" });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.queryByRole("button", { name: /claim/i })).not.toBeInTheDocument();
    });

    it("shows retry claim button for aborted certificate", () => {
        mockUsecase({ status: "aborted" });
        vi.mocked(useEvent).mockReturnValue({
            event: { isJoined: true },
            isLoadingEvent: false,
            isLoadingEventError: false,
        } as ReturnType<typeof useEvent>);
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByRole("button", { name: /retry claim/i })).toBeInTheDocument();
    });

    it("shows Go to Event button when event is not joined", () => {
        mockUsecase({ status: "pending" });
        vi.mocked(useEvent).mockReturnValue({
            event: { isJoined: false },
            isLoadingEvent: false,
            isLoadingEventError: false,
        } as ReturnType<typeof useEvent>);

        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByRole("link", { name: /go to event page/i })).toBeInTheDocument();
    });

    // ── Certificate image ─────────────────────────────────────────────────────

    it("renders image when url is loaded", () => {
        mockUsecase({});
        vi.mocked(useCertificateImage).mockReturnValue({
            imageUrl: "blob:http://localhost/img",
            imageData: null,
            isLoading: false,
            error: null,
        });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByRole("img")).toHaveAttribute("src", "blob:http://localhost/img");
    });

    it("renders image loading state", () => {
        mockUsecase({});
        vi.mocked(useCertificateImage).mockReturnValue({
            imageUrl: null,
            imageData: null,
            isLoading: true,
            error: null,
        });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText(/loading certificate preview/i)).toBeInTheDocument();
    });

    it("renders image error state", () => {
        mockUsecase({});
        vi.mocked(useCertificateImage).mockReturnValue({
            imageUrl: null,
            imageData: null,
            isLoading: false,
            error: new Error("failed"),
        });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByText("Failed to load certificate image")).toBeInTheDocument();
    });

    // ── BottomNav variant ─────────────────────────────────────────────────────

    it("uses certificate-details-shared variant when shareable exists", () => {
        mockUsecase({
            shareable: {
                id: "s1",
                handle: "abc",
                active: true,
                hasPassword: false,
                eventCertificateId: "cert-1",
                createdAt: "2024-01-01T00:00:00Z",
                updatedAt: "2024-01-01T00:00:00Z",
            },
        });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByTestId("bottom-nav")).toHaveAttribute(
            "data-variant",
            "certificate-details-shared",
        );
    });

    it("uses certificate-detail variant when no shareable", () => {
        mockUsecase({ shareable: undefined });
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByTestId("bottom-nav")).toHaveAttribute(
            "data-variant",
            "certificate-detail",
        );
    });

    // ── Share modal ───────────────────────────────────────────────────────────

    it("renders share modal in the page", () => {
        mockUsecase({});
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByTestId("share-modal")).toBeInTheDocument();
    });

    // ── Password dialog ───────────────────────────────────────────────────────

    it("does not render password dialog when closed", () => {
        mockUsecase({});
        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.queryByTestId("password-dialog")).not.toBeInTheDocument();
    });

    it("renders password dialog when open", () => {
        mockUsecase({});
        vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
            setOnChangePublish: mockSetOnChangePublish,
            setIsPublished: mockSetIsPublished,
            setIsPasswordProtected: vi.fn(),
            setShareableHandle: mockSetShareableHandle,
            isPasswordDialogOpen: true,
            setIsPasswordDialogOpen: mockSetIsPasswordDialogOpen,
        } as ReturnType<typeof useCertificateDetailsSharedNavStore>);

        render(<CertificateDetail certificateId="cert-1" />);
        expect(screen.getByTestId("password-dialog")).toBeInTheDocument();
    });
});

// ─── BYOK Tests ──────────────────────────────────────────────────────────────

describe("CertificateDetail – BYOK user claiming", () => {
    const mockClaimCertificate = vi.fn().mockResolvedValue({ status: "queued" });
    const mockOpenPasswordPrompt = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();

        // BYOK user
        vi.mocked(useAuth).mockReturnValue({
            user: { solutionStatus: "BYOK" } as ReturnType<typeof useAuth>["user"],
            isFetching: false,
            isAuthenticated: true,
            refetch: vi.fn(),
        });

        vi.mocked(usePasswordPrompt).mockReturnValue({
            mutateAsync: mockOpenPasswordPrompt,
        } as ReturnType<typeof usePasswordPrompt>);

        vi.mocked(useClaimCertificate).mockReturnValue({
            claimCertificate: mockClaimCertificate,
            isClaiming: false,
            claimError: null,
        });

        vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
            setOnChangePublish: mockSetOnChangePublish,
            setIsPublished: mockSetIsPublished,
            setIsPasswordProtected: vi.fn(),
            setShareableHandle: mockSetShareableHandle,
            isPasswordDialogOpen: false,
            setIsPasswordDialogOpen: mockSetIsPasswordDialogOpen,
        } as ReturnType<typeof useCertificateDetailsSharedNavStore>);

        vi.mocked(useCertificateImage).mockReturnValue({
            imageUrl: null,
            imageData: null,
            isLoading: false,
            error: null,
        });

        vi.mocked(useEvent).mockReturnValue({
            event: { isJoined: true },
            isLoadingEvent: false,
            isLoadingEventError: false,
        } as ReturnType<typeof useEvent>);

        mockGetClaimCertificateSignMessage.mockResolvedValue({
            signMessage: "signer,contract,99999",
        });
        mockWalletSignMessage.mockResolvedValue("0xabc123signature");
    });

    it("does not open password prompt when BYOK user claims", async () => {
        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(mockOpenPasswordPrompt).not.toHaveBeenCalled();
        });
    });

    it("calls walletSigningService.signMessage with the fetched sign message", async () => {
        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(mockGetClaimCertificateSignMessage).toHaveBeenCalledWith("cert-1");
            expect(mockWalletSignMessage).toHaveBeenCalledWith("signer,contract,99999");
        });
    });

    it("calls claimCertificate with signature params (not accountPassword) for BYOK", async () => {
        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(mockClaimCertificate).toHaveBeenCalledWith({
                certificateId: "cert-1",
                signature: "0xabc123signature",
                signMessage: "signer,contract,99999",
            });
        });

        expect(mockClaimCertificate).not.toHaveBeenCalledWith(
            expect.objectContaining({ accountPassword: expect.anything() }),
        );
    });

    it("shows wallet-not-connected error toast when wallet is not connected", async () => {
        const { WalletNotConnectedError } = await import(
            "@/services/WalletSigningService/WalletSigningService"
        );
        mockWalletSignMessage.mockRejectedValue(new WalletNotConnectedError());

        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(vi.mocked(toast.error)).toHaveBeenCalledWith(
                expect.stringMatching(/wallet.*not.*connected/i),
            );
        });
        expect(mockClaimCertificate).not.toHaveBeenCalled();
    });

    it("shows signing-rejected error toast when user rejects wallet signing", async () => {
        const { WalletSigningRejectedError } = await import(
            "@/services/WalletSigningService/WalletSigningService"
        );
        mockWalletSignMessage.mockRejectedValue(new WalletSigningRejectedError());

        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(vi.mocked(toast.error)).toHaveBeenCalledWith(expect.stringMatching(/rejected/i));
        });
        expect(mockClaimCertificate).not.toHaveBeenCalled();
    });

    it("does not call claimCertificate when fetching sign message fails", async () => {
        mockGetClaimCertificateSignMessage.mockRejectedValue(new Error("network error"));

        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(mockWalletSignMessage).not.toHaveBeenCalled();
            expect(mockClaimCertificate).not.toHaveBeenCalled();
        });
    });
});

// ─── SYSTEM_MANAGED regression Tests ─────────────────────────────────────────

describe("CertificateDetail – SYSTEM_MANAGED user claiming", () => {
    const mockClaimCertificate = vi.fn().mockResolvedValue({ status: "queued" });
    const mockOpenPasswordPrompt = vi.fn().mockResolvedValue("correct-password");

    beforeEach(() => {
        vi.clearAllMocks();

        vi.mocked(useAuth).mockReturnValue({
            user: { solutionStatus: "SYSTEM_MANAGED" } as ReturnType<typeof useAuth>["user"],
            isFetching: false,
            isAuthenticated: true,
            refetch: vi.fn(),
        });

        vi.mocked(usePasswordPrompt).mockReturnValue({
            mutateAsync: mockOpenPasswordPrompt,
        } as ReturnType<typeof usePasswordPrompt>);

        vi.mocked(useClaimCertificate).mockReturnValue({
            claimCertificate: mockClaimCertificate,
            isClaiming: false,
            claimError: null,
        });

        vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
            setOnChangePublish: mockSetOnChangePublish,
            setIsPublished: mockSetIsPublished,
            setIsPasswordProtected: vi.fn(),
            setShareableHandle: mockSetShareableHandle,
            isPasswordDialogOpen: false,
            setIsPasswordDialogOpen: mockSetIsPasswordDialogOpen,
        } as ReturnType<typeof useCertificateDetailsSharedNavStore>);

        vi.mocked(useCertificateImage).mockReturnValue({
            imageUrl: null,
            imageData: null,
            isLoading: false,
            error: null,
        });

        vi.mocked(useEvent).mockReturnValue({
            event: { isJoined: true },
            isLoadingEvent: false,
            isLoadingEventError: false,
        } as ReturnType<typeof useEvent>);
    });

    it("opens password prompt (not wallet signing) for SYSTEM_MANAGED user", async () => {
        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(mockOpenPasswordPrompt).toHaveBeenCalled();
            expect(mockGetClaimCertificateSignMessage).not.toHaveBeenCalled();
            expect(mockWalletSignMessage).not.toHaveBeenCalled();
        });
    });

    it("calls claimCertificate with accountPassword for SYSTEM_MANAGED user", async () => {
        vi.mocked(useCertificateDetailUsecase).mockReturnValue({
            certificate: { ...baseCertificate, status: "pending" },
            formattedDate: "Jun 15, 2024",
            isLoading: false,
            isError: false,
        });

        render(<CertificateDetail certificateId="cert-1" />);
        fireEvent.click(screen.getByRole("button", { name: /claim certificate/i }));

        await waitFor(() => {
            expect(mockClaimCertificate).toHaveBeenCalledWith({
                certificateId: "cert-1",
                accountPassword: "correct-password",
            });
        });

        expect(mockClaimCertificate).not.toHaveBeenCalledWith(
            expect.objectContaining({ signature: expect.anything() }),
        );
    });
});
