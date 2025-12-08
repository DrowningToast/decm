import { useNavigate } from "@/router";
import { useTranslation } from "react-i18next";
import { useSignup } from "../useSignup";
import { OnboardRegistrationMethod } from "@decm/api";
import { coreApiClient } from "@/lib/api/api";
import { handleUniversalError } from "@/common/Err";
import { toast } from "sonner";
import { USECASE_IDS } from "@/constants/usecase";
import type { OnboardCheckOnboardStatusResponse } from "@decm/api";
import type { Profile } from "../ProfilePage";
import { useMyProfile } from "@/hooks/useMyProfile";

export const useOAuthOnboardUsecase = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { createAccount, upsertProfile, isLoading } = useSignup();
    const { refetch: refetchMyProfile } = useMyProfile();

    const usecaseAsync = async (
        accessToken: string,
        expiresIn: number,
        password: string,
        profile: Profile,
    ) => {
        let status: OnboardCheckOnboardStatusResponse | undefined = undefined;
        try {
            status = await coreApiClient.v1.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                access_token: accessToken,
                expires_in: expiresIn,
            });
        } catch (error) {
            if (error instanceof Error) {
                handleUniversalError(t, error);
            }
            return;
        }

        let credential_id: string | undefined = undefined;
        if (!status?.authentication_credential_id) {
            try {
                const { credential_id: _credential_id } = await createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: accessToken ?? "",
                    expiresIn: expiresIn ?? 0,
                    password: password,
                });
                credential_id = _credential_id;
            } catch (error) {
                if (error instanceof Error) {
                    handleUniversalError(
                        t,
                        error,
                        {
                            onDuplicateEntry: () => {
                                toast.error(
                                    t("onboard.flow.oauth_google.create_account_error_duplicate"),
                                );
                                return;
                            },
                            onUnauthorized: () => {
                                toast.error(
                                    t(
                                        "onboard.flow.oauth_google.create_account_error_expired_token",
                                    ),
                                );
                                return;
                            },
                            onInternalServerError: () => {
                                toast.error(
                                    t("onboard.flow.oauth_google.create_account_error_generic"),
                                );
                                return;
                            },
                        },
                        USECASE_IDS.OAUTH_GOOGLE_CREATE_ACCOUNT,
                    );
                    return;
                }
            }
        }

        try {
            await upsertProfile({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: accessToken ?? "",
                expiresIn: expiresIn ?? 0,
                password: password,
                profile: {
                    authentication_credential_id: credential_id ?? "",
                    email: profile.email,
                    phone_number: profile.phoneNumber,
                    first_name: profile.firstName,
                    last_name: profile.lastName,
                    is_email_public: profile.isEmailPublic,
                    is_phone_number_public: profile.isPhoneNumberPublic,
                    is_first_name_public: profile.isFirstNamePublic,
                    is_last_name_public: profile.isLastNamePublic,
                },
            });
        } catch (error) {
            if (error instanceof Error) {
                handleUniversalError(
                    t,
                    error,
                    {
                        onInternalServerError: () => {
                            toast.error(
                                t("onboard.flow.oauth_google.create_account_error_generic"),
                            );
                            return;
                        },
                        onDuplicateEntry: () => {
                            toast.error(
                                t("onboard.flow.oauth_google.create_account_error_duplicate"),
                            );
                            return;
                        },
                    },
                    USECASE_IDS.OAUTH_GOOGLE_CREATE_PROFILE,
                );
            }
            return;
        }

        toast.success(t("onboard.flow.oauth_google.create_profile_success"));
        await refetchMyProfile();
        await navigate("/app");
    };

    return {
        usecaseAsync,
        isLoading,
    };
};
