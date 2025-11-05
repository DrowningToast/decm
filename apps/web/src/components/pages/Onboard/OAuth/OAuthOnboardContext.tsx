import { createContext, useContext, useMemo } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslation } from "react-i18next";
import { Form } from "@/components/ui/form";
import { ErrorPage } from "../../Error";
import { toast } from "sonner";
import { OnboardPageContext } from "../../../../pages/onboard/[method]";
import { AccountAlreadyExistsPage } from "../AccountAlreadyExistPage";
import type { Profile } from "../ProfilePage";
import { useOAuthOnboardUsecase } from "../Usecase/useOAuthOnboardUsecase";

type OAuthOnboardContextType = {
    form: UseFormReturn<OAuthOnboardForm>;
    handleSubmit: () => Promise<void>;
};

//eslint-disable-next-line react-refresh/only-export-components
export const OAuthOnboardContext = createContext<OAuthOnboardContextType>(
    undefined as unknown as OAuthOnboardContextType,
);

const createOAuthOnboardFormSchema = (t: (key: string) => string) => {
    return z.object({
        password: z
            .string()
            .min(6, { message: t("validation.passwordMin6") })
            .max(32, { message: t("validation.passwordMax32") }),
    });
};

export type OAuthOnboardForm = z.infer<ReturnType<typeof createOAuthOnboardFormSchema>>;

type OAuthOnboardErrorType = "missingAccessToken" | "missingExpiresIn" | "alreadyOnboarded";

const OAuthOnboardProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
    const { t } = useTranslation();
    const { accessToken, expiresIn, onboardStatus } = useContext(OnboardPageContext);

    const OAuthOnboardFormSchema = createOAuthOnboardFormSchema(t);
    const oauthForm = useForm<OAuthOnboardForm>({
        resolver: zodResolver(OAuthOnboardFormSchema),
        mode: "onChange",
    });
    const { profileForm } = useContext(OnboardPageContext);
    const { usecaseAsync, isLoading } = useOAuthOnboardUsecase();
    const onSubmit = async (data: OAuthOnboardForm, profile: Profile) => {
        if (isLoading) {
            toast.error(t("errors.generic"));
            return;
        }
        if (!accessToken || !expiresIn) {
            return;
        }

        await usecaseAsync(accessToken, expiresIn, data.password, profile);
    };

    const handleSubmit = async () => {
        const data = oauthForm.getValues();
        const profile = profileForm.getValues();
        if (!accessToken || !expiresIn) {
            return;
        }
        console.log(`profile`, profile);
        if (!onboardStatus?.authentication_credential_id) {
            if (!data.password) {
                toast.error(t("flow.oauth_google.create_account_error_password"));
                return;
            }
        }

        await onSubmit(data, profile);
    };

    const errorType: OAuthOnboardErrorType | undefined = useMemo(() => {
        if (!accessToken) {
            return "missingAccessToken";
        }
        if (!expiresIn) {
            return "missingExpiresIn";
        }
        if (onboardStatus?.authentication_credential_id && onboardStatus?.profile_id) {
            return "alreadyOnboarded";
        }
        return undefined;
    }, [
        accessToken,
        expiresIn,
        onboardStatus?.authentication_credential_id,
        onboardStatus?.profile_id,
    ]);

    if (errorType) {
        switch (errorType) {
            case "alreadyOnboarded":
                return <AccountAlreadyExistsPage />;
            case "missingAccessToken":
                return (
                    <ErrorPage
                        title={t("flow.oauth_google.missing_access_token")}
                        description={t("flow.oauth_google.missing_access_token_description")}
                    />
                );
            case "missingExpiresIn":
                return (
                    <ErrorPage
                        title={t("flow.oauth_google.missing_expires_in")}
                        description={t("flow.oauth_google.missing_expires_in_description")}
                    />
                );
        }
    }

    return (
        <OAuthOnboardContext.Provider value={{ form: oauthForm, handleSubmit }}>
            <Form {...oauthForm}>
                <form>{children}</form>
            </Form>
        </OAuthOnboardContext.Provider>
    );
};

export { OAuthOnboardProvider };
