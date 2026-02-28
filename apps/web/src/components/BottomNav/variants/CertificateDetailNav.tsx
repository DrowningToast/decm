import { ChevronLeft, Download } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useCertificateDetailNavStore } from "../stores/certificates";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { Typography } from "@/components/typography/typography";
import { certificateService } from "@/services/services";
import { toast } from "sonner";
import { useState } from "react";

interface CertificateDetailNavProps {
    className?: string;
}

export const CertificateDetailNav = ({ className: propClassName }: CertificateDetailNavProps) => {
    const { certificateId, isClaimed } = useCertificateDetailNavStore();
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
                "flex items-center gap-1.5 h-13 bg-primary rounded-xl p-1.5",
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

            {/* Download Button - only shown when certificate is claimed */}
            {isClaimed && (
                <button
                    onClick={handleDownload}
                    disabled={isDownloading}
                    className="cursor-pointer flex items-center justify-center gap-2 px-4 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0 disabled:opacity-50 disabled:cursor-not-allowed"
                    aria-label={t(
                        "participant.certificates.detail.downloadAsImage",
                        "Download as an image",
                    )}
                >
                    <Download
                        className={cn(
                            "w-5 h-5 text-background-alt",
                            isDownloading && "animate-pulse",
                        )}
                    />
                    <Typography
                        variant="text"
                        tag="span"
                        color="background-alt"
                        className="text-xs font-normal leading-normal tracking-[0.06px] whitespace-nowrap"
                    >
                        {isDownloading
                            ? t("participant.certificates.detail.downloading", "Downloading...")
                            : t(
                                  "participant.certificates.detail.downloadAsImage",
                                  "Download as an image",
                              )}
                    </Typography>
                </button>
            )}
        </div>
    );
};
