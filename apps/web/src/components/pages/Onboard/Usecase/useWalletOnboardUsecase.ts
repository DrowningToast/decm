import {
    OnboardRegistrationMethod,
    type OnboardCheckOnboardStatusResponse,
    type OnboardRegisterResponse,
} from "@decm/api";
import { handleUniversalError } from "@/common/Err";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { useSignup } from "../useSignup";
import { onboardService } from "@/services/OnboardService";
import { authService } from "@/services/AuthService";
import type { Profile } from "../ProfilePage";
import { useMyProfile } from "@/hooks/useMyProfile";

export const useWalletOnboardUsecase = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { upsertProfile } = useSignup();
    const { refetch: refetchMyProfile } = useMyProfile();

    const usecaseAsync = async (signSignature: string, profile: Profile) => {
        let onboardStatus: OnboardCheckOnboardStatusResponse | null = null;
        try {
            onboardStatus = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: signSignature,
            });
        } catch (error) {
            if (error instanceof Error) {
                handleUniversalError(t, error);
            }
            return;
        }
        let account: OnboardRegisterResponse | undefined = undefined;
        try {
            if (!onboardStatus?.authentication_credential_id) {
                account = await authService.createAccount({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: signSignature,
                });
            }
        } catch (error) {
            if (error instanceof Error) {
                handleUniversalError(t, error);
            }
            return;
        }

        // Derive credential ID from either newly created account or existing onboard status
        const credentialId = account?.credential_id ?? onboardStatus?.authentication_credential_id;

        if (!credentialId) {
            // This is a critical error - we should have a credential ID at this point
            toast.error(t("flow.wallet.create_profile_error"));
            console.error("Missing credential ID in wallet onboard flow");
            return;
        }

        try {
            await upsertProfile({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: signSignature,
                profile: {
                    authentication_credential_id: credentialId,
                    first_name: profile.firstName,
                    is_first_name_public: profile.isFirstNamePublic,
                    last_name: profile.lastName,
                    is_last_name_public: profile.isLastNamePublic,
                    email: profile.email,
                    is_email_public: profile.isEmailPublic,
                    phone_number: profile.phoneNumber,
                    is_phone_number_public: profile.isPhoneNumberPublic,
                },
            });
        } catch (error) {
            if (error instanceof Error) {
                handleUniversalError(t, error);
            }
            return;
        }

        // Only show success and navigate if profile was successfully created/updated
        toast.success(t("flow.wallet.create_profile_success"));
        await refetchMyProfile();
        await navigate("/app");
    };

    return { usecaseAsync };
};
