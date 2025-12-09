import { describe, it, expect, vi, beforeEach } from "vitest";
import { OnboardRegistrationMethod, CommonSolutionStatus } from "@decm/api";
import type { CoreApiType } from "@/lib/api/api";
import type { OnboardService } from "../OnboardService/OnboardService";
import type { QueryClient } from "@tanstack/react-query";
import type { Config } from "wagmi";
const wagmiMock = vi.hoisted(() => ({
    disconnect: vi.fn(),
    getAccount: vi.fn(),
}));

vi.mock("@wagmi/core", () => wagmiMock);

import { disconnect, getAccount } from "@wagmi/core";

const mockDisconnect = wagmiMock.disconnect;
const mockGetAccount = wagmiMock.getAccount;

// Mock i18next
vi.mock("i18next", () => ({
    t: (key: string) => key,
}));

// Mock the coreApiClient at the module level - must use factory function
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            registerWithGoogleOauth: vi.fn(),
            registerWithWallet: vi.fn(),
            createProfile: vi.fn(),
            updateProfileByCredentialId: vi.fn(),
            logout: vi.fn(),
            checkRole: vi.fn(),
            getMyProfile: vi.fn(),
            verifyPassword: vi.fn(),
        },
    },
}));

// Import the actual mocked coreApiClient after vi.mock
import { coreApiClient } from "@/lib/api/api";

import { AuthService, type CreateAccountParams, type CreateProfileParams } from "./AuthService";

