import { Typography } from "@/components/typography/typography";
import { useCheckOnboardStatus } from "@/components/pages/Onboard/useCheckOnboardStatus";
import { OnboardRegistrationMethod } from "@decm/api";
import { useContext, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { Error } from "@/components/pages/Error";
import { useNavigate } from "@/router";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { AxiosError } from "axios";
import { useTranslation } from "react-i18next";

type CheckStatusError = "unauthenticated_response" | "internal_error_response" | 'missing_access_token' | 'expired_token'

export const OAuthOnboardLoadingPage = () => {
    const { t } = useTranslation();
    const { setStep } = useContext(OnboardPageContext)
    const [searchParams] = useSearchParams()
    const accessToken = searchParams.get("access_token")
    const expiresIn = searchParams.get("expires_in")

    const { checkOnboardStatus } = useCheckOnboardStatus()
    const navigate = useNavigate()

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
            if (!accessToken) {
                setError("missing_access_token");
                return;
            }

            try {
                const response = await checkOnboardStatus({
                    method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                    accessToken: accessToken,
                    expiresIn: expiresIn ? parseInt(expiresIn) : undefined,
                });

                if (response.is_exists) {
                    return navigate("/app");
                }

                setStep(1);
            } catch (error) {
                if (error instanceof AxiosError) {
                    switch (error.response?.status) {
                        case 401:
                            setError("unauthenticated_response");
                            return
                        case 400:
                            setError("expired_token");
                            return
                        default:
                            setError("internal_error_response");
                            return
                    }
                }
            };
        }

        init();
    }, [accessToken, checkOnboardStatus, expiresIn, navigate, setStep, error]);

    if (error) {
        return <Error title={errorTitle} description={errorMessage} />;
    }

    // TODO: Implement
    return (
        <div>
            <Typography variant="header" tag="h1">Loading...</Typography>
        </div>
    )
}