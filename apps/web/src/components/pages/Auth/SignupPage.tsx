import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import type React from "react";
import { Link } from "@/router";
import { WalletConnectButton } from "@/components/WalletConnectButton";
import { Separator } from "@/components/ui/separator";

interface SignupPageProps {
    onGoogleOAuthClick?: () => void;
}

export const SignupPage: React.FC<SignupPageProps> = ({ onGoogleOAuthClick }) => {
    const { t } = useTranslation();

    return (
        <div className="min-h-screen flex flex-col bg-foreground-alt text-background-alt">
            {/* Main Content */}
            <div className="flex-1 flex justify-center">
                <div className="flex flex-col gap-y-24 max-w-md w-full px-6">
                    <div className="text-center space-y-2">
                        <h1 className="text-4xl font-bold">{t('signup.title')}</h1>
                        <p className="text-muted-foreground">
                            {t('signup.subtitle')}
                        </p>
                    </div>

                    <div className="space-y-2.5">
                        <WalletConnectButton >
                            <Button className="w-full" size="lg" variant="primary" >
                                {t('signup.walletButton')}
                            </Button>
                        </WalletConnectButton>

                        <Separator className="w-full bg-foreground-alt text-foreground-alt" />

                        <Button variant="secondary-dark" className="w-full" size="lg" onClick={onGoogleOAuthClick}>
                            {t('signup.googleButton')}
                        </Button>
                        <p className="text-xs text-muted-foreground">
                            <Link to="/signin" className="italic underline text-left">
                                {t('auth.hasAccount')}
                            </Link>
                        </p>

                    </div>

                    <p className="text-xs text-center text-muted-foreground mt-4">
                        {t('signup.terms')}
                    </p>
                </div>
            </div>
        </div>
    )
}