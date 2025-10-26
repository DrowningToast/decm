import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { BaseProfilePage } from "../../BaseProfilePage";
import { OnboardRegistrationMethod } from "@decm/api";

export const OAuthOnboardProfilePage: React.FC = () => {

    // if already has account, disable back button
    const { onboardStatus } = useContext(OnboardPageContext)
    const hasAccount = onboardStatus?.authentication_credential_id !== undefined

    const { step, setStep } = useContext(OnboardPageContext);
    const { t } = useTranslation();

    const onConfirm = () => {
        setStep(step + 1)
    }

    const onBack = () => {
        setStep(step - 1)
    }

    return (
        <BaseProfilePage
            t={t}
            onConfirm={onConfirm}
            onBack={onBack}
            hasAccount={hasAccount}
            method={OnboardRegistrationMethod.RegistrationMethodGoogle}
        />
    )
};