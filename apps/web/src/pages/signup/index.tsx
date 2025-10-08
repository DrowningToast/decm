import { Button } from "@/components/ui/button";
import { Link, } from "@/router";
import { useTranslation } from "react-i18next";
import { LanguageSwitcher } from "@/components/LanguageSwitcher";
import { env } from "@/config/env";
import { LOCAL_STORAGE_KEYS, setLocalStorageItem } from "@/lib/constants/localStorage";
import { OnboardMethods } from "../onboard/[method]";

const SignUpPage = () => {
    const { t } = useTranslation();

    const handleRequestGoogleOAuthUrl = async () => {
        // open new tab with the url
        setLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT, `/onboard/${OnboardMethods.GOOGLE}`);
        window.location.href = `${env.VITE_CORE_BACKEND_API}/api/v1/auth/request-google-oauth`;
    }

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
                        <h1 className="text-4xl font-bold">{t('signup.title')}</h1>
                        <p className="text-muted-foreground">
                            {t('signup.subtitle')}
                        </p>
                    </div>

                    <div className="space-y-4">
                        <Button className="w-full" size="lg" variant="primary">
                            {t('signup.walletButton')}
                        </Button>

                        <div className="relative">
                            <div className="absolute inset-0 flex items-center">
                                <span className="w-full border-t" />
                            </div>
                            <div className="relative flex justify-center text-xs uppercase">
                                <span className="bg-background px-2 text-muted-foreground">
                                    {t('signup.divider')}
                                </span>
                            </div>
                        </div>

                        <Button variant="secondary-dark" className="w-full" size="lg" onClick={handleRequestGoogleOAuthUrl}>
                            {t('signup.googleButton')}
                        </Button>
                    </div>

                    <p className="text-xs text-center text-muted-foreground">
                        <Link to="/signin" className="hover:underline">
                            {t('auth.hasAccount')}
                        </Link>
                    </p>

                    <p className="text-xs text-center text-muted-foreground mt-4">
                        {t('signup.terms')}
                    </p>
                </div>
            </div>
        </div>
    );
};

export default SignUpPage;