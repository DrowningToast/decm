import { BaseConfirmPage, type ConfirmationItem } from "../BaseConfirmPage";
import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { OnboardPageContext } from "@/pages/onboard/[method]";
import { WalletOnboardContext } from "./WalletOnboardContext";

export const WalletOnboardConfirmPage = () => {

    const { t } = useTranslation();
    const { setStep } = useContext(OnboardPageContext);
    const { handleSubmit: onSubmit } = useContext(WalletOnboardContext)

    const confirmations: ConfirmationItem[] = [
        {
            id: "accept-risk",
            message: t("onboard.confirm.wallet.checkboxes.acceptRisk"),
        },
        {
            id: "acknowledge-wallet-loss",
            message: t("onboard.confirm.wallet.checkboxes.acknowledgeWalletLoss"),
        },
    ];

    const handleConfirm = () => {
        onSubmit();
    };

    const handleBack = () => {
        setStep(1); // Go back to profile page
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
                __html: t("onboard.confirm.wallet.description.part1")
            }} />
            <p className="mb-0">&nbsp;</p>
            <p dangerouslySetInnerHTML={{
                __html: t("onboard.confirm.wallet.description.part2")
            }} />
        </BaseConfirmPage>
    )
}