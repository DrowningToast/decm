import { useContext, useState } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Eye, EyeOff } from "lucide-react";
import { FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { OAuthOnboardContext } from "../OAuthOnboardContext";
import { OnboardPageContext } from "@/pages/onboard/[method]";
interface PasswordInputPageProps {
    onSwitchToPin: () => void;
    onLogout: () => void;
}

export const PasswordInputPage: React.FC<PasswordInputPageProps> = ({
    onSwitchToPin,
    onLogout,
}) => {
    const { form } = useContext(OAuthOnboardContext);
    const { setStep } = useContext(OnboardPageContext);
    const { t } = useTranslation();
    const [showPassword, setShowPassword] = useState(false);

    const isValid = !form.formState.isValidating && !form.getFieldState('password').invalid

    const handleConfirm = () => {
        if (isValid) {
            setStep(2)
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

                {/* Password Input Section */}
                <div className="space-y-6">
                    <div className="flex flex-col gap-y-3">
                        {/* Password Input */}
                        <div className="relative">
                            <FormField
                                control={form.control}
                                name="password"
                                render={({ field }) => (
                                    <>
                                        <FormItem>
                                            <FormControl>
                                                <Input
                                                    type={showPassword ? "text" : "password"}
                                                    placeholder={t("onboard.enterPassword")}
                                                    {...field}
                                                    className="w-full h-12 border-2 border-[#362927]/20 rounded-xl text-foreground placeholder:text-foreground-alt pr-12 bg-primary"
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    </>

                                )}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-[#362927]/50 hover:text-[#362927] transition-colors"
                            >
                                {showPassword ? (
                                    <Eye className="w-5 h-5" />
                                ) : (
                                    <EyeOff className="w-5 h-5" />
                                )}
                            </button>
                        </div>

                        {/* Switch to PIN Link */}
                        <button
                            type="button"
                            onClick={onSwitchToPin}
                            className="text-start h-[14.5px] inline-block"
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                color="background-alt"
                                className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors"
                            >
                                {t("onboard.preferPin")}
                            </Typography>
                        </button>

                        {/* Logout Link */}
                        <button
                            type="button"
                            onClick={onLogout}
                            className="text-start h-[14.5px] inline-block"
                        >
                            <Typography
                                variant="text"
                                tag="span"
                                color="background-alt"
                                className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors"
                            >
                                {t("onboard.logout")}
                            </Typography>
                        </button>
                    </div>

                    {/* Confirm Button */}
                    <Button
                        onClick={handleConfirm}
                        disabled={!isValid}
                        className="w-full h-12 bg-primary hover:bg-primary/90 text-[#e9dede] rounded-xl text-base font-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {t("common.confirm")}
                    </Button>
                </div>
            </div>
        </div>
    );
};

