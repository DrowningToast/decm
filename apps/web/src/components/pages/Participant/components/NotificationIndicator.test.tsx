import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { NotificationIndicator } from "./NotificationIndicator";
import { useMyProfile } from "@/hooks/useMyProfile";
import type { ProfileWithAuth } from "@/services/AuthService/AuthService";
vi.mock("@/hooks/useMyProfile");

describe("NotificationIndicator", () => {
    const defaultProfile: ProfileWithAuth = {
        id: "test-id",
        profileId: "test-profile-id",
        authenticationCredentialId: "test-auth-id",
        firstName: "Test",
        lastName: "User",
        email: "test@example.com",
        unreadInboxMessageCount: 0,
        solutionStatus: "SYSTEM_MANAGED",
        walletAddress: "0x123",
        isAcademicEmailPublic: false,
        isAcademicInstitutionPublic: false,
        isAddressPublic: false,
        isBioPublic: false,
        isEmailPublic: false,
        isFirstNamePublic: false,
        isLastNamePublic: false,
        isPhoneNumberPublic: false,
        isProfilePicturePublic: false,
    };

    const mockUseMyProfile = (overrides?: {
        data?: Partial<ProfileWithAuth>;
        isLoading?: boolean;
        isError?: boolean;
        isSuccess?: boolean;
        error?: Error | null;
    }) => {
        const profileData = overrides?.data
            ? { ...defaultProfile, ...overrides.data }
            : defaultProfile;

        const result = {
            data: profileData,
            isLoading: overrides?.isLoading ?? false,
            isError: overrides?.isError ?? false,
            isSuccess: overrides?.isSuccess ?? true,
            error: overrides?.error ?? null,
            refetch: vi.fn(),
            status: overrides?.isLoading ? "pending" : overrides?.isError ? "error" : "success",
            fetchStatus: "idle",
        } as unknown as ReturnType<typeof useMyProfile>;

        vi.mocked(useMyProfile).mockReturnValue(result);
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should return null when loading", () => {
        mockUseMyProfile({ isLoading: true });
        const { container } = render(<NotificationIndicator />);
        expect(container).toBeEmptyDOMElement();
    });

    it("should return null when error", () => {
        mockUseMyProfile({ isError: true });
        const { container } = render(<NotificationIndicator />);
        expect(container).toBeEmptyDOMElement();
    });

    it("should return null when unread count is 0", () => {
        mockUseMyProfile({ data: { unreadInboxMessageCount: 0 } });
        const { container } = render(<NotificationIndicator />);
        expect(container).toBeEmptyDOMElement();
    });

    it("should return null when profile is undefined", () => {
        vi.mocked(useMyProfile).mockReturnValue({
            data: undefined,
            isLoading: false,
            isError: false,
            isSuccess: true,
            error: null,
        } as unknown as ReturnType<typeof useMyProfile>);
        const { container } = render(<NotificationIndicator />);
        expect(container).toBeEmptyDOMElement();
    });

    it("should display count when > 0", () => {
        mockUseMyProfile({ data: { unreadInboxMessageCount: 5 } });
        render(<NotificationIndicator />);
        expect(screen.getByText("5")).toBeInTheDocument();
    });

    it("should display 99+ when count > 99", () => {
        mockUseMyProfile({ data: { unreadInboxMessageCount: 100 } });
        render(<NotificationIndicator />);
        expect(screen.getByText("99+")).toBeInTheDocument();
    });

    it("should render 99 exactly when count is 99", () => {
        mockUseMyProfile({ data: { unreadInboxMessageCount: 99 } });
        render(<NotificationIndicator />);
        expect(screen.getByText("99")).toBeInTheDocument();
        expect(screen.queryByText("99+")).not.toBeInTheDocument();
    });

    it("should have animate-pulse class for visual feedback", () => {
        mockUseMyProfile({ data: { unreadInboxMessageCount: 3 } });
        const { container } = render(<NotificationIndicator />);
        const animatedElement = container.querySelector(".animate-pulse");
        expect(animatedElement).toBeInTheDocument();
    });

    it("should have accent color styling", () => {
        mockUseMyProfile({ data: { unreadInboxMessageCount: 7 } });
        const { container } = render(<NotificationIndicator />);
        const badge = container.querySelector(".bg-accent");
        const text = container.querySelector(".text-accent");
        expect(badge).toBeInTheDocument();
        expect(text).toBeInTheDocument();
    });
});
