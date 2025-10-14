import { OAuthOnboardPasswordPage } from "@/components/pages/Onboard/OAuth/Password/OAuthOnboardPasswordPage";
import { OAuthOnboardProvider } from "@/components/pages/Onboard/OAuth/OAuthOnboardContext";
import NotFoundPage from "@/pages/404";
import { useParams } from "@/router";
import { createContext, useEffect, useState } from "react";
import { OAuthOnboardLoadingPage } from "@/components/pages/Onboard/OAuth/Loading/OAuthOnboardLoadingPage";
import { ProfilePage } from "@/components/pages/Onboard/ProfilePage";
import { OAuthOnboardConfirmPage } from "@/components/pages/Onboard/OAuth/Confirm/OAuthOnboardConfirmPage";
import { OnboardRegistrationMethod, type OnboardCheckOnboardStatusResponse } from "@decm/api";
import { useCheckOnboardStatus } from "@/components/pages/Onboard/useCheckOnboardStatus";
import { useNavigate, useSearchParams } from "react-router-dom";
import { handleUniversalError } from "@/common/Err";
import { useTranslation } from "react-i18next";
import { USECASE_IDS } from "@/constants/usecase";

export const OnboardMethods = {
    WALLET: "wallet",
    GOOGLE: "google",
} as const

interface OnboardStep {
    render: React.FC<{ step: number }>;
}

type OnboardMethod = typeof OnboardMethods[keyof typeof OnboardMethods];

type OnboardSteps = Record<OnboardMethod, Record<number, OnboardStep> & { Parent: React.FC<React.PropsWithChildren> }>

const OnboardSteps: OnboardSteps = {
    [OnboardMethods.WALLET]: {
        0: {
            // Request the user to sign the message
            // Action: Request to sign message, then check the onboard status
            render: () => <div>Sign Message</div>,
        },
        1: {
            // Find if the account is already created
            // If not, lands here
            // UI: Show profile form
            render: () => <ProfilePage />,
        },
        2: {
            // UI: Show confirmation
            // Action: Create account + create profile
            render: () => <div></div>,
        },
        Parent: OAuthOnboardProvider,
    },
    [OnboardMethods.GOOGLE]: {
        // Precondition: Lands here from oauth-success/google
        // UI: Loading page
        // Action: Authenticate the access token, then check the onboard status
        0: {
            render: () => <OAuthOnboardLoadingPage />,
        },
        // If the account isn't created, lands here
        // UI: Show password page
        // Action: Set local password value
        1: {
            render: () => <OAuthOnboardPasswordPage />,
        },
        // If the account is already created, lands here
        // UI: Show profile form
        // Action: Set local profile value
        2: {
            render: () => <ProfilePage />,
        },
        // UI: Show confirmation
        // Action: Create profile + create account + redirect to home
        3: {
            render: () => <OAuthOnboardConfirmPage />,
        },
        Parent: OAuthOnboardProvider
    },
}

type OnboardPageContextType = {
    step: number;
    setStep: (step: number) => void;

    accessToken?: string;
    expiresIn?: number;

    onboardStatus?: OnboardCheckOnboardStatusResponse;
    isStatusLoading: boolean
}

const OnboardPageContext = createContext<OnboardPageContextType>({} as OnboardPageContextType);

export { OnboardPageContext }


const OnboardingPage = () => {
    const { method } = useParams<"/onboard/:method">("/onboard/:method");
    const [step, setStep] = useState<number>(0);
    const navigate = useNavigate();

    const { t } = useTranslation();
    const [searchParams] = useSearchParams();
    const accessToken = searchParams.get("access_token");
    const expiresIn = searchParams.get("expires_in");

    const { onboardStatus, isLoading, error } = useCheckOnboardStatus({
        method: OnboardRegistrationMethod.RegistrationMethodGoogle,
        accessToken: accessToken ?? "",
        expiresIn: expiresIn ? parseInt(expiresIn) : 0,
    });
    // handle error
    useEffect(() => {
        if (!error) {
            return
        }

        handleUniversalError(t, error, {
            onUnauthorized: () => {
                navigate("/");
            },
        }, USECASE_IDS.CHECK_ONBOARD_STATUS);


    }, [error, navigate, t])

    if (!method || !Object.values(OnboardMethods).includes(method as OnboardMethod)) {
        return <NotFoundPage />
    }

    const steps = OnboardSteps[method as OnboardMethod];

    if (!steps) {
        return <NotFoundPage />
    }

    const renderStep = () => {
        const StepComponent = steps[step].render;
        return <StepComponent step={step} />;
    }

    return (
        <OnboardPageContext.Provider value={{
            step, setStep,
            accessToken: accessToken ?? undefined,
            expiresIn: expiresIn ? parseInt(expiresIn) : undefined,
            onboardStatus,
            isStatusLoading: isLoading,
        }}>
            <steps.Parent>
                {
                    renderStep()
                }
            </steps.Parent>
        </OnboardPageContext.Provider>
    )
};

export default OnboardingPage;