import { Typography } from "@/components/typography/typography";
import { useContext, useEffect, useMemo, useState } from "react";
import { ErrorPage } from "@/components/pages/Error";
import { useNavigate } from "@/router";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { useTranslation } from "react-i18next";

type CheckStatusError =
    | "unauthenticated_response"
    | "internal_error_response"
    | "missing_access_token"
    | "expired_token";

export const OAuthOnboardLoadingPage = () => {
    const { t } = useTranslation();
    const { accessToken, expiresIn, setStep, onboardStatus, isStatusLoading } =
        useContext(OnboardPageContext);

    const navigate = useNavigate();

    const [error, setError] = useState<CheckStatusError | null>(null);
    const errorTitle = useMemo(() => {
        switch (error) {
            case "missing_access_token":
                return t("oauth.google.error.missingAccess");
            case "unauthenticated_response":
                return t("oauth.google.error.unauthenticated");
            case "internal_error_response":
                return t("oauth.google.error.internal");
            case "expired_token":
                return t("oauth.google.error.expired");
            default:
                return;
        }
    }, [error, t]);

    const errorMessage = useMemo(() => {
        switch (error) {
            case "missing_access_token":
                return t("oauth.google.error.reasons.missingAccess");
            case "unauthenticated_response":
                return t("oauth.google.error.reasons.unauthenticated");
            case "internal_error_response":
                return t("oauth.google.error.reasons.internal");
            case "expired_token":
                return t("oauth.google.error.reasons.expired");
            default:
                return;
        }
    }, [error, t]);

    useEffect(() => {
        const init = async () => {
            if (isStatusLoading) {
                return;
            }

            if (!accessToken) {
                setError("missing_access_token");
                return;
            }

            if (onboardStatus?.authentication_credential_id && onboardStatus?.profile_id) {
                return navigate("/app");
            }

            if (onboardStatus?.authentication_credential_id) {
                setStep(2);
                return;
            }

            setStep(1);
            return;
        };

        init();
    }, [
        accessToken,
        onboardStatus,
        expiresIn,
        navigate,
        setStep,
        error,
        isStatusLoading,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
    ]);

    if (error) {
        return <ErrorPage title={errorTitle} description={errorMessage} />;
    }

    // TODO: Implement
    return (
        <div>
            <Typography variant="header" tag="h1">
                Loading...
            </Typography>
        </div>
    );
};
