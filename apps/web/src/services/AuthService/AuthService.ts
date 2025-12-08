import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import { type OnboardService, defaultOnboardService } from "../OnboardService/OnboardService";
import { OnboardRegistrationMethod, type ProfileUpdateProfileRequest } from "@decm/api";
import { t } from "i18next";
import { Err } from "@/common/Err";
import type { QueryClient } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";
import { queryClient } from "@/lib/api/queryClient";
import { toast } from "sonner";
import { TOAST_USECASE_VIEWMODEL } from "@/constants/toast";
import { USECASE_IDS } from "@/constants/usecase";
import { LOCAL_STORAGE_KEYS, removeLocalStorageItem } from "@/lib/constants/localStorage";
import { wagmiConfig } from "@/config/walletConnect";
import { disconnect, getAccount, type Config } from "@wagmi/core";
import { mapProfileWithAuthViewModel } from "./mapper";

export type CreateAccountParams =
    | {
          method: OnboardRegistrationMethod.RegistrationMethodGoogle;
          accessToken: string;
          password: string;
          expiresIn: number;
      }
    | {
          method: OnboardRegistrationMethod.RegistrationMethodWallet;
          signSignature: string;
      };

export interface CreateProfileParams {
    academicEmail?: string;
    academicInstitution?: string;
    address?: string;
    bio?: string;
    email?: string;
    firstName?: string;
    lastName?: string;
    phoneNumber?: string;
    profilePictureUrl?: string;
    isAcademicEmailPublic?: boolean;
    isAcademicInstitutionPublic?: boolean;
    isAddressPublic?: boolean;
    isBioPublic?: boolean;
    isEmailPublic?: boolean;
    isFirstNamePublic?: boolean;
    isLastNamePublic?: boolean;
    isPhoneNumberPublic?: boolean;
    isProfilePicturePublic?: boolean;
}

export type UpdateProfileParams = CreateProfileParams;

export interface SignOutParams {
    showSuccessToast?: boolean;
}

interface CheckRolesParams {
    isAuthenticated?: boolean;
    requireHost?: boolean;
    requireIssuer?: boolean;
}

export type SolutionStatus = "BYOK" | "SYSTEM_MANAGED";

export interface Profile {
    academicEmail?: string;
    academicInstitution?: string;
    address?: string;
    authenticationCredentialId: string;
    bio?: string;
    email?: string;
    firstName?: string;
    id: string;
    isAcademicEmailPublic: boolean;
    isAcademicInstitutionPublic: boolean;
    isAddressPublic: boolean;
    isBioPublic: boolean;
    isEmailPublic: boolean;
    isFirstNamePublic: boolean;
    isLastNamePublic: boolean;
    isPhoneNumberPublic: boolean;
    isProfilePicturePublic: boolean;
    lastName?: string;
    phoneNumber?: string;
    profilePictureUrl?: string;
    googleConnectorRef?: string;
    walletAddress?: string;
}

export interface ProfileWithAuth extends Profile {
    walletAddress: string;
    solutionStatus: SolutionStatus;
    googleConnectorRef?: string;
    githubConnectorRef?: string;

    profileId: string;
    authenticationCredentialId: string;
}

export class AuthService {
    private _coreApi: CoreApiType;
    private _onboardService: OnboardService;
    private _queryClient: QueryClient;
    private _wagmiConfig: Config;

    constructor(
        coreApi: CoreApiType,
        queryClient: QueryClient,
        onboardService: OnboardService,
        wagmiConfig: Config,
    ) {
        this._coreApi = coreApi;
        this._queryClient = queryClient;
        this._onboardService = onboardService;
        this._wagmiConfig = wagmiConfig;
    }

    public async createAccount(params: CreateAccountParams) {
        const status = await this._onboardService.checkOnboardStatus(params);
        if (status?.authentication_credential_id) {
            throw new Error(t("onboard.api.error.authentication_credential_id_already_exists"));
        }

        switch (params.method) {
            case OnboardRegistrationMethod.RegistrationMethodGoogle: {
                if (!params.accessToken || !params.password) {
                    throw new Err("Invalid access token or password");
                }
                const response = await coreApiClient.v1.registerWithGoogleOauth({
                    access_token: params.accessToken,
                    password: params.password,
                });
                await this._queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
                return response;
            }

            case OnboardRegistrationMethod.RegistrationMethodWallet: {
                if (!params.signSignature) {
                    throw new Error("Invalid sign message");
                }
                const response = await coreApiClient.v1.registerWithWallet({
                    signed_message: params.signSignature,
                });
                await this._queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
                return response;
            }

            default:
                throw new Error("Invalid method");
        }
    }

