import { vi } from "vitest";
import type { ICoreApiV1Client, ICoreApiClient } from "./interfaces";
import type {
    CheckOnboardStatusResponse,
    LoginWithGoogleOauthResponse,
    LoginWithWalletResponse,
    LogoutResponse,
    ProfileCreateProfileResponse,
    ProfileUpdateProfileByCredentialIdResponse,
    RegisterWithGoogleOauthResponse,
    RegisterWithWalletResponse,
    OnboardRegistrationMethod,
} from "@decm/api";

// Default mock responses
export const mockResponses = {
    checkOnboardStatus: {
        authentication_credential_id: "mock-credential-id",
        profile_id: null,
        message: "User needs to complete onboarding",
    } as CheckOnboardStatusResponse,

    checkOnboardStatusWithProfile: {
        authentication_credential_id: "mock-credential-id",
        profile_id: "mock-profile-id",
        message: "User has completed onboarding",
    } as CheckOnboardStatusResponse,

    registerWithGoogleOauth: {
        message: "Successfully registered with Google",
        authentication_credential_id: "google-credential-id",
    } as RegisterWithGoogleOauthResponse,

    registerWithWallet: {
        message: "Successfully registered with wallet",
        authentication_credential_id: "wallet-credential-id",
    } as RegisterWithWalletResponse,

    loginWithGoogleOauth: {
        message: "Successfully logged in with Google",
    } as LoginWithGoogleOauthResponse,

    loginWithWallet: {
        message: "Successfully logged in with wallet",
    } as LoginWithWalletResponse,

    logout: {
        message: "Successfully logged out",
    } as LogoutResponse,

    createProfile: {
        id: "profile-id",
        authentication_credential_id: "mock-credential-id",
        email: "test@example.com",
        first_name: "Test",
        last_name: "User",
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
    } as ProfileCreateProfileResponse,

    updateProfile: {
        id: "profile-id",
        authentication_credential_id: "mock-credential-id",
        email: "updated@example.com",
        first_name: "Updated",
        last_name: "User",
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
    } as ProfileUpdateProfileByCredentialIdResponse,
};

// Create mock V1 client
export const createMockCoreApiV1Client = (
    overrides: Partial<ICoreApiV1Client> = {},
): ICoreApiV1Client => {
    return {
        registerWithGoogleOauth: vi.fn().mockResolvedValue(mockResponses.registerWithGoogleOauth),
        registerWithWallet: vi.fn().mockResolvedValue(mockResponses.registerWithWallet),
        loginWithGoogleOauth: vi.fn().mockResolvedValue(mockResponses.loginWithGoogleOauth),
        loginWithWallet: vi.fn().mockResolvedValue(mockResponses.loginWithWallet),
        logout: vi.fn().mockResolvedValue(mockResponses.logout),
        checkOnboardStatus: vi.fn().mockResolvedValue(mockResponses.checkOnboardStatus),
        createProfile: vi.fn().mockResolvedValue(mockResponses.createProfile),
        updateProfileByCredentialId: vi.fn().mockResolvedValue(mockResponses.updateProfile),
        ...overrides,
    };
};

// Create mock Core API client
export const createMockCoreApiClient = (
    v1Overrides: Partial<ICoreApiV1Client> = {},
): ICoreApiClient => {
    return {
        v1: createMockCoreApiV1Client(v1Overrides),
    };
};

// Factory for creating different mock scenarios
export const mockScenarios = {
    // New user scenario (no profile)
    newUser: () =>
        createMockCoreApiClient({
            checkOnboardStatus: vi.fn().mockResolvedValue(mockResponses.checkOnboardStatus),
        }),

    // Existing user with profile
    existingUser: () =>
        createMockCoreApiClient({
            checkOnboardStatus: vi
                .fn()
                .mockResolvedValue(mockResponses.checkOnboardStatusWithProfile),
        }),

    // Failed registration scenario
    failedRegistration: () =>
        createMockCoreApiClient({
            registerWithGoogleOauth: vi.fn().mockRejectedValue(new Error("Registration failed")),
            registerWithWallet: vi.fn().mockRejectedValue(new Error("Registration failed")),
        }),

    // Network error scenario
    networkError: () => {
        const networkError = new Error("Network error");
        return createMockCoreApiClient({
            registerWithGoogleOauth: vi.fn().mockRejectedValue(networkError),
            registerWithWallet: vi.fn().mockRejectedValue(networkError),
            loginWithGoogleOauth: vi.fn().mockRejectedValue(networkError),
            loginWithWallet: vi.fn().mockRejectedValue(networkError),
            logout: vi.fn().mockRejectedValue(networkError),
            checkOnboardStatus: vi.fn().mockRejectedValue(networkError),
            createProfile: vi.fn().mockRejectedValue(networkError),
            updateProfileByCredentialId: vi.fn().mockRejectedValue(networkError),
        });
    },
};
