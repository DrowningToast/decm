import { useState, useMemo } from "react";
import { useNavigate } from "@/router";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { useEvent } from "@/hooks/events/useEvent";
import { useEventIssuers } from "@/components/pages/HostPages/EventPages/useEventIssuers";
import { useEventCertificates } from "@/hooks/useEventCertificates";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useEventContract } from "@/hooks/events/useEventContracts";
import { useSignEventCertificates } from "@/hooks/useSignEventCertificates";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { ChevronLeft, ExternalLinkIcon, CheckCircle, Clock } from "lucide-react";
import { PasswordPinModal } from "@/components/ui/password-pin-modal";
import { TextLabelValue } from "@/components/ui/text-label-value";
import { DataTable } from "@/components/ui/data-table";
import { useCertificateColumns } from "@/components/pages/HostPages/EventsPage/columns/CertificateColumns";
import { formatEthereumAddress } from "@/lib/utils";
import SectionContainer from "@/components/container/SectionContainer";
import { IssuerStatusBadge } from "./IssuerStatusBadge";
import { IssuersStatus } from "./IssuersStatus";
import { useAuth } from "@/context/AuthContext";
import { CertificatePreview } from "@/components/CertificatePreview";
import { useFetchedCertificateSvg } from "@/components/pages/HostPages/EventPages/useFetchedCertificateSvg";
import type { CertificateFontConfig } from "@/lib/certificateRenderer";
import type { DetectedKeyword } from "@/hooks/useCertificateTemplate";

interface IssuerSignPageProps {
    eventId: string;
}

