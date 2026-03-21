import { useTranslation } from "react-i18next";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Typography } from "@/components/typography/typography";

export interface PasswordUnlockAreaProps {
    password: string;
    onPasswordChange: (value: string) => void;
    passwordError: string;
    onUnlock: () => void;
    isUnlocking: boolean;
}

export function PasswordUnlockArea({
    password,
    onPasswordChange,
    passwordError,
    onUnlock,
    isUnlocking,
}: PasswordUnlockAreaProps) {
    const { t } = useTranslation();

    return (
        <div className="flex flex-col gap-3">
            <Typography variant="text" tag="p" color="muted" className="text-sm text-center">
                {t("certificateVerify.passwordProtected")}
            </Typography>
            <div className="flex gap-2">
                <Input
                    type="password"
                    placeholder={t("certificateVerify.passwordPlaceholder")}
                    value={password}
                    onChange={(e) => onPasswordChange(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && onUnlock()}
                    aria-label={t("certificateVerify.passwordPlaceholder")}
                    className="flex-1"
                />
                <Button
                    variant="primary"
                    onClick={onUnlock}
                    disabled={isUnlocking}
                    loading={isUnlocking}
                >
                    {t("certificateVerify.unlockButton")}
                </Button>
            </div>
            {passwordError && (
                <Typography
                    variant="text"
                    tag="p"
                    color="destructive"
                    className="text-xs"
                    role="alert"
                >
                    {passwordError}
                </Typography>
            )}
        </div>
    );
}
