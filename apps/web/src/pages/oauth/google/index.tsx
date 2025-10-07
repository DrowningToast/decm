import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useVerifyGoogleOAuth } from "@/components/pages/GoogleOAuth/useVerifyGoogleOAuth";
import { FaviconHelmet } from '@/components/providers/helmets/FaviconHelmet';
import { GoogleOAuth } from '@/components/pages/GoogleOAuth/GoogleOAuth';
import { GoogleOAuthError } from '@/components/pages/GoogleOAuth/GoogleOAuthError';
import { useTranslation } from 'react-i18next';

const GoogleOAuthPage = () => {
    const { t } = useTranslation();

    // Get query params
    const [isError, setIsError] = useState(false);
    const [searchParams] = useSearchParams();
    const code = searchParams.get("code");
    const state = searchParams.get("state");

    const { verifyGoogleOAuth, isPending } = useVerifyGoogleOAuth();

    useEffect(() => {
        const init = async () => {
            if (isPending || !code || !state || isError) {
                return;
            }

            try {
                await verifyGoogleOAuth({ code, state }).then((response) => {
                    console.log("OAuth verification successful:", response);
                    setIsError(false);
                });
            } catch (error) {
                console.error("OAuth verification failed:", error);
                // Error handling is done by the mutation
                setIsError(true);
            }
        }

        init();
    }, [code, isPending, state, verifyGoogleOAuth, isError]);

    // Determine which component to render
    const renderPage = () => {
        // Case 1: Missing both parameters
        if (!code && !state) {
            return <GoogleOAuthError errorType="invalidParams" />;
        }

        // Case 2: Missing code
        if (!code) {
            return <GoogleOAuthError errorType="missingCode" />;
        }

        // Case 3: Missing state
        if (!state) {
            return <GoogleOAuthError errorType="missingState" />;
        }

        // Case 4: Verification failed (API error)
        if (isError) {
            return <GoogleOAuthError errorType="verificationFailed" />;
        }

        // Case 5: Valid parameters - show loading state
        return <GoogleOAuth />;
    };

    // Dynamic helmet based on state
    const getHelmetTitle = () => {
        if (!code || !state || isError) {
            return `${t('oauth.google.error.title')} | ${t('common.appName')}`;
        }
        return `${t('oauth.google.title')} | ${t('common.appName')}`;
    };

    const getHelmetDescription = () => {
        if (isError) {
            return t('oauth.google.error.verificationFailed');
        }
        if (!code || !state) {
            return t('oauth.google.error.invalidParams');
        }
        return t('oauth.google.subtitle');
    };

    return (
        <>
            <FaviconHelmet
                title={getHelmetTitle()}
                description={getHelmetDescription()}
            />
            {renderPage()}
        </>
    );
}

export default GoogleOAuthPage;