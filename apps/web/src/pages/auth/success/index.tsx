import { useEffect } from "react";
import {
    getLocalStorageItem,
    LOCAL_STORAGE_KEYS,
    removeLocalStorageItem,
} from "@/lib/constants/localStorage";
import { ErrorPage } from "@/components/pages/Error";
import { useTranslation } from "react-i18next";
import { useSearchParams } from "react-router-dom";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";

const AuthSuccessPage = () => {
    const { t } = useTranslation();
    const [searchParams] = useSearchParams();

    const redirectUrl = getLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);

    useEffect(() => {
        if (!redirectUrl) {
            return;
        }
        // Ensure the redirect URL starts with '/' for absolute path
        const absoluteRedirectUrl = redirectUrl.startsWith("/") ? redirectUrl : `/${redirectUrl}`;
        window.location.href = `${absoluteRedirectUrl}?${searchParams.toString()}`;
        removeLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);
        queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
    }, [redirectUrl, searchParams]);

    if (!redirectUrl) {
        return (
            <ErrorPage
                title={t("errors.notFound.title")}
                description={t("errors.notFound.description")}
            />
        );
    }

    return (
        <div className="min-h-screen flex flex-col">
            <h1>You're being redirected. Please wait...</h1>
        </div>
    );
};

export default AuthSuccessPage;
