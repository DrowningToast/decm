import { ChevronLeft, Download } from "lucide-react";
import { useBottomContainerContext } from "../context";
import { useCertificateDetailNavStore } from "../stores/certificates";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { Typography } from "@/components/typography/typography";

interface CertificateDetailNavProps {
    className?: string;
}

export const CertificateDetailNav = ({ className: propClassName }: CertificateDetailNavProps) => {
    const { certificateId, isClaimed, imageUrl } = useCertificateDetailNavStore();
    const { onBack, className: contextClassName } = useBottomContainerContext();
    const { t } = useTranslation();

    const handleDownload = async () => {
        if (!imageUrl || !certificateId) {
            console.error("Certificate image not available for download");
            // TODO: Show error toast notification
            return;
        }

        try {
            // Fetch the image as a blob
            const response = await fetch(imageUrl);
            if (!response.ok) {
                throw new Error("Failed to fetch certificate image");
            }

            const blob = await response.blob();

            // Create a download link
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement("a");
            link.href = url;
            link.download = `certificate-${certificateId}.png`;

            // Trigger download
            document.body.appendChild(link);
            link.click();

            // Cleanup
            document.body.removeChild(link);
            window.URL.revokeObjectURL(url);

            // TODO: Show success toast notification
            console.log("Certificate downloaded successfully");
        } catch (error) {
            console.error("Failed to download certificate:", error);
            // TODO: Show error toast notification
        }
    };

    const handleCopyShareableUrl = () => {
        // Copy current URL to clipboard
        navigator.clipboard
            .writeText(window.location.href)
            .then(() => {
                // TODO: Show success toast notification
                console.log("URL copied to clipboard");
            })
            .catch((err) => {
                console.error("Failed to copy URL:", err);
            });
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

            {/* Download Button - Hidden for now */}
            {false && isClaimed && (
                <button
                    onClick={handleDownload}
                    className="cursor-pointer flex items-center justify-center gap-2 px-4 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0"
                    aria-label={t(
                        "participant.certificates.detail.downloadAsImage",
                        "Download as an image",
                    )}
                >
                    <Download className="w-5 h-5 text-background-alt" />
                    <Typography
                        variant="text"
                        tag="span"
                        color="background-alt"
                        className="text-xs font-normal leading-normal tracking-[0.06px] whitespace-nowrap"
                    >
                        {t(
                            "participant.certificates.detail.downloadAsImage",
                            "Download as an image",
                        )}
                    </Typography>
                </button>
            )}

            {/* Copy Shareable URL Button - Only show if certificate is claimed */}
            {isClaimed && (
                <button
                    onClick={handleCopyShareableUrl}
                    className="cursor-pointer flex items-center justify-center flex-1 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors"
                    aria-label={t("participant.certificates.detail.copyShareableUrl")}
                >
                    <Typography
                        variant="text"
                        tag="span"
                        color="background-alt"
                        className="text-xs font-normal leading-normal tracking-[0.06px] text-center whitespace-nowrap"
                    >
                        {t("participant.certificates.detail.copyShareableUrl")}
                    </Typography>
                </button>
            )}
        </div>
    );
};
