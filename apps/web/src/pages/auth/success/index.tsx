import { useEffect } from "react";
import { getLocalStorageItem, LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { useNavigate } from "react-router-dom";
import { ErrorPage } from "@/components/pages/Error";
import { useTranslation } from "react-i18next";
import { useSearchParams } from "react-router-dom";

const AuthSuccessPage = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();

    const redirectUrl = getLocalStorageItem(LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT);

    useEffect(() => {
        if (!redirectUrl) {
            return;
        }
        navigate(`${redirectUrl}?${searchParams.toString()}`);
    }, [navigate, redirectUrl, searchParams]);

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
