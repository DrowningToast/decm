import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { useCertificateDetailUsecase } from "./useCertificateDetailUsecase";
import { BottomNav } from "@/components/BottomNav/BottomNav";
import { CircleCheckBig, ExternalLink, Award, Loader2 } from "lucide-react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { usePasswordPrompt } from "@/hooks/usePassowordPrompt";
import { useClaimCertificate } from "@/hooks/useClaimCertificate";
import { useCertificateImage } from "@/hooks/useCertificateImage";

interface CertificateDetailProps {
    certificateId: string;
}

export const CertificateDetail = ({ certificateId }: CertificateDetailProps) => {
    const { t } = useTranslation();
    const { certificate, formattedDate, isLoading, isError } =
        useCertificateDetailUsecase(certificateId);
    const { mutateAsync: openPasswordPrompt } = usePasswordPrompt();
    const { claimCertificate, isClaiming } = useClaimCertificate();
    const [isProcessing, setIsProcessing] = useState(false);

    // Fetch certificate image with authentication
    const {
        imageUrl: certificateImageUrl,
        isLoading: isImageLoading,
        error: imageError,
    } = useCertificateImage({
        certificateId,
        enabled: !!certificate,
    });

    const handleClaimCertificate = async () => {
        if (!certificate) return;

        try {
            setIsProcessing(true);

            // Open password/PIN prompt
            const accountPassword = await openPasswordPrompt({
                eventContractAddress: certificate.eventContractAddress || "",
                transactionType: "Certificate Claim",
                title: t("participant.certificates.claimTitle", "Claim Certificate"),
                description: t(
                    "participant.certificates.claimDescription",
                    "Sign to claim your certificate on the blockchain",
                ),
                details: `Claiming certificate: ${certificate.name}`,
            });

            // Claim the certificate with the verified password
            await claimCertificate({
                certificateId: certificate.id,
                eventId: certificate.eventId,
                accountPassword,
            });
        } catch (error) {
            console.error("Failed to claim certificate:", error);
        } finally {
            setIsProcessing(false);
        }
    };

    const isCertificateClaimed = certificate?.status === "completed";
    const canClaimCertificate = !isCertificateClaimed && !isClaiming && !isProcessing;

    // Loading state
    if (isLoading) {
        return (
            <div className="relative w-full overflow-hidden flex items-center justify-center pb-24 md:pb-12 min-h-[400px]">
                <Loader2 className="w-8 h-8 animate-spin text-primary" />
            </div>
        );
    }

    // Error state
    if (isError) {
        return (
            <div className="relative w-full overflow-hidden flex flex-col items-center justify-center pb-24 md:pb-12 min-h-[400px] gap-4">
                <Typography variant="text" tag="p" color="destructive" className="text-lg">
                    {t("common.error", "An error occurred")}
                </Typography>
                <Typography variant="text" tag="p" color="muted" className="text-sm">
                    {t(
                        "participant.certificates.errorLoadingCertificate",
                        "Failed to load certificate details",
                    )}
                </Typography>
            </div>
        );
    }

    // Certificate not found
    if (!certificate) {
        return (
            <div className="relative w-full overflow-hidden flex flex-col items-center justify-center pb-24 md:pb-12 min-h-[400px] gap-4">
                <Typography variant="text" tag="p" color="muted" className="text-lg">
                    {t("participant.certificates.notFound", "Certificate not found")}
                </Typography>
                <Typography variant="text" tag="p" color="muted" className="text-sm">
                    {t(
                        "participant.certificates.notFoundDescription",
                        "The certificate you're looking for doesn't exist or you don't have access to it.",
                    )}
                </Typography>
            </div>
        );
    }

    return (
        <div className="relative w-full overflow-hidden">
            {/* Main content */}
            <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 py-4 md:py-16 flex flex-col gap-y-4">
                {/* Certificate Header with Check Icon */}
                <div className="flex flex-col gap-y-1.5">
                    <div className="flex items-center gap-2">
                        <CircleCheckBig className="w-6 h-6 text-primary shrink-0" />
                        <Typography
                            variant="header"
                            tag="h1"
                            color="foreground"
                            className="text-[28px] font-header leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {certificate.name}
                        </Typography>
                    </div>
                    {/* Certificate Title and Subtitle */}
                    {(certificate.certificateTitle || certificate.certificateSubtitle) && (
                        <Typography
                            variant="text"
                            tag="p"
                            color="foreground"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {[certificate.certificateTitle, certificate.certificateSubtitle]
                                .filter(Boolean)
                                .join(" ")}
                        </Typography>
                    )}
                    {/* Event Name Link */}
                    <Link
                        to={`/app/events/${certificate.eventId}`}
                        className="underline text-muted"
                    >
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {certificate.event}
                        </Typography>
                    </Link>
                </div>

                {/* Certificate Image */}
                <div className="w-full h-[172px] bg-muted rounded-lg overflow-hidden relative">
                    {isImageLoading ? (
                        <div className="w-full h-full flex items-center justify-center bg-muted/50">
                            <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
                        </div>
                    ) : certificateImageUrl && !imageError ? (
                        <img
                            src={certificateImageUrl}
                            alt={`${certificate.name} - ${certificate.event}`}
                            className="w-full h-full object-contain bg-white"
                            loading="lazy"
                        />
                    ) : (
                        <div className="w-full h-full bg-[#d9d9d9] flex items-center justify-center">
                            <Typography variant="text" tag="p" color="muted" className="text-sm">
                                {imageError
                                    ? t(
                                          "participant.certificates.detail.imageLoadError",
                                          "Failed to load certificate image",
                                      )
                                    : t(
                                          "participant.certificates.detail.certificatePreview",
                                          "Certificate Preview",
                                      )}
                            </Typography>
                        </div>
                    )}
                </div>

                {/* Academic Institution */}
                {certificate.academicInstitution && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t(
                                "participant.certificates.detail.academicInstitution",
                                "Academic Institution",
                            )}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="foreground"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {certificate.academicInstitution}
                        </Typography>
                    </div>
                )}

                {/* Certificate issued/created date - always show */}
                <div className="flex flex-col gap-y-1">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {isCertificateClaimed
                            ? t(
                                  "participant.certificates.detail.certificateIssuedOn",
                                  "Certificate issued on",
                              )
                            : t(
                                  "participant.certificates.detail.certificateAvailableFrom",
                                  "Available to claim from",
                              )}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="foreground"
                        className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {formattedDate}
                    </Typography>
                </div>

                {/* Certificate Contract Address */}
                {certificate.certificateContractAddress && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.certificates.detail.certificateContractAddress")}
                        </Typography>
                        <div className="flex items-center gap-2.5">
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {certificate.certificateContractAddress}
                            </Typography>
                            <ExternalLink className="w-4 h-4 text-foreground shrink-0" />
                        </div>
                    </div>
                )}

                {/* Event Contract Address */}
                {certificate.eventContractAddress && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.certificates.detail.eventContractAddress")}
                        </Typography>
                        <div className="flex items-center gap-2.5">
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {certificate.eventContractAddress}
                            </Typography>
                            <ExternalLink className="w-4 h-4 text-foreground shrink-0" />
                        </div>
                    </div>
                )}

                {/* Verifiable Credential Attributes */}
                {certificate.verifiableCredentialUrl && (
                    <div className="flex flex-col gap-y-1">
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("participant.certificates.detail.verifiableCredentialAttributes")}
                        </Typography>
                        <a
                            href={certificate.verifiableCredentialUrl}
                            className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] text-foreground underline"
                        >
                            {t("participant.certificates.detail.downloadAttributes")}
                        </a>
                    </div>
                )}

                {/* Claim Certificate Button - Only show if not claimed yet */}
                {!isCertificateClaimed && (
                    <div className="mt-6">
                        <Button
                            variant="primary"
                            size="lg"
                            className="w-full"
                            onClick={handleClaimCertificate}
                            disabled={!canClaimCertificate}
                        >
                            {isClaiming || isProcessing ? (
                                <>
                                    <Loader2 className="w-5 h-5 animate-spin" />
                                    {t("participant.certificates.claiming", "Claiming...")}
                                </>
                            ) : (
                                <>
                                    <Award className="w-5 h-5" />
                                    {t("participant.certificates.claimButton", "Claim Certificate")}
                                </>
                            )}
                        </Button>
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-xs text-center mt-2"
                        >
                            {t(
                                "participant.certificates.claimNote",
                                "Sign this certificate to add it to your blockchain credentials",
                            )}
                        </Typography>
                    </div>
                )}

                {/* Claimed Status */}
                {isCertificateClaimed && (
                    <div className="mt-6 p-4 rounded-lg bg-primary/10 border border-primary/20">
                        <div className="flex items-center gap-2 justify-center">
                            <CircleCheckBig className="w-5 h-5 text-primary" />
                            <Typography
                                variant="text"
                                tag="p"
                                color="primary"
                                className="text-sm font-medium"
                            >
                                {t("participant.certificates.claimed", "Certificate Claimed")}
                            </Typography>
                        </div>
                    </div>
                )}
            </div>

            {/* Bottom Navigation */}
            <BottomNav variant="certificate-detail" onBack={() => window.history.back()} />
        </div>
    );
};
