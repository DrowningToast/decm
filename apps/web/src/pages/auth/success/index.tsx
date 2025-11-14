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
import { Typography } from "@/components/typography/typography";

const AuthSuccessPage = () => {
    const { t } = useTranslation();
    const [searchParams] = useSearchParams();

    const redirectUrl =
        getLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT) ?? "";

    useEffect(() => {
        const init = async () => {
            if (!redirectUrl) {
                return;
            }
            try {
                const absoluteRedirectUrl = redirectUrl.startsWith("/")
                    ? redirectUrl
                    : `/${redirectUrl}`;

                const url = new URL(absoluteRedirectUrl, window.location.origin);

                searchParams.forEach((value, key) => {
                    url.searchParams.set(key, value);
                });

                await queryClient.invalidateQueries({ queryKey: QUERY_KEY.user.profile });
                removeLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);

                window.location.href = `${url.pathname}${url.search}`;
            } catch (error) {
                console.error(error);
                if (error instanceof Error) {
                    handleUniversalError(t, error);
                }
            }
        };

        void init();
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
            <Typography variant="header" tag="h1">
                {t("auth.success.title")}
            </Typography>
        </div>
    );
};

export default AuthSuccessPage;
