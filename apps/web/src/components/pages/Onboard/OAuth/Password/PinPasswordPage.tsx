import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp";
import { Button } from "@/components/ui/button";
import { LogoutButton } from "../../../../auth/LogoutButton";

interface PinPasswordPageProps {
    onPasswordSet: (password: string) => void;
    onSwitchToPassword: () => void;
}

export const PinPasswordPage: React.FC<PinPasswordPageProps> = ({
    onPasswordSet,
    onSwitchToPassword,
}) => {
    const { t } = useTranslation();
    const [pin, setPin] = useState("");

    const handleConfirm = () => {
        if (pin.length === 6) {
            onPasswordSet(pin);
        }
    };

    return (
        <div className="min-h-screen bg-[#e9dede] flex flex-col items-center px-6 py-16 md:py-24">
            <div className="w-full max-w-[420px] space-y-16 text-background">
                {/* Header Section */}
                <div className="space-y-2">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-[36px] leading-[40px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                    >
                        {t("onboard.title")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="background-alt"
                        className="text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("onboard.subtitle")}
                    </Typography>
                </div>

                {/* PIN Input Section */}
                <div className="space-y-6">
                    <div className="flex flex-col gap-y-3">
                        {/* PIN Input */}
                        <div className="flex justify-center">
                            <InputOTP
                                maxLength={6}
                                value={pin}
                                onChange={(value) => {
                                    setPin(value);
                                }}
                                containerClassName="gap-4"
                            >
                                <InputOTPGroup className="gap-4">
                                    <InputOTPSlot
                                        index={0}
                                        className="w-10 h-16 md:w-12 md:h-16 rounded-xl bg-primary text-[#e9dede] border-0 text-2xl font-semibold shadow-none"
                                    />
                                    <InputOTPSlot
                                        index={1}
                                        className="w-10 h-16 md:w-12 md:h-16 rounded-xl bg-primary text-[#e9dede] border-0 text-2xl font-semibold shadow-none"
                                    />
                                    <InputOTPSlot
                                        index={2}
                                        className="w-10 h-16 md:w-12 md:h-16 rounded-xl bg-primary text-[#e9dede] border-0 text-2xl font-semibold shadow-none"
                                    />
                                    <InputOTPSlot
                                        index={3}
                                        className="w-10 h-16 md:w-12 md:h-16 rounded-xl bg-primary text-[#e9dede] border-0 text-2xl font-semibold shadow-none"
                                    />
                                    <InputOTPSlot
                                        index={4}
                                        className="w-10 h-16 md:w-12 md:h-16 rounded-xl bg-primary text-[#e9dede] border-0 text-2xl font-semibold shadow-none"
                                    />
                                    <InputOTPSlot
                                        index={5}
                                        className="w-10 h-16 md:w-12 md:h-16 rounded-xl bg-primary text-[#e9dede] border-0 text-2xl font-semibold shadow-none"
                                    />
                                </InputOTPGroup>
                            </InputOTP>
                        </div>

                        {/* Switch to Password Link */}
                        <button
                            type="button"
                            onClick={onSwitchToPassword}
                            className="text-start h-[14.5px] inline-block"
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                color="background-alt"
                                className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors"
                            >
                                {t("onboard.preferPassword")}
                            </Typography>
                        </button>

                        {/* Logout Link */}
                        <LogoutButton type="signout" />
                    </div>

                    {/* Confirm Button */}
                    <Button
                        onClick={handleConfirm}
                        disabled={pin.length !== 6}
                        className="w-full h-12 bg-primary hover:bg-primary/90 text-[#e9dede] rounded-xl text-base font-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("common.confirm")}
                    </Button>
                </div>
            </div>
        </div>
    );
};
