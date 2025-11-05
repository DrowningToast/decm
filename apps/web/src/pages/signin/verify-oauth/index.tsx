import { useSignout } from "@/components/useSignout";
import { TOAST_USECASE_VIEWMODEL } from "@/constants/toast";
import { USECASE_IDS } from "@/constants/usecase";
import { coreApiClient } from "@/lib/api/api";
import { authService } from "@/services/AuthService";
import { OnboardRegistrationMethod } from "@decm/api";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router-dom";
import { toast } from "sonner";

const VerifyOauthPage = () => {
    const [searchParams] = useSearchParams();
    const accessToken = searchParams.get("access_token");
    const expiresIn = searchParams.get("expires_in");
    const navigate = useNavigate();
    const { t } = useTranslation();

    useEffect(() => {
        const init = async () => {
            if (!accessToken || !expiresIn) {
                return;
            }

            const status = await coreApiClient.v1.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                access_token: accessToken,
                expires_in: parseInt(expiresIn),
            });

            if (!status) {
                console.error("Failed to check onboard status");
                signout();
                navigate("/signup");
                return;
            }

            if (status.authentication_credential_id) {
                navigate("/app");
            } else {
                toast.error(t(TOAST_USECASE_VIEWMODEL[USECASE_IDS.SIGN_IN].NOTFOUND));
                await authService.signOut({ showSuccessToast: false });
                navigate("/signup");
            }
        };
        init();
    }, [accessToken, expiresIn, navigate, t]);

    return <span>Loading...</span>;
};

export default VerifyOauthPage;
