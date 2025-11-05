import { useEffect } from "react";
import {
    getLocalStorageItem,
    LOCAL_STORAGE_KEYS,
    removeLocalStorageItem,
} from "@/lib/constants/localStorage";
import { useTranslation } from "react-i18next";
import { useSearchParams } from "react-router-dom";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { handleUniversalError } from "@/common/Err";
import { ErrorPage } from "@/components/pages/Error";

const AuthSuccessPage = () => {
    const { t } = useTranslation();
    const [searchParams] = useSearchParams();

    const redirectUrl = getLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);

    useEffect(() => {
        const init = async () => {
            try {
                if (!redirectUrl) {
                    return;
                }
                // Ensure the redirect URL starts with '/' for absolute path
                const absoluteRedirectUrl = redirectUrl.startsWith("/")
                    ? redirectUrl
                    : `/${redirectUrl}`;
                await queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
                removeLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);
                window.location.href = `${absoluteRedirectUrl}?${searchParams.toString()}`;
            } catch (error) {
                console.error(error);
                if (error instanceof Error) {
                    handleUniversalError(t, error);
                }
            }
        };
        init();
    }, [redirectUrl, searchParams, t]);

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
