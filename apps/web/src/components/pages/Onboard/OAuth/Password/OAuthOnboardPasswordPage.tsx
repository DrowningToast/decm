import { useContext, useState } from "react";
import { OAuthOnboardContext } from "../OAuthOnboardProvider";
import { PinPasswordPage } from "./PinPasswordPage";
import { PasswordInputPage } from "./PasswordInputPage";
import type { OAuthOnboardPasswordType } from "../constants";

export const OAuthOnboardPasswordPage = () => {
    const { form } = useContext(OAuthOnboardContext);
    const [passwordType, setPasswordType] = useState<OAuthOnboardPasswordType>("PINS");

    const handlePasswordSet = (password: string) => {
        form.setValue("password", password);

        // TODO: Proceed to next step
    };

    const handleLogout = () => {
        // TODO: Implement logout logic
        console.log("Logout clicked");
    };

    if (passwordType === "PINS") {
        return (
            <PinPasswordPage
                onPasswordSet={handlePasswordSet}
                onSwitchToPassword={() => setPasswordType("PASSWORD")}
                onLogout={handleLogout}
            />
        );
    }

    if (passwordType === "PASSWORD") {
        return (
            <PasswordInputPage
                onPasswordSet={handlePasswordSet}
                onSwitchToPin={() => setPasswordType("PINS")}
                onLogout={handleLogout}
            />
        );
    }

    return null;
};
