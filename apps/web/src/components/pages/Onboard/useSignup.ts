import { coreApi } from "@/lib/api/api";
import {
    OnboardRegistrationMethod,
    type ProfileUpdateProfileRequest,
    type UpdateProfileByCredentialIdParams,
} from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import type { Profile } from "./Profile";

type CreateAccountParams = {
    method: OnboardRegistrationMethod;
    accessToken?: string;
    signMessage?: string;
    password?: string;
};

const createAccountAsync = async ({
    method,
    accessToken,
    signMessage,
    password,
}: CreateAccountParams) => {
    switch (method) {
        case OnboardRegistrationMethod.RegistrationMethodGoogle:
            if (!accessToken || !password) {
                throw new Error("Invalid access token or password");
            }
            return coreApi.v1.registerWithGoogleOauth({
                access_token: accessToken,
                password: password,
            });
        case OnboardRegistrationMethod.RegistrationMethodWallet:
            if (!signMessage) {
                throw new Error("Invalid sign message");
            }
            return coreApi.v1.registerWithWallet({
                signed_message: signMessage,
            });
        default:
            throw new Error("Invalid method");
    }
};

type CreateProfileParams = Profile;

const updateAccountAsync = async (
    params: UpdateProfileByCredentialIdParams,
    body: ProfileUpdateProfileRequest,
) => {
    return await coreApi.v1.updateProfileByCredentialId(params, body);
};

const createProfileAsync = async (
    authenticationCredentialId: string,
    {
        email,
        phoneNumber,
        firstName,
        lastName,
        isEmailPublic,
        isPhoneNumberPublic,
        isFirstNamePublic,
        isLastNamePublic,
    }: CreateProfileParams,
) => {
    return await coreApi.v1.createProfile({
        authentication_credential_id: authenticationCredentialId,
        email: email,
        phone_number: phoneNumber,
        first_name: firstName,
        last_name: lastName,
        is_email_public: isEmailPublic,
        is_phone_number_public: isPhoneNumberPublic,
        is_first_name_public: isFirstNamePublic,
        is_last_name_public: isLastNamePublic,
    });
};

export const useSignup = () => {
    const { mutateAsync: signup, isPending } = useMutation({
        mutationKey: ["signup"],
        mutationFn: async ({
            method,
            accessToken,
            signMessage,
            profile,
            password,
            expiresIn,
        }: {
            method: OnboardRegistrationMethod;
            accessToken?: string;
            expiresIn?: number;
            signMessage?: string;
            profile: Profile;
            password: string;
        }) => {
            // validate
            if (!accessToken && !signMessage) {
                return;
            }

            const { authentication_credential_id, profile_id } =
                await coreApi.v1.checkOnboardStatus({
                    method,
                    access_token: accessToken,
                    sign_message: signMessage,
                    expires_in: expiresIn,
                });
            if (authentication_credential_id && profile_id) {
                return;
            }

            let accountId: string | undefined = undefined;
            if (!authentication_credential_id) {
                const accountResponse = await createAccountAsync({
                    method,
                    accessToken,
                    signMessage,
                    password,
                });
                accountId = accountResponse.authentication_credential_id;
            }

            if (!accountId) {
                throw new Error("Invalid account id");
            }
            if (!profile_id) {
                await createProfileAsync(accountId, {
                    email: profile.email,
                    phoneNumber: profile.phoneNumber,
                    firstName: profile.firstName,
                    lastName: profile.lastName,
                    isEmailPublic: profile.isEmailPublic,
                    isPhoneNumberPublic: profile.isPhoneNumberPublic,
                    isFirstNamePublic: profile.isFirstNamePublic,
                    isLastNamePublic: profile.isLastNamePublic,
                });
            } else {
                await updateAccountAsync(
                    {
                        credentialId: profile_id,
                    },
                    {
                        email: profile.email,
                        phone_number: profile.phoneNumber,
                        first_name: profile.firstName,
                        last_name: profile.lastName,
                        is_email_public: profile.isEmailPublic,
                        is_phone_number_public: profile.isPhoneNumberPublic,
                        is_first_name_public: profile.isFirstNamePublic,
                        is_last_name_public: profile.isLastNamePublic,
                    },
                );
            }

            return {
                authentication_credential_id,
                profile_id,
            };
        },
    });

    return {
        signup,
        isPending,
    };
};
