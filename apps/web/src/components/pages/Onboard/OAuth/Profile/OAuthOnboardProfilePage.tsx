import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { FormField, FormItem, FormControl, FormMessage } from "@/components/ui/form";
import { OAuthOnboardContext } from "../OAuthOnboardContext";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { LogoutButton } from "../../../../LogoutButton";
import { BaseProfilePage } from "../../BaseProfilePage";
import type { ProfileSchema } from "../../ProfilePage";

export const OAuthOnboardProfilePage: React.FC = () => {

    // if already has account, disable back button
    const { onboardStatus, method } = useContext(OnboardPageContext)
    const hasAccount = onboardStatus?.authentication_credential_id !== undefined

    const { setStep } = useContext(OnboardPageContext);
    const { form } = useContext(OAuthOnboardContext);
    const { t } = useTranslation();

    const onConfirm = () => {
        setStep(3)
    }

    const onBack = () => {
        setStep(1)
    }

    return (
        <BaseProfilePage<OAuthOnboardPr>
            form={form}
            t={t}
            onConfirm={onConfirm}
            onBack={onBack}
            hasAccount={hasAccount}
            method={method}
        />
    )
};