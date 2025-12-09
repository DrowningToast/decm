import { TOAST_USECASE_VIEWMODEL } from "@/constants/toast";
import { USECASE_IDS } from "@/constants/usecase";
import { authService } from "@/services/services";
import { OnboardRegistrationMethod } from "@decm/api";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import { AxiosError } from "axios";
import { Typography } from "@/components/typography/typography";
import { useOAuthOnboardUsecase } from "@/components/pages/Onboard/Usecase/useOAuthOnboardUsecase";
import { useMyProfile } from "@/hooks/useMyProfile";
import { useCheckOnboardStatus } from "@/components/pages/Onboard/useCheckOnboardStatus";

const VerifyOauthPage = () => {
    const [searchParams] = useSearchParams();
    const accessToken = searchParams.get("access_token");
    const expiresIn = searchParams.get("expires_in");
    const navigate = useNavigate();
    const { t } = useTranslation();
    const { refetch: refetchMyProfile } = useMyProfile();
    const { checkOnboardStatus } = useCheckOnboardStatus(undefined, false);

    useEffect(() => {
        const init = async () => {
            try {
                // Validate query parameters
                if (!accessToken || !expiresIn) {
                    toast.error(t("oauth.google.error.invalidParams"));
                    await authService.signOut({ showSuccessToast: false });
                    navigate("/signup", { replace: true });
                    return;
                }

                // Validate and parse expiresIn
                const expiresInNumber = parseInt(expiresIn, 10);
                if (isNaN(expiresInNumber)) {
                    toast.error(t("oauth.google.error.invalidParams"));
                    await authService.signOut({ showSuccessToast: false });
                    navigate("/signup", { replace: true });
                    return;
                }

                // Call API to check onboard status
                const status = await checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: accessToken,
                    expiresIn: expiresInNumber,
                });

                if (!status) {
                    toast.error(t("errors.serverError"));
                    await authService.signOut({ showSuccessToast: false });
                    navigate("/signup", { replace: true });
                    return;
                }
                await useOAuthOnboardUsecase;
                if (status.authentication_credential_id) {
                    await refetchMyProfile();
                    await navigate("/app");
                } else {
                    toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.SIGN_IN].NOTFOUND));
                    await authService.signOut({ showSuccessToast: false });
                    navigate("/signup", { replace: true });
                }
            } catch (error) {
                // Map known HTTP status codes to localized toast messages
                if (error instanceof AxiosError && error.response) {
                    switch (error.response.status) {
                        case 400:
                            toast.error(t("errors.invalidInput"));
                            break;
                        case 401:
                            toast.error(t("errors.unauthorized"));
                            break;
                        case 403:
                            toast.error(t("errors.forbidden"));
                            break;
                        case 404:
                            toast.error(t("errors.notFound"));
                            break;
                        case 409:
                            toast.error(t("errors.conflict"));
                            break;
                        case 500:
                            toast.error(t("errors.serverError"));
                            break;
                        default:
                            toast.error(t("errors.generic"));
                    }
                } else {
                    // Network or unknown error
                    toast.error(t("errors.network"));
                }

                await authService.signOut({ showSuccessToast: false });
                navigate("/signup", { replace: true });
            }
        };
        init();
    }, [accessToken, expiresIn, navigate, refetchMyProfile, t, checkOnboardStatus]);

    return (
        <Typography variant="text" tag="span">
            Loading...
        </Typography>
    );
};

export default VerifyOauthPage;
