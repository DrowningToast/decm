import { OAuthOnboardPasswordPage } from "@/components/pages/Onboard/OAuth/Password/OAuthOnboardPasswordPage";
import { OAuthOnboardProvider } from "@/components/pages/Onboard/OAuth/OAuthOnboardProvider";
import NotFoundPage from "@/pages/404";
import { useParams } from "@/router";
import { createContext, useState } from "react";
// import { NotFoundPage } from "@/components/pages/NotFoundPage";

export const OnboardMethods = {
    WALLET: "wallet",
    GOOGLE: "google",
} as const

interface OnboardStep {
    render: React.FC<{ step: number }>;
}

type OnboardMethod = typeof OnboardMethods[keyof typeof OnboardMethods];

type OnboardSteps = Record<OnboardMethod, Record<number, OnboardStep> & { Parent: React.FC }>

const OnboardSteps: OnboardSteps = {
    [OnboardMethods.WALLET]: {
        0: {
            // Request the user to sign the message
            // Action: Authenticate with wallet
            render: () => <div>Sign Message</div>,
        },
        1: {
            // Find if the account is already created
            // If not, lands here
            // UI: Show profile form
            render: () => <div>Profile Info</div>,
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
        // Action: Checks if the account is already created
        0: {
            render: () => <div>Loading...</div>,
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
            render: () => <div>Profile Info</div>,
        },
        // UI: Show confirmation
        // Action: Create profile + create account + redirect to home
        3: {
            render: () => <div>Confirmation</div>,
        },
        Parent: OAuthOnboardProvider
    },
}

type OnboardPageContextType = {
    step: number;
    setStep: (step: number) => void;
}

const OnboardPageContext = createContext<OnboardPageContextType>({} as OnboardPageContextType);

export { OnboardPageContext }


const OnboardingPage = () => {
    const { method } = useParams<"/onboard/:method">("/onboard/:method");
    const [step, setStep] = useState<number>(0);

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
        <OnboardPageContext.Provider value={{ step, setStep }}>
            {
                renderStep()
            }
        </OnboardPageContext.Provider>
    )
};

export default OnboardingPage;