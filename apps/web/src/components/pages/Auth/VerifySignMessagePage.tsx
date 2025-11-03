import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import type React from "react";
import { Typography } from "@/components/typography/typography";
import { ThemisWithScale } from "@/components/assets/ThemisWithScale";
import { Link, useNavigate } from "@/router";
import { useGetSignMessage } from "../Onboard/useGetSignMessage";
import { useMemo } from "react";
import { useSignMessage } from "wagmi";
import { toast } from "sonner";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { onboardService } from "@/services/OnboardService";
import { OnboardRegistrationMethod } from "@decm/api";
import { handleUniversalError } from "@/common/Err";

import { OnboardMethods } from "@/pages/onboard/[method]";
interface VerifySignMessagePageProps {
    onDisconnectWallet?: () => void;
    isLoading?: boolean;
}

export const VerifySignMessagePage: React.FC<VerifySignMessagePageProps> = ({
    onDisconnectWallet,
    isLoading = false,
}) => {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [, setSignSignature] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
        undefined,
    );

    const { signMessageAsync } = useSignMessage();

    const { signMessage, isPending: isGetMessagePending } = useGetSignMessage();
    const isEnabled = useMemo(() => {
        return !isGetMessagePending && signMessage && !isLoading;
    }, [isGetMessagePending, isLoading, signMessage]);

    const handleRequestSigningMessage = async () => {
        if (!isEnabled || !signMessage || isLoading) {
            return;
        }

        const loadingToastId = toast.loading(t("verify.signMessageLoading"));
        try {
            const signature = await signMessageAsync({ message: signMessage });
            const response = await onboardService.checkOnboardStatus({
                method: OnboardRegistrationMethod.RegistrationMethodWallet,
                signSignature: signature,
            });
            if (response?.profile_id) {
                await navigate("/app");
            }
            if (response?.authentication_credential_id) {
                await navigate("/onboard/:method", {
                    params: {
                        method: OnboardMethods.WALLET,
                    },
                });
            }
            setSignSignature(signature);
            toast.dismiss(loadingToastId);
        } catch (error) {
            toast.dismiss(loadingToastId);
            console.error(error);
            if (error instanceof Error) {
                handleUniversalError(t, error);
            }
        }
    };

    return (
        <div className="relative min-h-screen flex flex-col text-foreground overflow-hidden">
            {/* Main Content */}
            <div className="relative z-10 flex-1 flex justify-center px-6">
                <div className="w-full flex flex-col gap-y-24 md:gap-y-9 md:w-[420px]">
                    {/* Title Section */}
                    <div className="flex flex-col gap-1.5">
                        <Typography
                            variant="header"
                            tag="h1"
                            color="primary"
                            className="font-['Cormorant_Garamond'] text-[36px] leading-[40px] font-bold [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                        >
                            {t("verify.title")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="foreground-alt"
                            className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] max-w-[294px]"
                        >
                            {t("verify.subtitle")}
                        </Typography>
                    </div>

                    {/* Button Section */}
                    <div className="flex flex-col gap-2.5">
                        {/* Request Signing Message Button */}
                        <Button
                            className="w-full h-12 bg-primary hover:bg-primary/90 rounded-[12px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            size="lg"
                            onClick={handleRequestSigningMessage}
                            disabled={!isEnabled}
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                color="foreground-alt"
                                className="text-base"
                            >
                                {t("verify.requestButton")}
                            </Typography>
                        </Button>

                        {/* Disconnect wallet link */}
                        <Link
                            to="/signout"
                            onClick={onDisconnectWallet}
                            className="mt-0.5 text-left"
                        >
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground-alt"
                                className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] decoration-solid [text-decoration-skip-ink:none] [text-underline-position:from-font]"
                            >
                                {t("verify.disconnectLink")}
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
