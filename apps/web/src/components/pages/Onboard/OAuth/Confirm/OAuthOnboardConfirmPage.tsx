import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { BaseConfirmPage, type ConfirmationItem } from "../../BaseConfirmPage";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { OAuthOnboardContext } from "../OAuthOnboardContext";

export const OAuthOnboardConfirmPage: React.FC = () => {
    const { t } = useTranslation();
    const { setStep } = useContext(OnboardPageContext);
    const { handleSubmit: onSubmit } = useContext(OAuthOnboardContext)

    const confirmations: ConfirmationItem[] = [
        {
            id: "accept-risk",
            message: t("onboard.confirm.oauth.checkboxes.acceptRisk"),
        },
        {
            id: "acknowledge-password-loss",
            message: t("onboard.confirm.oauth.checkboxes.acknowledgePasswordLoss"),
        },
    ];

    const handleConfirm = () => {
        onSubmit();
    };

    const handleBack = () => {
        setStep(2); // Go back to profile page
    };

    return (
        <BaseConfirmPage
            title={t("onboard.confirm.title")}
            requiredConfirmations={confirmations}
            confirmButtonText={t("common.confirm")}
            backButtonText={t("common.back")}
            onConfirm={handleConfirm}
            onBack={handleBack}
            requireAllChecked={true}
        >
            <p className="mb-0" dangerouslySetInnerHTML={{
                __html: t("onboard.confirm.oauth.description.part1")
            }} />
            <p className="mb-0">&nbsp;</p>
            <p dangerouslySetInnerHTML={{
                __html: t("onboard.confirm.oauth.description.part2")
            }} />
        </BaseConfirmPage>
    );
};