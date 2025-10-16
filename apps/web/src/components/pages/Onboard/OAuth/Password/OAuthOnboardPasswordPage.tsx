import { useContext, useState } from "react";
import { OAuthOnboardContext } from "../OAuthOnboardContext";
import { PinPasswordPage } from "./PinPasswordPage";
import { PasswordInputPage } from "./PasswordInputPage";
import type { OAuthOnboardPasswordType } from "../constants";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { useNavigate } from "react-router-dom";

export const OAuthOnboardPasswordPage = () => {
    const navigate = useNavigate();
    const { setStep, accessToken, expiresIn } = useContext(OnboardPageContext)
    const { form } = useContext(OAuthOnboardContext);
    const [passwordType, setPasswordType] = useState<OAuthOnboardPasswordType>("PINS");

    console.log("accessToken", accessToken);
    console.log("expiresIn", expiresIn);

    const handlePasswordSet = (password: string) => {
        form.setValue("password", password);
        setStep(2);
    };

    if (passwordType === "PINS") {
        return (
            <PinPasswordPage
                onPasswordSet={handlePasswordSet}
                onSwitchToPassword={() => setPasswordType("PASSWORD")}
            />
        );
    }

    if (passwordType === "PASSWORD") {
        return (
            <PasswordInputPage
                onSwitchToPin={() => setPasswordType("PINS")}
                onLogout={() => navigate("/signout")}
            />
        );
    }

    return null;
};
