import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import { onboardService, type OnboardService } from "./OnboardService";
import { OnboardRegistrationMethod } from "@decm/api";
import { t } from "i18next";
import { Err } from "@/common/Err";

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

export type UpdateProfileParams = Omit<CreateProfileParams, "email">;

export class AuthService {
    private _coreApi: CoreApiType;
    private _onboardService: OnboardService;

    constructor(coreApi: CoreApiType, onboardService: OnboardService) {
        this._coreApi = coreApi;
        this._onboardService = onboardService;
    }

    public async createAccount(params: CreateAccountParams) {
        const status = await this._onboardService.checkOnboardStatus(params);
        if (status?.authentication_credential_id) {
            throw new Error(t("onboard.api.error.authentication_credential_id_already_exists"));
        }

        switch (params.method) {
            case OnboardRegistrationMethod.RegistrationMethodGoogle:
                if (!params.accessToken || !params.password) {
                    throw new Err(t("Invalid access token or password"));
                }
                return coreApiClient.v1.registerWithGoogleOauth({
                    access_token: params.accessToken,
                    password: params.password,
                });
            case OnboardRegistrationMethod.RegistrationMethodWallet:
                if (!params.signSignature) {
                    throw new Error(t("Invalid sign message"));
                }
                return coreApiClient.v1.registerWithWallet({
                    signed_message: params.signSignature,
                });
            default:
                throw new Error(t("Invalid method"));
        }
    }

    public async createProfile(authenticationCredentialId: string, profile: CreateProfileParams) {
        return await this._coreApi.v1.createProfile({
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
    }

    public async updateProfile(authenticationCredentialId: string, profile: UpdateProfileParams) {
        return await this._coreApi.v1.updateProfileByCredentialId(
            {
                credentialId: authenticationCredentialId,
            },
            {
                academic_email: profile.academicEmail,
                academic_institution: profile.academicInstitution,
                address: profile.address,
                bio: profile.bio,
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
            },
        );
    }

    public async signOut() {
        return await this._coreApi.v1.logout();
    }
}

// Default instance
export const authService = new AuthService(coreApiClient, onboardService);
