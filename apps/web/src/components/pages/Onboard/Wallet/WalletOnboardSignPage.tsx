import { Button } from "@/components/ui/button";
import { useCallback, useContext } from "react";
import { WalletOnboardContext } from "./WalletOnboardContext";
import { useSignMessage, useWalletClient } from "wagmi";
import { ErrorPage } from "../../Error";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { LOCAL_STORAGE_KEYS } from "@/lib/constants/localStorage";
import { LogoutButton } from "@/components/auth/LogoutButton";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

export const WalletOnboardSignPage = () => {
    const { t } = useTranslation();
    const { data: walletClient } = useWalletClient();
    const { signMessageAsync } = useSignMessage();
    const { signMessage, isPending } = useContext(WalletOnboardContext);
    const [, setSignSignature] = useLocalStorage<string | undefined>(
        LOCAL_STORAGE_KEYS.AUTH_SIGN_SIGNATURE,
        undefined,
    );

    const handleSignSignature = useCallback(async () => {
        if (!walletClient || !signMessage) {
            return;
        }

        try {
            const signature = await signMessageAsync({ message: signMessage });
            setSignSignature(signature);
        } catch (error) {
            console.error(error);
            toast.error(t("verify.signMessageError"));
        }
    }, [walletClient, signMessage, signMessageAsync, setSignSignature, t]);

    if (isPending) {
        return (
            <div className="min-h-screen bg-[#e9dede] flex flex-col items-center justify-center px-6">
                <Typography variant="text" tag="p" color="foreground" className="animate-pulse">
                    {t("common.loading")}
                </Typography>
            </div>
        );
    }

    if (!walletClient) {
        return <ErrorPage />;
    }

    return (
        <div className="min-h-screen bg-[#e9dede] flex flex-col items-center px-6 py-16 md:py-24 relative overflow-hidden">
            {/* Background decorative image - positioned absolutely */}
            <div className="absolute inset-0 pointer-events-none opacity-30 overflow-hidden">
                <div className="absolute top-[53%] left-1/2 -translate-x-1/2 w-[90%] md:w-[45%] h-auto">
                    {/* Lady Justice placeholder - using background color for now */}
                    <div className="w-full aspect-[3/4] bg-gradient-to-b from-transparent via-muted/20 to-transparent rounded-full blur-3xl" />
                </div>
            </div>

            {/* Main content - positioned relative to stay above background */}
            <div className="relative z-10 w-full max-w-[420px] space-y-8 md:space-y-9">
                {/* Header Section */}
                <div className="space-y-1.5">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-[36px] leading-[40px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] tracking-[0.06px]"
                    >
                        {t("verify.title")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="background-alt"
                        className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px]"
                    >
                        {t("verify.subtitle")}
                    </Typography>
                </div>

                {/* Action Buttons Section */}
                <div className="space-y-3 pt-6 md:pt-8">
                    {/* Sign Message Button */}
                    <Button
                        type="button"
                        onClick={handleSignSignature}
                        variant="primary"
                        size="xl"
                        className="w-full"
                    >
                        {t("verify.requestButton")}
                    </Button>

                    <LogoutButton type="disconnect" />
                </div>
            </div>
        </div>
    );
};
