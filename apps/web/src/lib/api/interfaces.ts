import type {
    CheckOnboardStatusResponse,
    LoginWithGoogleOauthRequest,
    LoginWithGoogleOauthResponse,
    LoginWithWalletRequest,
    LoginWithWalletResponse,
    LogoutResponse,
    ProfileCreateProfileRequest,
    ProfileCreateProfileResponse,
    ProfileUpdateProfileByCredentialIdRequest,
    ProfileUpdateProfileByCredentialIdResponse,
    RegisterWithGoogleOauthRequest,
    RegisterWithGoogleOauthResponse,
    RegisterWithWalletRequest,
    RegisterWithWalletResponse,
    CheckOnboardStatusRequest,
} from "@decm/api";

// Interface for the Core API V1 client
export interface ICoreApiV1Client {
    // Auth endpoints
    registerWithGoogleOauth(
        request: RegisterWithGoogleOauthRequest,
    ): Promise<RegisterWithGoogleOauthResponse>;
    registerWithWallet(request: RegisterWithWalletRequest): Promise<RegisterWithWalletResponse>;
    loginWithGoogleOauth(
        request: LoginWithGoogleOauthRequest,
    ): Promise<LoginWithGoogleOauthResponse>;
    loginWithWallet(request: LoginWithWalletRequest): Promise<LoginWithWalletResponse>;
    logout(): Promise<LogoutResponse>;

    // Onboarding endpoints
    checkOnboardStatus(request?: CheckOnboardStatusRequest): Promise<CheckOnboardStatusResponse>;

    // Profile endpoints
    createProfile(request: ProfileCreateProfileRequest): Promise<ProfileCreateProfileResponse>;
    updateProfileByCredentialId(
        pathParams: { credentialId: string },
        request: ProfileUpdateProfileByCredentialIdRequest,
    ): Promise<ProfileUpdateProfileByCredentialIdResponse>;
}

// Interface for the Core API client
export interface ICoreApiClient {
    v1: ICoreApiV1Client;
}

// Factory interface for creating API clients
export interface IApiClientFactory {
    createCoreApiClient(config?: { baseURL?: string; withCredentials?: boolean }): ICoreApiClient;
}
