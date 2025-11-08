import { describe, it, expect, vi, beforeEach } from "vitest";
import { OnboardRegistrationMethod } from "@decm/api";
import type { CoreApiType } from "@/lib/api/api";
import type { OnboardService } from "./OnboardService";
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
        it("should update profile successfully", async () => {
            const mockResponse = { id: "updated-profile-id" };
            (
                mockCoreApi.v1.updateProfileByCredentialId as ReturnType<typeof vi.fn>
            ).mockResolvedValue(mockResponse);

            const updateData = {
                firstName: "Jane",
                lastName: "Smith",
                bio: "Updated bio",
            };

            const result = await authService.updateProfile("cred-id", updateData);

            expect(mockCoreApi.v1.updateProfileByCredentialId).toHaveBeenCalledWith(
                { credentialId: "cred-id" },
                {
                    first_name: "Jane",
                    last_name: "Smith",
                    bio: "Updated bio",
                    academic_email: undefined,
                    academic_institution: undefined,
                    address: undefined,
                    phone_number: undefined,
                    profile_picture_url: undefined,
                    is_academic_email_public: undefined,
                    is_academic_institution_public: undefined,
                    is_address_public: undefined,
                    is_bio_public: undefined,
                    is_email_public: undefined,
                    is_first_name_public: undefined,
                    is_last_name_public: undefined,
                    is_phone_number_public: undefined,
                    is_profile_picture_public: undefined,
                },
            );
            expect(result).toEqual(mockResponse);
            expect(mockQueryClient.invalidateQueries).toHaveBeenCalled();
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

            await authService.signOut();

            expect(disconnect).toHaveBeenCalledWith(mockWagmiConfig);
        });
    });
});
