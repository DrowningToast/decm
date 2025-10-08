import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface PasswordInputPageProps {
    onPasswordSet: (password: string) => void;
    onSwitchToPin: () => void;
    onLogout: () => void;
}

export const PasswordInputPage: React.FC<PasswordInputPageProps> = ({
    onPasswordSet,
    onSwitchToPin,
    onLogout,
}) => {
    const { t } = useTranslation();
    const [password, setPassword] = useState("");

    const handleConfirm = () => {
        if (password.length >= 6) {
            onPasswordSet(password);
        }
    };

    return (
        <div className="min-h-screen bg-[#e9dede] flex flex-col items-center px-6 py-16 md:py-24">
            <div className="w-full max-w-md space-y-16">
                {/* Header Section */}
                <div className="space-y-2">
                    <Typography
                        variant="header"
                        tag="h1"
                        className="text-[36px] leading-[40px] text-primary [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                    >
                        {t("onboard.title")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        className="text-base text-[#362927] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("onboard.subtitle")}
                    </Typography>
                </div>

                {/* Password Input Section */}
                <div className="space-y-6">
                    <div className="space-y-4">
                        {/* Password Input */}
                        <Input
                            type="password"
                            placeholder={t("onboard.enterPassword")}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="w-full h-12 bg-white border-2 border-[#362927]/20 rounded-xl text-[#362927] placeholder:text-[#362927]/50"
                        />

                        {/* Switch to PIN Link */}
                        <button
                            type="button"
                            onClick={onSwitchToPin}
                            className="w-full text-center"
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                className="text-xs italic underline text-[#362927] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors"
                            >
                                {t("onboard.preferPin")}
                            </Typography>
                        </button>

                        {/* Logout Link */}
                        <button
                            type="button"
                            onClick={onLogout}
                            className="w-full text-center"
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                className="text-xs italic underline text-[#362927] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors"
                            >
                                {t("onboard.logout")}
                            </Typography>
                        </button>
                    </div>

                    {/* Confirm Button */}
                    <Button
                        onClick={handleConfirm}
                        disabled={password.length < 6}
                        className="w-full h-12 bg-primary hover:bg-primary/90 text-[#e9dede] rounded-xl text-base font-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("common.confirm")}
                    </Button>
                </div>
            </div>
        </div>
    );
};

