import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import type React from "react";
import { Link } from "@/router";
import { WalletConnectButton } from "@/components/WalletConnectButton";
import { Typography } from "@/components/typography/typography";
import { ThemisWithScale } from "@/components/assets/ThemisWithScale";

interface SignupPageProps {
    onGoogleOAuthClick?: () => void;
    isLoading?: boolean;
}

export const SignupPage: React.FC<SignupPageProps> = ({
    onGoogleOAuthClick,
    isLoading = false,
}) => {
    const { t } = useTranslation();

    return (
        <div className="relative min-h-screen flex flex-col overflow-hidden">
            {/* Main Content */}
            <div className="relative z-10 flex-1 flex justify-center px-6">
                <div className="w-full flex flex-col gap-y-24 md:gap-y-9 md:w-[420px]">
                    {/* Title Section */}
                    <div className="flex flex-col gap-1.5">
                        <Typography
                            variant="header"
                            tag="h1"
                            color="primary"
                            size="header"
                            className="leading-[40px] font-bold [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                        >
                            {t("signup.title")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="background-alt"
                            className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] max-w-[294px]"
                        >
                            {t("signup.subtitle")}
                        </Typography>
                    </div>

                    {/* Buttons Section */}
                    <div className="flex flex-col gap-2.5">
                        {/* Web3 Wallet Button */}
                        <WalletConnectButton
                            className="w-full h-12 bg-primary hover:bg-primary/90 rounded-[12px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            size="lg"
                            isLoading={isLoading}
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                color="foreground-alt"
                                className="text-base"
                            >
                                {t("signup.walletButton")}
                            </Typography>
                        </WalletConnectButton>

                        {/* Divider Line */}
                        <div className="w-full h-0 border-t border-[#b8b8b8]" />

                        {/* Google Button */}
                        <Button
                            variant="secondary-dark"
                            className="w-full h-12 bg-background-alt hover:bg-background-alt/90 rounded-[12px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            size="lg"
                            onClick={onGoogleOAuthClick}
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                color="foreground-alt"
                                className="text-base"
                            >
                                {t("signup.googleButton")}
                            </Typography>
                        </Button>

                        {/* Already have account link */}
                        <Link to="/signin" className="mt-0.5">
                            <Typography
                                variant="text"
                                tag="p"
                                color="background"
                                className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] decoration-solid [text-decoration-skip-ink:none] [text-underline-position:from-font]"
                            >
                                {t("auth.hasAccount")}
                            </Typography>
                        </Link>
                    </div>
                </div>
            </div>

            {/* Background Image - Themis Statue */}
            <div className="md:w-[674px] w-[428px] absolute left-1/2 -translate-x-2/3 md:-translate-x-1/2 bottom-0 translate-y-1/2">
                <ThemisWithScale className="w-full h-full" />
            </div>
        </div>
    );
};
