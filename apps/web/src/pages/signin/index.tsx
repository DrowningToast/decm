import { Button } from "@/components/ui/button";
import { Link } from "@/router";
import { useTranslation } from "react-i18next";
import { LanguageSwitcher } from "@/components/LanguageSwitcher";
import { setLocalStorageItem } from "@/lib/constants/localStorage";
import { env } from "@/config/env";
import { LOCAL_STORAGE_KEYS } from "../../lib/constants/localStorage";
import { useCheckOnboardStatus } from "@/components/pages/Onboard/useCheckOnboardStatus";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { authService } from "@/services/AuthService";

const SignInPage = () => {
    const { t } = useTranslation();

    const navigate = useNavigate();
    const { onboardStatus } = useCheckOnboardStatus();
    // Handle if the user is already authenticated
    useEffect(() => {
        const init = async () => {
            // If the user is not authenticated, do nothing
            if (!onboardStatus?.authentication_credential_id) {
                return;
            }
            // check if the user has a profile or not
            if (!onboardStatus?.profile_id) {
                authService.createProfile(onboardStatus.authentication_credential_id, {});
            }
            navigate("/app");
        };

        init();
    }, [navigate, onboardStatus]);

    const handleRequestGoogleOAuthUrl = async () => {
        // open new tab with the url
        setLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT, `/app`);
        window.location.href = `${env.VITE_CORE_BACKEND_API}/api/v1/auth/request-google-oauth`;
    };

    return (
        <div className="min-h-screen flex flex-col">
            {/* Header with Language Switcher */}
            <div className="flex justify-end p-4">
                <LanguageSwitcher />
            </div>

            {/* Main Content */}
            <div className="flex-1 flex items-center justify-center">
                <div className="flex flex-col gap-y-6 max-w-md w-full px-6">
                    <div className="text-center space-y-2">
                        <h1 className="text-4xl font-bold">{t("signin.title")}</h1>
                        <p className="text-muted-foreground">{t("signin.subtitle")}</p>
                    </div>

                    <div className="space-y-4">
                        <Button className="w-full" size="lg" variant="primary">
                            {t("signup.walletButton")}
                        </Button>

                        <div className="relative">
                            <div className="absolute inset-0 flex items-center">
                                <span className="w-full border-t" />
                            </div>
                            <div className="relative flex justify-center text-xs uppercase">
                                <span className="bg-background px-2 text-muted-foreground">
                                    {t("signup.divider")}
                                </span>
                            </div>
                        </div>

                        <Button
                            variant="secondary-dark"
                            className="w-full"
                            size="lg"
                            onClick={handleRequestGoogleOAuthUrl}
                        >
                            {t("signup.googleButton")}
                        </Button>
                    </div>

                    <p className="text-xs text-center text-muted-foreground">
                        <Link to="/signin" className="hover:underline">
                            {t("auth.hasAccount")}
                        </Link>
                    </p>

                    <p className="text-xs text-center text-muted-foreground mt-4">
                        {t("signup.terms")}
                    </p>
                </div>
            </div>
        </div>
    );
};

export default SignInPage;
