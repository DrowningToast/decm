import { vi } from "vitest";
import type { ICoreApiV1Client } from "./interfaces";
import type { OnboardCheckOnboardStatusResponse, OnboardRegisterResponse } from "@decm/api";

// Default mock responses
export const mockResponses = {
    checkOnboardStatus: {
        authentication_credential_id: "mock-credential-id",
        profile_id: undefined,
    } as OnboardCheckOnboardStatusResponse,

    checkOnboardStatusWithProfile: {
        authentication_credential_id: "mock-credential-id",
        profile_id: "mock-profile-id",
    } as OnboardCheckOnboardStatusResponse,

    registerWithGoogleOauth: {
        credential_id: "google-credential-id",
        jwt: "google-jwt",
        authentication_credential_id: "google-credential-id",
    } as OnboardRegisterResponse,

    registerWithWallet: {
        credential_id: "wallet-credential-id",
        jwt: "wallet-jwt",
        authentication_credential_id: "wallet-credential-id",
    } as OnboardRegisterResponse,

    loginWithGoogleOauth: {
        message: "Successfully logged in with Google",
        credential_id: "google-credential-id",
        jwt: "google-jwt",
        authentication_credential_id: "google-credential-id",
    } as OnboardRegisterResponse,

    loginWithWallet: {
        message: "Successfully logged in with wallet",
        credential_id: "wallet-credential-id",
        jwt: "wallet-jwt",
        authentication_credential_id: "wallet-credential-id",
    } as OnboardRegisterResponse,

    logout: {
        message: "Successfully logged out",
    },

    createProfile: {
        id: "profile-id",
        authentication_credential_id: "mock-credential-id",
        email: "test@example.com",
        first_name: "Test",
        last_name: "User",
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
    } as OnboardCheckOnboardStatusResponse,

    updateProfile: {
        id: "profile-id",
        authentication_credential_id: "mock-credential-id",
        email: "updated@example.com",
        first_name: "Updated",
        last_name: "User",
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
    } as OnboardCheckOnboardStatusResponse,
};

// Create mock V1 client
export const createMockCoreApiV1Client = (
    overrides: Partial<ICoreApiV1Client> = {},
): ICoreApiV1Client => {
    const mockClient: ICoreApiV1Client = {
        registerWithGoogleOauth: vi.fn().mockResolvedValue(mockResponses.registerWithGoogleOauth),
        registerWithWallet: vi.fn().mockResolvedValue(mockResponses.registerWithWallet),
        logout: vi.fn().mockResolvedValue(mockResponses.logout),
        checkOnboardStatus: vi.fn().mockResolvedValue(mockResponses.checkOnboardStatus),
        createProfile: vi.fn().mockResolvedValue(mockResponses.createProfile),
        updateProfileByCredentialId: vi.fn().mockResolvedValue(mockResponses.updateProfile),
        deleteProfileByCredentialId: vi.fn().mockResolvedValue(undefined),
        getProfileByCredentialId: vi.fn().mockResolvedValue(mockResponses.createProfile),
        ...overrides,
    };
    return mockClient;
};

// Factory for creating different mock scenarios
export const mockScenarios = {
    // New user scenario (no profile)
    newUser: () =>
        createMockCoreApiV1Client({
            checkOnboardStatus: vi.fn().mockResolvedValue(mockResponses.checkOnboardStatus),
        }),

    // Existing user with profile
    existingUser: () =>
        createMockCoreApiV1Client({
            checkOnboardStatus: vi
                .fn()
                .mockResolvedValue(mockResponses.checkOnboardStatusWithProfile),
        }),

    // Failed registration scenario
    failedRegistration: () =>
        createMockCoreApiV1Client({
            registerWithGoogleOauth: vi.fn().mockRejectedValue(new Error("Registration failed")),
            registerWithWallet: vi.fn().mockRejectedValue(new Error("Registration failed")),
        }),

    // Network error scenario
    networkError: () => {
        const networkError = new Error("Network error");
        return createMockCoreApiV1Client({
            registerWithGoogleOauth: vi.fn().mockRejectedValue(networkError),
            registerWithWallet: vi.fn().mockRejectedValue(networkError),
            logout: vi.fn().mockRejectedValue(networkError),
            checkOnboardStatus: vi.fn().mockRejectedValue(networkError),
            createProfile: vi.fn().mockRejectedValue(networkError),
            updateProfileByCredentialId: vi.fn().mockRejectedValue(networkError),
            deleteProfileByCredentialId: vi.fn().mockRejectedValue(networkError),
            getProfileByCredentialId: vi.fn().mockRejectedValue(networkError),
        });
    },
};
