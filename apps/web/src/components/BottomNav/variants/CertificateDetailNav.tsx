import { ChevronLeft, Download, Share2 } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useCertificateDetailNavStore } from "../stores/certificates";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { certificateService } from "@/services/services";
import { toast } from "sonner";
import React, { useState } from "react";
import type { ClassValue } from "clsx";
import { Button } from "@/components/ui/button";

interface CertificateDetailNavProps extends React.PropsWithChildren {
    className?: ClassValue;
    overrideShowCreateShareButton?: boolean;
}

export const CertificateDetailNav = ({
    className: propClassName,
    overrideShowCreateShareButton,
    children,
}: CertificateDetailNavProps) => {
    const { certificateId, isClaimed, onClickShareable, isShareableLoading, setIsShareModalOpen } =
        useCertificateDetailNavStore();
    const { onBack, className: contextClassName } = useBottomContainerContext();
    const { t } = useTranslation();
    const [isDownloading, setIsDownloading] = useState(false);

    const handleDownload = async () => {
        if (!certificateId) {
            toast.error(
                t(
                    "participant.certificates.detail.downloadError",
                    "Certificate ID is missing. Please try again.",
                ),
            );
            return;
        }

        setIsDownloading(true);

        try {
            // Fetch the certificate image using the service (includes authentication)
            const certificateImage = await certificateService.getCertificateImage(certificateId);

            // Create a download link from the blob URL
            const link = document.createElement("a");
            link.href = certificateImage.url;
            link.download = `certificate-${certificateId}.png`;

            // Trigger download
            document.body.appendChild(link);
            link.click();

            // Cleanup
            document.body.removeChild(link);

            toast.success(
                t(
                    "participant.certificates.detail.downloadSuccess",
                    "Certificate downloaded successfully!",
                ),
            );
        } catch (error) {
            console.error("Failed to download certificate:", error);
            const errorMessage =
                error instanceof Error
                    ? error.message
                    : t(
                          "participant.certificates.detail.downloadError",
                          "Failed to download certificate. Please try again.",
                      );
            toast.error(errorMessage);
        } finally {
            setIsDownloading(false);
        }
    };

    return (
        <div
            className={cn(
                contextClassName,
                propClassName,
                "flex items-center justify-center gap-1.5 bg-primary rounded-xl p-1.5",
            )}
        >
            {/* Back Button */}
            <button
                onClick={onBack}
                className="cursor-pointer flex items-center justify-center w-10 h-10 rounded-lg hover:opacity-90 transition-opacity flex-shrink-0"
                aria-label={t("common.back")}
            >
                <ChevronLeft className="w-5 h-5 text-white" />
            </button>
            <div className="flex gap-1.5 flex-wrap justify-center items-stretch">
                {/* Shareable Link Button - only shown when certificate is claimed */}
                {isClaimed && overrideShowCreateShareButton !== false && (
                    <Button
                        loading={isShareableLoading}
                        onClick={onClickShareable}
                        className="cursor-pointer flex items-center justify-center gap-2 px-3 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0"
                        aria-label={t(
                            "participant.certificates.detail.createShareableLink",
                            "Create shareable link",
                        )}
                    >
                        <Share2 className="w-5 h-5 text-background-alt flex-shrink-0" />
                        <span className="text-xs text-background-alt whitespace-nowrap">
                            {t("participant.certificates.detail.share", "Share")}
                        </span>
                    </Button>
                )}

                {/* Download Button - only shown when certificate is claimed */}
                {isClaimed && (
                    <button
                        onClick={handleDownload}
                        disabled={isDownloading}
                        className="cursor-pointer flex items-center justify-center gap-2 px-3 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0 disabled:opacity-50 disabled:cursor-not-allowed"
                        aria-label={t(
                            "participant.certificates.detail.downloadAsImage",
                            "Download as an image",
                        )}
                    >
                        <Download
                            className={cn(
                                "w-5 h-5 text-background-alt flex-shrink-0",
                                isDownloading && "animate-pulse",
                            )}
                        />
                        <span className="text-xs text-background-alt whitespace-nowrap">
                            {t("participant.certificates.detail.download", "Download")}
                        </span>
                    </button>
                )}

                {/* Share Modal Button - only shown when certificate is claimed */}
                {isClaimed && (
                    <button
                        onClick={() => setIsShareModalOpen(true)}
                        className="cursor-pointer flex items-center justify-center gap-2 px-3 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0"
                        aria-label={t(
                            "participant.certificates.detail.openShareModal",
                            "Share certificate",
                        )}
                    >
                        <Share2 className="w-5 h-5 text-background-alt flex-shrink-0" />
                        <span className="text-xs text-background-alt whitespace-nowrap">
                            {t(
                                "participant.certificates.detail.openShareModal",
                                "Share certificate",
                            )}
                        </span>
                    </button>
                )}

                {children}
            </div>
        </div>
    );
};
