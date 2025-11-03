import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { VerifySignMessagePage } from "@/components/pages/Auth/VerifySignMessagePage";
import { useTranslation } from "react-i18next";
import { useWalletClient } from "wagmi";
import { useEffect } from "react";
import { useNavigate } from "@/router";
import { toast } from "sonner";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { OnboardRegistrationMethod } from "@decm/api";
import { handleUniversalError } from "@/common/Err";
import { USECASE_IDS } from "@/constants/usecase";
import { useCheckOnboardStatus } from "@/components/pages/Onboard/useCheckOnboardStatus";
import { OnboardMethods } from "@/pages/onboard/[method]";
import { TOAST_USECASE_VIEWMODEL } from "@/constants/toast";
import { useSignout } from "@/components/useSignout";

const VerifyMessagePage = () => {
    const { t } = useTranslation();

    const navigate = useNavigate();
    const { data: walletClient } = useWalletClient();
    const { signout } = useSignout();

    useEffect(() => {
        if (!walletClient) {
            toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.SIGN_IN].WALLET_NOT_CONNECTED));
            navigate("/");
        }
    }, [navigate, t, walletClient]);

    const [authSignSignature] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
        undefined,
    );

    const { onboardStatus, isLoading, error } = useCheckOnboardStatus(
        {
            method: OnboardRegistrationMethod.RegistrationMethodWallet,
            signSignature: authSignSignature,
        },
        !!authSignSignature,
    );

    useEffect(() => {
        if (isLoading) {
            return;
        }
        if (error) {
            handleUniversalError(
                t,
                error,
                {
                    onUnauthorized: () => {
                        navigate("/");
                    },
                },
                USECASE_IDS.CHECK_ONBOARD_STATUS,
            );
        }
        const hasAuthenticationCredentialId = !!onboardStatus?.authentication_credential_id;
        const hasProfileId = !!onboardStatus?.profile_id;
        if (hasAuthenticationCredentialId && hasProfileId) {
            navigate("/app");
        }
        if (hasAuthenticationCredentialId) {
            navigate("/onboard/:method", {
                params: {
                    method: OnboardMethods.WALLET,
                },
            });
        }
        if (!hasAuthenticationCredentialId && !hasProfileId) {
            toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.SIGN_IN].NOTFOUND));
            signout();
            navigate("/signup");
        }
    }, [
        error,
        navigate,
        t,
        isLoading,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
        signout,
    ]);

    return (
        <>
            <FaviconHelmet
                title={`${t("verify.title")} | ${t("common.appName")}`}
                description={t("verify.subtitle")}
            />
            <VerifySignMessagePage />
        </>
    );
};

export default VerifyMessagePage;
