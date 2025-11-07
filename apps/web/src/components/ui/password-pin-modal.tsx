import { useState, useRef, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ExternalLinkIcon, Eye, EyeOff } from "lucide-react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { useAuth } from "@/context/AuthContext";
import type { ProfileVerifyPasswordRequest } from "@decm/api";
import { toast } from "sonner";

interface PasswordPinModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSuccess?: (result: { type: "pin" | "password"; value: string }) => void;
    onError?: (error: unknown) => void;
    title?: string;
    description?: string;
    showModeToggle?: boolean;
    showSigningDetails?: boolean;
    signingDetails?: {
        contractAddress?: string;
        transactionType: string;
        details: string;
    };
}

const pinSchema = z.object({
    pin: z.string().length(6, "PIN must be exactly 6 digits"),
});

const passwordSchema = z.object({
    password: z.string().min(6, "Password must be at least 6 characters"),
});

export const PasswordPinModal: React.FC<PasswordPinModalProps> = ({
    isOpen,
    onClose,
    onSuccess,
    onError,
    title,
    description,
    showModeToggle = true,
    showSigningDetails = false,
    signingDetails,
}) => {
    const { t } = useTranslation();
    const [mode, setMode] = useState<"pin" | "password">("pin");
    const [showPassword, setShowPassword] = useState(false);
    const [pinValues, setPinValues] = useState<string[]>(["", "", "", "", "", ""]);
    const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

    const { user } = useAuth();

    const {
        mutateAsync: _verifyPassword,
        isPending: isVerifying,
        error: apiError,
    } = useMutation({
        mutationFn: async (data: ProfileVerifyPasswordRequest) => {
            return await coreApiClient.v1.verifyPassword(data);
        },
    });

    const pinForm = useForm<z.infer<typeof pinSchema>>({
        resolver: zodResolver(pinSchema),
        mode: "onChange",
    });

    const passwordForm = useForm<z.infer<typeof passwordSchema>>({
        resolver: zodResolver(passwordSchema),
        mode: "onChange",
    });

    // Reset forms when modal opens/closes
    useEffect(() => {
        if (isOpen) {
            setMode("pin");
            setShowPassword(false);
            setPinValues(["", "", "", "", "", ""]);
            pinForm.reset();
            passwordForm.reset();
        }
    }, [isOpen, pinForm, passwordForm]);

    // Auto-focus first PIN input when switching to PIN mode
    useEffect(() => {
        if (mode === "pin" && inputRefs.current[0]) {
            inputRefs.current[0]?.focus();
        }
    }, [mode]);

    const handlePinChange = (index: number, value: string) => {
        // Only allow digits
        const digitValue = value.replace(/[^0-9]/g, "").slice(0, 1);

        const newPinValues = [...pinValues];
        newPinValues[index] = digitValue;
        setPinValues(newPinValues);

        // Update form value
        const completePin = newPinValues.join("");
        pinForm.setValue("pin", completePin, { shouldValidate: true });

        // Auto-focus next input
        if (digitValue && index < 5 && inputRefs.current[index + 1]) {
            inputRefs.current[index + 1]?.focus();
        }

        // Auto-submit if all digits are filled
        if (completePin.length === 6) {
            handleConfirm({ type: "pin", value: completePin });
        }
    };

    const handlePinKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Backspace" && !pinValues[index] && index > 0) {
            // Focus previous input if current is empty
            inputRefs.current[index - 1]?.focus();
        } else if (e.key === "Backspace") {
            // Clear current input
            handlePinChange(index, "");
        }
    };

    const handleConfirm = async (data: { type: "pin" | "password"; value: string }) => {
        try {
            const apiResponse = await _verifyPassword({
                authentication_credential_id: user?.authentication_credential_id,
                password: data.value,
            });

            if (apiResponse.is_success) {
                // Success
                onSuccess?.(data);
                onClose();
            } else {
                toast.error(apiResponse.message);
            }
        } catch (error) {
            // Network or unexpected error
            toast.error(t("common.networkError"));
            onError?.(error);
        }
    };

    const handleCancel = () => {
        onClose();
    };

    const handleModeSwitch = (newMode: "pin" | "password") => {
        setMode(newMode);
        if (newMode === "pin") {
            setTimeout(() => inputRefs.current[0]?.focus(), 100);
        }
    };

    const isPinValid = pinValues.every((val) => val !== "");
    const isPasswordValid = passwordForm.formState.isValid;

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent
                className="bg-[#e9dede] max-w-[420px] border-none shadow-2xl"
                showCloseButton={false}
            >
                {/* Header Section */}
                <DialogHeader className="text-center ">
                    <DialogTitle className="text-2xl text-primary [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]">
                        {title || "Signature Request"}
                    </DialogTitle>
                    <DialogDescription className="text-base text-background-alt [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]">
                        {description || "Review request details before you confirm."}
                    </DialogDescription>
                </DialogHeader>

                {/* Signing Details */}
                {showSigningDetails && signingDetails && (
                    <div className="bg-[#f5f5f5] rounded-xl p-4 space-y-3 border border-[#362927]/10 my-2 mb-4">
                        <Typography
                            variant="header"
                            tag="h3"
                            color="foreground"
                            className="text-sm font-semibol text-black"
                        >
                            {t("signing.details.title")}
                        </Typography>

                        <div className="space-y-3">
                            {/* Contract Address (Optional) */}
                            {signingDetails.contractAddress && (
                                <div className="flex flex-col gap-1">
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        color="foreground-alt"
                                        className="text-xs text-black"
                                    >
                                        {t("signing.details.contractAddress")}
                                    </Typography>
                                    <div className="flex items-center gap-1">
                                        <a
                                            href={`https://www.etherscan.io/address/${signingDetails.contractAddress}`}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="underline"
                                        >
                                            <Typography
                                                variant="text"
                                                tag="p"
                                                color="foreground"
                                                className="text-xs font-mono break-all text-black underline"
                                            >
                                                {signingDetails.contractAddress}
                                            </Typography>
                                        </a>
                                        <ExternalLinkIcon size={12} color="#000" />
                                    </div>
                                </div>
                            )}

                            {/* Transaction Type */}
                            <div className="flex flex-col gap-1">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="foreground-alt"
                                    className="text-xs text-black"
                                >
                                    {t("signing.details.transactionType")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="foreground"
                                    className="text-xs text-black font-mono break-all"
                                >
                                    {signingDetails.transactionType}
                                </Typography>
                            </div>

                            {/* Details */}
                            <div className="flex flex-col gap-1">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    color="foreground-alt"
                                    className="text-xs text-black"
                                >
                                    {t("signing.details.details")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="foreground"
                                    className="text-xs text-black font-mono break-all"
                                >
                                    {signingDetails.details}
                                </Typography>
                            </div>
                        </div>
                    </div>
                )}

                {/* Input Section */}
                <form className="space-y-6">
                    <Form {...pinForm}>
                        {mode === "pin" ? (
                            /* PIN Input Mode */
                            <div className="flex justify-center gap-2">
                                {pinValues.map((value, index) => (
                                    <Input
                                        key={index}
                                        ref={(el) => {
                                            inputRefs.current[index] = el;
                                        }}
                                        type="text"
                                        inputMode="numeric"
                                        pattern="[0-9]*"
                                        maxLength={1}
                                        value={value}
                                        onChange={(e) => handlePinChange(index, e.target.value)}
                                        onKeyDown={(e) => handlePinKeyDown(index, e)}
                                        className="w-12 h-12 text-center text-xl border-2 border-[#362927]/20 rounded-xl text-foreground bg-primary focus:border-primary/50"
                                    />
                                ))}
                            </div>
                        ) : (
                            /* Password Input Mode */
                            <div className="relative">
                                <FormField
                                    control={passwordForm.control}
                                    name="password"
                                    render={({ field }) => (
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
                                    )}
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    className="absolute right-3 top-1/2 -translate-y-1/2 text-[#362927]/50 hover:text-[#362927] transition-colors"
                                >
                                    {showPassword ? (
                                        <EyeOff className="w-5 h-5" />
                                    ) : (
                                        <Eye className="w-5 h-5" />
                                    )}
                                </button>
                            </div>
                        )}

                        {/* API Error Display */}
                        {apiError && (
                            <div className="text-center">
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="destructive"
                                    className="text-sm"
                                >
                                    {apiError.message}
                                </Typography>
                            </div>
                        )}

                        {/* Mode Toggle */}
                        {showModeToggle && (
                            <div className="text-center">
                                <button
                                    type="button"
                                    onClick={() =>
                                        handleModeSwitch(mode === "pin" ? "password" : "pin")
                                    }
                                    className="text-center inline-block"
                                    disabled={isVerifying}
                                >
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        color="background-alt"
                                        className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        {mode === "pin"
                                            ? t("onboard.preferPassword")
                                            : t("onboard.preferPin")}
                                    </Typography>
                                </button>
                            </div>
                        )}
                    </Form>
                </form>

                {/* Action Buttons */}
                <DialogFooter className="flex-col gap-3">
                    <Button
                        variant="ghost"
                        onClick={handleCancel}
                        className=" h-12 text-foreground hover:text-background-alt rounded-xl text-base font-normal"
                    >
                        {t("common.cancel")}
                    </Button>

                    <Button
                        onClick={() => {
                            if (mode === "pin") {
                                const pinValue = pinValues.join("");
                                if (pinValue.length === 6) {
                                    handleConfirm({ type: "pin", value: pinValue });
                                }
                            } else {
                                const passwordValue = passwordForm.getValues("password");
                                if (passwordValue) {
                                    handleConfirm({ type: "password", value: passwordValue });
                                }
                            }
                        }}
                        disabled={isVerifying || (mode === "pin" ? !isPinValid : !isPasswordValid)}
                        className=" h-12 bg-primary hover:bg-primary/90 text-[#e9dede] rounded-xl text-base font-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {isVerifying ? (
                            <div className="flex items-center justify-center gap-2">
                                <div className="w-4 h-4 border-2 border-[#e9dede] border-t-transparent rounded-full animate-spin" />
                                <span>{t("common.verifying")}</span>
                            </div>
                        ) : (
                            t("common.confirm")
                        )}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};
