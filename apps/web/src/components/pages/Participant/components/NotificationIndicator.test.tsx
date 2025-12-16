import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { NotificationIndicator } from "./NotificationIndicator";

// Mock the useMyProfile hook
vi.mock("@/hooks/useMyProfile", () => ({
    useMyProfile: vi.fn(),
}));

import { useMyProfile } from "@/hooks/useMyProfile";

describe("NotificationIndicator", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("should render notification badge when there are unread messages", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 5,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        render(<NotificationIndicator />);

        // Assert
        expect(screen.getByText("5")).toBeInTheDocument();
    });

    it("should display 99+ when unread count exceeds 99", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 150,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        render(<NotificationIndicator />);

        // Assert
        expect(screen.getByText("99+")).toBeInTheDocument();
    });

    it("should not render when unread count is zero", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 0,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        const { container } = render(<NotificationIndicator />);

        // Assert
        expect(container.firstChild).toBeNull();
    });

    it("should not render when profile data is not available", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: undefined,
            isLoading: false,
            isError: false,
            isSuccess: false,
        } as any);

        // Act
        const { container } = render(<NotificationIndicator />);

        // Assert
        expect(container.firstChild).toBeNull();
    });

    it("should not render when profile is loading", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: undefined,
            isLoading: true,
            isError: false,
            isSuccess: false,
        } as any);

        // Act
        const { container } = render(<NotificationIndicator />);

        // Assert
        expect(container.firstChild).toBeNull();
    });

    it("should not render when there is an error", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: undefined,
            isLoading: false,
            isError: true,
            isSuccess: false,
            error: new Error("Failed to fetch profile"),
        } as any);

        // Act
        const { container } = render(<NotificationIndicator />);

        // Assert
        expect(container.firstChild).toBeNull();
    });

    it("should render correct number for single digit count", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 1,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        render(<NotificationIndicator />);

        // Assert
        expect(screen.getByText("1")).toBeInTheDocument();
    });

    it("should render correct number for double digit count", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 42,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        render(<NotificationIndicator />);

        // Assert
        expect(screen.getByText("42")).toBeInTheDocument();
    });

    it("should render 99 exactly when count is 99", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 99,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        render(<NotificationIndicator />);

        // Assert
        expect(screen.getByText("99")).toBeInTheDocument();
        expect(screen.queryByText("99+")).not.toBeInTheDocument();
    });

    it("should render 99+ when count is 100", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 100,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        render(<NotificationIndicator />);

        // Assert
        expect(screen.getByText("99+")).toBeInTheDocument();
    });

    it("should have animate-pulse class for visual feedback", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 3,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        const { container } = render(<NotificationIndicator />);

        // Assert
        const animatedElement = container.querySelector(".animate-pulse");
        expect(animatedElement).toBeInTheDocument();
    });

    it("should have accent color styling", () => {
        // Arrange
        vi.mocked(useMyProfile).mockReturnValue({
            data: {
                id: "profile-123",
                unreadInboxMessageCount: 7,
                profileId: "profile-123",
                authenticationCredentialId: "cred-456",
                walletAddress: "0x1234",
                solutionStatus: "SYSTEM_MANAGED",
                isFirstNamePublic: false,
                isLastNamePublic: false,
                isEmailPublic: false,
                isPhoneNumberPublic: false,
                isAddressPublic: false,
                isBioPublic: false,
                isAcademicInstitutionPublic: false,
                isAcademicEmailPublic: false,
                isProfilePicturePublic: false,
            },
            isLoading: false,
            isError: false,
            isSuccess: true,
        } as any);

        // Act
        const { container } = render(<NotificationIndicator />);

        // Assert
        const badge = container.querySelector(".bg-accent");
        const text = container.querySelector(".text-accent");
        expect(badge).toBeInTheDocument();
        expect(text).toBeInTheDocument();
    });
});