describe("AuthService", () => {
    let mockCoreApi: CoreApiType;
    let mockOnboardService: OnboardService;
    let mockQueryClient: QueryClient;
    let mockWagmiConfig: Config;
    let authService: AuthService;

    beforeEach(() => {
        vi.clearAllMocks();
        mockDisconnect.mockReset();
        mockGetAccount.mockReset();
        mockGetAccount.mockImplementation(() => ({ isConnected: false }));

        mockCoreApi = coreApiClient as unknown as CoreApiType;

        mockOnboardService = {
            checkOnboardStatus: vi.fn(),
        } as unknown as OnboardService;

        mockQueryClient = {
            invalidateQueries: vi.fn().mockResolvedValue(undefined),
        } as unknown as QueryClient;

        mockWagmiConfig = {} as unknown as Config;

        // Constructor order: coreApi, queryClient, onboardService, wagmiConfig
        authService = new AuthService(
            mockCoreApi,
            mockQueryClient,
            mockOnboardService,
            mockWagmiConfig,
        );
    });

    describe("createAccount", () => {
        it("should create account with Google OAuth", async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: null,
            });

            const mockResponse = { id: "google-user-id" };
            (
                coreApiClient.v1.registerWithGoogleOauth as ReturnType<typeof vi.fn>
            ).mockResolvedValue(mockResponse);

            const params = {
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: "google-token",
                password: "secure-password",
                expiresIn: 3600,
            } satisfies CreateAccountParams;

            const result = await authService.createAccount(params);

            expect(mockOnboardService.checkOnboardStatus).toHaveBeenCalledWith(params);
            expect(coreApiClient.v1.registerWithGoogleOauth).toHaveBeenCalledWith({
                access_token: "google-token",
                password: "secure-password",
            });
            expect(result).toEqual(mockResponse);
        });

        it("should create account with Wallet", async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: null,
            });

            const mockResponse = { id: "wallet-user-id" };
            (coreApiClient.v1.registerWithWallet as ReturnType<typeof vi.fn>).mockResolvedValue(
                mockResponse,
            );

            const params: CreateAccountParams = {
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: "0xsignature",
            };

            const result = await authService.createAccount(params);

            expect(mockOnboardService.checkOnboardStatus).toHaveBeenCalledWith(params);
            expect(coreApiClient.v1.registerWithWallet).toHaveBeenCalledWith({
                signed_message: "0xsignature",
            });
            expect(result).toEqual(mockResponse);
        });

        it("should throw error if credential already exists", async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: "existing-id",
            });

            const params: CreateAccountParams = {
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: "google-token",
                password: "secure-password",
                expiresIn: 3600,
            };

            await expect(authService.createAccount(params)).rejects.toThrow(
                "onboard.api.error.authentication_credential_id_already_exists",
            );
        });

        it("should throw error for Google OAuth without access token", async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: null,
            });

            const params: CreateAccountParams = {
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: "",
                password: "password",
                expiresIn: 3600,
            };

            await expect(authService.createAccount(params)).rejects.toThrow();
        });

        it("should throw error for Wallet without signature", async () => {
            (mockOnboardService.checkOnboardStatus as ReturnType<typeof vi.fn>).mockResolvedValue({
                authentication_credential_id: null,
            });

            const params: CreateAccountParams = {
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: "",
            };

            await expect(authService.createAccount(params)).rejects.toThrow("Invalid sign message");
        });
    });

    describe("createProfile", () => {
        it("should create profile with all fields", async () => {
            const mockResponse = { id: "profile-id" };
            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue(
                mockResponse,
            );

            const profileData: CreateProfileParams = {
                email: "test@example.com",
                firstName: "John",
                lastName: "Doe",
                phoneNumber: "+1234567890",
                bio: "Test bio",
                address: "123 Test St",
                academicEmail: "john@university.edu",
                academicInstitution: "Test University",
                profilePictureUrl: "https://example.com/pic.jpg",
                isEmailPublic: true,
                isFirstNamePublic: true,
                isLastNamePublic: false,
            };

            const result = await authService.createProfile("cred-id", profileData);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith({
                authentication_credential_id: "cred-id",
                email: "test@example.com",
                first_name: "John",
                last_name: "Doe",
                phone_number: "+1234567890",
                bio: "Test bio",
                address: "123 Test St",
                academic_email: "john@university.edu",
                academic_institution: "Test University",
                profile_picture_url: "https://example.com/pic.jpg",
                is_email_public: true,
                is_first_name_public: true,
                is_last_name_public: false,
                is_academic_email_public: undefined,
                is_academic_institution_public: undefined,
                is_address_public: undefined,
                is_bio_public: undefined,
                is_phone_number_public: undefined,
                is_profile_picture_public: undefined,
            });
            expect(result).toEqual(mockResponse);
            expect(mockQueryClient.invalidateQueries).toHaveBeenCalled();
        });

        it("should create profile with minimal fields", async () => {
            const mockResponse = { id: "profile-id" };
            (mockCoreApi.v1.createProfile as ReturnType<typeof vi.fn>).mockResolvedValue(
                mockResponse,
            );

            const profileData: CreateProfileParams = {
                email: "minimal@example.com",
            };

            const result = await authService.createProfile("cred-id", profileData);

            expect(mockCoreApi.v1.createProfile).toHaveBeenCalledWith({
                authentication_credential_id: "cred-id",
                email: "minimal@example.com",
                first_name: undefined,
                last_name: undefined,
                phone_number: undefined,
                bio: undefined,
                address: undefined,
                academic_email: undefined,
                academic_institution: undefined,
                profile_picture_url: undefined,
                is_email_public: undefined,
                is_first_name_public: undefined,
                is_last_name_public: undefined,
                is_academic_email_public: undefined,
                is_academic_institution_public: undefined,
                is_address_public: undefined,
                is_bio_public: undefined,
                is_phone_number_public: undefined,
                is_profile_picture_public: undefined,
            });
            expect(result).toEqual(mockResponse);
            expect(mockQueryClient.invalidateQueries).toHaveBeenCalled();
        });
    });

    describe("updateProfile", () => {
        it("should update profile with only non-empty string fields and explicit boolean values", async () => {
            const mockResponse = { id: "updated-profile-id" };
            (
                mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>
            ).mockResolvedValue(mockResponse);

            const updateData = {
                firstName: "Jane",
                lastName: "Smith",
                bio: "Updated bio",
                isFirstNamePublic: true,
                isLastNamePublic: false,
            };

            const result = await authService.updateProfile("cred-id", updateData);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: "cred-id" },
                {
                    first_name: "Jane",
                    last_name: "Smith",
                    bio: "Updated bio",
                    is_first_name_public: true,
                    is_last_name_public: false,
                    is_academic_email_public: false,
                    is_academic_institution_public: false,
                    is_address_public: false,
                    is_bio_public: false,
                    is_email_public: false,
                    is_phone_number_public: false,
                    is_profile_picture_public: false,
                },
            );
            expect(result).toEqual(mockResponse);
            expect(mockQueryClient.invalidateQueries).toHaveBeenCalled();
        });

        it("should exclude empty string fields from update request", async () => {
            const mockResponse = { id: "updated-profile-id" };
            (
                mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>
            ).mockResolvedValue(mockResponse);

            const updateData = {
                firstName: "Jane",
                lastName: "",
                email: "",
                bio: "Test bio",
            };

            const result = await authService.updateProfile("cred-id", updateData);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: "cred-id" },
                {
                    first_name: "Jane",
                    bio: "Test bio",
                    is_academic_email_public: false,
                    is_academic_institution_public: false,
                    is_address_public: false,
                    is_bio_public: false,
                    is_email_public: false,
                    is_first_name_public: false,
                    is_last_name_public: false,
                    is_phone_number_public: false,
                    is_profile_picture_public: false,
                },
            );
            expect(result).toEqual(mockResponse);
        });

        it("should handle all string fields being provided", async () => {
            const mockResponse = { id: "updated-profile-id" };
            (
                mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>
            ).mockResolvedValue(mockResponse);

            const updateData = {
                email: "new@example.com",
                firstName: "Jane",
                lastName: "Smith",
                phoneNumber: "+9876543210",
                bio: "New bio",
                address: "456 New St",
                academicEmail: "jane@university.edu",
                academicInstitution: "New University",
                profilePictureUrl: "https://example.com/new.jpg",
                isEmailPublic: true,
                isProfilePicturePublic: true,
            };

            const result = await authService.updateProfile("cred-id", updateData);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: "cred-id" },
                {
                    email: "new@example.com",
                    first_name: "Jane",
                    last_name: "Smith",
                    phone_number: "+9876543210",
                    bio: "New bio",
                    address: "456 New St",
                    academic_email: "jane@university.edu",
                    academic_institution: "New University",
                    profile_picture_url: "https://example.com/new.jpg",
                    is_email_public: true,
                    is_profile_picture_public: true,
                    is_academic_email_public: false,
                    is_academic_institution_public: false,
                    is_address_public: false,
                    is_bio_public: false,
                    is_first_name_public: false,
                    is_last_name_public: false,
                    is_phone_number_public: false,
                },
            );
            expect(result).toEqual(mockResponse);
        });
    });

    describe("signOut", () => {
        it("should sign out successfully", async () => {
            const mockResponse = { success: true };
            (coreApiClient.v1.logout as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await authService.signOut();

            expect(coreApiClient.v1.logout).toHaveBeenCalled();
            expect(mockQueryClient.invalidateQueries).toHaveBeenCalled();
            expect(getAccount).toHaveBeenCalledWith(mockWagmiConfig);
            expect(disconnect).not.toHaveBeenCalled();
        });

        it("should disconnect wallet when connected", async () => {
            mockGetAccount.mockReturnValueOnce({ isConnected: true });
            const mockResponse = { success: true };
            (coreApiClient.v1.logout as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

            await authService.signOut();

            expect(coreApiClient.v1.logout).toHaveBeenCalled();
            expect(mockQueryClient.invalidateQueries).toHaveBeenCalled();
            expect(getAccount).toHaveBeenCalledWith(mockWagmiConfig);
            expect(disconnect).toHaveBeenCalledWith(mockWagmiConfig);
        });
    });

    describe("checkRoles", () => {
        it("should check roles with all params", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: true,
                is_issuer: false,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                isAuthenticated: true,
                requireHost: true,
                requireIssuer: false,
            });

            expect(coreApiClient.v1.checkRole).toHaveBeenCalledWith({
                is_authenticated: true,
                is_host: true,
                is_issuer: undefined,
            });
            expect(result).toEqual(mockResponse);
        });

        it("should check roles with only authentication", async () => {
            const mockResponse = {
                is_authenticated: true,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                isAuthenticated: true,
            });

            expect(coreApiClient.v1.checkRole).toHaveBeenCalledWith({
                is_authenticated: true,
                is_host: undefined,
                is_issuer: undefined,
            });
            expect(result).toEqual(mockResponse);
        });

        it("should check roles with only host requirement", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: true,
                is_issuer: false,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                requireHost: true,
            });

            expect(coreApiClient.v1.checkRole).toHaveBeenCalledWith({
                is_authenticated: undefined,
                is_host: true,
                is_issuer: undefined,
            });
            expect(result).toEqual(mockResponse);
        });

        it("should check roles with only issuer requirement", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: false,
                is_issuer: true,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                requireIssuer: true,
            });

            expect(coreApiClient.v1.checkRole).toHaveBeenCalledWith({
                is_authenticated: undefined,
                is_host: undefined,
                is_issuer: true,
            });
            expect(result).toEqual(mockResponse);
        });

        it("should check roles with both host and issuer requirements", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: true,
                is_issuer: true,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                requireHost: true,
                requireIssuer: true,
            });

            expect(coreApiClient.v1.checkRole).toHaveBeenCalledWith({
                is_authenticated: undefined,
                is_host: true,
                is_issuer: true,
            });
            expect(result).toEqual(mockResponse);
        });

        it("should check roles when user has host but not issuer", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: true,
                is_issuer: false,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                requireHost: true,
                requireIssuer: true,
            });

            expect(coreApiClient.v1.checkRole).toHaveBeenCalledWith({
                is_authenticated: undefined,
                is_host: true,
                is_issuer: true,
            });
            expect(result.is_host).toBe(true);
            expect(result.is_issuer).toBe(false);
        });

        it("should check roles when user has issuer but not host", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: false,
                is_issuer: true,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                requireHost: true,
                requireIssuer: true,
            });

            expect(result.is_host).toBe(false);
            expect(result.is_issuer).toBe(true);
        });

        it("should check roles when user has neither host nor issuer", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: false,
                is_issuer: false,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                requireHost: true,
                requireIssuer: true,
            });

            expect(result.is_host).toBe(false);
            expect(result.is_issuer).toBe(false);
        });

        it("should check roles when user has both roles", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: true,
                is_issuer: true,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({
                requireHost: true,
                requireIssuer: true,
            });

            expect(result.is_host).toBe(true);
            expect(result.is_issuer).toBe(true);
        });

        it("should handle API errors gracefully", async () => {
            const error = new Error("Role check failed");
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockRejectedValue(error);

            await expect(
                authService.checkRoles({
                    requireHost: true,
                }),
            ).rejects.toThrow("Role check failed");
        });

        it("should check roles with no requirements (all undefined)", async () => {
            const mockResponse = {
                is_authenticated: true,
                is_host: false,
                is_issuer: false,
            };
            (coreApiClient.v1.checkRole as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.checkRoles({});

            expect(coreApiClient.v1.checkRole).toHaveBeenCalledWith({
                is_authenticated: undefined,
                is_host: undefined,
                is_issuer: undefined,
            });
            expect(result).toEqual(mockResponse);
        });
    });

    describe("getMyProfile", () => {
        it("should fetch and map user profile", async () => {
            const mockApiResponse = {
                profile_id: "profile-123",
                authentication_credential_id: "auth-456",
                email: "test@example.com",
                first_name: "John",
                last_name: "Doe",
                phone_number: "+1234567890",
                bio: "Test bio",
                address: "123 Test St",
                academic_email: "john@university.edu",
                academic_institution: "Test University",
                profile_picture_url: "https://example.com/pic.jpg",
                is_email_public: true,
                is_first_name_public: true,
                is_last_name_public: false,
                is_phone_number_public: false,
                is_bio_public: true,
                is_address_public: false,
                is_academic_email_public: false,
                is_academic_institution_public: false,
                is_profile_picture_public: true,
                wallet_address: "0x1234567890abcdef",
                solution_status: CommonSolutionStatus.SolutionStatusBYOK,
                google_connector_ref: "google-123",
                github_connector_ref: "github-456",
            };

            (coreApiClient.v1.getMyProfile as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockApiResponse);

            const result = await authService.getMyProfile();

            expect(coreApiClient.v1.getMyProfile).toHaveBeenCalled();
            expect(result).toEqual({
                id: "profile-123",
                profileId: "profile-123",
                authenticationCredentialId: "auth-456",
                email: "test@example.com",
                firstName: "John",
                lastName: "Doe",
                phoneNumber: "+1234567890",
                bio: "Test bio",
                address: "123 Test St",
                academicEmail: "john@university.edu",
                academicInstitution: "Test University",
                profilePictureUrl: "https://example.com/pic.jpg",
                isEmailPublic: true,
                isFirstNamePublic: true,
                isLastNamePublic: false,
                isPhoneNumberPublic: false,
                isBioPublic: true,
                isAddressPublic: false,
                isAcademicEmailPublic: false,
                isAcademicInstitutionPublic: false,
                isProfilePicturePublic: true,
                walletAddress: "0x1234567890abcdef",
                solutionStatus: "BYOK",
                googleConnectorRef: "google-123",
                githubConnectorRef: "github-456",
            });
        });
    });

    describe("verifyPassword", () => {
        it("should verify password successfully", async () => {
            const mockResponse = {
                is_success: true,
                message: "Password verified",
            };
            (coreApiClient.v1.verifyPassword as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.verifyPassword("auth-123", "correct-password");

            expect(coreApiClient.v1.verifyPassword).toHaveBeenCalledWith({
                authentication_credential_id: "auth-123",
                password: "correct-password",
            });
            expect(result).toEqual({
                isSuccess: true,
                message: "Password verified",
            });
        });

        it("should handle incorrect password", async () => {
            const mockResponse = {
                is_success: false,
                message: "Invalid password",
            };
            (coreApiClient.v1.verifyPassword as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockResolvedValue(mockResponse);

            const result = await authService.verifyPassword("auth-123", "wrong-password");

            expect(result).toEqual({
                isSuccess: false,
                message: "Invalid password",
            });
        });

        it("should throw error on verification failure", async () => {
            const error = new Error("Verification failed");
            (coreApiClient.v1.verifyPassword as ReturnType<typeof vi.fn>) = vi
                .fn()
                .mockRejectedValue(error);

            await expect(authService.verifyPassword("auth-123", "password")).rejects.toThrow(
                "Verification failed",
            );
        });
    });
});
