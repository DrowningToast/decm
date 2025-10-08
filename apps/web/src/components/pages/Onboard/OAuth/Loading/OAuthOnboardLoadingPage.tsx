import { Typography } from "@/components/typography/typography";
import { useCheckOnboardStatus } from "@/components/pages/Onboard/useCheckOnboardStatus";
import { OnboardRegistrationMethod } from "@decm/api";
import { useContext, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { Error } from "@/components/pages/Error";
import { useNavigate } from "@/router";
import { OnboardPageContext } from "@/pages/onboard/[method]";

export const OAuthOnboardLoadingPage = () => {

    const { setStep } = useContext(OnboardPageContext)
    const [searchParams] = useSearchParams()
    const accessToken = searchParams.get("access_token")
    const expiresIn = searchParams.get("expires_in")

    const { checkOnboardStatus } = useCheckOnboardStatus()
    const navigate = useNavigate()

    useEffect(() => {
        const init = async () => {
            if (!accessToken) {
                return;
            }

            const response = await checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodGoogle,
                accessToken: accessToken,
                expiresIn: expiresIn ? parseInt(expiresIn) : undefined,
            });

            if (response.is_exists) {
                return navigate("/app");
            }

            setStep(1);

        };

        init();
    }, [accessToken, checkOnboardStatus, expiresIn, navigate, setStep]);

    if (accessToken) {
        return <Error />;
    }

    // TODO: Implement
    return (
        <div>
            <Typography variant="header" tag="h1">Loading...</Typography>
        </div>
    )
}