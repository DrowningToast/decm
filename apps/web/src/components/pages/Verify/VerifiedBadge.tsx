import { CheckCircle2, Loader2, XCircle } from "lucide-react";
import { useTranslation } from "react-i18next";

export function VerifiedBadge({ verified }: { verified: boolean | null }) {
    const { t } = useTranslation();
    if (verified === null) {
        return <Loader2 className="w-3 h-3 animate-spin text-muted-foreground shrink-0" />;
    }
    return verified ? (
        <span className="inline-flex items-center gap-0.5 text-[10px] font-medium text-green-600 bg-green-500/10 rounded px-1 py-0.5 shrink-0">
            <CheckCircle2 className="w-2.5 h-2.5" />
            {t("certificateVerify.verified")}
        </span>
    ) : (
        <span className="inline-flex items-center gap-0.5 text-[10px] font-medium text-destructive bg-destructive/10 rounded px-1 py-0.5 shrink-0">
            <XCircle className="w-2.5 h-2.5" />
            {t("certificateVerify.invalid")}
        </span>
    );
}
