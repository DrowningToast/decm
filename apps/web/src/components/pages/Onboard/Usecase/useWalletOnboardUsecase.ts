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

export const useWalletOnboardUsecase = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { upsertProfile } = useSignup();

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
        if (account?.credential_id) {
            try {
                await upsertProfile({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: signSignature,
                    profile: {
                        authentication_credential_id: account.credential_id,
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
                    handleUniversalError(t, error);
                }
                return;
            }
        }

        toast.success(t("flow.wallet.create_profile_success"));
        navigate("/app");
    };

    return {
        usecaseAsync,
    };
};
