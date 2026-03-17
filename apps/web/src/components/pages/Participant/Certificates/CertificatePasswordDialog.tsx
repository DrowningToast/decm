import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogCancel,
    AlertDialogAction,
} from "@/components/ui/alert-dialog";
import { buttonVariants } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

interface CertificatePasswordDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export const CertificatePasswordDialog = ({
    open,
    onOpenChange,
}: CertificatePasswordDialogProps) => {
    const { t } = useTranslation();

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>
                        {t(
                            "participant.certificates.passwordDialog.title",
                            "Password placeholder title",
                        )}
                    </AlertDialogTitle>
                    <AlertDialogDescription>
                        {t(
                            "participant.certificates.passwordDialog.description",
                            "Password placeholder description",
                        )}
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel
                        className={buttonVariants({ variant: "secondary-light" })}
                        onClick={() => onOpenChange(false)}
                    >
                        {t("common.cancel", "Cancel")}
                    </AlertDialogCancel>
                    <AlertDialogAction
                        className={buttonVariants({ variant: "primary" })}
                        onClick={() => {}}
                    >
                        {t("common.submit", "Submit")}
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
};
