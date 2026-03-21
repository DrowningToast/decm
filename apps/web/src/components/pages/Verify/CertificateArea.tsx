import { Award, LockKeyhole, Loader2 } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import type {
    GetCertificateShareDataResult,
    CertificateShareViewStatus,
} from "@/services/CertificateService/mapper";

function AspectFrame({ children }: { children: React.ReactNode }) {
    return (
        <div className="w-full relative" style={{ paddingBottom: "75%" }}>
            <div className="absolute inset-0 rounded-xl overflow-hidden">{children}</div>
        </div>
    );
}

function EmptyCertificateFrame() {
    const { t } = useTranslation();
    return (
        <AspectFrame>
            <div className="absolute inset-0 rounded-xl border-2 border-dashed border-muted-foreground/20 bg-muted/10 flex flex-col items-center justify-center gap-3">
                <Award className="w-12 h-12 text-muted-foreground/30" />
                <Typography variant="text" tag="p" color="muted" className="text-sm">
                    {t("certificateVerify.emptyFrame")}
                </Typography>
            </div>
        </AspectFrame>
    );
}

function LockedCertificateFrame() {
    const { t } = useTranslation();
    return (
        <AspectFrame>
            <div className="absolute inset-0 rounded-xl border-2 border-dashed border-amber-400/40 bg-amber-50/10 dark:bg-amber-900/10 flex flex-col items-center justify-center gap-3">
                <div className="rounded-full bg-amber-100/60 dark:bg-amber-900/30 p-4">
                    <LockKeyhole className="w-8 h-8 text-amber-500" />
                </div>
                <div className="flex flex-col items-center gap-1">
                    <Typography
                        variant="text"
                        tag="p"
                        color="foreground"
                        className="text-sm font-medium"
                    >
                        {t("certificateVerify.lockedFrameTitle")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-xs text-center px-6"
                    >
                        {t("certificateVerify.lockedFrameHint")}
                    </Typography>
                </div>
            </div>
        </AspectFrame>
    );
}

export interface CertificateAreaProps {
    handle: string | null;
    shareStatus: CertificateShareViewStatus | null;
    shareData: GetCertificateShareDataResult | null;
    isLoading: boolean;
    isError: boolean;
    errorStatus: number | undefined;
    imageUrl: string | undefined;
    isImageLoading: boolean;
}

export function CertificateArea({
    handle,
    shareStatus,
    shareData,
    isLoading,
    isError,
    errorStatus,
    imageUrl,
    isImageLoading,
}: CertificateAreaProps) {
    const { t } = useTranslation();

    if (!handle) {
        return <EmptyCertificateFrame />;
    }

    if (isLoading) {
        return (
            <AspectFrame>
                <div className="absolute inset-0 rounded-xl border border-muted/30 bg-muted/10 flex items-center justify-center">
                    <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
                </div>
            </AspectFrame>
        );
    }

    if (isError) {
        const isRateLimited = errorStatus === 429;
        return (
            <AspectFrame>
                <div className="absolute inset-0 rounded-xl border-2 border-dashed border-destructive/20 bg-destructive/5 flex flex-col items-center justify-center gap-2">
                    <Typography variant="text" tag="p" color="destructive" className="text-sm">
                        {isRateLimited
                            ? t("certificateVerify.rateLimited")
                            : t("certificateVerify.notFound")}
                    </Typography>
                </div>
            </AspectFrame>
        );
    }

    if (shareStatus === "PASSWORD_LOCKED" && !shareData) {
        return <LockedCertificateFrame />;
    }

    if (shareStatus === "VALID_BUT_PENDING") {
        return (
            <AspectFrame>
                <div className="absolute inset-0 rounded-xl border border-muted/30 bg-muted/10 flex flex-col items-center justify-center gap-2">
                    <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
                    <Typography variant="text" tag="p" color="muted" className="text-sm">
                        {t("certificateVerify.pending")}
                    </Typography>
                </div>
            </AspectFrame>
        );
    }

    if (shareStatus === "READY" && shareData) {
        return (
            <AspectFrame>
                {isImageLoading && (
                    <div className="absolute inset-0 rounded-xl border border-muted/30 bg-muted/10 flex items-center justify-center">
                        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
                    </div>
                )}
                {imageUrl && (
                    <img
                        src={imageUrl}
                        alt={shareData.payload.data.certificateTitle ?? "Certificate"}
                        className="absolute inset-0 w-full h-full object-contain rounded-xl"
                    />
                )}
                {!isImageLoading && !imageUrl && (
                    <div className="absolute inset-0 rounded-xl border border-primary/30 bg-primary/5 flex items-center justify-center">
                        <Typography
                            variant="header"
                            tag="h2"
                            color="foreground"
                            className="text-xl font-header"
                        >
                            {shareData.payload.data.certificateTitle}
                        </Typography>
                    </div>
                )}
            </AspectFrame>
        );
    }

    return <EmptyCertificateFrame />;
}
