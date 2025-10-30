import type {
    ProfileCreateProfileRequest,
    ProfileCreateProfileResponse,
    OnboardCheckOnboardStatusRequest,
    OnboardCheckOnboardStatusResponse,
    OnboardRegisterResponse,
    ProfileUpdateProfileRequest,
    ProfileUpdateProfileResponse,
    OnboardRegisterWithGoogleOAuthRequest,
    OnboardRegisterWithWalletRequest,
} from "@decm/api";

// Interface for the Core API V1 client
export interface ICoreApiV1Client {
    // Auth endpoints
    registerWithGoogleOauth(
        request: OnboardRegisterWithGoogleOAuthRequest,
    ): Promise<OnboardRegisterResponse>;
    registerWithWallet(request: OnboardRegisterWithWalletRequest): Promise<OnboardRegisterResponse>;
    // loginWithGoogleOauth(request: Request): Promise<OnboardRegisterResponse>;
    // loginWithWallet(request: OnboardRegisterRequest): Promise<OnboardRegisterResponse>;
    logout(): Promise<void>;

    // Onboarding endpoints
    checkOnboardStatus(
        request?: OnboardCheckOnboardStatusRequest,
    ): Promise<OnboardCheckOnboardStatusResponse>;

    // Profile endpoints
    createProfile(request: ProfileCreateProfileRequest): Promise<ProfileCreateProfileResponse>;
    updateProfileByCredentialId(
        pathParams: { credentialId: string },
        request: ProfileUpdateProfileRequest,
    ): Promise<ProfileUpdateProfileResponse>;
    deleteProfileByCredentialId(pathParams: { credentialId: string }): Promise<void>;
    getProfileByCredentialId(pathParams: {
        credentialId: string;
    }): Promise<ProfileCreateProfileResponse>;
}

// Interface for the Core API client
export interface ICoreApiClient {
    v1: ICoreApiV1Client;
}

// Factory interface for creating API clients
export interface IApiClientFactory {
    createCoreApiClient(config?: { baseURL?: string; withCredentials?: boolean }): ICoreApiClient;
}