export default function IssuerSignPage({ eventId }: IssuerSignPageProps) {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [showPinModal, setShowPinModal] = useState(false);
    const { user: currentUser } = useAuth();

    // Fetch event data
    const { event, isLoadingEventError } = useEvent(eventId);
    const { eventIssuers } = useEventIssuers(eventId);
    const { certificates: eventCertificates, isLoading: isLoadingCertificates } =
        useEventCertificates(eventId);
    const { data: eventCertificateConfig } = useEventCertificateConfig(eventId);
    const { data: eventContract } = useEventContract(eventId);
    const { signEventCertificates, isSigning } = useSignEventCertificates();

    // Fetch SVG template through the backend proxy (avoids CORS on S3 presigned URLs)
    const fetchedSvg = useFetchedCertificateSvg(eventId, false);

    // Derive font config from the saved certificate config
    const fontConfig = useMemo<CertificateFontConfig | undefined>(() => {
        if (!eventCertificateConfig) return undefined;
        return {
            event_name_font_family_id: eventCertificateConfig.event_name_font_family_id,
            event_name_font_weight: eventCertificateConfig.event_name_font_weight,
            name_font_family_id: eventCertificateConfig.name_font_family_id,
            name_font_weight: eventCertificateConfig.name_font_weight,
            academic_institution_font_family_id:
                eventCertificateConfig.academic_institution_font_family_id,
            academic_institution_font_weight:
                eventCertificateConfig.academic_institution_font_weight,
            certificate_title_font_family_id:
                eventCertificateConfig.certificate_title_font_family_id,
            certificate_title_font_weight: eventCertificateConfig.certificate_title_font_weight,
            certificate_subtitle_font_family_id:
                eventCertificateConfig.certificate_subtitle_font_family_id,
            certificate_subtitle_font_weight:
                eventCertificateConfig.certificate_subtitle_font_weight,
        };
    }, [eventCertificateConfig]);

    // Reconstruct keyword positions from the saved config (same logic as CertificateSettingsPage)
    const detectedKeywords = useMemo<DetectedKeyword[]>(() => {
        if (!eventCertificateConfig) return [];
        const keywords: DetectedKeyword[] = [];

        if (
            eventCertificateConfig.name_pos_x !== undefined &&
            eventCertificateConfig.name_pos_y !== undefined &&
            (eventCertificateConfig.name_pos_x !== 0 || eventCertificateConfig.name_pos_y !== 0)
        ) {
            keywords.push({
                keyword: "{{ name }}",
                x: eventCertificateConfig.name_pos_x,
                y: eventCertificateConfig.name_pos_y,
                count: 1,
            });
        }
        if (
            eventCertificateConfig.event_name_pos_x !== undefined &&
            eventCertificateConfig.event_name_pos_y !== undefined &&
            (eventCertificateConfig.event_name_pos_x !== 0 ||
                eventCertificateConfig.event_name_pos_y !== 0)
        ) {
            keywords.push({
                keyword: "{{ eventName }}",
                x: eventCertificateConfig.event_name_pos_x,
                y: eventCertificateConfig.event_name_pos_y,
                count: 1,
            });
        }
        if (
            eventCertificateConfig.academic_institution_pos_x !== undefined &&
            eventCertificateConfig.academic_institution_pos_y !== undefined
        ) {
            keywords.push({
                keyword: "{{ academicInstitutionName }}",
                x: eventCertificateConfig.academic_institution_pos_x,
                y: eventCertificateConfig.academic_institution_pos_y,
                count: 1,
            });
        }
        if (
            eventCertificateConfig.certificate_title_pos_x !== undefined &&
            eventCertificateConfig.certificate_title_pos_y !== undefined
        ) {
            keywords.push({
                keyword: "{{ certificateTitle }}",
                x: eventCertificateConfig.certificate_title_pos_x,
                y: eventCertificateConfig.certificate_title_pos_y,
                count: 1,
            });
        }
        if (
            eventCertificateConfig.certificate_subtitle_pos_x !== undefined &&
            eventCertificateConfig.certificate_subtitle_pos_y !== undefined
        ) {
            keywords.push({
                keyword: "{{ certificateSubtitle }}",
                x: eventCertificateConfig.certificate_subtitle_pos_x,
                y: eventCertificateConfig.certificate_subtitle_pos_y,
                count: 1,
            });
        }

        return keywords;
    }, [eventCertificateConfig]);

    // Get current issuer's information
    const currentIssuer = eventIssuers?.find(
        (issuer) => issuer.issuerCredentialId === currentUser?.authenticationCredentialId,
    );

    // Determine if current issuer has already signed
    const hasCurrentIssuerSigned = currentIssuer && currentIssuer.isSigned;
    const isCurrentIssuerPending = currentIssuer && !currentIssuer.isSigned;

    // Calculate total certificates to be signed
    const certificatesToSign =
        eventCertificates?.filter((cert) => !cert.revokedAt && currentIssuer) || [];

    // Get certificate columns
    const certificateColumns = useCertificateColumns();

    const handleSignCertificates = () => {
        setShowPinModal(true);
    };

    const handleGoBack = () => {
        navigate(-1);
    };

    const handlePinModalClose = () => {
        setShowPinModal(false);
    };

    const handleSignWithPin = (pin: string) => {
        signEventCertificates({
            eventId,
            issuerPin: pin,
        });
    };

    if (isLoadingEventError || !event) {
        return (
            <ProtectedRoute>
                <div className="flex items-center justify-center min-h-screen">
                    <div className="text-center">
                        <Typography variant="text" tag="p" color="destructive">
                            {t("events.errors.generic")}
                        </Typography>
                        <Button onClick={handleGoBack} className="mt-4">
                            {t("common.back")}
                        </Button>
                    </div>
                </div>
            </ProtectedRoute>
        );
    }

    return (
        <ProtectedRoute>
            <div className="min-h-screen text-white flex flex-col gap-y-4 mb-8">
                {/* Header */}
                <header className="bg-[#ff6a39] shadow-lg">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                        <div className="flex justify-between h-16 items-center">
                            <div
                                className="flex items-center gap-2 cursor-pointer"
                                onClick={handleGoBack}
                            >
                                <ChevronLeft className="h-5 w-5" />
                                <Typography variant="text" tag="span" color="foreground-alt">
                                    {t("issuer.sign.backToDashboard")}
                                </Typography>
                            </div>
                        </div>
                    </div>
                </header>

                {/* Main Content */}
                <SectionContainer className="mt-4">
                    <main className="flex flex-col gap-y-8">
                        {/* Page Header */}
                        <div>
                            <Typography
                                variant="header"
                                tag="h1"
                                className="text-4xl font-bold text-white mb-2"
                            >
                                {t("issuer.sign.pageTitle")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="muted-foreground"
                                className="text-lg"
                            >
                                {t("issuer.sign.pageDescription")}
                            </Typography>
                        </div>
                        {/* Issuers Status Section */}
                        {eventIssuers && eventIssuers.length > 0 && (
                            <IssuersStatus
                                issuers={
                                    eventIssuers?.map((issuer) => ({
                                        id: issuer.id || "",
                                        issuer_credential_id: issuer.issuerCredentialId || "",
                                        is_signed: issuer.isSigned ? 1 : 0,
                                        issuerProfile: issuer.issuerProfile
                                            ? {
                                                  firstName: issuer.issuerProfile.firstName,
                                                  lastName: issuer.issuerProfile.lastName,
                                                  walletAddress: issuer.issuerProfile.walletAddress,
                                                  googleConnectorRef:
                                                      issuer.issuerProfile.googleConnectorRef,
                                              }
                                            : undefined,
                                    })) || []
                                }
                                currentIssuerId={currentUser?.profileId}
                            />
                        )}

                        {/* Signing Details Section */}
                        <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-6">
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-semibold text-white mb-4"
                            >
                                {t("issuer.sign.signingDetails")}
                            </Typography>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="muted-foreground"
                                        className="text-sm"
                                    >
                                        {t("issuer.sign.signingStatus")}
                                    </Typography>
                                    <div className="mt-1">
                                        {currentIssuer ? (
                                            <div className="flex items-center space-x-2">
                                                <IssuerStatusBadge
                                                    isSigned={!!currentIssuer.isSigned}
                                                />
                                                {hasCurrentIssuerSigned && (
                                                    <Typography
                                                        variant="text"
                                                        tag="span"
                                                        className="text-sm text-green-500"
                                                    >
                                                        Completed
                                                    </Typography>
                                                )}
                                            </div>
                                        ) : (
                                            "N/A"
                                        )}
                                    </div>
                                </div>
                                <TextLabelValue
                                    label={t("issuer.sign.contractAddress")}
                                    value={
                                        eventContract?.certificate_contract_address
                                            ? formatEthereumAddress(
                                                  eventContract.certificate_contract_address,
                                              )
                                            : "N/A"
                                    }
                                    endIcon={<ExternalLinkIcon size={16} />}
                                    valueClassName="cursor-pointer underline"
                                    href={`https://www.etherscan.io/address/${eventContract?.certificate_contract_address}`}
                                />
                                <TextLabelValue
                                    label={t("issuer.sign.certificatesCount")}
                                    value={`${certificatesToSign.length} ${t("issuer.sign.certificates")}`}
                                />

                                <TextLabelValue
                                    label="Event"
                                    value={event.title ?? ""}
                                    endIcon={<ExternalLinkIcon size={16} />}
                                    valueClassName="cursor-pointer underline"
                                    href={`/events/${eventId}`}
                                />

                                <TextLabelValue
                                    label="Start Date"
                                    value={new Date(event.startDate ?? "").toLocaleDateString()}
                                />
                                <TextLabelValue
                                    label="End Date"
                                    value={new Date(event.endDate ?? "").toLocaleDateString()}
                                />
                            </div>
                        </div>

                        {/* Main Content Grid */}
                        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                            {/* Left Column - All Sections */}
                            <div className="lg:col-span-2 space-y-8">
                                {/* Certificate Preview Section */}
                                <div>
                                    <Typography
                                        variant="header"
                                        tag="h2"
                                        className="text-xl font-semibold text-white mb-4"
                                    >
                                        {t("issuer.sign.certificatePreview")}
                                    </Typography>
                                    <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-4">
                                        {eventCertificateConfig ? (
                                            <CertificatePreview
                                                svgPreview={fetchedSvg || undefined}
                                                imageUrl={
                                                    eventCertificateConfig.base_certificate_presigned_url ??
                                                    undefined
                                                }
                                                fontConfig={fontConfig}
                                                detectedKeywords={detectedKeywords}
                                                previewData={{
                                                    eventName: event?.title ?? undefined,
                                                }}
                                            />
                                        ) : (
                                            <div className="flex items-center justify-center border border-dashed border-[#333333] rounded aspect-video">
                                                <div className="text-center text-[#a0a0a0]">
                                                    <Typography variant="text" tag="p">
                                                        {t("issuer.sign.templateNotAvailable")}
                                                    </Typography>
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                </div>

                                {/* Event Certificates Section - Desktop: Below Preview */}
                                <div className="hidden lg:block w-full">
                                    <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-6 w-full">
                                        <Typography
                                            variant="header"
                                            tag="h2"
                                            className="text-xl font-semibold text-white mb-4"
                                        >
                                            {isCurrentIssuerPending
                                                ? t("issuer.sign.certificatesToSign")
                                                : hasCurrentIssuerSigned
                                                  ? t("issuer.sign.signedCertificates")
                                                  : t("issuer.sign.eventCertificates")}{" "}
                                            (
                                            {eventCertificates?.filter((cert) => !cert.revokedAt)
                                                .length || 0}
                                            )
                                        </Typography>

                                        <div className="w-full overflow-hidden">
                                            <DataTable
                                                columns={certificateColumns}
                                                data={
                                                    eventCertificates
                                                        ?.filter((cert) => {
                                                            // Filter out revoked certificates and ensure required fields exist
                                                            return (
                                                                cert.id &&
                                                                cert.eventId &&
                                                                !cert.revokedAt &&
                                                                cert.createdAt
                                                            );
                                                        })
                                                        .map((cert) => {
                                                            const firstName =
                                                                cert.name?.split(" ")[0] || "";
                                                            const lastName =
                                                                cert.name
                                                                    ?.split(" ")
                                                                    .slice(1)
                                                                    .join(" ") || "";

                                                            return {
                                                                id: cert.id!,
                                                                event_id: cert.eventId!,
                                                                receiver_credential_id:
                                                                    cert.receiverCredentialId,
                                                                receiver_email: cert.receiverEmail,
                                                                name: cert.name,
                                                                academic_institution:
                                                                    cert.academicInstitution,
                                                                certificate_title:
                                                                    cert.certificateTitle,
                                                                certificate_subtitle:
                                                                    cert.certificateSubtitle,
                                                                event_contract_address:
                                                                    cert.eventContractAddress || "",
                                                                event_certificate_address:
                                                                    cert.eventCertificateAddress,
                                                                created_at: cert.createdAt!,
                                                                revoked_at: cert.revokedAt,
                                                                firstName,
                                                                lastName,
                                                                email: cert.receiverEmail || "",
                                                                academicInstitution:
                                                                    cert.academicInstitution || "",
                                                                issuedAt: cert.createdAt!,
                                                                status: cert.revokedAt
                                                                    ? "rejected"
                                                                    : isCurrentIssuerPending
                                                                      ? "pending_signature"
                                                                      : "received",
                                                            };
                                                        }) || []
                                                }
                                                totalItems={
                                                    eventCertificates?.filter(
                                                        (cert) =>
                                                            cert.id &&
                                                            cert.eventId &&
                                                            !cert.revokedAt &&
                                                            cert.createdAt,
                                                    ).length || 0
                                                }
                                                currentPage={1}
                                                pageSize={10}
                                                onPageChange={() => {}}
                                                onPageSizeChange={() => {}}
                                                searchValue=""
                                                onSearchChange={() => {}}
                                                searchPlaceholder="Search certificates..."
                                                sorting={[]}
                                                onSortingChange={() => {}}
                                                isLoading={isLoadingCertificates}
                                                disablePagination
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Right Column - Action Card */}
                            <div className="lg:col-span-1 flex flex-col gap-y-8">
                                {/* Event Details Section */}
                                <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-6 flex flex-col">
                                    <Typography
                                        variant="header"
                                        tag="h2"
                                        className="text-xl font-semibold text-white mb-4"
                                    >
                                        {t("issuer.sign.eventDetails")}
                                    </Typography>
                                    <div className="space-y-4">
                                        <TextLabelValue
                                            label={t("events.details.status")}
                                            value={event.eventStatus?.toUpperCase() ?? "NA"}
                                        />
                                        <TextLabelValue
                                            label={t("events.details.participationRequest")}
                                            value={
                                                event.isBookingRequestRequired
                                                    ? t("common.required")
                                                    : t("common.notRequired")
                                            }
                                        />
                                        <TextLabelValue
                                            label={t("events.details.seatsCount")}
                                            value={`${event.attendeesCount ?? 0} / ${event.maxAttendees}`}
                                        />
                                        <TextLabelValue
                                            label={t("events.details.eventContractAddress")}
                                            value={
                                                eventContract?.certificate_contract_address
                                                    ? formatEthereumAddress(
                                                          eventContract.certificate_contract_address,
                                                      )
                                                    : "NA"
                                            }
                                            endIcon={<ExternalLinkIcon size={16} />}
                                            valueClassName="cursor-pointer underline"
                                            href={`https://www.etherscan.io/address/${eventContract?.certificate_contract_address}`}
                                        />
                                    </div>
                                </div>

                                {/* Event Certificates Section - Mobile: In Right Column */}
                                <div className="lg:hidden w-full">
                                    <div className="bg-[#1a1a1a] border border-[#333333] rounded-lg p-6 w-full">
                                        <Typography
                                            variant="header"
                                            tag="h2"
                                            className="text-xl font-semibold text-white mb-4"
                                        >
                                            {isCurrentIssuerPending
                                                ? t("issuer.sign.certificatesToSign")
                                                : hasCurrentIssuerSigned
                                                  ? t("issuer.sign.signedCertificates")
                                                  : t("issuer.sign.eventCertificates")}{" "}
                                            (
                                            {eventCertificates?.filter((cert) => !cert.revokedAt)
                                                .length || 0}
                                            )
                                        </Typography>

                                        <div className="w-full overflow-hidden">
                                            <DataTable
                                                columns={certificateColumns}
                                                data={
                                                    eventCertificates
                                                        ?.filter((cert) => {
                                                            // Filter out revoked certificates and ensure required fields exist
                                                            return (
                                                                cert.id &&
                                                                cert.eventId &&
                                                                !cert.revokedAt &&
                                                                cert.createdAt
                                                            );
                                                        })
                                                        .map((cert) => {
                                                            const firstName =
                                                                cert.name?.split(" ")[0] || "";
                                                            const lastName =
                                                                cert.name
                                                                    ?.split(" ")
                                                                    .slice(1)
                                                                    .join(" ") || "";

                                                            return {
                                                                id: cert.id!,
                                                                event_id: cert.eventId!,
                                                                receiver_credential_id:
                                                                    cert.receiverCredentialId,
                                                                receiver_email: cert.receiverEmail,
                                                                name: cert.name,
                                                                academic_institution:
                                                                    cert.academicInstitution,
                                                                certificate_title:
                                                                    cert.certificateTitle,
                                                                certificate_subtitle:
                                                                    cert.certificateSubtitle,
                                                                event_contract_address:
                                                                    cert.eventContractAddress || "",
                                                                event_certificate_address:
                                                                    cert.eventCertificateAddress,
                                                                created_at: cert.createdAt!,
                                                                revoked_at: cert.revokedAt,
                                                                firstName,
                                                                lastName,
                                                                email: cert.receiverEmail || "",
                                                                academicInstitution:
                                                                    cert.academicInstitution || "",
                                                                issuedAt: cert.createdAt!,
                                                                status: cert.revokedAt
                                                                    ? "rejected"
                                                                    : isCurrentIssuerPending
                                                                      ? "pending_signature"
                                                                      : "received",
                                                            };
                                                        }) || []
                                                }
                                                totalItems={
                                                    eventCertificates?.filter(
                                                        (cert) =>
                                                            cert.id &&
                                                            cert.eventId &&
                                                            !cert.revokedAt &&
                                                            cert.createdAt,
                                                    ).length || 0
                                                }
                                                currentPage={1}
                                                pageSize={10}
                                                onPageChange={() => {}}
                                                onPageSizeChange={() => {}}
                                                searchValue=""
                                                onSearchChange={() => {}}
                                                searchPlaceholder="Search certificates..."
                                                sorting={[]}
                                                onSortingChange={() => {}}
                                                isLoading={isLoadingCertificates}
                                                disablePagination
                                            />
                                        </div>
                                    </div>
                                </div>

                                <div className="bg-[#1a1a1a] border border-[#333333] p-6 rounded-lg shadow-lg space-y-6 sticky top-8">
                                    <div>
                                        <Typography
                                            variant="header"
                                            tag="h3"
                                            className="text-xl font-semibold text-white mb-4"
                                        >
                                            {t("issuer.sign.signCertificateBatch")}
                                        </Typography>
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            color="muted-foreground"
                                            className="text-sm mb-4"
                                        >
                                            {t("issuer.sign.signCertificateBatchDescription")}
                                        </Typography>

                                        {/* Sign Button */}
                                        {isCurrentIssuerPending ? (
                                            <Button
                                                onClick={handleSignCertificates}
                                                disabled={
                                                    isSigning || certificatesToSign.length === 0
                                                }
                                                className="w-full bg-[#ff6a39] text-white font-semibold py-3 px-6 rounded-lg mb-3"
                                            >
                                                {isSigning
                                                    ? t("issuer.sign.signing")
                                                    : `${t("issuer.sign.signAndApprove")} (${certificatesToSign.length})`}
                                            </Button>
                                        ) : hasCurrentIssuerSigned ? (
                                            <div className="w-full bg-green-600 text-white font-semibold py-3 px-6 rounded-lg mb-3 flex items-center justify-center">
                                                <CheckCircle className="mr-2 h-5 w-5" />
                                                {t("issuer.sign.alreadySigned")}
                                            </div>
                                        ) : (
                                            <div className="w-full bg-gray-600 text-white font-semibold py-3 px-6 rounded-lg mb-3 flex items-center justify-center">
                                                <Clock className="mr-2 h-5 w-5" />
                                                {t("issuer.sign.notAnIssuer")}
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </main>
                </SectionContainer>

                {/* Password/PIN Modal */}
                <PasswordPinModal
                    isOpen={showPinModal}
                    onClose={handlePinModalClose}
                    onSuccess={(result: { type: "pin" | "password"; value: string }) => {
                        handleSignWithPin(result.value);
                    }}
                    title={t("issuer.sign.signatureRequest")}
                    description={t("issuer.sign.signatureRequestDescription")}
                    showSigningDetails={true}
                    signingDetails={{
                        contractAddress: eventContract?.event_contract_address,
                        transactionType: "Certificate Signing",
                        details: `Signing ${certificatesToSign.length} certificates for ${event?.title}`,
                    }}
                />
            </div>
        </ProtectedRoute>
    );
}
