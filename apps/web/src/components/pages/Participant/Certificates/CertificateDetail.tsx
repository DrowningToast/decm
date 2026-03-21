import { useState, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { useCertificateDetailUsecase } from "./useCertificateDetailUsecase";
import { BottomNav, type BottomNavVariant } from "@/components/BottomNav/BottomNav";
import { CircleCheckBig, Award, Loader2, UserPlus, CircleX } from "lucide-react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { usePasswordPrompt } from "@/hooks/usePassowordPrompt";
import { useClaimCertificate } from "@/hooks/useClaimCertificate";
import { useCertificateImage } from "@/hooks/useCertificateImage";
import {
    useCertificateDetailNavStore,
    useCertificateDetailsSharedNavStore,
} from "@/components/BottomNav/stores/certificates";
import { CertificatePasswordDialog } from "./CertificatePasswordDialog";
import { CertificateShareModal } from "./CertificateShareModal";
import { EthExplorerLink } from "@/components/common/EthscanLink";
import { useEvent } from "@/hooks/events/useEvent";
import { useCertificateCreateShareLinkUsecase } from "./useCertificateCreateShareLinkUsecase";
import { useCertificateUpdateShareVisibilityUsecase } from "./useCertificateUpdateShareVisibilityUsecase";

interface CertificateDetailProps {
    certificateId: string;
}

export const CertificateDetail = ({ certificateId }: CertificateDetailProps) => {
    const { t, i18n } = useTranslation();
    const { certificate, formattedDate, isLoading, isError } =
        useCertificateDetailUsecase(certificateId);
    const { mutateAsync: openPasswordPrompt } = usePasswordPrompt();
    const { claimCertificate, isClaiming } = useClaimCertificate();
    const [isProcessing, setIsProcessing] = useState(false);
    const { setOnClickShareable, setImageUrl } = useCertificateDetailNavStore();
    const {
        setOnChangePublish,
        setIsPublished,
        setIsPasswordProtected,
        setShareableHandle,
        isPasswordDialogOpen,
        setIsPasswordDialogOpen,
    } = useCertificateDetailsSharedNavStore();
    const { createCertificateShareLink } = useCertificateCreateShareLinkUsecase();
    const { updateShareVisibility } = useCertificateUpdateShareVisibilityUsecase();

    useEffect(() => {
        setOnClickShareable(() => {
            void createCertificateShareLink(certificateId);
        });
    }, [certificateId, createCertificateShareLink, setOnClickShareable]);

    const shareId = certificate?.shareable?.id ?? null;

    useEffect(() => {
        if (!shareId) return;
        setOnChangePublish(async (isPublished: boolean) => {
            setIsPublished(isPublished);
            await updateShareVisibility(shareId, isPublished);
        });
    }, [shareId, updateShareVisibility, setOnChangePublish, setIsPublished]);

    // Sync shareable state into the shared nav store when certificate data loads
    useEffect(() => {
        if (certificate?.shareable) {
            setIsPublished(certificate.shareable.active);
            setIsPasswordProtected(certificate.shareable.hasPassword);
            setShareableHandle(certificate.shareable.handle);
        }
        return () => {
            setIsPublished(false);
            setIsPasswordProtected(false);
            setShareableHandle(null);
        };
    }, [certificate?.shareable, setIsPublished, setIsPasswordProtected, setShareableHandle]);

    // Fetch event details to check if user has joined
    const { event, isLoadingEvent } = useEvent(certificate?.eventId || "");

    // Fetch certificate image with authentication
    const {
        imageUrl: certificateImageUrl,
        isLoading: isImageLoading,
        error: imageError,
    } = useCertificateImage({
        certificateId,
        enabled: !!certificate,
    });

    // Update the store with the image URL when available
    useEffect(() => {
        if (certificateImageUrl) {
            setImageUrl(certificateImageUrl);
        }
        return () => {
            setImageUrl(null);
        };
    }, [certificateImageUrl, setImageUrl]);

    // Store image URL in nav store for download functionality
    useEffect(() => {
        setImageUrl(certificateImageUrl);
        return () => {
            setImageUrl(null);
        };
    }, [certificateImageUrl, setImageUrl]);

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
                    "Enter your account password to claim your certificate on the blockchain",
                ),
                details: t("participant.certificates.claimingCertificate", {
                    name: certificate.name,
                }),
            });

            // Claim the certificate with the verified password (PIN flow)
            await claimCertificate({
                certificateId: certificate.id,
                accountPassword,
            });
        } catch (error) {
            console.error("Failed to claim certificate:", error);
            // Error toast is already shown by the hook
        } finally {
            setIsProcessing(false);
        }
    };

    const isCertificateClaimed = certificate?.status === "completed";
    const isCertificateQueued = certificate?.status === "queued";
    const isCertificateExpired = certificate?.status === "expired";
    const isCertificateAborted = certificate?.status === "aborted";
    const canClaimCertificate =
        !isCertificateClaimed &&
        !isCertificateQueued &&
        !isCertificateExpired &&
        !isClaiming &&
        !isProcessing;

    const bottomNavVariant: BottomNavVariant = useMemo(() => {
        if (certificate?.shareable) {
            return "certificate-details-shared";
        }
        return "certificate-detail";
    }, [certificate?.shareable]);

    const formattedDeadline = certificate?.estimatedDeadline
        ? new Date(certificate.estimatedDeadline).toLocaleDateString(
              i18n.language === "th" ? "th-TH" : "en-US",
              {
                  year: "numeric",
                  month: "short",
                  day: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
              },
          )
        : undefined;

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
            <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 pt-4 md:pt-6 pb-24 md:pb-16 flex flex-col gap-y-4">
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
                            {certificate.certificateTitle ||
                                t(
                                    "participant.certificates.detail.certificateFromEvent",
                                    "Certificate from {{eventName}}",
                                    { eventName: certificate.event },
                                )}
                        </Typography>
                    </div>
                    {/* Participant Name and Certificate Subtitle */}
                    <Typography
                        variant="text"
                        tag="p"
                        color="foreground"
                        className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                    >
                        {certificate.name}
                        {certificate.certificateSubtitle && ` • ${certificate.certificateSubtitle}`}
                    </Typography>
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
                <div className="w-full xl:max-w-2xl 2xl:max-w-3xl mx-auto rounded-lg overflow-hidden relative">
                    {/* Aspect ratio container using padding (4:3 ratio) */}
                    <div className="relative w-full" style={{ paddingBottom: "75%" }}>
                        {isImageLoading ? (
                            <div className="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-muted/50 p-4">
                                <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
                                <Typography
                                    variant="text"
                                    size="small"
                                    tag="p"
                                    color="muted"
                                    className="text-center"
                                >
                                    {t(
                                        "participant.certificates.detail.loadingImageHint",
                                        "Loading certificate preview... This may take a while on the first load.",
                                    )}
                                </Typography>
                            </div>
                        ) : certificateImageUrl && !imageError ? (
                            <img
                                src={certificateImageUrl}
                                alt={`${certificate.name} - ${certificate.event}`}
                                className="absolute inset-0 w-full h-full object-contain"
                                loading="lazy"
                            />
                        ) : (
                            <div className="absolute inset-0 bg-muted/30 flex items-center justify-center">
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-sm"
                                >
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
                        <EthExplorerLink
                            address={certificate.certificateContractAddress}
                            className="flex items-center gap-2.5"
                        >
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {certificate.certificateContractAddress}
                            </Typography>
                        </EthExplorerLink>
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
                        <EthExplorerLink
                            address={certificate.eventContractAddress}
                            className="flex items-center gap-2.5"
                        >
                            <Typography
                                variant="text"
                                tag="p"
                                color="foreground"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {certificate.eventContractAddress}
                            </Typography>
                        </EthExplorerLink>
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

                {/* Queued Status — minting in progress */}
                {isCertificateQueued && (
                    <div className="p-4 rounded-lg bg-muted/20 border border-muted/40">
                        <div className="flex flex-col items-center gap-2">
                            <div className="flex items-center gap-2">
                                <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-sm font-medium"
                                >
                                    {t(
                                        "participant.certificates.queued",
                                        "Certificate minting in progress...",
                                    )}
                                </Typography>
                            </div>
                            {formattedDeadline && (
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-xs text-center"
                                >
                                    {t(
                                        "participant.certificates.mintingBy",
                                        "Expected on blockchain by {{deadline}}",
                                        { deadline: formattedDeadline },
                                    )}
                                </Typography>
                            )}
                        </div>
                    </div>
                )}

                {/* Expired Status — minting deadline passed without broadcast */}
                {isCertificateExpired && (
                    <div className="p-4 rounded-lg bg-destructive/10 border border-destructive/20">
                        <div className="flex items-center gap-2 justify-center">
                            <CircleX className="w-5 h-5 text-destructive" />
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-sm font-medium text-destructive"
                            >
                                {t(
                                    "participant.certificates.claimExpired",
                                    "Certificate claim expired",
                                )}
                            </Typography>
                        </div>
                    </div>
                )}

                {/* Aborted Status — claim was aborted (e.g. not joined on-chain at the time) */}
                {isCertificateAborted && (
                    <div className="p-4 rounded-lg bg-destructive/10 border border-destructive/20">
                        <div className="flex flex-col items-center gap-2">
                            <div className="flex items-center gap-2">
                                <CircleX className="w-5 h-5 text-destructive" />
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-sm font-medium text-destructive"
                                >
                                    {t(
                                        "participant.certificates.claimAborted",
                                        "Certificate claim could not be processed",
                                    )}
                                </Typography>
                            </div>
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted"
                                className="text-xs text-center"
                            >
                                {t(
                                    "participant.certificates.claimAbortedDescription",
                                    "This claim was cancelled because you have not joined the event on the blockchain. Please ensure you have joined the event and try claiming again.",
                                )}
                            </Typography>
                        </div>
                    </div>
                )}

                {/* Claim / Retry Certificate Button */}
                {!isCertificateClaimed && !isCertificateQueued && !isCertificateExpired && (
                    <div className="mt-6">
                        {/* Check if user has joined the event */}
                        {!isLoadingEvent && event && !event.isJoined ? (
                            <>
                                <div className="p-4 rounded-lg bg-yellow-500/10 border border-yellow-500/20 mb-4">
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="foreground"
                                        className="text-sm text-center"
                                    >
                                        {t(
                                            "participant.certificates.mustJoinEventFirst",
                                            "You need to join this event as a participant before claiming the certificate",
                                        )}
                                    </Typography>
                                </div>
                                <Link to={`/app/events/${certificate.eventId}`}>
                                    <Button variant="primary" size="lg" className="w-full">
                                        <UserPlus className="w-5 h-5" />
                                        {t(
                                            "participant.certificates.goToEvent",
                                            "Go to Event Page",
                                        )}
                                    </Button>
                                </Link>
                            </>
                        ) : (
                            <>
                                <Button
                                    variant="primary"
                                    size="lg"
                                    className="w-full"
                                    onClick={handleClaimCertificate}
                                    disabled={!canClaimCertificate || isLoadingEvent}
                                >
                                    {isClaiming || isProcessing ? (
                                        <>
                                            <Loader2 className="w-5 h-5 animate-spin" />
                                            {t("participant.certificates.claiming", "Claiming...")}
                                        </>
                                    ) : isCertificateAborted ? (
                                        <>
                                            <Award className="w-5 h-5" />
                                            {t(
                                                "participant.certificates.retryClaimButton",
                                                "Retry Claim",
                                            )}
                                        </>
                                    ) : (
                                        <>
                                            <Award className="w-5 h-5" />
                                            {t(
                                                "participant.certificates.claimButton",
                                                "Claim Certificate",
                                            )}
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
                            </>
                        )}
                    </div>
                )}

                {/* Claimed Status */}
                {isCertificateClaimed && (
                    <div className="p-4 rounded-lg bg-primary/10 border border-primary/20">
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
            <BottomNav variant={bottomNavVariant} onBack={() => window.history.back()} />

            {/* Share modal */}
            <CertificateShareModal />

            {/* Password dialog */}
            <CertificatePasswordDialog
                open={isPasswordDialogOpen}
                onOpenChange={setIsPasswordDialogOpen}
                shareId={shareId}
                hasPassword={certificate.shareable?.hasPassword ?? false}
            />
        </div>
    );
};
