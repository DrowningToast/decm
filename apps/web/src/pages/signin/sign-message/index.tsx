import { FaviconHelmet } from "@/components/providers/helmets/FaviconHelmet";
import { VerifySignMessagePage } from "@/components/pages/Auth/VerifySignMessagePage";
import { useTranslation } from "react-i18next";
import { useAccount } from "wagmi";
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
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { onboardService } from "@/services/OnboardService";

const VerifyMessagePage = () => {
    const { t } = useTranslation();

    const navigate = useNavigate();
    const { isReconnecting, isDisconnected, isConnecting, address } = useAccount();
    const { signout } = useSignout();

    const [authSignSignature] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
        undefined,
    );

    const {
        onboardStatus,
        isQueryFetching: isOnboardStatusLoading,
        error,
    } = useCheckOnboardStatus(
        {
            method: OnboardRegistrationMethod.RegistrationMethodWallet,
            signSignature: authSignSignature,
        },
        !!authSignSignature,
    );

    const isLoading = isReconnecting || isOnboardStatusLoading || isDisconnected || isConnecting;

    // check if the user is onboarded
    useEffect(() => {
        const init = async () => {
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
                return;
            }
            if (!address) {
                toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.SIGN_IN].WALLET_NOT_CONNECTED));
                navigate("/");
            }
            const hasAuthenticationCredentialId = !!onboardStatus?.authentication_credential_id;
            const hasProfileId = !!onboardStatus?.profile_id;
            if (hasAuthenticationCredentialId && hasProfileId) {
                await onboardService.checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodWallet,
                    signSignature: authSignSignature ?? "",
                });
                await queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
                await navigate("/app");
                return;
            }
            if (hasAuthenticationCredentialId) {
                toast.info(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.SIGN_IN].NOTFOUND));
                await navigate(`/onboard/:method`, {
                    params: {
                        method: OnboardMethods.WALLET,
                    },
                });
                return;
            }
        };
        init();
    }, [
        error,
        navigate,
        t,
        isLoading,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
        signout,
        address,
        authSignSignature,
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
