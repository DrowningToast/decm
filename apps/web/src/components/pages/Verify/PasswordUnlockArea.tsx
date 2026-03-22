import { useTranslation } from "react-i18next";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Typography } from "@/components/typography/typography";

export interface PasswordUnlockAreaProps {
    password: string;
    onPasswordChange: (value: string) => void;
    passwordError: string;
    onUnlock: () => void;
    onCancel: () => void;
    isUnlocking: boolean;
}

export function PasswordUnlockArea({
    password,
    onPasswordChange,
    passwordError,
    onUnlock,
    onCancel,
    isUnlocking,
}: PasswordUnlockAreaProps) {
    const { t } = useTranslation();

    return (
        <div className="flex flex-col gap-3">
            <Typography variant="text" tag="p" color="muted" className="text-sm text-center">
                {t("certificateVerify.passwordProtected")}
            </Typography>
            <div className="flex items-stretch gap-2">
                <Input
                    type="password"
                    placeholder={t("certificateVerify.passwordPlaceholder")}
                    value={password}
                    onChange={(e) => onPasswordChange(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && onUnlock()}
                    aria-label={t("certificateVerify.passwordPlaceholder")}
                    aria-invalid={!!passwordError}
                    aria-describedby={passwordError ? "password-error" : undefined}
                    className="flex-1 h-auto outline-primary"
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
                    id="password-error"
                    variant="text"
                    tag="p"
                    color="destructive"
                    className="text-xs"
                    role="alert"
                >
                    {passwordError}
                </Typography>
            )}
            <button
                type="button"
                onClick={onCancel}
                className="cursor-pointer text-xs text-muted-foreground hover:text-foreground transition-colors text-center underline"
            >
                {t("certificateVerify.cancelButton")}
            </button>
        </div>
    );
}
