import { useRef } from "react";
import { useTranslation } from "react-i18next";
import { QRCodeCanvas } from "qrcode.react";
import { Download, ExternalLink, Copy, Eye, EyeOff, Lock, LockOpen } from "lucide-react";
import { Link } from "react-router-dom";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
    useCertificateDetailNavStore,
    useCertificateDetailsSharedNavStore,
} from "@/components/BottomNav/stores/certificates";

export const CertificateShareModal = () => {
    const { t } = useTranslation();
    const { isShareModalOpen, setIsShareModalOpen } = useCertificateDetailNavStore();
    const { shareableHandle, isPublished, isPasswordProtected } =
        useCertificateDetailsSharedNavStore();

    const qrRef = useRef<HTMLCanvasElement>(null);

    const shareUrl = shareableHandle
        ? `${window.location.origin}/verify?handle=${shareableHandle}`
        : null;

    const verifyPageUrl = shareableHandle ? `/verify?handle=${shareableHandle}` : null;

    const handleCopyUrl = async () => {
        if (!shareUrl) return;
        await navigator.clipboard.writeText(shareUrl);
    };

    const handleCopyCode = async () => {
        if (!shareableHandle) return;
        await navigator.clipboard.writeText(shareableHandle);
    };

    const handleCopyQr = async () => {
        const canvas = qrRef.current;
        if (!canvas) return;
        canvas.toBlob(async (blob) => {
            if (!blob) return;
            try {
                await navigator.clipboard.write([new ClipboardItem({ "image/png": blob })]);
            } catch {
                // fallback: copy URL instead
                if (shareUrl) await navigator.clipboard.writeText(shareUrl);
            }
        });
    };

    const handleDownloadQr = () => {
        const canvas = qrRef.current;
        if (!canvas) return;
        const url = canvas.toDataURL("image/png");
        const link = document.createElement("a");
        link.href = url;
        link.download = `certificate-qr-${shareableHandle ?? "code"}.png`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <Dialog open={isShareModalOpen} onOpenChange={setIsShareModalOpen}>
            <DialogContent className="w-full max-w-sm sm:max-w-md overflow-y-auto max-h-[90vh] overflow-x-hidden">
                {/* Header */}
                <DialogHeader>
                    <DialogTitle className="text-primary">
                        {t("participant.certificates.detail.shareModal.title", "Share Certificate")}
                    </DialogTitle>
                    <DialogDescription>
                        {t(
                            "participant.certificates.detail.shareModal.description",
                            "Share this certificate with anyone using the link or QR code below.",
                        )}
                    </DialogDescription>
                </DialogHeader>

                {shareableHandle && shareUrl ? (
                    <div className="flex flex-col gap-4 min-w-0">
                        {/* Status badges */}
                        <div className="flex flex-wrap gap-2">
                            {isPublished ? (
                                <span
                                    data-testid="badge-published"
                                    className="inline-flex items-center gap-1 text-xs font-medium px-2.5 py-1 rounded-full bg-green-500/15 text-green-600 border border-green-500/30"
                                >
                                    <Eye className="w-3 h-3" />
                                    {t(
                                        "participant.certificates.detail.shareModal.statusVisible",
                                        "Published",
                                    )}
                                </span>
                            ) : (
                                <span
                                    data-testid="badge-hidden"
                                    className="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 rounded-lg bg-yellow-500/20 text-yellow-700 border-2 border-yellow-500/60 w-full"
                                >
                                    <EyeOff className="w-3.5 h-3.5 shrink-0" />
                                    {t(
                                        "participant.certificates.detail.shareModal.statusHidden",
                                        "Hidden — recipients cannot view this yet",
                                    )}
                                </span>
                            )}

                            {isPasswordProtected ? (
                                <span
                                    data-testid="badge-password-protected"
                                    className="inline-flex items-center gap-1 text-xs font-medium px-2.5 py-1 rounded-full bg-blue-500/15 text-blue-600 border border-blue-500/30"
                                >
                                    <Lock className="w-3 h-3" />
                                    {t(
                                        "participant.certificates.detail.shareModal.passwordProtected",
                                        "Password protected",
                                    )}
                                </span>
                            ) : (
                                <span
                                    data-testid="badge-no-password"
                                    className="inline-flex items-center gap-1 text-xs font-medium px-2.5 py-1 rounded-full bg-muted/50 text-muted-foreground border border-border"
                                >
                                    <LockOpen className="w-3 h-3" />
                                    {t(
                                        "participant.certificates.detail.shareModal.noPassword",
                                        "No password",
                                    )}
                                </span>
                            )}
                        </div>

                        {/* QR Code (main focus) */}
                        <div className="flex flex-col items-center gap-3">
                            <div className="p-4 bg-white rounded-xl border border-border shadow-sm w-full flex justify-center">
                                <QRCodeCanvas
                                    ref={qrRef}
                                    value={shareUrl}
                                    size={200}
                                    marginSize={1}
                                    style={{ maxWidth: "100%", height: "auto" }}
                                />
                            </div>
                            <div className="flex gap-2">
                                <Button
                                    variant="secondary-light"
                                    size="sm"
                                    onClick={handleCopyQr}
                                    aria-label={t(
                                        "participant.certificates.detail.shareModal.copyQr",
                                        "Copy QR",
                                    )}
                                    className="gap-1.5"
                                >
                                    <Copy className="w-3.5 h-3.5" />
                                    {t(
                                        "participant.certificates.detail.shareModal.copyQr",
                                        "Copy QR",
                                    )}
                                </Button>
                                <Button
                                    variant="secondary-light"
                                    size="sm"
                                    onClick={handleDownloadQr}
                                    aria-label={t(
                                        "participant.certificates.detail.shareModal.downloadQr",
                                        "Download QR",
                                    )}
                                    className="gap-1.5"
                                >
                                    <Download className="w-3.5 h-3.5" />
                                    {t(
                                        "participant.certificates.detail.shareModal.downloadQr",
                                        "Download QR",
                                    )}
                                </Button>
                            </div>
                        </div>

                        {/* Share URL */}
                        <div className="flex flex-col gap-1.5">
                            <span className="text-xs text-primary font-semibold uppercase tracking-wide">
                                {t(
                                    "participant.certificates.detail.shareModal.shareUrl",
                                    "Share URL",
                                )}
                            </span>
                            <div className="flex items-center gap-2 p-2.5 rounded-lg bg-muted/40 border border-border min-w-0">
                                <span
                                    data-testid="share-url"
                                    className="flex-1 min-w-0 text-sm font-mono break-all leading-snug select-all"
                                >
                                    {shareUrl}
                                </span>
                                <button
                                    data-testid="copy-url-btn"
                                    onClick={handleCopyUrl}
                                    aria-label={t(
                                        "participant.certificates.detail.shareModal.copyUrl",
                                        "Copy URL",
                                    )}
                                    className="shrink-0 p-1.5 rounded hover:bg-foreground/10 transition-colors"
                                >
                                    <Copy className="w-3.5 h-3.5 text-muted-foreground" />
                                </button>
                            </div>
                        </div>

                        {/* Share code */}
                        <div className="flex flex-col gap-1.5">
                            <span className="text-xs text-primary font-semibold uppercase tracking-wide">
                                {t(
                                    "participant.certificates.detail.shareModal.shareCode",
                                    "Share Code",
                                )}
                            </span>
                            <div className="flex items-center gap-2 p-2.5 rounded-lg bg-muted/40 border border-border min-w-0">
                                <span
                                    data-testid="share-code"
                                    className="flex-1 min-w-0 text-sm font-mono break-all leading-snug select-all"
                                >
                                    {shareableHandle}
                                </span>
                                <button
                                    data-testid="copy-code-btn"
                                    onClick={handleCopyCode}
                                    aria-label={t(
                                        "participant.certificates.detail.shareModal.copyCode",
                                        "Copy code",
                                    )}
                                    className="shrink-0 p-1.5 rounded hover:bg-foreground/10 transition-colors"
                                >
                                    <Copy className="w-3.5 h-3.5 text-muted-foreground" />
                                </button>
                            </div>
                        </div>

                        {/* Open verify page */}
                        <Link
                            to={verifyPageUrl!}
                            aria-label={t(
                                "participant.certificates.detail.shareModal.openVerifyPage",
                                "Open verify page",
                            )}
                        >
                            <Button variant="secondary-light" className="w-full gap-2">
                                <ExternalLink className="w-4 h-4" />
                                {t(
                                    "participant.certificates.detail.shareModal.openVerifyPage",
                                    "Open verify page",
                                )}
                            </Button>
                        </Link>
                    </div>
                ) : (
                    <p className="text-sm text-muted-foreground text-center py-4">
                        {t(
                            "participant.certificates.shareModal.noHandle",
                            "No share link available.",
                        )}
                    </p>
                )}
            </DialogContent>
        </Dialog>
    );
};
