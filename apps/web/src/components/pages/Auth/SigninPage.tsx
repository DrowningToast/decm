import type React from "react";
import { useTranslation } from "react-i18next";
import { LanguageSwitcher } from "@/components/LanguageSwitcher";
import { Button } from "@/components/ui/button";
import { Link } from "@/router";
import { WalletConnectButton } from "@/components/WalletConnectButton";

interface SigninPageProps {
    onGoogleOAuthClick?: () => void;
}

export const SigninPage: React.FC<SigninPageProps> = ({ onGoogleOAuthClick }) => {
    const { t } = useTranslation();

    return (
        <div className="min-h-screen flex flex-col">
            {/* Header with Language Switcher */}
            <div className="flex justify-end p-4">
                <LanguageSwitcher />
            </div>

            {/* Main Content */}
            <div className="flex-1 flex items-center justify-center">
                <div className="flex flex-col gap-y-6 max-w-md w-full px-6">
                    <div className="text-center space-y-2">
                        <h1 className="text-4xl font-bold">{t('signin.title')}</h1>
                        <p className="text-muted-foreground">
                            {t('signin.subtitle')}
                        </p>
                    </div>

                    <div className="space-y-4">
                        <WalletConnectButton >
                            <Button className="w-full" size="lg" variant="primary">
                                {t('signup.walletButton')}
                            </Button>
                        </WalletConnectButton>

                        <div className="relative">
                            <div className="absolute inset-0 flex items-center">
                                <span className="w-full border-t" />
                            </div>
                            <div className="relative flex justify-center text-xs uppercase">
                                <span className="bg-background px-2 text-muted-foreground">
                                    {t('signup.divider')}
                                </span>
                            </div>
                        </div>

                        <Button variant="secondary-dark" className="w-full" size="lg" onClick={onGoogleOAuthClick}>
                            {t('signup.googleButton')}
                        </Button>
                    </div>

                    <p className="text-xs text-center text-muted-foreground">
                        <Link to="/signup" className="hover:underline">
                            {t('auth.noAccount')}
                        </Link>
                    </p>

                    <p className="text-xs text-center text-muted-foreground mt-4">
                        {t('signup.terms')}
                    </p>
                </div>
            </div>
        </div>
    )
}