    public async createProfile(authenticationCredentialId: string, profile: CreateProfileParams) {
        const createdProfile = await this._coreApi.v1.createProfile({
            authentication_credential_id: authenticationCredentialId,
            academic_email: profile.academicEmail,
            academic_institution: profile.academicInstitution,
            address: profile.address,
            bio: profile.bio,
            email: profile.email,
            first_name: profile.firstName,
            last_name: profile.lastName,
            phone_number: profile.phoneNumber,
            profile_picture_url: profile.profilePictureUrl,
            is_academic_email_public: profile.isAcademicEmailPublic,
            is_academic_institution_public: profile.isAcademicInstitutionPublic,
            is_address_public: profile.isAddressPublic,
            is_bio_public: profile.isBioPublic,
            is_email_public: profile.isEmailPublic,
            is_first_name_public: profile.isFirstNamePublic,
            is_last_name_public: profile.isLastNamePublic,
            is_phone_number_public: profile.isPhoneNumberPublic,
            is_profile_picture_public: profile.isProfilePicturePublic,
        });

        await this._queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
        return createdProfile;
    }

    public async updateProfile(authenticationCredentialId: string, profile: UpdateProfileParams) {
        try {
            // Build request body - only include fields that have actual values
            // For string fields: include only if value is provided and not empty
            // For boolean fields: always include with explicit value
            const requestBody: ProfileUpdateProfileRequest = {};

            // String fields - only include if value is provided and not empty
            if (profile.academicEmail !== undefined && profile.academicEmail !== "") {
                requestBody.academic_email = profile.academicEmail;
            }
            if (profile.academicInstitution !== undefined && profile.academicInstitution !== "") {
                requestBody.academic_institution = profile.academicInstitution;
            }
            if (profile.address !== undefined && profile.address !== "") {
                requestBody.address = profile.address;
            }
            if (profile.bio !== undefined && profile.bio !== "") {
                requestBody.bio = profile.bio;
            }
            if (profile.email !== undefined && profile.email !== "") {
                requestBody.email = profile.email;
            }
            if (profile.firstName !== undefined && profile.firstName !== "") {
                requestBody.first_name = profile.firstName;
            }
            if (profile.lastName !== undefined && profile.lastName !== "") {
                requestBody.last_name = profile.lastName;
            }
            if (profile.phoneNumber !== undefined && profile.phoneNumber !== "") {
                requestBody.phone_number = profile.phoneNumber;
            }
            if (profile.profilePictureUrl !== undefined && profile.profilePictureUrl !== "") {
                requestBody.profile_picture_url = profile.profilePictureUrl;
            }

            // Boolean fields - always include with explicit values
            requestBody.is_academic_email_public = profile.isAcademicEmailPublic ?? false;
            requestBody.is_academic_institution_public =
                profile.isAcademicInstitutionPublic ?? false;
            requestBody.is_address_public = profile.isAddressPublic ?? false;
            requestBody.is_bio_public = profile.isBioPublic ?? false;
            requestBody.is_email_public = profile.isEmailPublic ?? false;
            requestBody.is_first_name_public = profile.isFirstNamePublic ?? false;
            requestBody.is_last_name_public = profile.isLastNamePublic ?? false;
            requestBody.is_phone_number_public = profile.isPhoneNumberPublic ?? false;
            requestBody.is_profile_picture_public = profile.isProfilePicturePublic ?? false;

            const response = await this._coreApi.v1.updateProfileByCredentialId(
                {
                    credentialId: authenticationCredentialId,
                },
                requestBody,
            );
            await this._queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
            return response;
        } catch (error) {
            console.error(error);
            throw error;
        }
    }

    public async signOut({ showSuccessToast: showToast = true }: SignOutParams | undefined = {}) {
        try {
            await this._coreApi.v1.logout();
            removeLocalStorageItem(LOCAL_STORAGE_KEYS.ACCESS_TOKEN);
            removeLocalStorageItem(LOCAL_STORAGE_KEYS.EXPIRES_IN);
            removeLocalStorageItem(LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE);
            // if wallet is connected, disconnect it
            const account = getAccount(this._wagmiConfig);
            if (account?.isConnected) {
                await disconnect(this._wagmiConfig);
            }
            if (showToast) {
                toast.info(TOAST_USECASE_VIEWMODEL[USECASE_IDS.GENERIC].SIGN_OUT_SUCCESS);
            }
        } catch (error) {
            console.error(error);
            throw error;
        }
        await this._queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
    }

    public async checkRoles(params: CheckRolesParams) {
        const response = await this._coreApi.v1.checkRole({
            is_authenticated: params.isAuthenticated || undefined,
            is_host: params.requireHost || undefined,
            is_issuer: params.requireIssuer || undefined,
        });
        return response;
    }

    public async getMyProfile() {
        const response = await this._coreApi.v1.getMyProfile();
        return mapProfileWithAuthViewModel(response);
    }

    public async verifyPassword(
        authenticationCredentialId: string,
        password: string,
    ): Promise<{ isSuccess: boolean; message: string }> {
        try {
            const response = await this._coreApi.v1.verifyPassword({
                authentication_credential_id: authenticationCredentialId,
                password,
            });
            return {
                isSuccess: response.is_success,
                message: response.message,
            };
        } catch (error) {
            console.error("Password verification failed:", error);
            throw error;
        }
    }
}

// Default instance
export const defaultAuthService = new AuthService(
    coreApiClient,
    queryClient,
    defaultOnboardService,
    wagmiConfig,
);
