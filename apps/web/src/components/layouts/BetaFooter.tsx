import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";

export const BetaFooter = () => {
    const { t } = useTranslation();

    return (
        <footer className="fixed bottom-0 left-0 right-0 z-30 pointer-events-none">
            <div className="flex justify-center items-center p-2">
                <Typography
                    variant="text"
                    tag="span"
                    color="foreground-alt"
                    className="text-[9px] md:text-xs opacity-70"
                >
                    {t("nav.betaAccess")}
                </Typography>
            </div>
        </footer>
    );
};
