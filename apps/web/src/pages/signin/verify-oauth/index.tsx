import { TOAST_USECASE_VIEWMODEL } from "@/constants/toast";
import { USECASE_IDS } from "@/constants/usecase";
import { coreApiClient } from "@/lib/api/api";
import { authService } from "@/services/AuthService/AuthService";
import { OnboardRegistrationMethod } from "@decm/api";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import { AxiosError } from "axios";
import { Typography } from "@/components/typography/typography";

const VerifyOauthPage = () => {
    const [searchParams] = useSearchParams();
    const accessToken = searchParams.get("access_token");
    const expiresIn = searchParams.get("expires_in");
    const navigate = useNavigate();
    const { t } = useTranslation();

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
                const status = await coreApiClient.v1.checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    access_token: accessToken,
                    expires_in: expiresInNumber,
                });

                if (!status) {
                    toast.error(t("errors.serverError"));
                    await authService.signOut({ showSuccessToast: false });
                    navigate("/signup", { replace: true });
                    return;
                }

                if (status.authentication_credential_id) {
                    navigate("/app");
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
    }, [accessToken, expiresIn, navigate, t]);

    return (
        <Typography variant="text" tag="span">
            Loading...
        </Typography>
    );
};

export default VerifyOauthPage;
