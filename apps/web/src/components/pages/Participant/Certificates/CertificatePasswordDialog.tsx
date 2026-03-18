import { useState, useEffect } from "react";
import { LockKeyhole } from "lucide-react";
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogCancel,
} from "@/components/ui/alert-dialog";
import { Button, buttonVariants } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { useTranslation } from "react-i18next";
import { useCertificateUpdateSharePasswordUsecase } from "./useCertificateUpdateSharePasswordUsecase";

interface CertificatePasswordDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    shareId: string | null;
    hasPassword: boolean;
}

export const CertificatePasswordDialog = ({
    open,
    onOpenChange,
    shareId,
    hasPassword,
}: CertificatePasswordDialogProps) => {
    const { t } = useTranslation();
    const { updateSharePassword, isPending } = useCertificateUpdateSharePasswordUsecase();

    const [isPasswordProtected, setIsPasswordProtected] = useState(false);
    const [password, setPassword] = useState("");

    useEffect(() => {
        if (open) {
            setIsPasswordProtected(hasPassword);
            setPassword("");
        } else {
            setIsPasswordProtected(false);
            setPassword("");
        }
    }, [open, hasPassword]);

    const handleCheckedChange = (checked: boolean) => {
        setIsPasswordProtected(checked);
        if (!checked) {
            setPassword("");
        }
    };

    const handleSave = async () => {
        if (!shareId) return;
        try {
            await updateSharePassword(shareId, isPasswordProtected, password);
            onOpenChange(false);
        } catch {
            // keep dialog open on error
        }
    };

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>
                        {t(
                            "participant.certificates.passwordDialog.title",
                            "Certificate Password Settings",
                        )}
                    </AlertDialogTitle>
                    <AlertDialogDescription>
                        {t(
                            "participant.certificates.passwordDialog.description",
                            "Enable password protection to restrict access to this certificate share.",
                        )}
                    </AlertDialogDescription>
                </AlertDialogHeader>

                <div className="flex flex-col gap-4 py-2">
                    {hasPassword && (
                        <div className="flex items-start gap-2.5 rounded-lg border border-amber-500/40 bg-amber-500/10 px-3 py-2.5">
                            <LockKeyhole className="mt-0.5 h-4 w-4 shrink-0 text-amber-500" />
                            <p className="text-sm text-amber-600 dark:text-amber-400">
                                {t(
                                    "participant.certificates.passwordDialog.passwordCurrentlySet",
                                    "A password is currently set. Enter a new password to replace it, or uncheck the box to remove protection.",
                                )}
                            </p>
                        </div>
                    )}

                    <div className="flex items-center gap-3">
                        <Checkbox
                            id="password-protection"
                            checked={isPasswordProtected}
                            onCheckedChange={handleCheckedChange}
                        />
                        <label htmlFor="password-protection" className="text-sm cursor-pointer">
                            {t(
                                "participant.certificates.passwordDialog.enablePasswordProtection",
                                "Enable password protection",
                            )}
                        </label>
                    </div>

                    <Input
                        type="password"
                        placeholder={
                            hasPassword
                                ? t(
                                      "participant.certificates.passwordDialog.newPasswordPlaceholder",
                                      "Enter new password",
                                  )
                                : t(
                                      "participant.certificates.passwordDialog.passwordPlaceholder",
                                      "Enter password",
                                  )
                        }
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        disabled={!isPasswordProtected}
                    />
                </div>

                <AlertDialogFooter>
                    <AlertDialogCancel
                        className={buttonVariants({ variant: "secondary-light" })}
                        onClick={() => onOpenChange(false)}
                    >
                        {t("common.cancel", "Cancel")}
                    </AlertDialogCancel>
                    <Button
                        variant="primary"
                        onClick={handleSave}
                        disabled={
                            isPending || !shareId || (isPasswordProtected && password.length < 1)
                        }
                    >
                        {t("common.save", "Save")}
                    </Button>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
};
