import { useEffect, useState } from "react";
import { redirect, useSearchParams } from "react-router-dom";
import { useVerifyGoogleOAuth } from "@/components/pages/GoogleOAuth/useVerifyGoogleOAuth";
import { FaviconHelmet } from '@/components/providers/helmets/FaviconHelmet';
import { GoogleOAuth } from '@/components/pages/GoogleOAuth/GoogleOAuth';
import { GoogleOAuthError } from '@/components/pages/GoogleOAuth/GoogleOAuthError';
import { useTranslation } from 'react-i18next';
import { getLocalStorageItem, LOCAL_STORAGE_KEYS, removeLocalStorageItem } from "@/lib/constants/localStorage";

const GoogleOAuthPage = () => {
    const { t } = useTranslation();

    // Get query params
    const [errorType, setErrorType] = useState<string | null>(null);
    const [searchParams] = useSearchParams();

    const accessToken = searchParams.get("access_token");
    const expiresIn = searchParams.get("expires_in");


    const { verifyGoogleOAuth, isPending } = useVerifyGoogleOAuth();

    useEffect(() => {
        const init = async () => {

            if (isPending || !accessToken || errorType || !expiresIn) {
                return;
            }

            // Check for redirect url in localstorage
            const redirectUrl = getLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);
            console.log(redirectUrl)
            if (!redirectUrl) {
                setErrorType("missingRedirect");
                return;
            }

            // Remove redirect url from localstorage    
            removeLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);

            const queryString = new URLSearchParams({
                access_token: accessToken,
                expires_in: expiresIn,
            }).toString();

            // Redirect to redirect url
            window.location.href = redirectUrl + "?" + queryString;
        }

        init();
    }, [isPending, verifyGoogleOAuth, errorType, accessToken, expiresIn]);

    // Determine which component to render
    const renderPage = () => {
        // Case 1: Missing both parameters
        if (!accessToken) {
            return <GoogleOAuthError errorType="invalidParams" />;
        }

        // Case 2: Missing redirect URL
        if (errorType === "missingRedirect") {
            return <GoogleOAuthError errorType="missingRedirect" />;
        }

        // Case 5: Valid parameters - show loading state
        return <GoogleOAuth />;
    };

    // Dynamic helmet based on state
    const getHelmetTitle = () => {
        if (!accessToken || errorType) {
            return `${t('oauth.google.error.title')} | ${t('common.appName')}`;
        }
        return `${t('oauth.google.title')} | ${t('common.appName')}`;
    };

    const getHelmetDescription = () => {
        if (errorType === "missingRedirect") {
            return t('oauth.google.error.missingRedirect');
        }
        if (!accessToken) {
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